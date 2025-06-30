package main

import (
	"github.com/gofiber/fiber/v2"

	// my packages
	"github.com/trentwiles/hackernews/internal/db"
	"github.com/trentwiles/hackernews/internal/email"
	"github.com/trentwiles/hackernews/internal/jwt"
	"github.com/trentwiles/hackernews/internal/captcha"

	_ "github.com/lib/pq"
)

type BasicResponse struct {
	Message string
	Status  int
}

type LoginRequest struct {
    Email string `json:"email"`
	Username string `json:"username"`
    CaptchaToken string `json:"captchaToken"`
}

var version string = "/api/v1"

func main() {
	// create web app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}
		return c.JSON(BasicResponse{Message: "Logged in as " + username, Status: 200})
	})

	app.Post(version + "/login", func(c *fiber.Ctx) error {
		var req LoginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		if req.Email == "" || req.Username == "" || req.CaptchaToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "email, username, and captchaToken are required",
			})
		}

		// logic to verify Google Captcha
		if !captcha.ValidateToken(req.CaptchaToken) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Google Captcha response, try again later",
			})
		}

		// database calls:
		// check if email exists in the database:
		// 		exists?  => check if email/username combo matches, otherwise return error
		//		doesn't? => send magic link email, once user has verified the magic link, marry them together in the database

		var databaseUser db.CompleteUser = db.SearchUser(db.User{Email: req.Email})
		if databaseUser.User.Username != req.Username {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email with this username already exists, and the submitted email doesn't said username",
			})
		}


		var token string = db.CreateMagicLink(db.User{Email: req.Email})

		_ = email.MagicLinkEmail(email.MagicLinkEmail{To: req.Email, Token: token})

		return c.Status(fiber.StatusConflict).JSON(fiber.Map{})
	})

	app.Get(version+"/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Healthy", "status": 200})
	})

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found", "status": 404})
	})

	app.Listen("127.0.0.1:3000") //ok
}
