package main

import (
	"fmt"
	"os"

	"github.com/tigerbig1242/evacuation-planning/config"
	"github.com/tigerbig1242/evacuation-planning/logger"
	"github.com/tigerbig1242/evacuation-planning/models"
	"github.com/tigerbig1242/evacuation-planning/routes"
)

func main() {

	config.ConnectDatabase()

	// Auto migrate models
	config.DB.AutoMigrate(&models.Evacuation_zone{}, &models.Vehicle{}, &models.Evacuation_plan{},
		&models.Evacuation_status{}, models.EvacuationLog{})

	// Initialize logger
	// logger.InitLogger()

	err := logger.InitLogger()
	if err != nil {
		fmt.Printf("Fielded to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Log.Sync()

	// Set up routes
	app := routes.SetRoutes()
	app.Listen(":8080")

	fmt.Println("Finished calculation.")

	fmt.Println("Connected Successful from database file!")
}
