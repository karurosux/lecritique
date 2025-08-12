package authmodel

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	Name            string `json:"name" validate:"required_without=InvitationToken"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	InvitationToken string `json:"invitation_token,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SendEmailVerificationRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}

type ConfirmEmailChangeRequest struct {
	Token string `json:"token" validate:"required"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=1"`
	Phone string `json:"phone,omitempty" validate:"omitempty"`
}

type InviteMemberRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=OWNER ADMIN MANAGER VIEWER"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=OWNER ADMIN MANAGER VIEWER"`
}

type AcceptInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}