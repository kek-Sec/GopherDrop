package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Auto-migrate the Send model
	if err := db.AutoMigrate(&models.Send{}).Error; err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	// Insert a test record
	send := models.Send{
		Hash:      "testhash",
		Type:      "text",
		Data:      "testdata",
		Password:  "",
		OneTime:   false,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	db.Create(&send)

	return db
}

func setupTestRouter() *gin.Engine {
	limiter = rate.NewLimiter(1, 6)
	cfg := config.Config{
		SecretKey:   "supersecretkey",
		MaxFileSize: 1024 * 1024,
	}
	db := setupTestDB()
	return SetupRouter(cfg, db)
}

func TestRoutesExist(t *testing.T) {
	router := setupTestRouter()

	tests := []struct {
		method   string
		endpoint string
		payload  string
		status   int
	}{
		{"POST", "/send", "type=text&data=test", http.StatusOK},
		{"POST", "/send/text", "data=test", http.StatusOK},
		{"POST", "/send/file", "test file payload", http.StatusBadRequest},
		{"GET", "/send/testhash", "", http.StatusNotFound},
		{"GET", "/send/testhash/check", "", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.endpoint, func(t *testing.T) {
			var req *http.Request
			if tt.method == "POST" {
				req = httptest.NewRequest(tt.method, tt.endpoint, strings.NewReader(tt.payload))
				if tt.endpoint == "/send/file" {
					req.Header.Set("Content-Type", "application/octet-stream")
				} else {
					req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				}
			} else {
				req = httptest.NewRequest(tt.method, tt.endpoint, nil)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.NotEqual(t, http.StatusNotFound, w.Code, "Route %s %s should exist", tt.method, tt.endpoint)
		})
	}
}

func TestCORSHeaders(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest("OPTIONS", "/send", nil)
	req.Header.Set("Origin", "*")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code, "CORS preflight should return 204 No Content")
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestRateLimiter(t *testing.T) {
	router := setupTestRouter()

	// Define a payload for the POST request
	payload := "type=text&data=test"

	// Simulate 6 requests (the burst capacity) in quick succession.
	for i := 0; i < 6; i++ {
		req := httptest.NewRequest("POST", "/send", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed within the burst capacity", i+1)
	}

	// The 7th request should be rate limited and return a 429 status
	req := httptest.NewRequest("POST", "/send", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "7th request should be rate limited")
}
