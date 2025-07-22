package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
	"kyooar/internal/subscription/models"
	"kyooar/internal/subscription/services"
	organizationServices "kyooar/internal/organization/services"
)

type SubscriptionMiddleware struct {
	subscriptionService services.SubscriptionService
	usageService        services.UsageService
	organizationService   organizationServices.OrganizationService
}

func NewSubscriptionMiddleware(
	subscriptionService services.SubscriptionService,
	usageService services.UsageService,
	organizationService organizationServices.OrganizationService,
) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{
		subscriptionService: subscriptionService,
		usageService:        usageService,
		organizationService:   organizationService,
	}
}

// RequireActiveSubscription ensures user has an active subscription
func (m *SubscriptionMiddleware) RequireActiveSubscription() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accountID, ok := c.Get("account_id").(uuid.UUID)
			if !ok {
				return response.Error(c, errors.BadRequest("Invalid or missing account information"))
			}

			subscription, err := m.subscriptionService.GetUserSubscription(c.Request().Context(), accountID)
			if err != nil || subscription == nil {
				return response.Error(c, errors.ErrNoSubscriptionFound)
			}

			if !subscription.IsActive() {
				return response.Error(c, errors.ErrSubscriptionNotActive)
			}

			// Store subscription in context for later use
			c.Set("subscription", subscription)
			return next(c)
		}
	}
}

// RequireFeature checks if the subscription has a specific feature flag
func (m *SubscriptionMiddleware) RequireFeature(feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			subscription, ok := c.Get("subscription").(*models.Subscription)
			if !ok {
				// Try to fetch subscription if not in context
				accountID, ok := c.Get("account_id").(uuid.UUID)
				if !ok {
					return response.Error(c, errors.ErrUnauthorized)
				}

				var err error
				subscription, err = m.subscriptionService.GetUserSubscription(c.Request().Context(), accountID)
				if err != nil || subscription == nil {
					return response.Error(c, errors.ErrNoSubscriptionFound)
				}
				c.Set("subscription", subscription)
			}

			// Check feature availability using the new column-based approach
			hasFeature := subscription.Plan.GetFlag(feature)

			if !hasFeature {
				return response.Error(c, errors.Forbidden("This feature is not available in your current plan"))
			}

			return next(c)
		}
	}
}

// CheckResourceLimit checks if user can add a resource type
func (m *SubscriptionMiddleware) CheckResourceLimit(resourceType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			subscription, ok := c.Get("subscription").(*models.Subscription)
			if !ok {
				accountID, ok := c.Get("account_id").(uuid.UUID)
				if !ok {
					return response.Error(c, errors.ErrUnauthorized)
				}

				var err error
				subscription, err = m.subscriptionService.GetUserSubscription(c.Request().Context(), accountID)
				if err != nil || subscription == nil {
					return response.Error(c, errors.ErrNoSubscriptionFound)
				}
				c.Set("subscription", subscription)
			}

			// For organizations, check actual count instead of usage tracking
			if resourceType == models.ResourceTypeOrganization {
				accountID, ok := c.Get("account_id").(uuid.UUID)
				if !ok {
					return response.Error(c, errors.ErrUnauthorized)
				}

				currentCount, err := m.organizationService.CountByAccountID(c.Request().Context(), accountID)
				if err != nil {
					return response.Error(c, errors.BadRequest("Failed to check organization count"))
				}

				if int(currentCount) >= subscription.Plan.MaxOrganizations {
					return response.Error(c, errors.Forbidden(fmt.Sprintf("Organization limit reached (%d/%d)", int(currentCount), subscription.Plan.MaxOrganizations)))
				}
			} else {
				// Check usage limits for other resources
				canAdd, reason, err := m.usageService.CanAddResource(c.Request().Context(), subscription.ID, resourceType)
				if err != nil {
					return response.Error(c, errors.BadRequest("Failed to check resource limits"))
				}

				if !canAdd {
					return response.Error(c, errors.Forbidden(reason))
				}
			}

			// Store resource type for tracking after successful creation
			c.Set("track_resource_type", resourceType)
			return next(c)
		}
	}
}

// GetSubscriptionFromContext retrieves the subscription from echo context
func GetSubscriptionFromContext(c echo.Context) (*models.Subscription, error) {
	subscription, ok := c.Get("subscription").(*models.Subscription)
	if !ok {
		return nil, errors.ErrNoSubscriptionFound
	}
	return subscription, nil
}

// TrackUsageAfterSuccess tracks usage after successful resource creation
func (m *SubscriptionMiddleware) TrackUsageAfterSuccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Execute the handler first
			err := next(c)

			// Only track if handler was successful
			if err == nil && c.Response().Status < 400 {
				resourceType, ok := c.Get("track_resource_type").(string)
				if ok {
					subscription, ok := c.Get("subscription").(*models.Subscription)
					if ok {
						// Track usage asynchronously to not block response
						go func() {
							_ = m.usageService.TrackUsage(c.Request().Context(), subscription.ID, resourceType, 1)

							// Record usage event
							event := &models.UsageEvent{
								SubscriptionID: subscription.ID,
								EventType:      models.EventTypeCreate,
								ResourceType:   resourceType,
							}
							_ = m.usageService.RecordUsageEvent(c.Request().Context(), event)
						}()
					}
				}
			}

			return err
		}
	}
}
