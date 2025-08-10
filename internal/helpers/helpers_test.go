package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestCreateDefaultConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tgcli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test_config.yml")

	// Reset viper before test
	viper.Reset()

	// Test creating default config
	err = CreateDefaultConfig(configFile)
	if err != nil {
		t.Fatalf("CreateDefaultConfig failed: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// Read and verify config content by reading the file and checking viper state
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	// Verify default values
	user := viper.GetString("tgcloud.user")
	if user != "mail@domain.com" {
		t.Errorf("Expected default user 'mail@domain.com', got '%s'", user)
	}

	password := viper.GetString("tgcloud.password")
	if password != "" {
		t.Errorf("Expected empty password, got '%s'", password)
	}

	defaultAlias := viper.GetString("default")
	if defaultAlias != "" {
		t.Errorf("Expected empty default alias, got '%s'", defaultAlias)
	}

	machines := viper.GetStringMap("machines")
	if len(machines) != 0 {
		t.Errorf("Expected empty machines map, got %d entries", len(machines))
	}
}

func TestCreateDefaultConfigInvalidPath(t *testing.T) {
	// Test with invalid path (non-existent directory)
	invalidPath := "/nonexistent/directory/config.yml"

	err := CreateDefaultConfig(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}

	// Verify error contains expected message
	if !strings.Contains(err.Error(), "unable to create default config file") {
		t.Errorf("Expected error message about config file, got: %v", err)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "tgcli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "test_config.yml")

	// Setup viper with test config
	viper.SetConfigFile(configFile)
	viper.Set("test.key", "test_value")
	viper.Set("test.number", 42)

	// Test saving config
	err = SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file was not saved")
	}

	// Reset viper and read config to verify save worked
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	// Verify saved values
	testKey := viper.GetString("test.key")
	if testKey != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", testKey)
	}

	testNumber := viper.GetInt("test.number")
	if testNumber != 42 {
		t.Errorf("Expected 42, got %d", testNumber)
	}
}

func TestSaveConfigNoConfigSet(t *testing.T) {
	// Reset viper to ensure no config is set
	viper.Reset()

	// This should return an error since no config file is set
	err := SaveConfig()
	if err == nil {
		t.Error("Expected error when no config file is set, got nil")
	}
}

func TestGracefulShutdown(t *testing.T) {
	// This test is tricky because GracefulShutdown sets up signal handlers
	// We can test that it doesn't panic and sets up the handlers

	// Save original args to restore later
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test that GracefulShutdown doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GracefulShutdown panicked: %v", r)
		}
	}()

	GracefulShutdown()

	// Give it a moment to set up handlers
	time.Sleep(10 * time.Millisecond)

}

func TestGracefulShutdownSignalHandling(t *testing.T) {

	if testing.Short() {
		t.Skip("Skipping signal test in short mode")
	}

	// Test multiple calls don't cause issues
	GracefulShutdown()
	GracefulShutdown() // Should be safe to call multiple times

	// Give it a moment to set up handlers
	time.Sleep(10 * time.Millisecond)

	t.Log("GracefulShutdown signal handler setup completed without issues")
}

func TestCheckForUpdates(t *testing.T) {
	// Test the placeholder implementation
	version, err := CheckForUpdates()

	if err != nil {
		t.Errorf("CheckForUpdates returned error: %v", err)
	}

	if version != "N/A" {
		t.Errorf("Expected 'N/A', got '%s'", version)
	}
}

func TestCheckForUpdatesConsistency(t *testing.T) {
	// Test that multiple calls return consistent results
	version1, err1 := CheckForUpdates()
	version2, err2 := CheckForUpdates()

	if err1 != err2 {
		t.Error("CheckForUpdates should return consistent errors")
	}

	if version1 != version2 {
		t.Error("CheckForUpdates should return consistent versions")
	}
}

// Helper function to test viper configuration
func setupTestViper(t *testing.T) (string, func()) {
	tempDir, err := os.MkdirTemp("", "tgcli_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	configFile := filepath.Join(tempDir, "test_config.yml")

	// Save original viper state
	originalConfigFile := viper.ConfigFileUsed()
	originalSettings := viper.AllSettings()

	// Setup test viper
	viper.Reset()
	viper.SetConfigFile(configFile)

	cleanup := func() {
		viper.Reset()
		if originalConfigFile != "" {
			viper.SetConfigFile(originalConfigFile)
			for key, value := range originalSettings {
				viper.Set(key, value)
			}
		}
		os.RemoveAll(tempDir)
	}

	return configFile, cleanup
}

func TestCreateDefaultConfigWithViperState(t *testing.T) {
	configFile, cleanup := setupTestViper(t)
	defer cleanup()

	// Set some existing viper values
	viper.Set("existing.key", "existing_value")

	err := CreateDefaultConfig(configFile)
	if err != nil {
		t.Fatalf("CreateDefaultConfig failed: %v", err)
	}

	// Read the created config file to verify defaults were written
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read created config: %v", err)
	}

	// Verify that default values are set
	user := viper.GetString("tgcloud.user")
	if user != "mail@domain.com" {
		t.Errorf("Expected default user 'mail@domain.com', got '%s'", user)
	}

	// Verify existing values are preserved (if the function preserves them)
	existingKey := viper.GetString("existing.key")
	if existingKey != "" {
		t.Log("Note: CreateDefaultConfig overwrites existing viper values")
	}
}

func TestSaveConfigWithComplexData(t *testing.T) {
	configFile, cleanup := setupTestViper(t)
	defer cleanup()

	// Set complex nested data
	viper.Set("tgcloud.user", "test@example.com")
	viper.Set("tgcloud.password", "testpass")
	viper.Set("machines.prod.host", "https://prod.example.com")
	viper.Set("machines.prod.user", "admin")
	viper.Set("machines.dev.host", "http://localhost")
	viper.Set("machines.dev.user", "tigergraph")
	viper.Set("default", "prod")

	err := SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Reset and reload to verify
	viper.Reset()
	viper.SetConfigFile(configFile)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read saved config: %v", err)
	}

	// Verify complex data was saved correctly
	if viper.GetString("tgcloud.user") != "test@example.com" {
		t.Error("TGCloud user not saved correctly")
	}

	if viper.GetString("machines.prod.host") != "https://prod.example.com" {
		t.Error("Machine host not saved correctly")
	}

	if viper.GetString("default") != "prod" {
		t.Error("Default value not saved correctly")
	}
}

func TestErrorHandling(t *testing.T) {
	// Test error handling in CreateDefaultConfig with read-only directory
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Create a read-only directory
	tempDir, err := os.MkdirTemp("", "tgcli_readonly_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Make directory read-only
	err = os.Chmod(tempDir, 0444)
	if err != nil {
		t.Fatalf("Failed to make directory read-only: %v", err)
	}

	configFile := filepath.Join(tempDir, "config.yml")

	err = CreateDefaultConfig(configFile)
	if err == nil {
		t.Error("Expected error when writing to read-only directory")
	}
}

func TestGracefulShutdownMultipleCalls(t *testing.T) {
	// Test that calling GracefulShutdown multiple times is safe
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Multiple GracefulShutdown calls caused panic: %v", r)
		}
	}()

	// Should be safe to call multiple times
	GracefulShutdown()
	GracefulShutdown()
	GracefulShutdown()

	// Give handlers time to set up
	time.Sleep(10 * time.Millisecond)
}

func TestHelperFunctionsSafety(t *testing.T) {
	// Test that helper functions can be called without side effects

	// Test CheckForUpdates multiple times
	for i := 0; i < 3; i++ {
		version, err := CheckForUpdates()
		if err != nil {
			t.Errorf("CheckForUpdates call %d failed: %v", i+1, err)
		}
		if version != "N/A" {
			t.Errorf("CheckForUpdates call %d returned unexpected version: %s", i+1, version)
		}
	}
}
