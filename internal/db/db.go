package db

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/drhodes/golorem"
	"github.com/trentwiles/hackernews/internal/config"

	_ "github.com/lib/pq"
)

const DEFAULT_SELECT_LIMIT = 10

// sql database pool
var (
	db   *sql.DB
	once sync.Once
)

// Go representations of database objects
type User struct {
	Username      string
	Email         string
	Created_at    string
	Registered_ip string
	Score         int
}

type UserMetadata struct {
	Username  string
	Full_name string
	Birthdate string
	Bio_text  string
	IsAdmin   bool
}

type CompleteUser struct {
	User     User
	Metadata UserMetadata
}

type Submission struct {
	Id         string
	Title      string
	Username   string
	Link       string
	Body       string
	Flagged    bool
	Created_at string
	Votes      int
}

type BasicSubmission struct {
	Id         string
	Title      string
	Link       string
	Created_at string
}

type VoteMetrics struct {
	Upvotes   int
	Downvotes int
}

type BasicSubmissionAndVote struct {
	Title      string
	Link       string
	Body       string
	Created_at string
	Username   string
	Id         string
	IsUpvoted  bool
}

type AdminMetrics struct {
	TodayPosts           int
	TodayMinusOnePosts   int
	TodayMinusTwoPosts   int
	TodayMinusThreePosts int
	TodayMinusFourPosts  int
	TodayMinusFivePosts  int
	TodayMinusSixPosts   int

	TotalAllTimeSubmissions int

	TotalAllTimeUsers int
	TotalActiveUsers  int
}

type Comment struct {
	Id            string
	InResponseTo  string // uuid of the post the comment is being made on
	Content       string // body of the comment
	Author        string // username
	ParentComment string // <OPTIONAL> uuid of the parent comment
	Flagged       bool   // is this comment flagged for review?
	CreatedAt     string // timestamp in string format, typescript can interpret this as a Date object
	Upvotes       int
	Downvotes     int
	HasUpvoted    bool // has the user in question upvoted this post? TRUE if so...
	HasDownvoted  bool // has the user in question downvoted this post? TRUE if so...
}

// enum equiv in Go for audit log events
// ('login', 'logout', 'failed_login', 'post', 'comment', 'post_click', 'sent_email')
type AuditEvent string

const (
	Login       AuditEvent = "login"
	Logout      AuditEvent = "logout"
	FailedLogin AuditEvent = "failed_login"
	Post        AuditEvent = "post"
	PostClick   AuditEvent = "post_click"
	SentEmail   AuditEvent = "sent_email"
)

type SortMethod string

const (
	Latest SortMethod = "latest"
	Oldest SortMethod = "oldest"
	Best   SortMethod = "best"
)

func InitDB() error {
	var err error
	once.Do(func() {
		config.LoadEnv()
		connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			config.GetEnv("POSTGRES_USERNAME"),
			config.GetEnv("POSTGRES_PASSWORD"),
			config.GetEnv("POSTGRES_HOST"),
			config.GetEnv("POSTGRES_DB"))

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("[WARN] Issue establishing PostgreSQL connection: %s\n", err.Error())
			return
		}

		// connection pool config -- future, read this from an .env
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		// test connection
		err = db.Ping()
	})

	log.Println("[INFO] PostgreSQL connection pool created")
	return err
}

// creates the pool via InitDB()
func GetDB() *sql.DB {
	if db == nil {
		if err := InitDB(); err != nil {
			log.Fatal("Failed to initialize database:", err)
		}
	}

	log.Println("[INFO] PostgreSQL connection pool retrieved")
	return db
}

// backwards compatibility because i was lazy
func Connect() (*sql.DB, error) {
	return GetDB(), nil
}

func CreateUser(user User) {
	query := `INSERT INTO users (username, email, registered_ip) VALUES ($1, $2, $3)`

	_, err := GetDB().Exec(query, user.Username, user.Email, user.Registered_ip)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Create user %s with email %s from IP address %s\n", user.Username, user.Email, user.Registered_ip)
}

func UpsertUserMetadata(metadata UserMetadata) {
	// Validate input
	if metadata.Username == "" {
		log.Fatal("Please provide a username")
	}

	// Check if metadata exists for this username
	var exists bool
	err := GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM bio WHERE username = $1)", metadata.Username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Check that user %s exists in the database (%t)\n", metadata.Username, exists)

	if exists {
		// Update existing metadata
		_, err = GetDB().Exec(
			"UPDATE bio SET full_name = $1, birthdate = $2, bio_text = $3 WHERE username = $4",
			metadata.Full_name, metadata.Birthdate, metadata.Bio_text, metadata.Username,
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[INFO] User exists, updated user %s\n", metadata.Username)
	} else {
		// Insert new metadata
		_, err = GetDB().Exec(
			"INSERT INTO bio (username, full_name, birthdate, bio_text) VALUES ($1, $2, $3, $4)",
			metadata.Username, metadata.Full_name, metadata.Birthdate, metadata.Bio_text,
		)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] User did not exist, inserted user %s\n", metadata.Username)
	}
}

func SearchUser(user User) CompleteUser {
	// two cases: search by username and search by email
	if user.Email == "" && user.Username == "" {
		log.Fatal("To select a user, you must be either an email or username")
	}

	// username VARCHAR(100) PRIMARY KEY,
	// email VARCHAR(100) NOT NULL,
	// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	// registered_ip

	qUsername := `
	SELECT users.username, users.email, users.created_at, users.registered_ip,
				SUM(CASE 
					WHEN votes.positive = true THEN 1 
					WHEN votes.positive = false THEN -1 
					ELSE 0 
				END) AS score
	FROM users
	LEFT JOIN submissions ON users.username = submissions.username
	LEFT JOIN votes ON submissions.id = votes.submission_id
	WHERE users.username = $1
	GROUP BY users.username
	LIMIT 1
	`

	qEmail := `
	SELECT users.username, users.email, users.created_at, users.registered_ip,
				SUM(CASE 
					WHEN votes.positive = true THEN 1 
					WHEN votes.positive = false THEN -1 
					ELSE 0 
				END) AS score
	FROM users
	LEFT JOIN submissions ON users.username = submissions.username
	LEFT JOIN votes ON submissions.id = votes.submission_id
	WHERE users.email = $1
	GROUP BY users.username
	LIMIT 1
	`

	var rows *sql.Rows
	var err error

	if user.Username != "" {
		rows, err = GetDB().Query(qUsername, user.Username)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] Queried user %s via username\n", user.Username)
	} else {
		rows, err = GetDB().Query(qEmail, user.Email)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[INFO] Queried user email %s to get user\n", user.Email)
	}
	defer rows.Close()

	var tempUser = User{}
	for rows.Next() {
		err := rows.Scan(&tempUser.Username, &tempUser.Email, &tempUser.Created_at, &tempUser.Registered_ip, &tempUser.Score)
		if err != nil {
			log.Fatal(err)
		}
	}

	if tempUser.Username == "" {
		log.Printf("[INFO] User query search did not result in any user(s)\n")
	} else {
		log.Printf("[INFO] User query search resulted in user %s created @ %s\n", tempUser.Username, tempUser.Created_at)
	}

	// now that we've got the user themselves, let's grab their metadata
	var tempMetadata = UserMetadata{}
	if tempUser.Username != "" {
		query := `
		SELECT bio.username, bio.full_name, bio.birthdate, bio.bio_text, 
		CASE 
			WHEN admins.username IS NOT NULL THEN true
			ELSE false
		END AS isAdmin

		FROM bio
		LEFT JOIN admins ON bio.username = admins.username
		WHERE bio.username = $1;
		`
		rows, err = GetDB().Query(query, tempUser.Username)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&tempMetadata.Username, &tempMetadata.Full_name, &tempMetadata.Birthdate, &tempMetadata.Bio_text, &tempMetadata.IsAdmin)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Printf("[INFO] User query metadata search success, admin status: %t\n", tempMetadata.IsAdmin)
	}

	return CompleteUser{User: tempUser, Metadata: tempMetadata}
}

func DeleteUser(user User) {
	// two cases: search by username and search by email
	if user.Email == "" && user.Username == "" {
		log.Fatal("To delete a user, you must be either an email or username")
	}

	var err error
	if user.Username != "" {
		_, err = GetDB().Exec("DELETE FROM users WHERE username = $1", user.Username)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] Deleted user %s (via username)\n", user.Username)
	} else {
		_, err = GetDB().Exec("DELETE FROM users WHERE email = $1", user.Email)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] Deleted user via email, %s\n", user.Email)
	}

	// additional note: user bios are cascading, so Postgres will delete them automatically
	// future: find out if a user was deleted, via RETURNING?
}

// validation for the correct user is done in the API business logic
func DeleteSubmission(submission Submission) {
	if submission.Id == "" {
		log.Fatal("To delete a submission, you must pass a submission ID")
	}

	_, err := GetDB().Exec("DELETE FROM submissions WHERE id = $1", submission.Id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Deleted submission %s", submission.Id)
}

func SearchSubmission(stub Submission) Submission {
	if stub.Id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	rows, err := GetDB().Query("SELECT * FROM submissions WHERE id = $1", stub.Id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tempBody sql.NullString
		err := rows.Scan(&stub.Id, &stub.Username, &stub.Title, &stub.Link, &tempBody, &stub.Flagged, &stub.Created_at)
		if err != nil {
			log.Fatal(err)
		}

		if tempBody.String != "" {
			stub.Body = tempBody.String
		} else {
			stub.Body = ""
		}

		log.Printf("[INFO] Succesful query for submission %s created at %s", stub.Id, stub.Created_at)

		return stub
	}

	return Submission{}
}

func AllSubmissions(sort SortMethod, offset int) []Submission {
	// determine how to do the sorting itself
	var order string
	switch sort {
	case Latest:
		order = "ORDER BY created_at DESC"
		log.Printf("[INFO] Attempting all submissions sort query for filter 'latest'\n")
	case Oldest:
		order = "ORDER BY created_at ASC"
		log.Printf("[INFO] Attempting all submissions sort query for filter 'oldest'\n")
	case Best:
		order = `ORDER BY score DESC`
		log.Printf("[INFO] Attempting all submissions sort query for filter 'best'\n")
	}

	query := `
			SELECT submissions.id, username, title, link, body, created_at, 
				SUM(CASE 
					WHEN votes.positive = true THEN 1 
					WHEN votes.positive = false THEN -1 
					ELSE 0 
				END) AS score
			FROM submissions
			LEFT JOIN votes ON submissions.id = votes.submission_id
			WHERE flagged = false
			GROUP BY submissions.id
			` + order + `
			LIMIT $1 OFFSET $2`

	rows, err := GetDB().Query(query, DEFAULT_SELECT_LIMIT, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	log.Printf("[INFO] Queried all submissions with limit of %d, offset of %d\n", DEFAULT_SELECT_LIMIT, offset)

	var submissions []Submission
	for rows.Next() {
		var tempBody sql.NullString
		var current Submission

		if err := rows.Scan(&current.Id, &current.Username, &current.Title, &current.Link, &tempBody, &current.Created_at, &current.Votes); err != nil {
			log.Fatal(err)
		}

		if tempBody.Valid {
			current.Body = tempBody.String
		} else {
			current.Body = ""
		}

		submissions = append(submissions, current)
	}

	log.Printf("[INFO] Query resulted in %d submissions\n", len(submissions))

	return submissions
}

func LatestUserSubmissions(offset int, user User) []BasicSubmission {
	query := `
		SELECT id, title, link, created_at
		FROM submissions
		WHERE username = $1 AND flagged = false
		LIMIT $2 OFFSET $3
	`

	rows, err := GetDB().Query(query, user.Username, DEFAULT_SELECT_LIMIT, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var submissions []BasicSubmission
	for rows.Next() {
		var current BasicSubmission

		if err := rows.Scan(&current.Id, &current.Title, &current.Link, &current.Created_at); err != nil {
			log.Fatal(err)
		}

		submissions = append(submissions, current)
	}

	log.Printf("[INFO] Latest user submissions query for %s resulted in %d, using a limit of %d\n", user.Username, len(submissions), offset)

	return submissions
}

func CreateSubmission(submission Submission) string {
	query := `
		INSERT INTO submissions (username, title, link, body, flagged)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var id string
	err := GetDB().QueryRow(query, submission.Username, submission.Title, submission.Link, submission.Body, submission.Flagged).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] New submission authored by %s with ID %s created\n", submission.Username, id)

	return id
}

func UpdateSubmission(stub Submission) {
	if stub.Id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	_, err := GetDB().Exec("UPDATE submissions SET link = $1, title = $2, body = $3, flagged = $4 WHERE id = $5", stub.Link, stub.Title, stub.Body, stub.Flagged, stub.Id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Updated submission %s\n", stub.Id)
}

// true on success (new insert or update)
// false on failure (attempting to "double vote")
func Vote(user User, submission Submission, isUpvote bool) bool {
	// check that we have a valid username + submission id combo
	if user.Username == "" || submission.Id == "" {
		log.Fatal("username or submission id is blank (required to vote on a submission)")
	}

	// check if a vote already exists
	// if so run an update instead
	var wasPositive bool
	err := GetDB().QueryRow("SELECT positive FROM votes WHERE submission_id = $1 AND voter_username = $2", submission.Id, user.Username).Scan(&wasPositive)

	if err == sql.ErrNoRows {
		// if we enter this, there was no record found, so we need to do an insert
		query := `
			INSERT INTO votes (submission_id, voter_username, positive)
			VALUES ($1, $2, $3)
		`

		_, err = GetDB().Exec(query, submission.Id, user.Username, isUpvote)
		if err != nil {
			log.Fatal(err)
		}

		var voteType string
		if isUpvote {
			voteType = "upvote"
		} else {
			voteType = "downvote"
		}

		log.Printf("[INFO] Inserted new %s for post ID %s from %s\n", voteType, submission.Id, user.Username)

		return true
	} else if err != nil {
		log.Fatal(err)
	}

	// if we hit this point, a record was found, and now we just need to update it
	// "can't vote twice"
	if isUpvote == wasPositive {
		log.Printf("[INFO] Double vote attempted by %s\n", user.Username)
		return false
	}

	_, err = GetDB().Exec("UPDATE votes SET positive = $1 WHERE voter_username = $2 AND submission_id = $3", isUpvote, user.Username, submission.Id)
	if err != nil {
		log.Fatal(err)
	}

	var updated string

	if isUpvote {
		updated = "downvote to upvote"
	} else {
		updated = "upvote to downvote"
	}

	log.Printf("[INFO] Updated vote from a %s by user %s\n", updated, user.Username)

	return true
}

// Response meaning:
// Boolean #1: did the user vote on the post?
// Boolean #2: if so, did they upvote (true) or downvote (false)?
func GetUserVote(user User, submission Submission) (bool, bool) {
	if user.Username == "" {
		log.Fatal("missing username")
	}

	if submission.Id == "" {
		log.Fatal("missing submission ID")
	}

	var didUpvote bool
	err := GetDB().QueryRow("SELECT positive FROM votes WHERE voter_username = $1 AND submission_id = $2", user.Username, submission.Id).Scan(&didUpvote)

	if err == sql.ErrNoRows {
		fmt.Printf("[INFO] No vote found for user %s on post %s\n", user.Username, submission.Id)
		return false, false
	}

	if err != nil {
		log.Fatal(err)
	}

	var voteType string

	// i am PRAYING the next golang update includes ternaries
	if didUpvote {
		voteType = "upvote"
	} else {
		voteType = "downvote"
	}

	log.Printf("[INFO] Search found a %s on post %s by user %s\n", voteType, submission.Id, user.Username)

	return true, didUpvote
}

func GetAllUserVotes(user User) []BasicSubmissionAndVote {
	if user.Username == "" {
		log.Fatal("user's username cannot be blank")
	}
	// future: maybe instead of a string of IDs, use a string of submissions?
	query := `
		SELECT title, link, body, created_at, username, submission_id, positive
		FROM votes
		INNER JOIN submissions ON submission_id = submissions.id
		WHERE voter_username = $1
		LIMIT $2
	`

	rows, err := GetDB().Query(query, user.Username, DEFAULT_SELECT_LIMIT)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var submissions []BasicSubmissionAndVote
	for rows.Next() {
		var tempBody sql.NullString
		var current BasicSubmissionAndVote

		if err := rows.Scan(&current.Title, &current.Link, &tempBody, &current.Created_at, &current.Username, &current.Id, &current.IsUpvoted); err != nil {
			log.Fatal(err)
		}

		if tempBody.Valid {
			current.Body = tempBody.String
		} else {
			current.Body = ""
		}

		fmt.Println(current.IsUpvoted)

		submissions = append(submissions, current)
	}

	log.Printf("[INFO] Query for all user votes on user %s resulted in %d voted posts, w/ limit of %d\n", user.Username, len(submissions), DEFAULT_SELECT_LIMIT)

	return submissions
}

func CreateMagicLink(user User) string {
	if user.Username == "" || user.Email == "" {
		log.Fatal("to create a magic link, user must have a username and email")
	}

	// first, delete all old magic links for given user
	_, err := GetDB().Exec("DELETE FROM magic_links WHERE username = $1", user.Username)
	if err != nil {
		log.Fatal(err)
	}

	// next generate the secure token
	var token string = SecureToken(100)
	query := `
		INSERT INTO magic_links (username, email, token)
		VALUES ($1, $2, $3)
	`

	_, err = GetDB().Exec(query, user.Username, user.Email, token)
	if err != nil {
		log.Fatal(err)
	}

	// given a string 's'
	// len(s) --> returns the length of ASCII code
	// len([]rune(s)) --> returns the length of char count in string

	// []rune() returns an array of ascii values for each char in a string

	log.Printf("[INFO] Magic link for username %s and email %s created, length %d", user.Username, user.Email, len([]rune(token)))

	return token
}

func DeleteMagicLink(token string) {
	_, err := GetDB().Exec("DELETE FROM magic_links WHERE token = $1", token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Magic link for of length %d deleted", len([]rune(token)))
}

func ValidateMagicLink(token string, ip string) User {
	if token == "" {
		log.Fatal("You must pass in a token to validate")
	}

	var username string
	var email string

	err := GetDB().QueryRow("SELECT username, email FROM magic_links WHERE token = $1", token).Scan(&username, &email)

	log.Printf("[INFO] Magic link search found user %s and email %s for token of length %d\n", username, email, len([]rune(token)))

	if err == sql.ErrNoRows {
		log.Printf("[WARN] Magic link search found no user for token length %d\n", len([]rune(token)))
		// Empty user = invalid token
		return User{}
	}

	if username == "" {
		return User{}
	}

	if ip == "" {
		log.Fatal("Registration IP address required")
	}

	DeleteMagicLink(token)

	// determine if we need to insert the user into the database or not
	var searchedUser User = SearchUser(User{Username: username}).User

	if searchedUser.Username == "" {
		var toInsert User = User{Username: username, Email: email, Registered_ip: ip}
		CreateUser(toInsert)
		log.Printf("[INFO] User %s registration via magic link completed\n", username)
		return toInsert
	} else {
		log.Printf("[INFO] User %s login via magic link completed\n", username)
		return searchedUser
	}
}

// in the future make this into a single query, rather than counting positive, then negative votes
func CountVotes(post Submission) (VoteMetrics, error) {
	if post.Id == "" {
		return VoteMetrics{}, fmt.Errorf("cannot query votes with a blank submission ID")
	}

	var upvotes int
	var downvotes int
	err := GetDB().QueryRow("SELECT count(*) as ct FROM votes WHERE submission_id = $1 AND positive = $2", post.Id, true).Scan(&upvotes)
	if err != nil {
		return VoteMetrics{}, fmt.Errorf("error from SQL: %s", err)
	}

	log.Printf("[INFO] %d upvotes counted for submission ID %s\n", upvotes, post.Id)

	err = GetDB().QueryRow("SELECT count(*) as ct FROM votes WHERE submission_id = $1 AND positive = $2", post.Id, false).Scan(&downvotes)
	if err != nil {
		return VoteMetrics{}, fmt.Errorf("error from SQL: %s", err)
	}

	log.Printf("[INFO] %d downvotes counted for submission ID %s\n", downvotes, post.Id)

	return VoteMetrics{Upvotes: upvotes, Downvotes: downvotes}, nil
}

func SearchSubmissionByQuery(query string, offset int) []Submission {
	if offset < 0 {
		log.Printf("[WARN] Offset in SearchSubmissionByQuery %d is <0, set to 0\n", offset)
		offset = 0
	}

	if query == "" {
		log.Printf("[WARN] Unable to search for a blank query. Returning empty list.\n")
		return []Submission{}
	}

	q := `
		SELECT * FROM submissions
		WHERE flagged = false
		AND (title ILIKE $1 OR body ILIKE $1)
		LIMIT $2 OFFSET $3
	`

	rows, err := GetDB().Query(q, "%"+query+"%", DEFAULT_SELECT_LIMIT, offset)
	if err != nil {
		log.Fatal(err)
	}

	var resultList []Submission
	var tempResult Submission
	for rows.Next() {
		var tempBody sql.NullString

		err := rows.Scan(&tempResult.Id, &tempResult.Username, &tempResult.Title, &tempResult.Link, &tempBody, &tempResult.Flagged, &tempResult.Created_at)
		if err != nil {
			log.Fatal(err)
		}

		if tempBody.String != "" {
			tempResult.Body = tempBody.String
		} else {
			tempResult.Body = ""
		}

		resultList = append(resultList, tempResult)
	}

	log.Printf("[INFO] Submission search query returned %d results, with a limit of %d\n", len(resultList), DEFAULT_SELECT_LIMIT)

	return resultList
}

// ATTN: consider caching this route, it's expensive
func GetAdminMetrics() AdminMetrics {

	// Last seven days, how many posts per day?
	query := `
		SELECT
			days_between((NOW() - INTERVAL '1 day')::TIMESTAMP, NOW()::TIMESTAMP) as today,
			days_between((NOW() - INTERVAL '2 days')::TIMESTAMP, (NOW() - INTERVAL '1 day')::TIMESTAMP) as todayMinusOne,
			days_between((NOW() - INTERVAL '3 days')::TIMESTAMP, (NOW() - INTERVAL '2 days')::TIMESTAMP) as todayMinusTwo,
			days_between((NOW() - INTERVAL '4 days')::TIMESTAMP, (NOW() - INTERVAL '3 days')::TIMESTAMP) as todayMinusThree,
			days_between((NOW() - INTERVAL '5 days')::TIMESTAMP, (NOW() - INTERVAL '4 days')::TIMESTAMP) as todayMinusFour,
			days_between((NOW() - INTERVAL '6 days')::TIMESTAMP, (NOW() - INTERVAL '5 days')::TIMESTAMP) as todayMinusFive,
			days_between((NOW() - INTERVAL '7 days')::TIMESTAMP, (NOW() - INTERVAL '6 days')::TIMESTAMP) as todayMinusSix;	
	`
	var today int
	var todayMinusOne int
	var todayMinusTwo int
	var todayMinusThree int
	var todayMinusFour int
	var todayMinusFive int
	var todayMinusSix int

	err := GetDB().QueryRow(query).Scan(&today, &todayMinusOne, &todayMinusTwo, &todayMinusThree, &todayMinusFour, &todayMinusFive, &todayMinusSix)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Database made admin query for # of submissions over the last 7 days\n")

	var totalPosts int
	query = `
		SELECT count(*)
		FROM submissions
	`
	err = GetDB().QueryRow(query).Scan(&totalPosts)
	if err != nil {
		log.Fatal(err)
	}

	var totalUsers int
	query = `
		SELECT count(*)
		FROM users
	`
	err = GetDB().QueryRow(query).Scan(&totalUsers)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Database made admin query for # of total users\n")

	// goal of the actives users query: how many users have made a post/voted in the last week?
	// when comments become avaialble, this should include comments?
	var totalActiveUsers int
	query = `
		SELECT COUNT(DISTINCT username) AS active_users
		FROM (
			SELECT username FROM submissions
			WHERE created_at BETWEEN (NOW() - INTERVAL '7 days') AND NOW()
			
			UNION

			SELECT voter_username AS username FROM votes
			WHERE ts BETWEEN (NOW() - INTERVAL '7 days') AND NOW()
		) AS active;
	`
	err = GetDB().QueryRow(query).Scan(&totalActiveUsers)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Database made admin query for # of active users over the last 7 days\n")

	return AdminMetrics{
		TodayPosts:              today,
		TodayMinusOnePosts:      todayMinusOne,
		TodayMinusTwoPosts:      todayMinusTwo,
		TodayMinusThreePosts:    todayMinusThree,
		TodayMinusFourPosts:     todayMinusFour,
		TodayMinusFivePosts:     todayMinusFive,
		TodayMinusSixPosts:      todayMinusSix,
		TotalAllTimeSubmissions: totalPosts,
		TotalAllTimeUsers:       totalUsers,
		TotalActiveUsers:        totalActiveUsers,
	}
}

func CheckAdminStatus(user User) bool {
	if user.Username == "" {
		log.Fatal("[FATAL] Unable to check admin status of user with blank username\n")
	}

	query := `
		SELECT EXISTS (
			SELECT 1 FROM admins WHERE username = $1
		);
		`

	var exists bool
	err := GetDB().QueryRow(query, user.Username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	return exists
}

func GenerateNonsenseData(userCount int, postCount int) {
	// fake names
	var firstNames []string = []string{"Alice", "Bob", "Charlie", "Diana", "Emma", "Frank", "Grace", "Henry", "Isabella", "Jack", "Katherine", "Liam", "Maria", "Noah", "Olivia", "Peter", "Quinn", "Rachel", "Samuel", "Tara", "Uma", "Victor", "Wendy", "Xavier", "Yasmine", "Zachary"}
	var lastNames []string = []string{"Anderson", "Brown", "Chen", "Davis", "Evans", "Fisher", "Garcia", "Harris", "Johnson", "Kim", "Lopez", "Miller", "Nguyen", "O'Brien", "Patel", "Quinn", "Rodriguez", "Smith", "Taylor", "Upton", "Vasquez", "Williams", "Xavier", "Young", "Zhang"}

	// fake users
	var chars string = "qwertyuiopadfghjklzxcvbnm1234567890"
	var usernames []string
	// create 10 users
	for i := 0; i < userCount; i++ {
		var username string = ""
		// each user has a 10 char username
		for a := 0; a < 10; a++ {
			username = username + string(chars[rand.Intn(len(chars))])
		}
		usernames = append(usernames, username)
	}

	for _, value := range usernames {
		CreateUser(User{Username: value, Email: value + "@gmail.com", Registered_ip: "0.0.0.0"})
		UpsertUserMetadata(UserMetadata{Username: value, Full_name: firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))], Birthdate: "01/01/2004", Bio_text: "this is a fake user generated by an automated script... don't listen to anything they say!"})
	}

	for i := 0; i < postCount; i++ {
		var s Submission = Submission{Title: lorem.Sentence(3, 6), Body: lorem.Paragraph(50, 500), Username: usernames[rand.Intn(len(usernames))], Link: "http://www.example.com"}
		CreateSubmission(s)
	}
}

func InsertNewComment(comment Comment) string {
	// bare minimum requirements for a new comment
	if comment.InResponseTo == "" || comment.Author == "" || comment.Content == "" {
		log.Fatal("[FATAL] Attempted to insert new comment without one or more of the following: InResponseTo, Author, Content")
	}

	var id string

	if comment.ParentComment != "" {
		query := `
			INSERT INTO comments (in_response_to, content, author, parent_comment)
			VALUES ($1, $2, $3, $4)
			RETURNING id;
		`

		err := GetDB().QueryRow(query, comment.InResponseTo, comment.Content, comment.Author, comment.ParentComment).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] Database made comment insertion in response to %s WITH a parent comment\n", comment.InResponseTo)
	} else {
		query := `
			INSERT INTO comments (in_response_to, content, author)
			VALUES ($1, $2, $3)
			RETURNING id;
		`

		err := GetDB().QueryRow(query, comment.InResponseTo, comment.Content, comment.Author).Scan(&id)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[INFO] Database made comment insertion in response to %s WITHOUT a parent comment\n", comment.InResponseTo)
	}

	return id
}

func GetCommentsOnSubmission(submission Submission, contextUser User) []Comment {
	if submission.Id == "" {
		log.Fatal("Please use an ID when searching for a submission's comments")
	}

	// if contextUser.Username == "" {
	// 	log.Fatal("Please provide a valid username, so we can determine which comments the user has voted on")
	// }

	query := `
		SELECT
			c.id,
			c.in_response_to,
			c.content,
			c.author,
			c.parent_comment,
			c.flagged,
			c.created_at,
			COUNT(CASE WHEN cv.positive = TRUE THEN 1 END) AS upvotes,
			COUNT(CASE WHEN cv.positive = FALSE THEN 1 END) AS downvotes,
			
			-- TRUE if has voted, FALSE if hasn't voted, FALSE if no results in comment_votes (edge case)
			COALESCE(BOOL_OR(cv.voter_username = $2 AND cv.positive = TRUE), FALSE) AS has_upvoted,
			COALESCE(BOOL_OR(cv.voter_username = $2 AND cv.positive = FALSE), FALSE) AS has_downvoted
		FROM comments c
		LEFT JOIN comment_votes cv ON c.id = cv.comment_id
		WHERE c.in_response_to = $1
		GROUP BY c.id, c.in_response_to, c.content, c.author, c.parent_comment, c.flagged, c.created_at
		ORDER BY c.created_at;	
	`

	// no limits/offset here at the moment, do this in a future update
	rows, err := GetDB().Query(query, submission.Id, contextUser.Username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var commentHolder []Comment

	for rows.Next() {
		// the following fields may be NULL:
		var parentComment sql.NullString

		var tempComment Comment
		err := rows.Scan(&tempComment.Id, &tempComment.InResponseTo, &tempComment.Content, &tempComment.Author, &parentComment, &tempComment.Flagged, &tempComment.CreatedAt, &tempComment.Upvotes, &tempComment.Downvotes, &tempComment.HasUpvoted, &tempComment.HasDownvoted)
		if err != nil {
			log.Fatal(err)
		}

		if !parentComment.Valid {
			tempComment.ParentComment = ""
		} else {
			tempComment.ParentComment = parentComment.String
		}

		commentHolder = append(commentHolder, tempComment)

	}

	return commentHolder
}

// get the comments on a post, plus if the user has voted on the comments

func DeleteComment(comment Comment) {
	if comment.Id == "" {
		log.Fatal("Please provide a comment ID to delete a comment")
	}

	_, err := GetDB().Exec("DELETE FROM comments WHERE id = $1", comment.Id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Deleted (or attempted to delete) comment ID %s\n", comment.Id)
}

func VoteOnComment(user User, comment Comment, isUpvote bool) bool {
	// check that we have a valid username + comment id combo
	if user.Username == "" || comment.Id == "" {
		log.Fatal("username or comment id is blank (required to vote on a comment)")
	}

	// check if a vote already exists
	var wasPositive bool
	err := GetDB().QueryRow(
		"SELECT positive FROM comment_votes WHERE comment_id = $1 AND voter_username = $2",
		comment.Id, user.Username,
	).Scan(&wasPositive)

	if err == sql.ErrNoRows {
		// insert new vote
		query := `
			INSERT INTO comment_votes (comment_id, voter_username, positive)
			VALUES ($1, $2, $3)
		`
		_, err = GetDB().Exec(query, comment.Id, user.Username, isUpvote)
		if err != nil {
			log.Fatal(err)
		}

		voteType := "downvote"
		if isUpvote {
			voteType = "upvote"
		}
		log.Printf("[INFO] Inserted new %s for comment ID %s from %s\n", voteType, comment.Id, user.Username)
		return true
	} else if err != nil {
		log.Fatal(err)
	}

	// update existing vote if changed
	if isUpvote == wasPositive {
		log.Printf("[INFO] Double vote attempted on comment by %s\n", user.Username)
		return false
	}

	_, err = GetDB().Exec(
		"UPDATE comment_votes SET positive = $1 WHERE voter_username = $2 AND comment_id = $3",
		isUpvote, user.Username, comment.Id,
	)
	if err != nil {
		log.Fatal(err)
	}

	updated := "upvote to downvote"
	if isUpvote {
		updated = "downvote to upvote"
	}
	log.Printf("[INFO] Updated comment vote from a %s by user %s\n", updated, user.Username)
	return true
}
