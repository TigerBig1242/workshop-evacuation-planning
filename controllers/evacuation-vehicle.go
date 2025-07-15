package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
)

type EvacuationVehicleResponse struct {
	ID           uint    `json:"id"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Capacity     int     `json:"capacity"`
	Vehicle_type string  `json:"vehicle_type"`
	Speed        int     `json:"speed"`
	CreatedDate  string  `json:"createdAt"`
}

// Create handles the creation of a new evacuation zone.
func AddEvacuationVehicle(c *fiber.Ctx) error {
	evacuationVehicleRequest := new(models.Vehicle)

	// check if the request body is empty
	err := c.BodyParser(evacuationVehicleRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error not parsing body evacuation vehicle",
		})
	}

	// check if the request body fields is empty
	if evacuationVehicleRequest.Capacity <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Capacity of vehicle must be greater than 0 and must be number",
		})
	} else if evacuationVehicleRequest.Vehicle_type == "" {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Vehicle type must be empty",
		})
	} else if evacuationVehicleRequest.Latitude <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Latitude of vehicle must be greater than 0 and must be number",
		})
	} else if evacuationVehicleRequest.Longitude <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Longitude of vehicle must be greater than 0 and must be number",
		})
	} else if evacuationVehicleRequest.Speed <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Speed of vehicle must be greater than 0 and must be number",
		})
	}

	// Check creating evacuation vehicle
	evacuationVehicle, err := services.CreateEvacuationVehicle(evacuationVehicleRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error not creating evacuation vehicle",
		})
	}

	response := EvacuationVehicleResponse{
		ID:           evacuationVehicle.ID,
		Latitude:     evacuationVehicle.Latitude,
		Longitude:    evacuationVehicle.Longitude,
		Capacity:     evacuationVehicle.Capacity,
		Vehicle_type: evacuationVehicle.Vehicle_type,
		Speed:        evacuationVehicle.Speed,
		CreatedDate:  ThaiTimeFormat(evacuationVehicle.CreatedAt),
	}

	return c.Status(201).JSON(fiber.Map{
		"Message":      "Created evacuation vehicle successfully",
		"vehicle data": response,
	})
}

// Get evacuation vehicles retrieves all evacuation vehicles from the database.
func GetEvacuationVehicles(c *fiber.Ctx) error {
	vehicle, err := services.GetAllVehicles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving evacuation vehicles",
		})
	}

	var vehiclesData []EvacuationVehicleResponse
	for i := range vehicle {
		vehiclesData = append(vehiclesData, EvacuationVehicleResponse{
			ID:           vehicle[i].ID,
			Latitude:     vehicle[i].Latitude,
			Longitude:    vehicle[i].Longitude,
			Capacity:     vehicle[i].Capacity,
			Vehicle_type: vehicle[i].Vehicle_type,
			Speed:        vehicle[i].Speed,
			CreatedDate:  ThaiTimeFormat(vehicle[i].CreatedAt),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"Message": "Get all evacuation vehicles successfully",
		"Vehicle": vehiclesData,
	})
}

// GetEvacuationVehicle retrieves a vehicle by its ID from the database.
func GetEvacuationVehicle(c *fiber.Ctx) error {
	idParams, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Error not parsing id parameter",
		})
	}

	vehicle, err := services.GetVehicle(uint(idParams))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Vehicle evacuation not found",
		})
	}

	response := EvacuationVehicleResponse{
		ID:           vehicle.ID,
		Latitude:     vehicle.Latitude,
		Longitude:    vehicle.Longitude,
		Capacity:     vehicle.Capacity,
		Vehicle_type: vehicle.Vehicle_type,
		Speed:        vehicle.Speed,
		CreatedDate:  ThaiTimeFormat(vehicle.CreatedAt),
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":      "Get evacuation vehicle successfully",
		"vehicle data": response,
	})
}
