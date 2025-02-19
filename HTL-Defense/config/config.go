package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort string `yaml:"server_port"`
	LogPath    string `yaml:"log_path"`
	// Additional fields can be added as needed.
}

var configInstance *Config

// GetConfig reads config.yaml, validates it and returns the configuration.
func GetConfig() (*Config, error) {
	if configInstance != nil {
		return configInstance, nil
	}

	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %v", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	// Validate required fields.
	if cfg.ServerPort == "" || cfg.LogPath == "" {
		return nil, fmt.Errorf("missing required configuration fields")
	}

	configInstance = &cfg
	return configInstance, nil
}
