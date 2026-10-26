package data

import (
	_ "embed"
	"encoding/json"
	"log"

	"github.com/trikunc/go_vibe/models"
	"gorm.io/gorm"
)

//go:embed seed_proteins.json
var SeedProteinsJSON []byte

// SeedProteins populates the database with initial scientific food protein data if the table is empty
func SeedProteins(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Protein{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Printf("[INFO] Database already contains %d proteins. Skipping seeder.\n", count)
		return nil
	}

	var proteins []models.Protein
	if err := json.Unmarshal(SeedProteinsJSON, &proteins); err != nil {
		return err
	}

	if err := db.Create(&proteins).Error; err != nil {
		return err
	}

	log.Printf("[INFO] Successfully seeded %d food proteins into the database.\n", len(proteins))
	return nil
}
