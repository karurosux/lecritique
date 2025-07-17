package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/feedback/models"
	"github.com/lecritique/api/internal/feedback/repositories"
	"github.com/lecritique/api/internal/feedback/services"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/lecritique/api/internal/shared/middleware"
	sharedModels "github.com/lecritique/api/internal/shared/models"
	"github.com/sirupsen/logrus"
)

type FeedbackHandler struct {
	feedbackService services.FeedbackService
}

func NewFeedbackHandler(feedbackService services.FeedbackService) *FeedbackHandler {
	return &FeedbackHandler{
		feedbackService: feedbackService,
	}
}

// GetByRestaurant gets feedback for a restaurant with optional filters
// @Summary Get restaurant feedback with filters
// @Description Get all feedback for a specific restaurant with pagination and optional filters
// @Tags feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param search query string false "Search in comments, customer name, or email"
// @Param rating_min query int false "Minimum rating (1-5)"
// @Param rating_max query int false "Maximum rating (1-5)"
// @Param date_from query string false "Start date (YYYY-MM-DD format)"
// @Param date_to query string false "End date (YYYY-MM-DD format)"
// @Param dish_id query string false "Filter by specific dish ID"
// @Param is_complete query boolean false "Filter by completion status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/feedback [get]
func (h *FeedbackHandler) GetByRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID := middleware.GetResourceAccountID(c)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Parse filter parameters
	filters := repositories.FeedbackFilter{
		Search: c.QueryParam("search"),
	}

	// Parse rating filters
	if ratingMinStr := c.QueryParam("rating_min"); ratingMinStr != "" {
		if ratingMin, err := strconv.Atoi(ratingMinStr); err == nil && ratingMin >= 1 && ratingMin <= 5 {
			filters.RatingMin = &ratingMin
		}
	}
	if ratingMaxStr := c.QueryParam("rating_max"); ratingMaxStr != "" {
		if ratingMax, err := strconv.Atoi(ratingMaxStr); err == nil && ratingMax >= 1 && ratingMax <= 5 {
			filters.RatingMax = &ratingMax
		}
	}

	// Parse date filters
	if dateFromStr := c.QueryParam("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filters.DateFrom = &dateFrom
		}
	}
	if dateToStr := c.QueryParam("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			filters.DateTo = &dateTo
		}
	}

	// Parse dish filter
	if dishIDStr := c.QueryParam("dish_id"); dishIDStr != "" {
		if dishID, err := uuid.Parse(dishIDStr); err == nil {
			filters.DishID = &dishID
		}
	}

	// Parse completion filter
	if isCompleteStr := c.QueryParam("is_complete"); isCompleteStr != "" {
		if isComplete, err := strconv.ParseBool(isCompleteStr); err == nil {
			filters.IsComplete = &isComplete
		}
	}

	// Use filtered service method if any filters are provided
	hasFilters := filters.Search != "" || filters.RatingMin != nil || filters.RatingMax != nil || 
		filters.DateFrom != nil || filters.DateTo != nil || filters.DishID != nil || filters.IsComplete != nil

	var feedbacks interface{}
	if hasFilters {
		feedbacks, err = h.feedbackService.GetByRestaurantIDWithFilters(ctx, accountID, restaurantID, page, limit, filters)
	} else {
		feedbacks, err = h.feedbackService.GetByRestaurantID(ctx, accountID, restaurantID, page, limit)
	}

	if err != nil {
		logger.Error("Failed to get feedbacks", err, logrus.Fields{
			"account_id":    accountID,
			"restaurant_id": restaurantID,
			"page":          page,
			"limit":         limit,
			"filters":       filters,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get feedbacks")
	}

	// Type assertion for response
	if response, ok := feedbacks.(*sharedModels.PageResponse[models.Feedback]); ok {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    response.Data,
			"meta": map[string]interface{}{
				"total":        response.Total,
				"page":         response.Page,
				"limit":        response.Limit,
				"total_pages":  response.TotalPages,
			},
		})
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to process response")
}

// GetStats gets feedback statistics for a restaurant
// @Summary Get feedback statistics
// @Description Get feedback analytics and statistics for a restaurant
// @Tags feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/analytics [get]
func (h *FeedbackHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID := middleware.GetResourceAccountID(c)

	stats, err := h.feedbackService.GetStats(ctx, accountID, restaurantID)
	if err != nil {
		logger.Error("Failed to get feedback stats", err, logrus.Fields{
			"account_id":    accountID,
			"restaurant_id": restaurantID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get feedback statistics")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}