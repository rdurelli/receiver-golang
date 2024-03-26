package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"receiver/handlers"
	"receiver/models"
	"receiver/routes"
)

func Init() *gorm.DB {
	dbURL := "postgres://pg:pass@database:5432/crud"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	err = db.AutoMigrate(&models.FileToConvert{})
	if err != nil {
		return nil
	}

	return db
}

func main() {
	// Create a background context
	ctx := context.Background()
	DB := Init()
	h := handlers.New(DB)
	app := fiber.New(fiber.Config{
		BodyLimit: 15 * 1024 * 1024, // this is the default limit of 15MB
	})
	// Initialize default config
	app.Use(logger.New())
	routes.SetupRoutes(app, h, ctx)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}
