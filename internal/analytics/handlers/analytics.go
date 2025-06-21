package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	feedbackRepo   feedbackRepos.FeedbackRepository
	dishRepo       menuRepos.DishRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewAnalyticsHandler(
	feedbackRepo feedbackRepos.FeedbackRepository,
	dishRepo menuRepos.DishRepository,
	restaurantRepo restaurantRepos.RestaurantRepository,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		feedbackRepo:   feedbackRepo,
		dishRepo:       dishRepo,
		restaurantRepo: restaurantRepo,
	}
}

type DishAnalytics struct {
	DishID        uuid.UUID `json:"dish_id"`
	DishName      string    `json:"dish_name"`
	AverageRating float64   `json:"average_rating"`
	TotalFeedback int64     `json:"total_feedback"`
}

type RestaurantAnalytics struct {
	RestaurantID      uuid.UUID       `json:"restaurant_id"`
	RestaurantName    string          `json:"restaurant_name"`
	TotalFeedback     int64           `json:"total_feedback"`
	AverageRating     float64         `json:"average_rating"`
	FeedbackToday     int64           `json:"feedback_today"`
	FeedbackThisWeek  int64           `json:"feedback_this_week"`
	FeedbackThisMonth int64           `json:"feedback_this_month"`
	TopRatedDishes    []DishAnalytics `json:"top_rated_dishes"`
	LowestRatedDishes []DishAnalytics `json:"lowest_rated_dishes"`
}

func (h *AnalyticsHandler) GetRestaurantAnalytics(c echo.Context) error {
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	// Verify restaurant ownership
	restaurant, err := h.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}
	if restaurant.AccountID != accountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get overall restaurant stats
	totalFeedback, _ := h.feedbackRepo.CountByRestaurantID(restaurantID, time.Time{})
	feedbackToday, _ := h.feedbackRepo.CountByRestaurantID(restaurantID, time.Now().Truncate(24*time.Hour))
	feedbackThisWeek, _ := h.feedbackRepo.CountByRestaurantID(restaurantID, time.Now().AddDate(0, 0, -7))
	feedbackThisMonth, _ := h.feedbackRepo.CountByRestaurantID(restaurantID, time.Now().AddDate(0, -1, 0))
	averageRating, _ := h.feedbackRepo.GetAverageRating(restaurantID, nil)

	analytics := RestaurantAnalytics{
		RestaurantID:      restaurantID,
		RestaurantName:    restaurant.Name,
		TotalFeedback:     totalFeedback,
		AverageRating:     averageRating,
		FeedbackToday:     feedbackToday,
		FeedbackThisWeek:  feedbackThisWeek,
		FeedbackThisMonth: feedbackThisMonth,
	}

	// Get dish analytics
	dishes, err := h.dishRepo.FindByRestaurantID(restaurantID)
	if err == nil && len(dishes) > 0 {
		dishAnalytics := make([]DishAnalytics, 0, len(dishes))

		for _, dish := range dishes {
			avgRating, _ := h.feedbackRepo.GetAverageRating(restaurantID, &dish.ID)
			count, _ := h.feedbackRepo.CountByDishID(dish.ID)

			if count > 0 {
				dishAnalytics = append(dishAnalytics, DishAnalytics{
					DishID:        dish.ID,
					DishName:      dish.Name,
					AverageRating: avgRating,
					TotalFeedback: count,
				})
			}
		}

		// Sort and get top/bottom dishes
		if len(dishAnalytics) > 0 {
			// Sort by rating (descending)
			for i := 0; i < len(dishAnalytics)-1; i++ {
				for j := i + 1; j < len(dishAnalytics); j++ {
					if dishAnalytics[i].AverageRating < dishAnalytics[j].AverageRating {
						dishAnalytics[i], dishAnalytics[j] = dishAnalytics[j], dishAnalytics[i]
					}
				}
			}

			// Get top 5
			topCount := 5
			if len(dishAnalytics) < topCount {
				topCount = len(dishAnalytics)
			}
			analytics.TopRatedDishes = dishAnalytics[:topCount]

			// Get bottom 5
			bottomStart := len(dishAnalytics) - 5
			if bottomStart < 0 {
				bottomStart = 0
			}
			if bottomStart < topCount {
				analytics.LowestRatedDishes = []DishAnalytics{}
			} else {
				analytics.LowestRatedDishes = dishAnalytics[bottomStart:]
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analytics,
	})
}

func (h *AnalyticsHandler) GetDishAnalytics(c echo.Context) error {
	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	// Get dish and verify ownership
	dish, err := h.dishRepo.FindByID(dishID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	restaurant, err := h.restaurantRepo.FindByID(dish.RestaurantID)
	if err != nil || restaurant.AccountID != accountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get dish stats
	totalFeedback, _ := h.feedbackRepo.CountByDishID(dishID)
	averageRating, _ := h.feedbackRepo.GetAverageRating(dish.RestaurantID, &dishID)

	// Get recent feedback
	recentFeedback, err := h.feedbackRepo.FindByDishID(dishID, models.PageRequest{Page: 1, Limit: 10})
	if err != nil {
		logger.Error("Failed to get recent feedback", err, logrus.Fields{
			"dish_id": dishID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"dish_id":         dishID,
			"dish_name":       dish.Name,
			"total_feedback":  totalFeedback,
			"average_rating":  averageRating,
			"recent_feedback": recentFeedback.Data,
		},
	})
}

