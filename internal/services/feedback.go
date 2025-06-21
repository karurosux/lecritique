package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/models"
	"github.com/lecritique/api/internal/repositories"
)

type FeedbackService interface {
	Submit(feedback *models.Feedback) error
	GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*repositories.PageResponse[models.Feedback], error)
	GetStats(accountID uuid.UUID, restaurantID uuid.UUID) (*FeedbackStats, error)
}

type feedbackService struct {
	feedbackRepo   repositories.FeedbackRepository
	restaurantRepo repositories.RestaurantRepository
	qrCodeRepo     repositories.QRCodeRepository
}

type FeedbackStats struct {
	TotalFeedbacks    int64   `json:"total_feedbacks"`
	AverageRating     float64 `json:"average_rating"`
	FeedbacksToday    int64   `json:"feedbacks_today"`
	FeedbacksThisWeek int64   `json:"feedbacks_this_week"`
}

func NewFeedbackService(feedbackRepo repositories.FeedbackRepository, restaurantRepo repositories.RestaurantRepository, qrCodeRepo repositories.QRCodeRepository) FeedbackService {
	return &feedbackService{
		feedbackRepo:   feedbackRepo,
		restaurantRepo: restaurantRepo,
		qrCodeRepo:     qrCodeRepo,
	}
}

func (s *feedbackService) Submit(feedback *models.Feedback) error {
	// Validate QR code
	qrCode, err := s.qrCodeRepo.FindByID(feedback.QRCodeID)
	if err != nil {
		return err
	}

	if !qrCode.IsValid() {
		return repositories.ErrRecordNotFound
	}

	// Set restaurant ID from QR code
	feedback.RestaurantID = qrCode.RestaurantID

	// Submit feedback
	return s.feedbackRepo.Create(feedback)
}

func (s *feedbackService) GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID, page, limit int) (*repositories.PageResponse[models.Feedback], error) {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, repositories.ErrRecordNotFound
	}

	return s.feedbackRepo.FindByRestaurantID(restaurantID, repositories.PageRequest{
		Page:  page,
		Limit: limit,
	})
}

func (s *feedbackService) GetStats(accountID uuid.UUID, restaurantID uuid.UUID) (*FeedbackStats, error) {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, repositories.ErrRecordNotFound
	}

	// Get stats
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekAgo := today.AddDate(0, 0, -7)

	totalFeedbacks, _ := s.feedbackRepo.CountByRestaurantID(restaurantID, time.Time{})
	feedbacksToday, _ := s.feedbackRepo.CountByRestaurantID(restaurantID, today)
	feedbacksThisWeek, _ := s.feedbackRepo.CountByRestaurantID(restaurantID, weekAgo)
	averageRating, _ := s.feedbackRepo.GetAverageRating(restaurantID, nil)

	return &FeedbackStats{
		TotalFeedbacks:    totalFeedbacks,
		AverageRating:     averageRating,
		FeedbacksToday:    feedbacksToday,
		FeedbacksThisWeek: feedbacksThisWeek,
	}, nil
}
