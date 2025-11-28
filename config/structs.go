package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig       `mapstructure:"app"`
	User      UserConfig      `mapstructure:"user"`
	Inventory InventoryConfig `mapstructure:"inventory"`
	Order     OrderConfig     `mapstructure:"order"`
	Gateway   GatewayConfig   `mapstructure:"gateway"`
	Log       LogConfig       `mapstructure:"log"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type UserConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type InventoryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type OrderConfig struct {
	Host string `mapstrucutre:"host"`
	Port int    `mapstrucutre:"port"`
}

type GatewayConfig struct {
	User      ServiceConfig `mapstructure:"user"`
	Inventory ServiceConfig `mapstructure:"inventory"`
}

type ServiceConfig struct {
	Host string `mapstrucutre:"host"`
	Port int    `mapstrucutre:"port"`
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

	// Inventory config
	loadEnvToStruct("inventory.host", "INVENTORY_HOST")
	loadEnvToStruct("inventory.port", "INVENTORY_PORT")

	// Order config
	loadEnvToStruct("order.host", "ORDER_HOST")
	loadEnvToStruct("order.port", "ORDER_PORT")

	// Gateway config
	loadEnvToStruct("gateway.user.host", "GATEWAY_USER_HOST")
	loadEnvToStruct("gateway.user.port", "GATEWAY_USER_PORT")

	loadEnvToStruct("gateway.inventory.host", "GATEWAY_INVENTORY_HOST")
	loadEnvToStruct("gateway.inventory.port", "GATEWAY_INVENTORY_PORT")

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
