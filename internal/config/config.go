package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	LogServer     ServerConfig
	MetricsServer ServerConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		LogServer: ServerConfig{
			Port:         getEnvInt("LOG_SERVER_PORT", 5000),
			ReadTimeout:  getEnvDuration("LOG_SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getEnvDuration("LOG_SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getEnvDuration("LOG_SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		MetricsServer: ServerConfig{
			Port:         getEnvInt("METRICS_SERVER_PORT", 9102),
			ReadTimeout:  getEnvDuration("METRICS_SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getEnvDuration("METRICS_SERVER_WRITE_TIMEOUT", 5*time.Second),
			IdleTimeout:  getEnvDuration("METRICS_SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
	}

	return cfg, nil
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
