package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var envFileCandidates = []string{
	// High priority: .env.local, then .env
	".env.local",
	".env",
}

// Load loads the configuration from the .env file and environment variables.
func Load() (*Config, error) {
	// Set default values for all configuration fields
	setDefaults()
	
	// Load .env file if it exists
	loadEnvFile()

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

	// Logging
	viper.SetDefault("log_level", "info")
}

// findUpwards searches for the target file starting from the startDir and moving up the directory tree.
func findUpwards(startDir, targetFile string) string {
	currentDir := startDir
	for {
		candidatePath := filepath.Join(currentDir, targetFile)
		if info, err := os.Stat(candidatePath); err == nil && !info.IsDir() {
			return candidatePath
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return ""
		}
		currentDir = parentDir
	}
}

// loadEnvFile loads the .env file using godotenv.
func loadEnvFile() {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	for _, envFile := range envFileCandidates {
		if envPath := findUpwards(cwd, envFile); envPath != "" {
			if err := godotenv.Load(envPath); err != nil {
				log.Printf("config: fail to load %s: %v", envPath, err)
			}
		}
	}
}

// bindEnvVars binds environment variables to viper keys.
func bindEnvVars() {
	viper.AutomaticEnv()
}