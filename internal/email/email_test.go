package email

import (
	"flag"
	"os"
	"testing"
	"time"

	"tourism-backend/internal/domain"
)

var (
	apiKey       = flag.String("mailersend-api-key", "", "MailerSend API key")
	fromEmail    = flag.String("mailersend-from-email", "", "From email address")
	fromName     = flag.String("mailersend-from-name", "", "From name")
	toEmail      = flag.String("to-email", "", "Recipient email for testing")
	adminEmail   = flag.String("admin-email", "", "Admin email to receive notification")
	status       = flag.String("status", "accepted", "Application status (accepted/rejected)")
	adminNote    = flag.String("admin-note", "", "Admin note")
	userName     = flag.String("user-name", "Test User", "Applicant name")
	userEmail    = flag.String("user-email", "test@example.com", "Applicant email")
	phone        = flag.String("phone", "+992000000000", "Phone number")
	providerType = flag.String("provider-type", "private", "Provider type (private/organization)")
	instagram    = flag.String("instagram", "https://instagram.com/test", "Instagram URL")
	telegram     = flag.String("telegram", "https://t.me/test", "Telegram URL")
	facebook     = flag.String("facebook", "https://facebook.com/test", "Facebook URL")
	bio          = flag.String("bio", "Test bio for the provider application", "Bio")
)

func TestSendApplicationStatusEmail(t *testing.T) {

	flag.Parse()

	if *apiKey == "" {
		t.Skip("Skipping test: -mailersend-api-key not provided")
	}
	if *fromEmail == "" {
		t.Skip("Skipping test: -mailersend-from-email not provided")
	}
	if *fromName == "" {
		t.Skip("Skipping test: -mailersend-from-name not provided")
	}
	if *toEmail == "" {
		t.Skip("Skipping test: -to-email not provided")
	}

	os.Setenv("MAILERSEND_API_KEY", *apiKey)
	os.Setenv("MAILERSEND_FROM_EMAIL", *fromEmail)
	os.Setenv("MAILERSEND_FROM_NAME", *fromName)

	svc := NewEmailService()

	t.Logf("Sending test email to: %s", *toEmail)
	t.Logf("Status: %s", *status)
	if *adminNote != "" {
		t.Logf("Admin note: %s", *adminNote)
	}

	err := svc.SendApplicationStatusEmail(*toEmail, *status, *adminNote)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("Email sent successfully")
}

func TestSendProviderApplicationEmail(t *testing.T) {
	flag.Parse()

	if *apiKey == "" {
		t.Skip("Skipping test: -mailersend-api-key not provided")
	}
	if *fromEmail == "" {
		t.Skip("Skipping test: -mailersend-from-email not provided")
	}
	if *fromName == "" {
		t.Skip("Skipping test: -mailersend-from-name not provided")
	}
	if *adminEmail == "" {
		t.Skip("Skipping test: -admin-email not provided")
	}

	os.Setenv("MAILERSEND_API_KEY", *apiKey)
	os.Setenv("MAILERSEND_FROM_EMAIL", *fromEmail)
	os.Setenv("MAILERSEND_FROM_NAME", *fromName)
	os.Setenv("ADMIN_EMAIL", *adminEmail)

	svc := NewEmailService()

	instagramPtr := instagram
	if *instagram == "" {
		instagramPtr = nil
	}
	telegramPtr := telegram
	if *telegram == "" {
		telegramPtr = nil
	}
	facebookPtr := facebook
	if *facebook == "" {
		facebookPtr = nil
	}
	bioPtr := bio
	if *bio == "" {
		bioPtr = nil
	}

	app := &domain.ProviderApplication{
		ID:              999,
		Phone:           *phone,
		ProviderType:    domain.ProviderType(*providerType),
		InstagramURL:    instagramPtr,
		TelegramURL:     telegramPtr,
		FacebookURL:     facebookPtr,
		YearsExperience: 5,
		Bio:             bioPtr,
		CreatedAt:       time.Now(),
	}

	t.Logf("Sending provider application notification to admin: %s", *adminEmail)
	t.Logf("Applicant: %s (%s)", *userName, *userEmail)

	err := svc.SendProviderApplicationEmail(*userName, *userEmail, app)
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	t.Log("Email sent successfully")
}

func TestServiceInitialization(t *testing.T) {
	t.Run("with environment variables", func(t *testing.T) {
		os.Setenv("MAILERSEND_API_KEY", "test-key")
		os.Setenv("MAILERSEND_FROM_EMAIL", "test@example.com")
		os.Setenv("MAILERSEND_FROM_NAME", "Test Name")
		os.Setenv("ADMIN_EMAIL", "admin@example.com")
		defer func() {
			os.Unsetenv("MAILERSEND_API_KEY")
			os.Unsetenv("MAILERSEND_FROM_EMAIL")
			os.Unsetenv("MAILERSEND_FROM_NAME")
			os.Unsetenv("ADMIN_EMAIL")
		}()

		svc := NewEmailService()

		if svc.apiKey != "test-key" {
			t.Errorf("Expected apiKey 'test-key', got '%s'", svc.apiKey)
		}
		if svc.fromEmail != "test@example.com" {
			t.Errorf("Expected fromEmail 'test@example.com', got '%s'", svc.fromEmail)
		}
		if svc.fromName != "Test Name" {
			t.Errorf("Expected fromName 'Test Name', got '%s'", svc.fromName)
		}
		if svc.adminEmail != "admin@example.com" {
			t.Errorf("Expected adminEmail 'admin@example.com', got '%s'", svc.adminEmail)
		}
	})

	t.Run("with empty environment variables", func(t *testing.T) {
		os.Unsetenv("MAILERSEND_API_KEY")
		os.Unsetenv("MAILERSEND_FROM_EMAIL")
		os.Unsetenv("MAILERSEND_FROM_NAME")
		os.Unsetenv("ADMIN_EMAIL")

		svc := NewEmailService()

		if svc.apiKey != "" {
			t.Errorf("Expected empty apiKey, got '%s'", svc.apiKey)
		}
		if svc.fromEmail != "" {
			t.Errorf("Expected empty fromEmail, got '%s'", svc.fromEmail)
		}
	})
}

func TestNullableString(t *testing.T) {
	t.Run("nil string", func(t *testing.T) {
		result := nullableString(nil)
		expected := "Not provided"
		if result != expected {
			t.Errorf("Expected '%s', got '%s'", expected, result)
		}
	})

	t.Run("valid string", func(t *testing.T) {
		str := "test string"
		result := nullableString(&str)
		if result != str {
			t.Errorf("Expected '%s', got '%s'", str, result)
		}
	})
}

func TestSendSkipsWithEmptyApiKey(t *testing.T) {
	os.Setenv("MAILERSEND_API_KEY", "")
	os.Setenv("MAILERSEND_FROM_EMAIL", "test@example.com")
	os.Setenv("MAILERSEND_FROM_NAME", "Test")
	defer func() {
		os.Unsetenv("MAILERSEND_API_KEY")
		os.Unsetenv("MAILERSEND_FROM_EMAIL")
		os.Unsetenv("MAILERSEND_FROM_NAME")
	}()

	svc := NewEmailService()

	err := svc.SendApplicationStatusEmail("test@example.com", "accepted", "")
	if err != nil {
		t.Errorf("Expected no error (silent skip), got: %v", err)
	}
}
