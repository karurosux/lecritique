package organizationservice

import (
	"context"

	"github.com/google/uuid"
	organizationinterface "kyooar/internal/organization/interface"
	organizationmodel "kyooar/internal/organization/model"
	"kyooar/internal/shared/errors"
	sharedRepos "kyooar/internal/shared/repositories"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type organizationService struct {
	organizationRepo   organizationinterface.OrganizationRepository
	subscriptionRepo   subscriptioninterface.SubscriptionRepository
}

func NewOrganizationService(
	organizationRepo organizationinterface.OrganizationRepository,
	subscriptionRepo subscriptioninterface.SubscriptionRepository,
) organizationinterface.OrganizationService {
	return &organizationService{
		organizationRepo:   organizationRepo,
		subscriptionRepo:   subscriptionRepo,
	}
}

func (s *organizationService) Create(ctx context.Context, accountID uuid.UUID, organization *organizationmodel.Organization) error {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
	}

	currentCount, err := s.organizationRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to verify organization count", 500)
	}

	if subscription != nil && !subscription.CanAddOrganization(int(currentCount)) {
		return errors.NewWithDetails("SUBSCRIPTION_LIMIT",
			"Organization limit exceeded for current subscription plan",
			402,
			map[string]interface{}{
				"current_count": currentCount,
				"max_allowed":   subscription.Plan.MaxOrganizations,
			})
	}

	organization.AccountID = accountID
	if err := s.organizationRepo.Create(ctx, organization); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to create organization", 500)
	}

	return nil
}

func (s *organizationService) Update(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, updates map[string]interface{}) error {
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
		case "address":
			if v, ok := value.(string); ok {
				organization.Address = v
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

func (s *organizationService) GetByID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*organizationmodel.Organization, error) {
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

func (s *organizationService) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]organizationmodel.Organization, error) {
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

func (s *organizationService) GetByIDForAnalytics(ctx context.Context, organizationID uuid.UUID) (*organizationmodel.Organization, error) {
	return s.organizationRepo.FindByID(ctx, organizationID)
}
