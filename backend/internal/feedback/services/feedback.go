package services

import (
	"context"
	"fmt"
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

	// Calculate overall rating from numeric question responses
	feedback.OverallRating = s.calculateOverallRating(feedback.Responses)

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

func (s *feedbackService) GetByRestaurantIDWithFilters(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, page, limit int, filters feedbackRepos.FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error) {
	// Verify ownership
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

// calculateOverallRating computes an overall rating from individual question responses
func (s *feedbackService) calculateOverallRating(responses models.Responses) int {
	var totalScore float64
	var count int
	
	// Debug logging
	fmt.Printf("Calculating overall rating from %d responses\n", len(responses))
	
	for i, response := range responses {
		fmt.Printf("Response %d: Answer=%v, Type=%T\n", i, response.Answer, response.Answer)
		
		// Only consider numeric responses for overall rating calculation
		switch v := response.Answer.(type) {
		case float64:
			// Normalize different scales to 1-5 range
			normalizedScore := s.normalizeScore(v)
			totalScore += normalizedScore
			count++
			fmt.Printf("  Added float64: %v -> normalized: %v\n", v, normalizedScore)
		case int:
			// Convert int to float64 and normalize
			normalizedScore := s.normalizeScore(float64(v))
			totalScore += normalizedScore
			count++
			fmt.Printf("  Added int: %v -> normalized: %v\n", v, normalizedScore)
		default:
			fmt.Printf("  Skipped non-numeric response: %T\n", v)
		}
	}
	
	fmt.Printf("Total score: %v, Count: %d\n", totalScore, count)
	
	if count == 0 {
		fmt.Printf("No numeric responses found, returning 0\n")
		return 0 // No numeric responses
	}
	
	// Calculate average and round to nearest integer (1-5 range)
	average := totalScore / float64(count)
	fmt.Printf("Average before clamping: %v\n", average)
	
	// Ensure result is within 1-5 range
	if average < 1 {
		fmt.Printf("Clamping to 1 (was %v)\n", average)
		return 1
	} else if average > 5 {
		fmt.Printf("Clamping to 5 (was %v)\n", average)
		return 5
	}
	
	// Round to nearest integer
	result := int(average + 0.5)
	fmt.Printf("Final overall rating: %d\n", result)
	return result
}

// normalizeScore converts different rating scales to 1-5 range
func (s *feedbackService) normalizeScore(score float64) float64 {
	// Detect scale and normalize to 1-5
	if score <= 5 {
		// Already 1-5 scale, ensure minimum of 1
		if score < 1 {
			return 1
		}
		return score
	} else if score <= 10 {
		// 1-10 scale, convert to 1-5
		return ((score - 1) / 9 * 4) + 1
	} else if score <= 100 {
		// 1-100 or percentage scale, convert to 1-5
		return ((score - 1) / 99 * 4) + 1
	}
	
	// Unknown scale, assume it's already correct but clamp to 5
	return 5
}
