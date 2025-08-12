package analyticscontroller

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	analyticsconstants "kyooar/internal/analytics/constants"
	analyticsinterface "kyooar/internal/analytics/interface"
	models "kyooar/internal/analytics/model"
	organizationRepos "kyooar/internal/organization/repositories"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"

	"github.com/sirupsen/logrus"
)

type TimeSeriesController struct {
	timeSeriesService analyticsinterface.TimeSeriesService
	organizationRepo  organizationRepos.OrganizationRepository
}

func NewTimeSeriesController(
	timeSeriesService analyticsinterface.TimeSeriesService,
	organizationRepo organizationRepos.OrganizationRepository,
) *TimeSeriesController {
	return &TimeSeriesController{
		timeSeriesService: timeSeriesService,
		organizationRepo:  organizationRepo,
	}
}

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
func (c *TimeSeriesController) GetTimeSeries(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	
	organizationID, err := uuid.Parse(ctx.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidOrganizationID)
	}

	resourceAccountID := middleware.GetResourceAccountID(ctx)

	organization, err := c.organizationRepo.FindByID(requestCtx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, analyticsconstants.ErrOrganizationNotFound)
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, analyticsconstants.ErrAccessDenied)
	}

	startDateStr := ctx.QueryParam("start_date")
	endDateStr := ctx.QueryParam("end_date")
	
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
	
	granularity := ctx.QueryParam("granularity")
	if granularity == "" {
		granularity = models.GranularityDaily
	}
	
	validGranularities := map[string]bool{
		models.GranularityHourly:  true,
		models.GranularityDaily:   true,
		models.GranularityWeekly:  true,
		models.GranularityMonthly: true,
	}
	if !validGranularities[granularity] {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidGranularity)
	}
	
	metricTypes := ctx.QueryParams()["metric_types"]
	
	if len(metricTypes) == 0 {
		metricTypes = ctx.QueryParams()["metric_types[]"]
	}
	
	if len(metricTypes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one metric_type is required")
	}
	
	request := models.TimeSeriesRequest{
		OrganizationID: organizationID,
		MetricTypes:    metricTypes,
		StartDate:      startDate,
		EndDate:        endDate,
		Granularity:    granularity,
	}
	
	if productIDStr := ctx.QueryParam("product_id"); productIDStr != "" {
		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidProductID)
		}
		request.ProductID = &productID
	}
	
	if questionIDStr := ctx.QueryParam("question_id"); questionIDStr != "" {
		questionID, err := uuid.Parse(questionIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid question_id")
		}
		request.QuestionID = &questionID
	}
	
	response, err := c.timeSeriesService.GetTimeSeries(requestCtx, request)
	if err != nil {
		logger.Error("Failed to get time series data", err, logrus.Fields{
			"organization_id": organizationID,
			"metric_types":    metricTypes,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, analyticsconstants.ErrFailedToGetMetrics)
	}
	
	return ctx.JSON(http.StatusOK, response)
}

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
func (c *TimeSeriesController) CompareTimePeriods(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	
	organizationID, err := uuid.Parse(ctx.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidOrganizationID)
	}

	resourceAccountID := middleware.GetResourceAccountID(ctx)

	organization, err := c.organizationRepo.FindByID(requestCtx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, analyticsconstants.ErrOrganizationNotFound)
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, analyticsconstants.ErrAccessDenied)
	}

	var request models.ComparisonRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	
	request.OrganizationID = organizationID
	
	if request.Period1Start.IsZero() || request.Period1End.IsZero() ||
		request.Period2Start.IsZero() || request.Period2End.IsZero() {
		return echo.NewHTTPError(http.StatusBadRequest, "All period dates are required")
	}
	
	if request.Period1Start.After(request.Period1End) || request.Period2Start.After(request.Period2End) {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidDateRange)
	}
	
	if len(request.MetricTypes) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "At least one metric type is required")
	}
	
	response, err := c.timeSeriesService.GetComparison(requestCtx, request)
	if err != nil {
		logger.Error("Failed to get comparison data", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, analyticsconstants.ErrFailedToGetMetrics)
	}
	
	return ctx.JSON(http.StatusOK, response)
}

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
func (c *TimeSeriesController) CollectMetrics(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	
	organizationID, err := uuid.Parse(ctx.Param("organizationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidOrganizationID)
	}

	resourceAccountID := middleware.GetResourceAccountID(ctx)

	organization, err := c.organizationRepo.FindByID(requestCtx, organizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, analyticsconstants.ErrOrganizationNotFound)
	}
	if organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, analyticsconstants.ErrAccessDenied)
	}

	if err := c.timeSeriesService.CollectMetrics(requestCtx, organizationID); err != nil {
		logger.Error("Failed to collect metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, analyticsconstants.ErrFailedToCollectMetrics)
	}
	
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Metrics collected successfully",
	})
}