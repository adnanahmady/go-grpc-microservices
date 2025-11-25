package config

import (
	"bytes"
	"embed"
	"errors"
	"log"

	"github.com/adnanahmady/go-grpc-microservices/pkg/app"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var cfg *Config

func LoadConfig() {
	loadDotEnvFile()
	readConfigFile()
	loadEnvVariables()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
}

func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}
	return cfg
}

func loadDotEnvFile() {
	root := app.GetRootDir()
	if err := godotenv.Load(
		root + "/.env",
	); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}
}

//go:embed config.yaml
var configFile embed.FS

func readConfigFile() {
	fileBytes, err := configFile.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read embedded config file: %v", err)
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewReader(fileBytes)); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			log.Fatalf("config file not found: %v", err)
		}
		log.Fatalf("failed to read config from embedded file: %v", err)
	}
}

func loadEnvVariables() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	mapToStructs()
}
