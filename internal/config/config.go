// ./internal/config/config.go
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
	App   AppConfig   `mapstructure:"app"`
	MySQL MySQLConfig `mapstructure:"mysql"`
}

type AppConfig struct {
	Port int         `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mode Environment `mapstructure:"mode" validate:"required,oneof=debug test release"`
}

type MySQLConfig struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	User            string `mapstructure:"user" validate:"required"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name" validate:"required"`
	Params         string `mapstructure:"params"`
	MaxIdle         int    `mapstructure:"max_idle"`
	MaxOpen         int    `mapstructure:"max_open"`
	ConnMaxLifetimeSecond int32  `mapstructure:"conn_max_lifetime_second"`
}

// IsReleaseMode checks if the server is running in release mode.
func (c *Config) IsReleaseMode() bool {
	return c.App.Mode == Release
}

// IsDebugMode checks if the server is running in debug mode.
func (c *Config) IsDebugMode() bool {
	return c.App.Mode == Debug
}

// ConnMaxLifetime returns the connection max lifetime as a time.Duration.
func (c *MySQLConfig) ConnMaxLifetime() time.Duration {
	return time.Duration(c.ConnMaxLifetimeSecond) * time.Second
}

// MySQLDSN constructs the MySQL Data Source Name (DSN) from the configuration.
func (c *MySQLConfig) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Params)
}