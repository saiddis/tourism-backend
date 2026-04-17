package service

import (
	"context"
	"errors"
	"tourism-backend/internal/domain"
	"tourism-backend/internal/repository"
)

type ProviderService struct {
	providerRepo repository.ProviderRepository
	userRepo     repository.UserRepository
}

func NewProviderService(providerRepo repository.ProviderRepository, userRepo repository.UserRepository) *ProviderService {
	return &ProviderService{
		providerRepo: providerRepo,
		userRepo:     userRepo,
	}
}

func (s *ProviderService) Create(ctx context.Context, providerData *domain.Provider) (*domain.Provider, error) {
	provider := &domain.Provider{
		UserID:          providerData.UserID,
		Phone:           providerData.Phone,
		ProviderType:    providerData.ProviderType,
		InstagramURL:    providerData.InstagramURL,
		TelegramURL:     providerData.TelegramURL,
		FacebookURL:     providerData.FacebookURL,
		YearsExperience: providerData.YearsExperience,
		Bio:             providerData.Bio,
	}

	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, err
	}

	if err := s.userRepo.UpdateRole(ctx, provider.UserID, domain.RoleProvider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *ProviderService) GetByID(ctx context.Context, id int) (*domain.Provider, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, errors.New("provider not found")
	}
	return provider, nil
}

func (s *ProviderService) GetByUserID(ctx context.Context, userID int) (*domain.Provider, error) {
	return s.providerRepo.GetByUserID(ctx, userID)
}

func (s *ProviderService) GetActiveProviders(ctx context.Context) ([]*domain.Provider, error) {
	return s.providerRepo.GetActive(ctx)
}

func (s *ProviderService) GetAll(ctx context.Context) ([]*domain.Provider, error) {
	return s.providerRepo.GetAll(ctx)
}

func (s *ProviderService) Update(ctx context.Context, provider *domain.Provider) (*domain.Provider, error) {
	existing, err := s.providerRepo.GetByID(ctx, provider.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("provider not found")
	}

	provider.UserID = existing.UserID
	if err := s.providerRepo.Update(ctx, provider); err != nil {
		return nil, err
	}

	return s.providerRepo.GetByID(ctx, provider.ID)
}

func (s *ProviderService) ToggleActive(ctx context.Context, id int, active bool) error {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if provider == nil {
		return errors.New("provider not found")
	}
	return s.providerRepo.UpdateActive(ctx, id, active)
}
