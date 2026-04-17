package service

import (
	"context"
	"errors"
	"log"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/email"
	"tourism-backend/internal/repository"
)

type ProviderApplicationService struct {
	appRepo      repository.ProviderApplicationRepository
	providerRepo repository.ProviderRepository
	userRepo     repository.UserRepository
	emailSvc     *email.EmailService
}

func NewProviderApplicationService(
	appRepo repository.ProviderApplicationRepository,
	providerRepo repository.ProviderRepository,
	userRepo repository.UserRepository,
	emailSvc *email.EmailService,
) *ProviderApplicationService {
	return &ProviderApplicationService{
		appRepo:      appRepo,
		providerRepo: providerRepo,
		userRepo:     userRepo,
		emailSvc:     emailSvc,
	}
}

func (s *ProviderApplicationService) Submit(ctx context.Context, userID int, data *domain.ProviderApplication) (*domain.ProviderApplication, error) {
	existingProvider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingProvider != nil {
		return nil, errors.New("you are already a provider")
	}

	existingApp, err := s.appRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingApp != nil && existingApp.Status == domain.ApplicationStatusPending {
		return nil, errors.New("you already have a pending application")
	}

	app := &domain.ProviderApplication{
		UserID:          userID,
		Phone:           data.Phone,
		ProviderType:    data.ProviderType,
		InstagramURL:    data.InstagramURL,
		TelegramURL:     data.TelegramURL,
		FacebookURL:     data.FacebookURL,
		YearsExperience: data.YearsExperience,
		Bio:             data.Bio,
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user != nil && s.emailSvc != nil {
		go func() {
			if err := s.emailSvc.SendProviderApplicationEmail(user.Name, user.Email, app); err != nil {
				log.Printf("Failed to send provider application email: %v", err)
			}
		}()
	}

	return app, nil
}

func (s *ProviderApplicationService) GetByID(ctx context.Context, id int) (*domain.ProviderApplication, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("application not found")
	}
	return app, nil
}

func (s *ProviderApplicationService) GetByUserID(ctx context.Context, userID int) (*domain.ProviderApplication, error) {
	return s.appRepo.GetByUserID(ctx, userID)
}

func (s *ProviderApplicationService) GetAll(ctx context.Context) ([]*domain.ProviderApplication, error) {
	return s.appRepo.GetAll(ctx)
}

func (s *ProviderApplicationService) Accept(ctx context.Context, id int, adminNote string) (*domain.ProviderApplication, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("application not found")
	}
	if app.Status != domain.ApplicationStatusPending {
		return nil, errors.New("application is not pending")
	}

	provider := &domain.Provider{
		UserID:          app.UserID,
		Phone:           app.Phone,
		ProviderType:    app.ProviderType,
		InstagramURL:    app.InstagramURL,
		TelegramURL:     app.TelegramURL,
		FacebookURL:     app.FacebookURL,
		YearsExperience: app.YearsExperience,
		Bio:             app.Bio,
	}

	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateRole(ctx, app.UserID, domain.RoleProvider); err != nil {
		return nil, err
	}

	var note *string
	if adminNote != "" {
		note = &adminNote
	}
	if err := s.appRepo.UpdateStatus(ctx, id, domain.ApplicationStatusAccepted, note); err != nil {
		return nil, err
	}

	user, _ := s.userRepo.GetByID(ctx, app.UserID)
	if user != nil && s.emailSvc != nil {
		go func() {
			if err := s.emailSvc.SendApplicationStatusEmail(user.Email, "accepted", adminNote); err != nil {
				log.Printf("Failed to send acceptance email: %v", err)
			}
		}()
	}

	return s.appRepo.GetByID(ctx, id)
}

func (s *ProviderApplicationService) Reject(ctx context.Context, id int, adminNote string) (*domain.ProviderApplication, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, errors.New("application not found")
	}
	if app.Status != domain.ApplicationStatusPending {
		return nil, errors.New("application is not pending")
	}

	var note *string
	if adminNote != "" {
		note = &adminNote
	}
	if err := s.appRepo.UpdateStatus(ctx, id, domain.ApplicationStatusRejected, note); err != nil {
		return nil, err
	}

	user, _ := s.userRepo.GetByID(ctx, app.UserID)
	if user != nil && s.emailSvc != nil {
		go func() {
			if err := s.emailSvc.SendApplicationStatusEmail(user.Email, "rejected", adminNote); err != nil {
				log.Printf("Failed to send rejection email: %v", err)
			}
		}()
	}

	return s.appRepo.GetByID(ctx, id)
}
