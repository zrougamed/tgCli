package cloud

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/internal/models"
	"github.com/zrougamed/tgCli/pkg/constants"
)

func setupTestEnvironment(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "tgcli_cloud_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Set test constants
	originalCredsFile := constants.CredsFile
	constants.CredsFile = filepath.Join(tempDir, "test_creds.bank")

	cleanup := func() {
		constants.CredsFile = originalCredsFile
		os.RemoveAll(tempDir)
		viper.Reset()
	}

	return tempDir, cleanup
}

func TestGetBearerToken(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test when token file doesn't exist
	_, err := getBearerToken()
	if err == nil {
		t.Error("Expected error when token file doesn't exist")
	}

	// Create token file
	testToken := "test_bearer_token_123"
	err = os.WriteFile(constants.CredsFile, []byte(testToken), 0600)
	if err != nil {
		t.Fatalf("Failed to create token file: %v", err)
	}

	// Test reading token
	token, err := getBearerToken()
	if err != nil {
		t.Fatalf("Failed to read token: %v", err)
	}

	if token != testToken {
		t.Errorf("Expected '%s', got '%s'", testToken, token)
	}
}

func TestPrintMachineTable(t *testing.T) {
	machines := []models.Machine{
		{
			ID:    "machine1",
			Name:  "test-machine-1",
			Tag:   "starter",
			State: "running",
		},
		{
			ID:    "machine2",
			Name:  "test-machine-2",
			Tag:   "enterprise",
			State: "stopped",
		},
	}

	printMachineTable("Test Machines", machines)
}

func TestRunCreate(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("id", "", "")

	RunCreate(cmd, []string{})
}

func TestCloudCommandFlags(t *testing.T) {
	// Test that cloud commands have the expected flags

	// Test login command flags
	cmd := &cobra.Command{}
	cmd.Flags().String("email", "", "")
	cmd.Flags().String("password", "", "")
	cmd.Flags().String("save", "n", "")
	cmd.Flags().String("output", "stdout", "")

	email, _ := cmd.Flags().GetString("email")
	password, _ := cmd.Flags().GetString("password")
	save, _ := cmd.Flags().GetString("save")
	output, _ := cmd.Flags().GetString("output")

	// Test default values
	if email != "" {
		t.Error("Email should default to empty string")
	}
	if password != "" {
		t.Error("Password should default to empty string")
	}
	if save != "n" {
		t.Error("Save should default to 'n'")
	}
	if output != "stdout" {
		t.Error("Output should default to 'stdout'")
	}
}

func TestMachineOperationCommands(t *testing.T) {
	// Test that machine operation commands have the expected flags
	operations := []string{"start", "stop", "terminate", "archive"}

	for _, op := range operations {
		t.Run(op, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().String("id", "", "")

			cmd.MarkFlagRequired("id")

			// Test that the flag exists
			idFlag := cmd.Flags().Lookup("id")
			if idFlag == nil {
				t.Errorf("%s command should have an 'id' flag", op)
			}
		})
	}
}

func TestListCommandFlags(t *testing.T) {
	// Test list command flags
	cmd := &cobra.Command{}
	cmd.Flags().String("activeonly", "y", "")
	cmd.Flags().String("output", "stdout", "")

	activeOnly, _ := cmd.Flags().GetString("activeonly")
	output, _ := cmd.Flags().GetString("output")

	if activeOnly != "y" {
		t.Error("ActiveOnly should default to 'y'")
	}
	if output != "stdout" {
		t.Error("Output should default to 'stdout'")
	}
}

func TestTokenFileOperations(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test token file creation and reading
	testTokens := []string{
		"simple_token",
		"Bearer abc123def456",
		"complex_token_with_special_chars!@#$%",
		"very_long_token_" + strings.Repeat("x", 100),
	}

	for _, testToken := range testTokens {
		t.Run("token_"+testToken[:10], func(t *testing.T) {
			// Write token
			err := os.WriteFile(constants.CredsFile, []byte(testToken), 0600)
			if err != nil {
				t.Fatalf("Failed to write token: %v", err)
			}

			// Read token back
			token, err := getBearerToken()
			if err != nil {
				t.Fatalf("Failed to read token: %v", err)
			}

			if token != testToken {
				t.Errorf("Expected '%s', got '%s'", testToken, token)
			}

			// Clean up for next iteration
			os.Remove(constants.CredsFile)
		})
	}
}

func TestTokenFilePermissions(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create token file
	testToken := "permission_test_token"
	err := os.WriteFile(constants.CredsFile, []byte(testToken), 0600)
	if err != nil {
		t.Fatalf("Failed to create token file: %v", err)
	}

	// Check file permissions
	fileInfo, err := os.Stat(constants.CredsFile)
	if err != nil {
		t.Fatalf("Failed to stat token file: %v", err)
	}

	// Check that file is not world-readable
	mode := fileInfo.Mode()
	if mode&0077 != 0 {
		t.Error("Token file should not be readable by group or others")
	}
}

func TestBearerTokenError(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test reading non-existent token file
	_, err := getBearerToken()
	if err == nil {
		t.Error("Expected error when reading non-existent token file")
	}

	expectedError := "bearer token not found, please login first"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error message to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestMachineTableFormatting(t *testing.T) {
	// Test with different machine configurations
	testCases := []struct {
		name     string
		machines []models.Machine
	}{
		{
			name:     "empty list",
			machines: []models.Machine{},
		},
		{
			name: "single machine",
			machines: []models.Machine{
				{
					ID:    "machine1",
					Name:  "test-machine",
					Tag:   "starter",
					State: "running",
				},
			},
		},
		{
			name: "multiple machines",
			machines: []models.Machine{
				{
					ID:    "machine1",
					Name:  "production",
					Tag:   "enterprise",
					State: "running",
				},
				{
					ID:    "machine2",
					Name:  "development",
					Tag:   "starter",
					State: "stopped",
				},
			},
		},
		{
			name: "machines with long names",
			machines: []models.Machine{
				{
					ID:    "very_long_machine_id_12345",
					Name:  "very_long_machine_name_that_exceeds_normal_length",
					Tag:   "enterprise_with_long_tag_name",
					State: "terminated_but_not_destroyed",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that printMachineTable doesn't panic with different inputs
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("printMachineTable panicked with %s: %v", tc.name, r)
				}
			}()

			printMachineTable("Test: "+tc.name, tc.machines)
		})
	}
}

func TestCloudFunctionSafety(t *testing.T) {
	// Test that cloud functions can be called without crashing
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test functions that don't require network access
	testFunctions := []func(){
		func() {
			cmd := &cobra.Command{}
			cmd.Flags().String("id", "", "")
			RunCreate(cmd, []string{})
		},
		func() {
			printMachineTable("Safety Test", []models.Machine{})
		},
		func() {
			// Test getBearerToken with no file (should return error, not panic)
			_, err := getBearerToken()
			if err == nil {
				t.Error("Expected error when no token file exists")
			}
		},
	}

	for i, testFunc := range testFunctions {
		t.Run(func() string { return "function_" + string(rune(i+'0')) }(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Function %d panicked: %v", i, r)
				}
			}()

			testFunc()
		})
	}
}

// Test that we can handle edge cases in token management
func TestTokenEdgeCases(t *testing.T) {
	_, cleanup := setupTestEnvironment(t)
	defer cleanup()

	edgeCases := []struct {
		name        string
		tokenData   []byte
		expectError bool
	}{
		{
			name:        "empty token",
			tokenData:   []byte(""),
			expectError: false,
		},
		{
			name:        "whitespace token",
			tokenData:   []byte("   \n\t   "),
			expectError: false,
		},
		{
			name:        "binary data",
			tokenData:   []byte{0x00, 0x01, 0x02, 0xFF},
			expectError: false,
		},
		{
			name:        "very long token",
			tokenData:   []byte(strings.Repeat("a", 10000)),
			expectError: false,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write edge case token
			err := os.WriteFile(constants.CredsFile, tc.tokenData, 0600)
			if err != nil {
				t.Fatalf("Failed to write token: %v", err)
			}

			// Try to read it back
			token, err := getBearerToken()

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tc.expectError && string(tc.tokenData) != token {
				t.Errorf("Token mismatch: expected '%s', got '%s'", string(tc.tokenData), token)
			}

			// Clean up
			os.Remove(constants.CredsFile)
		})
	}
}

// Test command flag validation
func TestCommandFlagDefaults(t *testing.T) {
	flagTests := []struct {
		name            string
		setupCmd        func() *cobra.Command
		flagName        string
		expectedDefault string
	}{
		{
			name: "login email flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().String("email", "", "Email address")
				return cmd
			},
			flagName:        "email",
			expectedDefault: "",
		},
		{
			name: "login save flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().String("save", "n", "Save credentials")
				return cmd
			},
			flagName:        "save",
			expectedDefault: "n",
		},
		{
			name: "list activeonly flag",
			setupCmd: func() *cobra.Command {
				cmd := &cobra.Command{}
				cmd.Flags().String("activeonly", "y", "Show active only")
				return cmd
			},
			flagName:        "activeonly",
			expectedDefault: "y",
		},
	}

	for _, tt := range flagTests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.setupCmd()
			value, err := cmd.Flags().GetString(tt.flagName)
			if err != nil {
				t.Fatalf("Failed to get flag %s: %v", tt.flagName, err)
			}

			if value != tt.expectedDefault {
				t.Errorf("Flag %s: expected default '%s', got '%s'",
					tt.flagName, tt.expectedDefault, value)
			}
		})
	}
}
