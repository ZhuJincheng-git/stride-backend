// ./internal/config/config.go
package config

type Environment string

const (
	Debug   Environment = "debug"
	Release Environment = "release"
	Test    Environment = "test"
)

// Config holds the configuration values for the application.
type Config struct {
	// Server configuration
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
}

type ServerConfig struct {
	Port int         `mapstructure:"port" validate:"required,min=1,max=65535"`
	Mode Environment `mapstructure:"mode" validate:"required,oneof=debug test release"`
}

type MySQLConfig struct {
	Host            string `mapstructure:"host" validate:"required"`
	Port            int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	User            string `mapstructure:"user" validate:"required"`
	Password        string `mapstructure:"password"`
	Name            string `mapstructure:"name" validate:"required"`
	Charset         string `mapstructure:"charset"`
	MaxIdle         int    `mapstructure:"max_idle"`
	MaxOpen         int    `mapstructure:"max_open"`
	ConnMaxLifetime int64  `mapstructure:"conn_max_lifetime"`
}

// IsReleaseMode checks if the server is running in release mode.
func (c *Config) IsReleaseMode() bool {
	return c.Server.Mode == Release
}

// IsDebugMode checks if the server is running in debug mode.
func (c *Config) IsDebugMode() bool {
	return c.Server.Mode == Debug
}