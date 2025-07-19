package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "lecritique/internal/auth/models"
	authServices "lecritique/internal/auth/services"
	sharedErrors "lecritique/internal/shared/errors"
	"lecritique/internal/shared/middleware"
	sharedRepos "lecritique/internal/shared/repositories"
	"lecritique/internal/shared/response"
	"lecritique/internal/shared/validator"
	"lecritique/internal/subscription/services"
	"github.com/samber/do"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
	usageService        services.UsageService
	teamMemberService   authServices.TeamMemberServiceV2
	validator           *validator.Validator
}

func NewSubscriptionHandler(i *do.Injector) (*SubscriptionHandler, error) {
	return &SubscriptionHandler{
		subscriptionService: do.MustInvoke[services.SubscriptionService](i),
		usageService:        do.MustInvoke[services.UsageService](i),
		teamMemberService:   do.MustInvoke[authServices.TeamMemberServiceV2](i),
		validator:           validator.New(),
	}, nil
}

type CreateSubscriptionRequest struct {
	PlanID string `json:"plan_id" validate:"required,uuid"`
}


// GetAvailablePlans godoc
// @Summary Get available subscription plans
// @Description Retrieve all available subscription plans with their features and pricing
// @Tags subscription
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]models.SubscriptionPlan}
// @Failure 500 {object} response.Response
// @Router /api/v1/plans [get]
func (h *SubscriptionHandler) GetAvailablePlans(c echo.Context) error {
	ctx := c.Request().Context()

	plans, err := h.subscriptionService.GetAvailablePlans(ctx)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, plans)
}

// GetUserSubscription godoc
// @Summary Get user's current subscription
// @Description Retrieve the current subscription details for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.Subscription}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [get]
func (h *SubscriptionHandler) GetUserSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	log.Printf("GetUserSubscription called for account: %s", accountID)

	// Check if this user is a team member of another account using service
	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	var externalMemberships []authModels.TeamMember
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		externalMemberships = []authModels.TeamMember{*teamMember}
	}
	
	log.Printf("Found %d external team memberships for account %s", len(externalMemberships), accountID)
	
	// If user is a team member of another organization, get that organization's subscription
	if len(externalMemberships) > 0 {
		// Use the organization's account ID
		orgAccountID := externalMemberships[0].AccountID
		log.Printf("Team member detected, fetching subscription for organization: %s", orgAccountID)
		subscription, err := h.subscriptionService.GetUserSubscription(ctx, orgAccountID)
		if err != nil {
			log.Printf("Failed to get organization subscription: %v", err)
			// If no subscription found, return null instead of error
			if errors.Is(err, sharedRepos.ErrRecordNotFound) {
				return response.Success(c, nil)
			}
			return response.Error(c, err)
		}
		log.Printf("Returning organization subscription: %+v", subscription)
		return response.Success(c, subscription)
	}

	// Otherwise, get the user's own subscription
	log.Printf("Not a team member, fetching own subscription for account: %s", accountID)
	subscription, err := h.subscriptionService.GetUserSubscription(ctx, accountID)
	if err != nil {
		log.Printf("Failed to get user subscription: %v", err)
		// If no subscription found, return null instead of error
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			return response.Success(c, nil)
		}
		return response.Error(c, err)
	}

	log.Printf("Returning user subscription: %+v", subscription)
	return response.Success(c, subscription)
}

// GetUserUsage godoc
// @Summary Get user's current subscription usage
// @Description Retrieve the current subscription usage details for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=models.SubscriptionUsage}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription/usage [get]
func (h *SubscriptionHandler) GetUserUsage(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	log.Printf("GetUserUsage called for account: %s", accountID)

	// Check if this user is a team member of another account using service
	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	var targetAccountID uuid.UUID
	
	// If user is a team member of another organization, use that organization's account ID
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		targetAccountID = teamMember.AccountID
		log.Printf("Team member detected, fetching usage for organization: %s", targetAccountID)
	} else {
		targetAccountID = accountID
		log.Printf("Not a team member, fetching own usage for account: %s", accountID)
	}

	// Get the subscription first
	subscription, err := h.subscriptionService.GetUserSubscription(ctx, targetAccountID)
	if err != nil {
		log.Printf("Failed to get subscription: %v", err)
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			return response.Success(c, nil)
		}
		return response.Error(c, err)
	}

	// Get usage data
	usage, err := h.usageService.GetCurrentUsage(ctx, subscription.ID)
	if err != nil {
		log.Printf("Failed to get usage: %v", err)
		return response.Error(c, err)
	}

	log.Printf("Returning usage data: %+v", usage)
	return response.Success(c, usage)
}

// CanUserCreateRestaurant godoc
// @Summary Check if user can create more restaurants
// @Description Check if the authenticated user can create more restaurants based on their subscription plan
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=services.PermissionResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/can-create-restaurant [get]
func (h *SubscriptionHandler) CanUserCreateRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	// Check if this user is a team member of another account using service
	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	
	// If user is a team member, use the organization's account ID
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		orgAccountID := teamMember.AccountID
		permission, err := h.subscriptionService.CanUserCreateRestaurant(ctx, orgAccountID)
		if err != nil {
			return response.Error(c, err)
		}
		return response.Success(c, permission)
	}

	// Otherwise, use the user's own account ID
	permission, err := h.subscriptionService.CanUserCreateRestaurant(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, permission)
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} response.Response{data=models.Subscription}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [post]
func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	var req CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, sharedErrors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, sharedErrors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return response.Error(c, sharedErrors.ErrBadRequest)
	}

	createReq := &services.CreateSubscriptionRequest{
		AccountID: accountID,
		PlanID:    planID,
	}

	subscription, err := h.subscriptionService.CreateSubscription(ctx, createReq)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, subscription)
}

// CancelSubscription godoc
// @Summary Cancel user's subscription
// @Description Cancel the current subscription for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [delete]
func (h *SubscriptionHandler) CancelSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	err = h.subscriptionService.CancelSubscription(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{"message": "Subscription cancelled successfully"})
}