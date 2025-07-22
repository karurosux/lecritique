package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/feedback/models"
	"lecritique/internal/feedback/repositories"
	"lecritique/internal/feedback/services"
	"lecritique/internal/shared/logger"
	"lecritique/internal/shared/middleware"
	sharedModels "lecritique/internal/shared/models"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type FeedbackHandler struct {
	feedbackService services.FeedbackService
}

func NewFeedbackHandler(i *do.Injector) (*FeedbackHandler, error) {
	return &FeedbackHandler{
		feedbackService: do.MustInvoke[services.FeedbackService](i),
	}, nil
}

// GetByOrganization gets feedback for a organization with optional filters
// @Summary Get organization feedback with filters
// @Description Get all feedback for a specific organization with pagination and optional filters
// @Tags feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 20, max: 100)"
// @Param search query string false "Search in comments, customer name, or email"
// @Param rating_min query int false "Minimum rating (1-5)"
// @Param rating_max query int false "Maximum rating (1-5)"
// @Param date_from query string false "Start date (YYYY-MM-DD format)"
// @Param date_to query string false "End date (YYYY-MM-DD format)"
// @Param product_id query string false "Filter by specific product ID"
// @Param is_complete query boolean false "Filter by completion status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/feedback [get]
func (h *FeedbackHandler) GetByOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
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

	// Parse product filter
	if productIDStr := c.QueryParam("product_id"); productIDStr != "" {
		if productID, err := uuid.Parse(productIDStr); err == nil {
			filters.ProductID = &productID
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
		filters.DateFrom != nil || filters.DateTo != nil || filters.ProductID != nil || filters.IsComplete != nil

	var feedbacks interface{}
	if hasFilters {
		feedbacks, err = h.feedbackService.GetByOrganizationIDWithFilters(ctx, accountID, organizationID, page, limit, filters)
	} else {
		feedbacks, err = h.feedbackService.GetByOrganizationID(ctx, accountID, organizationID, page, limit)
	}

	if err != nil {
		logger.Error("Failed to get feedbacks", err, logrus.Fields{
			"account_id":    accountID,
			"organization_id": organizationID,
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

// GetStats gets feedback statistics for a organization
// @Summary Get feedback statistics
// @Description Get feedback analytics and statistics for a organization
// @Tags feedback
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/analytics [get]
func (h *FeedbackHandler) GetStats(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	accountID := middleware.GetResourceAccountID(c)

	stats, err := h.feedbackService.GetStats(ctx, accountID, organizationID)
	if err != nil {
		logger.Error("Failed to get feedback stats", err, logrus.Fields{
			"account_id":    accountID,
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get feedback statistics")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}
