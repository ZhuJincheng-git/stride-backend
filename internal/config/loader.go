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
	// Set default values for all configuration fields
	setDefaults()
	
	// Load .env file if it exists
	if err := loadEnvFile(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Bind environment variables to viper
	bindEnvVars() 

	// Unmarshal the configuration into the Config struct
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	// Validate the configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.mode", Release)

	// Database defaults
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "")
	viper.SetDefault("mysql.name", "stride_db")
	viper.SetDefault("mysql.params", "charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("mysql.max_idle", 10)
	viper.SetDefault("mysql.max_open", 100)
	viper.SetDefault("mysql.conn_max_lifetime_second", 3600)
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