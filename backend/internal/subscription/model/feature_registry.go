package subscriptionmodel

import (
	"fmt"
	"strings"
	subscriptionconstants "kyooar/internal/subscription/constants"
)

type FeatureType = subscriptionconstants.FeatureType

type FeatureDefinition struct {
	Key           string                 `json:"key"`
	Type          FeatureType            `json:"type"`
	DisplayName   string                 `json:"display_name"`
	Description   string                 `json:"description"`
	Unit          string                 `json:"unit,omitempty"`
	UnlimitedText string                 `json:"unlimited_text,omitempty"`
	Format        string                 `json:"format,omitempty"`
	Icon          string                 `json:"icon,omitempty"`
	Category      string                 `json:"category,omitempty"`
	SortOrder     int                    `json:"sort_order"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

var FeatureRegistry = map[string]FeatureDefinition{
	subscriptionconstants.LimitOrganizations: {
		Key:           subscriptionconstants.LimitOrganizations,
		Type:          subscriptionconstants.FeatureTypeLimit,
		DisplayName:   "Organizations",
		Description:   "Maximum number of organizations",
		Unit:          "organizations",
		UnlimitedText: "Unlimited organizations",
		Format:        "{value} organization(s)",
		Icon:          "store",
		Category:      "core",
		SortOrder:     1,
	},
	subscriptionconstants.LimitQRCodes: {
		Key:           subscriptionconstants.LimitQRCodes,
		Type:          subscriptionconstants.FeatureTypeLimit,
		DisplayName:   "QR Codes",
		Description:   "Total QR codes across all organizations",
		Unit:          "QR codes",
		UnlimitedText: "Unlimited QR codes",
		Format:        "{value} QR codes",
		Icon:          "qr-code",
		Category:      "core",
		SortOrder:     2,
	},
	subscriptionconstants.LimitFeedbacksPerMonth: {
		Key:           subscriptionconstants.LimitFeedbacksPerMonth,
		Type:          subscriptionconstants.FeatureTypeLimit,
		DisplayName:   "Monthly Feedbacks",
		Description:   "Maximum feedbacks per month",
		Unit:          "feedbacks/month",
		UnlimitedText: "Unlimited feedbacks",
		Format:        "{value} feedbacks/month",
		Icon:          "message-square",
		Category:      "core",
		SortOrder:     3,
	},
	subscriptionconstants.LimitTeamMembers: {
		Key:           subscriptionconstants.LimitTeamMembers,
		Type:          subscriptionconstants.FeatureTypeLimit,
		DisplayName:   "Team Members",
		Description:   "Maximum team members",
		Unit:          "members",
		UnlimitedText: "Unlimited team members",
		Format:        "{value} team member(s)",
		Icon:          "users",
		Category:      "collaboration",
		SortOrder:     4,
	},

	subscriptionconstants.FlagBasicAnalytics: {
		Key:         subscriptionconstants.FlagBasicAnalytics,
		Type:        subscriptionconstants.FeatureTypeFlag,
		DisplayName: "Basic Analytics",
		Description: "View feedback analytics and insights",
		Icon:        "bar-chart-2",
		Category:    "analytics",
		SortOrder:   20,
	},
	subscriptionconstants.FlagAdvancedAnalytics: {
		Key:         subscriptionconstants.FlagAdvancedAnalytics,
		Type:        subscriptionconstants.FeatureTypeFlag,
		DisplayName: "Advanced Analytics",
		Description: "Detailed insights and reporting",
		Icon:        "bar-chart",
		Category:    "analytics",
		SortOrder:   21,
	},
	subscriptionconstants.FlagFeedbackExplorer: {
		Key:         subscriptionconstants.FlagFeedbackExplorer,
		Type:        subscriptionconstants.FeatureTypeFlag,
		DisplayName: "Feedback Explorer",
		Description: "Browse and search all feedback",
		Icon:        "search",
		Category:    "analytics",
		SortOrder:   22,
	},
	subscriptionconstants.FlagCustomBranding: {
		Key:         subscriptionconstants.FlagCustomBranding,
		Type:        subscriptionconstants.FeatureTypeFlag,
		DisplayName: "Custom Branding",
		Description: "Customize with your brand",
		Icon:        "palette",
		Category:    "customization",
		SortOrder:   23,
	},
	subscriptionconstants.FlagPrioritySupport: {
		Key:         subscriptionconstants.FlagPrioritySupport,
		Type:        subscriptionconstants.FeatureTypeFlag,
		DisplayName: "Priority Support",
		Description: "24/7 priority customer support",
		Icon:        "headphones",
		Category:    "support",
		SortOrder:   24,
	},
}

func GetFeatureDefinition(key string) (FeatureDefinition, bool) {
	def, exists := FeatureRegistry[key]
	return def, exists
}

func GetFeaturesByCategory(category string) []FeatureDefinition {
	var features []FeatureDefinition
	for _, def := range FeatureRegistry {
		if def.Category == category {
			features = append(features, def)
		}
	}
	return features
}

func FormatFeatureValue(key string, value interface{}) string {
	def, exists := FeatureRegistry[key]
	if !exists {
		return ""
	}

	switch def.Type {
	case subscriptionconstants.FeatureTypeLimit:
		limitValue, ok := value.(int64)
		if !ok {
			if intVal, ok := value.(int); ok {
				limitValue = int64(intVal)
			} else {
				return ""
			}
		}
		
		if limitValue == -1 {
			return def.UnlimitedText
		}
		
		result := def.Format
		if result == "" {
			result = "{value} {unit}"
		}
		result = replaceValue(result, "{value}", limitValue)
		result = replaceValue(result, "{unit}", def.Unit)
		return result
		
	case subscriptionconstants.FeatureTypeFlag:
		if boolVal, ok := value.(bool); ok && boolVal {
			return def.DisplayName
		}
		return ""
		
	default:
		return ""
	}
}

func replaceValue(s string, old string, value interface{}) string {
	return strings.ReplaceAll(s, old, fmt.Sprintf("%v", value))
}