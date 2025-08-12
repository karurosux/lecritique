package analyticsconstants

type ContextKey string

const (
	OrganizationIDKey ContextKey = "organization_id"
	ProductIDKey      ContextKey = "product_id"
	MetricsKey        ContextKey = "metrics"
	AnalyticsKey      ContextKey = "analytics"
)