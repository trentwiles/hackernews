package utils

import (
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func IsValidURL(url string) bool {
	regex := regexp.MustCompile(`(?i)^(https?|ftp):\/\/[^\s/$.?#].[^\s]*$`)
	return regex.MatchString(url)
}