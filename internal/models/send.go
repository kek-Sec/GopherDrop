package models

import "time"

// Send represents a stored secret send, either text or file.
type Send struct {
	Hash      string    `gorm:"primary_key"` // Random hash as the unique identifier
	Type      string
	Data      string
	FilePath  string
	FileName  string
	Password  string
	OneTime   bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
