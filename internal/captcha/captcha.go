package captcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"log"
	"strconv"
	"strings"

	"github.com/trentwiles/hackernews/internal/config"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
	ErrorCodes  []string `json:"error-codes"`
}

func ValidateToken(token string) bool {
	// for user privacy I don't pass the request IP address back to Google
	resp, err := sendCaptchaToken(token)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatal("JSON decode error:", err)
	}

	if !result.Success {
		fmt.Println("erorr from google captcha API v3")
		for num, error_message := range result.ErrorCodes {
			fmt.Printf("%d. %s\n", num, error_message)
		}
		log.Fatal("\n\nGoogle captcha API request failed, see message(s) above.")
	}

	fmt.Printf("Debug: Google Captcha Response was a success. Score for action %s was a %.1f/1.0\n", result.Action, result.Score)
	
	// Is the score provided by google higher or lower than the cutoff?
	cutoff, err := strconv.ParseFloat(config.GetEnv("GOOGLE_CUTOFF"), 64)
	if err != nil {
		log.Fatal("Error parsing float:", err)
	}

	return result.Score >= cutoff
}


func sendCaptchaToken(token string) (*http.Response, error) {
	config.LoadEnv()

	endpoint := "https://www.google.com/recaptcha/api/siteverify"
	data := url.Values{
		"secret":   {config.GetEnv("GOOGLE_SECRET_KEY")},
		"response": {token},
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")


	client := &http.Client{}
	resp, err := client.Do(req)

	return resp, err
}