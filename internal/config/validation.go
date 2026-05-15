// ./internal/config/validation.go
package config

import (
	"github.com/go-playground/validator/v10"
)

var configValidator = newConfigValidator()

func newConfigValidator() *validator.Validate {
	validate := validator.New()
	return validate
}

func validateConfig(cfg *Config) error {
	return configValidator.Struct(cfg)
}
