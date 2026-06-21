package config

import (
	"fmt"
	"time"
)

type Environment string

const (
	Debug   Environment = "debug"
	Release Environment = "release"
	Test    Environment = "test"
)

// Config holds the configuration values for the application.
type Config struct {
	// App
	AppPort int         `mapstructure:"app_port" validate:"required,min=1,max=65535"`
	AppMode Environment `mapstructure:"app_mode" validate:"required,oneof=debug test release"`

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
}

// IsReleaseMode checks if the server is running in release mode.
func (c *Config) IsReleaseMode() bool {
	return c.AppMode == Release
}

// IsDebugMode checks if the server is running in debug mode.
func (c *Config) IsDebugMode() bool {
	return c.AppMode == Debug
}

// ConnMaxLifetime returns the connection max lifetime as a time.Duration.
func (c *Config) ConnMaxLifetime() time.Duration {
	return time.Duration(c.DBConnMaxLifetimeSecond) * time.Second
}

// MySQLDSN constructs the MySQL Data Source Name (DSN) from the configuration.
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBParams)
}
