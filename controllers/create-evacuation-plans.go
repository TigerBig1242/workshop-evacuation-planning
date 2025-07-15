package controllers

import (
	"fmt"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
	"github.com/tigerbig1242/evacuation-planning/utils"
)

type VehicleCapacity struct {
	ID       uint
	Capacity int
	Type     string
}

func CreateEvacuationPlan(c *fiber.Ctx) error {

	getZones, err := services.GetEvacuationZones()
	if err != nil {
		fmt.Println("Error retrieving zone urgency level")
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving zone urgency level",
		})
	}

	getVehicles, err := services.GetAllVehicles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving vehicles",
		})
	}

	urgencyLevel := 1
	for _, zone := range getZones {
		if zone.Urgency_level >= urgencyLevel {
			urgencyLevel = zone.Urgency_level
		}
	}

	// Filter find zone follow urgency level
	urgentZones := filterZoneUrgencyLevel(getZones, urgencyLevel)
	if len(urgentZones) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No zones found",
		})
	}
	// Sort amount people zone follow urgency
	// and choose most urgency with most amount people of that zone
	urgentZones = SortPeople(urgentZones)
	selectedZone := urgentZones[0]

	// Sort vehicle capacity from most to least
	for i := 0; i < len(getVehicles)-1; i++ {
		for j := 0; j < len(getVehicles)-i-1; j++ {
			if getVehicles[j].Capacity < getVehicles[j+1].Capacity {
				getVehicles[j], getVehicles[j+1] = getVehicles[j+1], getVehicles[j]
			}
		}
	}

	var selectedVehicle *models.Vehicle
	minCapacity := math.MaxInt32

	// Loop for find vehicle capacity fit amount evacuees
	// and find vehicle capacity that has nearby amount evacuees
	for i, vehicle := range getVehicles {
		if vehicle.Capacity <= 0 {
			continue
		}

		if vehicle.Capacity == selectedZone.Number_of_people {
			selectedVehicle = &getVehicles[i]
			break
		}
		fitCapacity := math.Abs(float64(vehicle.Capacity) - float64(selectedZone.Number_of_people))
		if int(fitCapacity) < int(minCapacity) {
			minCapacity = int(fitCapacity)
			selectedVehicle = &getVehicles[i]
		}
	}
	if selectedVehicle != nil {
		distance := int(utils.HaversineFormula(selectedVehicle.Latitude, selectedVehicle.Longitude,
			selectedZone.Latitude, selectedZone.Longitude))

		evacuationPlan := models.Evacuation_plan{
			Zone_id:               int(selectedZone.ID),
			Vehicle_id:            int(selectedVehicle.ID),
			Estimated_time_arrive: CalculateETAS(float64(distance), float64(selectedVehicle.Speed)),
			// People_evacuated:      selectedVehicle.Capacity,
		}

		plan, err := services.CreateEvacuationPlan(&evacuationPlan)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "error created evacuation plan",
			})
		}

		// Time format
		timeInfo := map[string]string{
			"date_time": ThaiTimeFormat(plan.Estimated_time_arrive),
			"ETA":       FormatRemainingTime(plan.Estimated_time_arrive),
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Create plan success",
			"details": fiber.Map{
				"plan ID":          plan.ID,
				"zone ID":          selectedZone.ID,
				"vehicle ID":       selectedVehicle.ID,
				"zone urgency":     selectedZone.Urgency_level,
				"people evacuees":  selectedZone.Number_of_people,
				"vehicle capacity": selectedVehicle.Capacity,
				"vehicle type":     selectedVehicle.Vehicle_type,
				"distance":         distance,
				"time Info":        timeInfo,
			},
		})
	}

	return c.Status(404).JSON(fiber.Map{
		"message": "No suitable vehicle found",
	})
}

// Filter zones by urgency level
func filterZoneUrgencyLevel(zones []models.Evacuation_zone, urgencyLevel int) []models.Evacuation_zone {
	var result []models.Evacuation_zone
	for _, zone := range zones {
		if zone.Urgency_level == urgencyLevel {
			result = append(result, zone)
		}
	}
	return result
}

// Sort people in zones by most to least
func SortPeople(zones []models.Evacuation_zone) []models.Evacuation_zone {
	amount := len(zones)
	for i := 0; i < amount-1; i++ {
		for j := 0; j < amount-i-1; j++ {
			if zones[j].Number_of_people < zones[j+1].Number_of_people {
				zones[j], zones[j+1] = zones[j+1], zones[j]
			}
		}
	}
	return zones
}

func UpdateSinglePlan(c *fiber.Ctx) error {
	idParams, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Invalid plan ID",
			"error":   err.Error(),
		})
	}
	fmt.Printf("ID Params %d:", idParams)
	// Get single evacuation plan by id
	getPlan, err := services.GetEvacuationPlan(uint(idParams))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Error plan not found",
			"error":   err.Error(),
		})
	}
	fmt.Printf("Get Plan : %v\n", getPlan)

	// Parse update people evacuees from request body
	var updateData models.Evacuation_plan
	errBody := c.BodyParser(&updateData)
	if errBody != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   errBody.Error(),
		})

	}

	if updateData.People_evacuated < 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "People evacuated cannot be negative",
		})
	}

	vehicleDetails, ok := getPlan["vehicle_details"].(map[string]interface{})
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"message": "Invalid vehicle details format",
		})
	}

	capacity, ok := vehicleDetails["capacity"].(int)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"message": "Invalid capacity format",
		})
	}

	if updateData.People_evacuated > capacity {
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Sprintf("People evacuated (%d) exceeds vehicle capacity (%d)",
				updateData.People_evacuated, capacity),
		})
	}

	// Update only the fields that are provided and id matched
	updateData.ID = uint(idParams)
	fmt.Println("updateData ID", updateData.ID)

	updatePlan, err := services.UpdatePlan(updateData.ID, updateData.People_evacuated)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed updated evacuation plan",
			"error":   err.Error(),
		})
	}

	// Format time
	timeInfo := map[string]string{
		"date time": ThaiTimeFormat(updatePlan.Estimated_time_arrive),
		"ETA":       FormatRemainingTime(updatePlan.Estimated_time_arrive),
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Update evacuees success",
		"update plan": fiber.Map{
			"id":               updatePlan.ID,
			"zone id":          updatePlan.Zone_id,
			"vehicle id":       updatePlan.Vehicle_id,
			"people evacuated": updatePlan.People_evacuated,
			"ETA":              timeInfo,
			"vehicle details": fiber.Map{
				"type":     updatePlan.Vehicle.Vehicle_type,
				"capacity": updatePlan.Vehicle.Capacity,
			},
		},
	})
}
