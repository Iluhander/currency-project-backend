package config

import (
	"encoding/json"
	"os"
)

type ServiceConfig struct {
	ServePort    uint16 `json:"SERVE_PORT"`
	ServeEnpoint string `json:"SERVE_ENDPOINT"`
	DBHost string `json:"DB_HOST"`
	DBPort int `json:"DB_PORT"`
	DBUser string `json:"DB_USER"`
	DBPassword string `json:"DB_PASSWORD"`
	DBName string `json:"DB_NAME"`
}

func Init(prod bool) (config *ServiceConfig, err error) {
	configPath := "config.json"
	if !prod {
		configPath = "config.dev.json"
	}

	jsonContents, readErr := os.ReadFile(configPath)

	if readErr != nil {
		return nil, readErr
	}

	parseErr := json.Unmarshal(jsonContents, &config)

	if parseErr != nil {
		return nil, parseErr
	}

	return config, nil
}
