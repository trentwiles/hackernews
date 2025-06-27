package main

import (
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Response struct {
	Message string
	Status  int
}

var version string = "/api/v1"

func main() {
	// load envs
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to postgres
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// create web app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(Response{Message: "ok", Status: 200})
	})

	app.Get(version+"/status", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT $1", 1)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(rows)

		return c.JSON(Response{Message: "ok", Status: 200})

	})

	app.Listen(":3000")
}
