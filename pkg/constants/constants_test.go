package constants

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"VERSION_CLI", VERSION_CLI, "0.1.1"},
		{"GSQL_PATH", GSQL_PATH, "/gsqlserver/gsql/"},
		{"GSQL_SEPARATOR", GSQL_SEPARATOR, "__GSQL__"},
		{"GSQL_COOKIES", GSQL_COOKIES, "__GSQL__COOKIES__"},
		{"COMMAND_ENDPOINT", COMMAND_ENDPOINT, "command"},
		{"FILE_ENDPOINT", FILE_ENDPOINT, "file"},
		{"LOGIN_ENDPOINT", LOGIN_ENDPOINT, "login"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %s, got %s", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// Test URL constants separately if they're const
func TestURLConstantsValues(t *testing.T) {
	if TGCLOUD_BASE_URL != "https://tgcloud.io/api" {
		t.Errorf("Expected TGCLOUD_BASE_URL to be 'https://tgcloud.io/api', got '%s'", TGCLOUD_BASE_URL)
	}

	if TIGERTOOL_URL != "https://tigertool.tigergraph.com" {
		t.Errorf("Expected TIGERTOOL_URL to be 'https://tigertool.tigergraph.com', got '%s'", TIGERTOOL_URL)
	}
}

func TestURLConstants(t *testing.T) {
	urls := []string{TGCLOUD_BASE_URL, TIGERTOOL_URL}

	for _, url := range urls {
		if !strings.HasPrefix(url, "https://") {
			t.Errorf("URL %s should start with https://", url)
		}
		if strings.HasSuffix(url, "/") {
			t.Errorf("URL %s should not end with /", url)
		}
	}
}

func TestPathConstants(t *testing.T) {
	paths := []string{GSQL_PATH}

	for _, path := range paths {
		if !strings.HasPrefix(path, "/") {
			t.Errorf("Path %s should start with /", path)
		}
		if !strings.HasSuffix(path, "/") {
			t.Errorf("Path %s should end with /", path)
		}
	}
}

func TestGlobalVariables(t *testing.T) {
	// Test that global variables can be set
	originalHomeDir := HomeDir
	originalConfigDir := ConfigDir
	originalConfigFile := ConfigFile
	originalCredsFile := CredsFile
	originalDebug := Debug

	// Set test values
	HomeDir = "/test/home"
	ConfigDir = "/test/config"
	ConfigFile = "/test/config.yml"
	CredsFile = "/test/creds"
	Debug = true

	// Verify values were set
	if HomeDir != "/test/home" {
		t.Error("HomeDir was not set correctly")
	}
	if ConfigDir != "/test/config" {
		t.Error("ConfigDir was not set correctly")
	}
	if ConfigFile != "/test/config.yml" {
		t.Error("ConfigFile was not set correctly")
	}
	if CredsFile != "/test/creds" {
		t.Error("CredsFile was not set correctly")
	}
	if Debug != true {
		t.Error("Debug was not set correctly")
	}

	// Restore original values
	HomeDir = originalHomeDir
	ConfigDir = originalConfigDir
	ConfigFile = originalConfigFile
	CredsFile = originalCredsFile
	Debug = originalDebug
}

func TestDirectoryPaths(t *testing.T) {
	// Simulate typical path construction
	testHomeDir := "/home/testuser"
	testConfigDir := filepath.Join(testHomeDir, ".tgcli")
	testConfigFile := filepath.Join(testConfigDir, "config.yml")
	testCredsFile := filepath.Join(testConfigDir, "creds.bank")

	// Verify path construction
	expectedConfigDir := "/home/testuser/.tgcli"
	if testConfigDir != expectedConfigDir {
		t.Errorf("Expected config dir %s, got %s", expectedConfigDir, testConfigDir)
	}

	expectedConfigFile := "/home/testuser/.tgcli/config.yml"
	if testConfigFile != expectedConfigFile {
		t.Errorf("Expected config file %s, got %s", expectedConfigFile, testConfigFile)
	}

	expectedCredsFile := "/home/testuser/.tgcli/creds.bank"
	if testCredsFile != expectedCredsFile {
		t.Errorf("Expected creds file %s, got %s", expectedCredsFile, testCredsFile)
	}
}

func TestEnvironmentVariables(t *testing.T) {
	// Test environment variable handling (if any)
	originalDebug := Debug

	// Test debug mode
	os.Setenv("TGCLI_DEBUG", "true")

	// Clean up
	os.Unsetenv("TGCLI_DEBUG")
	Debug = originalDebug
}

func TestVersionFormat(t *testing.T) {
	// Test that version follows semantic versioning pattern
	if !strings.Contains(VERSION_CLI, ".") {
		t.Error("Version should contain dots for semantic versioning")
	}

	parts := strings.Split(VERSION_CLI, ".")
	if len(parts) != 3 {
		t.Errorf("Version should have 3 parts, got %d: %s", len(parts), VERSION_CLI)
	}
}

func TestEndpointConstants(t *testing.T) {
	endpoints := []string{COMMAND_ENDPOINT, FILE_ENDPOINT, LOGIN_ENDPOINT}

	for _, endpoint := range endpoints {
		if endpoint == "" {
			t.Error("Endpoint should not be empty")
		}
		if strings.Contains(endpoint, "/") {
			t.Errorf("Endpoint %s should not contain slashes", endpoint)
		}
		if strings.Contains(endpoint, " ") {
			t.Errorf("Endpoint %s should not contain spaces", endpoint)
		}
	}
}

func TestSeparatorConstants(t *testing.T) {
	separators := []string{GSQL_SEPARATOR, GSQL_COOKIES}

	for _, sep := range separators {
		if sep == "" {
			t.Error("Separator should not be empty")
		}
		if !strings.Contains(sep, "__") {
			t.Errorf("Separator %s should contain double underscores", sep)
		}
	}
}
