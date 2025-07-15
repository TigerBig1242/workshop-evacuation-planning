package services

import (
	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/logger"
	"github.com/tigerbig1242/evacuation-planning/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EvacuationLogService struct {
	DB *gorm.DB
}

// Create EvacuationLog creates a new evacuation log in the database.
func NewEvacuationLogService(db *gorm.DB) *EvacuationLogService {
	return &EvacuationLogService{DB: db}
}

// Create log for vehicle assignment
func (s *EvacuationLogService) LogVehicleAssignment(log *models.EvacuationLog) error {
	result := config.DB.Create(log)

	if result.Error != nil {
		logger.Log.Error("Failed to create evacuation log",
			zap.Int("operation_id", log.Operation_id),
			zap.Int("vehicle_id", log.Vehicle_id),
			zap.Error(result.Error))
		return result.Error
	}

	logger.Log.Info("Vehicle assignment for evacuation operation logged successfully",
		zap.Int("operation_id", log.Operation_id),
		zap.Int("vehicle_id", log.Vehicle_id),
		zap.Time("eta", log.EstimatedETA))

	return nil
}
