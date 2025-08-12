package organizationconstants

type ContextKey string

const (
	AccountIDKey      ContextKey = "account_id"
	OrganizationIDKey ContextKey = "organization_id"
	ResourceKey       ContextKey = "resource"
	UserIDKey         ContextKey = "user_id"
)