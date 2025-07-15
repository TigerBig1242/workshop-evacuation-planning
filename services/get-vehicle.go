package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// GetVehicle retrieves a vehicle by its ID from the database.
func GetVehicle(id uint) (models.Vehicle, error) {
	var vehicle models.Vehicle
	result := config.DB.First(&vehicle, id)
	if result.Error != nil {
		return models.Vehicle{}, result.Error
	}

	return vehicle, nil
}

// GetAllVehicles retrieves all vehicles from the database.
func GetAllVehicles() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	result := config.DB.Find(&vehicles)
	if result.Error != nil {
		return nil, result.Error
	}

	return vehicles, nil
}
