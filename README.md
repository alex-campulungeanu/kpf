## TODO:
[] add some unit tests
[] add a config file, this should pe checked and created if not exists when the program is started
    - should have the following fields:
        - namespace
        - port_forward_rules
    - the app should be able to read the config file and use the values if not provided as arguments

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() error = %v", err)
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("getConfigPath() = %v, want file name config.json", filepath.Base(path))
	}

	expectedParent := filepath.Join(os.Getenv("HOME"), ".config", "kpf")
	if filepath.Dir(path) != expectedParent {
		t.Errorf("getConfigPath() parent dir = %v, want %v", filepath.Dir(path), expectedParent)
	}
}

func TestCreateConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	path, err := CreateConfigFile()
	if err != nil {
		t.Fatalf("CreateConfigFile() error = %v", err)
	}

	if filepath.Base(path) != "config.json" {
		t.Errorf("CreateConfigFile() = %v, want config.json", filepath.Base(path))
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	var cfg ConfigStructure
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("Config file is not valid JSON: %v", err)
	}

	if cfg.Namespace != "<namespace_name>" {
		t.Errorf("Config namespace = %v, want <namespace_name>", cfg.Namespace)
	}

	if len(cfg.PortForwardRules) != 2 {
		t.Errorf("Config PortForwardRules count = %v, want 2", len(cfg.PortForwardRules))
	}
}

func TestCreateConfigFileAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	existingContent := `{"namespace": "test-ns", "port_forward_rules": []}`
	configDir := filepath.Join(tmpDir, ".config", "kpf")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing config: %v", err)
	}

	path, err := CreateConfigFile()
	if err != nil {
		t.Fatalf("CreateConfigFile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if string(data) != existingContent {
		t.Errorf("CreateConfigFile() should not overwrite existing file, got %v, want %v", string(data), existingContent)
	}
}

func TestReadConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	expectedConfig := ConfigStructure{
		Namespace: "my-namespace",
		PortForwardRules: []PortForwardRule{
			{Prefix: "pod1", Port: "8080"},
			{Prefix: "pod2", Port: "3000"},
		},
	}

	configDir := filepath.Join(tmpDir, ".config", "kpf")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")
	data, _ := json.Marshal(expectedConfig)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	cfg, err := ReadConfigFile()
	if err != nil {
		t.Fatalf("ReadConfigFile() error = %v", err)
	}

	if cfg.Namespace != expectedConfig.Namespace {
		t.Errorf("Namespace = %v, want %v", cfg.Namespace, expectedConfig.Namespace)
	}

	if len(cfg.PortForwardRules) != len(expectedConfig.PortForwardRules) {
		t.Errorf("PortForwardRules count = %v, want %v", len(cfg.PortForwardRules), len(expectedConfig.PortForwardRules))
	}

	for i, rule := range cfg.PortForwardRules {
		if rule.Prefix != expectedConfig.PortForwardRules[i].Prefix {
			t.Errorf("PortForwardRules[%d].Prefix = %v, want %v", i, rule.Prefix, expectedConfig.PortForwardRules[i].Prefix)
		}
		if rule.Port != expectedConfig.PortForwardRules[i].Port {
			t.Errorf("PortForwardRules[%d].Port = %v, want %v", i, rule.Port, expectedConfig.PortForwardRules[i].Port)
		}
	}
}

func TestReadConfigFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	_, err := ReadConfigFile()
	if err == nil {
		t.Error("ReadConfigFile() expected error for non-existent file")
	}
}

func TestReadConfigFileInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	configDir := filepath.Join(tmpDir, ".config", "kpf")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to write invalid config: %v", err)
	}

	_, err := ReadConfigFile()
	if err == nil {
		t.Error("ReadConfigFile() expected error for invalid JSON")
	}
}

func TestEditConfigFileNoEditor(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	originalEditor := os.Getenv("EDITOR")
	os.Unsetenv("EDITOR")
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		}
	}()

	err := EditConfigFile()
	if err == nil {
		t.Error("EditConfigFile() expected error when EDITOR is not set")
	}
}

func TestEditConfigFileWithEditor(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir = func() (string, error) {
		return tmpDir, nil
	}
	defer func() { homeDir = os.UserHomeDir }()

	configDir := filepath.Join(tmpDir, ".config", "kpf")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}
	configPath := filepath.Join(configDir, "config.json")
	if err := os.WriteFile(configPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	originalEditor := os.Getenv("EDITOR")
	os.Setenv("EDITOR", "echo")
	defer func() {
		if originalEditor != "" {
			os.Setenv("EDITOR", originalEditor)
		} else {
			os.Unsetenv("EDITOR")
		}
	}()

	err := EditConfigFile()
	if err != nil {
		t.Errorf("EditConfigFile() error = %v", err)
	}
}
