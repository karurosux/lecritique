package services

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/auth/models"
	"kyooar/internal/auth/repositories"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/services"
	"github.com/samber/do"
)

type TeamMemberServiceV2 interface {
	ListMembers(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error)
	GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error)
	GetMemberByMemberID(ctx context.Context, memberId uuid.UUID) (*models.TeamMember, error)
	UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error
	RemoveMember(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) error

	InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamInvitation, error)
	ResendInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error
	CancelInvitation(ctx context.Context, accountID uuid.UUID, invitationID uuid.UUID) error
	ListPendingInvitations(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error)
	GetInvitationByToken(ctx context.Context, token string) (*models.TeamInvitation, error)

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
	members, err := s.teamMemberRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	invitations, err := s.invitationRepo.FindPendingByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, inv := range invitations {
		member := models.TeamMember{
			AccountID: inv.AccountID,
			Role:      inv.Role,
			InvitedBy: inv.InvitedBy,
			InvitedAt: inv.CreatedAt,
			MemberAccount: models.Account{
				Email: inv.Email,
			},
		}
		member.ID = inv.ID
		member.CreatedAt = inv.CreatedAt
		member.UpdatedAt = inv.UpdatedAt

		members = append(members, member)
	}

	return members, nil
}

func (s *teamMemberServiceV2) GetMemberByID(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID) (*models.TeamMember, error) {
	member, err := s.teamMemberRepo.FindByID(ctx, memberID, "MemberAccount")
	if err != nil {
		return nil, err
	}

	if member.AccountID != accountID {
		return nil, errors.New("NOT_FOUND", "Member not found", http.StatusNotFound)
	}

	return member, nil
}

func (s *teamMemberServiceV2) InviteMember(ctx context.Context, accountID uuid.UUID, inviterID uuid.UUID, email string, role models.MemberRole) (*models.TeamInvitation, error) {
	if role == models.RoleOwner {
		return nil, errors.New("FORBIDDEN", "Cannot invite another owner", http.StatusForbidden)
	}

	existingAccount, _ := s.accountRepo.FindByEmail(ctx, email)
	if existingAccount != nil {
		existingMember, _ := s.teamMemberRepo.FindByMemberAndAccount(ctx, existingAccount.ID, accountID)
		if existingMember != nil {
			return nil, errors.New("CONFLICT", "User is already a team member", http.StatusConflict)
		}
	}

	existingInvite, _ := s.invitationRepo.FindByAccountAndEmail(ctx, accountID, email)
	if existingInvite != nil && existingInvite.IsValid() {
		return nil, errors.New("CONFLICT", "Invitation already sent to this email", http.StatusConflict)
	}

	token, err := models.GenerateToken()
	if err != nil {
		return nil, err
	}

	invitation := &models.TeamInvitation{
		AccountID: accountID,
		Email:     email,
		Role:      role,
		InvitedBy: inviterID,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	org, _ := s.accountRepo.FindByID(ctx, accountID)
	companyName := "Kyooar"
	if org != nil {
		companyName = org.Name
	}

	if err := s.emailService.SendTeamInviteEmail(ctx, email, token, companyName); err != nil {
		fmt.Printf("Failed to send invitation email to %s: %v\n", email, err)
	}

	return invitation, nil
}

func (s *teamMemberServiceV2) UpdateRole(ctx context.Context, accountID uuid.UUID, memberID uuid.UUID, newRole models.MemberRole) error {
	member, err := s.GetMemberByID(ctx, accountID, memberID)
	if err != nil {
		return err
	}

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

	org, _ := s.accountRepo.FindByID(ctx, invitation.AccountID)
	companyName := "Kyooar"
	if org != nil {
		companyName = org.Name
	}

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

	if !invitation.IsValid() {
		return nil, errors.New("EXPIRED_TOKEN", "This invitation has expired", http.StatusBadRequest)
	}

	return invitation, nil
}

func (s *teamMemberServiceV2) CheckPendingInvitations(ctx context.Context, email string) ([]*models.TeamInvitation, error) {
	return s.invitationRepo.FindByEmail(ctx, email)
}

func (s *teamMemberServiceV2) AcceptInvitation(ctx context.Context, invitationToken string, memberAccountID uuid.UUID) error {
	invitation, err := s.invitationRepo.FindByToken(ctx, invitationToken)
	if err != nil {
		fmt.Printf("Error finding invitation by token: %v\n", err)
		return errors.New("INVALID_TOKEN", "Invalid invitation token", http.StatusBadRequest)
	}

	fmt.Printf("Found invitation: ID=%s, Email=%s, AcceptedAt=%v\n", invitation.ID, invitation.Email, invitation.AcceptedAt)

	if !invitation.IsValid() {
		return errors.New("EXPIRED_TOKEN", "Invitation has expired", http.StatusBadRequest)
	}

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

	existingMember, err := s.teamMemberRepo.FindByMemberAndAccount(ctx, memberAccountID, invitation.AccountID)
	if err == nil && existingMember != nil {
		fmt.Printf("User is already a member of this team\n")
		now := time.Now()
		invitation.AcceptedAt = &now
		return s.invitationRepo.Update(ctx, invitation)
	}

	member := &models.TeamMember{
		AccountID:  invitation.AccountID,
		MemberID:   memberAccountID,
		Role:       invitation.Role,
		InvitedBy:  invitation.InvitedBy,
		InvitedAt:  invitation.CreatedAt,
		AcceptedAt: func() *time.Time { t := time.Now(); return &t }(),
	}

	fmt.Printf("Creating team member: AccountID=%s, MemberID=%s, Role=%s\n", member.AccountID, member.MemberID, member.Role)

	if err := s.teamMemberRepo.Create(ctx, member); err != nil {
		fmt.Printf("Error creating team member: %v\n", err)
		return err
	}

	now := time.Now()
	invitation.AcceptedAt = &now
	fmt.Printf("Updating invitation with AcceptedAt=%v\n", now)

	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		fmt.Printf("Error updating invitation: %v\n", err)
		return err
	}

	return nil
}

func (s *teamMemberServiceV2) UpdateInvitation(ctx context.Context, invitation *models.TeamInvitation) error {
	return s.invitationRepo.Update(ctx, invitation)
}
