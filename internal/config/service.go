package config

import (
	"log"
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

const (
	DefaultListenedPort = 8787
	DefaultDBPort = 5432
)

func Init(prod bool) (config *ServiceConfig, err error) {
	configPath := ".env"
	if !prod {
		configPath = ".env.development"
	}
	
	_, fileErr := os.Stat(configPath)
	if fileErr != nil {
		log.Printf("Warning: Failed opening %s file\n", configPath)
	}

	if fileErr == nil {
		err = godotenv.Load(configPath)
		if err != nil {
			return nil, err
		}
	}

	portStr := os.Getenv("SERVE_PORT")
	port, portParseErr := strconv.Atoi(portStr)
	if portParseErr != nil {
		log.Printf("Warning: Failed reading the SERVE_PORT env. Using %d.\n", DefaultListenedPort)

		port = DefaultListenedPort
	}

	config = &ServiceConfig{}

	config.ServePort = uint16(port)
	config.ServeEnpoint = os.Getenv("SERVE_ENDPOINT")
	config.DBHost = os.Getenv("DB_HOST")

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, dbPortParseErr := strconv.Atoi(dbPortStr)
	if dbPortParseErr != nil {
		log.Printf("Warning: Failed reading the DB_PORT env. Using %d.\n", DefaultDBPort)

		dbPort = DefaultDBPort
	}

	config.DBPort = uint16(dbPort)
	config.DBUser = os.Getenv("DB_USER")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBName = os.Getenv("DB_NAME")

	return config, nil
}
