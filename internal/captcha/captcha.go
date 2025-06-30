package captcha

import (
	"fmt"
)

func ValidateToken(token string) bool {
	fmt.Println("Debug: validated token via the dummy function!")

	return true
}