package middleware

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
	organizationinterface "kyooar/internal/organization/interface"
)

type SubscriptionMiddleware struct {
	subscriptionService subscriptioninterface.SubscriptionService
	usageService        subscriptioninterface.UsageService
	organizationService   organizationinterface.OrganizationService
}

func NewSubscriptionMiddleware(
	subscriptionService subscriptioninterface.SubscriptionService,
	usageService subscriptioninterface.UsageService,
	organizationService organizationinterface.OrganizationService,
) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{
		subscriptionService: subscriptionService,
		usageService:        usageService,
		organizationService:   organizationService,
	}
}

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

			c.Set("subscription", subscription)
			return next(c)
		}
	}
}

func (m *SubscriptionMiddleware) RequireFeature(feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			subscription, ok := c.Get("subscription").(*subscriptionmodel.Subscription)
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

			hasFeature := subscription.Plan.GetFlag(feature)

			if !hasFeature {
				return response.Error(c, errors.Forbidden("This feature is not available in your current plan"))
			}

			return next(c)
		}
	}
}

func (m *SubscriptionMiddleware) CheckResourceLimit(resourceType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			subscription, ok := c.Get("subscription").(*subscriptionmodel.Subscription)
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

			if resourceType == subscriptionmodel.ResourceTypeOrganization {
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
				canAdd, reason, err := m.usageService.CanAddResource(c.Request().Context(), subscription.ID, resourceType)
				if err != nil {
					return response.Error(c, errors.BadRequest("Failed to check resource limits"))
				}

				if !canAdd {
					return response.Error(c, errors.Forbidden(reason))
				}
			}

			c.Set("track_resource_type", resourceType)
			return next(c)
		}
	}
}

func GetSubscriptionFromContext(c echo.Context) (*subscriptionmodel.Subscription, error) {
	subscription, ok := c.Get("subscription").(*subscriptionmodel.Subscription)
	if !ok {
		return nil, errors.ErrNoSubscriptionFound
	}
	return subscription, nil
}

func (m *SubscriptionMiddleware) TrackUsageAfterSuccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if err == nil && c.Response().Status < 400 {
				resourceType, ok := c.Get("track_resource_type").(string)
				if ok {
					subscription, ok := c.Get("subscription").(*subscriptionmodel.Subscription)
					if ok {
						go func() {
							_ = m.usageService.TrackUsage(c.Request().Context(), subscription.ID, resourceType, 1)

							event := &subscriptionmodel.UsageEvent{
								SubscriptionID: subscription.ID,
								EventType:      subscriptionmodel.EventTypeCreate,
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
