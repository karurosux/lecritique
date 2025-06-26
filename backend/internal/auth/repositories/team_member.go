package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type TeamMemberRepository interface {
	Create(ctx context.Context, member *models.TeamMember) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.TeamMember, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error)
	FindByUserAndAccount(ctx context.Context, userID, accountID uuid.UUID) (*models.TeamMember, error)
	Update(ctx context.Context, member *models.TeamMember) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
}

type teamMemberRepository struct {
	*repositories.BaseRepository[models.TeamMember]
}

func NewTeamMemberRepository(db *gorm.DB) TeamMemberRepository {
	return &teamMemberRepository{
		BaseRepository: repositories.NewBaseRepository[models.TeamMember](db),
	}
}

func (r *teamMemberRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.DB.WithContext(ctx).
		Preload("User").
		Where("account_id = ?", accountID).
		Find(&members).Error
	return members, err
}

func (r *teamMemberRepository) FindByUserAndAccount(ctx context.Context, userID, accountID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND account_id = ? AND deleted_at IS NULL", userID, accountID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *teamMemberRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.TeamMember{}).
		Where("account_id = ?", accountID).
		Count(&count).Error
	return count, err
}