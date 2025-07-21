package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"lecritique/internal/feedback/models"
	feedbackRepos "lecritique/internal/feedback/repositories"
	qrcodeRepos "lecritique/internal/qrcode/repositories"
	restaurantRepos "lecritique/internal/restaurant/repositories"
	sharedModels "lecritique/internal/shared/models"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
)

type FeedbackService interface {
	Submit(ctx context.Context, feedback *models.Feedback) error
	GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error)
	GetByRestaurantIDWithFilters(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int, filters feedbackRepos.FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error)
	GetStats(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*FeedbackStats, error)
}

type feedbackService struct {
	feedbackRepo   feedbackRepos.FeedbackRepository
	restaurantRepo restaurantRepos.RestaurantRepository
	qrCodeRepo     qrcodeRepos.QRCodeRepository
}

type FeedbackStats struct {
	TotalFeedbacks    int64   `json:"total_feedbacks"`
	AverageRating     float64 `json:"average_rating"`
	FeedbacksToday    int64   `json:"feedbacks_today"`
	FeedbacksThisWeek int64   `json:"feedbacks_this_week"`
}

func NewFeedbackService(i *do.Injector) (FeedbackService, error) {
	return &feedbackService{
		feedbackRepo:   do.MustInvoke[feedbackRepos.FeedbackRepository](i),
		restaurantRepo: do.MustInvoke[restaurantRepos.RestaurantRepository](i),
		qrCodeRepo:     do.MustInvoke[qrcodeRepos.QRCodeRepository](i),
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

	feedback.RestaurantID = qrCode.RestaurantID

	feedback.OverallRating = s.calculateOverallRating(feedback.Responses)

	return s.feedbackRepo.Create(ctx, feedback)
}

func (s *feedbackService) GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.feedbackRepo.FindByRestaurantID(ctx, restaurantID, sharedModels.PageRequest{
		Page:  page,
		Limit: limit,
	})
}

func (s *feedbackService) GetByRestaurantIDWithFilters(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int, filters feedbackRepos.FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.feedbackRepo.FindByRestaurantIDWithFilters(ctx, restaurantID, sharedModels.PageRequest{
		Page:  page,
		Limit: limit,
	}, filters)
}

func (s *feedbackService) GetStats(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*FeedbackStats, error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)

	totalFeedbacks, _ := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Time{})
	feedbacksToday, _ := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, today)
	feedbacksThisWeek, _ := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, weekAgo)
	averageRating, _ := s.feedbackRepo.GetAverageRating(ctx, restaurantID, nil)

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
	
	
	for i, response := range responses {
		
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
