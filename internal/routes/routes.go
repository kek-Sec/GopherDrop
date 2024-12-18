package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/kek-Sec/gopherdrop/internal/handlers"
	"github.com/kek-Sec/gopherdrop/internal/config"
)

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

	r.POST("/send", handlers.CreateSend(cfg, db))
	r.GET("/send/:id", handlers.GetSend(cfg, db))
	r.GET("/send/:id/check", handlers.CheckPasswordProtection(db)) // New route for checking password protection

	return r
}
