// Package main starts the GopherDrop server.
package main

import (
	"log"
	"net/http"

	"github.com/kek-Sec/gopherdrop/internal/config"
	"github.com/kek-Sec/gopherdrop/internal/database"
	"github.com/kek-Sec/gopherdrop/internal/models"
	"github.com/kek-Sec/gopherdrop/internal/routes"
)

var version = "dev"

// main initializes configuration, database, and the HTTP server.
func main() {
	log.Printf("Starting GopherDrop: %s", version)
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	db.AutoMigrate(&models.Send{})
	go database.CleanupExpired(db)
	r := routes.SetupRouter(cfg, db)
	log.Fatal(http.ListenAndServe(cfg.ListenAddr, r))
}
