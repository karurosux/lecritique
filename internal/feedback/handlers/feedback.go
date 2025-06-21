package handlers

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/feedback/services"
	"github.com/lecritique/api/internal/shared/logger"
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

func (h *FeedbackHandler) GetByRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	feedbacks, err := h.feedbackService.GetByRestaurantID(ctx, accountID, restaurantID, page, limit)
	if err != nil {
		logger.Error("Failed to get feedbacks", err, logrus.Fields{
			"account_id":    accountID,
			"restaurant_id": restaurantID,
			"page":          page,
			"limit":         limit,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get feedbacks")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    feedbacks.Data,
		"meta": map[string]interface{}{
			"total":        feedbacks.Total,
			"page":         feedbacks.Page,
			"limit":        feedbacks.Limit,
			"total_pages":  feedbacks.TotalPages,
		},
	})
}

func (h *FeedbackHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

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