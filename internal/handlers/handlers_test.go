package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/models"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Send{})
	return db
}

func setupTestRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/send", CreateSend(cfg, db))
	r.GET("/send/:id", GetSend(cfg, db))
	r.GET("/send/:id/check", CheckPasswordProtection(db))
	return r
}

func createMultipartRequest(fieldName, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create form field
	_ = writer.WriteField("type", "text")
	_ = writer.WriteField(fieldName, content)

	// Close writer to finalize boundary
	writer.Close()
	return body, writer.FormDataContentType()
}

func createMultipartFileRequest(fieldName, filename, content string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field
	_ = writer.WriteField("type", "file")

	// Add file field
	part, _ := writer.CreateFormFile(fieldName, filename)
	io.WriteString(part, content)

	writer.Close()
	return body, writer.FormDataContentType()
}

func TestCheckPasswordProtection(t *testing.T) {
	db := setupTestDB()
	cfg := config.Config{
		SecretKey:   "supersecretkeysupersecretkey32",
		MaxFileSize: 1024 * 1024, // 1MB
	}
	r := setupTestRouter(cfg, db)

	// Create a send with a password
	sendWithPassword := models.Send{
		Hash:      "protectedhash",
		Type:      "text",
		Data:      "encryptedDataHere",
		Password:  "password123",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	db.Create(&sendWithPassword)

	// Create a send without a password
	sendWithoutPassword := models.Send{
		Hash:      "unprotectedhash",
		Type:      "text",
		Data:      "encryptedDataHere",
		Password:  "",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	db.Create(&sendWithoutPassword)

	// Test for send with password
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/send/protectedhash/check", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	expectedBody := `{"requiresPassword":true}`
	if w.Body.String() != expectedBody {
		t.Fatalf("expected body %s got %s", expectedBody, w.Body.String())
	}

	// Test for send without password
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/send/unprotectedhash/check", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	expectedBody = `{"requiresPassword":false}`
	if w.Body.String() != expectedBody {
		t.Fatalf("expected body %s got %s", expectedBody, w.Body.String())
	}
}

func TestCreateSendText(t *testing.T) {
	db := setupTestDB()
	cfg := config.Config{
		SecretKey:   "supersecretkeysupersecretkey32",
		MaxFileSize: 1024 * 1024, // 1MB
	}
	r := setupTestRouter(cfg, db)

	body, contentType := createMultipartRequest("data", "This is a test message.")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/send", body)
	req.Header.Set("Content-Type", contentType)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestCreateSendFileTooLarge(t *testing.T) {
	db := setupTestDB()
	cfg := config.Config{
		SecretKey:   "supersecretkeysupersecretkey32",
		MaxFileSize: 10, // Only allow 10 bytes
	}
	r := setupTestRouter(cfg, db)

	// Create a file with more than 10 bytes
	body, contentType := createMultipartFileRequest("file", "test.txt", "This file is too large.")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/send", body)
	req.Header.Set("Content-Type", contentType)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413 got %d", w.Code)
	}
}

func TestGetNonExistentSend(t *testing.T) {
	db := setupTestDB()
	cfg := config.Config{
		SecretKey:   "supersecretkeysupersecretkey32",
		MaxFileSize: 1024 * 1024,
	}
	r := setupTestRouter(cfg, db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/send/unknownhash", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}

func TestExpiredSend(t *testing.T) {
	db := setupTestDB()
	cfg := config.Config{
		SecretKey:   "supersecretkeysupersecretkey32",
		MaxFileSize: 1024 * 1024,
	}
	r := setupTestRouter(cfg, db)

	// Create expired send
	send := models.Send{
		Hash:      "expiredhash",
		Type:      "text",
		Data:      "encryptedDataHere",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	db.Create(&send)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/send/expiredhash", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", w.Code)
	}
}
