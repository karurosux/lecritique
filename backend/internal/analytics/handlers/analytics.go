package handlers

import (
	"fmt"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/models"
	"net/http"
	"time"

	analyticsServices "kyooar/internal/analytics/services"
	feedbackRepos "kyooar/internal/feedback/repositories"
	menuRepos "kyooar/internal/menu/repositories"
	organizationRepos "kyooar/internal/organization/repositories"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type AnalyticsHandler struct {
	feedbackRepo     feedbackRepos.FeedbackRepository
	productRepo      menuRepos.ProductRepository
	organizationRepo organizationRepos.OrganizationRepository
	analyticsService analyticsServices.AnalyticsService
}

func NewAnalyticsHandler(i *do.Injector) (*AnalyticsHandler, error) {
	return &AnalyticsHandler{
		feedbackRepo:     do.MustInvoke[feedbackRepos.FeedbackRepository](i),
		productRepo:      do.MustInvoke[menuRepos.ProductRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
		analyticsService: do.MustInvoke[analyticsServices.AnalyticsService](i),
	}, nil
}

type ProductAnalytics struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	AverageRating float64   `json:"average_rating"`
	TotalFeedback int64     `json:"total_feedback"`
}

type OrganizationAnalytics struct {
	OrganizationID      uuid.UUID          `json:"organization_id"`
	OrganizationName    string             `json:"organization_name"`
	TotalFeedback       int64              `json:"total_feedback"`
	AverageRating       float64            `json:"average_rating"`
	FeedbackToday       int64              `json:"feedback_today"`
	FeedbackThisWeek    int64              `json:"feedback_this_week"`
	FeedbackThisMonth   int64              `json:"feedback_this_month"`
	TopRatedProducts    []ProductAnalytics `json:"top_rated_products"`
	LowestRatedProducts []ProductAnalytics `json:"lowest_rated_products"`
}

// GetOrganizationAnalytics gets analytics for a organization
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
func (h *AnalyticsHandler) GetOrganizationAnalytics(c echo.Context) error {
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

	// Get overall organization stats
	totalFeedback, _ := h.feedbackRepo.CountByOrganizationID(ctx, organizationID, time.Time{})
	feedbackToday, _ := h.feedbackRepo.CountByOrganizationID(ctx, organizationID, time.Now().Truncate(24*time.Hour))
	feedbackThisWeek, _ := h.feedbackRepo.CountByOrganizationID(ctx, organizationID, time.Now().AddDate(0, 0, -7))
	feedbackThisMonth, _ := h.feedbackRepo.CountByOrganizationID(ctx, organizationID, time.Now().AddDate(0, -1, 0))
	averageRating, err := h.feedbackRepo.GetAverageRating(ctx, organizationID, nil)
	if err != nil {
		logger.Error("Failed to get average rating", err, logrus.Fields{
			"organization_id": organizationID,
		})
	}

	// Debug logging
	logger.Info("Analytics Debug - BEFORE creating struct", logrus.Fields{
		"organization_id": organizationID,
		"total_feedback":  totalFeedback,
		"average_rating":  averageRating,
		"feedback_today":  feedbackToday,
	})

	analytics := OrganizationAnalytics{
		OrganizationID:    organizationID,
		OrganizationName:  organization.Name,
		TotalFeedback:     totalFeedback,
		AverageRating:     averageRating,
		FeedbackToday:     feedbackToday,
		FeedbackThisWeek:  feedbackThisWeek,
		FeedbackThisMonth: feedbackThisMonth,
	}

	logger.Info("Analytics Debug - AFTER creating struct", logrus.Fields{
		"analytics.AverageRating": analytics.AverageRating,
	})

	// Get product analytics using batch method to avoid N+1 queries
	products, err := h.productRepo.FindByOrganizationID(ctx, organizationID)
	if err == nil && len(products) > 0 {
		// Extract product IDs
		var productIDs []uuid.UUID
		for _, product := range products {
			productIDs = append(productIDs, product.ID)
		}

		// Get all product analytics in a single batch call
		productAnalyticsMap, err := h.analyticsService.GetProductAnalyticsBatch(ctx, organizationID, productIDs)
		if err != nil {
			logger.Error("Failed to get product analytics batch", err, logrus.Fields{
				"organization_id": organizationID,
			})
		} else {
			productAnalytics := make([]ProductAnalytics, 0, len(products))

			for _, analytics := range productAnalyticsMap {
				if analytics.TotalFeedback > 0 {
					productAnalytics = append(productAnalytics, ProductAnalytics{
						ProductID:     analytics.ProductID,
						ProductName:   analytics.ProductName,
						AverageRating: analytics.AverageRating,
						TotalFeedback: analytics.TotalFeedback,
					})
				}
			}

			// Sort and get top/bottom products
			if len(productAnalytics) > 0 {
				// Sort by rating (descending)
				for i := 0; i < len(productAnalytics)-1; i++ {
					for j := i + 1; j < len(productAnalytics); j++ {
						if productAnalytics[i].AverageRating < productAnalytics[j].AverageRating {
							productAnalytics[i], productAnalytics[j] = productAnalytics[j], productAnalytics[i]
						}
					}
				}

				// Get top 5
				topCount := 5
				if len(productAnalytics) < topCount {
					topCount = len(productAnalytics)
				}
				analytics.TopRatedProducts = productAnalytics[:topCount]

				// Get bottom 5
				bottomStart := len(productAnalytics) - 5
				if bottomStart < 0 {
					bottomStart = 0
				}
				if bottomStart < topCount {
					analytics.LowestRatedProducts = []ProductAnalytics{}
				} else {
					analytics.LowestRatedProducts = productAnalytics[bottomStart:]
				}
			}
		}
	}

	logger.Info("Analytics Debug - FINAL before response", logrus.Fields{
		"analytics.AverageRating": analytics.AverageRating,
		"full_analytics":          analytics,
	})

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    analytics,
	})
}

// GetProductAnalytics gets analytics for a specific product
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
func (h *AnalyticsHandler) GetProductAnalytics(c echo.Context) error {
	ctx := c.Request().Context()

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Get product and verify ownership
	product, err := h.productRepo.FindByID(ctx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := h.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil || organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get product stats
	totalFeedback, _ := h.feedbackRepo.CountByProductID(ctx, productID)
	averageRating, _ := h.feedbackRepo.GetAverageRating(ctx, product.OrganizationID, &productID)

	// Get recent feedback
	recentFeedback, err := h.feedbackRepo.FindByProductID(ctx, productID, models.PageRequest{Page: 1, Limit: 10})
	if err != nil {
		logger.Error("Failed to get recent feedback", err, logrus.Fields{
			"product_id": productID,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
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

// GetDashboardMetrics gets basic analytics metrics for dashboard
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
func (h *AnalyticsHandler) GetDashboardMetrics(c echo.Context) error {
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

	// Get dashboard metrics
	metrics, err := h.analyticsService.GetDashboardMetrics(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get dashboard metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get metrics")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success": true,
		"data":    metrics,
	})
}

// GetProductInsights gets detailed insights for a specific product
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
func (h *AnalyticsHandler) GetProductInsights(c echo.Context) error {
	ctx := c.Request().Context()

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid product ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	// Get product and verify ownership
	product, err := h.productRepo.FindByID(ctx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := h.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil || organization.AccountID != resourceAccountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	// Get product insights
	insights, err := h.analyticsService.GetProductInsights(ctx, productID)
	if err != nil {
		logger.Error("Failed to get product insights", err, logrus.Fields{
			"product_id": productID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get insights")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    insights,
	})
}

// GetOrganizationChartData gets aggregated chart data for organization analytics
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
func (h *AnalyticsHandler) GetOrganizationChartData(c echo.Context) error {
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

	// Build filters from query parameters
	filters := make(map[string]interface{})
	if dateFrom := c.QueryParam("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}
	if dateTo := c.QueryParam("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}
	if productID := c.QueryParam("product_id"); productID != "" {
		filters["product_id"] = productID
	}

	// Get chart data
	logger.Info("Getting organization chart data", logrus.Fields{
		"organization_id":     organizationID,
		"filters":             filters,
		"resource_account_id": resourceAccountID,
	})

	chartData, err := h.analyticsService.GetOrganizationChartData(ctx, organizationID, filters)
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    chartData,
	})
}
