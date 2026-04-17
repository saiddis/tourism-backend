package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"tourism-backend/internal/domain"
)

type EmailService struct {
	host       string
	port       string
	username   string
	password   string
	adminEmail string
}

func NewEmailService() *EmailService {
	return &EmailService{
		host:       os.Getenv("SMTP_HOST"),
		port:       os.Getenv("SMTP_PORT"),
		username:   os.Getenv("SMTP_USER"),
		password:   os.Getenv("SMTP_PASSWORD"),
		adminEmail: os.Getenv("ADMIN_EMAIL"),
	}
}

func (s *EmailService) SendProviderApplicationEmail(userName, userEmail string, app *domain.ProviderApplication) error {
	if s.host == "" || s.adminEmail == "" {
		return nil
	}

	subject := fmt.Sprintf("New Provider Application from %s", userName)

	body := fmt.Sprintf(`New Provider Application

Name: %s
Email: %s
Phone: %s
Provider Type: %s
Years of Experience: %d

Bio:
%s

Social Links:
Instagram: %s
Telegram: %s
Facebook: %s

Application ID: %d
Submitted: %s
`,
		userName,
		userEmail,
		app.Phone,
		app.ProviderType,
		app.YearsExperience,
		nullableString(app.Bio),
		nullableString(app.InstagramURL),
		nullableString(app.TelegramURL),
		nullableString(app.FacebookURL),
		app.ID,
		app.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	return s.send(s.adminEmail, subject, body)
}

func (s *EmailService) SendApplicationStatusEmail(toEmail, status, adminNote string) error {
	if s.host == "" {
		return nil
	}

	subject := fmt.Sprintf("Your Provider Application Has Been %s", strings.Title(status))

	body := fmt.Sprintf(`Provider Application Update

Your provider application has been %s.

`, strings.ToLower(status))

	if adminNote != "" {
		body += fmt.Sprintf("Note from admin:\n%s\n", adminNote)
	}

	body += "\nBest regards,\nTourism Platform Team"

	return s.send(toEmail, subject, body)
}

func (s *EmailService) send(to, subject, body string) error {
	from := s.username

	headers := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n", from, to, subject)

	message := headers + body

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
	return err
}

func nullableString(s *string) string {
	if s == nil {
		return "Not provided"
	}
	return *s
}
