package authmodel

import (
	"time"
	"kyooar/internal/auth/models"
)

type AuthResponse struct {
	Token   string      `json:"token"`
	Account interface{} `json:"account"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type ProfileResponse struct {
	Account interface{} `json:"account"`
	Message string      `json:"message,omitempty"`
}

type DeactivationResponse struct {
	Message          string     `json:"message"`
	DeactivationDate *time.Time `json:"deactivation_date,omitempty"`
}

type MemberListResponse struct {
	Members []*models.TeamMember `json:"members"`
}

type InvitationResponse struct {
	Invitation interface{} `json:"invitation"`
	Message    string      `json:"message"`
}