package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/trikunc/go_vibe/config"
	"github.com/trikunc/go_vibe/data"
	"github.com/trikunc/go_vibe/models"
	"github.com/trikunc/go_vibe/routes"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] No .env file found or using system environment variables")
	}

	// Initialize Database connection
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("[FATAL] Database connection failed: %v", err)
	}

	// Auto-migrate schema
	log.Println("[INFO] Running database auto-migration...")
	if err := db.AutoMigrate(&models.Protein{}, &models.User{}); err != nil {
		log.Fatalf("[FATAL] AutoMigrate failed: %v", err)
	}

	// Seed scientific protein data if empty
	if err := data.SeedProteins(db); err != nil {
		log.Printf("[WARN] Seeder warning: %v", err)
	}

	// Setup Gin router and routes
	router := routes.SetupRouter(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[INFO] Server Protein Pangan running on http://localhost:%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Server failed to start: %v", err)
	}
}
