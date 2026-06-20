package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
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
	// App
	viper.SetDefault("app_port", 8080)
	viper.SetDefault("app_mode", Release)

	// Database
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_port", 3306)
	viper.SetDefault("db_user", "root")
	viper.SetDefault("db_password", "")
	viper.SetDefault("db_name", "stride_db")
	viper.SetDefault("db_params", "charset=utf8mb4&parseTime=True&loc=Local")
	viper.SetDefault("db_max_idle", 10)
	viper.SetDefault("db_max_open", 100)
	viper.SetDefault("db_conn_max_lifetime_second", 3600)
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
	viper.AutomaticEnv()
}