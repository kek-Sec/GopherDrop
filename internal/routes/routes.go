package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/time/rate"

	"github.com/kek-Sec/gopherdrop/internal/handlers"
	"github.com/kek-Sec/gopherdrop/internal/config"
)

// Define a rate limiter with 1 request per second and a burst of 5.
var limiter = rate.NewLimiter(1, 6)

// rateLimiterMiddleware applies rate limiting to the endpoint.
func rateLimiterMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

// SetupRouter initializes the Gin router with routes and middleware.
func SetupRouter(cfg config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Enable CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
	}))

	// Apply rate limiting only to the POST /send endpoint
	r.POST("/send", rateLimiterMiddleware, handlers.CreateSend(cfg, db))

	// Other routes without rate limiting
	r.GET("/send/:id", handlers.GetSend(cfg, db))
	r.GET("/send/:id/check", handlers.CheckPasswordProtection(db))

	return r
}
