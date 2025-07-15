package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// CreateEvacuationZone creates a new evacuation zone in the database.
func CreateEvacuationZone(zone *models.Evacuation_zone) (models.Evacuation_zone, error) {
	result := config.DB.Create(&zone)
	if result.Error != nil {
		return models.Evacuation_zone{}, result.Error
	}

	return *zone, nil
}
