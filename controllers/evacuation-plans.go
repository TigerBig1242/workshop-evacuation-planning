package controllers

import (
	"fmt"
	// "math"
	// "sort"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/services"
	// "github.com/tigerbig1242/evacuation-planning/utils"
)

// Create evacuation plans
// func CreateEvacuationPlan(c *fiber.Ctx) error {

// 	// Get zone urgency by level
// 	zoneUrgencyLevel, err := services.GetAllZones()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Error retrieving zone urgency level",
// 		})
// 	}

// 	// Get vehicle
// 	vehicles, err := services.GetAllVehicles()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"Message": "Error retrieving vehicle",
// 		})
// 	}

// 	// Validate zone urgency level and vehicle not found
// 	if len(zoneUrgencyLevel) == 0 || len(vehicles) == 0 {
// 		return c.Status(400).JSON(fiber.Map{
// 			"Message": "No zone urgency level or vehicle found",
// 		})
// 	}

// 	// Sort by urgency level most to least
// 	// if that zone has the same urgency level, sort by number of people most to least
// 	sort.Slice(zoneUrgencyLevel, func(i, j int) bool {
// 		if zoneUrgencyLevel[i].Urgency_level == zoneUrgencyLevel[j].Urgency_level {
// 			return zoneUrgencyLevel[i].Number_of_people > zoneUrgencyLevel[j].Number_of_people
// 		}
// 		return zoneUrgencyLevel[i].Urgency_level > zoneUrgencyLevel[j].Urgency_level
// 	})

// 	// Sort by vehicle capacity most to least with appropriate vehicle type
// 	sort.Slice(vehicles, func(i, j int) bool {
// 		return vehicles[i].Capacity > vehicles[j].Capacity
// 	})

// 	// Create evacuation plan by urgency level
// 	var createPlans []models.Evacuation_plan
// 	var lastDistance float64
// 	// var allDistance []float64
// 	var zoneVehicleDistance []map[string]interface{}
// 	var peopleRemaining int

// 	// Loop through each zone urgency level
// 	for _, zone := range zoneUrgencyLevel {
// 		remainingPeople := zone.Number_of_people
// 		// Check if the vehicle is available for evacuation
// 		if len(vehicles) == 0 {
// 			break
// 		}

// 		for remainingPeople > 0 && len(vehicles) > 0 {
// 			// for urgentLevel := 5; urgentLevel >= 1; urgentLevel-- {
// 			// urgentZones := findZoneUrgency(zoneUrgencyLevel, urgencyLevel)
// 			// Search for the vehicle with the shortest distance to the zone
// 			var bestVehicleIndex int = 0
// 			var shortestDistance float64 = math.MaxFloat64

// 			for i, vehicle := range vehicles {
// 				distance := utils.HaversineFormula(vehicle.Latitude, vehicle.Longitude, zone.Latitude, zone.Longitude)

// 				// for urgencyLevel := 5; urgencyLevel >= 1; urgencyLevel-- {
// 				// 	urgentZones := findZoneUrgency(zoneUrgencyLevel, urgencyLevel)
// 				// 	sort.Slice(urgentZones, func(i, j int) bool {
// 				// 		return urgentZones[i].Number_of_people > urgentZones[j].Number_of_people
// 				// 	})

// 				// 	if zone.Number_of_people >= vehicle.Capacity {
// 				// 		vehicleForPeople := vehicle.Capacity
// 				// 	}
// 				// }

// 				// Find the vehicle with the shortest distance
// 				if distance < shortestDistance {
// 					shortestDistance = distance
// 					bestVehicleIndex = i
// 				}
// 			}
// 			// lastDistance = shortestDistance
// 			bestVehicle := vehicles[bestVehicleIndex]
// 			capacity := bestVehicle.Capacity

// 			// Number of people evacuated
// 			peopleToEvacuate := capacity
// 			if remainingPeople < capacity {
// 				peopleToEvacuate = remainingPeople
// 			}
// 			remainingPeople -= peopleToEvacuate
// 			fmt.Println("Remaining people:", remainingPeople)
// 			// peopleToEvacuate = remainingPeople

// 			// Create new evacuation plan
// 			evacuationPlanRequest := &models.Evacuation_plan{
// 				Zone_id:               int(zone.ID),
// 				Vehicle_id:            int(vehicles[bestVehicleIndex].ID),
// 				Estimated_time_arrive: calculateETA(shortestDistance, float64(vehicles[bestVehicleIndex].Speed)),
// 				People_evacuated:      peopleToEvacuate,
// 			}

// 			// Call function CreateEvacuationPlan to create evacuation plan and save to database
// 			plan, err := services.CreateEvacuationPlan(evacuationPlanRequest)
// 			if err != nil {
// 				// If error pass through, continue to next vehicle zone
// 				continue
// 			}

// 			// Save distance between vehicle and zone
// 			zoneVehicleDistance = append(zoneVehicleDistance, map[string]interface{}{
// 				"ZoneID":    zone.ID,
// 				"vehicleID": vehicles[bestVehicleIndex].ID,
// 				"Distance":  fmt.Sprintf("%.2f km", shortestDistance),
// 			})

// 			createPlans = append(createPlans, plan)

// 			// Delate vehicles is used from listed
// 			vehicles = append(vehicles[:bestVehicleIndex], vehicles[bestVehicleIndex+1:]...)
// 			// remainingPeople -= peopleToEvacuate

// 			// Keep all distance
// 			// allDistance = append(allDistance, shortestDistance)
// 		}

// 	}

// 	if len(createPlans) == 0 {
// 		return c.Status(500).JSON(fiber.Map{
// 			"Message": "Could not create any evacuation plan",
// 		})
// 	}

// 	// Copy original vehicles
// 	// originalVehicles := make([]models.Vehicle, len(vehicles))
// 	// copy(originalVehicles, vehicles)

// 	fmt.Println("Evacuation plan found:", createPlans)
// 	fmt.Printf("Distance between evacuation zone and vehicle evacuate: %.2f km\n", lastDistance)
// 	return c.Status(201).JSON(fiber.Map{
// 		"Message":          "Evacuation plan created successfully",
// 		"Evacuation plans": createPlans,
// 		// "Get zone urgency": zoneUrgencyLevel,
// 		// "Get vehicle":      vehicles,
// 		// "All vehicles":          originalVehicles,
// 		"Haversine": lastDistance,
// 		// "Zone-Vehicle distance": zoneVehicleDistance,
// 		"Remaining people": peopleRemaining,
// 	})
// }

// Get evacuation all plans
func GetEvacuationPlanss(c *fiber.Ctx) error {
	evacuationPlans, err := services.GetAllEvacuationPlans()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving evacuation plans not found",
		})
	}

	// Convert estimated time to Thai time
	for i := range evacuationPlans {
		if estimatedTime, ok := evacuationPlans[i]["estimated_time_arrive"].(time.Time); ok {
			// thaiDate := ThaiTimeFormat(evacuationPlans[i]["estimated_time_arrive"].(time.Time))
			// evacuationPlans[i]["estimated_time_arrive"] = thaiDate

			timeInfo := map[string]string{
				"thai_format": ThaiTimeFormat(estimatedTime),
				"ETA":         FormatRemainingTime(estimatedTime),
			}

			// timeDate := evacuationPlans[i]["estimated_time_arrive"]
			evacuationPlans[i]["estimated_time_arrive"] = timeInfo
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":          "Get all evacuation plans successfully",
		"Evacuation plans": evacuationPlans,
	})
}

// Get evacuation plan is single by ID
func GetEvacuationPlanID(c *fiber.Ctx) error {
	idParams, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"MEssage": "Error not parsing id parameter",
		})
	}

	evacuationPlan, err := services.GetEvacuationPlan(uint(idParams))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"Message": "Evacuation plan not found",
		})
	}

	// Convert estimated time to Thai time is single evacuation plan
	if estimatedTime, ok := evacuationPlan["estimated_time_arrive"].(time.Time); ok {
		// thaiDate := ThaiTimeFormat(estimatedTime)
		// evacuationPlan["estimated_time_arrive"] = thaiDate,

		// timeTravel := formatRemainingTime(estimatedTime)
		// evacuationPlan["estimated_time_arrive"] = timeTravel

		timeInfo := map[string]string{
			"thai_format": ThaiTimeFormat(estimatedTime),
			"ETA":         FormatRemainingTime(estimatedTime),
		}
		evacuationPlan["estimated_time_arrive"] = timeInfo
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":         "Get evacuation plan successfully",
		"Evacuation plan": evacuationPlan,
	})
}

// Get zone urgency by level
func GetZoneByUrgency(c *fiber.Ctx) error {
	zoneUrgencyLevel, err := services.GetZoneUrgencyLevel()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving zone urgency level",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":      "get zone urgency level successfully",
		"zone urgency": zoneUrgencyLevel,
	})
}

// Calculate ETA
func CalculateETA(distance, speed float64) time.Time {
	if speed <= 0 {
		fmt.Println("Error: speed must be greater than 0")
		return time.Now()
		// return time.Time{}, fmt.Errorf("speed must be greater than 0, got: %.2f", speed)
	}

	if distance < 0 {
		fmt.Println("Error :  distance cannot be negative")
		return time.Now()
		// return time.Time{}, fmt.Errorf("distance cannot be negative, got: %.2f", distance)
	}

	// Calculate time in hours
	hours := distance / speed

	// Convert hours to time
	minutes := hours * 60

	// Convert minutes to duration
	travelTime := time.Duration(minutes) * time.Minute

	// Get current time
	return time.Now().Add(travelTime)
}

// Convert time to thai time
func ThaiTimeFormat(date time.Time) string {
	thaiMonths := []string{
		"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
		"กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤษจิกายน", "ธันวาคม",
	}

	thaiYear := date.Year() + 543

	// Format the date to Thai format
	thaiDate := fmt.Sprintf("%d %s %d เวลา %02d:%02d น.",
		date.Day(), thaiMonths[date.Month()-1], thaiYear, date.Hour(), date.Minute())

	return thaiDate
}

// Format remaining time to string
func FormatRemainingTime(eta time.Time) string {
	now := time.Now()

	duration := eta.Sub(now)
	totalMinutes := int(duration.Minutes())

	if duration < 0 {
		return "เสร็จสิ้นแล้ว"

		// Calculate duration
		// delay := -duration
		// hours := int(delay.Hours())
		// minutes := int(delay.Minutes()) % 60

		// display delayed is hours and minutes
		// if hours > 0 {
		// 	return fmt.Sprintf("ล่าช้า %d ชั่วโมง %d นาที", hours, minutes)
		// }
		// return fmt.Sprintf("ล่าช้า %d นาที", minutes)
	}

	if totalMinutes < 60 {
		return fmt.Sprintf("%d นาที", totalMinutes)
	}

	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	if minutes > 0 {
		return fmt.Sprintf("%d ชั่วโมง %d นาที", hours, minutes)
	}

	return fmt.Sprintf("%d ชั่วโมง", hours)
}
