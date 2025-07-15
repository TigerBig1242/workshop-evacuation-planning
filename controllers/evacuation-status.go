package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/services"
)

// Get_evacuation_status from service
// func Get_evacuation_status(c *fiber.Ctx) error {
// 	// evacuationStatus, err := services.GetEvacuationStatus()
// 	var results []map[string]interface{}
// 	var lastVehicle int
// 	// if err != nil {
// 	// 	return c.Status(500).JSON(fiber.Map{
// 	// 		"Message": "Error fetching evacuation status",
// 	// 	})
// 	// }

// 	// return c.Status(200).JSON(fiber.Map{
// 	// 	"Message":     "Evacuation status fetched successfully",
// 	// 	"Data status": evacuationStatus,
// 	// })

// 	// Get all zones
// 	getZone, err := services.GetAllZones()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"Message": "Error fetching evacuation zones",
// 		})
// 	}
// 	// Get all evacuation plans
// 	getPlan, err := services.GetAllEvacuationPlans()
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"Message": "Error fetching evacuation plan",
// 		})
// 	}

// 	// Get all vehicles
// 	// getVehicles, err := services.GetAllVehicles()
// 	// if err != nil {
// 	// 	return c.Status(500).JSON(fiber.Map{
// 	// 		"Message": "Error fetching vehicles",
// 	// 	})
// 	// }

// 	for _, zone := range getZone {
// 		var remaining int
// 		var totalEvacuated int

// 		for _, plan := range getPlan {
// 			if plan["zone_id"] == zone.ID {
// 				evacuated := plan["people_evacuated"].(int)
// 				remaining = evacuated

// 				if plan["vehicle_id"] != nil {
// 					lastVehicle = plan["vehicle_id"].(int)
// 				}
// 				totalEvacuated = evacuated
// 				break

// 				// for _, vehicle := range getVehicles {
// 				// 	if vehicle.ID == uint(lastVehicle) {

// 				// 	}
// 				// }
// 			}
// 		}

// 		results = append(results, map[string]interface{}{
// 			"ZoneID":            zone.ID,
// 			"TotalEvacuated":    totalEvacuated,
// 			"Remaining":         remaining,
// 			"Last vehicle used": lastVehicle,
// 		})
// 	}
// 	fmt.Println("Evacuation Zones:", getZone)
// 	fmt.Println("Evacuation Plan:", getPlan)

// 	return c.Status(200).JSON(fiber.Map{
// 		"Message": "Fetched zones and plan successfully",
// 		"Results": results,
// 		// "Data status": evacuationStatus,
// 		// "Get zones success": getZone,
// 		// "Get Plans success": getPlan,
// 	})
// }

func EvacuationStatus(c *fiber.Ctx) error {
	evacuationPlans, err := services.GetEvacuationPlans()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving evacuation plans",
		})
	}

	evacuationZones, err := services.GetEvacuationZones()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"Message": "Error retrieving evacuation zones",
		})
	}

	evacuatedPeople := make(map[uint]int)
	var zoneStatus []map[string]interface{}
	// Loop to evacuees amount of plan evacuation.
	// If there is more than one round, add them to that round.
	for _, plan := range evacuationPlans {
		evacuatedPeople[uint(plan.Zone_id)] += plan.People_evacuated
		// fmt.Println("evacuated people:", evacuatedPeople)
	}

	// Loop to calculate remaining and evacuated people
	for i, zone := range evacuationZones {
		totalEvacuated := evacuatedPeople[zone.ID]
		// fmt.Println("total evacuated:", totalEvacuated)
		remaining := zone.Number_of_people - totalEvacuated

		// Protected not negative value
		if remaining < 0 {
			remaining = 0
		}

		zoneResult := map[string]interface{}{
			"zone_id":          zone.ID,
			"original_people":  zone.Number_of_people,
			"evacuated_people": totalEvacuated,
			"remaining_people": remaining,
			"urgency_level":    zone.Urgency_level,
		}

		if totalEvacuated > 0 {
			fmt.Printf("Evacuated zone id: %d Original people: %d Evacuated people: %d Remaining: %d\n",
				zone.ID, zone.Number_of_people, totalEvacuated, remaining)
		} else {
			fmt.Printf("No evacuation zone id: %d Original people: %d Evacuated people: %d Remaining: %d\n",
				zone.ID, zone.Number_of_people, totalEvacuated, remaining)
		}
		zoneStatus = append(zoneStatus, zoneResult)
		evacuationZones[i].Number_of_people = remaining
	}

	return c.Status(200).JSON(fiber.Map{
		"Message":      "Evacuation plans retrieved successfully",
		"zones status": zoneStatus,
	})
}
