package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
)

type EvacuationZoneResponse struct {
	ID             uint    `json:"id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	NumberOfPeople int     `json:"number_of_people"`
	UrgencyLevel   int     `json:"urgency_level"`
	CreatedDate    string  `json:"createdAt"`
}

// AddEvacuationZone handles the creation of a new evacuation zone.
func AddEvacuationZone(c *fiber.Ctx) error {
	evacuationZoneRequest := new(models.Evacuation_zone)

	// check if the request body is empty
	err := c.BodyParser(evacuationZoneRequest)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error not parsing body evacuation zone",
		})
	}

	// check if the request body fields is empty
	if evacuationZoneRequest.Number_of_people <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Number of people must be greater than 0",
		})
	} else if evacuationZoneRequest.Urgency_level <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Urgency level must be greater than 0",
		})
	} else if evacuationZoneRequest.Latitude <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Latitude must be greater than 0",
		})
	} else if evacuationZoneRequest.Longitude <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Longitude must be greater than 0",
		})
	}

	// check creating evacuation zone
	evacuationZone, err := services.CreateEvacuationZone(evacuationZoneRequest)
	// evacuationZone, err := services.CreateEvacuationZone(req)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error not creating evacuation zone",
		})
	}

	response := EvacuationZoneResponse{
		ID:             evacuationZone.ID,
		Latitude:       evacuationZone.Latitude,
		Longitude:      evacuationZone.Longitude,
		NumberOfPeople: evacuationZone.Number_of_people,
		UrgencyLevel:   evacuationZone.Urgency_level,
		CreatedDate:    ThaiTimeFormat(evacuationZone.CreatedAt),
	}

	return c.Status(201).JSON(fiber.Map{
		"Message":   "Created evacuation zone successfully",
		"zone data": response,
	})
}

// Get evacuation zones retrieves all evacuation zones from the database.
func GetEvacuationZones(c *fiber.Ctx) error {
	zones, err := services.GetAllZones()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving evacuation zones",
		})
	}

	var zonesData []EvacuationZoneResponse
	for i := range zones {
		zonesData = append(zonesData, EvacuationZoneResponse{
			ID:             zones[i].ID,
			Latitude:       zones[i].Latitude,
			Longitude:      zones[i].Longitude,
			NumberOfPeople: zones[i].Number_of_people,
			UrgencyLevel:   zones[i].Urgency_level,
			CreatedDate:    ThaiTimeFormat(zones[i].CreatedAt),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Message": "Get all evacuation zones successfully",
		"Zones":   zonesData,
	})
}

// Get evacuation zone retrieves a specific evacuation zone by its ID.
func GetEvacuationZone(c *fiber.Ctx) error {
	idParam, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error not parsing id parameter",
		})
	}

	zone, err := services.GetZone(uint(idParam))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Zone evacuation not found",
		})
	}

	responseZone := EvacuationZoneResponse{
		ID:             zone.ID,
		Latitude:       zone.Latitude,
		Longitude:      zone.Longitude,
		NumberOfPeople: zone.Number_of_people,
		UrgencyLevel:   zone.Urgency_level,
		CreatedDate:    ThaiTimeFormat(zone.CreatedAt),
	}

	return c.Status(200).JSON(fiber.Map{
		"Message": "Get evacuation zone successfully",
		"Zone":    responseZone,
	})
}
