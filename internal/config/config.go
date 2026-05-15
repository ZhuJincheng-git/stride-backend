package config

import (
	"os"
	"path/filepath"
	"strings"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	envPath = filepath.Join(".", ".env")
)

// Config holds the configuration values for the application.
type Config struct {
	// Server configuration
	ServerPort string `mapstructure:"SERVER_PORT"`
	ServerMode string `mapstructure:"SERVER_MODE"`

	// Database configuration
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBCharset  string `mapstructure:"DB_CHARSET"`
	DBMaxIdle int `mapstructure:"DB_MAX_IDLE"`
	DBMaxOpen int `mapstructure:"DB_MAX_OPEN"`
	DBConnMaxLifetime int `mapstructure:"DB_CONN_MAX_LIFETIME"`
}

var AppConfig *Config

// Load loads the configuration from the .env file and environment variables.
func Load() (*Config, error) {
	setDefaults() // Set default values for all configuration fields
	loadEnvFile() // Load .env file if it exists
	bindEnvVars() // Bind environment variables to viper

	config := &Config{
		ServerPort: viper.GetString("server.port"),
		ServerMode: viper.GetString("server.mode"),
		DBHost:     viper.GetString("db.host"),
		DBPort:     viper.GetString("db.port"),
		DBUser: viper.GetString("db.user"),
		DBPassword: viper.GetString("db.password"),
		DBName: viper.GetString("db.name"),
		DBCharset: viper.GetString("db.charset"),
		DBMaxIdle: viper.GetInt("db.max_idle"),
		DBMaxOpen: viper.GetInt("db.max_open"),
		DBConnMaxLifetime: viper.GetInt("db.conn_max_lifetime"),
	}
	AppConfig = config
	return config, nil
}

func setDefaults() {
	// server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// database defaults
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", "3306")
	viper.SetDefault("db.user", "root")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.name", "stride_db")
	viper.SetDefault("db.charset", "utf8mb4")
	viper.SetDefault("db.max_idle", 10)
	viper.SetDefault("db.max_open", 100)
	viper.SetDefault("db.conn_max_lifetime", 3600)
}

// loadEnvFile loads the .env file using godotenv.
func loadEnvFile() error {
	if _, err := os.Stat(envPath); err == nil {
		return godotenv.Load(envPath)
	}
	return nil
}

// bindEnvVars binds environment variables to viper keys.
func bindEnvVars() {
	viper.SetEnvPrefix("STRIDE_BACKEND")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}