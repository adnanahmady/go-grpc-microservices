package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App  AppConfig  `mapstructure:"app"`
	User UserConfig `mapstructure:"user"`
	Log  LogConfig  `mapstructure:"log"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type UserConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Level            string `mapstructure:"level"`
	Dir              string `mapstructure:"dir"`
	WriteToFile      bool   `mapstructure:"write_to_file"`
	CompressLogFiles bool   `mapstructure:"compress_log_files"`
	MaxAge           int    `mapstructure:"max_age"`  // days
	MaxSize          int    `mapstructure:"max_size"` // MB
}

func mapToStructs() {
	// App config
	loadEnvToStruct("app.env", "APP_ENV")

	// User config
	loadEnvToStruct("user.host", "USER_HOST")
	loadEnvToStruct("user.port", "USER_PORT")

	// Log config
	loadEnvToStruct("log.level", "LOG_LEVEL")
	loadEnvToStruct("log.dir", "LOG_DIR")
	loadEnvToStruct("log.write_to_file", "LOG_WRITE_TO_FILE")
	loadEnvToStruct("log.compress_log_files", "LOG_COMPRESS_LOG_FILES")
	loadEnvToStruct("log.max_age", "LOG_MAX_AGE")
	loadEnvToStruct("log.max_size", "LOG_MAX_SIZE")
}

func loadEnvToStruct(field, key string) {
	if err := viper.BindEnv(field, key); err != nil {
		log.Fatalf("failed to bind env to struct field: %v", err)
	}
}
