package feedbackinterface

import (
	"context"
	"time"

	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
	sharedModels "kyooar/internal/shared/models"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *feedbackmodel.Feedback) error
	FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Feedback, error)
	Update(ctx context.Context, feedback *feedbackmodel.Feedback) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[feedbackmodel.Feedback], error)
	FindByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackmodel.FeedbackFilter) (*sharedModels.PageResponse[feedbackmodel.Feedback], error)
	GetStatsByOrganization(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*feedbackmodel.FeedbackStats, error)
	FindByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error)
	FindByQuestionInPeriod(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time) ([]feedbackmodel.Feedback, error)
	CountByOrganizationID(ctx context.Context, organizationID uuid.UUID, since time.Time) (int64, error)
	CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error)
	CountByQRCodeID(ctx context.Context, qrCodeID uuid.UUID) (int64, error)
	CountByQRCodeIDs(ctx context.Context, qrCodeIDs []uuid.UUID) (map[uuid.UUID]int64, error)
	GetAverageRating(ctx context.Context, organizationID uuid.UUID, productID *uuid.UUID) (float64, error)
	FindByProductID(ctx context.Context, productID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[feedbackmodel.Feedback], error)
	FindByProductIDForAnalytics(ctx context.Context, productID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error)
	GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]feedbackmodel.Question, error)
}

type FeedbackService interface {
	Submit(ctx context.Context, feedback *feedbackmodel.Feedback) error
	GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[feedbackmodel.Feedback], error)
	GetByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackmodel.FeedbackFilter) (*sharedModels.PageResponse[feedbackmodel.Feedback], error)
	GetStats(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*feedbackmodel.FeedbackStats, error)
	GetByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error)
	GetByQuestionInPeriod(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time) ([]feedbackmodel.Feedback, error)
}