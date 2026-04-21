// Package handlers contains logic for creating and retrieving sends.
package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	return createSendWithType(cfg, db, "")
}

// CreateTextSend handles text send creation without requiring a type field.
func CreateTextSend(cfg config.Config, db *gorm.DB) gin.HandlerFunc {
	return createSendWithType(cfg, db, "text")
}

// CreateFileSend handles file send creation without requiring a type field.
func CreateFileSend(cfg config.Config, db *gorm.DB) gin.HandlerFunc {
	return createSendWithType(cfg, db, "file")
}

func createSendWithType(cfg config.Config, db *gorm.DB, forcedType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		stype := strings.TrimSpace(forcedType)
		if stype == "" {
			stype = strings.TrimSpace(c.PostForm("type"))
			if stype == "" {
				stype = strings.TrimSpace(c.Query("type"))
			}
		}

		pw := firstNonEmpty(c.PostForm("password"), c.Query("password"))
		ot := firstNonEmpty(c.PostForm("onetime"), c.Query("onetime"))
		exp := firstNonEmpty(c.PostForm("expires"), c.Query("expires"))

		log.Println("CreateSend called with type:", stype)

		if stype == "" {
			log.Println("Error: 'type' field is missing")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Type field is required"})
			return
		}

		oneTime, err := strconv.ParseBool(ot)
		if ot == "" {
			oneTime = false
		} else if err != nil {
			log.Println("Error parsing onetime flag:", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid onetime value"})
			return
		}
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
			text, err := readTextPayload(c, cfg.MaxFileSize)
			if err != nil {
				log.Println("Error reading text payload:", err)
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Text size exceeds the maximum allowed limit"})
				return
			}
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
			fileData, fileName, err := readFilePayload(c, cfg.MaxFileSize)
			if err != nil {
				log.Println("Error retrieving file from request:", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from form data"})
				return
			}
			if int64(len(fileData)) > cfg.MaxFileSize {
				log.Printf("Error: File size (%d bytes) exceeds maximum allowed size (%d bytes)\n", len(fileData), cfg.MaxFileSize)
				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds the maximum allowed limit"})
				return
			}

			log.Println("Received file:", fileName, "Size:", len(fileData))

			enc, err := security.EncryptData(fileData, key)
			if err != nil {
				log.Println("Error encrypting file data:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt file data"})
				return
			}

			storagePath := cfg.StoragePath
			if strings.TrimSpace(storagePath) == "" {
				storagePath = os.TempDir()
			}
			if err := os.MkdirAll(storagePath, 0700); err != nil {
				log.Println("Error preparing storage path:", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare storage path"})
				return
			}

			fp := filepath.Join(storagePath, hash)
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
				FileName:  sanitizeFilename(fileName),
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func readTextPayload(c *gin.Context, maxSize int64) (string, error) {
	if text, ok := c.GetPostForm("data"); ok {
		return text, nil
	}

	limited := io.LimitReader(c.Request.Body, maxSize+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return "", err
	}
	if int64(len(raw)) > maxSize {
		return "", fmt.Errorf("text payload exceeds max size")
	}

	return string(raw), nil
}

func readFilePayload(c *gin.Context, maxSize int64) ([]byte, string, error) {
	if fileHeader, err := c.FormFile("file"); err == nil {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, "", err
		}
		defer file.Close()

		limited := io.LimitReader(file, maxSize+1)
		raw, err := io.ReadAll(limited)
		if err != nil {
			return nil, "", err
		}
		return raw, fileHeader.Filename, nil
	}

	limited := io.LimitReader(c.Request.Body, maxSize+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", err
	}
	if len(raw) == 0 {
		return nil, "", fmt.Errorf("empty file payload")
	}

	fileName := firstNonEmpty(c.Query("filename"), c.GetHeader("X-Filename"))
	if fileName == "" {
		fileName = "upload.bin"
	}

	return raw, fileName, nil
}

func sanitizeFilename(name string) string {
	base := filepath.Base(strings.TrimSpace(name))
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "upload.bin"
	}
	return base
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
