package config

import (
	"fmt"
	"log"
	"os"
	"time"

	// "github.com/tigerbig1242/evacuation-planning/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "tigerbig"
	password = "tigerbig1242"
	dbname   = "evacuation-planning-db"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Connect to the database
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	// Open a connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		// log.Fatal(err)
		panic("Failed to connect database !")
	}

	DB = db

	fmt.Println(DB)
	fmt.Println("Database Migrated Successfully !")
}
