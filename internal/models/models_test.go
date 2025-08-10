package models

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	config := Config{
		TGCloud: TGCloudConfig{
			User:     "test@example.com",
			Password: "testpass",
		},
		Machines: map[string]MachineConfig{
			"test": {
				Host:     "http://localhost",
				User:     "tigergraph",
				Password: "tigergraph",
				GSPort:   "14240",
				RestPort: "9000",
			},
		},
		Default: "test",
	}

	// Test that config can be created
	if config.TGCloud.User != "test@example.com" {
		t.Error("TGCloud user not set correctly")
	}

	if config.Default != "test" {
		t.Error("Default not set correctly")
	}

	if len(config.Machines) != 1 {
		t.Error("Machines map should have 1 entry")
	}
}

func TestTGCloudConfig(t *testing.T) {
	tgConfig := TGCloudConfig{
		User:     "user@domain.com",
		Password: "securepass",
	}

	if tgConfig.User != "user@domain.com" {
		t.Error("User not set correctly")
	}

	if tgConfig.Password != "securepass" {
		t.Error("Password not set correctly")
	}
}

func TestMachineConfig(t *testing.T) {
	machine := MachineConfig{
		Host:     "https://cluster.tgcloud.io",
		User:     "admin",
		Password: "adminpass",
		GSPort:   "14240",
		RestPort: "9000",
	}

	if machine.Host != "https://cluster.tgcloud.io" {
		t.Error("Host not set correctly")
	}

	if machine.User != "admin" {
		t.Error("User not set correctly")
	}

	if machine.GSPort != "14240" {
		t.Error("GSPort not set correctly")
	}

	if machine.RestPort != "9000" {
		t.Error("RestPort not set correctly")
	}
}

func TestGSQLCookie(t *testing.T) {
	cookie := GSQLCookie{
		ClientCommit:    "abc123",
		FromGsqlClient:  true,
		FromGraphStudio: false,
		GShellTest:      true,
		FromGsqlServer:  false,
	}

	if cookie.ClientCommit != "abc123" {
		t.Error("ClientCommit not set correctly")
	}

	if !cookie.FromGsqlClient {
		t.Error("FromGsqlClient should be true")
	}

	if cookie.FromGraphStudio {
		t.Error("FromGraphStudio should be false")
	}

	if !cookie.GShellTest {
		t.Error("GShellTest should be true")
	}

	if cookie.FromGsqlServer {
		t.Error("FromGsqlServer should be false")
	}
}

func TestGSQLCookieJSON(t *testing.T) {
	cookie := GSQLCookie{
		ClientCommit:    "test123",
		FromGsqlClient:  true,
		FromGraphStudio: false,
		GShellTest:      true,
		FromGsqlServer:  true,
	}

	// Test JSON marshaling
	data, err := json.Marshal(cookie)
	if err != nil {
		t.Fatalf("Failed to marshal GSQLCookie: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled GSQLCookie
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal GSQLCookie: %v", err)
	}

	// Verify data integrity
	if unmarshaled.ClientCommit != cookie.ClientCommit {
		t.Error("ClientCommit not preserved during JSON round-trip")
	}

	if unmarshaled.FromGsqlClient != cookie.FromGsqlClient {
		t.Error("FromGsqlClient not preserved during JSON round-trip")
	}

	if unmarshaled.FromGraphStudio != cookie.FromGraphStudio {
		t.Error("FromGraphStudio not preserved during JSON round-trip")
	}

	if unmarshaled.GShellTest != cookie.GShellTest {
		t.Error("GShellTest not preserved during JSON round-trip")
	}

	if unmarshaled.FromGsqlServer != cookie.FromGsqlServer {
		t.Error("FromGsqlServer not preserved during JSON round-trip")
	}
}

func TestTGCloudResponse(t *testing.T) {
	response := TGCloudResponse{
		Error:   false,
		Message: "Success",
		Result:  map[string]string{"status": "ok"},
		Token:   "Bearer abc123",
	}

	if response.Error {
		t.Error("Error should be false")
	}

	if response.Message != "Success" {
		t.Error("Message not set correctly")
	}

	if response.Token != "Bearer abc123" {
		t.Error("Token not set correctly")
	}

	// Test with different result types
	response.Result = []string{"item1", "item2"}
	if result, ok := response.Result.([]string); !ok || len(result) != 2 {
		t.Error("Result should accept different types")
	}
}

func TestTGCloudResponseJSON(t *testing.T) {
	response := TGCloudResponse{
		Error:   false,
		Message: "Login successful",
		Result:  "user_data",
		Token:   "Bearer token123",
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal TGCloudResponse: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled TGCloudResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal TGCloudResponse: %v", err)
	}

	// Verify data integrity
	if unmarshaled.Error != response.Error {
		t.Error("Error not preserved during JSON round-trip")
	}

	if unmarshaled.Message != response.Message {
		t.Error("Message not preserved during JSON round-trip")
	}

	if unmarshaled.Token != response.Token {
		t.Error("Token not preserved during JSON round-trip")
	}
}

func TestMachine(t *testing.T) {
	machine := Machine{
		ID:        "machine123",
		Name:      "test-cluster",
		Tag:       "starter",
		State:     "running",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	if machine.ID != "machine123" {
		t.Error("ID not set correctly")
	}

	if machine.Name != "test-cluster" {
		t.Error("Name not set correctly")
	}

	if machine.Tag != "starter" {
		t.Error("Tag not set correctly")
	}

	if machine.State != "running" {
		t.Error("State not set correctly")
	}

	if machine.CreatedAt != "2024-01-01T00:00:00Z" {
		t.Error("CreatedAt not set correctly")
	}
}

func TestMachineJSON(t *testing.T) {
	machine := Machine{
		ID:        "test123",
		Name:      "production",
		Tag:       "enterprise",
		State:     "stopped",
		CreatedAt: "2024-01-15T10:30:00Z",
	}

	// Test JSON marshaling
	data, err := json.Marshal(machine)
	if err != nil {
		t.Fatalf("Failed to marshal Machine: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled Machine
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal Machine: %v", err)
	}

	// Verify data integrity
	if unmarshaled.ID != machine.ID {
		t.Error("ID not preserved during JSON round-trip")
	}

	if unmarshaled.Name != machine.Name {
		t.Error("Name not preserved during JSON round-trip")
	}

	if unmarshaled.Tag != machine.Tag {
		t.Error("Tag not preserved during JSON round-trip")
	}

	if unmarshaled.State != machine.State {
		t.Error("State not preserved during JSON round-trip")
	}

	if unmarshaled.CreatedAt != machine.CreatedAt {
		t.Error("CreatedAt not preserved during JSON round-trip")
	}
}

func TestComplexConfig(t *testing.T) {
	// Test complex configuration with multiple machines
	config := Config{
		TGCloud: TGCloudConfig{
			User:     "admin@company.com",
			Password: "complexpass123",
		},
		Machines: map[string]MachineConfig{
			"production": {
				Host:     "https://prod.tgcloud.io",
				User:     "admin",
				Password: "prodpass",
				GSPort:   "14240",
				RestPort: "9000",
			},
			"staging": {
				Host:     "https://staging.tgcloud.io",
				User:     "staginguser",
				Password: "stagingpass",
				GSPort:   "14241",
				RestPort: "9001",
			},
			"development": {
				Host:     "http://localhost",
				User:     "tigergraph",
				Password: "tigergraph",
				GSPort:   "14240",
				RestPort: "9000",
			},
		},
		Default: "production",
	}

	// Verify multiple machines
	if len(config.Machines) != 3 {
		t.Errorf("Expected 3 machines, got %d", len(config.Machines))
	}

	// Verify specific machine configurations
	prod := config.Machines["production"]
	if prod.Host != "https://prod.tgcloud.io" {
		t.Error("Production host not set correctly")
	}

	staging := config.Machines["staging"]
	if staging.GSPort != "14241" {
		t.Error("Staging GSPort not set correctly")
	}

	dev := config.Machines["development"]
	if dev.Host != "http://localhost" {
		t.Error("Development host not set correctly")
	}
}

func TestEmptyStructs(t *testing.T) {
	// Test that empty structs can be created
	var config Config
	var tgConfig TGCloudConfig
	var machine MachineConfig
	var cookie GSQLCookie
	var response TGCloudResponse
	var machineInstance Machine

	// These should not panic
	if config.Default != "" {
		t.Error("Empty config should have empty default")
	}

	if tgConfig.User != "" {
		t.Error("Empty TGCloudConfig should have empty user")
	}

	if machine.Host != "" {
		t.Error("Empty MachineConfig should have empty host")
	}

	if cookie.ClientCommit != "" {
		t.Error("Empty GSQLCookie should have empty client commit")
	}

	if response.Message != "" {
		t.Error("Empty TGCloudResponse should have empty message")
	}

	if machineInstance.ID != "" {
		t.Error("Empty Machine should have empty ID")
	}
}
