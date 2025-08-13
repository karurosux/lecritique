package subscriptioncontroller

import (
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	authModels "kyooar/internal/auth/models"
	authinterface "kyooar/internal/auth/interface"
	sharedErrors "kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	sharedRepos "kyooar/internal/shared/repositories"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type SubscriptionController struct {
	subscriptionService subscriptioninterface.SubscriptionService
	usageService        subscriptioninterface.UsageService
	teamMemberService   authinterface.TeamMemberService
	validator           *validator.Validator
}

func NewSubscriptionController(
	subscriptionService subscriptioninterface.SubscriptionService,
	usageService subscriptioninterface.UsageService,
	teamMemberService authinterface.TeamMemberService,
) *SubscriptionController {
	return &SubscriptionController{
		subscriptionService: subscriptionService,
		usageService:        usageService,
		teamMemberService:   teamMemberService,
		validator:           validator.New(),
	}
}


// @Summary Get available subscription plans
// @Description Retrieve all available subscription plans with their features and pricing
// @Tags subscription
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]subscriptionmodel.SubscriptionPlan}
// @Failure 500 {object} response.Response
// @Router /api/v1/plans [get]
func (h *SubscriptionController) GetAvailablePlans(c echo.Context) error {
	ctx := c.Request().Context()

	plans, err := h.subscriptionService.GetAvailablePlans(ctx)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, plans)
}

// @Summary Get user's current subscription
// @Description Retrieve the current subscription details for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=subscriptionmodel.Subscription}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [get]
func (h *SubscriptionController) GetUserSubscription(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	log.Printf("GetUserSubscription called for account: %s", accountID)

	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	var externalMemberships []authModels.TeamMember
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		externalMemberships = []authModels.TeamMember{*teamMember}
	}
	
	log.Printf("Found %d external team memberships for account %s", len(externalMemberships), accountID)
	
	if len(externalMemberships) > 0 {
		orgAccountID := externalMemberships[0].AccountID
		log.Printf("Team member detected, fetching subscription for organization: %s", orgAccountID)
		subscription, err := h.subscriptionService.GetUserSubscription(ctx, orgAccountID)
		if err != nil {
			log.Printf("Failed to get organization subscription: %v", err)
			if errors.Is(err, sharedRepos.ErrRecordNotFound) {
				return response.Success(c, nil)
			}
			return response.Error(c, err)
		}
		log.Printf("Returning organization subscription: %+v", subscription)
		return response.Success(c, subscription)
	}

	log.Printf("Not a team member, fetching own subscription for account: %s", accountID)
	subscription, err := h.subscriptionService.GetUserSubscription(ctx, accountID)
	if err != nil {
		log.Printf("Failed to get user subscription: %v", err)
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			return response.Success(c, nil)
		}
		return response.Error(c, err)
	}

	log.Printf("Returning user subscription: %+v", subscription)
	return response.Success(c, subscription)
}

// @Summary Get user's current subscription usage
// @Description Retrieve the current subscription usage details for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=subscriptionmodel.SubscriptionUsage}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription/usage [get]
func (h *SubscriptionController) GetUserUsage(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	log.Printf("GetUserUsage called for account: %s", accountID)

	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	var targetAccountID uuid.UUID
	
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		targetAccountID = teamMember.AccountID
		log.Printf("Team member detected, fetching usage for organization: %s", targetAccountID)
	} else {
		targetAccountID = accountID
		log.Printf("Not a team member, fetching own usage for account: %s", accountID)
	}

	subscription, err := h.subscriptionService.GetUserSubscription(ctx, targetAccountID)
	if err != nil {
		log.Printf("Failed to get subscription: %v", err)
		if errors.Is(err, sharedRepos.ErrRecordNotFound) {
			return response.Success(c, nil)
		}
		return response.Error(c, err)
	}

	usage, err := h.usageService.GetCurrentUsage(ctx, subscription.ID)
	if err != nil {
		log.Printf("Failed to get usage: %v", err)
		return response.Error(c, err)
	}

	log.Printf("Returning usage data: %+v", usage)
	return response.Success(c, usage)
}

// @Summary Check if user can create more organizations
// @Description Check if the authenticated user can create more organizations based on their subscription plan
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=subscriptioninterface.PermissionResponse}
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/can-create-organization [get]
func (h *SubscriptionController) CanUserCreateOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	teamMember, err := h.teamMemberService.GetMemberByMemberID(ctx, accountID)
	
	if err == nil && teamMember != nil && teamMember.AccountID != accountID {
		orgAccountID := teamMember.AccountID
		permission, err := h.subscriptionService.CanUserCreateOrganization(ctx, orgAccountID)
		if err != nil {
			return response.Error(c, err)
		}
		return response.Success(c, permission)
	}

	permission, err := h.subscriptionService.CanUserCreateOrganization(ctx, accountID)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, permission)
}

// @Summary Create a new subscription
// @Description Create a new subscription for the authenticated user
// @Tags subscription
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} response.Response{data=subscriptionmodel.Subscription}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/user/subscription [post]
func (h *SubscriptionController) CreateSubscription(c echo.Context) error {
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

	createReq := &subscriptioninterface.CreateSubscriptionRequest{
		AccountID: accountID,
		PlanID:    planID,
	}

	subscription, err := h.subscriptionService.CreateSubscription(ctx, createReq)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, subscription)
}

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
func (h *SubscriptionController) CancelSubscription(c echo.Context) error {
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