package config

import (
	"strings"
	"github.com/go-playground/validator/v10"
)

var configValidator = newConfigValidator()

func newConfigValidator() *validator.Validate {
	validate := validator.New()
	return validate
}

func validateConfig(cfg *Config) error {
	configValidator.RegisterStructValidation(validation, cfg)
	return configValidator.Struct(cfg)
}

func validation(fl validator.StructLevel) {
	cfg := fl.Current().Interface().(Config)
	// DBPassword
	if cfg.AppMode == Release && (cfg.DBPassword == "" || strings.TrimSpace(cfg.DBPassword) == "" || strings.Contains(cfg.DBPassword, "password")) {
		fl.ReportError(cfg.DBPassword, "DBPassword", "db_password", "required", "")
	}
}