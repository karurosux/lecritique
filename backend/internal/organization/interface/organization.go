package organizationinterface

import (
	"context"

	"github.com/google/uuid"
	organizationmodel "kyooar/internal/organization/model"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization *organizationmodel.Organization) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*organizationmodel.Organization, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]organizationmodel.Organization, error)
	Update(ctx context.Context, organization *organizationmodel.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
}

type OrganizationService interface {
	Create(ctx context.Context, accountID uuid.UUID, organization *organizationmodel.Organization) error
	Update(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) error
	GetByID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*organizationmodel.Organization, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]organizationmodel.Organization, error)
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
	GetByIDForAnalytics(ctx context.Context, organizationID uuid.UUID) (*organizationmodel.Organization, error)
}