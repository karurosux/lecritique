package models

import (
	"fmt"
	"strings"
)

// FeatureType represents the type of feature
type FeatureType string

const (
	FeatureTypeLimit FeatureType = "limit"
	FeatureTypeFlag  FeatureType = "flag"
	FeatureTypeCustom FeatureType = "custom"
)

// FeatureDefinition defines metadata for a feature
type FeatureDefinition struct {
	Key           string                 `json:"key"`
	Type          FeatureType            `json:"type"`
	DisplayName   string                 `json:"display_name"`
	Description   string                 `json:"description"`
	Unit          string                 `json:"unit,omitempty"`          // e.g., "GB", "per month"
	UnlimitedText string                 `json:"unlimited_text,omitempty"` // e.g., "Unlimited storage"
	Format        string                 `json:"format,omitempty"`         // e.g., "{value} {unit}"
	Icon          string                 `json:"icon,omitempty"`           // Icon identifier
	Category      string                 `json:"category,omitempty"`       // For grouping features
	SortOrder     int                    `json:"sort_order"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// FeatureRegistry holds all feature definitions
var FeatureRegistry = map[string]FeatureDefinition{
	// Limits
	LimitRestaurants: {
		Key:           LimitRestaurants,
		Type:          FeatureTypeLimit,
		DisplayName:   "Restaurants",
		Description:   "Maximum number of restaurants",
		Unit:          "restaurants",
		UnlimitedText: "Unlimited restaurants",
		Format:        "{value} restaurant(s)",
		Icon:          "store",
		Category:      "core",
		SortOrder:     1,
	},
	LimitQRCodes: {
		Key:           LimitQRCodes,
		Type:          FeatureTypeLimit,
		DisplayName:   "QR Codes",
		Description:   "Total QR codes across all restaurants",
		Unit:          "QR codes",
		UnlimitedText: "Unlimited QR codes",
		Format:        "{value} QR codes",
		Icon:          "qr-code",
		Category:      "core",
		SortOrder:     2,
	},
	LimitFeedbacksPerMonth: {
		Key:           LimitFeedbacksPerMonth,
		Type:          FeatureTypeLimit,
		DisplayName:   "Monthly Feedbacks",
		Description:   "Maximum feedbacks per month",
		Unit:          "feedbacks/month",
		UnlimitedText: "Unlimited feedbacks",
		Format:        "{value} feedbacks/month",
		Icon:          "message-square",
		Category:      "core",
		SortOrder:     3,
	},
	LimitTeamMembers: {
		Key:           LimitTeamMembers,
		Type:          FeatureTypeLimit,
		DisplayName:   "Team Members",
		Description:   "Maximum team members",
		Unit:          "members",
		UnlimitedText: "Unlimited team members",
		Format:        "{value} team member(s)",
		Icon:          "users",
		Category:      "collaboration",
		SortOrder:     4,
	},

	// Flags
	FlagBasicAnalytics: {
		Key:         FlagBasicAnalytics,
		Type:        FeatureTypeFlag,
		DisplayName: "Basic Analytics",
		Description: "View feedback analytics and insights",
		Icon:        "bar-chart-2",
		Category:    "analytics",
		SortOrder:   20,
	},
	FlagAdvancedAnalytics: {
		Key:         FlagAdvancedAnalytics,
		Type:        FeatureTypeFlag,
		DisplayName: "Advanced Analytics",
		Description: "Detailed insights and reporting",
		Icon:        "bar-chart",
		Category:    "analytics",
		SortOrder:   21,
	},
	FlagFeedbackExplorer: {
		Key:         FlagFeedbackExplorer,
		Type:        FeatureTypeFlag,
		DisplayName: "Feedback Explorer",
		Description: "Browse and search all feedback",
		Icon:        "search",
		Category:    "analytics",
		SortOrder:   22,
	},
	FlagCustomBranding: {
		Key:         FlagCustomBranding,
		Type:        FeatureTypeFlag,
		DisplayName: "Custom Branding",
		Description: "Customize with your brand",
		Icon:        "palette",
		Category:    "customization",
		SortOrder:   23,
	},
	FlagPrioritySupport: {
		Key:         FlagPrioritySupport,
		Type:        FeatureTypeFlag,
		DisplayName: "Priority Support",
		Description: "24/7 priority customer support",
		Icon:        "headphones",
		Category:    "support",
		SortOrder:   24,
	},
}

// GetFeatureDefinition returns the definition for a feature key
func GetFeatureDefinition(key string) (FeatureDefinition, bool) {
	def, exists := FeatureRegistry[key]
	return def, exists
}

// GetFeaturesByCategory returns all features in a category
func GetFeaturesByCategory(category string) []FeatureDefinition {
	var features []FeatureDefinition
	for _, def := range FeatureRegistry {
		if def.Category == category {
			features = append(features, def)
		}
	}
	return features
}

// FormatFeatureValue formats a feature value for display
func FormatFeatureValue(key string, value interface{}) string {
	def, exists := FeatureRegistry[key]
	if !exists {
		return ""
	}

	switch def.Type {
	case FeatureTypeLimit:
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
		
		// Simple format replacement
		result := def.Format
		if result == "" {
			result = "{value} {unit}"
		}
		result = replaceValue(result, "{value}", limitValue)
		result = replaceValue(result, "{unit}", def.Unit)
		return result
		
	case FeatureTypeFlag:
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