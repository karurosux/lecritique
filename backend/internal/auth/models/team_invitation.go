package models

import (
	"time"

	"github.com/google/uuid"
	"kyooar/internal/shared/models"
)

type TeamInvitation struct {
	models.BaseModel
	AccountID       uuid.UUID  `gorm:"not null" json:"account_id"`
	Account         Account    `json:"account,omitempty"`
	Email           string     `gorm:"not null;index" json:"email"`
	Role            MemberRole `gorm:"not null" json:"role"`
	InvitedBy       uuid.UUID  `gorm:"not null" json:"invited_by"`
	InvitedByUser   Account    `gorm:"foreignKey:InvitedBy" json:"invited_by_user,omitempty"`
	Token           string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt       time.Time  `gorm:"not null" json:"expires_at"`
	AcceptedAt      *time.Time `json:"accepted_at"`
	EmailAcceptedAt *time.Time `json:"email_accepted_at"`
}

func (ti *TeamInvitation) IsValid() bool {
	return ti.AcceptedAt == nil && time.Now().Before(ti.ExpiresAt)
}

func (ti *TeamInvitation) IsExpired() bool {
	return time.Now().After(ti.ExpiresAt)
}

func (ti *TeamInvitation) IsAccepted() bool {
	return ti.AcceptedAt != nil
}