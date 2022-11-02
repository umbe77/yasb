package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/umbe77/yasb/handlers"
)

func main() {
	fmt.Println("Hello, world")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.All("/service", handlers.ExecuteService)

	log.Fatal(app.Listen(":3000"))

}
