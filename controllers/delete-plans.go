package controllers

import (
	// "fmt"

	"github.com/gofiber/fiber/v2"
	// "github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
)

func DeletePlan(c *fiber.Ctx) error {
	idParams, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Invalid plan ID",
			"error":   err.Error(),
		})
	}

	err = services.DeleteEvacuationPlan(uint(idParams))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Error plan not found",
			// "error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "deleted plan success",
		// "delete plan": plan,
	})
}

func DeletePlans(c *fiber.Ctx) error {
	// plans := services.DeleteEvacuationPlans()

	// if plans != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": "Failed to delete all plans",
	// 	})
	// }

	if err := services.DeleteEvacuationPlans(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete all plans",
			"error":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Deleted all current evacuation plans success",
		// "delete plans": plans,
	})
}
