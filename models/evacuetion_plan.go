package models

import (
	"time"

	"gorm.io/gorm"
)

type Evacuation_plan struct {
	gorm.Model
	Zone_id               int       `json:"zone_id"`
	Vehicle_id            int       `json:"vehicle_id"`
	Estimated_time_arrive time.Time `json:"estimated_time_arrive"`
	People_evacuated      int       `json:"people_evacuated"`
	Vehicle               Vehicle   `gorm:"foreignKey:Vehicle_id"`
}
