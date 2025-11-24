package config

import (
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

func readConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(app.GetRootDir())

	log.Printf("reading config file: %s", app.GetRootDir()+"/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) {
			log.Fatalf("config file not found: %v", err)
			return
		}
		log.Fatalf("failed to read config file: %v", err)
	}
}

func loadEnvVariables() {
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	mapToStructs()
}
