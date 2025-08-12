package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
)

type TeamInvitationRepository struct {
	db *gorm.DB
}

func NewTeamInvitationRepository(db *gorm.DB) authinterface.TeamInvitationRepository {
	return &TeamInvitationRepository{db: db}
}

func (r *TeamInvitationRepository) Create(ctx context.Context, invitation *models.TeamInvitation) error {
	return r.db.WithContext(ctx).Create(invitation).Error
}

func (r *TeamInvitationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.TeamInvitation, error) {
	var invitation models.TeamInvitation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *TeamInvitationRepository) FindByToken(ctx context.Context, token string) (*models.TeamInvitation, error) {
	var invitation models.TeamInvitation
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *TeamInvitationRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error) {
	var invitations []*models.TeamInvitation
	err := r.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&invitations).Error
	return invitations, err
}

func (r *TeamInvitationRepository) Update(ctx context.Context, invitation *models.TeamInvitation) error {
	return r.db.WithContext(ctx).Save(invitation).Error
}

func (r *TeamInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.TeamInvitation{}, id).Error
}