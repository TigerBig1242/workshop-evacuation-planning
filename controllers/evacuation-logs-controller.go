package controllers

import (
	// "time"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/services"
	// "go.uber.org/zap"
)

type EvacuationLogController struct {
	LogService *services.EvacuationLogService
}

func (c *EvacuationLogController) LoggingEvacuationController(ctx *fiber.Ctx) error {
	return nil
}

// allPlans := GetEvacuationPlans()

// Not complete yet. ยังไม่เสร็จเพราะยังไม่รู้ว่าจะทำยังไง ยังต้องหาวิธีการทำต่อให้เสร็จ
