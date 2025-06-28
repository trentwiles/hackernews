package db

import (
    "crypto/rand"
    "math/big"
)

// source: https://www.slingacademy.com/article/how-to-generate-secure-random-numbers-in-go/
func SecureToken(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

    password := make([]byte, length)
    for i := range password {
        n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
        password[i] = letters[n.Int64()]
    }

	return string(password)
}