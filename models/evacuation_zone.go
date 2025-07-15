package models

import (
	"gorm.io/gorm"
)

type Evacuation_zone struct {
	gorm.Model
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Number_of_people int     `json:"number_of_people"`
	Urgency_level    int     `json:"urgency_level"`
}
