package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	models.BaseModel
	Email                   string        `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash            string        `gorm:"not null" json:"-"`
	CompanyName             string        `gorm:"not null" json:"company_name"`
	Phone                   string        `json:"phone"`
	IsActive                bool          `gorm:"default:true" json:"is_active"`
	EmailVerified           bool          `gorm:"default:false" json:"email_verified"`
	EmailVerifiedAt         *time.Time    `json:"email_verified_at"`
	DeactivationRequestedAt *time.Time    `json:"deactivation_requested_at"`
	SubscriptionID          *uuid.UUID    `json:"subscription_id"`
	// Subscription     *Subscription `json:"subscription,omitempty"` // TODO: Add when subscription domain is ready
	// Restaurants      []Restaurant  `json:"restaurants,omitempty"`  // TODO: Add when restaurant domain is ready
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

// IsPendingDeactivation checks if the account has a pending deactivation request
func (a *Account) IsPendingDeactivation() bool {
	return a.DeactivationRequestedAt != nil
}

// GetDeactivationDate returns the date when the account will be deactivated (15 days after request)
func (a *Account) GetDeactivationDate() *time.Time {
	if a.DeactivationRequestedAt == nil {
		return nil
	}
	deactivationDate := a.DeactivationRequestedAt.Add(15 * 24 * time.Hour)
	return &deactivationDate
}

// ShouldBeDeactivated checks if the account should be deactivated (15 days have passed)
func (a *Account) ShouldBeDeactivated() bool {
	if a.DeactivationRequestedAt == nil {
		return false
	}
	return time.Now().After(a.DeactivationRequestedAt.Add(15 * 24 * time.Hour))
}

type TeamMember struct {
	models.BaseModel
	AccountID  uuid.UUID  `gorm:"not null" json:"account_id"`
	Account    Account    `json:"account,omitempty"`
	UserID     uuid.UUID  `gorm:"not null" json:"user_id"`
	User       User       `json:"user,omitempty"`
	Role       MemberRole `gorm:"not null" json:"role"`
	InvitedBy  uuid.UUID  `json:"invited_by"`
	InvitedAt  time.Time  `json:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at"`
}

type MemberRole string

const (
	RoleOwner   MemberRole = "OWNER"
	RoleAdmin   MemberRole = "ADMIN"
	RoleManager MemberRole = "MANAGER"
	RoleViewer  MemberRole = "VIEWER"
)

type User struct {
	models.BaseModel
	Email        string       `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string       `gorm:"not null" json:"-"`
	FirstName    string       `json:"first_name"`
	LastName     string       `json:"last_name"`
	IsActive     bool         `gorm:"default:true" json:"is_active"`
	TeamMembers  []TeamMember `json:"team_members,omitempty"`
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) FullName() string {
	if u.FirstName == "" && u.LastName == "" {
		return u.Email
	}
	return u.FirstName + " " + u.LastName
}
