package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// GetZone retrieves an evacuation zone by its ID from the database.
func GetZone(id uint) (models.Evacuation_zone, error) {
	var zone models.Evacuation_zone
	result := config.DB.First(&zone, id)
	if result.Error != nil {
		return models.Evacuation_zone{}, result.Error
	}

	return zone, nil
}

// GetAllZones retrieves all evacuation zones from the database.
func GetAllZones() ([]models.Evacuation_zone, error) {
	var zones []models.Evacuation_zone
	result := config.DB.Order("id ASC").Find(&zones)
	if result.Error != nil {
		return nil, result.Error
	}

	return zones, nil
}

// Get zone by urgency level
func GetZoneUrgencyLevel() ([]models.Evacuation_zone, error) {
	var zoneUrgency []models.Evacuation_zone
	result := config.DB.Order("urgency_level DESC").Find(&zoneUrgency)
	if result.Error != nil {
		return nil, result.Error
	}

	return zoneUrgency, nil
}
