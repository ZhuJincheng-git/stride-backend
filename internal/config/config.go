package config

import (
	"fmt"
	"time"
)

type Environment string

const (
	Develop Environment = "develop"
	Release Environment = "release"
	Test    Environment = "test" // Test environment allows Config validation failures; use it in tests.
)

// Config holds the configuration values for the application.
type Config struct {
	// App
	AppPort int         `mapstructure:"app_port" validate:"required,min=1,max=65535"`
	AppEnv  Environment `mapstructure:"app_env" validate:"required,oneof=develop test release"`

	// Database
	DBHost                  string `mapstructure:"db_host" validate:"required"`
	DBPort                  int    `mapstructure:"db_port" validate:"required,min=1,max=65535"`
	DBUser                  string `mapstructure:"db_user" validate:"required"`
	DBPassword              string `mapstructure:"db_password"`
	DBName                  string `mapstructure:"db_name" validate:"required"`
	DBParams                string `mapstructure:"db_params"`
	DBMaxIdle               int    `mapstructure:"db_max_idle"`
	DBMaxOpen               int    `mapstructure:"db_max_open"`
	DBConnMaxLifetimeSecond int32  `mapstructure:"db_conn_max_lifetime_second"`

	// Auth
	JWTSecret       string `mapstructure:"jwt_secret" validate:"required"`
	JWTExpiresHours int    `mapstructure:"jwt_expires_hours"`

	// Logging
	LogLevel string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

// IsReleaseEnv checks if the server is running in release environment.
func (c *Config) IsReleaseEnv() bool {
	return c.AppEnv == Release
}

// IsDevelopEnv checks if the server is running in develop environment.
func (c *Config) IsDevelopEnv() bool {
	return c.AppEnv == Develop
}

// ConnMaxLifetime returns the connection max lifetime as a time.Duration.
func (c *Config) ConnMaxLifetime() time.Duration {
	return time.Duration(c.DBConnMaxLifetimeSecond) * time.Second
}

// MySQLDSN constructs the MySQL Data Source Name (DSN) from the configuration.
func (c *Config) MySQLDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
	if c.DBParams != "" {
		dsn += "?" + c.DBParams
	}
	return dsn
}

func (c *Config) JWTExpires() time.Duration {
	return time.Duration(c.JWTExpiresHours) * time.Hour
}
