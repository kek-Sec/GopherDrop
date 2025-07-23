// Package handlers contains logic for creating and retrieving sends.
package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/models"
	"github.com/kek-Sec/gopherdrop/internal/security"
)

// CreateSend handles creation of a new send.
// It accepts form data for type (text/file), optional password, one-time use, and expiration.
func CreateSend(cfg config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Explicitly parse the multipart form before accessing form values.
		// This is crucial for handling mixed file/text forms reliably in Gin.
		if err := c.Request.ParseMultipartForm(cfg.MaxFileSize + 1024*1024); err != nil { // Add buffer to max size
			log.Println("Error parsing multipart form:", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid form data or file too large"})
			return
		}

		// Use c.Request.FormValue now that the form is parsed.
		stype := c.Request.FormValue("type")
		pw := c.Request.FormValue("password")
		ot := c.Request.FormValue("onetime")
		exp := c.Request.FormValue("expires")

		log.Println("CreateSend called with type:", stype)

		if stype == "" {
			log.Println("Error: 'type' field is missing")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Type field is required"})
			return
		}

		oneTime := (ot == "true")
		log.Println("One-Time:", oneTime)

		var expiresAt time.Time
		if exp != "" {
			d, err := time.ParseDuration(exp)
			if err != nil {
				log.Println("Error parsing expiration duration:", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid expiration duration"})
				return
			}
			expiresAt = time.Now().Add(d)
		} else {
			expiresAt = time.Now().Add(24 * time.Hour)
		}
		log.Println("Expires At:", expiresAt)

		hash, err := security.GenerateHash(16)
		if err != nil {
			log.Println("Error generating hash:", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate hash"})
			return
		}
		log.Println("Generated Hash:", hash)

		key := deriveKey(pw, cfg)

		if stype == "text" {
			text := c.Request.FormValue("data")
			if text == "" {
				log.Println("Error: 'data' field is empty for text type")
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Data field is required for text type"})
				return
			}

			enc, err := security.EncryptData([]byte(text), key)
			if err != nil {
				log.Println("Error encrypting text data:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt text data"})
				return
			}

			s := models.Send{
				Hash:      hash,
				Type:      "text",
				Data:      enc,
				Password:  pw,
				OneTime:   oneTime,
				ExpiresAt: expiresAt,
			}
			db.Create(&s)
			log.Println("Text send created successfully with hash:", hash)
			c.JSON(http.StatusOK, gin.H{"hash": s.Hash})
			return
		}

		if stype == "file" {
			// Use c.Request.FormFile now that the form is parsed.
			file, header, err := c.Request.FormFile("file")
			if err != nil {
				log.Println("Error retrieving file from form data:", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from form data"})
				return
			}
			defer file.Close()

			log.Println("Received file:", header.Filename, "Size:", header.Size)

			if header.Size > cfg.MaxFileSize {
				log.Printf("Error: File size (%d bytes) exceeds maximum allowed size (%d bytes)\n", header.Size, cfg.MaxFileSize)
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds the maximum allowed limit"})
				return
			}

			data, err := io.ReadAll(file)
			if err != nil {
				log.Println("Error reading file data:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file data"})
				return
			}

			enc, err := security.EncryptData(data, key)
			if err != nil {
				log.Println("Error encrypting file data:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file data"})
				return
			}

			fp := filepath.Join(cfg.StoragePath, hash)
			if err := os.WriteFile(fp, []byte(enc), 0600); err != nil {
				log.Println("Error writing encrypted file to storage:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to write encrypted file to storage"})
				return
			}

			log.Println("File saved successfully to:", fp)

			s := models.Send{
				Hash:      hash,
				Type:      "file",
				FilePath:  fp,
				FileName:  header.Filename,
				Password:  pw,
				OneTime:   oneTime,
				ExpiresAt: expiresAt,
			}
			db.Create(&s)
			log.Println("File send created successfully with hash:", hash)
			c.JSON(http.StatusOK, gin.H{"hash": s.Hash})
			return
		}

		log.Println("Error: Unsupported send type:", stype)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Unsupported send type"})
	}
}

// GetSend handles retrieving and decrypting a send by its hash.
func GetSend(cfg config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("id")
		var s models.Send

		if db.First(&s, "hash = ?", hash).RecordNotFound() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if time.Now().After(s.ExpiresAt) {
			deleteSendAndFile(db, &s)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		pw := c.Query("password")
		if s.Password != "" && s.Password != pw {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		key := deriveKey(pw, cfg)

		if s.Type == "text" {
			d, err := security.DecryptData(s.Data, key)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.String(http.StatusOK, string(d))
		} else {
			d, err := os.ReadFile(s.FilePath)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			dec, err := security.DecryptData(string(d), key)
			if err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, s.FileName))
			c.Data(http.StatusOK, "application/octet-stream", dec)
		}

		if s.OneTime {
			deleteSendAndFile(db, &s)
		}
	}
}

func deriveKey(pw string, cfg config.Config) []byte {
	if pw != "" {
		return []byte(security.PadKey(pw))
	}
	return []byte(security.PadKey(cfg.SecretKey))
}

func deleteSendAndFile(db *gorm.DB, s *models.Send) {
	if s.Type == "file" && s.FilePath != "" {
		os.Remove(s.FilePath)
	}
	db.Delete(&s)
}

// CheckPasswordProtection checks if a send requires a password.
func CheckPasswordProtection(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		hash := c.Param("id")
		var s models.Send

		if db.First(&s, "hash = ?", hash).RecordNotFound() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if time.Now().After(s.ExpiresAt) {
			deleteSendAndFile(db, &s)
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Return whether the send requires a password
		c.JSON(http.StatusOK, gin.H{"requiresPassword": s.Password != ""})
	}
}