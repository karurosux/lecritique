package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/shared/models"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	models.BaseModel
	Email                   string        `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash            string        `gorm:"not null" json:"-"`
	Name                    string        `gorm:"not null" json:"name"`
	FirstName               string        `json:"first_name"`
	LastName                string        `json:"last_name"`
	Phone                   string        `json:"phone"`
	IsActive                bool          `gorm:"default:true" json:"is_active"`
	EmailVerified           bool          `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt         *time.Time    `json:"email_verified_at"`
	DeactivationRequestedAt *time.Time    `json:"deactivation_requested_at"`
	SubscriptionID          *uuid.UUID    `json:"subscription_id"`
	Subscription            interface{}   `gorm:"-" json:"subscription,omitempty"`
	TeamMembers             []TeamMember  `json:"team_members,omitempty"`
}

func (a *Account) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	a.PasswordHash = string(hash)
	return nil
}

func (a *Account) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(password))
	return err == nil
}

func (a *Account) DisplayName() string {
	if a.FirstName != "" || a.LastName != "" {
		return strings.TrimSpace(a.FirstName + " " + a.LastName)
	}
	if a.Name != "" {
		return a.Name
	}
	return a.Email
}

func (a *Account) IsPendingDeactivation() bool {
	return a.DeactivationRequestedAt != nil
}

func (a *Account) GetDeactivationDate() *time.Time {
	if a.DeactivationRequestedAt == nil {
		return nil
	}
	deactivationDate := a.DeactivationRequestedAt.Add(15 * 24 * time.Hour)
	return &deactivationDate
}

func (a *Account) ShouldBeDeactivated() bool {
	if a.DeactivationRequestedAt == nil {
		return false
	}
	return time.Now().After(a.DeactivationRequestedAt.Add(15 * 24 * time.Hour))
}

type TeamMember struct {
	models.BaseModel
	AccountID      uuid.UUID  `gorm:"not null" json:"account_id"`
	Account        Account    `json:"account,omitempty"`
	MemberID       uuid.UUID  `gorm:"not null" json:"member_id"`
	MemberAccount  Account    `gorm:"foreignKey:MemberID" json:"member,omitempty"`
	Role           MemberRole `gorm:"not null" json:"role"`
	InvitedBy      uuid.UUID  `json:"invited_by"`
	InvitedAt      time.Time  `json:"invited_at"`
	AcceptedAt     *time.Time `json:"accepted_at"`
}

type MemberRole string

const (
	RoleOwner   MemberRole = "OWNER"
	RoleAdmin   MemberRole = "ADMIN"
	RoleManager MemberRole = "MANAGER"
	RoleViewer  MemberRole = "VIEWER"
)

