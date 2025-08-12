package authinterface

import (
	"context"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/auth/models"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error)
	FindByEmail(ctx context.Context, email string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	UpdateEmailVerification(ctx context.Context, accountID uuid.UUID, verified bool) error
	FindAccountsPendingDeactivation(ctx context.Context) ([]models.Account, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token *models.VerificationToken) error
	FindByToken(ctx context.Context, token string) (*models.VerificationToken, error)
	FindByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) ([]*models.VerificationToken, error)
	Update(ctx context.Context, token *models.VerificationToken) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	MarkAsUsed(ctx context.Context, id uuid.UUID) error
	DeleteByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) error
}

type TeamInvitationRepository interface {
	Create(ctx context.Context, invitation *models.TeamInvitation) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.TeamInvitation, error)
	FindByToken(ctx context.Context, token string) (*models.TeamInvitation, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error)
	Update(ctx context.Context, invitation *models.TeamInvitation) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TeamMemberRepository interface {
	Create(ctx context.Context, member *models.TeamMember) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.TeamMember, error)
	FindByAccountAndMember(ctx context.Context, accountID, memberID uuid.UUID) (*models.TeamMember, error)
	FindByMemberID(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamMember, error)
	Update(ctx context.Context, member *models.TeamMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByMemberID(ctx context.Context, memberID uuid.UUID) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type TokenGenerator interface {
	GenerateToken() (string, error)
	GenerateSecureToken() (string, error)
}

type EmailSender interface {
	SendVerificationEmail(email, token string) error
	SendPasswordResetEmail(email, token string) error
	SendEmailChangeVerification(email, token string) error
	SendDeactivationRequest(email, deactivationDate string) error
	SendDeactivationCancelled(email string) error
	SendAccountDeactivated(email string) error
}

type AuthConfig interface {
	GetJWTSecret() string
	GetJWTExpiration() time.Duration
	IsEmailVerificationRequired() bool
	IsDevMode() bool
	IsSMTPConfigured() bool
}

type AuthService interface {
	Register(ctx context.Context, data RegisterData) (*models.Account, error)
	Login(ctx context.Context, email, password string) (string, *models.Account, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	SendEmailVerification(ctx context.Context, accountID uuid.UUID) error
	ResendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	SendPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	RequestEmailChange(ctx context.Context, accountID uuid.UUID, newEmail string) (string, error)
	ConfirmEmailChange(ctx context.Context, token string) (string, error)
	RequestDeactivation(ctx context.Context, accountID uuid.UUID) error
	CancelDeactivation(ctx context.Context, accountID uuid.UUID) error
	ProcessPendingDeactivations(ctx context.Context) error
	UpdateProfile(ctx context.Context, accountID uuid.UUID, updates map[string]interface{}) (*models.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
}

type TeamMemberService interface {
	GetMemberByMemberID(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error)
	ListMembers(ctx context.Context, accountID uuid.UUID) ([]*models.TeamMember, error)
	InviteMember(ctx context.Context, accountID uuid.UUID, email string, role models.MemberRole, invitedBy uuid.UUID) (*models.TeamInvitation, error)
	ResendInvitation(ctx context.Context, accountID, invitationID uuid.UUID) error
	UpdateRole(ctx context.Context, accountID, memberID uuid.UUID, role models.MemberRole) error
	UpdateRoleByID(ctx context.Context, teamMemberID uuid.UUID, role models.MemberRole) error
	RemoveMember(ctx context.Context, accountID, memberID uuid.UUID) error
	RemoveMemberByID(ctx context.Context, teamMemberID uuid.UUID) error
	AcceptInvitation(ctx context.Context, token string) (*models.TeamMember, error)
}

type RegisterData struct {
	Email     string
	Password  string
	Name      string
	FirstName string
	LastName  string
}

type Claims struct {
	AccountID            uuid.UUID                 `json:"account_id"`
	MemberID             uuid.UUID                 `json:"member_id"`
	Name                 string                    `json:"name"`
	Email                string                    `json:"email"`
	Role                 models.MemberRole         `json:"role"`
	SubscriptionFeatures *SubscriptionFeatures     `json:"subscription_features,omitempty"`
}

type SubscriptionFeatures struct {
	MaxOrganizations       int  `json:"max_organizations"`
	MaxQRCodes           int  `json:"max_qr_codes"`
	MaxFeedbacksPerMonth int  `json:"max_feedbacks_per_month"`
	MaxTeamMembers       int  `json:"max_team_members"`
	HasBasicAnalytics    bool `json:"has_basic_analytics"`
	HasAdvancedAnalytics bool `json:"has_advanced_analytics"`
	HasFeedbackExplorer  bool `json:"has_feedback_explorer"`
	HasCustomBranding    bool `json:"has_custom_branding"`
	HasPrioritySupport   bool `json:"has_priority_support"`
}