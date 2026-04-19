package email

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"tourism-backend/internal/domain"

	"github.com/mailersend/mailersend-go"
)

type EmailService struct {
	apiKey     string
	fromEmail  string
	fromName   string
	adminEmail string
	serverURL  string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey:     os.Getenv("MAILERSEND_API_KEY"),
		fromEmail:  os.Getenv("MAILERSEND_FROM_EMAIL"),
		fromName:   os.Getenv("MAILERSEND_FROM_NAME"),
		adminEmail: os.Getenv("ADMIN_EMAIL"),
		serverURL:  os.Getenv("SERVER_URL"),
	}
}

func (s *EmailService) SendProviderApplicationEmail(userName, userEmail string, app *domain.ProviderApplication) error {
	if s.apiKey == "" || s.adminEmail == "" {
		return nil
	}

	subject := fmt.Sprintf("New Provider Application from %s", userName)

	serverURL := s.serverURL
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}
	acceptURL := fmt.Sprintf("%s/provider-applications/%d/accept?token=%s", serverURL, app.ID, app.AdminToken)
	rejectURL := fmt.Sprintf("%s/provider-applications/%d/reject?token=%s", serverURL, app.ID, app.AdminToken)

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

func (s *EmailService) SendBookingConfirmationEmail(toEmail, userName, tourName string, startDate time.Time) error {
	if s.apiKey == "" {
		return nil
	}

	subject := "Your Tour Booking Has Been Confirmed!"

	body := fmt.Sprintf(`Hello %s,

Great news! Your booking for the tour "%s" has been confirmed!

Tour Details:
- Tour: %s
- Start Date: %s

Your spot is now secured. Please make sure to arrive on time for the tour.

If you have any questions, please don't hesitate to contact us.

Best regards,
Tourism Platform Team
`, userName, tourName, tourName, startDate.Format("2006-01-02"))

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
