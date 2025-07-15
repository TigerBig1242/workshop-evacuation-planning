package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
)

// Get evacuation status from the database
func GetEvacuationStatus() ([]models.Evacuation_status, error) {
	var evacuationStatus []models.Evacuation_status
	result := config.DB.Find(&evacuationStatus)
	if result.Error != nil {
		return nil, result.Error
	}

	return evacuationStatus, nil
}

// func GetEvacuationPlans() ([]map[uint]interface{}, error) {
func GetEvacuationPlans() ([]models.Evacuation_plan, error) {
	// var evacuationPlans []map[uint]interface{}
	var evacuationPlans []models.Evacuation_plan
	result := config.DB.Table("evacuation_plans").Select("id, zone_id, vehicle_id, people_evacuated, estimated_time_arrive").
		Scan(&evacuationPlans)
	if result.Error != nil {
		return nil, result.Error
	}
	return evacuationPlans, nil
}

// func GetEvacuationZones() ([]map[uint]interface{}, error) {
func GetEvacuationZones() ([]models.Evacuation_zone, error) {
	// var evacuationZones []map[uint]interface{}
	var evacuationZones []models.Evacuation_zone
	// result := config.DB.Table("evacuation_zones").Select("id, number_of_people, urgency_level, latitude, longitude").
	// 	Scan(&evacuationZones)
	result := config.DB.Find(&evacuationZones)
	if result.Error != nil {
		return nil, result.Error
	}
	return evacuationZones, nil
}
