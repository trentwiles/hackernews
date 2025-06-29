package email


import (
	"testing"

	//"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	SendEmail(Email{to: "me@trentwil.es", subject: "Hey! this is a test", message: "hi this is a test"})
}

func TestMagicLink(t *testing.T) {
	SendEmailTemplate(MagicLinkEmail{to: "me@trentwil.es", token: "123123123"})
}