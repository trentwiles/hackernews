package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	CreateUser(User{username: "james", email: "test@example.com", registered_ip: "127.0.0.1"})
	var res User = SearchUser(User{username: "james"})
	
	assert.Equal(t, res.username, "james", "select statement for sample user (username)")
	assert.Equal(t, res.email, "test@example.com", "select statement for sample user (email)")
	assert.Equal(t, res.registered_ip, "127.0.0.1", "select statement for sample user (ip address)")

	DeleteUser(User{username: "james"})
}

func TestUserDeletion(t *testing.T) {
	DeleteUser(User{username: "james"})
}