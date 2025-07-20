package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"lecritique/internal/auth/models"
	"lecritique/internal/auth/repositories"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/services"
	"github.com/samber/do"
)

type TeamMemberServiceV2 interface {
	// Member management
	ListMembers(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error)
	GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error)
	GetMemberByMemberID(ctx context.Context, memberId uuid.UUID) (*models.TeamMember, error)
	UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error
	RemoveMember(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) error

	// Invitation management
	InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamInvitation, error)
	ResendInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error
	CancelInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error
	ListPendingInvitations(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error)

	// For registration flow
	CheckPendingInvitations(ctx context.Context, email string) ([]*models.TeamInvitation, error)
	AcceptInvitation(ctx context.Context, invitationToken string, memberAccountID uuid.UUID) error
	UpdateInvitation(ctx context.Context, invitation *models.TeamInvitation) error
}

type teamMemberServiceV2 struct {
	teamMemberRepo repositories.TeamMemberRepository
	invitationRepo repositories.TeamInvitationRepository
	accountRepo    repositories.AccountRepository
	emailService   services.EmailService
}

func NewTeamMemberServiceV2(i *do.Injector) (TeamMemberServiceV2, error) {
	return &teamMemberServiceV2{
		teamMemberRepo: do.MustInvoke[repositories.TeamMemberRepository](i),
		invitationRepo: do.MustInvoke[repositories.TeamInvitationRepository](i),
		accountRepo:    do.MustInvoke[repositories.AccountRepository](i),
		emailService:   do.MustInvoke[services.EmailService](i),
	}, nil
}

func (s *teamMemberServiceV2) GetMemberByMemberID(ctx context.Context, memberId uuid.UUID) (*models.TeamMember, error) {
	members, err := s.teamMemberRepo.FindByMemberIDNotOwner(ctx, memberId)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, nil
	}

	return &members[0], nil
}

func (s *teamMemberServiceV2) ListMembers(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error) {
	fmt.Println("acccount id => ", accountID)
	// Get active members
	members, err := s.teamMemberRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Get pending invitations and convert them to TeamMember format
	invitations, err := s.invitationRepo.FindPendingByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Convert invitations to TeamMember format with a pending status
	for _, inv := range invitations {
		member := models.TeamMember{
			AccountID: inv.AccountID,
			Role:      inv.Role,
			InvitedBy: inv.InvitedBy,
			InvitedAt: inv.CreatedAt,
			// Set MemberAccount with email for display
			MemberAccount: models.Account{
				Email: inv.Email,
			},
		}
		member.ID = inv.ID
		member.CreatedAt = inv.CreatedAt
		member.UpdatedAt = inv.UpdatedAt

		// Add to members list
		members = append(members, member)
	}

	return members, nil
}

func (s *teamMemberServiceV2) GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error) {
	member, err := s.teamMemberRepo.FindByID(ctx, memberID, "MemberAccount")
	if err != nil {
		return nil, err
	}

	// Verify member belongs to the account
	if member.AccountID != accountID {
		return nil, errors.New("NOT_FOUND", "Member not found", http.StatusNotFound)
	}

	return member, nil
}

func (s *teamMemberServiceV2) InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamInvitation, error) {
	// Validate role
	if role == models.RoleOwner {
		return nil, errors.New("FORBIDDEN", "Cannot invite another owner", http.StatusForbidden)
	}

	// Check if already a member
	existingAccount, _ := s.accountRepo.FindByEmail(ctx, email)
	if existingAccount != nil {
		existingMember, _ := s.teamMemberRepo.FindByMemberAndAccount(ctx, existingAccount.ID, accountID)
		if existingMember != nil {
			return nil, errors.New("CONFLICT", "User is already a team member", http.StatusConflict)
		}
	}

	// Check for existing pending invitation
	existingInvite, _ := s.invitationRepo.FindByAccountAndEmail(ctx, accountID, email)
	if existingInvite != nil && existingInvite.IsValid() {
		return nil, errors.New("CONFLICT", "Invitation already sent to this email", http.StatusConflict)
	}

	// Generate invitation token
	token, err := models.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Create invitation
	invitation := &models.TeamInvitation{
		AccountID: accountID,
		Email:     email,
		Role:      role,
		InvitedBy: inviterID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	// Get organization name for email
	org, _ := s.accountRepo.FindByID(ctx, accountID)
	companyName := "LeCritique"
	if org != nil {
		companyName = org.Name
	}

	// Send invitation email
	if err := s.emailService.SendTeamInviteEmail(ctx, email, token, companyName); err != nil {
		// Log error but don't fail the invitation
		fmt.Printf("Failed to send invitation email to %s: %v\n", email, err)
	}

	return invitation, nil
}

func (s *teamMemberServiceV2) UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error {
	member, err := s.GetMemberByID(ctx, accountID, memberID)
	if err != nil {
		return err
	}

	// Cannot change owner role
	if member.Role == models.RoleOwner || newRole == models.RoleOwner {
		return errors.New("FORBIDDEN", "Cannot change owner role", http.StatusForbidden)
	}

	member.Role = newRole
	return s.teamMemberRepo.Update(ctx, member)
}

func (s *teamMemberServiceV2) RemoveMember(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) error {
	member, err := s.GetMemberByID(ctx, accountID, memberID)
	if err != nil {
		return err
	}

	// Cannot remove owner
	if member.Role == models.RoleOwner {
		return errors.New("FORBIDDEN", "Cannot remove owner", http.StatusForbidden)
	}

	return s.teamMemberRepo.Delete(ctx, member.ID)
}

func (s *teamMemberServiceV2) ResendInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error {
	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return errors.New("NOT_FOUND", "Invitation not found", http.StatusNotFound)
	}

	if invitation.AccountID != accountID {
		return errors.New("FORBIDDEN", "Cannot resend invitation from another account", http.StatusForbidden)
	}

	if !invitation.IsValid() {
		return errors.New("BAD_REQUEST", "Invitation is no longer valid", http.StatusBadRequest)
	}

	// Get organization name for email
	org, _ := s.accountRepo.FindByID(ctx, invitation.AccountID)
	companyName := "LeCritique"
	if org != nil {
		companyName = org.Name
	}

	// Send invitation email again
	return s.emailService.SendTeamInviteEmail(ctx, invitation.Email, invitation.Token, companyName)
}

func (s *teamMemberServiceV2) CancelInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error {
	invitation, err := s.invitationRepo.FindByToken(ctx, invitationID.String())
	if err != nil {
		return errors.New("NOT_FOUND", "Invitation not found", http.StatusNotFound)
	}

	if invitation.AccountID != accountID {
		return errors.New("FORBIDDEN", "Cannot cancel invitation from another account", http.StatusForbidden)
	}

	return s.invitationRepo.Delete(ctx, invitation.ID)
}

func (s *teamMemberServiceV2) ListPendingInvitations(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error) {
	return s.invitationRepo.FindPendingByAccountID(ctx, accountID)
}

func (s *teamMemberServiceV2) GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error) {
	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN", "Invalid or expired invitation token", http.StatusBadRequest)
	}

	// Check if invitation is still valid
	if !invitation.IsValid() {
		return nil, errors.New("EXPIRED_TOKEN", "This invitation has expired", http.StatusBadRequest)
	}

	return invitation, nil
}

func (s *teamMemberServiceV2) CheckPendingInvitations(ctx context.Context, email string) ([]*models.TeamInvitation, error) {
	return s.invitationRepo.FindByEmail(ctx, email)
}

func (s *teamMemberServiceV2) AcceptInvitation(ctx context.Context, invitationToken string, memberAccountID uuid.UUID) error {
	// Find invitation
	invitation, err := s.invitationRepo.FindByToken(ctx, invitationToken)
	if err != nil {
		fmt.Printf("Error finding invitation by token: %v\n", err)
		return errors.New("INVALID_TOKEN", "Invalid invitation token", http.StatusBadRequest)
	}

	fmt.Printf("Found invitation: ID=%s, Email=%s, AcceptedAt=%v\n", invitation.ID, invitation.Email, invitation.AcceptedAt)

	// Validate invitation
	if !invitation.IsValid() {
		return errors.New("EXPIRED_TOKEN", "Invitation has expired", http.StatusBadRequest)
	}

	// Get member account to verify email matches
	memberAccount, err := s.accountRepo.FindByID(ctx, memberAccountID)
	if err != nil {
		fmt.Printf("Error finding member account: %v\n", err)
		return err
	}

	fmt.Printf("Member account: ID=%s, Email=%s\n", memberAccount.ID, memberAccount.Email)

	if memberAccount.Email != invitation.Email {
		fmt.Printf("Email mismatch: invitation=%s, account=%s\n", invitation.Email, memberAccount.Email)
		return errors.New("EMAIL_MISMATCH", "Invitation email does not match account email", http.StatusBadRequest)
	}

	// Check if already a member
	existingMember, err := s.teamMemberRepo.FindByMemberAndAccount(ctx, memberAccountID, invitation.AccountID)
	if err == nil && existingMember != nil {
		fmt.Printf("User is already a member of this team\n")
		// Already a member, just mark invitation as accepted
		now := time.Now()
		invitation.AcceptedAt = &now
		return s.invitationRepo.Update(ctx, invitation)
	}

	// Create team membership
	member := &models.TeamMember{
		AccountID:  invitation.AccountID,
		MemberID:   memberAccountID,
		Role:       invitation.Role,
		InvitedBy:  invitation.InvitedBy,
		InvitedAt:  invitation.CreatedAt, // Use CreatedAt from BaseModel
		AcceptedAt: func() *time.Time { t := time.Now(); return &t }(),
	}

	fmt.Printf("Creating team member: AccountID=%s, MemberID=%s, Role=%s\n", member.AccountID, member.MemberID, member.Role)

	if err := s.teamMemberRepo.Create(ctx, member); err != nil {
		fmt.Printf("Error creating team member: %v\n", err)
		return err
	}

	// Mark invitation as accepted
	now := time.Now()
	invitation.AcceptedAt = &now
	fmt.Printf("Updating invitation with AcceptedAt=%v\n", now)

	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		fmt.Printf("Error updating invitation: %v\n", err)
		return err
	}

	return nil
}

// UpdateInvitation updates an invitation
func (s *teamMemberServiceV2) UpdateInvitation(ctx context.Context, invitation *models.TeamInvitation) error {
	return s.invitationRepo.Update(ctx, invitation)
}
