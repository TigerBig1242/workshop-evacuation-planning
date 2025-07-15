package services

import (
	"fmt"
	"time"

	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

type PlanResponse struct {
	ID          uint      `json:"id"`
	ZoneID      int       `json:"zone_id"`
	VehicleID   int       `json:"vehicle_id"`
	ETA         time.Time `json:"estimated_time_arrive"`
	VehicleType string    `json:"vehicle_type"`
	Capacity    int       `json:"capacity"`
}

// Get evacuation plan by ID
func GetEvacuationPlan(id uint) (map[string]interface{}, error) {
	// var results map[string]interface{}
	var results models.Evacuation_plan
	query := config.DB.Model(models.Evacuation_plan{}).
		Select(`evacuation_plans.id, evacuation_plans.zone_id, evacuation_plans.vehicle_id,
			evacuation_plans.estimated_time_arrive, evacuation_plans.people_evacuated,
			vehicles.capacity, vehicles.vehicle_type
		`).
		Joins("JOIN vehicles ON evacuation_plans.vehicle_id = vehicles.id").
		Preload("Vehicle").
		Where("evacuation_plans.id = ?", id).First(&results)

	if query.Error != nil {
		return nil, query.Error
	}

	// Check ensure results is not empty
	if query.RowsAffected == 0 {
		return nil, fmt.Errorf("no evacuation plan found with ID %d", id)
	}

	response := map[string]interface{}{
		"id":                    results.ID,
		"zone_id":               results.Zone_id,
		"vehicle_id":            results.Vehicle_id,
		"estimated_time_arrive": results.Estimated_time_arrive,
		"people_evacuated":      results.People_evacuated,
		"vehicle_details": map[string]interface{}{
			"vehicle_type": results.Vehicle.Vehicle_type,
			"capacity":     results.Vehicle.Capacity,
		},
	}

	return response, nil
}

// Get evacuation all plans
func GetAllEvacuationPlans() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	// var plans []models.Evacuation_plan
	query := config.DB.Model(&models.Evacuation_plan{}).
		Select(`evacuation_plans.id, evacuation_plans.zone_id, evacuation_plans.vehicle_id, 
			evacuation_plans.estimated_time_arrive, 
			evacuation_plans.people_evacuated,
			vehicles.vehicle_type,
			vehicles.capacity
		`).
		Joins("JOIN vehicles ON evacuation_plans.vehicle_id = vehicles.id").
		Find(&results)

	if query.Error != nil {
		return nil, query.Error
	}

	return results, nil
}

func UpdatePlan(id uint, evacuees int) (*models.Evacuation_plan, error) {
	var existPlan models.Evacuation_plan
	if err := config.DB.First(&existPlan, id).Error; err != nil {
		return nil, fmt.Errorf("error finding evacuation plan: %v", err)
	}

	if err := config.DB.Model(&existPlan).
		Update("people_evacuated", evacuees).Error; err != nil {
		return nil, fmt.Errorf("error updating plan: %v", err)
	}

	if err := config.DB.Preload("Vehicle").
		First(&existPlan, id).Error; err != nil {
		return nil, fmt.Errorf("error retrieving updated plan: %v", err)
	}
	return &existPlan, nil
}

// func UpdatePlans() (*models.Evacuation_plan, error) {
// 	var getPlans models.Evacuation_plan
// 	err := config.DB.First(&getPlans)
// 	if err != nil {
// 		return nil, fmt.Errorf("error finding evacuation plans %v", err)
// 	}

//		return &getPlans, nil
//	}
func UpdatePlans(updates []map[string]interface{}) ([]map[string]interface{}, error) {
	transaction := config.DB.Begin()
	var updatePlans []map[string]interface{}

	for _, update := range updates {
		planID, ok := update["id"].(float64)
		if !ok {
			transaction.Rollback()
			return nil, fmt.Errorf("invalid plan id format")
		}

		peopleEvacuated, ok := update["people_evacuated"].(float64)
		if !ok {
			transaction.Rollback()
			return nil, fmt.Errorf("invalid people evacuated format")
		}

		result := transaction.Exec(`
		UPDATE evacuation_plans SET people_evacuated = ?
		WHERE id = ?`, int(peopleEvacuated), uint(planID))

		if result.Error != nil {
			transaction.Rollback()
			return nil, fmt.Errorf("failed to update plans %v: %v", planID, result.Error)
		}

		var updatePlan map[string]interface{}
		if err := transaction.Raw(`
			SELECT 
				evacuation_plans.id, 
				evacuation_plans.zone_id, 
				evacuation_plans.vehicle_id,
				evacuation_plans.people_evacuated, 
				evacuation_plans.estimated_time_arrive,
				vehicles.vehicle_type, 
				vehicles.capacity 
			FROM evacuation_plans
			JOIN vehicles ON evacuation_plans.vehicle_id = vehicles.id
			WHERE evacuation_plans.id = ?`, uint(planID)).Scan(&updatePlan).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}
		updatePlans = append(updatePlans, updatePlan)
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit update: %v", err)
	}

	return updatePlans, nil
}
