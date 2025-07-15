package models

import (
	"gorm.io/gorm"
)

type Evacuation_status struct {
	gorm.Model
	Zone_id           int    `json:"zone_id"`
	Total_evacuated   int    `json:"total_evacuated"`
	Remaining_people  int    `json:"remaining_people"`
	Last_vehicle_used string `json:"last_vehicle_used"`
}
