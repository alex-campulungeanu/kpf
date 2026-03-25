package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

var homeDir = os.UserHomeDir
var execCommand = exec.Command
var runCmd = func(cmd *exec.Cmd) error {
	return cmd.Run()
}

var getConfigPath = func() (string, error) {
	home, err := homeDir()
	if err != nil {
		return "", fmt.Errorf("while trying to fetch config from home dir %w", err)
	}
	configDirPath := filepath.Join(filepath.Join(home, ".config"), configFileDir)
	configFilePath := filepath.Join(configDirPath, configFileName)
	slog.Info("Config file path", "path", configFilePath)
	return configFilePath, nil
}

func CreateConfigFile() (string, error) {
	configFilePath, err := getConfigPath()
	if err != nil {
		return "", err
	}
	configDirPath := filepath.Dir(configFilePath)
	err = os.MkdirAll(configDirPath, 0755)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(configFilePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(configFilePath)
		if err != nil {
			return "", err
		}
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		if err := encoder.Encode(configTemplate); err != nil {
			return "", err
		}
		defer file.Close()
	}

	return configFilePath, nil
}

func ReadConfigFile() (ConfigStructure, error) {
	configFilePath, err := getConfigPath()
	if err != nil {
		return ConfigStructure{}, err
	}
	rawData, err := os.ReadFile(configFilePath)
	configData := ConfigStructure{}
	if err := json.Unmarshal(rawData, &configData); err != nil {
		return ConfigStructure{}, err
	}
	if err != nil {
		return ConfigStructure{}, err
	}
	return configData, nil
}

func EditConfigFile() error {
	configFilePath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("unable to get the config path %w", err)
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("no editor found")
	}
	slog.Info("Trying to edit the config file using", "editor", editor, "configFilePath", configFilePath)
	cmd := execCommand(editor, configFilePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return runCmd(cmd)
}

func Init() {
	filePath, err := CreateConfigFile()
	if err != nil {
		slog.Error("error create config file %v", "err", err)
	}

	configFile, err := os.Open(filePath)
	if err != nil {
		slog.Error("error open config file %v", "err", err)
	}
	defer configFile.Close()
	byteValue, _ := io.ReadAll(configFile)
	var Config ConfigStructure
	if err := json.Unmarshal(byteValue, &Config); err != nil {
		slog.Error("error unmarshal config file %v", "err", err)
	}

}
