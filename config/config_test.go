package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func newTestPathProvider(t *testing.T) (OSPathProvider, string) {
	tmp := t.TempDir()
	return OSPathProvider{
		HomeDirFunc: func() (string, error) {
			return tmp, nil
		},
	}, tmp
}

func newTestFileStore(t *testing.T) FileStore {
	p, _ := newTestPathProvider(t)
	return FileStore{
		PathProvider: p,
	}
}

type MockRunner struct {
	Called bool
	Err    error
}

func (m *MockRunner) Run(cmd *exec.Cmd) error {
	m.Called = true
	return m.Err
}
func TestConfigPath(t *testing.T) {
	p, tmpDir := newTestPathProvider(t)
	path, err := p.GetConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath error %v", err)
	}
	if filepath.Base(path) != "config.json" {
		t.Errorf("getConfigpath() = %v, want file namc config.json", filepath.Base(path))
	}
	expectedParent := filepath.Join(tmpDir, ".config", "kpf")
	if filepath.Dir(path) != expectedParent {
		t.Errorf("getConfigPath() parent dir = %v, want %v", filepath.Dir(path), expectedParent)
	}
}

func TestCreateConfigFile(t *testing.T) {
	fs := newTestFileStore(t)
	path, err := fs.Create()
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
		t.Errorf("Config PortForwardRules count=%v, want 2", len(cfg.PortForwardRules))
	}

}

func TestCreateConfigFileAlreadyExists(t *testing.T) {
	pp, _ := newTestPathProvider(t)
	fs := FileStore{PathProvider: pp}

	existingContent := `{"namespace": "test-ns", "port_forward_rules": []}`
	configPath, err := pp.GetConfigPath()

	if err != nil {
		t.Fatalf("getConfigPath error %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	if err := os.WriteFile(configPath, []byte(existingContent), 0644); err != nil {
		t.Fatalf("Failed to write existing config: %v", err)
	}

	path, err := fs.Create()
	if err != nil {
		t.Fatalf("CreateConfigFile() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	if string(data) != existingContent {
		t.Errorf("CreateConfigFile should not overwrite existing file, got %v, want %v", string(data), existingContent)
	}
}

func TestEditor_RunCalled(t *testing.T) {
	p, _ := newTestPathProvider(t)
	mockRunnner := &MockRunner{}
	editor := OSEditor{
		PathProvider: p,
		Runner:       mockRunnner,
	}
	err := editor.Edit()
	if err != nil {
		t.Errorf("Edit() error = %v", err)
	}
	if !mockRunnner.Called {
		t.Error("expected Runner.Run to be called, but i wasn't")
	}
}

func TestEdit_RunFails(t *testing.T) {
	p, _ := newTestPathProvider(t)
	mockRunner := &MockRunner{
		Err: fmt.Errorf("failed to run"),
	}
	editor := OSEditor{
		PathProvider: p,
		Runner:       mockRunner,
	}
	err := editor.Edit()
	if err == nil {
		t.Fatal("expected error")
	}
}
