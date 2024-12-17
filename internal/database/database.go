// Package database handles database connection and maintenance routines.
package database

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/models"
)

// InitDB connects to PostgreSQL using environment variables.
func InitDB(cfg config.Config) *gorm.DB {
	dsn := "host=" + cfg.DBHost + " user=" + cfg.DBUser + " password=" + cfg.DBPass + " dbname=" + cfg.DBName + " sslmode=" + cfg.DBSSLMode
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// CleanupExpired periodically removes expired sends from the database.
func CleanupExpired(db *gorm.DB) {
	for {
		time.Sleep(time.Hour)
		var sends []models.Send
		db.Where("expires_at < ?", time.Now()).Find(&sends)
		for _, s := range sends {
			if s.Type == "file" && s.FilePath != "" {
				removeFile(s.FilePath)
			}
			db.Delete(&s)
		}
	}
}
