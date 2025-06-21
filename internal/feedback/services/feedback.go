package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/feedback/models"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	sharedModels "github.com/lecritique/api/internal/shared/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
)

type FeedbackService interface {
	Submit(ctx context.Context, feedback *models.Feedback) error
	GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error)
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

func NewFeedbackService(feedbackRepo feedbackRepos.FeedbackRepository, restaurantRepo restaurantRepos.RestaurantRepository, qrCodeRepo qrcodeRepos.QRCodeRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo:   feedbackRepo,
		restaurantRepo: restaurantRepo,
		qrCodeRepo:     qrCodeRepo,
	}
}

func (s *feedbackService) Submit(ctx context.Context, feedback *models.Feedback) error {
	// Validate QR code
	qrCode, err := s.qrCodeRepo.FindByID(ctx, feedback.QRCodeID)
	if err != nil {
		return err
	}

	if !qrCode.IsValid() {
		return sharedRepos.ErrRecordNotFound
	}

	// Set restaurant ID from QR code
	feedback.RestaurantID = qrCode.RestaurantID

	// Submit feedback
	return s.feedbackRepo.Create(ctx, feedback)
}

func (s *feedbackService) GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*sharedModels.PageResponse[models.Feedback], error) {
	// Verify ownership
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

func (s *feedbackService) GetStats(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*FeedbackStats, error) {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	// Get stats
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
