package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	DefaultConfigPath       = "termigochi-config.json"
	DefaultPetStateFilePath = "termigochi-pet-state.json"
)

type Config struct {
	IsFirstRun       bool   `json:"is_first_run"`
	PlayerName       string `json:"player_name"`
	PetStateFilePath string `json:"pet_state_file_path"`
	ConfigPath       string `json:"config_path"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{
		IsFirstRun: true,
		PlayerName: "",
	}

	file, err := os.Create(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to write to config file: %w", err)
	}

	return config, nil
}

func LoadConfig(configPath string) (*Config, error, bool) {
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			config, err := NewConfig(configPath)
			if err != nil {
				return nil, err, false
			}

			return config, nil, true
		}
		return nil, fmt.Errorf("failed to open config file: %w", err), false
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err), false
	}

	return &config, nil, false
}

func (c *Config) SaveConfig() error {
	if c.ConfigPath == "" {
		c.ConfigPath = DefaultConfigPath
		c.PetStateFilePath = DefaultPetStateFilePath
	}

	file, err := os.Create(c.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(c)
	if err != nil {
		return fmt.Errorf("failed to write to config file: %w", err)
	}

	return nil
}
