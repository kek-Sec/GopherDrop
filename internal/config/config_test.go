package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Clear environment variables to test defaults
	os.Clearenv()

	cfg := LoadConfig()

	assert.Equal(t, ":8080", cfg.ListenAddr, "ListenAddr should default to :8080")
	assert.Equal(t, "/tmp", cfg.StoragePath, "StoragePath should default to /tmp")
	assert.Equal(t, int64(10*1024*1024), cfg.MaxFileSize, "MaxFileSize should default to 10MB")
}

func TestLoadConfig_EnvVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "disable")
	os.Setenv("SECRET_KEY", "supersecret")
	os.Setenv("LISTEN_ADDR", ":9090")
	os.Setenv("STORAGE_PATH", "/var/data")
	os.Setenv("MAX_FILE_SIZE", "5242880") // 5MB

	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.DBHost, "DBHost should match the environment variable")
	assert.Equal(t, "testuser", cfg.DBUser, "DBUser should match the environment variable")
	assert.Equal(t, "testpass", cfg.DBPass, "DBPass should match the environment variable")
	assert.Equal(t, "testdb", cfg.DBName, "DBName should match the environment variable")
	assert.Equal(t, "disable", cfg.DBSSLMode, "DBSSLMode should match the environment variable")
	assert.Equal(t, "supersecret", cfg.SecretKey, "SecretKey should match the environment variable")
	assert.Equal(t, ":9090", cfg.ListenAddr, "ListenAddr should match the environment variable")
	assert.Equal(t, "/var/data", cfg.StoragePath, "StoragePath should match the environment variable")
	assert.Equal(t, int64(5242880), cfg.MaxFileSize, "MaxFileSize should match the environment variable")
}

func TestLoadConfig_InvalidMaxFileSize(t *testing.T) {
	// Set invalid MAX_FILE_SIZE
	os.Setenv("MAX_FILE_SIZE", "invalid")

	cfg := LoadConfig()

	assert.Equal(t, int64(10*1024*1024), cfg.MaxFileSize, "MaxFileSize should default to 10MB on invalid input")
}
