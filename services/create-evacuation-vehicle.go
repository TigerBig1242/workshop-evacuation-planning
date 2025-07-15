package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// CreateEvacuationVehicle creates a new evacuation vehicle in the database.
func CreateEvacuationVehicle(vehicle *models.Vehicle) (models.Vehicle, error) {
	result := config.DB.Create(&vehicle)
	if result.Error != nil {
		return models.Vehicle{}, result.Error
	}

	return *vehicle, nil
}
