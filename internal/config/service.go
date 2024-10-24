package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServiceConfig struct {
	ServePort    uint16 `json:"SERVE_PORT"`
	ServeEnpoint string `json:"SERVE_ENDPOINT"`
	DBHost string `json:"DB_HOST"`
	DBPort uint16 `json:"DB_PORT"`
	DBUser string `json:"DB_USER"`
	DBPassword string `json:"DB_PASSWORD"`
	DBName string `json:"DB_NAME"`
}

func Init(prod bool) (config *ServiceConfig, err error) {
	configPath := ".env"
	if !prod {
		configPath = ".env.development"
	}
	
	err = godotenv.Load(configPath)
	if err != nil {
		return nil, err
	}

	portStr := os.Getenv("SERVE_PORT")
	port, portParseErr := strconv.Atoi(portStr)
	if portParseErr != nil {
		return nil, portParseErr
	}

	config = &ServiceConfig{}

	config.ServePort = uint16(port)
	config.ServeEnpoint = os.Getenv("SERVE_ENDPOINT")
	config.DBHost = os.Getenv("DB_HOST")

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, dbPortParseErr := strconv.Atoi(dbPortStr)
	if dbPortParseErr != nil {
		return nil, dbPortParseErr
	}

	config.DBPort = uint16(dbPort)
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")

	return config, nil
}
