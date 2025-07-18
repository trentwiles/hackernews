package email

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"

	"github.com/trentwiles/hackernews/internal/config"
)

type Email struct {
	To      string
	Subject string
	Message string
}

type MagicLinkEmail struct {
	To    string
	Token string
}

// unused testing function
func SendEmail(email Email) {
	config.LoadEnv()

	from := config.GetEnv("EMAIL_USERNAME")     // <your_email>@gmail.com
	password := config.GetEnv("EMAIL_PASSWORD") // google app password

	toEmailAddress := email.To
	to := []string{toEmailAddress}

	host := config.GetEnv("EMAIL_HOST")
	port := "587"
	address := host + ":" + port

	subject := email.Subject
	body := email.Message

	message := []byte("To: " + email.To + "\r\n" +
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

	// golang note: since message is a byte array, there's no need to convert to a rune
	log.Printf("[INFO] Sent plaintext email to %s with msg length of %d\n", to, len(message))
}

func SendEmailTemplate(magic MagicLinkEmail) {
	config.LoadEnv()

	from := config.GetEnv("EMAIL_USERNAME")     // <your_email>@gmail.com
	password := config.GetEnv("EMAIL_PASSWORD") // google app password

	toEmailAddress := magic.To
	to := []string{toEmailAddress}

	host := config.GetEnv("EMAIL_HOST")
	port := "587"
	address := host + ":" + port

	// TEMPLATE HANDLING
	tmpl := template.Must(template.ParseFiles("internal/templates/magic-link.html"))
	var buf bytes.Buffer
	tmpl.Execute(&buf, struct {
		Token string
		Title string
		Url string
	}{
		Token: magic.Token,
		Title: config.GetEnv("VITE_SERVICE_NAME"),
		Url: config.GetEnv("VITE_SERVICE_NAME"),
	})
	// END TEMPLATE HANDLING

	// Future:
	subject := "Magic Login Link | " + config.GetEnv("VITE_SERVICE_NAME")

	message := []byte("To: " + magic.To + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		buf.String() + "\r\n")

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] Sent HTML email to %s with content length of %d\n", to, len(message))
}
