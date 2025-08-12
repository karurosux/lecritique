package authconstants

type ContextKey string

const (
	AccountIDKey ContextKey = "account_id"
	MemberIDKey  ContextKey = "member_id"
	UserKey      ContextKey = "user"
	ClaimsKey    ContextKey = "claims"
	SessionKey   ContextKey = "session"
)