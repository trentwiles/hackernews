package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trentwiles/hackernews/internal/config"
)

func generateJWT(username string, expiresIn int) (string, error) {
	config.LoadEnv()
	claims := jwt.MapClaims{
		"username": username,
		"nbf": time.Now().Add(-1 * time.Minute).Unix(), // Valid starting 1 minute ago
		"exp": time.Now().Add(time.Duration(expiresIn) * time.Minute).Unix(),  // Expires in 5 minutes
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetEnv("JWT_TOKEN")))
}

func verifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetEnv("JWT_TOKEN")), nil
	})

	if err != nil {
		return "", fmt.Errorf("parse error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// `nbf` = not valid before
	if nbfFloat, ok := claims["nbf"].(float64); ok {
		nbf := time.Unix(int64(nbfFloat), 0)
		if time.Now().Before(nbf) {
			return "", fmt.Errorf("token not valid before: %v", nbf)
		}
	}

	// `exp` = expires time
	if expFloat, ok := claims["exp"].(float64); ok {
		exp := time.Unix(int64(expFloat), 0)
		if time.Now().After(exp) {
			return "", fmt.Errorf("token expired at: %v", exp)
		}
	}

	return claims["username"].(string), nil
}
