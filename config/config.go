package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ENV string

const (
	ENV_LOCAL ENV = "local"
	ENV_PROD  ENV = "production"
)

type Config struct {
	Environment    ENV    `default:"local" envconfig:"APP_ENV"`
	APIHost        string `default:"127.0.0.1:8080" envconfig:"API_HOST"`
	AdminJWTSecret string `default:"bruh" envconfig:"ADMIN_JWT_SECRET"`
	DB             struct {
		DBUser string `default:"postgres" envconfig:"DB_USER"`
		DBPass string `default:"postgres" envconfig:"DB_PASS"`
		DBURI  string `default:"localhost:5432" envconfig:"DB_URI"`
		DBName string `default:"syncspace" envconfig:"DB_NAME"`
	}
}

var (
	cfg *Config
)

func Get() (*Config, error) {
	if cfg == nil {
		godotenv.Load(".env")
		var config Config
		if err := envconfig.Process("", &config); err != nil {
			return nil, fmt.Errorf("[db.Get] failed to process env vars: %w", err)
		}
		cfg = &config
	}
	return cfg, nil
}
