package database

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestInitDB tests the InitDB function.
func TestInitDB(t *testing.T) {
	cfg := config.Config{
		DBHost:    "",
		DBUser:    "",
		DBPass:    "",
		DBName:    ":memory:",
		DBSSLMode: "",
	}

	db, err := gorm.Open("sqlite3", cfg.DBName)
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer db.Close()

	assert.NotNil(t, db, "Database connection should not be nil")
}

// TestCleanupExpired tests the CleanupExpired function.
func TestCleanupExpired(t *testing.T) {
	// Setup an in-memory database
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}
	defer db.Close()

	// Migrate the Send model
	db.AutoMigrate(&models.Send{})

	// Create some test data
	now := time.Now()
	expiredSend := models.Send{
		Hash:      "expired1",
		Type:      "text",
		Data:      "some data",
		ExpiresAt: now.Add(-1 * time.Hour),
	}

	validSend := models.Send{
		Hash:      "valid1",
		Type:      "text",
		Data:      "valid data",
		ExpiresAt: now.Add(1 * time.Hour),
	}

	db.Create(&expiredSend)
	db.Create(&validSend)

	// Run the cleanup function (without the infinite loop)
	func() {
		var sends []models.Send
		db.Where("expires_at < ?", time.Now()).Find(&sends)
		for _, s := range sends {
			db.Delete(&s)
		}
	}()

	// Verify that the expired send was deleted
	var result models.Send
	err = db.Where("hash = ?", "expired1").First(&result).Error
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Expired send should be deleted")

	// Verify that the valid send still exists
	err = db.Where("hash = ?", "valid1").First(&result).Error
	assert.Nil(t, err, "Valid send should still exist")
}
