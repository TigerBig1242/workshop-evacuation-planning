package controllers

import (
	"fmt"
	"math"

	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
	"github.com/tigerbig1242/evacuation-planning/utils"
)

type RemainingZones struct {
	ZoneID          uint `json:"zone_id"`
	Urgency_level   int  `json:"urgency_level"`
	RemainingPeople uint `json:"number_of_people"`
}

func GenerateEvacuationPlans(c *fiber.Ctx) error {
	getAllZones, err := services.GetEvacuationZones()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving zone urgency level",
		})
	}

	// Get all vehicles
	vehicles, err := services.GetAllVehicles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving vehicles",
		})
	}

	// Check if there are no zones
	if len(getAllZones) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "No zones available for evacuation",
		})
	}

	// Check if there are no vehicles
	if len(vehicles) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"Message": "No vehicles available for evacuation",
		})
	}

	var createPlans []models.Evacuation_plan
	var zoneRemaining []RemainingZones

	usedVehicles := make(map[uint]bool)

	// Loop through all zones from most urgent(5) to least urgent(1)
	for urgencyLevel := 5; urgencyLevel >= 1; urgencyLevel-- {

		// Filter zones by urgency level and sort by number of people most to least
		urgentZones := FilterZoneUrgency(getAllZones, urgencyLevel)
		sort.Slice(urgentZones, func(i, j int) bool {
			return urgentZones[i].Number_of_people > urgentZones[j].Number_of_people
		})
		fmt.Println("Urgent Zones: ", len(urgentZones))

		// Loop to find evacuate from each zone
		for _, zone := range urgentZones {
			evacueesPeople := zone.Number_of_people

			// Sort vehicles by capacity most to least
			sort.Slice(vehicles, func(i, j int) bool {
				return vehicles[i].Capacity > vehicles[j].Capacity
			})

			// Loop for manage evacuation follow amount evacuees
			// and vehicles not available
			for evacueesPeople > 0 {
				var selectedVehicle *models.Vehicle
				matchedVehicle, bestFitVehicle := findBestVehicle(vehicles, evacueesPeople, usedVehicles)

				// Check available vehicle for remaining zones
				if !CheckAvailableVehicles(vehicles, usedVehicles) {
					zoneRemaining = append(zoneRemaining, RemainingZones{
						ZoneID:          zone.ID,
						Urgency_level:   urgencyLevel,
						RemainingPeople: uint(evacueesPeople),
					})
					break
				}

				// Select vehicle that has capacity fit evacuees
				// If there is vehicle capacity, select vehicle that has nearby capacity
				if matchedVehicle != nil {
					selectedVehicle = matchedVehicle
				} else if bestFitVehicle != nil {
					selectedVehicle = bestFitVehicle
				} else {
					break
				}

				// Calculate amount evacuees in each round
				peopleToEvacuate := min(evacueesPeople, selectedVehicle.Capacity)
				distance := int(utils.HaversineFormula(selectedVehicle.Latitude, selectedVehicle.Longitude,
					zone.Latitude, zone.Longitude))

				evacuationPlans := &models.Evacuation_plan{
					Zone_id:               int(zone.ID),
					Vehicle_id:            int(selectedVehicle.ID),
					Estimated_time_arrive: CalculateETA(float64(distance), float64(selectedVehicle.Speed)),
					// People_evacuated:      peopleToEvacuate,
				}

				plans, err := services.CreateEvacuationPlan(evacuationPlans)
				if err != nil {
					fmt.Printf("Error creating evacuation plan: %v\n", err)
					continue
				}
				createPlans = append(createPlans, plans)
				usedVehicles[uint(selectedVehicle.ID)] = true
				evacueesPeople -= peopleToEvacuate

				fmt.Printf("Zone %d: Using vehicle %d (capacity: %d) to evacuate %d people. Remaining: %d\n",
					zone.ID, selectedVehicle.ID, selectedVehicle.Capacity, peopleToEvacuate, evacueesPeople)
			}
		}
	}

	formatPlans := make([]fiber.Map, 0)
	for _, plans := range createPlans {
		timeInfo := map[string]string{
			"ETA":       FormatRemainingTime(plans.Estimated_time_arrive),
			"date_time": ThaiTimeFormat(plans.Estimated_time_arrive),
		}

		formatPlans = append(formatPlans, fiber.Map{
			"plan id":               plans.ID,
			"zone id":               plans.Zone_id,
			"vehicle id":            plans.Vehicle_id,
			"people_evacuated":      plans.People_evacuated,
			"vehicle capacity":      plans.Vehicle.Capacity,
			"vehicle_type":          plans.Vehicle.Vehicle_type,
			"Estimated time arrive": timeInfo,
		})
	}

	vehiclesDetails := make([]fiber.Map, 0)
	for _, vehicle := range vehicles {
		if usedVehicles[uint(vehicle.ID)] {
			vehiclesDetails = append(vehiclesDetails, fiber.Map{
				"id":       vehicle.ID,
				"type":     vehicle.Vehicle_type,
				"capacity": vehicle.Capacity,
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"Message": "Generated evacuation plan successfully",
		// "Created plans":    formatPlans,
		"Created plans":    formatPlans,
		"Vehicles Details": vehiclesDetails,
		"Don't has available vehicle for zone remaining": fiber.Map{
			"Used vehicles":   usedVehicles,
			"remaining zones": zoneRemaining,
		},
	})
}

func FilterZoneUrgency(zones []models.Evacuation_zone, urgencyLevel int) []models.Evacuation_zone {
	var result []models.Evacuation_zone
	for _, zone := range zones {
		if zone.Urgency_level == urgencyLevel {
			result = append(result, zone)
		}
	}
	return result
}

// Calculate ETA
func CalculateETAS(distance, speed float64) time.Time {

	// Calculate time in hours
	hours := distance / speed

	// Convert hours to time
	minutes := hours * 60

	// Convert minutes to duration
	travelTime := time.Duration(minutes) * time.Minute

	eta := time.Now().Add(travelTime)

	// Get current time
	return eta
}

func GetEvacuationPlans(c *fiber.Ctx) error {
	evacuationPlans, err := services.GetAllEvacuationPlans()
	// evacuationPlans, err := services.EvacuationPlans()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"Message": "Error retrieving evacuation plans not found",
		})
	}

	if evacuationPlans == nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Not found evacuation plans",
		})
	}

	vehicles, errVehicles := services.GetAllVehicles()
	if errVehicles != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Error vehicles not found",
		})
	}

	// Convert estimated time to Thai time
	for i := range evacuationPlans {
		if estimatedTime, ok := evacuationPlans[i]["estimated_time_arrive"].(time.Time); ok {
			timeInfo := map[string]string{
				"date time": ThaiTimeFormat(estimatedTime),
				"ETA":       FormatRemainingTime(estimatedTime),
			}

			evacuationPlans[i]["estimated_time_arrive"] = timeInfo
		}

		if vehicleID, ok := evacuationPlans[i]["Vehicles Details"].(uint); ok {
			for _, vehicle := range vehicles {
				if vehicleID == vehicle.ID {
					evacuationPlans[i]["vehicle_type"] = map[string]interface{}{
						"type":     vehicle.Vehicle_type,
						"capacity": vehicle.Capacity,
					}
					break
				}
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":          "Found evacuation plans",
		"Evacuation plans": evacuationPlans,
	})
}

// Get evacuation plan is single by ID
func GetEvacuationPlan(c *fiber.Ctx) error {
	idParams, errID := c.ParamsInt("id")
	if errID != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Error not parsing id parameter",
		})
	}

	evacuationPlan, errPlan := services.GetEvacuationPlan(uint(idParams))
	if errPlan != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Evacuation plan not found",
		})
	}

	// Convert estimated time to Thai time is single evacuation plan
	if estimatedTime, ok := evacuationPlan["estimated_time_arrive"].(time.Time); ok {
		timeInfo := map[string]string{
			"date_time": ThaiTimeFormat(estimatedTime),
			"ETA":       FormatRemainingTime(estimatedTime),
		}
		evacuationPlan["estimated_time_arrive"] = timeInfo
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":         "Get evacuation plan successfully",
		"Evacuation plan": evacuationPlan,
	})
}

// Format date time
func FormatDateTime(plans models.Evacuation_plan) fiber.Map {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		location = time.UTC
	}
	thaiTime := plans.Estimated_time_arrive.In(location)

	return fiber.Map{
		"id":         plans.ID,
		"zone id":    plans.Zone_id,
		"vehicle id": plans.Vehicle_id,
		"ETA": fiber.Map{
			"date time": thaiTime.Format("02/01/2006 15:04:05"),
		},
		// "people evacuated": plans.People_evacuated,
	}
}

func findBestVehicle(vehicles []models.Vehicle, evacueesPeople int, usedVehicles map[uint]bool) (*models.Vehicle, *models.Vehicle) {
	var matchedVehicle, bestVehicle *models.Vehicle
	minCapacity := math.MaxInt32

	for i, vehicle := range vehicles {
		if usedVehicles[uint(vehicle.ID)] {
			continue
		}

		switch {
		case vehicle.Capacity == evacueesPeople:
			return &vehicles[i], nil
		case evacueesPeople > vehicle.Capacity:
			if vehicle.Capacity < minCapacity {
				minCapacity = vehicle.Capacity
				bestVehicle = &vehicles[i]
			}
		case vehicle.Capacity > evacueesPeople:
			diff := vehicle.Capacity - evacueesPeople
			if diff < minCapacity {
				minCapacity = diff
				bestVehicle = &vehicles[i]
			}
		}
	}
	return matchedVehicle, bestVehicle
}

func CheckAvailableVehicles(vehicles []models.Vehicle, usedVehicles map[uint]bool) bool {
	for _, vehicle := range vehicles {
		if !usedVehicles[uint(vehicle.ID)] {
			return true
		}
	}
	return false
}

func UpdatePlans(c *fiber.Ctx) error {
	plans, err := services.GetAllEvacuationPlans()
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": "Error not found evacuation plans",
			"error":   err.Error(),
		})
	}
	fmt.Printf("found plans : %v\n", plans)

	var updateRequests []map[string]interface{}
	errBody := c.BodyParser(&updateRequests)
	if errBody != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   errBody.Error(),
		})
	}

	if len(updateRequests) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "No updates provided",
		})
	}

	for _, update := range updateRequests {
		peopleEvacuated, ok := update["people_evacuated"].(float64)
		if !ok {
			return c.Status(400).JSON(fiber.Map{
				"message": fmt.Sprintf("Invalid people evacuated format in update (%d)", int(peopleEvacuated)),
			})
		}

		if peopleEvacuated < 0 {
			return c.Status(400).JSON(fiber.Map{
				"message":       fmt.Sprintf("People evacuated cannot be negative in update %d", int(peopleEvacuated)),
				"invalid value": peopleEvacuated,
			})
		}

		vehicleCapacity := update["capacity"].(float64)
		if peopleEvacuated > vehicleCapacity {
			return c.Status(400).JSON(fiber.Map{
				"message":      fmt.Sprintf("People evacuated over vehicle capacity: '%s', which specified capacity is: %d", update["vehicle_type"], int(peopleEvacuated)),
				"max capacity": vehicleCapacity,
			})
		}
	}

	updatePlans, err := services.UpdatePlans(updateRequests)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update plans",
		})
	}

	var formattedPlans []fiber.Map
	for _, plans := range updatePlans {
		timeInfo := fiber.Map{
			"ETA":       FormatRemainingTime(plans["estimated_time_arrive"].(time.Time)),
			"date_time": ThaiTimeFormat(plans["estimated_time_arrive"].(time.Time)),
		}

		formattedPlans = append(formattedPlans, fiber.Map{
			"id":               plans["id"],
			"zone_id":          plans["zone_id"],
			"vehicle_id":       plans["vehicle_id"],
			"people_evacuated": plans["people_evacuated"],
			"capacity":         plans["capacity"],
			"vehicle_type":     plans["vehicle_type"],
			"ETA":              timeInfo,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "update plans success",
		"plans":   formattedPlans,
	})
}
