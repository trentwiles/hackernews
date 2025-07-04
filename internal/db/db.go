package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/trentwiles/hackernews/internal/config"

	_ "github.com/lib/pq"
)

const DEFAULT_SELECT_LIMIT = 10

// Go representations of database objects
type User struct {
	Username      string
	Email         string
	Created_at    string
	Registered_ip string
}

type UserMetadata struct {
	Username  string
	Full_name string
	Birthdate string
	Bio_text  string
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
}

type VoteMetrics struct {
	Upvotes   int
	Downvotes int
}

// enum equiv in Go for audit log events
// ('login', 'logout', 'failed_login', 'post', 'comment', 'post_click', 'sent_email')
type AuditEvent string

const (
	Login       AuditEvent = "login"
	Logout      AuditEvent = "logout"
	FailedLogin AuditEvent = "failed_login"
	Post        AuditEvent = "post"
	Comment     AuditEvent = "comment"
	PostClick   AuditEvent = "post_click"
	SentEmail   AuditEvent = "sent_email"
)

type SortMethod string

const (
	Latest SortMethod = "latest"
	Oldest SortMethod = "oldest"
	Best   SortMethod = "best"
)

func Connect() (*sql.DB, error) {
	config.LoadEnv()
	// connect to postgres
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", config.GetEnv("POSTGRES_USERNAME"), config.GetEnv("POSTGRES_PASSWORD"), config.GetEnv("POSTGRES_DB"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func CreateUser(user User) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	query := `INSERT INTO users (username, email, registered_ip) VALUES ($1, $2, $3)`

	_, err = db.Exec(query, user.Username, user.Email, user.Registered_ip)
	if err != nil {
		log.Fatal(err)
	}
}

func UpsertUserMetadata(metadata UserMetadata) {
	// Validate input
	if metadata.Username == "" {
		log.Fatal("Please provide a username")
	}

	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if metadata exists for this username
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM bio WHERE username = $1)", metadata.Username).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		// Update existing metadata
		_, err = db.Exec(
			"UPDATE bio SET full_name = $1, birthdate = $2, bio_text = $3 WHERE username = $4",
			metadata.Full_name, metadata.Birthdate, metadata.Bio_text, metadata.Username,
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Insert new metadata
		_, err = db.Exec(
			"INSERT INTO bio (username, full_name, birthdate, bio_text) VALUES ($1, $2, $3, $4)",
			metadata.Username, metadata.Full_name, metadata.Birthdate, metadata.Bio_text,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SearchUser(user User) CompleteUser {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	// two cases: search by username and search by email
	if user.Email == "" && user.Username == "" {
		log.Fatal("To select a user, you must be either an email or username")
	}

	var rows *sql.Rows

	if user.Username != "" {
		rows, err = db.Query("SELECT * FROM users WHERE username = $1", user.Username)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows, err = db.Query("SELECT * FROM users WHERE email = $1", user.Email)
		if err != nil {
			log.Fatal(err)
		}
	}

	var tempUser = User{}
	for rows.Next() {
		err := rows.Scan(&tempUser.Username, &tempUser.Email, &tempUser.Created_at, &tempUser.Registered_ip)
		if err != nil {
			log.Fatal(err)
		}
	}

	// now that we've got the user themselves, let's grab their metadata
	var tempMetadata = UserMetadata{}
	if tempUser.Username != "" {
		rows, err = db.Query("SELECT * FROM bio WHERE username = $1", tempUser.Username)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			err := rows.Scan(&tempMetadata.Username, &tempMetadata.Full_name, &tempMetadata.Birthdate, &tempMetadata.Bio_text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return CompleteUser{User: tempUser, Metadata: tempMetadata}
}

func DeleteUser(user User) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	// two cases: search by username and search by email
	if user.Email == "" && user.Username == "" {
		log.Fatal("To delete a user, you must be either an email or username")
	}

	if user.Username != "" {
		_, err = db.Exec("DELETE FROM users WHERE username = $1", user.Username)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = db.Exec("DELETE FROM users WHERE email = $1", user.Email)
		if err != nil {
			log.Fatal(err)
		}
	}

	// additional note: user bios are cascading, so Postgres will delete them automatically
}

// validation for the correct user is done in the API business logic
func DeleteSubmission(submission Submission) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if submission.Id == "" {
		log.Fatal("To delete a submission, you must pass a submission ID")
	}

	_, err = db.Exec("DELETE FROM submissions WHERE id = $1", submission.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func SearchSubmission(stub Submission) Submission {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if stub.Id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	rows, err := db.Query("SELECT * FROM submissions WHERE id = $1", stub.Id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&stub.Id, &stub.Username, &stub.Title, &stub.Link, &stub.Body, &stub.Flagged, &stub.Created_at)
		if err != nil {
			log.Fatal(err)
		}

		return stub
	}

	return Submission{}
}

func AllSubmissions(sort SortMethod, offset int) []Submission {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// determine how to do the sorting itself
	var order string
	switch sort {
	case Latest:
		order = "ORDER BY created_at DESC"
	case Oldest:
		order = "ORDER BY created_at ASC"
	case Best:
		fmt.Println("WARNING: BEST SELECT WAS USED!!! this has not been implemented, either don't use it or implement it")
		order = ""
	}

	query := `
		SELECT id, username, title, link, body, created_at
		FROM submissions
		WHERE flagged = false
		` + order + `
		LIMIT $1 OFFSET $2
	`

	rows, err := db.Query(query, DEFAULT_SELECT_LIMIT, offset)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var tempBody sql.NullString
		var current Submission

		if err := rows.Scan(&current.Id, &current.Username, &current.Title, &current.Link, &tempBody, &current.Created_at); err != nil {
			log.Fatal(err)
		}

		if tempBody.Valid {
			current.Body = tempBody.String
		} else {
			current.Body = ""
		}

		submissions = append(submissions, current)
	}

	return submissions
}

func CreateSubmission(submission Submission) string {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		INSERT INTO submissions (username, link, body, flagged)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	var id string
	err = db.QueryRow(query, submission.Username, submission.Title, submission.Link, submission.Body, submission.Flagged).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}

func UpdateSubmission(stub Submission) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if stub.Id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	_, err = db.Exec("UPDATE submissions SET link = $1, title = $2, body = $3, flagged = $4 WHERE id = $5", stub.Link, stub.Title, stub.Body, stub.Flagged, stub.Id)
	if err != nil {
		log.Fatal(err)
	}
}

// true on success (new insert or update)
// false on failure (attempting to "double vote")
func Vote(user User, submission Submission, isUpvote bool) bool {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	// check that we have a valid username + submission id combo
	if user.Username == "" || submission.Id == "" {
		log.Fatal("username or submission id is blank (required to vote on a submission)")
	}

	// check if a vote already exists
	// if so run an update instead
	var wasPositive bool
	err = db.QueryRow("SELECT positive FROM votes WHERE submission_id = $1 AND voter_username = $2", submission.Id, user.Username).Scan(&wasPositive)

	if err == sql.ErrNoRows {
		// if we enter this, there was no record found, so we need to do an insert
		query := `
			INSERT INTO votes (submission_id, voter_username, positive)
			VALUES ($1, $2, $3)
		`

		_, err = db.Exec(query, submission.Id, user.Username, isUpvote)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else if err != nil {
		log.Fatal(err)
	}

	// if we hit this point, a record was found, and now we just need to update it
	// "can't vote twice"
	if isUpvote == wasPositive {
		return false
	}

	_, err = db.Exec("UPDATE votes SET positive = $1 WHERE voter_username = $2 AND submission_id = $3", isUpvote, user.Username, submission.Id)
	if err != nil {
		log.Fatal(err)
	}

	return true
}

func CreateMagicLink(user User) string {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if user.Username == "" || user.Email == "" {
		log.Fatal("to create a magic link, user must have a username and email")
	}

	// first, delete all old magic links for given user
	_, err = db.Exec("DELETE FROM magic_links WHERE username = $1", user.Username)
	if err != nil {
		log.Fatal(err)
	}

	// next generate the secure token
	var token string = SecureToken(100)
	query := `
			INSERT INTO magic_links (username, email, token)
			VALUES ($1, $2, $3)
		`

	_, err = db.Exec(query, user.Username, user.Email, token)
	if err != nil {
		log.Fatal(err)
	}

	return token
}

func DeleteMagicLink(token string) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	_, err = db.Exec("DELETE FROM magic_links WHERE token = $1", token)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateMagicLink(token string, ip string) User {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if token == "" {
		log.Fatal("You must pass in a token to validate")
	}

	var username string
	var email string

	err = db.QueryRow("SELECT username, email FROM magic_links WHERE token = $1", token).Scan(&username, &email)

	fmt.Printf("Database search found user %s and email %s for token %s\n", username, email, token)

	if err == sql.ErrNoRows {
		fmt.Printf("No user found for token: %s\n", token)
		// Empty user = invalid token
		return User{}
	}

	if username == "" {
		return User{}
	}

	if ip == "" {
		log.Fatal("Registration IP address required")
	}

	fmt.Println("About to delete magic link")
	DeleteMagicLink(token)
	fmt.Println("Deleted magic link")

	fmt.Println("Checking if we need to insert a new user into the database")

	// determine if we need to insert the user into the database or not
	var searchedUser User = SearchUser(User{Username: username}).User

	if searchedUser.Username == "" {
		var toInsert User = User{Username: username, Email: email, Registered_ip: ip}
		CreateUser(toInsert)
		fmt.Println("Inserted new user into the database")
		return toInsert
	} else {
		return searchedUser
	}

}

func CountVotes(post Submission) (VoteMetrics, error) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// end connection via connection function

	if post.Id == "" {
		return VoteMetrics{}, fmt.Errorf("cannot query votes with a blank submission ID")
	}

	var upvotes int
	var downvotes int
	err = db.QueryRow("SELECT count(*) as ct FROM votes WHERE submission_id = $1 AND positive = $2", post.Id, true).Scan(&upvotes)
	if err != nil {
		return VoteMetrics{}, fmt.Errorf("error from SQL: %s", err)
	}

	err = db.QueryRow("SELECT count(*) as ct FROM votes WHERE submission_id = $1 AND positive = $2", post.Id, false).Scan(&downvotes)
	if err != nil {
		return VoteMetrics{}, fmt.Errorf("error from SQL: %s", err)
	}

	return VoteMetrics{Upvotes: upvotes, Downvotes: downvotes}, nil
}
