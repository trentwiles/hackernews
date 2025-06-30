package email


import (
	"testing"

	//"github.com/stretchr/testify/assert"
)

func TestUserCreation(t *testing.T) {
	SendEmail(Email{To: "me@trentwil.es", Subject: "Hey! this is a test", Message: "hi this is a test"})
}

func TestMagicLink(t *testing.T) {
	SendEmailTemplate(MagicLinkEmail{To: "me@trentwil.es", Token: "123123123"})
}