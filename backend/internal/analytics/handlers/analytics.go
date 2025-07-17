package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	analyticsServices "github.com/lecritique/api/internal/analytics/services"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/lecritique/api/internal/shared/middleware"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	feedbackRepo     feedbackRepos.FeedbackRepository
	dishRepo         menuRepos.DishRepository
	restaurantRepo   restaurantRepos.RestaurantRepository
	analyticsService analyticsServices.AnalyticsService
}

func NewAnalyticsHandler(
	feedbackRepo feedbackRepos.FeedbackRepository,
	dishRepo menuRepos.DishRepository,
	restaurantRepo restaurantRepos.RestaurantRepository,
	analyticsService analyticsServices.AnalyticsService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		feedbackRepo:     feedbackRepo,
		dishRepo:         dishRepo,
		restaurantRepo:   restaurantRepo,
		analyticsService: analyticsService,
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

// GetRestaurantAnalytics gets analytics for a restaurant
// @Summary Get restaurant analytics
// @Description Get comprehensive analytics data for a restaurant including ratings, feedback counts, and dish performance
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/restaurants/{restaurantId} [get]
func (h *AnalyticsHandler) GetRestaurantAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify restaurant ownership
	restaurant, err := h.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}
	if restaurant.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get overall restaurant stats
	totalFeedback, _ := h.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Time{})
	feedbackToday, _ := h.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Now().Truncate(24*time.Hour))
	feedbackThisWeek, _ := h.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Now().AddDate(0, 0, -7))
	feedbackThisMonth, _ := h.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Now().AddDate(0, -1, 0))
	averageRating, err := h.feedbackRepo.GetAverageRating(ctx, restaurantID, nil)
	if err != nil {
		logger.Error("Failed to get average rating", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	}
	
	// Debug logging
	logger.Info("Analytics Debug - BEFORE creating struct", logrus.Fields{
		"restaurant_id":    restaurantID,
		"total_feedback":   totalFeedback,
		"average_rating":   averageRating,
		"feedback_today":   feedbackToday,
	})

	analytics := RestaurantAnalytics{
		RestaurantID:      restaurantID,
		RestaurantName:    restaurant.Name,
		TotalFeedback:     totalFeedback,
		AverageRating:     averageRating,
		FeedbackToday:     feedbackToday,
		FeedbackThisWeek:  feedbackThisWeek,
		FeedbackThisMonth: feedbackThisMonth,
	}
	
	logger.Info("Analytics Debug - AFTER creating struct", logrus.Fields{
		"analytics.AverageRating": analytics.AverageRating,
	})

	// Get dish analytics
	dishes, err := h.dishRepo.FindByRestaurantID(ctx, restaurantID)
	if err == nil && len(dishes) > 0 {
		dishAnalytics := make([]DishAnalytics, 0, len(dishes))

		for _, dish := range dishes {
			avgRating, _ := h.feedbackRepo.GetAverageRating(ctx, restaurantID, &dish.ID)
			count, _ := h.feedbackRepo.CountByDishID(ctx, dish.ID)

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

	logger.Info("Analytics Debug - FINAL before response", logrus.Fields{
		"analytics.AverageRating": analytics.AverageRating,
		"full_analytics": analytics,
	})
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analytics,
	})
}

// GetDishAnalytics gets analytics for a specific dish
// @Summary Get dish analytics
// @Description Get detailed analytics data for a specific dish including ratings, feedback count, and recent feedback
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dishId path string true "Dish ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/dishes/{dishId} [get]
func (h *AnalyticsHandler) GetDishAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	
	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Get dish and verify ownership
	dish, err := h.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	restaurant, err := h.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil || restaurant.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get dish stats
	totalFeedback, _ := h.feedbackRepo.CountByDishID(ctx, dishID)
	averageRating, _ := h.feedbackRepo.GetAverageRating(ctx, dish.RestaurantID, &dishID)

	// Get recent feedback
	recentFeedback, err := h.feedbackRepo.FindByDishID(ctx, dishID, models.PageRequest{Page: 1, Limit: 10})
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

// GetDashboardMetrics gets basic analytics metrics for dashboard
// @Summary Get dashboard metrics
// @Description Get basic analytics metrics for the dashboard including satisfaction, recommendation rate, and recent feedback
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/dashboard/{restaurantId} [get]
func (h *AnalyticsHandler) GetDashboardMetrics(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify restaurant ownership
	restaurant, err := h.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}
	if restaurant.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get dashboard metrics
	metrics, err := h.analyticsService.GetDashboardMetrics(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get dashboard metrics", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get metrics")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    metrics,
	})
}

// GetDishInsights gets detailed insights for a specific dish
// @Summary Get dish insights
// @Description Get detailed insights for a specific dish including question-level analytics
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dishId path string true "Dish ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/dishes/{dishId}/insights [get]
func (h *AnalyticsHandler) GetDishInsights(c echo.Context) error {
	ctx := c.Request().Context()
	
	dishID, err := uuid.Parse(c.Param("dishId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid dish ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Get dish and verify ownership
	dish, err := h.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	restaurant, err := h.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil || restaurant.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get dish insights
	insights, err := h.analyticsService.GetDishInsights(ctx, dishID)
	if err != nil {
		logger.Error("Failed to get dish insights", err, logrus.Fields{
			"dish_id": dishID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get insights")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    insights,
	})
}

// GetRestaurantChartData gets aggregated chart data for restaurant analytics
// @Summary Get restaurant chart data
// @Description Get pre-aggregated chart data for all questions in a restaurant with optional filters
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param dish_id query string false "Filter by specific dish ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/restaurants/{restaurantId}/charts [get]
func (h *AnalyticsHandler) GetRestaurantChartData(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify restaurant ownership
	restaurant, err := h.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}
	if restaurant.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Build filters from query parameters
	filters := make(map[string]interface{})
	if dateFrom := c.QueryParam("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.QueryParam("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if dishID := c.QueryParam("dish_id"); dishID != "" {
		filters["dish_id"] = dishID
	}

	// Get chart data
	chartData, err := h.analyticsService.GetRestaurantChartData(ctx, restaurantID, filters)
	if err != nil {
		logger.Error("Failed to get restaurant chart data", err, logrus.Fields{
			"restaurant_id": restaurantID,
			"filters": filters,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get chart data")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    chartData,
	})
}

