package db

import (
	"testing"

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

	DeleteUser(james)

	var res2 UserMetadata = SearchUser(james).metadata
	assert.Equal(t, res2.full_name, "", "metadata ensure postgres cascade delete")
}

func TestUserDeletion(t *testing.T) {
	DeleteUser(User{username: "james"})
}