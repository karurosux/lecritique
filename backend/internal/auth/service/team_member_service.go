package authservice

import (
	"context"
	"time"

	"github.com/google/uuid"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/errors"
	sharedmodels "kyooar/internal/shared/models"
)

type TeamMemberService struct {
	teamMemberRepo   authinterface.TeamMemberRepository
	invitationRepo   authinterface.TeamInvitationRepository
	accountRepo      authinterface.AccountRepository
	tokenGenerator   authinterface.TokenGenerator
	emailSender      authinterface.EmailSender
}

func NewTeamMemberService(
	teamMemberRepo authinterface.TeamMemberRepository,
	invitationRepo authinterface.TeamInvitationRepository,
	accountRepo authinterface.AccountRepository,
	tokenGenerator authinterface.TokenGenerator,
	emailSender authinterface.EmailSender,
) authinterface.TeamMemberService {
	return &TeamMemberService{
		teamMemberRepo: teamMemberRepo,
		invitationRepo: invitationRepo,
		accountRepo:    accountRepo,
		tokenGenerator: tokenGenerator,
		emailSender:    emailSender,
	}
}

func (s *TeamMemberService) GetMemberByMemberID(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error) {
	return s.teamMemberRepo.FindByMemberID(ctx, memberID)
}

func (s *TeamMemberService) ListMembers(ctx context.Context, accountID uuid.UUID) ([]*models.TeamMember, error) {
	members, err := s.teamMemberRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	pendingInvitations, err := s.invitationRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, invitation := range pendingInvitations {
		if invitation.AcceptedAt != nil {
			continue
		}

		pendingMember := &models.TeamMember{
			BaseModel: sharedmodels.BaseModel{
				ID: invitation.ID,
			},
			AccountID: invitation.AccountID,
			Role:      invitation.Role,
			InvitedBy: invitation.InvitedBy,
			InvitedAt: invitation.InvitedAt,
			AcceptedAt: nil,
			MemberAccount: models.Account{
				Email: invitation.Email,
			},
		}
		members = append(members, pendingMember)
	}

	return members, nil
}

func (s *TeamMemberService) InviteMember(ctx context.Context, accountID uuid.UUID, email string, role models.MemberRole, invitedBy uuid.UUID) (*models.TeamInvitation, error) {
	token, err := s.tokenGenerator.GenerateSecureToken()
	if err != nil {
		return nil, err
	}

	invitation := &models.TeamInvitation{
		AccountID: accountID,
		Email:     email,
		Role:      role,
		Token:     token,
		InvitedBy: invitedBy,
		InvitedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
	}

	if err := s.invitationRepo.Create(ctx, invitation); err != nil {
		return nil, err
	}

	return invitation, nil
}

func (s *TeamMemberService) ResendInvitation(ctx context.Context, accountID, invitationID uuid.UUID) error {
	invitation, err := s.invitationRepo.FindByID(ctx, invitationID)
	if err != nil {
		return err
	}

	if invitation.AccountID != accountID {
		return errors.ErrUnauthorized
	}

	newToken, err := s.tokenGenerator.GenerateSecureToken()
	if err != nil {
		return err
	}

	invitation.Token = newToken
	invitation.InvitedAt = time.Now()
	invitation.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	return s.invitationRepo.Update(ctx, invitation)
}

func (s *TeamMemberService) UpdateRole(ctx context.Context, accountID, memberID uuid.UUID, role models.MemberRole) error {
	member, err := s.teamMemberRepo.FindByAccountAndMember(ctx, accountID, memberID)
	if err != nil {
		return err
	}

	member.Role = role
	return s.teamMemberRepo.Update(ctx, member)
}

func (s *TeamMemberService) UpdateRoleByID(ctx context.Context, teamMemberID uuid.UUID, role models.MemberRole) error {
	member, err := s.teamMemberRepo.FindByID(ctx, teamMemberID)
	if err != nil {
		return err
	}

	member.Role = role
	return s.teamMemberRepo.Update(ctx, member)
}

func (s *TeamMemberService) RemoveMember(ctx context.Context, accountID, memberID uuid.UUID) error {
	member, err := s.teamMemberRepo.FindByAccountAndMember(ctx, accountID, memberID)
	if err != nil {
		return err
	}

	return s.teamMemberRepo.Delete(ctx, member.ID)
}

func (s *TeamMemberService) RemoveMemberByID(ctx context.Context, teamMemberID uuid.UUID) error {
	// First try to find it as a team invitation
	_, err := s.invitationRepo.FindByID(ctx, teamMemberID)
	if err == nil {
		// It's a pending invitation
		return s.invitationRepo.Delete(ctx, teamMemberID)
	}

	// If not found as invitation, try as team member
	_, err = s.teamMemberRepo.FindByID(ctx, teamMemberID)
	if err != nil {
		return err
	}

	return s.teamMemberRepo.Delete(ctx, teamMemberID)
}

func (s *TeamMemberService) AcceptInvitation(ctx context.Context, token string) (*models.TeamMember, error) {
	invitation, err := s.invitationRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired invitation token", 400, nil)
	}

	if invitation.ExpiresAt.Before(time.Now()) {
		return nil, errors.NewWithDetails("INVITATION_EXPIRED", "Invitation has expired", 400, nil)
	}

	if invitation.AcceptedAt != nil {
		return nil, errors.NewWithDetails("INVITATION_ALREADY_ACCEPTED", "Invitation has already been accepted", 400, nil)
	}

	// Find or create the member account
	memberAccount, err := s.accountRepo.FindByEmail(ctx, invitation.Email)
	if err != nil {
		return nil, errors.NewWithDetails("ACCOUNT_NOT_FOUND", "Account not found for this email", 404, nil)
	}

	// Create team member relationship
	member := &models.TeamMember{
		AccountID: invitation.AccountID,
		MemberID:  memberAccount.ID,
		Role:      invitation.Role,
		InvitedBy: invitation.InvitedBy,
		InvitedAt: invitation.InvitedAt,
		AcceptedAt: &time.Time{},
	}
	now := time.Now()
	member.AcceptedAt = &now

	if err := s.teamMemberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	// Mark invitation as accepted
	invitation.AcceptedAt = &now
	if err := s.invitationRepo.Update(ctx, invitation); err != nil {
		return nil, err
	}

	// Load the member with relationships
	return s.teamMemberRepo.FindByID(ctx, member.ID)
}