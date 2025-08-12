package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	feedbackinterface "kyooar/internal/feedback/interface"
	feedbackmodel "kyooar/internal/feedback/model"
	organizationRepos "kyooar/internal/organization/repositories"
	qrcodeRepos "kyooar/internal/qrcode/repositories"
	sharedModels "kyooar/internal/shared/models"
)

type feedbackService struct {
	feedbackRepo     feedbackinterface.FeedbackRepository
	organizationRepo organizationRepos.OrganizationRepository
	qrCodeRepo       qrcodeRepos.QRCodeRepository
}

func NewFeedbackService(
	feedbackRepo feedbackinterface.FeedbackRepository,
	organizationRepo organizationRepos.OrganizationRepository,
	qrCodeRepo qrcodeRepos.QRCodeRepository,
) feedbackinterface.FeedbackService {
	return &feedbackService{
		feedbackRepo:     feedbackRepo,
		organizationRepo: organizationRepo,
		qrCodeRepo:       qrCodeRepo,
	}
}

func (s *feedbackService) Submit(ctx context.Context, feedback *feedbackmodel.Feedback) error {
	qrCode, err := s.qrCodeRepo.FindByID(ctx, feedback.QRCodeID)
	if err != nil {
		return err
	}

	feedback.OrganizationID = qrCode.OrganizationID
	// ProductID comes from the request payload, QRCodeID identifies the location

	return s.feedbackRepo.Create(ctx, feedback)
}

func (s *feedbackService) GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[feedbackmodel.Feedback], error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, fmt.Errorf("organization not found")
	}

	return s.feedbackRepo.FindByOrganizationID(ctx, accountID, organizationID, page, limit)
}

func (s *feedbackService) GetByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackmodel.FeedbackFilter) (*sharedModels.PageResponse[feedbackmodel.Feedback], error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, fmt.Errorf("organization not found")
	}

	return s.feedbackRepo.FindByOrganizationIDWithFilters(ctx, accountID, organizationID, page, limit, filters)
}

func (s *feedbackService) GetStats(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*feedbackmodel.FeedbackStats, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, fmt.Errorf("organization not found")
	}

	return s.feedbackRepo.GetStatsByOrganization(ctx, accountID, organizationID)
}

func (s *feedbackService) GetByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error) {
	return s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, limit)
}

func (s *feedbackService) GetByQuestionInPeriod(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time) ([]feedbackmodel.Feedback, error) {
	return s.feedbackRepo.FindByQuestionInPeriod(ctx, questionID, startDate, endDate)
}