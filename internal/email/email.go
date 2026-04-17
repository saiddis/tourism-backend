package email

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mailersend/mailersend-go"
	"tourism-backend/internal/domain"
)

type EmailService struct {
	apiKey     string
	fromEmail  string
	fromName   string
	adminEmail string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey:     os.Getenv("MAILERSEND_API_KEY"),
		fromEmail:  os.Getenv("MAILERSEND_FROM_EMAIL"),
		fromName:   os.Getenv("MAILERSEND_FROM_NAME"),
		adminEmail: os.Getenv("ADMIN_EMAIL"),
	}
}

func (s *EmailService) SendProviderApplicationEmail(userName, userEmail string, app *domain.ProviderApplication) error {
	if s.apiKey == "" || s.adminEmail == "" {
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
	if s.apiKey == "" {
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
	ms := mailersend.NewMailersend(s.apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := ms.Email.NewMessage()

	message.SetFrom(mailersend.From{
		Name:  s.fromName,
		Email: s.fromEmail,
	})

	message.SetRecipients([]mailersend.Recipient{
		{
			Email: to,
		},
	})

	message.SetSubject(subject)
	message.SetText(body)

	_, err := ms.Email.Send(ctx, message)
	return err
}

func nullableString(s *string) string {
	if s == nil {
		return "Not provided"
	}
	return *s
}
