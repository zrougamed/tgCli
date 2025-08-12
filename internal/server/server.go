package server

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/internal/models"
	"github.com/zrougamed/tgCli/pkg/constants"
)

var versionCommits = map[string]string{
	"3.6.2": "31716aa98a0d4bd3bd7c5488dfc795a82dfee80d",
	"3.6.1": "b77b8fc6c2ceadd457571fc0a6ce1fb243e5f31c",
	"3.6.0": "b77b8fc6c2ceadd457571fc0a6ce1fb243e5f31c",
	"3.5.3": "7edb256d9750ab4451d27eef605e58e9adcedc7a",
	"3.5.0": "375e661f96298db4df018037949827e16ee8df60",
	"3.4.0": "421b0740e4a9f61d6eb0e03a4d2079de625a3ffe",
	"3.3.0": "90cc0512851acca2be10044878815bc414876f23",
	"3.2.2": "a31220261f440f61ebf3edfb1a84efba62939177",
	"3.2.1": "986e09c5d17d303659bed1342506a1c458462e30",
	"3.2.0": "f451d8a9a66c7ca0d8d4d9046f440f21097cdd03",
	"3.1.6": "71b39b25e198f690e28113e7f8874dab7b4559ec",
	"3.1.5": "f91c690f375ecd4d600eb126ca5a920a5c9ad0f4",
	"3.1.2": "3887cbd1d67b58ba6f88c50a069b679e20743984",
	"3.1.1": "375a182bc03b0c78b489e18a0d6af222916a48d2",
	"3.1.0": "e9d3c5d98e7229118309f6d4bbc9446bad7c4c3d",
	"3.0.5": "a9f902e5c552780589a15ba458adb48984359165",
	"3.0.0": "c90ec746a7e77ef5b108554be2133dfd1e1ab1b2",
}

type GSQLSession struct {
	Host     string
	User     string
	Password string
	Version  string
	Cookie   models.GSQLCookie
	Client   *http.Client
}

func RunGSQL(cmd *cobra.Command, args []string) {
	alias, _ := cmd.Flags().GetString("alias")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	gsPort, _ := cmd.Flags().GetString("gsPort")

	// Get configuration if alias is provided
	if alias != "" {
		machineConfig := getMachineConfig(alias)
		if machineConfig != nil {
			host = machineConfig.Host
			user = machineConfig.User
			password = machineConfig.Password
			gsPort = machineConfig.GSPort
		} else {
			fmt.Printf("Alias %s not found. Try: tg conf list\n", alias)
			return
		}
	}

	fullHost := fmt.Sprintf("%s:%s", host, gsPort)

	session := &GSQLSession{
		Host:     fullHost,
		User:     user,
		Password: password,
		Client:   &http.Client{Timeout: 60 * time.Second},
	}

	if err := session.login(); err != nil {
		fmt.Printf("Error logging in to TigerGraph: %v\n", err)
		return
	}

	fmt.Printf("Connected to TigerGraph at %s\n", fullHost)

	// Start interactive GSQL session
	session.startInteractiveSession()
}

func (s *GSQLSession) login() error {
	for version, commit := range versionCommits {
		s.Cookie = models.GSQLCookie{
			ClientCommit:    commit,
			FromGsqlClient:  false,
			FromGraphStudio: false,
			GShellTest:      true,
			FromGsqlServer:  false,
		}

		if err := s.attemptLogin(version); err == nil {
			s.Version = version
			return nil
		}
	}
	return fmt.Errorf("unable to establish compatible connection")
}

func (s *GSQLSession) attemptLogin(version string) error {
	userPass := fmt.Sprintf("%s:%s", s.User, s.Password)
	b64Val := base64.StdEncoding.EncodeToString([]byte(userPass))

	cookieJSON, _ := json.Marshal(s.Cookie)

	req, err := http.NewRequest("POST", s.Host+constants.GSQL_PATH+constants.LOGIN_ENDPOINT, strings.NewReader(b64Val))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Language", "en-US")
	req.Header.Set("Authorization", "Basic "+b64Val)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", string(cookieJSON))
	req.Header.Set("User-Agent", "Java/1.8.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResp struct {
		IsClientCompatible bool   `json:"isClientCompatible"`
		Error              bool   `json:"error"`
		Message            string `json:"message"`
		WelcomeMessage     string `json:"welcomeMessage"`
	}

	if err := json.Unmarshal(body, &loginResp); err != nil {
		return err
	}

	if loginResp.IsClientCompatible {
		if loginResp.Error && s.User != "__GSQL__secret" {
			return fmt.Errorf("%s", loginResp.Message)
		}

		// Update cookies from response
		if setCookie := resp.Header.Get("Set-Cookie"); setCookie != "" {
			var updatedCookie models.GSQLCookie
			if err := json.Unmarshal([]byte(setCookie), &updatedCookie); err == nil {
				s.Cookie = updatedCookie
			}
		}

		fmt.Println(loginResp.WelcomeMessage)
		return nil
	}

	return fmt.Errorf("client not compatible with version %s", version)
}

func (s *GSQLSession) startInteractiveSession() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("GSQL > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}

		command = strings.TrimSpace(command)

		if command == "Quit" || command == "quit" || command == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		if command == "" {
			continue
		}

		if err := s.executeCommand(command); err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}

func (s *GSQLSession) executeCommand(command string) error {
	userPass := fmt.Sprintf("%s:%s", s.User, s.Password)
	b64Val := base64.StdEncoding.EncodeToString([]byte(userPass))

	cookieJSON, _ := json.Marshal(s.Cookie)

	req, err := http.NewRequest("POST", s.Host+constants.GSQL_PATH+constants.FILE_ENDPOINT, strings.NewReader(command))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Language", "en-US")
	req.Header.Set("Authorization", "Basic "+b64Val)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", string(cookieJSON))
	req.Header.Set("User-Agent", "Java/1.8.0")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response in chunks to handle streaming output
	buffer := make([]byte, 1024)
	progressRegex := regexp.MustCompile(`\[.*?\]\s*([0-9]\d*|0)+%.*\(([1-9]\d*|0)\/([1-9]\d*|0)\)`)

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			data := string(buffer[:n])

			if !strings.Contains(data, constants.GSQL_SEPARATOR) {
				// Check for progress bar
				if progressRegex.MatchString(data) {
					fmt.Print(data) // Print progress inline
				} else {
					fmt.Print(strings.TrimSpace(data))
					if !strings.HasSuffix(data, "\n") {
						fmt.Println()
					}
				}
			} else if strings.Contains(data, constants.GSQL_COOKIES) {
				// Update cookies
				parts := strings.Split(data, "__,")
				if len(parts) > 1 {
					var updatedCookie models.GSQLCookie
					if err := json.Unmarshal([]byte(parts[1]), &updatedCookie); err == nil {
						updatedCookie.FromGsqlClient = true
						updatedCookie.FromGraphStudio = false
						updatedCookie.GShellTest = true
						updatedCookie.FromGsqlServer = true
						s.Cookie = updatedCookie
					}
				}
			}
		}

		if err != nil {
			break
		}
	}

	return nil
}

func RunBackup(cmd *cobra.Command, args []string) {
	alias, _ := cmd.Flags().GetString("alias")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	gsPort, _ := cmd.Flags().GetString("gsPort")
	// restPort, _ := cmd.Flags().GetString("restPort")
	backupType, _ := cmd.Flags().GetString("type")

	// Get configuration if alias is provided
	if alias != "" {
		machineConfig := getMachineConfig(alias)
		if machineConfig != nil {
			host = machineConfig.Host
			user = machineConfig.User
			password = machineConfig.Password
			gsPort = machineConfig.GSPort
			// restPort = machineConfig.RestPort
		} else {
			fmt.Printf("Alias %s not found. Try: tg conf list\n", alias)
			return
		}
	}

	optionBKP := ""
	switch backupType {
	case "DATA":
		optionBKP = "-D"
	case "SCHEMA":
		optionBKP = "-S"
	}

	fmt.Printf("Starting backup with type: %s\n", optionBKP)

	// Authenticate and get session
	fullHost := fmt.Sprintf("%s:%s", host, gsPort)
	loginData := map[string]string{
		"username": user,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(fullHost+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error logging in: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Authentication failed with status: %d\n", resp.StatusCode)
		return
	}

	// Get session cookie
	cookie := resp.Header.Get("Set-Cookie")
	if cookie != "" {
		cookie = strings.Split(cookie, ";")[0]
	}

	// Get TigerGraph path
	req, _ := http.NewRequest("GET", fullHost+"/api/log", nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error getting log path: %v\n", err)
		return
	}
	defer resp.Body.Close()

	pathTG := "/home/tigergraph"
	if resp.StatusCode == 200 {
		var logResp struct {
			Error   bool `json:"error"`
			Results []struct {
				Path string `json:"path"`
			} `json:"results"`
		}

		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &logResp); err == nil && !logResp.Error && len(logResp.Results) > 0 {
			parts := strings.Split(logResp.Results[0].Path, "/log/")
			if len(parts) > 0 {
				pathTG = parts[0]
			}
		}
	}

	fmt.Printf("Using TigerGraph path: %s\n", pathTG)
	fmt.Println("Backup functionality requires integration with pyTigerGraph equivalent")
	fmt.Println("This is a placeholder for the full backup implementation")
}

func RunServices(cmd *cobra.Command, args []string) {
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	gsPort, _ := cmd.Flags().GetString("gsPort")
	ops, _ := cmd.Flags().GetString("ops")

	fullHost := fmt.Sprintf("%s:%s", host, gsPort)

	loginData := map[string]string{
		"username": user,
		"password": password,
	}

	jsonData, _ := json.Marshal(loginData)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(fullHost+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error logging in: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Authentication failed with status: %d\n", resp.StatusCode)
		return
	}

	cookie := resp.Header.Get("Set-Cookie")
	if cookie != "" {
		cookie = strings.Split(cookie, ";")[0]
	}

	// Perform service operation
	serviceURL := fmt.Sprintf("%s/api/service/%s?serviceName=gpe&serviceName=gse&serviceName=restpp", fullHost, ops)
	req, _ := http.NewRequest("POST", serviceURL, nil)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error performing service operation: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		var serviceResp struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(body, &serviceResp); err == nil {
			fmt.Println(serviceResp.Message)
		}
	} else {
		fmt.Printf("Service operation failed with status: %d\n", resp.StatusCode)
	}
}

func getMachineConfig(alias string) *models.MachineConfig {
	machines := viper.GetStringMap("machines")
	if machineData, exists := machines[alias]; exists {
		// Convert map[string]interface{} to MachineConfig
		if machineMap, ok := machineData.(map[string]interface{}); ok {
			config := &models.MachineConfig{}
			if host, ok := machineMap["host"].(string); ok {
				config.Host = host
			}
			if user, ok := machineMap["user"].(string); ok {
				config.User = user
			}
			if password, ok := machineMap["password"].(string); ok {
				config.Password = password
			}
			if gsPort, ok := machineMap["gsPort"].(string); ok {
				config.GSPort = gsPort
			}
			if restPort, ok := machineMap["restPort"].(string); ok {
				config.RestPort = restPort
			}
			return config
		}
	}
	return nil
}
