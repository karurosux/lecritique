package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/shared/models"
)

type TokenType string

const (
	TokenTypeEmailVerification TokenType = "EMAIL_VERIFICATION"
	TokenTypePasswordReset     TokenType = "PASSWORD_RESET"
	TokenTypeTeamInvite        TokenType = "TEAM_INVITE"
	TokenTypeEmailChange       TokenType = "EMAIL_CHANGE"
)

type VerificationToken struct {
	models.BaseModel
	AccountID uuid.UUID  `gorm:"not null" json:"account_id"`
	Account   Account    `json:"account,omitempty"`
	Token     string     `gorm:"uniqueIndex;not null" json:"token"`
	Type      TokenType  `gorm:"not null" json:"type"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	NewEmail  string     `json:"new_email,omitempty"`
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (vt *VerificationToken) IsValid() bool {
	return vt.UsedAt == nil && time.Now().Before(vt.ExpiresAt)
}

func (vt *VerificationToken) MarkAsUsed() {
	now := time.Now()
	vt.UsedAt = &now
}