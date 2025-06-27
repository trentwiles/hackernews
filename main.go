package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Response struct {
  Message string
  Status  int
}

var version string = "/api/v1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(Response{Message: "ok", Status: 200})
	})

	app.Get(version + "/status", func(c *fiber.Ctx) error {
		return c.JSON(Response{Message: "ok", Status: 200})
		// => Get request with value: hello world
	})

	app.Listen(":3000")
}
