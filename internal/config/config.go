package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/internal/helpers"
	"github.com/zrougamed/tgCli/internal/models"
	"github.com/zrougamed/tgCli/pkg/constants"
	"golang.org/x/term"
)

func RunConfAdd(cmd *cobra.Command, args []string) {
	alias, _ := cmd.Flags().GetString("alias")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	host, _ := cmd.Flags().GetString("host")
	gsPort, _ := cmd.Flags().GetString("gsPort")
	restPort, _ := cmd.Flags().GetString("restPort")
	defaultFlag, _ := cmd.Flags().GetString("default")

	reader := bufio.NewReader(os.Stdin)

	// Get inputs if not provided via flags
	if alias == "" {
		fmt.Print("What is your machine alias? ")
		alias, _ = reader.ReadString('\n')
		alias = strings.TrimSpace(alias)
	}

	if alias == "" {
		fmt.Println("Alias is required")
		return
	}

	// Check if alias already exists
	machines := viper.GetStringMap("machines")
	if _, exists := machines[alias]; exists {
		fmt.Printf("Alias '%s' already exists\n", alias)
		return
	}

	// Get other inputs if not provided
	if host == "http://127.0.0.1" {
		fmt.Print("What is your machine address? (http://127.0.0.1) ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			host = input
		}
	}

	if user == "tigergraph" {
		fmt.Print("What is your machine user? (tigergraph) ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			user = input
		}
	}

	if password == "tigergraph" {
		fmt.Print("What is your machine password? ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err == nil && len(bytePassword) > 0 {
			password = string(bytePassword)
		}
		fmt.Println() // New line after password input
	}

	if gsPort == "14240" {
		fmt.Print("What is your machine gsPort? [14240] ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			gsPort = input
		}
	}

	if restPort == "9000" {
		fmt.Print("What is your machine restPort? [9000] ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input != "" {
			restPort = input
		}
	}

	if defaultFlag == "n" {
		fmt.Print("Would you like to set this machine as default? (y/n) [n] ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "y" || input == "Y" {
			defaultFlag = "y"
		}
	}

	// Save the configuration
	machineConfig := models.MachineConfig{
		Host:     host,
		User:     user,
		Password: password,
		GSPort:   gsPort,
		RestPort: restPort,
	}

	viper.Set(fmt.Sprintf("machines.%s", alias), machineConfig)

	if defaultFlag == "y" {
		viper.Set("default", alias)
		fmt.Printf("Setting up the alias %s as default: success\n", alias)
	}

	if err := helpers.SaveConfig(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Printf("Saving alias %s: success\n", alias)
}

func RunConfDelete(cmd *cobra.Command, args []string) {
	alias, _ := cmd.Flags().GetString("alias")

	if alias == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("What is the machine alias to delete? ")
		alias, _ = reader.ReadString('\n')
		alias = strings.TrimSpace(alias)
	}

	if alias == "" {
		fmt.Println("Alias is required")
		return
	}

	machines := viper.GetStringMap("machines")
	if _, exists := machines[alias]; !exists {
		fmt.Println("Alias not found!")
		return
	}

	// Check if it's the default alias
	defaultAlias := viper.GetString("default")
	if defaultAlias == alias {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("⚠️  You are about to delete the default alias, proceed? (y/n) ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.TrimSpace(strings.ToLower(confirm))

		if confirm != "y" && confirm != "yes" {
			fmt.Println("Aborting...")
			return
		}

		viper.Set("default", "")
	}

	// Delete the machine configuration
	delete(machines, alias)
	viper.Set("machines", machines)

	if err := helpers.SaveConfig(); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		return
	}

	fmt.Println("Alias deleted!")
}

func RunConfList(cmd *cobra.Command, args []string) {
	fmt.Println("======= TGCloud Account ======")

	tgcloudUser := viper.GetString("tgcloud.user")
	tgcloudPassword := viper.GetString("tgcloud.password")

	if tgcloudUser == "mail@domain.com" || tgcloudUser == "" {
		fmt.Println("tgcloud user not set. Use: tg conf tgcloud")
	} else {
		fmt.Printf("tgcloud username: %s\n", tgcloudUser)
		fmt.Printf("tgcloud password: %s\n", maskPassword(tgcloudPassword))
	}

	fmt.Println("======= TigerGraph Instances ======")

	machines := viper.GetStringMap("machines")
	defaultAlias := viper.GetString("default")

	if len(machines) > 0 {
		for alias, machineData := range machines {
			defaultTag := ""
			if defaultAlias == alias {
				defaultTag = " (default)"
			}

			fmt.Printf("Machine: alias = %s%s\n", alias, defaultTag)

			if machineMap, ok := machineData.(map[string]interface{}); ok {
				if host, ok := machineMap["host"].(string); ok {
					fmt.Printf("   host: %s\n", host)
				}
				if user, ok := machineMap["user"].(string); ok {
					fmt.Printf("   user: %s\n", user)
				}
				if password, ok := machineMap["password"].(string); ok {
					fmt.Printf("   password: %s\n", maskPassword(password))
				}
				if gsPort, ok := machineMap["gsPort"].(string); ok {
					fmt.Printf("   GSQL Port: %s\n", gsPort)
				}
				if restPort, ok := machineMap["restPort"].(string); ok {
					fmt.Printf("   REST Port: %s\n", restPort)
				}
			}
			fmt.Println()
		}
	} else {
		fmt.Println("No conf available. Use: tg conf add")
	}
}

func RunConfTGCloud(cmd *cobra.Command, args []string) {
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")

	reader := bufio.NewReader(os.Stdin)

	if email == "" {
		fmt.Print("What is your tgcloud email? ")
		email, _ = reader.ReadString('\n')
		email = strings.TrimSpace(email)
	}

	if password == "" {
		fmt.Print("What is your tgcloud password? ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("Error reading password: %v\n", err)
			return
		}
		password = string(bytePassword)
		fmt.Println() // New line after password input
	}

	if email == "" || password == "" {
		fmt.Println("Email and password are required")
		return
	}

	// Test credentials
	fmt.Println("Trying your credentials...")

	loginData := map[string]string{
		"username": email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("Error marshaling login data: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(constants.TIGERTOOL_URL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error making login request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		var loginResp models.TGCloudResponse
		if err := json.Unmarshal(body, &loginResp); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			return
		}

		if loginResp.Token != "" {
			// Extract bearer token and save
			tokenParts := strings.Split(loginResp.Token, " ")
			if len(tokenParts) >= 2 {
				bearerToken := tokenParts[1]

				if err := os.WriteFile(constants.CredsFile, []byte(bearerToken), 0600); err != nil {
					fmt.Printf("Error saving credentials: %v\n", err)
					return
				}

				// Save credentials to config
				viper.Set("tgcloud.user", email)
				viper.Set("tgcloud.password", password)

				if err := helpers.SaveConfig(); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					return
				}

				fmt.Println("Login Successful! 😊")
				fmt.Println("Credentials saved to configuration")
			}
		}
	} else {
		fmt.Printf("Error logging in: %s\n", string(body))
	}
}

func maskPassword(password string) string {
	if password == "" {
		return ""
	}
	if len(password) <= 3 {
		return strings.Repeat("*", len(password))
	}
	return password[:1] + strings.Repeat("*", len(password)-2) + password[len(password)-1:]
}
