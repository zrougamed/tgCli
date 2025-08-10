package cloud

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func RunLogin(cmd *cobra.Command, args []string) {
	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")
	save, _ := cmd.Flags().GetString("save")
	output, _ := cmd.Flags().GetString("output")

	// Get credentials if not provided
	if email == "" {
		fmt.Print("What is your tgcloud email? ")
		reader := bufio.NewReader(os.Stdin)
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

	// Login request
	loginData := map[string]string{
		"username": email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("Error marshaling login data: %v\n", err)
		return
	}

	fmt.Println("Logging into your account...")

	resp, err := http.Post(constants.TIGERTOOL_URL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error making login request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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
			// Extract bearer token
			tokenParts := strings.Split(loginResp.Token, " ")
			if len(tokenParts) >= 2 {
				bearerToken := tokenParts[1]

				// Save token to file
				if err := ioutil.WriteFile(constants.CredsFile, []byte(bearerToken), 0600); err != nil {
					fmt.Printf("Error saving credentials: %v\n", err)
					return
				}

				// Save credentials to config if requested
				if save == "y" {
					viper.Set("tgcloud.user", email)
					viper.Set("tgcloud.password", password)
					if err := helpers.SaveConfig(); err != nil {
						fmt.Printf("Error saving config: %v\n", err)
					}
				}

				if output == "json" {
					fmt.Printf(`{"error":false,"message":"Login successful","token":"%s"}`, bearerToken)
				} else {
					fmt.Println("Login Successful! 😊")
				}
			}
		}
	} else {
		if output == "json" {
			fmt.Printf(`{"error":true,"message":"Login failed"}`)
		} else {
			fmt.Printf("Error logging in: %s\n", string(body))
		}
	}
}

func RunStart(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	performMachineOperation("start", id)
}

func RunStop(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	performMachineOperation("stop", id)
}

func RunTerminate(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	performMachineOperation("terminate", id)
}

func RunArchive(cmd *cobra.Command, args []string) {
	id, _ := cmd.Flags().GetString("id")
	performMachineOperation("archive", id)
}

func RunList(cmd *cobra.Command, args []string) {
	activeOnly, _ := cmd.Flags().GetString("activeonly")
	output, _ := cmd.Flags().GetString("output")

	bearerToken, err := getBearerToken()
	if err != nil {
		fmt.Printf("Error getting bearer token: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", constants.TGCLOUD_BASE_URL+"/solution", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		var response struct {
			Error  bool             `json:"Error"`
			Result []models.Machine `json:"Result"`
		}

		if err := json.Unmarshal(body, &response); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			return
		}

		if !response.Error {
			var machines []models.Machine
			for _, machine := range response.Result {
				if activeOnly == "y" && machine.State == "terminated" {
					continue
				}
				machines = append(machines, machine)
			}

			if output == "json" {
				result, _ := json.Marshal(map[string]interface{}{
					"error":  false,
					"result": machines,
				})
				fmt.Println(string(result))
			} else {
				printMachineTable("tgcloud solutions", machines)
			}
		}
	} else if resp.StatusCode == 401 {
		if output == "json" {
			fmt.Println(`{"error":true,"message":"Re-Login to tgcloud"}`)
		} else {
			fmt.Println("You should re-login using 'tg cloud login'")
		}
	}
}

func RunCreate(cmd *cobra.Command, args []string) {
	fmt.Println("tgcli Create Machine: 🚧 Work in progress 🚧 will be in next release 🙏 🚀 !")
}

func performMachineOperation(action, machineID string) {
	bearerToken, err := getBearerToken()
	if err != nil {
		fmt.Printf("Error getting bearer token: %v\n", err)
		return
	}

	client := &http.Client{Timeout: 30 * time.Second}

	var req *http.Request
	if action == "terminate" {
		req, err = http.NewRequest("DELETE", constants.TGCLOUD_BASE_URL+"/solution/destroy/"+machineID, nil)
	} else {
		req, err = http.NewRequest("POST", constants.TGCLOUD_BASE_URL+"/solution/"+action+"/"+machineID, nil)
	}

	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	if resp.StatusCode == 200 {
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err == nil {
			if message, ok := response["Message"].(string); ok {
				fmt.Printf("tgcloud response: %s\n", message)
			}
		}
	} else if resp.StatusCode == 401 {
		fmt.Println("tgcloud response: Please re-login")
	} else {
		fmt.Printf("Error: %s\n", string(body))
	}
}

func getBearerToken() (string, error) {
	data, err := ioutil.ReadFile(constants.CredsFile)
	if err != nil {
		return "", fmt.Errorf("bearer token not found, please login first")
	}
	return string(data), nil
}

func printMachineTable(title string, machines []models.Machine) {
	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
	fmt.Printf("%-15s %-20s %-15s %-10s\n", "ID", "Machine", "Solution", "Status")
	fmt.Println(strings.Repeat("-", 65))

	for _, machine := range machines {
		fmt.Printf("%-15s %-20s %-15s %-10s\n",
			machine.ID, machine.Name, machine.Tag, machine.State)
	}
	fmt.Println()
}
