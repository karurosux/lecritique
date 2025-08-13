package analyticscontroller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	analyticsconstants "kyooar/internal/analytics/constants"
	analyticsinterface "kyooar/internal/analytics/interface"
	analyticsmodel "kyooar/internal/analytics/model"
	feedbackinterface "kyooar/internal/feedback/interface"
	productRepos "kyooar/internal/product/repositories"
	organizationinterface "kyooar/internal/organization/interface"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/models"
	"kyooar/internal/shared/response"

	"github.com/sirupsen/logrus"
)

type AnalyticsController struct {
	feedbackRepo     feedbackinterface.FeedbackRepository
	productRepo      productRepos.ProductRepository
	organizationRepo organizationinterface.OrganizationRepository
	analyticsService analyticsinterface.AnalyticsService
}

func NewAnalyticsController(
	feedbackRepo feedbackinterface.FeedbackRepository,
	productRepo productRepos.ProductRepository,
	organizationRepo organizationinterface.OrganizationRepository,
	analyticsService analyticsinterface.AnalyticsService,
) *AnalyticsController {
	return &AnalyticsController{
		feedbackRepo:     feedbackRepo,
		productRepo:      productRepo,
		organizationRepo: organizationRepo,
		analyticsService: analyticsService,
	}
}

func (c *AnalyticsController) handleError(ctx echo.Context, err error) error {
	switch err.Error() {
	case analyticsconstants.ErrInvalidOrganizationID:
		return response.Error(ctx, errors.ErrBadRequest)
	case analyticsconstants.ErrInvalidProductID:
		return response.Error(ctx, errors.ErrBadRequest)
	case analyticsconstants.ErrOrganizationNotFound:
		return response.Error(ctx, errors.ErrNotFound)
	case analyticsconstants.ErrProductNotFound:
		return response.Error(ctx, errors.ErrNotFound)
	case analyticsconstants.ErrAccessDenied:
		return response.Error(ctx, errors.ErrForbidden)
	case analyticsconstants.ErrInvalidDateRange:
		return response.Error(ctx, errors.ErrBadRequest)
	case analyticsconstants.ErrInvalidGranularity:
		return response.Error(ctx, errors.ErrBadRequest)
	case analyticsconstants.ErrInvalidMetricType:
		return response.Error(ctx, errors.ErrBadRequest)
	case analyticsconstants.ErrMetricsNotFound:
		return response.Error(ctx, errors.ErrNotFound)
	case analyticsconstants.ErrFailedToGetMetrics:
		return response.Error(ctx, errors.ErrInternalServer)
	case analyticsconstants.ErrFailedToCollectMetrics:
		return response.Error(ctx, errors.ErrInternalServer)
	default:
		return response.Error(ctx, errors.ErrInternalServer)
	}
}

// @Summary Get organization analytics
// @Description Get comprehensive analytics data for a organization including ratings, feedback counts, and product performance
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/organizations/{organizationId} [get]
func (c *AnalyticsController) GetOrganizationAnalytics(ctx echo.Context) error {
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

	totalFeedback, _ := c.feedbackRepo.CountByOrganizationID(requestCtx, organizationID, time.Time{})
	feedbackToday, _ := c.feedbackRepo.CountByOrganizationID(requestCtx, organizationID, time.Now().Truncate(24*time.Hour))
	feedbackThisWeek, _ := c.feedbackRepo.CountByOrganizationID(requestCtx, organizationID, time.Now().AddDate(0, 0, -7))
	feedbackThisMonth, _ := c.feedbackRepo.CountByOrganizationID(requestCtx, organizationID, time.Now().AddDate(0, -1, 0))
	averageRating, err := c.feedbackRepo.GetAverageRating(requestCtx, organizationID, nil)
	if err != nil {
		logger.Error("Failed to get average rating", err, logrus.Fields{
			"organization_id": organizationID,
		})
	}

	analytics := analyticsmodel.OrganizationAnalytics{
		OrganizationID:    organizationID,
		OrganizationName:  organization.Name,
		TotalFeedback:     totalFeedback,
		AverageRating:     averageRating,
		FeedbackToday:     feedbackToday,
		FeedbackThisWeek:  feedbackThisWeek,
		FeedbackThisMonth: feedbackThisMonth,
	}

	products, err := c.productRepo.FindByOrganizationID(requestCtx, organizationID)
	if err == nil && len(products) > 0 {
		var productIDs []uuid.UUID
		for _, product := range products {
			productIDs = append(productIDs, product.ID)
		}

		productAnalyticsMap, err := c.analyticsService.GetProductAnalyticsBatch(requestCtx, organizationID, productIDs)
		if err != nil {
			logger.Error("Failed to get product analytics batch", err, logrus.Fields{
				"organization_id": organizationID,
			})
		} else {
			productAnalytics := make([]analyticsmodel.ProductAnalytics, 0, len(products))

			for _, analyticsData := range productAnalyticsMap {
				if analyticsData.TotalFeedback > 0 {
					productAnalytics = append(productAnalytics, analyticsmodel.ProductAnalytics{
						ProductID:     analyticsData.ProductID,
						ProductName:   analyticsData.ProductName,
						AverageRating: analyticsData.AverageRating,
						TotalFeedback: analyticsData.TotalFeedback,
					})
				}
			}

			if len(productAnalytics) > 0 {
				for i := 0; i < len(productAnalytics)-1; i++ {
					for j := i + 1; j < len(productAnalytics); j++ {
						if productAnalytics[i].AverageRating < productAnalytics[j].AverageRating {
							productAnalytics[i], productAnalytics[j] = productAnalytics[j], productAnalytics[i]
						}
					}
				}

				topCount := 5
				if len(productAnalytics) < topCount {
					topCount = len(productAnalytics)
				}
				analytics.TopRatedProducts = productAnalytics[:topCount]

				bottomStart := len(productAnalytics) - 5
				if bottomStart < 0 {
					bottomStart = 0
				}
				if bottomStart < topCount {
					analytics.LowestRatedProducts = []analyticsmodel.ProductAnalytics{}
				} else {
					analytics.LowestRatedProducts = productAnalytics[bottomStart:]
				}
			}
		}
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analytics,
	})
}

// @Summary Get product analytics
// @Description Get detailed analytics data for a specific product including ratings, feedback count, and recent feedback
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/products/{productId} [get]
func (c *AnalyticsController) GetProductAnalytics(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	productID, err := uuid.Parse(ctx.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidProductID)
	}

	resourceAccountID := middleware.GetResourceAccountID(ctx)

	product, err := c.productRepo.FindByID(requestCtx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, analyticsconstants.ErrProductNotFound)
	}

	organization, err := c.organizationRepo.FindByID(requestCtx, product.OrganizationID)
	if err != nil || organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, analyticsconstants.ErrAccessDenied)
	}

	totalFeedback, _ := c.feedbackRepo.CountByProductID(requestCtx, productID)
	averageRating, _ := c.feedbackRepo.GetAverageRating(requestCtx, product.OrganizationID, &productID)

	recentFeedback, err := c.feedbackRepo.FindByProductID(requestCtx, productID, models.PageRequest{Page: 1, Limit: 10})
	if err != nil {
		logger.Error("Failed to get recent feedback", err, logrus.Fields{
			"product_id": productID,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"product_id":      productID,
			"product_name":    product.Name,
			"total_feedback":  totalFeedback,
			"average_rating":  averageRating,
			"recent_feedback": recentFeedback.Data,
		},
	})
}

// @Summary Get dashboard metrics
// @Description Get basic analytics metrics for the dashboard including satisfaction, recommendation rate, and recent feedback
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/dashboard/{organizationId} [get]
func (c *AnalyticsController) GetDashboardMetrics(ctx echo.Context) error {
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

	metrics, err := c.analyticsService.GetDashboardMetrics(requestCtx, organizationID)
	if err != nil {
		logger.Error("Failed to get dashboard metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, analyticsconstants.ErrFailedToGetMetrics)
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    metrics,
	})
}

// @Summary Get product insights
// @Description Get detailed insights for a specific product including question-level analytics
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param productId path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/products/{productId}/insights [get]
func (c *AnalyticsController) GetProductInsights(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()

	productID, err := uuid.Parse(ctx.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, analyticsconstants.ErrInvalidProductID)
	}

	resourceAccountID := middleware.GetResourceAccountID(ctx)

	product, err := c.productRepo.FindByID(requestCtx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, analyticsconstants.ErrProductNotFound)
	}

	organization, err := c.organizationRepo.FindByID(requestCtx, product.OrganizationID)
	if err != nil || organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, analyticsconstants.ErrAccessDenied)
	}

	insights, err := c.analyticsService.GetProductInsights(requestCtx, productID)
	if err != nil {
		logger.Error("Failed to get product insights", err, logrus.Fields{
			"product_id": productID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, analyticsconstants.ErrFailedToGetMetrics)
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    insights,
	})
}

// @Summary Get organization chart data
// @Description Get pre-aggregated chart data for all questions in a organization with optional filters
// @Tags analytics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param product_id query string false "Filter by specific product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/analytics/organizations/{organizationId}/charts [get]
func (c *AnalyticsController) GetOrganizationChartData(ctx echo.Context) error {
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

	filters := make(map[string]interface{})
	if dateFrom := ctx.QueryParam("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := ctx.QueryParam("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if productID := ctx.QueryParam("product_id"); productID != "" {
		filters["product_id"] = productID
	}

	logger.Info("Getting organization chart data", logrus.Fields{
		"organization_id":     organizationID,
		"filters":             filters,
		"resource_account_id": resourceAccountID,
	})

	chartData, err := c.analyticsService.GetOrganizationChartData(requestCtx, organizationID, filters)
	if err != nil {
		logger.Error("Failed to get organization chart data", err, logrus.Fields{
			"organization_id":     organizationID,
			"filters":             filters,
			"resource_account_id": resourceAccountID,
			"error_type":          fmt.Sprintf("%T", err),
			"error_message":       err.Error(),
		})
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to get chart data: %v", err))
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    chartData,
	})
}