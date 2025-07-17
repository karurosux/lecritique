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
	FindByMemberAndAccount(ctx context.Context, memberID, accountID uuid.UUID) (*models.TeamMember, error)
	FindByMemberID(ctx context.Context, memberID uuid.UUID) ([]models.TeamMember, error)
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
		Preload("Account").
		Preload("MemberAccount").
		Where("account_id = ?", accountID).
		Find(&members).Error
	return members, err
}

func (r *teamMemberRepository) FindByMemberAndAccount(ctx context.Context, memberID, accountID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.DB.WithContext(ctx).
		Where("member_id = ? AND account_id = ? AND deleted_at IS NULL", memberID, accountID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *teamMemberRepository) FindByMemberID(ctx context.Context, memberID uuid.UUID) ([]models.TeamMember, error) {
	var members []models.TeamMember
	err := r.DB.WithContext(ctx).
		Preload("Account").
		Where("member_id = ? AND accepted_at IS NOT NULL", memberID).
		Find(&members).Error
	return members, err
}

func (r *teamMemberRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).
		Model(&models.TeamMember{}).
		Where("account_id = ?", accountID).
		Count(&count).Error
	return count, err
}