package config

import (
	"flag"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type PostgresConfig struct {
	User     string `envconfig:"POSTGRES_USER" required:"true"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	Name     string `envconfig:"POSTGRES_NAME" required:"true"`
	Port     string `envconfig:"POSTGRES_PORT" required:"true"`
	Host     string `envconfig:"POSTGRES_HOST" default:"postgres"`
}

type Config struct {
	Port     int            `envconfig:"SERVICE_PORT" required:"true"`
	LogLevel string         `envconfig:"LOG_LEVEL" default:"DEBUG"`
	Postgres PostgresConfig `envconfig:"POSTGRES"`
}

func Load() (*Config, error) {
	loadEnvFile()

	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	slog.Info("config is loaded", slog.Any("config", cfg))

	return &cfg, nil
}

func loadEnvFile() {
	envFile := ".env"

	if checkLocal() {
		envFile = ".env.local"
	}

	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			slog.Error("failed to load .env file", "error", err.Error())
		}
	}
}

func checkLocal() bool {
	isLocal := flag.Bool("local", false, "check local development")
	flag.Parse()
	return *isLocal
}
