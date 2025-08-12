package gormrepo

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
)

type TeamMemberRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) authinterface.TeamMemberRepository {
	return &TeamMemberRepository{db: db}
}

func (r *TeamMemberRepository) Create(ctx context.Context, member *models.TeamMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *TeamMemberRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.WithContext(ctx).Preload("Account").Preload("MemberAccount").Where("id = ?", id).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) FindByAccountAndMember(ctx context.Context, accountID, memberID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("MemberAccount").
		Where("account_id = ? AND member_id = ?", accountID, memberID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) FindByMemberID(ctx context.Context, memberID uuid.UUID) (*models.TeamMember, error) {
	var member models.TeamMember
	err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("MemberAccount").
		Where("member_id = ?", memberID).
		First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *TeamMemberRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]*models.TeamMember, error) {
	var members []*models.TeamMember
	err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("MemberAccount").
		Where("account_id = ?", accountID).
		Find(&members).Error
	return members, err
}

func (r *TeamMemberRepository) Update(ctx context.Context, member *models.TeamMember) error {
	return r.db.WithContext(ctx).Save(member).Error
}

func (r *TeamMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.TeamMember{}, id).Error
}

func (r *TeamMemberRepository) DeleteByMemberID(ctx context.Context, memberID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("member_id = ?", memberID).Delete(&models.TeamMember{}).Error
}