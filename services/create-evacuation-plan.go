package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// Create evacuation plan a new evacuation plan in the database.
func CreateEvacuationPlan(evacuationPlan *models.Evacuation_plan) (models.Evacuation_plan, error) {
	result := config.DB.Create(&evacuationPlan)
	if result.Error != nil {
		return models.Evacuation_plan{}, result.Error
	}
	plan := config.DB.Preload("Vehicle").Preload("Zone").First(evacuationPlan, evacuationPlan.ID)
	if plan.Error != nil {
		return models.Evacuation_plan{}, plan.Error
	}
	return *evacuationPlan, nil
}
