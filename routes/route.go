package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/controllers"
	"github.com/tigerbig1242/evacuation-planning/logger"
	"go.uber.org/zap"

	// "github.com/tigerbig1242/evacuation-planning/logger"
	"github.com/tigerbig1242/evacuation-planning/middleware"
	// "go.uber.org/zap"/
)

func SetRoutes() *fiber.App {

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.Log.Error("Application error",
				zap.Error(err),
				zap.String("url", c.OriginalURL()),
				zap.String("method", c.Method()),
			)

			code := fiber.StatusInternalServerError

			fiberError, ok := err.(*fiber.Error)
			if ok {
				code = fiberError.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error0": err.Error(),
			})
		},
	})
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello Go Fiber")
	})

	app.Post("/create-evacuation-zone", controllers.AddEvacuationZone)
	app.Post("/create-evacuation-vehicle", controllers.AddEvacuationVehicle)
	app.Get("/get-evacuation-zones", controllers.GetEvacuationZones)
	app.Get("/get-evacuation-zone/:id", controllers.GetEvacuationZone)
	app.Get("/get-evacuation-vehicles", controllers.GetEvacuationVehicles)
	app.Get("/get-evacuation-vehicle/:id", controllers.GetEvacuationVehicle)
	// app.Post("/create-evacuation-plans", controllers.GenerateEvacuationPlans)
	// app.Post("/create-evacuation-plan", controllers.CreateEvacuationPlan)
	// app.Get("/get-evacuation-plans", controllers.GetEvacuationPlans)
	// app.Get("/get-evacuation-plan/:id", controllers.GetEvacuationPlan)
	app.Get("/get-zone-urgency", controllers.GetZoneByUrgency)
	app.Get("/get-evacuation-status", controllers.EvacuationStatus)
	// app.Put("/update-evacuation-plan/:id", controllers.UpdateSinglePlan)
	// app.Put("/update-evacuation-plans", controllers.UpdatePlans)

	// app.Delete("/delete-evacuation-plan/:id", controllers.DeletePlan) // Delete single plan
	app.Delete("/delete-evacuation-plans", controllers.DeletePlans) // Delete all current plans

	app.Get("/cause-error", func(c *fiber.Ctx) error {
		return fmt.Errorf("intentional error for testing")
	})
	app.Use(middleware.LoggingMiddleware)
	app.Post("/create-evacuation-plans", controllers.GenerateEvacuationPlans) // Generate multi evacuation plans
	app.Post("/create-evacuation-plan", controllers.CreateEvacuationPlan)     // Generate single plan
	app.Get("/get-evacuation-plans", controllers.GetEvacuationPlans)          // Get multi plans
	app.Get("/get-evacuation-plan/:id", controllers.GetEvacuationPlan)        // Get single plan
	app.Put("/update-evacuation-plan/:id", controllers.UpdateSinglePlan)      // Update single plan
	app.Put("/update-evacuation-plans", controllers.UpdatePlans)              // Update multi plans
	app.Delete("/delete-evacuation-plan/:id", controllers.DeletePlan)         // Delete single plan
	app.Delete("/delete-evacuation-plans", controllers.DeletePlans)           // Delete all current plans

	return app
}
