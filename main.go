package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/umbe77/yasb/config"
	"github.com/umbe77/yasb/handlers"
	"github.com/umbe77/yasb/store"
)

func main() {
	conf := config.GetConfig()
	db, err := store.NewStore(conf.MongoUri)
	if err != nil {
		log.Fatal(err)
	}
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(logger.New(logger.Config{
		Format: logger.ConfigDefault.Format,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!!!",
		})
	})

	// app.All("/service", handlers.ExecuteService)
	app.Get("/api/workflows", handlers.GetWorflows(db))
	app.Post("/api/workflows", handlers.AddWorkflow(db))

	log.Fatal(app.Listen(":8080"))

}
