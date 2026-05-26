package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
	Port    string        `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
	Timeout time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" env-default:"10s"`
}

type DatabaseConfig struct {
	Url string `yaml:"url" env:"DATABASE_URL" env-required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	configPath := "config.yaml"

	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	if _, err := os.Stat(configPath); err == nil {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("Couldn't load configuration from file: %w", err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("Couldn't load configuration from environment: %w", err)
		}
	}

	return &cfg, nil
}
