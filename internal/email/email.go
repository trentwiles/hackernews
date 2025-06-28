package email

import (
	"net/smtp"
	"log"
	"github.com/trentwiles/hackernews/internal/config"
)

type Email struct {
	to string
	subject string
	message string
}

func SendEmail(email Email) {
	config.LoadEnv()

	from := config.GetEnv("EMAIL_USERNAME") // <your_email>@gmail.com
	password := config.GetEnv("EMAIL_PASSWORD") // google app password

	toEmailAddress := email.to
	to := []string{toEmailAddress}

	host := config.GetEnv("EMAIL_HOST")
	port := "587"
	address := host + ":" + port

	subject := email.subject
	body := email.message
	
	message := []byte("To: " + email.to + "\r\n" +
    "Subject: " + subject + "\r\n" +
    "MIME-version: 1.0;\r\n" +
    "Content-Type: text/plain; charset=\"UTF-8\";\r\n" +
    "\r\n" +
    body + "\r\n")

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}
