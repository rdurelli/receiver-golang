package main

import (
	"context"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/grafana/loki-client-go/loki"
	slogloki "github.com/samber/slog-loki/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"os"
	"receiver/handlers"
	"receiver/models"
	"receiver/routes"
)

func Init() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	//create the url connection
	dbURL := "postgres://" + dbUser + ":" + dbPass + "@" + dbHost + ":" + dbPort + "/" + dbName

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

// init loki
func InitLoki() *slog.Logger {
	config, _ := loki.NewDefaultConfig("http://loki:3100/loki/api/v1/push")
	config.TenantID = "xyz"
	client, _ := loki.New(config)

	logger := slog.New(slogloki.Option{Level: slog.LevelDebug, Client: client}.NewLokiHandler())
	logger = logger.
		With("environment", "dev").
		With("release", "v1.0.0")
	return logger
}

func main() {
	// Create a background context
	ctx := context.Background()
	DB := Init()
	h := handlers.New(DB)
	app := fiber.New(fiber.Config{
		BodyLimit: 15 * 1024 * 1024, // this is the default limit of 15MB
	})

	prometheus := fiberprometheus.New("my-golang-receiver")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	routes.SetupRoutes(app, h, ctx, InitLoki())

	serverPort := os.Getenv("SERVER_PORT")

	err := app.Listen(":" + serverPort)
	if err != nil {
		return
	}
}
