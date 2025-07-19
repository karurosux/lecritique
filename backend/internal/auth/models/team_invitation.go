package models

import (
	"time"

	"github.com/google/uuid"
	"lecritique/internal/shared/models"
)

type TeamInvitation struct {
	models.BaseModel
	AccountID       uuid.UUID  `gorm:"not null" json:"account_id"`              // Organization sending invite
	Account         Account    `json:"account,omitempty"`
	Email           string     `gorm:"not null;index" json:"email"`             // Invited email
	Role            MemberRole `gorm:"not null" json:"role"`
	InvitedBy       uuid.UUID  `gorm:"not null" json:"invited_by"`
	InvitedByUser   Account    `gorm:"foreignKey:InvitedBy" json:"invited_by_user,omitempty"`
	Token           string     `gorm:"uniqueIndex;not null" json:"-"`           // Security: don't expose in JSON
	ExpiresAt       time.Time  `gorm:"not null" json:"expires_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`                             // When the invitation was fully accepted
	EmailAcceptedAt *time.Time `json:"email_accepted_at"`                       // When the email link was clicked
}

// IsValid checks if the invitation is still valid
func (ti *TeamInvitation) IsValid() bool {
	return ti.AcceptedAt == nil && time.Now().Before(ti.ExpiresAt)
}

// IsExpired checks if the invitation has expired
func (ti *TeamInvitation) IsExpired() bool {
	return time.Now().After(ti.ExpiresAt)
}

// IsAccepted checks if the invitation has been accepted
func (ti *TeamInvitation) IsAccepted() bool {
	return ti.AcceptedAt != nil
}