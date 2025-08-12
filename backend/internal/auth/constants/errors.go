package authconstants

const (
	ErrAccountAlreadyExists    = "account already exists"
	ErrAccountNotFound        = "account not found"
	ErrInvalidCredentials     = "invalid credentials"
	ErrInvalidToken          = "invalid token"
	ErrTokenExpired          = "token expired"
	ErrTokenAlreadyUsed      = "token already used"
	ErrSessionExpired        = "session expired"
	ErrPasswordTooWeak       = "password too weak"
	ErrInvalidPassword       = "invalid password"
	ErrEmailNotVerified      = "email not verified"
	ErrEmailAlreadyVerified  = "email already verified"
	ErrDeactivationExists    = "deactivation already exists"
	ErrNoDeactivation        = "no deactivation request found"
	ErrTeamMemberNotFound    = "team member not found"
	ErrInvitationNotFound    = "invitation not found"
	ErrInvitationExpired     = "invitation expired"
	ErrInsufficientPrivileges = "insufficient privileges"
	ErrSameEmail             = "new email must be different"
)