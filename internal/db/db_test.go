package db

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	CreateUser(User{username: "james", email: "test@example.com", registered_ip: "127.0.0.1"})
	var res CompleteUser = SearchUser(User{username: "james"})
	var coreUserData User = res.user
	
	assert.Equal(t, coreUserData.username, "james", "select statement for sample user (username)")
	assert.Equal(t, coreUserData.email, "test@example.com", "select statement for sample user (email)")
	assert.Equal(t, coreUserData.registered_ip, "127.0.0.1", "select statement for sample user (ip address)")

	DeleteUser(User{username: "james"})
}

func TestUserBioCreation(t *testing.T) {
	var james User = User{username: "james", email: "test@example.com", registered_ip: "127.0.0.1"}
	var jamesPersonalDetails UserMetadata = UserMetadata{username: "james", full_name: "James Johnson", birthdate: "02/03/2001", bio_text: "blahblahblahblahblahblah"}
	CreateUser(james)
	CreateUserMetadata(jamesPersonalDetails)

	var res UserMetadata = SearchUser(james).metadata

	assert.Equal(t, res.full_name, "James Johnson", "metadata full name")

	jamesPersonalDetails.bio_text = "this is my new bio"
	UpdateUserMetadata(jamesPersonalDetails)

	res = SearchUser(james).metadata

	assert.Equal(t, res.bio_text, "this is my new bio", "metadata test bio update")

	DeleteUser(james)

	var res2 UserMetadata = SearchUser(james).metadata
	assert.Equal(t, res2.full_name, "", "metadata ensure postgres cascade delete")
}

func TestSubmissionCreation(t *testing.T) {
	// need a test user in the database due to FKs
	var james User = User{username: "james", email: "test@example.com", registered_ip: "127.0.0.1"}
	CreateUser(james)

	var testSubmission Submission = Submission{username: "james", link: "https://www.google.com", body: "here's a search engine", flagged: true}
	var generatedID string = CreateSubmission(testSubmission)
	var searchedSubmission Submission = SearchSubmission(Submission{id: generatedID})

	assert.Equal(t, searchedSubmission.body, testSubmission.body, "submission insert, check body is the same")
	assert.Equal(t, searchedSubmission.flagged, testSubmission.flagged, "submission insert, check flagged is the same")
	assert.Equal(t, searchedSubmission.link, testSubmission.link, "submission insert, check link is the same")

	DeleteSubmission(Submission{id: generatedID})

	searchedSubmission = SearchSubmission(Submission{id: generatedID})
	assert.Equal(t, searchedSubmission.link, "", "submission delete, ensure link is blank")

	DeleteUser(james)
}

func TestCreateMagicLink(t *testing.T) {
	// need a test user in the database due to FKs
	var james User = User{username: "james", email: "test@example.com", registered_ip: "127.0.0.1"}
	CreateUser(james)

	var insertedToken string = CreateMagicLink(james)
	fmt.Println(insertedToken)

	DeleteUser(james)
}