package services

import (
	"context"

	"github.com/google/uuid"
	"kyooar/internal/organization/models"
	organizationRepos "kyooar/internal/organization/repositories"
	"kyooar/internal/shared/errors"
	sharedRepos "kyooar/internal/shared/repositories"
	subscriptionRepos "kyooar/internal/subscription/repositories"
	"github.com/samber/do"
)

type OrganizationService interface {
	Create(ctx context.Context, accountID uuid.UUID, organization *models.Organization) error
	Update(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) error
	GetByID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*models.Organization, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Organization, error)
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
}

type organizationService struct {
	organizationRepo   organizationRepos.OrganizationRepository
	subscriptionRepo subscriptionRepos.SubscriptionRepository
}

func NewOrganizationService(i *do.Injector) (OrganizationService, error) {
	return &organizationService{
		organizationRepo:   do.MustInvoke[organizationRepos.OrganizationRepository](i),
		subscriptionRepo: do.MustInvoke[subscriptionRepos.SubscriptionRepository](i),
	}, nil
}

func (s *organizationService) Create(ctx context.Context, accountID uuid.UUID, organization *models.Organization) error {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		// For now, allow creation without subscription (development mode)
		// TODO: Re-enable subscription check in production
		// return errors.New("SUBSCRIPTION_REQUIRED", "No active subscription found for account", 402)
	}

	currentCount, err := s.organizationRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to verify organization count", 500)
	}

	// Only check limits if subscription exists
	if subscription != nil && !subscription.CanAddOrganization(int(currentCount)) {
		return errors.NewWithDetails("SUBSCRIPTION_LIMIT",
			"Organization limit exceeded for current subscription plan",
			402,
			map[string]interface{}{
				"current_count": currentCount,
				"max_allowed":   subscription.Plan.MaxOrganizations,
			})
	}

	// Set account ID and create
	organization.AccountID = accountID
	if err := s.organizationRepo.Create(ctx, organization); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to create organization", 500)
	}

	return nil
}

func (s *organizationService) Update(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, updates map[string]interface{}) error {
	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return errors.NotFound("Organization")
		}
		return errors.ErrDatabaseOperation
	}

	if organization.AccountID != accountID {
		return errors.Forbidden("update this organization")
	}

	// Update fields
	for key, value := range updates {
		switch key {
		case "name":
			if v, ok := value.(string); ok {
				organization.Name = v
			}
		case "description":
			if v, ok := value.(string); ok {
				organization.Description = v
			}
		case "phone":
			if v, ok := value.(string); ok {
				organization.Phone = v
			}
		case "email":
			if v, ok := value.(string); ok {
				organization.Email = v
			}
		case "website":
			if v, ok := value.(string); ok {
				organization.Website = v
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				organization.IsActive = v
			}
		}
	}

	if err := s.organizationRepo.Update(ctx, organization); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to update organization", 500)
	}

	return nil
}

func (s *organizationService) Delete(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) error {
	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return errors.NotFound("Organization")
		}
		return errors.ErrDatabaseOperation
	}

	if organization.AccountID != accountID {
		return errors.Forbidden("delete this organization")
	}

	if err := s.organizationRepo.Delete(ctx, organizationID); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to delete organization", 500)
	}

	return nil
}

func (s *organizationService) GetByID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*models.Organization, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return nil, errors.NotFound("Organization")
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to fetch organization", 500)
	}

	if organization.AccountID != accountID {
		return nil, errors.Forbidden("access this organization")
	}

	return organization, nil
}

func (s *organizationService) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Organization, error) {
	organizations, err := s.organizationRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Unable to retrieve organizations", 500)
	}
	return organizations, nil
}

func (s *organizationService) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	count, err := s.organizationRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return 0, errors.Wrap(err, "DATABASE_ERROR", "Unable to count organizations", 500)
	}
	return count, nil
}
