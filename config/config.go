package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Service  ServiceConfig
}

type ServerConfig struct {
	Port            string        `json:"port"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout"`
	RequestTimeout  time.Duration `json:"request_timeout"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
}

type DatabaseConfig struct {
	URL               string
	MaxConns          int32         `json:"max_conns"`
	MinConns          int32         `json:"min_conns"`
	MaxConnLifetime   time.Duration `json:"max_conn_lifetime"`
	MaxConnIdleTime   time.Duration `json:"max_conn_idle_time"`
	HealthCheckPeriod time.Duration `json:"health_check_period"`
}

type ServiceConfig struct {
	MaxReviewersCount int `json:"max_reviewers_count"`
}

func LoadConfig() (*Config, error) {
	cfg, err := loadFromJSON("config/app.json")
	if err != nil {
		cfg = getDefaultConfig()
	}

	if err := loadFromEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadFromJSON(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &cfg, nil
}

func loadFromEnv(cfg *Config) error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://username:password@localhost:5432/pull_requests_reviewer?sslmode=disable"
	}
	cfg.Database.URL = databaseURL

	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}

	return nil
}

func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            "8080",
			ShutdownTimeout: 30 * time.Second,
			RequestTimeout:  5 * time.Second,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
		Database: DatabaseConfig{
			MaxConns:          20,
			MinConns:          5,
			MaxConnLifetime:   time.Hour,
			MaxConnIdleTime:   30 * time.Minute,
			HealthCheckPeriod: time.Minute,
		},
		Service: ServiceConfig{
			MaxReviewersCount: 2,
		},
	}
}