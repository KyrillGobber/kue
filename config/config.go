package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigType struct {
	BridgeAddress string
}

const filename = "kue.conf"

var config *ConfigType

func LoadConfig() error {
	var err error
	config, err = loadOrCreateConfig()
	if err != nil {
		return err
	}
    return err
}

func GetConfig() *ConfigType {
	return config
}

func getConfigPath() (string, error) {
    var dev bool
	flag.BoolVar(&dev, "dev", false, "Run in development mode")
	flag.Parse()
    if dev {
        return "/home/ky/projects/kue/kue.conf", nil
    }
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(exePath), filename), nil
}

func loadOrCreateConfig() (*ConfigType, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %v", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config doesn't exist, create a new one
		return createNewConfig(configPath)
	} else {
		// Config exists, load it
		return readConfig(configPath)
	}
}

func readConfig(configPath string) (*ConfigType, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %v", err)
	}
	defer file.Close()

	var config ConfigType
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %v", err)
	}

	return &config, nil
}

func createNewConfig(configPath string) (*ConfigType, error) {
    address := discoveryProcess()
    if address == "" {
        return nil, fmt.Errorf("no bridge selected")
    }
	config := &ConfigType{
        BridgeAddress: fmt.Sprintf("https://%s", address),
	}

	file, err := os.Create(configPath)
	if err != nil {
		return nil, fmt.Errorf("error creating config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print
	if err := encoder.Encode(config); err != nil {
		return nil, fmt.Errorf("error encoding config: %v", err)
	}

	fmt.Println("New configuration file created:", configPath)
	return config, nil
}
