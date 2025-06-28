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