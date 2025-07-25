package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	analyticsServices "kyooar/internal/analytics/services"
	"kyooar/internal/analytics/models"
	organizationRepos "kyooar/internal/organization/repositories"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type TimeSeriesHandler struct {
	timeSeriesService analyticsServices.TimeSeriesService
	organizationRepo  organizationRepos.OrganizationRepository
}

func NewTimeSeriesHandler(i *do.Injector) (*TimeSeriesHandler, error) {
	return &TimeSeriesHandler{
		timeSeriesService: do.MustInvoke[analyticsServices.TimeSeriesService](i),
		organizationRepo:  do.MustInvoke[organizationRepos.OrganizationRepository](i),
	}, nil
}

// GetTimeSeries gets time series data for analytics
// @Summary Get time series analytics data
// @Description Get time series data for various metrics with customizable granularity and date range
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Param metric_types query []string true "Metric types to retrieve" collectionFormat(csv)
// @Param start_date query string true "Start date (ISO 8601)"
// @Param end_date query string true "End date (ISO 8601)"
// @Param granularity query string true "Data granularity (hourly, daily, weekly, monthly)"
// @Param product_id query string false "Filter by product ID"
// @Param question_id query string false "Filter by question ID"
// @Success 200 {object} models.TimeSeriesResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/organizations/{organizationId}/time-series [get]
func (h *TimeSeriesHandler) GetTimeSeries(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify organization ownership
	organization, err := h.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Parse query parameters
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")
	
	if startDateStr == "" || endDateStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "start_date and end_date are required")
	}
	
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid start_date format")
	}
	
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid end_date format")
	}
	
	granularity := c.QueryParam("granularity")
	if granularity == "" {
		granularity = models.GranularityDaily
	}
	
	// Validate granularity
	validGranularities := map[string]bool{
		models.GranularityHourly:  true,
		models.GranularityDaily:   true,
		models.GranularityWeekly:  true,
		models.GranularityMonthly: true,
	}
	if !validGranularities[granularity] {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid granularity")
	}
	
	// Parse metric types - handle both formats
	metricTypes := c.QueryParams()["metric_types"]
	
	// If metric_types is empty, try metric_types[] (common from frontend frameworks)
	if len(metricTypes) == 0 {
		metricTypes = c.QueryParams()["metric_types[]"]
	}
	
	if len(metricTypes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one metric_type is required")
	}
	
	// Build request
	request := models.TimeSeriesRequest{
		OrganizationID: organizationID,
		MetricTypes:    metricTypes,
		StartDate:      startDate,
		EndDate:        endDate,
		Granularity:    granularity,
	}
	
	// Parse optional filters
	if productIDStr := c.QueryParam("product_id"); productIDStr != "" {
		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid product_id")
		}
		request.ProductID = &productID
	}
	
	if questionIDStr := c.QueryParam("question_id"); questionIDStr != "" {
		questionID, err := uuid.Parse(questionIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid question_id")
		}
		request.QuestionID = &questionID
	}
	
	// Get time series data
	response, err := h.timeSeriesService.GetTimeSeries(ctx, request)
	if err != nil {
		logger.Error("Failed to get time series data", err, logrus.Fields{
			"organization_id": organizationID,
			"metric_types":    metricTypes,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get time series data")
	}
	
	return c.JSON(http.StatusOK, response)
}

// CompareTimePeriods compares analytics between two time periods
// @Summary Compare analytics between two time periods
// @Description Compare metrics between two different time periods to identify trends and changes
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Param body body models.ComparisonRequest true "Comparison request"
// @Success 200 {object} models.ComparisonResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/organizations/{organizationId}/compare [post]
func (h *TimeSeriesHandler) CompareTimePeriods(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify organization ownership
	organization, err := h.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Parse request body
	var request models.ComparisonRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	// Set organization ID from path
	request.OrganizationID = organizationID
	
	// Validate request
	if request.Period1Start.IsZero() || request.Period1End.IsZero() ||
		request.Period2Start.IsZero() || request.Period2End.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "All period dates are required")
	}
	
	if request.Period1Start.After(request.Period1End) || request.Period2Start.After(request.Period2End) {
		return echo.NewHTTPError(http.StatusBadRequest, "Start date must be before end date")
	}
	
	if len(request.MetricTypes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one metric type is required")
	}
	
	// Get comparison data
	response, err := h.timeSeriesService.GetComparison(ctx, request)
	if err != nil {
		logger.Error("Failed to get comparison data", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get comparison data")
	}
	
	return c.JSON(http.StatusOK, response)
}

// CollectMetrics manually triggers metric collection for an organization
// @Summary Collect metrics for an organization
// @Description Manually trigger the collection of time series metrics for analytics
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/organizations/{organizationId}/collect-metrics [post]
func (h *TimeSeriesHandler) CollectMetrics(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Verify organization ownership
	organization, err := h.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Collect metrics
	if err := h.timeSeriesService.CollectMetrics(ctx, organizationID); err != nil {
		logger.Error("Failed to collect metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to collect metrics")
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Metrics collected successfully",
	})
}