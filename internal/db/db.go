package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/trentwiles/hackernews/internal/config"

	_ "github.com/lib/pq"
)

// Go representations of database objects
type User struct {
	username string
	email string
	created_at string
	registered_ip string
}

type UserMetadata struct {
	username string
	full_name string
    birthdate string
    bio_text string
}

type CompleteUser struct {
	user User
	metadata UserMetadata
}

type Submission struct {
	id string
    username string
    link string
    body string
    flagged bool
}

// enum equiv in Go for audit log events
// ('login', 'logout', 'failed_login', 'post', 'comment', 'post_click', 'sent_email')
type AuditEvent string

const (
	Login  AuditEvent = "login"
	Logout   AuditEvent = "logout"
	FailedLogin  AuditEvent = "failed_login"
	Post AuditEvent = "post"
	Comment AuditEvent = "comment"
	PostClick AuditEvent = "post_click"
	SentEmail AuditEvent = "sent_email"
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

    _, err = db.Exec(query, user.username, user.email, user.registered_ip)
    if err != nil {
        log.Fatal(err)
    }
}

func CreateUserMetadata(metadata UserMetadata) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
        log.Fatal(err)
    }
	defer db.Close()
	// end connection via connection function

	query := `INSERT INTO bio (username, full_name, birthdate, bio_text) VALUES ($1, $2, $3, $4)`

    _, err = db.Exec(query, metadata.username, metadata.full_name, metadata.birthdate, metadata.bio_text)
    if err != nil {
        log.Fatal(err)
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
	if (user.email == "" && user.username == "") {
		log.Fatal("To select a user, you must be either an email or username")
	}

	var rows *sql.Rows
	
	if (user.username != "") {
		rows, err = db.Query("SELECT * FROM users WHERE username = $1", user.username)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		rows, err = db.Query("SELECT * FROM users WHERE email = $1", user.email)
		if err != nil {
			log.Fatal(err)
		}
	}

	var tempUser = User{}
	for rows.Next() {
        err := rows.Scan(&tempUser.username, &tempUser.email, &tempUser.created_at, &tempUser.registered_ip)
        if err != nil {
            log.Fatal(err)
        }
    }

	// now that we've got the user themselves, let's grab their metadata
	var tempMetadata = UserMetadata{}
	if tempUser.username != "" {
		rows, err = db.Query("SELECT * FROM bio WHERE username = $1", tempUser.username)
		if err != nil {
            log.Fatal(err)
        }

		for rows.Next() {
			err := rows.Scan(&tempMetadata.username, &tempMetadata.full_name, &tempMetadata.birthdate, &tempMetadata.bio_text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return CompleteUser{user: tempUser, metadata: tempMetadata}
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
	if (user.email == "" && user.username == "") {
		log.Fatal("To delete a user, you must be either an email or username")
	}

	
	if (user.username != "") {
		_, err = db.Exec("DELETE FROM users WHERE username = $1", user.username)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = db.Exec("DELETE FROM users WHERE email = $1", user.email)
		if err != nil {
			log.Fatal(err)
		}
	}

	// additional note: user bios are cascading, so Postgres will delete them automatically
}

func DeleteSubmission(submission Submission) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
        log.Fatal(err)
    }
	defer db.Close()
	// end connection via connection function

	if submission.id == "" {
		log.Fatal("To delete a submission, you must pass a submission ID")
	}

	_, err = db.Exec("DELETE FROM submissions WHERE id = $1", submission.id)
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

	if stub.id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	rows, err := db.Query("SELECT * FROM submissions WHERE id = $1", stub.id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&stub.id, &stub.username, &stub.link, &stub.body, &stub.flagged)
		if err != nil {
			log.Fatal(err)
		}

		return stub
	}

	return Submission{}
}

func CreateSubmission(submission Submission) string {
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `
		INSERT INTO submissions (username, link, body, flagged)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id string
	err = db.QueryRow(query, submission.username, submission.link, submission.body, submission.flagged).Scan(&id)
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

	if stub.id == "" {
		log.Fatal("Please use an ID when searching for a submission")
	}

	_, err = db.Exec("UPDATE submissions SET link = $1, body = $2, flagged = $3 WHERE id = $4", stub.link, stub.body, stub.flagged, stub.id)
	if err != nil {
		log.Fatal(err)
	}
}


func UpdateUserMetadata(metadata UserMetadata) {
	// connection via connection function
	db, err := Connect()
	if err != nil {
        log.Fatal(err)
    }
	defer db.Close()
	// end connection via connection function

	if metadata.username == "" {
		log.Fatal("Please use a username updating user metadata")
	}

	_, err = db.Exec("UPDATE bio SET full_name = $1, birthdate = $2, bio_text = $3 WHERE username = $4", metadata.full_name, metadata.birthdate, metadata.bio_text, metadata.username)
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
	if user.username == "" || submission.id == "" {
		log.Fatal("username or submission id is blank (required to vote on a submission)")
	}

	// check if a vote already exists
	// if so run an update instead
	var wasPositive bool
	err = db.QueryRow("SELECT positive FROM votes WHERE submission_id = $1 AND voter_username = $2", submission.id, user.username).Scan(&wasPositive)

	if err == sql.ErrNoRows {
		// if we enter this, there was no record found, so we need to do an insert
		query := `
			INSERT INTO votes (submission_id, voter_username, positive)
			VALUES ($1, $2, $3)
		`

		_, err = db.Exec(query, user.username, submission.id, isUpvote)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else if err != nil {
		log.Fatal(err)
	}

	// if we hit this point, a record was found, and now we just need to update it
	// "can't vote twice"
	if (isUpvote == wasPositive) {
		return false
	}

	_, err = db.Exec("UPDATE votes SET positive = $1 WHERE voter_username = $2 AND submission_id = $3", isUpvote, user.username, submission.id)
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

	if user.username == "" || user.email == "" {
		log.Fatal("to create a magic link, user must have a username and email")
	}

	// first, delete all old magic links for given user
	_, err = db.Exec("DELETE FROM magic_links WHERE username = $1", user.username)
	if err != nil {
		log.Fatal(err)
	}

	// next generate the secure token
	var token string = SecureToken(100)
	query := `
			INSERT INTO magic_links (username, token)
			VALUES ($1, $2)
		`

	_, err = db.Exec(query, user.username, token)
    if err != nil {
        log.Fatal(err)
    }

	return token
}