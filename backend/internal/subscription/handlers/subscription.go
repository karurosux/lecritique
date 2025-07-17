package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "github.com/lecritique/api/internal/auth/models"
	sharedErrors "github.com/lecritique/api/internal/shared/errors"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
	"github.com/lecritique/api/internal/subscription/services"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	subscriptionService services.SubscriptionService
	validator          *validator.Validator
	db                 *gorm.DB
}

func NewSubscriptionHandler(subscriptionService services.SubscriptionService, db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validator:          validator.New(),
		db:                 db,
	}
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
	accountID := c.Get("account_id").(uuid.UUID)

	log.Printf("GetUserSubscription called for account: %s", accountID)

	// Check if this user is a team member of another account
	var teamMemberships []authModels.TeamMember
	h.db.WithContext(ctx).
		Preload("Account").
		Where("member_id = ? AND accepted_at IS NOT NULL", accountID).
		Find(&teamMemberships)
	
	log.Printf("Found %d team memberships for account %s", len(teamMemberships), accountID)
	
	// Filter out memberships where the user is a member of their own account
	var externalMemberships []authModels.TeamMember
	for _, tm := range teamMemberships {
		if tm.AccountID != accountID {
			externalMemberships = append(externalMemberships, tm)
		}
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
	accountID := c.Get("account_id").(uuid.UUID)

	// Check if this user is a team member of another account
	var teamMemberships []authModels.TeamMember
	h.db.WithContext(ctx).
		Preload("Account").
		Where("member_id = ? AND accepted_at IS NOT NULL", accountID).
		Find(&teamMemberships)
	
	// If user is a team member, use the organization's account ID
	if len(teamMemberships) > 0 {
		orgAccountID := teamMemberships[0].AccountID
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
	accountID := c.Get("account_id").(uuid.UUID)

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
	accountID := c.Get("account_id").(uuid.UUID)

	err := h.subscriptionService.CancelSubscription(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{"message": "Subscription cancelled successfully"})
}