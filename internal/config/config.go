package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer `yaml:"http-server"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:":8080"`
	Timeout         time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"HTTP_SERVER_IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"server_shutdown_timeout" env:"HTTP_SERVER_SHUTDOWN_TIMEOUT" env-default:"10s"`
}

func (c *Config) Validate() error {
	if c.HTTPServer.Timeout <= 0 {
		return fmt.Errorf("http server timeout must be positive")
	}
	if c.HTTPServer.ShutdownTimeout <= 0 {
		return fmt.Errorf("http server shutdown timeout must be positive")
	}

	return nil
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &cfg, nil
}
