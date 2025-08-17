package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/net/html"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	// my packages
	"github.com/trentwiles/hackernews/internal/captcha"
	"github.com/trentwiles/hackernews/internal/config"
	"github.com/trentwiles/hackernews/internal/db"
	"github.com/trentwiles/hackernews/internal/dump"
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

type CommentCreationRequest struct {
	InResponseTo string `json:"inResponseTo"`
	Content string `json:"content"`
	CaptchaToken string `json:"captchaToken"`
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
	app.Use(logger.New())

	app.Static("/", "./static")

	log.Println("[INFO] Started webserver with CORS & Logging middleware")

	config.LoadEnv()
	expiresString := config.GetEnv("TOKENS_EXPIRE_IN")
	TOKEN_EXPIRES_IN, err := strconv.Atoi(expiresString)
	if err != nil {
		log.Fatalf("Invalid TOKENS_EXPIRE_IN (parse error): %v", err)
	}

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

	// 	if !success {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
	// 	}
	// 	return c.JSON(BasicResponse{Message: "Logged in as " + username, Status: 200})
	// })

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
		jwtToken, err := jwt.GenerateJWT(user.Username, TOKEN_EXPIRES_IN)

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

	// determine if a user has voted on a post, and if they have
	// whether it is an upvote or downvote
	app.Get(version+"/vote", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please pass an id parameter"})
		}

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		didVote, didUpvote := db.GetUserVote(db.User{Username: username}, db.Submission{Id: id})

		if !didVote {
			return c.JSON(fiber.Map{
				"didVote": false,
			})
		}

		return c.JSON(fiber.Map{
			"didVote":   true,
			"didUpvote": didUpvote,
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

	app.Get(version+"/allUserVotes", func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please pass a username parameter"})
		}

		posts := db.GetAllUserVotes(db.User{Username: username})

		return c.JSON(fiber.Map{
			"results": posts,
		})
	})

	// aka get the metadata and votes on a post
	app.Get(version+"/submission", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Please pass an id parameter"})
		}

		var queriedSubmission db.Submission = db.SearchSubmission(db.Submission{Id: id})

		votes, err := db.CountVotes(db.Submission{Id: id})
		if err != nil {
			log.Fatal(err)
		}

		return c.JSON(fiber.Map{
			"id": queriedSubmission.Id,
			"metadata": fiber.Map{
				"title":     queriedSubmission.Title,
				"link":      queriedSubmission.Link,
				"body":      queriedSubmission.Body,
				"author":    queriedSubmission.Username,
				"isFlagged": queriedSubmission.Flagged,
				"createdAt": queriedSubmission.Created_at,
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
				"error":   true,
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
				"error":   true,
				"message": "invalid sort filter",
			})
		}

		// if we run the select query on the database, and there are less
		// than the max amount of submissions available, then we know we've
		// hit the end of the list; therefore, the next value should be null
		if len(selection) != db.DEFAULT_SELECT_LIMIT {
			return c.JSON(fiber.Map{
				"results": selection,
				"next":    nil,
			})
		}

		return c.JSON(fiber.Map{
			"results": selection,
			"next":    version + "/all?sort=" + sortType + "&offset=" + strconv.Itoa(offsetInt+10),
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

		if user.Username == "" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "No such user or account deleted",
			})
		}

		return c.JSON(fiber.Map{
			"username": user.Username,
			"email":    user.Email,
			"joined":   user.Created_at,
			"metadata": fiber.Map{
				"full_name": userMetadata.Full_name,
				"birthday":  userMetadata.Birthdate,
				"bio":       userMetadata.Bio_text,
				"isAdmin":   userMetadata.IsAdmin,
				"score":     user.Score,
			},
		})
	})

	app.Get(version+"/userSubmissions", func(c *fiber.Ctx) error {
		username := c.Query("username")
		if username == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please pass a username parameter",
			})
		}

		offset := c.Query("offset")
		if offset == "" {
			offset = "0"
		}

		var tempUser db.User = db.User{Username: username}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "error parsing 'offset', " + err.Error(),
			})
		}

		search := db.LatestUserSubmissions(offsetInt, tempUser)

		if search == nil {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"next":    nil,
				"results": []any{},
			})
		}

		// if we run the select query on the database, and there are less
		// than the max amount of submissions available, then we know we've
		// hit the end of the list; therefore, the next value should be null
		if len(search) != db.DEFAULT_SELECT_LIMIT {
			return c.JSON(fiber.Map{
				"results": search,
				"next":    nil,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"next":    fmt.Sprintf("%s/userSubmissions?username=%s&offset=%d", version, username, offsetInt+10),
			"results": search,
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

	app.Get(version+"/fetchWebsiteTitle", func(c *fiber.Ctx) error {
		// success, _ := jwt.ParseAuthHeader(c.Get("Authorization"))

		// if !success {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		// }

		url := c.Query("url")
		if url == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please pass a url parameter",
			})
		}

		// make the request "object"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		// set headers
		req.Header.Set("User-Agent", "HackerNewsClone (+https://github.com/trentwiles/hackernews)")

		// send the request object we just made
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}

		defer resp.Body.Close()

		tokenizer := html.NewTokenizer(resp.Body)

		for {
			tt := tokenizer.Next()
			switch tt {
			case html.ErrorToken:
				log.Println("no title found/blocked")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   true,
					"message": "unable to fetch title for provided URL",
				})
			case html.StartTagToken, html.SelfClosingTagToken:
				token := tokenizer.Token()
				if token.Data == "title" {
					tokenizer.Next()
					return c.JSON(fiber.Map{
						"title": string(tokenizer.Text()),
					})
				}
			}
		}

	})

	app.Get(version+"/searchSubmissions", func(c *fiber.Ctx) error {
		// success, _ := jwt.ParseAuthHeader(c.Get("Authorization"))

		// if !success {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		// }

		q := c.Query("q")
		if q == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Please pass a `q` parameter",
			})
		}

		page := c.Query("page")
		if page == "" {
			page = "0"
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "`page` parse error: " + err.Error(),
			})
		}

		var offset int = pageInt * db.DEFAULT_SELECT_LIMIT

		query := db.SearchSubmissionByQuery(q, offset)

		return c.JSON(fiber.Map{
			"results": query,
		})
	})

	app.Get(version+"/adminMetrics", func(c *fiber.Ctx) error {
		// success, _ := jwt.ParseAuthHeader(c.Get("Authorization"))

		// if !success {
		// 	return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		// }

		// again, in the future check if the user is an admin, but for MVP, it doesn't really matter

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"metrics": db.GetAdminMetrics(),
		})
	})

	app.Get(version+"/checkAdmin", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		return c.JSON(fiber.Map{
			"isAdmin": db.CheckAdminStatus(db.User{Username: username}),
		})
	})

	// POST /api/v1/comment?parent=123123123
	app.Post(version+"/comment", func(c *fiber.Ctx) error {
		parent := c.Query("parent")

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		var req CommentCreationRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse JSON",
			})
		}

		// VALIDATE CAPTCHA TOKEN
		if !captcha.ValidateToken(req.CaptchaToken) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Google Captcha response, try again later",
			})
		}
		// END VALIDATE CAPTCHA TOKEN

		var yourComment db.Comment = db.Comment{InResponseTo: req.InResponseTo, Content: req.Content, Author: username}
		if parent != "" {
			yourComment.ParentComment = parent
		}

		fmt.Printf("yourComment (full debug): %+v\n", yourComment)


		var commentId string = db.InsertNewComment(yourComment)

		return c.JSON(fiber.Map{
			"success": true,
			"commentID": commentId,
		})
	})

	app.Get(version + "/comments", func(c *fiber.Ctx) error {
		parent := c.Query("id") // submissionID
		username := c.Query("username") // has comment been upvoted by ...



		if parent == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing submission `id`",
			})
		}

		msg := ``

		if username == "" {
			msg = "`username` parameter is NULL, meaning that hasUpvoted and hasDownvoted will always be false"
		}

		if msg != "" {
			return c.JSON(fiber.Map{
				"notice": msg,
				"comments": db.GetCommentsOnSubmission(db.Submission{Id: parent}, db.User{Username: username}),
			})
		}

		return c.JSON(fiber.Map{
			"comments": db.GetCommentsOnSubmission(db.Submission{Id: parent}, db.User{Username: username}),
		})
	})

	app.Delete(version + "/comment", func(c *fiber.Ctx) error {
		id := c.Query("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing comment `id`",
			})
		}

		db.DeleteComment(db.Comment{Id: id})

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Deleted comment " + id,
		})
	})

	app.Post(version+"/commentVote", func(c *fiber.Ctx) error {
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
				"error": "missing valid comment `id` parameter",
			})
		}

		return c.JSON(fiber.Map{
			"success": db.VoteOnComment(db.User{Username: username}, db.Comment{Id: req.Id}, req.Upvote),
		})
	})

	app.Post(version + "/generateKey", func(c *fiber.Ctx) error {
		// here's the catch: to create an API key, you must already be authenticated,
		// that is, you must click on the code in your email, then log in
		// in the future, I'll consider developing a way to avoid email

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		var key string = db.CreateUserAPIKey(db.User{Username: username})

		// by this stage, we assume the username is valid
		log.Printf("[WARN] Created API key for user %s, but no check was made that this user has already created an API key (future: add this)\n", username)

		return c.JSON(fiber.Map{"username": username, "apiKey": key, "comment": "Store this API key in a safe place."})
	})

	app.Post(version + "/generateKey", func(c *fiber.Ctx) error {
		// here's the catch: to create an API key, you must already be authenticated,
		// that is, you must click on the code in your email, then log in
		// in the future, I'll consider developing a way to avoid email

		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		// by this stage, we assume the username is valid

		var key string = db.CreateUserAPIKey(db.User{Username: username})

		return c.JSON(fiber.Map{"username": username, "apiKey": key, "comment": "Store this API key in a safe place."})
	})

	app.Post(version + "/dump", func(c *fiber.Ctx) error {
		success, username := jwt.ParseAuthHeader(c.Get("Authorization"))

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		// future: add ratelimit/cooldown period

		var dumpLocation string = dump.DumpForUser(db.User{Username: username})
		// for the user example, the dump would be stored at exports\example\

		cmd := exec.Command("zip", "-r", "exports/" + username + ".zip", dumpLocation)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Println("[WARN] Error running zip on dump:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "internal zipping error, please contact site administrator if you see this message",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
		})
	})

	app.Get(version + "/dump", func(c *fiber.Ctx) error {
		// checks if the current logged in user has an available dump

		auth := c.Query("authToken")

		if auth == "" {
			return c.Status(fiber.StatusBadRequest).JSON(BasicResponse{Message: "missing authorization token", Status: fiber.StatusUnauthorized})
		}

		success, username := jwt.ParseAuthString(auth)

		if !success {
			return c.Status(fiber.StatusUnauthorized).JSON(BasicResponse{Message: "not signed in", Status: fiber.StatusUnauthorized})
		}

		_, err := os.Stat("exports/" + username + ".zip")
		if err != nil {
			return c.Status(fiber.StatusGone).JSON(fiber.Map{
				"success": false,
				"message": "Data dump was either already downloaded or never created",
			})
		}

		return c.SendFile("exports/" + username + ".zip")
	})


	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Route not found", "status": 404})
	})

	app.Listen("0.0.0.0:30000")
}
