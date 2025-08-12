package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/response"
)

func JWTAuth(authService authinterface.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return response.Error(c, errors.ErrUnauthorized)
			}

			token := parts[1]

			claims, err := authService.ValidateToken(token)
			if err != nil {
				return response.Error(c, err)
			}

			c.Set("account_id", claims.AccountID)
			c.Set("member_id", claims.MemberID)
			c.Set("user_id", claims.AccountID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims)
			
			if claims.SubscriptionFeatures != nil {
				c.Set("subscription_features", claims.SubscriptionFeatures)
			}

			return next(c)
		}
	}
}

func GetAccountID(c echo.Context) (uuid.UUID, error) {
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return accountID, nil
}

func GetMemberID(c echo.Context) (uuid.UUID, error) {
	memberID, ok := c.Get("member_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return memberID, nil
}

func GetUserID(c echo.Context) (uuid.UUID, error) {
	userID, ok := c.Get("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.ErrUnauthorized
	}
	return userID, nil
}

func GetEmail(c echo.Context) (string, error) {
	email, ok := c.Get("email").(string)
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return email, nil
}

func GetRole(c echo.Context) (models.MemberRole, error) {
	role, ok := c.Get("role").(models.MemberRole)
	if !ok {
		return "", errors.ErrUnauthorized
	}
	return role, nil
}

func GetClaims(c echo.Context) (*authinterface.Claims, error) {
	claims, ok := c.Get("claims").(*authinterface.Claims)
	if !ok {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}

func GetSubscriptionFeatures(c echo.Context) (*authinterface.SubscriptionFeatures, error) {
	features, ok := c.Get("subscription_features").(*authinterface.SubscriptionFeatures)
	if !ok {
		return nil, errors.New("SUBSCRIPTION_REQUIRED", "No subscription features available", 403)
	}
	return features, nil
}

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
