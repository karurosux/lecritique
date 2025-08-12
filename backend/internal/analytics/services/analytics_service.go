package analyticsservice

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
	analyticsModels "kyooar/internal/analytics/model"
	analyticsinterface "kyooar/internal/analytics/interface"
	feedbackmodel "kyooar/internal/feedback/model"
	feedbackinterface "kyooar/internal/feedback/interface"
	menuRepos "kyooar/internal/menu/repositories"
	qrcodemodel "kyooar/internal/qrcode/model"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	organizationinterface "kyooar/internal/organization/interface"
	"kyooar/internal/shared/logger"
	"github.com/sirupsen/logrus"
)

type AnalyticsService struct {
	analyticsRepo    analyticsinterface.AnalyticsRepository
	feedbackRepo     feedbackinterface.FeedbackRepository
	productRepo      menuRepos.ProductRepository
	qrCodeRepo       qrcodeinterface.QRCodeRepository
	organizationRepo organizationinterface.OrganizationRepository
}

func NewAnalyticsService(
	analyticsRepo analyticsinterface.AnalyticsRepository,
	feedbackRepo feedbackinterface.FeedbackRepository,
	productRepo menuRepos.ProductRepository,
	qrCodeRepo qrcodeinterface.QRCodeRepository,
	organizationRepo organizationinterface.OrganizationRepository,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:    analyticsRepo,
		feedbackRepo:     feedbackRepo,
		productRepo:      productRepo,
		qrCodeRepo:       qrCodeRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *AnalyticsService) GetDashboardMetrics(ctx context.Context, organizationID uuid.UUID) (*analyticsModels.DashboardMetrics, error) {
	metrics := &analyticsModels.DashboardMetrics{}
	
	feedbackCounts, err := s.analyticsRepo.GetFeedbackCounts(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get feedback counts", err, logrus.Fields{
			"organization_id": organizationID,
		})
	} else {
		metrics.TotalFeedbacks = feedbackCounts.Total
		metrics.TodaysFeedback = feedbackCounts.Today
		
		if feedbackCounts.Yesterday > 0 {
			metrics.TrendVsYesterday = float64(feedbackCounts.Today-feedbackCounts.Yesterday) / float64(feedbackCounts.Yesterday) * 100
		}
	}
	
	qrMetrics, err := s.analyticsRepo.GetQRCodeMetrics(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get QR code metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
	} else {
		metrics.ActiveQRCodes = qrMetrics.ActiveCount
		metrics.TotalQRScans = qrMetrics.TotalScans
		metrics.ScansToday = qrMetrics.ScansToday
	}
	
	if metrics.TotalQRScans > 0 && feedbackCounts != nil && feedbackCounts.Recent30Days > 0 {
		rate := float64(feedbackCounts.Recent30Days) / float64(metrics.TotalQRScans) * 100
		if rate > 100 {
			metrics.CompletionRate = 100.0
		} else {
			metrics.CompletionRate = rate
		}
	} else {
		metrics.CompletionRate = 0.0
	}
	
	feedback, err := s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, 1000)
	if err != nil {
		logger.Error("Failed to get feedback for dashboard metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		feedback = []feedbackmodel.Feedback{}
	}
	
	metrics.DeviceBreakdown = s.getDeviceMetricsFromFeedback(feedback)
	
	metrics.AverageResponseTime = s.getAverageResponseTimeFromFeedback(ctx, feedback)
	
	metrics.PeakHours = s.getPeakUsageHoursFromFeedback(feedback)
	
	qrPerformance, err := s.getQRCodePerformance(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get QR performance", err, logrus.Fields{
			"organization_id": organizationID,
		})
	} else {
		metrics.QRPerformance = qrPerformance
	}
	
	return metrics, nil
}

func (s *AnalyticsService) GetProductInsights(ctx context.Context, productID uuid.UUID) (*analyticsModels.ProductInsights, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	
	allFeedback, err := s.feedbackRepo.FindByProductIDForAnalytics(ctx, productID, 500)
	if err != nil {
		return nil, err
	}
	
	insights := &analyticsModels.ProductInsights{
		ProductID:        productID,
		ProductName:      product.Name,
		TotalFeedback: int64(len(allFeedback)),
	}
	
	if len(allFeedback) == 0 {
		return insights, nil
	}
	
	questions, err := s.feedbackRepo.GetQuestionsByProductID(ctx, productID)
	if err != nil {
		logger.Error("Failed to get questions", err, logrus.Fields{"product_id": productID})
	}
	
	questionMetrics := s.aggregateQuestionMetrics(allFeedback, questions)
	insights.Questions = questionMetrics
	
	var totalScore float64
	var scoreCount int
	for _, q := range questionMetrics {
		if q.AverageScore != nil {
			totalScore += *q.AverageScore
			scoreCount++
		}
	}
	if scoreCount > 0 {
		insights.OverallScore = totalScore / float64(scoreCount)
	}
	
	insights.BestAspects = s.identifyBestAspects(questionMetrics)
	insights.NeedsAttention = s.identifyNeedsAttention(questionMetrics)
	
	insights.CompletionRate = 100.0
	
	if len(allFeedback) > 0 {
		insights.LastFeedback = allFeedback[0].CreatedAt
	}
	
	weekAgo := time.Now().AddDate(0, 0, -7)
	weekAgoFeedback := filterFeedbackByDate(allFeedback, weekAgo)
	if len(weekAgoFeedback) > 0 {
		weekScore, _, _ := s.calculateOverallMetrics(weekAgoFeedback)
		insights.WeeklyChange = ((insights.OverallScore - weekScore) / weekScore) * 100
	}
	
	if insights.WeeklyChange > 5 {
		insights.ScoreTrend = "improving"
	} else if insights.WeeklyChange < -5 {
		insights.ScoreTrend = "declining"
	} else {
		insights.ScoreTrend = "stable"
	}
	
	return insights, nil
}

func (s *AnalyticsService) GetOrganizationInsights(ctx context.Context, organizationID uuid.UUID, period string) (*analyticsModels.OrganizationInsights, error) {
	return &analyticsModels.OrganizationInsights{
		OrganizationID:   organizationID,
		Period:         period,
		TotalFeedback:  0,
		ActiveProducts:   0,
	}, nil
}


func (s *AnalyticsService) calculateOverallMetrics(feedback []feedbackmodel.Feedback) (satisfaction, recommendRate, sentiment float64) {
	if len(feedback) == 0 {
		return 0, 0, 0
	}
	
	var totalScore float64
	var scoreCount int
	var recommendYes int
	var recommendTotal int
	var positiveCount int
	
	for _, f := range feedback {
		if f.OverallRating > 0 {
			totalScore += float64(f.OverallRating)
			scoreCount++
			if f.OverallRating >= 4 {
				positiveCount++
			}
		}
		
		for _, response := range f.Responses {
			switch v := response.Answer.(type) {
			case bool:
					recommendTotal++
				if v {
					recommendYes++
				}
			case string:
						lower := strings.ToLower(v)
				if strings.Contains(lower, "excellent") || strings.Contains(lower, "great") || 
				   strings.Contains(lower, "amazing") || strings.Contains(lower, "love") {
					positiveCount++
				}
			}
		}
	}
	
	if scoreCount > 0 {
		satisfaction = totalScore / float64(scoreCount)
	}
	
	if recommendTotal > 0 {
		recommendRate = float64(recommendYes) / float64(recommendTotal) * 100
	}
	
	totalResponses := len(feedback)
	if totalResponses > 0 {
		sentiment = float64(positiveCount) / float64(totalResponses) * 100
	}
	
	return
}

func (s *AnalyticsService) aggregateQuestionMetrics(feedback []feedbackmodel.Feedback, questions []feedbackmodel.Question) []analyticsModels.QuestionMetric {
	metricsMap := make(map[uuid.UUID]*analyticsModels.QuestionMetric)
	
	for _, q := range questions {
		metric := &analyticsModels.QuestionMetric{
			QuestionID:         q.ID,
			QuestionText:       q.Text,
			QuestionType:       string(q.Type),
			OptionDistribution: make(map[string]int64),
		}
		metricsMap[q.ID] = metric
	}
	
	for _, f := range feedback {
		for _, response := range f.Responses {
			metric, exists := metricsMap[response.QuestionID]
			if !exists {
				continue
			}
			
			metric.ResponseCount++
			
			switch v := response.Answer.(type) {
			case float64:
				if metric.AverageScore == nil {
					avg := v
					min := v
					max := v
					metric.AverageScore = &avg
					metric.MinScore = &min
					metric.MaxScore = &max
				} else {
					*metric.AverageScore = (*metric.AverageScore * float64(metric.ResponseCount-1) + v) / float64(metric.ResponseCount)
					if v < *metric.MinScore {
						*metric.MinScore = v
					}
					if v > *metric.MaxScore {
						*metric.MaxScore = v
					}
				}
				
				if v >= 4 {
					metric.PositiveRate++
				} else if v >= 3 {
					metric.NeutralRate++
				} else {
					metric.NegativeRate++
				}
				
			case string:
				if metric.QuestionType == string(feedbackmodel.QuestionTypeText) {
					metric.TextResponses = append(metric.TextResponses, v)
				} else {
					metric.OptionDistribution[v]++
				}
				
			case bool:
				if v {
					metric.OptionDistribution["Yes"]++
					metric.PositiveRate++
				} else {
					metric.OptionDistribution["No"]++
					metric.NegativeRate++
				}
			}
		}
	}
	
	for _, metric := range metricsMap {
		if metric.ResponseCount > 0 {
			total := float64(metric.ResponseCount)
			metric.PositiveRate = (metric.PositiveRate / total) * 100
			metric.NeutralRate = (metric.NeutralRate / total) * 100
			metric.NegativeRate = (metric.NegativeRate / total) * 100
		}
	}
	
	var metrics []analyticsModels.QuestionMetric
	for _, m := range metricsMap {
		metrics = append(metrics, *m)
	}
	
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ResponseCount > metrics[j].ResponseCount
	})
	
	return metrics
}

func (s *AnalyticsService) identifyTopIssues(feedback []feedbackmodel.Feedback) []analyticsModels.QuickIssue {
	issues := []analyticsModels.QuickIssue{}
	
	lowRatingCount := 0
	for _, f := range feedback {
		if f.OverallRating <= 2 {
			lowRatingCount++
		}
	}
	
	if lowRatingCount > 0 {
		issues = append(issues, analyticsModels.QuickIssue{
			Title:    "Low ratings",
			Count:    lowRatingCount,
			Severity: "critical",
			ActionLink: "/feedback/manage?filter=low",
		})
	}
	
	
	return issues
}

func (s *AnalyticsService) aggregateByProduct(feedback []feedbackmodel.Feedback) map[uuid.UUID]*analyticsModels.ProductSummary {
	productMap := make(map[uuid.UUID]*analyticsModels.ProductSummary)
	
	for _, f := range feedback {
		summary, exists := productMap[f.ProductID]
		if !exists {
			summary = &analyticsModels.ProductSummary{
				ProductID:   f.ProductID,
				ProductName: f.Product.Name,
			}
			productMap[f.ProductID] = summary
		}
		
		summary.FeedbackCount++
		summary.Score = (summary.Score*float64(summary.FeedbackCount-1) + float64(f.OverallRating)) / float64(summary.FeedbackCount)
	}
	
	return productMap
}

func (s *AnalyticsService) getTopProducts(productMap map[uuid.UUID]*analyticsModels.ProductSummary, limit int) []analyticsModels.ProductSummary {
	var products []analyticsModels.ProductSummary
	for _, d := range productMap {
		products = append(products, *d)
	}
	
	sort.Slice(products, func(i, j int) bool {
		return products[i].Score > products[j].Score
	})
	
	if len(products) > limit {
		products = products[:limit]
	}
	
	return products
}

func (s *AnalyticsService) getBottomProducts(productMap map[uuid.UUID]*analyticsModels.ProductSummary, limit int) []analyticsModels.ProductSummary {
	var products []analyticsModels.ProductSummary
	for _, d := range productMap {
		if d.Score < 3.5 {
			products = append(products, *d)
		}
	}
	
	sort.Slice(products, func(i, j int) bool {
		return products[i].Score < products[j].Score
	})
	
	if len(products) > limit {
		products = products[:limit]
	}
	
	return products
}

func (s *AnalyticsService) getRecentFeedbackSummaries(feedback []feedbackmodel.Feedback, limit int) []analyticsModels.FeedbackSummary {
	sort.Slice(feedback, func(i, j int) bool {
		return feedback[i].CreatedAt.After(feedback[j].CreatedAt)
	})
	
	summaries := []analyticsModels.FeedbackSummary{}
	for i, f := range feedback {
		if i >= limit {
			break
		}
		
		summary := analyticsModels.FeedbackSummary{
			FeedbackID:   f.ID,
			ProductName:     f.Product.Name,
			CustomerName: f.CustomerName,
			Score:        float64(f.OverallRating),
			CreatedAt:    f.CreatedAt,
		}
		
		if f.OverallRating >= 4 {
			summary.Sentiment = "positive"
		} else if f.OverallRating >= 3 {
			summary.Sentiment = "neutral"
		} else {
			summary.Sentiment = "negative"
		}
		
		for _, r := range f.Responses {
			if str, ok := r.Answer.(string); ok && len(str) > 20 {
				summary.Highlight = str
				if len(summary.Highlight) > 100 {
					summary.Highlight = summary.Highlight[:100] + "..."
				}
				break
			}
		}
		
		summaries = append(summaries, summary)
	}
	
	return summaries
}

func (s *AnalyticsService) identifyBestAspects(metrics []analyticsModels.QuestionMetric) []string {
	var best []string
	
	for _, m := range metrics {
		if m.AverageScore != nil && *m.AverageScore >= 4.5 {
			best = append(best, m.QuestionText)
		} else if m.PositiveRate >= 80 {
			best = append(best, m.QuestionText)
		}
		
		if len(best) >= 3 {
			break
		}
	}
	
	return best
}

func (s *AnalyticsService) identifyNeedsAttention(metrics []analyticsModels.QuestionMetric) []string {
	var needs []string
	
	for _, m := range metrics {
		if m.AverageScore != nil && *m.AverageScore < 3 {
			needs = append(needs, m.QuestionText)
		} else if m.NegativeRate >= 40 {
			needs = append(needs, m.QuestionText)
		}
		
		if len(needs) >= 3 {
			break
		}
	}
	
	return needs
}

func filterFeedbackByDate(feedback []feedbackmodel.Feedback, since time.Time) []feedbackmodel.Feedback {
	var filtered []feedbackmodel.Feedback
	for _, f := range feedback {
		if f.CreatedAt.After(since) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}


type QRMetrics struct {
	ActiveCount int64
	TotalScans  int64
	ScansToday  int64
}

func (s *AnalyticsService) getQRCodeMetrics(ctx context.Context, organizationID uuid.UUID) (*QRMetrics, error) {
	qrCodes, err := s.qrCodeRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	
	metrics := &QRMetrics{}
	todayStart := time.Now().Truncate(24 * time.Hour)
	
	for _, qr := range qrCodes {
		if qr.IsActive {
			metrics.ActiveCount++
		}
		metrics.TotalScans += int64(qr.ScansCount)
		
		if qr.LastScannedAt != nil && qr.LastScannedAt.After(todayStart) {
			metrics.ScansToday++
		}
	}
	
	return metrics, nil
}

func (s *AnalyticsService) getDeviceMetrics(ctx context.Context, organizationID uuid.UUID) (map[string]int64, error) {
	feedback, err := s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, 1000)
	if err != nil {
		return nil, err
	}
	
	return s.getDeviceMetricsFromFeedback(feedback), nil
}

func (s *AnalyticsService) getDeviceMetricsFromFeedback(feedback []feedbackmodel.Feedback) map[string]int64 {
	deviceBreakdown := make(map[string]int64)
	platformBreakdown := make(map[string]int64)
	
	for _, f := range feedback {
		if f.DeviceInfo.Platform != "" {
			platformBreakdown[f.DeviceInfo.Platform]++
		}
		
		if f.DeviceInfo.Browser != "" {
			deviceBreakdown[f.DeviceInfo.Browser]++
		}
	}
	
	result := make(map[string]int64)
	for k, v := range platformBreakdown {
		result[k] = v
	}
	for k, v := range deviceBreakdown {
		result[k+" Browser"] = v
	}
	
	return result
}

func (s *AnalyticsService) getAverageResponseTime(ctx context.Context, organizationID uuid.UUID) (float64, error) {
	feedback, err := s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, 500)
	if err != nil {
		return 0, err
	}
	
	return s.getAverageResponseTimeFromFeedback(ctx, feedback), nil
}

func (s *AnalyticsService) getAverageResponseTimeFromFeedback(ctx context.Context, feedback []feedbackmodel.Feedback) float64 {
	if len(feedback) == 0 {
		return 0
	}
	
	qrCodeIDMap := make(map[uuid.UUID]bool)
	for _, f := range feedback {
		if f.QRCodeID != uuid.Nil {
			qrCodeIDMap[f.QRCodeID] = true
		}
	}
	
	if len(qrCodeIDMap) == 0 {
		return 0
	}
	
	var qrCodeIDs []uuid.UUID
	for id := range qrCodeIDMap {
		qrCodeIDs = append(qrCodeIDs, id)
	}
	
	qrCodes, err := s.qrCodeRepo.FindByIDs(ctx, qrCodeIDs)
	if err != nil {
		logger.Error("Failed to get QR codes for response time calculation", err, logrus.Fields{
			"qr_code_count": len(qrCodeIDs),
		})
		return 0
	}
	
	qrCodeMap := make(map[uuid.UUID]*qrcodemodel.QRCode)
	for i := range qrCodes {
		qrCodeMap[qrCodes[i].ID] = &qrCodes[i]
	}
	
	var totalTime float64
	var count int
	
	for _, f := range feedback {
		qrCode, exists := qrCodeMap[f.QRCodeID]
		if !exists || qrCode.LastScannedAt == nil {
			continue
		}
		
		responseTime := f.CreatedAt.Sub(*qrCode.LastScannedAt)
		if responseTime > 0 && responseTime < 24*time.Hour {
			totalTime += responseTime.Minutes()
			count++
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return totalTime / float64(count)
}

func (s *AnalyticsService) getPeakUsageHours(ctx context.Context, organizationID uuid.UUID) ([]int, error) {
	feedback, err := s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, 1000)
	if err != nil {
		return nil, err
	}
	
	return s.getPeakUsageHoursFromFeedback(feedback), nil
}

func (s *AnalyticsService) getPeakUsageHoursFromFeedback(feedback []feedbackmodel.Feedback) []int {
	weekAgo := time.Now().AddDate(0, 0, -7)
	recentFeedback := filterFeedbackByDate(feedback, weekAgo)
	
	hourCounts := make(map[int]int64)
	for _, f := range recentFeedback {
		hour := f.CreatedAt.Hour()
		hourCounts[hour]++
	}
	
	type hourCount struct {
		hour  int
		count int64
	}
	
	var hours []hourCount
	for h, c := range hourCounts {
		hours = append(hours, hourCount{hour: h, count: c})
	}
	
	sort.Slice(hours, func(i, j int) bool {
		return hours[i].count > hours[j].count
	})
	
	var peakHours []int
	for i, h := range hours {
		if i >= 3 {
			break
		}
		peakHours = append(peakHours, h.hour)
	}
	
	return peakHours
}

func (s *AnalyticsService) getQRCodePerformance(ctx context.Context, organizationID uuid.UUID) ([]analyticsModels.QRCodePerformance, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	
	qrCodes, err := s.qrCodeRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	
	if len(qrCodes) == 0 {
		return []analyticsModels.QRCodePerformance{}, nil
	}
	
	var qrCodeIDs []uuid.UUID
	for _, qr := range qrCodes {
		qrCodeIDs = append(qrCodeIDs, qr.ID)
	}
	
	feedbackCounts, err := s.feedbackRepo.CountByQRCodeIDs(ctx, qrCodeIDs)
	if err != nil {
		logger.Error("Failed to count feedback for QR codes", err, logrus.Fields{
			"organization_id": organizationID,
		})
		feedbackCounts = make(map[uuid.UUID]int64)
	}
	
	var performance []analyticsModels.QRCodePerformance
	
	for _, qr := range qrCodes {
		feedbackCount := feedbackCounts[qr.ID]
		
		conversionRate := 0.0
		if qr.ScansCount > 0 && feedbackCount > 0 {
			rate := (float64(feedbackCount) / float64(qr.ScansCount)) * 100
			if rate > 100 {
				conversionRate = 100.0
			} else {
				conversionRate = rate
			}
		}
		
		perf := analyticsModels.QRCodePerformance{
			ID:             qr.ID,
			Label:          qr.Label,
			OrganizationID:   qr.OrganizationID,
			OrganizationName: organization.Name,
			ScansCount:     int64(qr.ScansCount),
			FeedbackCount:  feedbackCount,
			ConversionRate: conversionRate,
			LastScan:       qr.LastScannedAt,
			IsActive:       qr.IsActive,
		}
		
		if qr.Location != nil {
			perf.Location = *qr.Location
		}
		
		performance = append(performance, perf)
	}
	
	sort.Slice(performance, func(i, j int) bool {
		return performance[i].ScansCount > performance[j].ScansCount
	})
	
	return performance, nil
}


func (s *AnalyticsService) GetOrganizationChartData(ctx context.Context, organizationID uuid.UUID, filters map[string]interface{}) (*analyticsModels.OrganizationChartData, error) {
	logger.Info("Starting GetOrganizationChartData", logrus.Fields{
		"organization_id": organizationID,
		"filters": filters,
	})
	
	feedbackFilters := s.buildFeedbackFilters(filters)
	logger.Info("Built feedback filters", logrus.Fields{
		"feedback_filters": feedbackFilters,
	})
	
	feedback, err := s.feedbackRepo.FindByOrganizationIDWithFilters(ctx, uuid.Nil, organizationID, 1, 10000, feedbackFilters)
	if err != nil {
		logger.Error("Failed to fetch feedback data", err, logrus.Fields{
			"organization_id": organizationID,
			"feedback_filters": feedbackFilters,
			"error_type": fmt.Sprintf("%T", err),
		})
		return nil, fmt.Errorf("failed to fetch feedback data: %w", err)
	}
	
	logger.Info("Fetched feedback data", logrus.Fields{
		"organization_id": organizationID,
		"feedback_count": len(feedback.Data),
		"total_feedback": feedback.Total,
	})
	
	products, err := s.productRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		logger.Error("Failed to get products", err, logrus.Fields{
			"organization_id": organizationID,
			"error_type": fmt.Sprintf("%T", err),
		})
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	
	logger.Info("Fetched products", logrus.Fields{
		"organization_id": organizationID,
		"product_count": len(products),
	})
	
	productMap := make(map[uuid.UUID]string)
	if products != nil {
		for _, product := range products {
			productMap[product.ID] = product.Name
		}
	}
	
	chartData := &analyticsModels.OrganizationChartData{
		OrganizationID: organizationID,
		Charts:       []analyticsModels.ChartData{},
	}
	
	type questionProductKey struct {
		QuestionID uuid.UUID
		ProductID     uuid.UUID
	}
	
	questionProductResponses := make(map[questionProductKey][]feedbackmodel.Response)
	questionProductMeta := make(map[questionProductKey]struct {
		Text     string
		Type     string
		ProductID   uuid.UUID
		ProductName string
	})
	
	for _, f := range feedback.Data {
		productName := productMap[f.ProductID]
		for _, response := range f.Responses {
			key := questionProductKey{
				QuestionID: response.QuestionID,
				ProductID:     f.ProductID,
			}
			questionProductResponses[key] = append(questionProductResponses[key], response)
			if questionProductMeta[key].Text == "" {
				questionProductMeta[key] = struct {
					Text     string
					Type     string
					ProductID   uuid.UUID
					ProductName string
				}{
					Text:     response.QuestionText,
					Type:     string(response.QuestionType),
					ProductID:   f.ProductID,
					ProductName: productName,
				}
			}
		}
	}
	
	for key, responses := range questionProductResponses {
		meta := questionProductMeta[key]
		chart := s.aggregateQuestionResponses(key.QuestionID, meta.Text, meta.Type, responses)
		if meta.ProductID != uuid.Nil {
			chart.ProductID = &meta.ProductID
			chart.ProductName = meta.ProductName
		}
		chartData.Charts = append(chartData.Charts, chart)
	}
	
	chartData.Summary.TotalResponses = int64(len(feedback.Data))
	chartData.Summary.FiltersApplied = filters
	
	if len(feedback.Data) > 0 {
		chartData.Summary.DateRange.Start = feedback.Data[len(feedback.Data)-1].CreatedAt
		chartData.Summary.DateRange.End = feedback.Data[0].CreatedAt
	}
	
	return chartData, nil
}

func (s *AnalyticsService) GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*analyticsModels.ChartData, error) {
	feedbackFilters := s.buildFeedbackFilters(filters)
	
	
	allFeedback, err := s.feedbackRepo.FindByOrganizationIDWithFilters(ctx, uuid.Nil, uuid.Nil, 1, 10000, feedbackFilters)
	if err != nil {
		return nil, err
	}
	
	var responses []feedbackmodel.Response
	var questionText, questionType string
	
	for _, f := range allFeedback.Data {
		for _, response := range f.Responses {
			if response.QuestionID == questionID {
				responses = append(responses, response)
				if questionText == "" {
					questionText = response.QuestionText
					questionType = string(response.QuestionType)
				}
			}
		}
	}
	
	chartData := s.aggregateQuestionResponses(questionID, questionText, questionType, responses)
	return &chartData, nil
}


func (s *AnalyticsService) buildFeedbackFilters(filters map[string]interface{}) feedbackmodel.FeedbackFilter {
	feedbackFilters := feedbackmodel.FeedbackFilter{}
	
	if dateFrom, ok := filters["date_from"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", dateFrom); err == nil {
			feedbackFilters.DateFrom = &parsed
		}
	}
	
	if dateTo, ok := filters["date_to"].(string); ok {
		if parsed, err := time.Parse("2006-01-02", dateTo); err == nil {
			feedbackFilters.DateTo = &parsed
		}
	}
	
	if productIDStr, ok := filters["product_id"].(string); ok {
		if productID, err := uuid.Parse(productIDStr); err == nil {
			feedbackFilters.ProductID = &productID
		}
	}
	
	return feedbackFilters
}

func (s *AnalyticsService) aggregateQuestionResponses(questionID uuid.UUID, questionText, questionType string, responses []feedbackmodel.Response) analyticsModels.ChartData {
	chartData := analyticsModels.ChartData{
		QuestionID:   questionID,
		QuestionText: questionText,
		QuestionType: questionType,
		Data:         make(map[string]interface{}),
	}
	
	switch questionType {
	case "rating":
		chartData.ChartType = "rating"
		chartData.Data = s.aggregateRatingResponses(responses)
	case "scale":
		chartData.ChartType = "scale"
		chartData.Data = s.aggregateScaleResponses(responses)
	case "single_choice", "multi_choice":
		chartData.ChartType = "choice"
		chartData.Data = s.aggregateChoiceResponses(responses, questionType == "multi_choice")
	case "yes_no":
		chartData.ChartType = "choice"  
		chartData.Data = s.aggregateYesNoResponses(responses)
	case "text":
		chartData.ChartType = "text_sentiment"
		chartData.Data = s.aggregateTextResponses(responses)
	default:
		chartData.ChartType = "rating"
		chartData.Data = s.aggregateRatingResponses(responses)
	}
	
	return chartData
}

func (s *AnalyticsService) aggregateRatingResponses(responses []feedbackmodel.Response) map[string]interface{} {
	distribution := make(map[string]int64)
	var total int64
	var sum float64
	maxScale := 0
	
	for _, response := range responses {
		switch v := response.Answer.(type) {
		case float64:
			key := fmt.Sprintf("%.0f", v)
			distribution[key]++
			sum += v
			total++
			if int(v) > maxScale {
				maxScale = int(v)
			}
		case int:
			key := fmt.Sprintf("%d", v)
			distribution[key]++
			sum += float64(v)
			total++
			if v > maxScale {
				maxScale = v
			}
		}
	}
	
	percentages := make(map[string]float64)
	for key, count := range distribution {
		if total > 0 {
			percentages[key] = float64(count) / float64(total) * 100
		}
	}
	
	scale := 5
	if maxScale > 5 {
		scale = 10
	}
	
	var average float64
	if total > 0 {
		average = sum / float64(total)
	}
	
	return map[string]interface{}{
		"scale":        scale,
		"distribution": distribution,
		"average":      average,
		"total":        total,
		"percentages":  percentages,
	}
}

func (s *AnalyticsService) aggregateScaleResponses(responses []feedbackmodel.Response) map[string]interface{} {
	distribution := make(map[string]int64)
	var total int64
	var sum float64
	maxScale := 0
	
	for _, response := range responses {
		switch v := response.Answer.(type) {
		case float64:
			key := fmt.Sprintf("%.0f", v)
			distribution[key]++
			sum += v
			total++
			if int(v) > maxScale {
				maxScale = int(v)
			}
		case int:
			key := fmt.Sprintf("%d", v)
			distribution[key]++
			sum += float64(v)
			total++
			if v > maxScale {
				maxScale = v
			}
		}
	}
	
	percentages := make(map[string]float64)
	for key, count := range distribution {
		if total > 0 {
			percentages[key] = float64(count) / float64(total) * 100
		}
	}
	
	scale := 5
	if maxScale > 5 {
		scale = 10
	}
	
	var average float64
	if total > 0 {
		average = sum / float64(total)
	}
	
	return map[string]interface{}{
		"scale":        scale,
		"distribution": distribution,
		"average":      average,
		"total":        total,
		"percentages":  percentages,
		"is_scale":     true,
	}
}

func (s *AnalyticsService) aggregateChoiceResponses(responses []feedbackmodel.Response, isMultiChoice bool) map[string]interface{} {
	options := make(map[string]int64)
	combinations := make(map[string]int64)
	var total int64
	
	for _, response := range responses {
		total++
		
		switch v := response.Answer.(type) {
		case string:
			if isMultiChoice {
				choices := strings.Split(v, ",")
				for i, choice := range choices {
					choices[i] = strings.TrimSpace(choice)
					options[choices[i]]++
				}
				sort.Strings(choices)
				combKey := strings.Join(choices, " + ")
				combinations[combKey]++
			} else {
				options[v]++
			}
		case []interface{}:
			if isMultiChoice {
				var choices []string
				for _, item := range v {
					if str, ok := item.(string); ok {
						choices = append(choices, str)
						options[str]++
					}
				}
				sort.Strings(choices)
				combKey := strings.Join(choices, " + ")
				combinations[combKey]++
			}
		}
	}
	
	percentages := make(map[string]float64)
	for option, count := range options {
		if total > 0 {
			percentages[option] = float64(count) / float64(total) * 100
		}
	}
	
	result := map[string]interface{}{
		"options":        options,
		"total":          total,
		"percentages":    percentages,
		"is_multi_choice": isMultiChoice,
	}
	
	if isMultiChoice && len(combinations) > 0 {
		var combList []map[string]interface{}
		for combo, count := range combinations {
			combList = append(combList, map[string]interface{}{
				"options":    strings.Split(combo, " + "),
				"count":      count,
				"percentage": float64(count) / float64(total) * 100,
			})
		}
		result["combinations"] = combList
	}
	
	return result
}

func (s *AnalyticsService) GetProductAnalyticsBatch(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]analyticsModels.ProductAnalytics, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID]analyticsModels.ProductAnalytics), nil
	}

	metricsMap, err := s.analyticsRepo.GetProductRatingsAndCounts(ctx, organizationID, productIDs)
	if err != nil {
		return nil, err
	}

	products, err := s.productRepo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	productNameMap := make(map[uuid.UUID]string)
	for _, product := range products {
		productNameMap[product.ID] = product.Name
	}

	result := make(map[uuid.UUID]analyticsModels.ProductAnalytics)
	for _, productID := range productIDs {
		metrics := metricsMap[productID]
		productName := productNameMap[productID]
		
		result[productID] = analyticsModels.ProductAnalytics{
			ProductID:     productID,
			ProductName:   productName,
			AverageRating: metrics.AverageRating,
			TotalFeedback: metrics.FeedbackCount,
		}
	}

	return result, nil
}

func (s *AnalyticsService) aggregateYesNoResponses(responses []feedbackmodel.Response) map[string]interface{} {
	options := make(map[string]int64)
	var total int64
	
	for _, response := range responses {
		total++
		switch v := response.Answer.(type) {
		case bool:
			if v {
				options["Yes"]++
			} else {
				options["No"]++
			}
		case string:
			lower := strings.ToLower(v)
			if lower == "yes" || lower == "true" || lower == "1" {
				options["Yes"]++
			} else {
				options["No"]++
			}
		}
	}
	
	percentages := make(map[string]float64)
	for option, count := range options {
		if total > 0 {
			percentages[option] = float64(count) / float64(total) * 100
		}
	}
	
	return map[string]interface{}{
		"options":        options,
		"total":          total,
		"percentages":    percentages,
		"is_multi_choice": false,
	}
}

func (s *AnalyticsService) analyzeSentiment(text string) float64 {
	if strings.TrimSpace(text) == "" {
		return 0.0
	}

	analyzer := sentitext.Parse(text, lexicon.DefaultLexicon)
	sentiment := sentitext.PolarityScore(analyzer)

	return sentiment.Compound
}

func (s *AnalyticsService) aggregateTextResponses(responses []feedbackmodel.Response) map[string]interface{} {
	var positive, neutral, negative int64
	var samples []string
	var totalSentiment float64
	validCount := 0
	
	wordFreq := make(map[string]int)
	
	for _, response := range responses {
		if text, ok := response.Answer.(string); ok && strings.TrimSpace(text) != "" {
			sentiment := s.analyzeSentiment(text)
			totalSentiment += sentiment
			validCount++
			
			if sentiment > 0.05 {
				positive++
			} else if sentiment < -0.05 {
				negative++
			} else {
				neutral++
			}
			
			if len(samples) < 5 {
				samples = append(samples, text)
			}
			
			words := strings.Fields(strings.ToLower(text))
			for _, word := range words {
				if len(word) > 3 && !isCommonWord(word) {
					wordFreq[word]++
				}
			}
		}
	}
	
	type wordCount struct {
		word  string
		count int
	}
	var sortedWords []wordCount
	for word, count := range wordFreq {
		sortedWords = append(sortedWords, wordCount{word, count})
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		return sortedWords[i].count > sortedWords[j].count
	})
	
	var topKeywords []string
	for i := 0; i < len(sortedWords) && i < 5; i++ {
		topKeywords = append(topKeywords, sortedWords[i].word)
	}
	
	total := positive + neutral + negative
	averageSentiment := float64(0)
	if validCount > 0 {
		averageSentiment = totalSentiment / float64(validCount)
	}
	
	return map[string]interface{}{
		"positive":  positive,
		"neutral":   neutral,
		"negative":  negative,
		"total":     total,
		"samples":   samples,
		"keywords":  topKeywords,
		"average_sentiment": averageSentiment,
	}
}

func isCommonWord(word string) bool {
	commonWords := map[string]bool{
		"the": true, "and": true, "is": true, "it": true, "to": true,
		"of": true, "in": true, "for": true, "on": true, "with": true,
		"as": true, "was": true, "that": true, "are": true, "been": true,
		"have": true, "had": true, "were": true, "said": true, "each": true,
		"which": true, "she": true, "will": true, "their": true, "what": true,
		"there": true, "would": true, "this": true, "from": true, "they": true,
		"more": true, "when": true, "other": true, "into": true,
	}
	return commonWords[word]
}
