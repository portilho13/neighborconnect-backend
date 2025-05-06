package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
)

type EmailConfig struct {
	SenderEmail    string
	SenderPassword string
	SmtpHost       string
	SmtpPort       string
}

type Email struct {
	To      []string
	Subject string
}

// SendEmail sends an HTML email
func SendEmail(email Email, email_type string, type_struct any) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config := EmailConfig{
		SenderEmail:    os.Getenv("SMTP_EMAIL"),
		SenderPassword: os.Getenv("STMP_PASSWORD"),
		SmtpHost:       "smtp.gmail.com",
		SmtpPort:       "587",
	}
	smtpServer := config.SmtpHost + ":" + config.SmtpPort

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + email.Subject + "\n"
	from := "From: " + config.SenderEmail + "\n"
	to := "To: " + strings.Join(email.To, ",") + "\n"

	var htmlContent string

	switch email_type {
	case "reward":
		if event, ok := type_struct.(models.Community_Event); ok {
			htmlContent = CreateRewardEmailTemplate(event)
		} else {
			return fmt.Errorf("invalid data type for reward email")
		}
	default:
	}
	message := []byte(from + to + subject + mime + htmlContent)

	auth := smtp.PlainAuth("", config.SenderEmail, config.SenderPassword, config.SmtpHost)

	err = smtp.SendMail(smtpServer, auth, config.SenderEmail, email.To, message)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}
