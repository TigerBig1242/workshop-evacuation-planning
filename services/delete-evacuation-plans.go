package services

import (
	"fmt"

	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/models"
	"gorm.io/gorm"
)

func DeleteEvacuationPlan(id uint) error {
	var deletePlan models.Evacuation_plan
	query := config.DB.Delete(&deletePlan, id)
	if query.Error != nil {
		return query.Error
	}

	return nil
}

func DeleteEvacuationPlans() error {
	transaction := config.DB.Begin()
	var count int64
	if err := transaction.Model(&models.Evacuation_plan{}).Count(&count).Error; err != nil {
		transaction.Rollback()
		return fmt.Errorf("error counting records: %v", err)
	}
	// if err.Error != nil {
	// 	return fmt.Errorf("error counting records: %v", err)
	// }

	if count == 0 {
		transaction.Rollback()
		return fmt.Errorf("no records found to delete")
	}
	if err := transaction.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Evacuation_plan{}).Error; err != nil {
		transaction.Rollback()
		return fmt.Errorf("error deleting all records: %v", err)
	}

	if err := transaction.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	// if query.Error != nil {
	// 	return fmt.Errorf("error deleting all records: %v", err)
	// }
	// if query.Error != nil {
	// 	return query.Error
	// }

	return nil
}
