// ./internal/config/loader.go
package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// Load loads the configuration from the .env file and environment variables.
func Load() (*Config, error) {
	setDefaults() // Set default values for all configuration fields
	
	// Load .env file if it exists
	if err := loadEnvFile(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	bindEnvVars() // Bind environment variables to viper

	var cfg Config
	err := viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setDefaults() {
	// app defaults
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.mode", Debug)

	// database defaults
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", "3306")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "")
	viper.SetDefault("mysql.name", "stride_db")
	viper.SetDefault("mysql.charset", "utf8mb4")
	viper.SetDefault("mysql.max_idle", 10)
	viper.SetDefault("mysql.max_open", 100)
	viper.SetDefault("mysql.conn_max_lifetime", 3600)
}

// loadEnvFile loads the .env file using godotenv.
func loadEnvFile() error {
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		return godotenv.Load(envPath)
	}
	return nil
}

// bindEnvVars binds environment variables to viper keys.
func bindEnvVars() {
	viper.SetEnvPrefix("STRIDE_BACKEND")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
}
