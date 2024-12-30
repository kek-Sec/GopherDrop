// Package config provides application configuration loading from environment variables.
package config

import (
	"os"
	"strconv"
)

// Config holds all environment-based configuration.
type Config struct {
	DBHost      string
	DBUser      string
	DBPass      string
	DBName      string
	DBSSLMode   string
	SecretKey   string
	ListenAddr  string
	StoragePath string
	MaxFileSize int64
}

// LoadConfig returns a Config object populated from environment variables.
func LoadConfig() Config {
	cfg := Config{
		DBHost:      os.Getenv("DB_HOST"),
		DBUser:      os.Getenv("DB_USER"),
		DBPass:      os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBSSLMode:   os.Getenv("DB_SSLMODE"),
		SecretKey:   os.Getenv("SECRET_KEY"),
		ListenAddr:  os.Getenv("LISTEN_ADDR"),
		StoragePath: os.Getenv("STORAGE_PATH"),
	}

	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":8080"
	}
	if cfg.StoragePath == "" {
		cfg.StoragePath = "/tmp"
	}

	maxFileSizeStr := os.Getenv("MAX_FILE_SIZE")
	if maxFileSizeStr == "" {
		// Default 10MB if not set
		cfg.MaxFileSize = 10 * 1024 * 1024
	} else {
		size, err := strconv.ParseInt(maxFileSizeStr, 10, 64)
		if err != nil {
			cfg.MaxFileSize = 10 * 1024 * 1024
		} else {
			cfg.MaxFileSize = size
		}
	}

	return cfg
}
