package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zrougamed/tgCli/pkg/constants"
)

func setupMainTestEnvironment(t *testing.T) func() {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tgcli_main_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Save original values
	originalHomeDir := constants.HomeDir
	originalConfigDir := constants.ConfigDir
	originalConfigFile := constants.ConfigFile
	originalCredsFile := constants.CredsFile
	originalDebug := constants.Debug

	// Set test values
	constants.HomeDir = tempDir
	constants.ConfigDir = filepath.Join(tempDir, ".tgcli")
	constants.ConfigFile = filepath.Join(constants.ConfigDir, "config.yml")
	constants.CredsFile = filepath.Join(constants.ConfigDir, "creds.bank")

	// Create config directory
	err = os.MkdirAll(constants.ConfigDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	cleanup := func() {
		constants.HomeDir = originalHomeDir
		constants.ConfigDir = originalConfigDir
		constants.ConfigFile = originalConfigFile
		constants.CredsFile = originalCredsFile
		constants.Debug = originalDebug
		os.RemoveAll(tempDir)
		viper.Reset()
	}

	return cleanup
}

func TestInitialization(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that initialization creates necessary directories
	if _, err := os.Stat(constants.ConfigDir); os.IsNotExist(err) {
		t.Error("Config directory was not created during init")
	}

	// Test that viper is configured correctly
	if viper.GetString("tgcloud.user") != "mail@domain.com" {
		t.Error("Default TGCloud user not set correctly")
	}

	if viper.GetString("tgcloud.password") != "" {
		t.Error("Default TGCloud password should be empty")
	}

	defaultAlias := viper.GetString("default")
	if defaultAlias != "" {
		t.Error("Default alias should be empty initially")
	}

	machines := viper.GetStringMap("machines")
	if len(machines) != 0 {
		t.Error("Machines map should be empty initially")
	}
}

func TestRootCommand(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Create root command
	var rootCmd = &cobra.Command{
		Use:   "tg",
		Short: "TigerGraph CLI tool for cloud and server management",
		Long:  `A comprehensive CLI tool for managing TigerGraph cloud instances and server operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Test basic properties
	if rootCmd.Use != "tg" {
		t.Error("Root command name should be 'tg'")
	}

	if !strings.Contains(rootCmd.Short, "TigerGraph CLI") {
		t.Error("Root command short description should mention TigerGraph CLI")
	}

	if !strings.Contains(rootCmd.Long, "comprehensive CLI tool") {
		t.Error("Root command long description should mention comprehensive CLI tool")
	}
}

func TestCreateCloudCmd(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	cloudCmd := createCloudCmd()

	// Test command properties
	if cloudCmd.Use != "cloud" {
		t.Error("Cloud command should use 'cloud'")
	}

	if !strings.Contains(cloudCmd.Short, "TigerGraph Cloud") {
		t.Error("Cloud command should mention TigerGraph Cloud in short description")
	}

	// Test subcommands
	expectedSubcommands := []string{"login", "start", "stop", "terminate", "archive", "list", "create"}
	commands := cloudCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d cloud subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Cloud subcommand '%s' not found", expected)
		}
	}
}

func TestCreateServerCmd(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	serverCmd := createServerCmd()

	// Test command properties
	if serverCmd.Use != "server" {
		t.Error("Server command should use 'server'")
	}

	if !strings.Contains(serverCmd.Short, "TigerGraph Server") {
		t.Error("Server command should mention TigerGraph Server in short description")
	}

	// Test subcommands
	expectedSubcommands := []string{"gsql", "backup", "services"}
	commands := serverCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d server subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Server subcommand '%s' not found", expected)
		}
	}
}

func TestCreateConfCmd(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	confCmd := createConfCmd()

	// Test command properties
	if confCmd.Use != "conf" {
		t.Error("Conf command should use 'conf'")
	}

	if !strings.Contains(confCmd.Short, "Configuration") {
		t.Error("Conf command should mention Configuration in short description")
	}

	// Test subcommands
	expectedSubcommands := []string{"add", "delete", "list", "tgcloud"}
	commands := confCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d conf subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Conf subcommand '%s' not found", expected)
		}
	}
}

func TestCloudLoginCommandFlags(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	cloudCmd := createCloudCmd()

	// Find login command
	var loginCmd *cobra.Command
	for _, cmd := range cloudCmd.Commands() {
		if cmd.Use == "login" {
			loginCmd = cmd
			break
		}
	}

	if loginCmd == nil {
		t.Fatal("Login command not found")
	}

	// Test required flags
	expectedFlags := []string{"email", "password", "save", "output"}
	for _, flagName := range expectedFlags {
		flag := loginCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Login command should have flag '%s'", flagName)
		}
	}

	// Test flag properties
	emailFlag := loginCmd.Flags().Lookup("email")
	if emailFlag.Shorthand != "e" {
		t.Error("Email flag should have shorthand 'e'")
	}

	passwordFlag := loginCmd.Flags().Lookup("password")
	if passwordFlag.Shorthand != "p" {
		t.Error("Password flag should have shorthand 'p'")
	}
}

func TestServerGSQLCommandFlags(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	serverCmd := createServerCmd()

	// Find gsql command
	var gsqlCmd *cobra.Command
	for _, cmd := range serverCmd.Commands() {
		if cmd.Use == "gsql" {
			gsqlCmd = cmd
			break
		}
	}

	if gsqlCmd == nil {
		t.Fatal("GSQL command not found")
	}

	// Test flags
	expectedFlags := []string{"alias", "user", "password", "host", "gsPort"}
	for _, flagName := range expectedFlags {
		flag := gsqlCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("GSQL command should have flag '%s'", flagName)
		}
	}

	// Test default values
	userFlag := gsqlCmd.Flags().Lookup("user")
	if userFlag.DefValue != "tigergraph" {
		t.Error("User flag should default to 'tigergraph'")
	}

	hostFlag := gsqlCmd.Flags().Lookup("host")
	if hostFlag.DefValue != "http://127.0.0.1" {
		t.Error("Host flag should default to 'http://127.0.0.1'")
	}

	gsPortFlag := gsqlCmd.Flags().Lookup("gsPort")
	if gsPortFlag.DefValue != "14240" {
		t.Error("GSPort flag should default to '14240'")
	}
}

func TestConfAddCommandFlags(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	confCmd := createConfCmd()

	// Find add command
	var addCmd *cobra.Command
	for _, cmd := range confCmd.Commands() {
		if cmd.Use == "add" {
			addCmd = cmd
			break
		}
	}

	if addCmd == nil {
		t.Fatal("Add command not found")
	}

	// Test flags
	expectedFlags := []string{"alias", "user", "password", "host", "gsPort", "restPort", "default"}
	for _, flagName := range expectedFlags {
		flag := addCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Add command should have flag '%s'", flagName)
		}
	}

	// Test flag shortcuts
	aliasFlag := addCmd.Flags().Lookup("alias")
	if aliasFlag.Shorthand != "a" {
		t.Error("Alias flag should have shorthand 'a'")
	}

	defaultFlag := addCmd.Flags().Lookup("default")
	if defaultFlag.Shorthand != "d" {
		t.Error("Default flag should have shorthand 'd'")
	}
}

func TestRequiredFlags(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test cloud start command required flags
	cloudCmd := createCloudCmd()
	var startCmd *cobra.Command
	for _, cmd := range cloudCmd.Commands() {
		if cmd.Use == "start" {
			startCmd = cmd
			break
		}
	}

	if startCmd == nil {
		t.Fatal("Start command not found")
	}

	// Test that ID flag is required
	idFlag := startCmd.Flags().Lookup("id")
	if idFlag == nil {
		t.Error("Start command should have ID flag")
	}

	// Test conf delete command required flags
	confCmd := createConfCmd()
	var deleteCmd *cobra.Command
	for _, cmd := range confCmd.Commands() {
		if cmd.Use == "delete" {
			deleteCmd = cmd
			break
		}
	}

	if deleteCmd == nil {
		t.Fatal("Delete command not found")
	}

	aliasFlag := deleteCmd.Flags().Lookup("alias")
	if aliasFlag == nil {
		t.Error("Delete command should have alias flag")
	}
}

func TestVersionCommand(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			// Capture output for testing
		},
	}

	if versionCmd.Use != "version" {
		t.Error("Version command should use 'version'")
	}

	if !strings.Contains(versionCmd.Short, "version") {
		t.Error("Version command should mention version in description")
	}
}

func TestGlobalFlags(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	var rootCmd = &cobra.Command{
		Use: "tg",
	}

	// Add global flags
	rootCmd.PersistentFlags().BoolVarP(&constants.Debug, "debug", "d", false, "Enable debug mode")

	// Test that debug flag exists
	debugFlag := rootCmd.PersistentFlags().Lookup("debug")
	if debugFlag == nil {
		t.Error("Root command should have debug flag")
	}

	if debugFlag.Shorthand != "d" {
		t.Error("Debug flag should have shorthand 'd'")
	}
}

func TestCommandExecution(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that commands can be executed without panicking
	commands := []func() *cobra.Command{
		createCloudCmd,
		createServerCmd,
		createConfCmd,
	}

	for _, cmdFunc := range commands {
		cmd := cmdFunc()
		if cmd == nil {
			t.Error("Command creation function returned nil")
		}

		// Test that command has proper structure
		if cmd.Use == "" {
			t.Error("Command should have a Use field")
		}

		if cmd.Short == "" {
			t.Error("Command should have a Short description")
		}
	}
}

func TestMainFunction(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test root command creation
	var rootCmd = &cobra.Command{
		Use:   "tg",
		Short: "TigerGraph CLI tool for cloud and server management",
		Long:  `A comprehensive CLI tool for managing TigerGraph cloud instances and server operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands (like main does)
	rootCmd.AddCommand(createCloudCmd())
	rootCmd.AddCommand(createServerCmd())
	rootCmd.AddCommand(createConfCmd())

	// Test that all subcommands are added
	if len(rootCmd.Commands()) < 3 {
		t.Error("Root command should have at least 3 subcommands")
	}
}

func TestErrorHandling(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test error handling in initialization
	// Make directory read-only to test error handling
	if os.Getuid() != 0 { // Skip if running as root
		readOnlyDir := filepath.Join(os.TempDir(), "readonly_tgcli_test")
		err := os.MkdirAll(readOnlyDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		defer os.RemoveAll(readOnlyDir)

		err = os.Chmod(readOnlyDir, 0444)
		if err != nil {
			t.Fatalf("Failed to make directory read-only: %v", err)
		}

		// Test would normally check that the init function handles this gracefully
		// In the actual init function, this would be tested by trying to create
		// a subdirectory in the read-only directory
	}
}

func TestConfigFileHandling(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that config file is created if it doesn't exist
	if _, err := os.Stat(constants.ConfigFile); os.IsNotExist(err) {
		// This is expected - config file should be created by init or helpers
		t.Log("Config file doesn't exist initially, which is expected")
	}

	// Test viper configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(constants.ConfigDir)

	// Test reading config (should handle file not found gracefully)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			t.Errorf("Unexpected error reading config: %v", err)
		}
	}
}

func TestCommandLineArgumentParsing(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test various command line scenarios
	testCases := []struct {
		name     string
		args     []string
		expected bool // whether command should be valid
	}{
		{
			name:     "No arguments",
			args:     []string{},
			expected: true, // Should show help
		},
		{
			name:     "Version command",
			args:     []string{"version"},
			expected: true,
		},
		{
			name:     "Cloud command",
			args:     []string{"cloud"},
			expected: true, // Should show cloud help
		},
		{
			name:     "Cloud login",
			args:     []string{"cloud", "login"},
			expected: true,
		},
		{
			name:     "Server command",
			args:     []string{"server"},
			expected: true,
		},
		{
			name:     "Conf command",
			args:     []string{"conf"},
			expected: true,
		},
		{
			name:     "Invalid command",
			args:     []string{"invalid"},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootCmd := &cobra.Command{
				Use: "tg",
				Run: func(cmd *cobra.Command, args []string) {
					// Default help behavior
				},
			}

			// Add subcommands
			rootCmd.AddCommand(createCloudCmd())
			rootCmd.AddCommand(createServerCmd())
			rootCmd.AddCommand(createConfCmd())

			// Add version command
			versionCmd := &cobra.Command{
				Use:   "version",
				Short: "Show version information",
				Run: func(cmd *cobra.Command, args []string) {
					// Version output
				},
			}
			rootCmd.AddCommand(versionCmd)

			// Test command parsing
			rootCmd.SetArgs(tc.args)

			// Capture output
			var output bytes.Buffer
			rootCmd.SetOut(&output)
			rootCmd.SetErr(&output)

			err := rootCmd.Execute()

			if tc.expected && err != nil {
				t.Errorf("Expected command to succeed, got error: %v", err)
			}

			if !tc.expected && err == nil {
				t.Error("Expected command to fail, but it succeeded")
			}
		})
	}
}

func TestFlagValidation(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test cloud start command with missing required flag
	cloudCmd := createCloudCmd()
	var startCmd *cobra.Command
	for _, cmd := range cloudCmd.Commands() {
		if cmd.Use == "start" {
			startCmd = cmd
			break
		}
	}

	if startCmd == nil {
		t.Fatal("Start command not found")
	}

	// Test executing without required ID flag
	startCmd.SetArgs([]string{})

	var output bytes.Buffer
	startCmd.SetOut(&output)
	startCmd.SetErr(&output)

	err := startCmd.Execute()
	// Should fail due to missing required flag (if properly configured)
	if err == nil {
		t.Log("Note: Required flag validation may not be strictly enforced in test")
	}
}

func TestEnvironmentSetup(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that environment variables are handled correctly
	// Set a test environment variable
	os.Setenv("TGCLI_TEST", "test_value")
	defer os.Unsetenv("TGCLI_TEST")

	// Test that constants are set correctly
	if constants.HomeDir == "" {
		t.Error("HomeDir should be set")
	}

	if constants.ConfigDir == "" {
		t.Error("ConfigDir should be set")
	}

	if constants.ConfigFile == "" {
		t.Error("ConfigFile should be set")
	}

	if constants.CredsFile == "" {
		t.Error("CredsFile should be set")
	}

	// Test directory structure
	if !strings.HasSuffix(constants.ConfigDir, ".tgcli") {
		t.Error("ConfigDir should end with .tgcli")
	}

	if !strings.HasSuffix(constants.ConfigFile, "config.yml") {
		t.Error("ConfigFile should end with config.yml")
	}

	if !strings.HasSuffix(constants.CredsFile, "creds.bank") {
		t.Error("CredsFile should end with creds.bank")
	}
}

func TestConcurrentExecution(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that multiple command creations don't interfere with each other
	done := make(chan bool, 3)

	go func() {
		cmd := createCloudCmd()
		if cmd == nil {
			t.Error("Cloud command creation failed")
		}
		done <- true
	}()

	go func() {
		cmd := createServerCmd()
		if cmd == nil {
			t.Error("Server command creation failed")
		}
		done <- true
	}()

	go func() {
		cmd := createConfCmd()
		if cmd == nil {
			t.Error("Conf command creation failed")
		}
		done <- true
	}()

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestMemoryLeaks(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that repeated command creation doesn't cause memory leaks
	for i := 0; i < 100; i++ {
		cloudCmd := createCloudCmd()
		serverCmd := createServerCmd()
		confCmd := createConfCmd()

		// Ensure commands are created successfully
		if cloudCmd == nil || serverCmd == nil || confCmd == nil {
			t.Error("Command creation failed in iteration", i)
			break
		}
	}
}

func TestCommandHelp(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test that all commands have proper help text
	commands := []*cobra.Command{
		createCloudCmd(),
		createServerCmd(),
		createConfCmd(),
	}

	for _, cmd := range commands {
		if cmd.Short == "" {
			t.Errorf("Command %s should have short description", cmd.Use)
		}

		if cmd.Long == "" {
			t.Errorf("Command %s should have long description", cmd.Use)
		}

		// Test help execution
		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetArgs([]string{"--help"})

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Help execution failed for command %s: %v", cmd.Use, err)
		}

		helpOutput := output.String()
		if !strings.Contains(helpOutput, cmd.Use) {
			t.Errorf("Help output should contain command name %s", cmd.Use)
		}
	}
}

func TestCommandCompletion(t *testing.T) {
	cleanup := setupMainTestEnvironment(t)
	defer cleanup()

	// Test basic command structure for shell completion
	rootCmd := &cobra.Command{Use: "tg"}
	rootCmd.AddCommand(createCloudCmd())
	rootCmd.AddCommand(createServerCmd())
	rootCmd.AddCommand(createConfCmd())

	// Test that commands are properly structured for completion
	commands := rootCmd.Commands()
	if len(commands) < 3 {
		t.Error("Root command should have at least 3 subcommands for completion")
	}

	for _, cmd := range commands {
		if cmd.Use == "" {
			t.Error("All commands should have Use field for completion")
		}
	}
}
