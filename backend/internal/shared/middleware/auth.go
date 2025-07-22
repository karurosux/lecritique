package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"lecritique/internal/auth/services"
	"lecritique/internal/auth/models"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/response"
)

func JWTAuth(authService services.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			// Check Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			token := parts[1]

			// Validate token
			claims, err := authService.ValidateToken(token)
			if err != nil {
				return response.Error(c, err)
			}

			// Set all claims data in context
			c.Set("account_id", claims.AccountID)
			c.Set("member_id", claims.MemberID)
			c.Set("user_id", claims.AccountID) // For compatibility
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims) // Store full claims for easy access
			
			// Set subscription features if available
			if claims.SubscriptionFeatures != nil {
				c.Set("subscription_features", claims.SubscriptionFeatures)
			}

			return next(c)
		}
	}
}

// Context helper functions

// GetAccountID retrieves the account ID from context
func GetAccountID(c echo.Context) (uuid.UUID, error) {
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return accountID, nil
}

// GetMemberID retrieves the member ID from context
func GetMemberID(c echo.Context) (uuid.UUID, error) {
	memberID, ok := c.Get("member_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return memberID, nil
}

// GetUserID retrieves the user ID from context (alias for account_id)
func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return userID, nil
}

// GetEmail retrieves the email from context
func GetEmail(c echo.Context) (string, error) {
	email, ok := c.Get("email").(string)
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return email, nil
}

// GetRole retrieves the role from context
func GetRole(c echo.Context) (models.MemberRole, error) {
	role, ok := c.Get("role").(models.MemberRole)
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return role, nil
}

// GetClaims retrieves the full claims from context
func GetClaims(c echo.Context) (*services.Claims, error) {
	claims, ok := c.Get("claims").(*services.Claims)
	if !ok {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

// GetSubscriptionFeatures retrieves the subscription features from context
func GetSubscriptionFeatures(c echo.Context) (*services.SubscriptionFeatures, error) {
	features, ok := c.Get("subscription_features").(*services.SubscriptionFeatures)
	if !ok {
		return nil, errors.New("SUBSCRIPTION_REQUIRED", "No subscription features available", 403)
	}
	return features, nil
}

// HasSubscriptionFeature checks if a specific feature is enabled
func HasSubscriptionFeature(c echo.Context, feature string) bool {
	features, err := GetSubscriptionFeatures(c)
	if err != nil {
		return false
	}
	
	switch feature {
	case "basic_analytics":
		return features.HasBasicAnalytics
	case "advanced_analytics":
		return features.HasAdvancedAnalytics
	case "feedback_explorer":
		return features.HasFeedbackExplorer
	case "custom_branding":
		return features.HasCustomBranding
	case "priority_support":
		return features.HasPrioritySupport
	default:
		return false
	}
}

// GetSubscriptionLimit retrieves a specific subscription limit
func GetSubscriptionLimit(c echo.Context, limitType string) (int, error) {
	features, err := GetSubscriptionFeatures(c)
	if err != nil {
		return 0, err
	}
	
	switch limitType {
	case "max_organizations":
		return features.MaxOrganizations, nil
	case "max_qr_codes":
		return features.MaxQRCodes, nil
	case "max_feedbacks_per_month":
		return features.MaxFeedbacksPerMonth, nil
	case "max_team_members":
		return features.MaxTeamMembers, nil
	default:
		return 0, errors.BadRequest("Invalid limit type")
	}
}

// RequireSubscriptionFeature middleware that requires a specific feature
func RequireSubscriptionFeature(feature string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !HasSubscriptionFeature(c, feature) {
				return response.Error(c, errors.Forbidden("This feature is not available in your current plan"))
			}
			return next(c)
		}
	}
}
