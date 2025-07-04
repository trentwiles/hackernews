package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// my packages
	"github.com/trentwiles/hackernews/internal/captcha"
	"github.com/trentwiles/hackernews/internal/db"
	"github.com/trentwiles/hackernews/internal/email"
	"github.com/trentwiles/hackernews/internal/jwt"
	"github.com/trentwiles/hackernews/internal/utils"

	_ "github.com/lib/pq"
)

// consider moving these to a types.go file?
type BasicResponse struct {
	Message string
	Status  int
}

type LoginRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	CaptchaToken string `json:"captchaToken"`
}

type SubmissionRequest struct {
	Link         string `json:"link"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	CaptchaToken string `json:"captchaToken"`
}

type SubmissionDeleteRequest struct {
	Id string `json:"id"`
}

// username VARCHAR(100) PRIMARY KEY,
// full_name VARCHAR(100),
// birthdate DATE,
// bio_text TEXT,

type BioUpdateRequest struct {
	FullName  string `json:"fullName"`  // full_name
	Birthdate string `json:"birthdate"` // birthdate
	BioText   string `json:"bioText"`   // bio_text
}

type VoteRequest struct {
	Id     string `json:"id"`
	Upvote bool   `json:"upvote"`
}

var version string = "/api/v1"

func main() {
	// create web app
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}
		return c.JSON(BasicResponse{Message: "Logged in as " + username, Status: 200})
	})

	app.Post(version+"/login", func(c *fiber.Ctx) error {
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

		// is the email address even valid?
		if !utils.IsValidEmail(req.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email (failed regex)",
			})
		}

		// are they both under 100 chars (limit as defined in postgres)
		if len(req.Email) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email (over 100 chars)",
			})
		}

		if len(req.Username) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid username (over 100 chars)",
			})
		}

		// does username pass filter?
		// has only letters, underscores, and numbers
		// future: doesn't contain slurs/other forbidden words
		if !utils.IsValidUsername(req.Username) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid username (must contain only letters, numbers, and underscores)",
			})
		}

		// database calls:
		// check if email exists in the database:
		// 		exists?  => check if email/username combo matches, otherwise return error
		//		doesn't? => send magic link email, once user has verified the magic link, marry them together in the database

		var databaseUser db.CompleteUser = db.SearchUser(db.User{Email: req.Email})
		// case exists, but doesn't match
		if databaseUser.User.Username != "" && databaseUser.User.Username != req.Username {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email with this username already exists, and the submitted email doesn't said username",
			})
		}

		// case exists and does match
		if databaseUser.User.Username == req.Username {
			fmt.Printf("Attempted a sign in for %s\n", req.Username)
		}

		// email does not exist in the database
		if databaseUser.User.Username == "" {
			fmt.Printf("Email %s does not have a username tied to it in the database.\n", req.Email)
		}

		var token string = db.CreateMagicLink(db.User{Username: req.Username, Email: req.Email})
		email.SendEmailTemplate(email.MagicLinkEmail{To: req.Email, Token: token})

		return c.JSON(fiber.Map{"message": "Emailed a magic link to " + req.Email})
	})

	app.Get(version+"/magic", func(c *fiber.Ctx) error {
		token := c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please pass a valid token parameter"})
		}

		// verify token
		var user db.User = db.ValidateMagicLink(token, c.IP())
		var blankUser db.User = db.User{}
		if user == blankUser {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"message": "Magic link was not found. Maybe it expired?"})
		}

		var jwtToken string
		jwtToken, err := jwt.GenerateJWT(user.Username, 60)

		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(fiber.Map{"username": user.Username, "token": jwtToken})
	})

	app.Post(version+"/submit", func(c *fiber.Ctx) error {
		var req SubmissionRequest

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		if req.CaptchaToken == "" || req.Link == "" || req.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing one or more of the following parameters: captchaToken, title, link",
			})
		}

		// is the link valid (passes regex and length restriction?)
		if !utils.IsValidURL(req.Link) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid URL (failed regex)",
			})
		}

		if len(req.Link) > 255 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid URL (exceeds 255 char limit)",
			})
		}

		// passed all checks and restrictions now insert into database
		var id string = db.CreateSubmission(db.Submission{Title: req.Title, Username: username, Body: req.Body, Link: req.Link})

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	})

	app.Post(version+"/bio", func(c *fiber.Ctx) error {
		var req BioUpdateRequest

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		if req.BioText == "" {
			fmt.Println("debug note: bio text is empty")
		}

		if req.Birthdate == "" {
			fmt.Println("debug note: birth date is empty")
		}

		if !utils.IsValidDateFormat(req.Birthdate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Birth date fails regex. Not in American format, MM-DD-YYYY.",
			})
		}

		if req.FullName == "" {
			fmt.Println("debug note: full name is empty")
		}

		if len(req.FullName) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Full name cannot be longer than 100 chars",
			})
		}

		// all validations have passed
		var meta db.UserMetadata = db.UserMetadata{Username: username, Full_name: req.FullName, Birthdate: req.Birthdate, Bio_text: req.BioText}
		db.UpsertUserMetadata(meta)

		return c.JSON(fiber.Map{
			"message": "Updated metadata for user " + username,
		})
	})

	app.Post(version+"/vote", func(c *fiber.Ctx) error {
		var req VoteRequest

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		if req.Id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing valid id parameter",
			})
		}

		// all parameters have been validated
		var voteSuccess bool = db.Vote(db.User{Username: username}, db.Submission{Id: req.Id}, req.Upvote)

		// for now just return if the vote went through or not (false = trued to doublevote)
		// future return the count of votes
		return c.JSON(fiber.Map{"id": req.Id, "voteSuccess": voteSuccess})
	})

	// aka get the metadata and votes on a post
	app.Get(version+"/submission", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please pass an id parameter"})
		}

		if id == "" {
			return c.JSON(fiber.Map{"message": "Please pass an id parameter"})
		}

		var queriedSubmission db.Submission = db.SearchSubmission(db.Submission{Id: id})

		votes, err := db.CountVotes(db.Submission{Id: id})
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(fiber.Map{
			"id": queriedSubmission.Id,
			"metdata": fiber.Map{
				"title":     queriedSubmission.Title,
				"link":      queriedSubmission.Link,
				"body":      queriedSubmission.Body,
				"author":    queriedSubmission.Username,
				"isFlagged": queriedSubmission.Flagged,
			},
			"votes": fiber.Map{
				"upvotes":   votes.Upvotes,
				"downvotes": votes.Downvotes,
				"total":     votes.Upvotes - votes.Downvotes,
			},
		})
	})

	// grab all the submissions to display on the front page
	app.Get(version+"/all", func(c *fiber.Ctx) error {
		sortType := c.Query("sort")
		if sortType == "" {
			sortType = "latest"
		}

		offset := c.Query("offset")
		if offset == "" {
			offset = "0"
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"message": "error parsing 'offset', " + err.Error(),
			})
		}


		var selection []db.Submission

		switch sortType {
		case "latest":
			// ORDER BY created_time DESC
			fmt.Println("Latest placeholder")
			selection = db.AllSubmissions(db.Latest, offsetInt)
		case "best":
			// some sort of advanced SQL command to calculate all upvotes
			fmt.Println("Best placeholder")
			selection = db.AllSubmissions(db.Best, offsetInt)
		case "oldest":
			// ORDER BY created_time ASC
			selection = db.AllSubmissions(db.Oldest, offsetInt)
		default:
			fmt.Println("default placeholder")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"message": "invalid sort filter",
			})
		}

		return c.JSON(fiber.Map{
			"results": selection,
			"next": version + "/all?sort=" + sortType + "&offset=" + strconv.Itoa(offsetInt + 10),
		})


	})

	app.Get(version+"/user", func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please pass a username parameter",
			})
		}

		var user db.User
		var userMetadata db.UserMetadata

		complete := db.SearchUser(db.User{Username: username})

		user = complete.User
		userMetadata = complete.Metadata

		return c.JSON(fiber.Map{
			"username": user.Username,
			"email": user.Email,
			"joined":   user.Created_at,
			"metadata": fiber.Map{
				"full_name": userMetadata.Full_name,
				"birthday":  userMetadata.Birthdate,
				"bio":       userMetadata.Bio_text,
			},
		})
	})

	// shorthand for /api/v1/user?username=<authenticated_user>
	app.Get(version+"/me", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		return c.Redirect(version + "/user?username=" + username)
	})

	app.Delete(version+"/submission", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		var req SubmissionDeleteRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		if req.Id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "must pass a valid id parameter",
			})
		}

		// first check if submission matches the requested username
		query := db.SearchSubmission(db.Submission{Id: req.Id})

		if query.Id == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No such submission (has it already been deleted?)",
			})
		}

		if query.Username != username {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You do not own this post, and therefore cannot delete it",
			})
		}

		// prevent a flagged post from being deleted (can't let people destroy the evidence!)
		if query.Flagged {
			return c.Status(fiber.StatusLocked).JSON(fiber.Map{
				"error": "This post is currently under review, and cannot be deleted during this process",
			})
		}

		fmt.Printf("debug - deleted a post with id %s by user %s", req.Id, username)
		db.DeleteSubmission(query)

		return c.JSON(fiber.Map{
			"message": "OK",
		})
	})

	app.Get(version+"/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Healthy", "status": 200})
	})

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Route not found", "status": 404})
	})

	app.Listen("127.0.0.1:3000")
}
