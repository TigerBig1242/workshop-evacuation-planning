package models

import (
	"time"

	"gorm.io/gorm"
)

type EvacuationLog struct {
	gorm.Model
	Operation_id    int       `json:"operation_id"`
	Vehicle_id      int       `json:"vehicle_id"`
	Origin          string    `json:"origin"`
	Destination     string    `json:"destination"`
	EstimatedETA    time.Time `json:"estimated_time_arrive"`
	ActualETA       time.Time `json:"actual_time_arrive"`
	Status          string    `json:"status"`
	PeopleEvacuated int       `json:"people_evacuated"`
	Distance        float64   `json:"distance"`
	RequestPayload  string    `json:"request_payload"`
}
