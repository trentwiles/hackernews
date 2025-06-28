package email


import (
	"testing"

	//"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	SendEmail(Email{to: "", subject: "Hey! this is a test", message: "hi this is a test"})
}