package models

import (
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	Capacity     int     `json:"capacity"`
	Vehicle_type string  `json:"vehicle_type"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Speed        int     `json:"speed"`
}
