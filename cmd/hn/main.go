package main

import (
	"github.com/gofiber/fiber/v2"

	// my packages
	"github.com/trentwiles/hackernews/internal/db"

	_ "github.com/lib/pq"
)

type Response struct {
	Message string
	Status  int
}

var version string = "/api/v1"

func main() {
	// create web app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(Response{Message: "ok", Status: 200})
	})

	app.Get(version+"/status", func(c *fiber.Ctx) error {
		return c.JSON(Response{Message: "ok", Status: 200})

	})

	app.Listen(":3000")
}
