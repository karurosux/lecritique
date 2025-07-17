package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type TeamInvitationRepository interface {
	Create(ctx context.Context, invitation *models.TeamInvitation) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.TeamInvitation, error)
	FindByToken(ctx context.Context, token string) (*models.TeamInvitation, error)
	FindByEmail(ctx context.Context, email string) ([]*models.TeamInvitation, error)
	FindByAccountAndEmail(ctx context.Context, accountID uuid.UUID, email string) (*models.TeamInvitation, error)
	FindPendingByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error)
	Update(ctx context.Context, invitation *models.TeamInvitation) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type teamInvitationRepository struct {
	*repositories.BaseRepository[models.TeamInvitation]
}

func NewTeamInvitationRepository(db *gorm.DB) TeamInvitationRepository {
	return &teamInvitationRepository{
		BaseRepository: repositories.NewBaseRepository[models.TeamInvitation](db),
	}
}

func (r *teamInvitationRepository) Create(ctx context.Context, invitation *models.TeamInvitation) error {
	return r.DB.WithContext(ctx).Create(invitation).Error
}

func (r *teamInvitationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.TeamInvitation, error) {
	var invitation models.TeamInvitation
	err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *teamInvitationRepository) FindByToken(ctx context.Context, token string) (*models.TeamInvitation, error) {
	var invitation models.TeamInvitation
	err := r.DB.WithContext(ctx).
		Preload("Account").
		Preload("InvitedByUser").
		Where("token = ?", token).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *teamInvitationRepository) FindByEmail(ctx context.Context, email string) ([]*models.TeamInvitation, error) {
	var invitations []*models.TeamInvitation
	err := r.DB.WithContext(ctx).
		Preload("Account").
		Where("email = ? AND accepted_at IS NULL", email).
		Find(&invitations).Error
	return invitations, err
}

func (r *teamInvitationRepository) FindByAccountAndEmail(ctx context.Context, accountID uuid.UUID, email string) (*models.TeamInvitation, error) {
	var invitation models.TeamInvitation
	err := r.DB.WithContext(ctx).
		Where("account_id = ? AND email = ? AND accepted_at IS NULL", accountID, email).
		First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *teamInvitationRepository) Update(ctx context.Context, invitation *models.TeamInvitation) error {
	return r.DB.WithContext(ctx).Save(invitation).Error
}

func (r *teamInvitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.TeamInvitation{}, "id = ?", id).Error
}

func (r *teamInvitationRepository) FindPendingByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamInvitation, error) {
	var invitations []*models.TeamInvitation
	err := r.DB.WithContext(ctx).
		Where("account_id = ? AND accepted_at IS NULL AND expires_at > NOW()", accountID).
		Find(&invitations).Error
	return invitations, err
}

func (r *teamInvitationRepository) DeleteExpired(ctx context.Context) error {
	return r.DB.WithContext(ctx).
		Where("expires_at < NOW() AND accepted_at IS NULL").
		Delete(&models.TeamInvitation{}).Error
}