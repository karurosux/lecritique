package services

import (
	"context"
	"lecritique/internal/feedback/models"
	"time"

	feedbackRepos "lecritique/internal/feedback/repositories"
	organizationRepos "lecritique/internal/organization/repositories"
	qrcodeRepos "lecritique/internal/qrcode/repositories"
	sharedModels "lecritique/internal/shared/models"
	sharedRepos "lecritique/internal/shared/repositories"

	"github.com/google/uuid"
	"github.com/samber/do"
)

type FeedbackService interface {
	Submit(ctx context.Context, feedback *models.Feedback) error
	GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error)
	GetByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackRepos.FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error)
	GetStats(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*FeedbackStats, error)
}

type feedbackService struct {
	feedbackRepo     feedbackRepos.FeedbackRepository
	organizationRepo organizationRepos.OrganizationRepository
	qrCodeRepo       qrcodeRepos.QRCodeRepository
}

type FeedbackStats struct {
	TotalFeedbacks    int64   `json:"total_feedbacks"`
	AverageRating     float64 `json:"average_rating"`
	FeedbacksToday    int64   `json:"feedbacks_today"`
	FeedbacksThisWeek int64   `json:"feedbacks_this_week"`
}

func NewFeedbackService(i *do.Injector) (FeedbackService, error) {
	return &feedbackService{
		feedbackRepo:     do.MustInvoke[feedbackRepos.FeedbackRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
		qrCodeRepo:       do.MustInvoke[qrcodeRepos.QRCodeRepository](i),
	}, nil
}

func (s *feedbackService) Submit(ctx context.Context, feedback *models.Feedback) error {
	qrCode, err := s.qrCodeRepo.FindByID(ctx, feedback.QRCodeID)
	if err != nil {
		return err
	}

	if !qrCode.IsValid() {
		return sharedRepos.ErrRecordNotFound
	}

	feedback.OrganizationID = qrCode.OrganizationID

	feedback.OverallRating = s.calculateOverallRating(feedback.Responses)

	return s.feedbackRepo.Create(ctx, feedback)
}

func (s *feedbackService) GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.feedbackRepo.FindByOrganizationID(ctx, organizationID, sharedModels.PageRequest{
		Page:  page,
		Limit: limit,
	})
}

func (s *feedbackService) GetByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackRepos.FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.feedbackRepo.FindByOrganizationIDWithFilters(ctx, organizationID, sharedModels.PageRequest{
		Page:  page,
		Limit: limit,
	}, filters)
}

func (s *feedbackService) GetStats(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*FeedbackStats, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)

	totalFeedbacks, _ := s.feedbackRepo.CountByOrganizationID(ctx, organizationID, time.Time{})
	feedbacksToday, _ := s.feedbackRepo.CountByOrganizationID(ctx, organizationID, today)
	feedbacksThisWeek, _ := s.feedbackRepo.CountByOrganizationID(ctx, organizationID, weekAgo)
	averageRating, _ := s.feedbackRepo.GetAverageRating(ctx, organizationID, nil)

	return &FeedbackStats{
		TotalFeedbacks:    totalFeedbacks,
		AverageRating:     averageRating,
		FeedbacksToday:    feedbacksToday,
		FeedbacksThisWeek: feedbacksThisWeek,
	}, nil
}

// calculateOverallRating computes an overall rating from individual question responses
func (s *feedbackService) calculateOverallRating(responses models.Responses) int {
	var totalScore float64
	var count int

	for _, response := range responses {
		switch v := response.Answer.(type) {
		case float64:
			normalizedScore := s.normalizeScore(v)
			totalScore += normalizedScore
			count++
		case int:
			normalizedScore := s.normalizeScore(float64(v))
			totalScore += normalizedScore
			count++
		default:
		}
	}

	if count == 0 {
		return 0 // No numeric responses
	}

	average := totalScore / float64(count)

	if average < 1 {
		return 1
	} else if average > 5 {
		return 5
	}

	result := int(average + 0.5)
	return result
}

func (s *feedbackService) normalizeScore(score float64) float64 {
	if score <= 5 {
		if score < 1 {
			return 1
		}
		return score
	} else if score <= 10 {
		return ((score - 1) / 9 * 4) + 1
	} else if score <= 100 {
		return ((score - 1) / 99 * 4) + 1
	}

	// Unknown scale, assume it's already correct but clamp to 5
	return 5
}
