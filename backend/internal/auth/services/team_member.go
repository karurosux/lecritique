package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/services"
)

type TeamMemberService interface {
	ListMembers(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error)
	InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamMember, error)
	UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error
	RemoveMember(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) error
	GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error)
	AcceptInvitation(ctx context.Context, token string) error
}

type teamMemberService struct {
	teamMemberRepo  repositories.TeamMemberRepository
	userRepo        repositories.UserRepository
	accountRepo     repositories.AccountRepository
	tokenRepo       repositories.TokenRepository
	emailService    services.EmailService
}

func NewTeamMemberService(
	teamMemberRepo repositories.TeamMemberRepository,
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	tokenRepo repositories.TokenRepository,
	emailService services.EmailService,
) TeamMemberService {
	return &teamMemberService{
		teamMemberRepo: teamMemberRepo,
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		tokenRepo:      tokenRepo,
		emailService:   emailService,
	}
}

func (s *teamMemberService) ListMembers(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error) {
	return s.teamMemberRepo.FindByAccountID(ctx, accountID)
}

func (s *teamMemberService) InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamMember, error) {
	// Validate role
	if role == models.RoleOwner {
		return nil, errors.New("FORBIDDEN", "Cannot invite another owner", http.StatusForbidden)
	}

	// TODO: Check team member limit when usage service is ready

	// Check if user exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// If user doesn't exist, create a placeholder user
		user = &models.User{
			Email:    email,
			IsActive: false,
		}
		// Generate a temporary password
		if err := user.SetPassword(uuid.New().String()); err != nil {
			return nil, err
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		fmt.Printf("Created new user: %s (ID: %s)\n", user.Email, user.ID)
	} else {
		fmt.Printf("Found existing user: %s (ID: %s)\n", user.Email, user.ID)
	}

	// Check if already a member (including soft-deleted)
	existing, err := s.teamMemberRepo.FindByUserAndAccount(ctx, user.ID, accountID)
	if err == nil && existing != nil {
		fmt.Printf("User %s is already a team member with role: %s\n", user.Email, existing.Role)
		return nil, errors.New("CONFLICT", "User is already a team member", http.StatusConflict)
	}
	
	// For now, let's just check if this is the case by logging
	// TODO: Implement proper soft-delete handling in repository
	
	fmt.Printf("User %s is not a team member yet, proceeding with invitation\n", user.Email)

	// Create team member
	member := &models.TeamMember{
		AccountID:  accountID,
		UserID:     user.ID,
		Role:       role,
		InvitedBy:  inviterID,
		InvitedAt:  time.Now(),
		AcceptedAt: nil,
	}

	if err := s.teamMemberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// Generate invitation token
	tokenStr, err := models.GenerateToken()
	if err != nil {
		return nil, err
	}

	verificationToken := &models.VerificationToken{
		AccountID: accountID,
		Token:     tokenStr,
		Type:      models.TokenTypeTeamInvite,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		// Store member ID in NewEmail field as a workaround
		NewEmail:  member.ID.String(),
	}

	if err := s.tokenRepo.Create(ctx, verificationToken); err != nil {
		return nil, err
	}

	// Send invitation email
	if err := s.sendInvitationEmail(ctx, email, tokenStr, accountID); err != nil {
		// Log error but don't fail the invitation
		fmt.Printf("Failed to send invitation email: %v\n", err)
	}

	// Load user data
	member.User = *user

	// TODO: Update usage tracking when ready

	return member, nil
}

func (s *teamMemberService) UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error {
	// Get member
	member, err := s.teamMemberRepo.FindByID(ctx, memberID)
	if err != nil {
		return errors.New("NOT_FOUND", "Team member not found", http.StatusNotFound)
	}

	// Verify member belongs to account
	if member.AccountID != accountID {
		return errors.ErrForbidden
	}

	// Cannot change owner role
	if member.Role == models.RoleOwner {
		return errors.New("FORBIDDEN", "Cannot change owner role", http.StatusForbidden)
	}

	// Cannot set someone as owner
	if newRole == models.RoleOwner {
		return errors.New("FORBIDDEN", "Cannot assign owner role", http.StatusForbidden)
	}

	// Update role
	member.Role = newRole
	return s.teamMemberRepo.Update(ctx, member)
}

func (s *teamMemberService) RemoveMember(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) error {
	// Get member
	member, err := s.teamMemberRepo.FindByID(ctx, memberID)
	if err != nil {
		return errors.New("NOT_FOUND", "Team member not found", http.StatusNotFound)
	}

	// Verify member belongs to account
	if member.AccountID != accountID {
		return errors.ErrForbidden
	}

	// Cannot remove owner
	if member.Role == models.RoleOwner {
		return errors.New("FORBIDDEN", "Cannot remove owner", http.StatusForbidden)
	}

	// Delete member
	if err := s.teamMemberRepo.Delete(ctx, memberID); err != nil {
		return err
	}

	// TODO: Update usage tracking when ready

	return nil
}

func (s *teamMemberService) GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error) {
	member, err := s.teamMemberRepo.FindByID(ctx, memberID, "User")
	if err != nil {
		return nil, err
	}

	// Verify member belongs to account
	if member.AccountID != accountID {
		return nil, errors.ErrForbidden
	}

	return member, nil
}

func (s *teamMemberService) AcceptInvitation(ctx context.Context, token string) error {
	// Find token
	tokenData, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	}

	// Validate token
	if !tokenData.IsValid() {
		return errors.New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	}

	if tokenData.Type != models.TokenTypeTeamInvite {
		return errors.New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	}

	// Get member ID from NewEmail field (used as workaround)
	memberID, err := uuid.Parse(tokenData.NewEmail)
	if err != nil {
		return errors.New("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	}

	// Get and update member
	member, err := s.teamMemberRepo.FindByID(ctx, memberID)
	if err != nil {
		return err
	}

	// Mark as accepted
	now := time.Now()
	member.AcceptedAt = &now

	if err := s.teamMemberRepo.Update(ctx, member); err != nil {
		return err
	}

	// Mark token as used
	return s.tokenRepo.MarkAsUsed(ctx, tokenData.ID)
}

func (s *teamMemberService) sendInvitationEmail(ctx context.Context, email, token string, accountID uuid.UUID) error {
	// Get account details
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	return s.emailService.SendTeamInviteEmail(ctx, email, token, account.CompanyName)
}