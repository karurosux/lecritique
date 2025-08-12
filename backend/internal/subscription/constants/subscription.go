package subscriptionconstants

type SubscriptionStatus string

const (
	SubscriptionActive   SubscriptionStatus = "active"
	SubscriptionPending  SubscriptionStatus = "pending"
	SubscriptionCanceled SubscriptionStatus = "canceled"
	SubscriptionExpired  SubscriptionStatus = "expired"
)

const (
	LimitOrganizations       = "max_organizations"
	LimitQRCodes           = "max_qr_codes"
	LimitFeedbacksPerMonth = "max_feedbacks_per_month"
	LimitTeamMembers       = "max_team_members"
)

const (
	FlagBasicAnalytics    = "basic_analytics"
	FlagAdvancedAnalytics = "advanced_analytics"
	FlagFeedbackExplorer  = "feedback_explorer"
	FlagCustomBranding    = "custom_branding"
	FlagPrioritySupport   = "priority_support"
)

type FeatureType string

const (
	FeatureTypeLimit  FeatureType = "limit"
	FeatureTypeFlag   FeatureType = "flag"
	FeatureTypeCustom FeatureType = "custom"
)

type UsageType string

const (
	UsageTypeFeedback     UsageType = "feedback"
	UsageTypeQRScan       UsageType = "qr_scan"
	UsageTypeOrganization UsageType = "organization"
	UsageTypeTeamMember   UsageType = "team_member"
	UsageTypeQRCode       UsageType = "qr_code"
)

type PaymentProvider string

const (
	PaymentProviderStripe PaymentProvider = "stripe"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCanceled  PaymentStatus = "canceled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)