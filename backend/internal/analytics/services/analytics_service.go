package services

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	analyticsModels "github.com/lecritique/api/internal/analytics/models"
	feedbackModels "github.com/lecritique/api/internal/feedback/models"
	feedbackRepos "github.com/lecritique/api/internal/feedback/repositories"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/sirupsen/logrus"
)

type AnalyticsService interface {
	GetDashboardMetrics(ctx context.Context, restaurantID uuid.UUID) (*analyticsModels.DashboardMetrics, error)
	GetDishInsights(ctx context.Context, dishID uuid.UUID) (*analyticsModels.DishInsights, error)
	GetRestaurantInsights(ctx context.Context, restaurantID uuid.UUID, period string) (*analyticsModels.RestaurantInsights, error)
}

type analyticsService struct {
	feedbackRepo   feedbackRepos.FeedbackRepository
	dishRepo       menuRepos.DishRepository
	qrCodeRepo     qrcodeRepos.QRCodeRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewAnalyticsService(
	feedbackRepo feedbackRepos.FeedbackRepository,
	dishRepo menuRepos.DishRepository,
	qrCodeRepo qrcodeRepos.QRCodeRepository,
	restaurantRepo restaurantRepos.RestaurantRepository,
) AnalyticsService {
	return &analyticsService{
		feedbackRepo:   feedbackRepo,
		dishRepo:       dishRepo,
		qrCodeRepo:     qrCodeRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *analyticsService) GetDashboardMetrics(ctx context.Context, restaurantID uuid.UUID) (*analyticsModels.DashboardMetrics, error) {
	metrics := &analyticsModels.DashboardMetrics{}
	
	// Get total feedback count (operational metric)
	totalCount, err := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, time.Time{})
	if err != nil {
		logger.Error("Failed to get total feedback count", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	}
	metrics.TotalFeedbacks = totalCount
	
	// Get today's feedback count (operational metric)
	todayStart := time.Now().Truncate(24 * time.Hour)
	todayCount, err := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, todayStart)
	if err != nil {
		logger.Error("Failed to get today's feedback count", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	}
	metrics.TodaysFeedback = todayCount
	
	// Get yesterday's count for trend (operational metric)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	yesterdayCount, _ := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, yesterdayStart)
	if yesterdayCount > 0 {
		metrics.TrendVsYesterday = float64(todayCount-yesterdayCount) / float64(yesterdayCount) * 100
	}
	
	// Get QR code metrics (operational)
	qrMetrics, err := s.getQRCodeMetrics(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get QR code metrics", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	} else {
		metrics.ActiveQRCodes = qrMetrics.ActiveCount
		metrics.TotalQRScans = qrMetrics.TotalScans
		metrics.ScansToday = qrMetrics.ScansToday
	}
	
	// Calculate completion rate using recent data (last 30 days)
	// This avoids issues with historical data from before scan tracking
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	recentFeedbackCount, _ := s.feedbackRepo.CountByRestaurantID(ctx, restaurantID, thirtyDaysAgo)
	
	if metrics.TotalQRScans > 0 && recentFeedbackCount > 0 {
		// Use recent feedback vs total scans for a more realistic rate
		rate := float64(recentFeedbackCount) / float64(metrics.TotalQRScans) * 100
		if rate > 100 {
			// Still cap at 100% as safety measure
			metrics.CompletionRate = 100.0
		} else {
			metrics.CompletionRate = rate
		}
	} else {
		// Show 0% if insufficient data
		metrics.CompletionRate = 0.0
	}
	
	// Get device/platform analytics (operational)
	deviceMetrics, err := s.getDeviceMetrics(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get device metrics", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	} else {
		metrics.DeviceBreakdown = deviceMetrics
	}
	
	// Get response time analytics (operational)
	responseTime, err := s.getAverageResponseTime(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get response time", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	} else {
		metrics.AverageResponseTime = responseTime
	}
	
	// Get peak usage hours (operational)
	peakHours, err := s.getPeakUsageHours(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get peak hours", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	} else {
		metrics.PeakHours = peakHours
	}
	
	// Get QR code performance analytics
	qrPerformance, err := s.getQRCodePerformance(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get QR performance", err, logrus.Fields{
			"restaurant_id": restaurantID,
		})
	} else {
		metrics.QRPerformance = qrPerformance
	}
	
	return metrics, nil
}

func (s *analyticsService) GetDishInsights(ctx context.Context, dishID uuid.UUID) (*analyticsModels.DishInsights, error) {
	// Get dish info
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return nil, err
	}
	
	// Get all feedback for this dish
	allFeedback, err := s.feedbackRepo.FindByDishIDForAnalytics(ctx, dishID, 500)
	if err != nil {
		return nil, err
	}
	
	insights := &analyticsModels.DishInsights{
		DishID:        dishID,
		DishName:      dish.Name,
		TotalFeedback: int64(len(allFeedback)),
	}
	
	if len(allFeedback) == 0 {
		return insights, nil
	}
	
	// Get questionnaire to understand questions
	questions, err := s.feedbackRepo.GetQuestionsByDishID(ctx, dishID)
	if err != nil {
		logger.Error("Failed to get questions", err, logrus.Fields{"dish_id": dishID})
	}
	
	// Aggregate metrics by question
	questionMetrics := s.aggregateQuestionMetrics(allFeedback, questions)
	insights.Questions = questionMetrics
	
	// Calculate overall score from numeric questions
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
	
	// Identify best aspects and needs attention
	insights.BestAspects = s.identifyBestAspects(questionMetrics)
	insights.NeedsAttention = s.identifyNeedsAttention(questionMetrics)
	
	// Calculate completion rate
	// TODO: Track started vs completed questionnaires
	insights.CompletionRate = 100.0 // Placeholder
	
	// Get last feedback time
	if len(allFeedback) > 0 {
		insights.LastFeedback = allFeedback[0].CreatedAt
	}
	
	// Calculate weekly change
	weekAgo := time.Now().AddDate(0, 0, -7)
	weekAgoFeedback := filterFeedbackByDate(allFeedback, weekAgo)
	if len(weekAgoFeedback) > 0 {
		weekScore, _, _ := s.calculateOverallMetrics(weekAgoFeedback)
		insights.WeeklyChange = ((insights.OverallScore - weekScore) / weekScore) * 100
	}
	
	// Determine trend
	if insights.WeeklyChange > 5 {
		insights.ScoreTrend = "improving"
	} else if insights.WeeklyChange < -5 {
		insights.ScoreTrend = "declining"
	} else {
		insights.ScoreTrend = "stable"
	}
	
	return insights, nil
}

func (s *analyticsService) GetRestaurantInsights(ctx context.Context, restaurantID uuid.UUID, period string) (*analyticsModels.RestaurantInsights, error) {
	// Implementation would be similar to dashboard metrics but more comprehensive
	// This is a placeholder for now
	return &analyticsModels.RestaurantInsights{
		RestaurantID:   restaurantID,
		Period:         period,
		TotalFeedback:  0,
		ActiveDishes:   0,
	}, nil
}

// Helper methods

func (s *analyticsService) calculateOverallMetrics(feedback []feedbackModels.Feedback) (satisfaction, recommendRate, sentiment float64) {
	if len(feedback) == 0 {
		return 0, 0, 0
	}
	
	var totalScore float64
	var scoreCount int
	var recommendYes int
	var recommendTotal int
	var positiveCount int
	
	for _, f := range feedback {
		// Use the overall_rating field for satisfaction calculation
		if f.OverallRating > 0 {
			totalScore += float64(f.OverallRating)
			scoreCount++
			if f.OverallRating >= 4 {
				positiveCount++
			}
		}
		
		// Still process responses for recommendation rate and other metrics
		for _, response := range f.Responses {
			// Handle different response types
			switch v := response.Answer.(type) {
			case bool:
				// Yes/No responses - need to check question text from questionnaire
				recommendTotal++
				if v {
					recommendYes++
				}
			case string:
				// Text or choice responses - would need sentiment analysis
				// For now, use simple keyword matching
				lower := strings.ToLower(v)
				if strings.Contains(lower, "excellent") || strings.Contains(lower, "great") || 
				   strings.Contains(lower, "amazing") || strings.Contains(lower, "love") {
					positiveCount++
				}
			}
		}
	}
	
	if scoreCount > 0 {
		satisfaction = totalScore / float64(scoreCount) // Keep as 1-5 scale for consistency
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

func (s *analyticsService) aggregateQuestionMetrics(feedback []feedbackModels.Feedback, questions []feedbackModels.Question) []analyticsModels.QuestionMetric {
	metricsMap := make(map[uuid.UUID]*analyticsModels.QuestionMetric)
	
	// Initialize metrics for each question
	for _, q := range questions {
		metric := &analyticsModels.QuestionMetric{
			QuestionID:         q.ID,
			QuestionText:       q.Text,
			QuestionType:       string(q.Type),
			OptionDistribution: make(map[string]int64),
		}
		metricsMap[q.ID] = metric
	}
	
	// Aggregate responses
	for _, f := range feedback {
		for _, response := range f.Responses {
			metric, exists := metricsMap[response.QuestionID]
			if !exists {
				continue
			}
			
			metric.ResponseCount++
			
			// Process based on answer type
			switch v := response.Answer.(type) {
			case float64:
				// Numeric response (rating, scale)
				if metric.AverageScore == nil {
					avg := v
					min := v
					max := v
					metric.AverageScore = &avg
					metric.MinScore = &min
					metric.MaxScore = &max
				} else {
					// Update average
					*metric.AverageScore = (*metric.AverageScore * float64(metric.ResponseCount-1) + v) / float64(metric.ResponseCount)
					// Update min/max
					if v < *metric.MinScore {
						*metric.MinScore = v
					}
					if v > *metric.MaxScore {
						*metric.MaxScore = v
					}
				}
				
				// Calculate sentiment
				if v >= 4 {
					metric.PositiveRate++
				} else if v >= 3 {
					metric.NeutralRate++
				} else {
					metric.NegativeRate++
				}
				
			case string:
				// Text or choice response
				if metric.QuestionType == string(feedbackModels.QuestionTypeText) {
					metric.TextResponses = append(metric.TextResponses, v)
				} else {
					// Choice response
					metric.OptionDistribution[v]++
				}
				
			case bool:
				// Yes/No response
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
	
	// Convert sentiment counts to rates
	for _, metric := range metricsMap {
		if metric.ResponseCount > 0 {
			total := float64(metric.ResponseCount)
			metric.PositiveRate = (metric.PositiveRate / total) * 100
			metric.NeutralRate = (metric.NeutralRate / total) * 100
			metric.NegativeRate = (metric.NegativeRate / total) * 100
		}
	}
	
	// Convert map to slice
	var metrics []analyticsModels.QuestionMetric
	for _, m := range metricsMap {
		metrics = append(metrics, *m)
	}
	
	// Sort by display order or response count
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].ResponseCount > metrics[j].ResponseCount
	})
	
	return metrics
}

func (s *analyticsService) identifyTopIssues(feedback []feedbackModels.Feedback) []analyticsModels.QuickIssue {
	issues := []analyticsModels.QuickIssue{}
	
	// Count low ratings
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
	
	// More sophisticated issue detection would go here
	// - Trending down dishes
	// - Repeated complaints
	// - Sudden drops in specific questions
	
	return issues
}

func (s *analyticsService) aggregateByDish(feedback []feedbackModels.Feedback) map[uuid.UUID]*analyticsModels.DishSummary {
	dishMap := make(map[uuid.UUID]*analyticsModels.DishSummary)
	
	for _, f := range feedback {
		summary, exists := dishMap[f.DishID]
		if !exists {
			summary = &analyticsModels.DishSummary{
				DishID:   f.DishID,
				DishName: f.Dish.Name,
			}
			dishMap[f.DishID] = summary
		}
		
		summary.FeedbackCount++
		summary.Score = (summary.Score*float64(summary.FeedbackCount-1) + float64(f.OverallRating)) / float64(summary.FeedbackCount)
	}
	
	return dishMap
}

func (s *analyticsService) getTopDishes(dishMap map[uuid.UUID]*analyticsModels.DishSummary, limit int) []analyticsModels.DishSummary {
	var dishes []analyticsModels.DishSummary
	for _, d := range dishMap {
		dishes = append(dishes, *d)
	}
	
	sort.Slice(dishes, func(i, j int) bool {
		return dishes[i].Score > dishes[j].Score
	})
	
	if len(dishes) > limit {
		dishes = dishes[:limit]
	}
	
	return dishes
}

func (s *analyticsService) getBottomDishes(dishMap map[uuid.UUID]*analyticsModels.DishSummary, limit int) []analyticsModels.DishSummary {
	var dishes []analyticsModels.DishSummary
	for _, d := range dishMap {
		if d.Score < 3.5 { // Only show dishes that need attention
			dishes = append(dishes, *d)
		}
	}
	
	sort.Slice(dishes, func(i, j int) bool {
		return dishes[i].Score < dishes[j].Score
	})
	
	if len(dishes) > limit {
		dishes = dishes[:limit]
	}
	
	return dishes
}

func (s *analyticsService) getRecentFeedbackSummaries(feedback []feedbackModels.Feedback, limit int) []analyticsModels.FeedbackSummary {
	// Sort by date
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
			DishName:     f.Dish.Name,
			CustomerName: f.CustomerName,
			Score:        float64(f.OverallRating),
			CreatedAt:    f.CreatedAt,
		}
		
		// Determine sentiment
		if f.OverallRating >= 4 {
			summary.Sentiment = "positive"
		} else if f.OverallRating >= 3 {
			summary.Sentiment = "neutral"
		} else {
			summary.Sentiment = "negative"
		}
		
		// Extract highlight from text responses
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

func (s *analyticsService) identifyBestAspects(metrics []analyticsModels.QuestionMetric) []string {
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

func (s *analyticsService) identifyNeedsAttention(metrics []analyticsModels.QuestionMetric) []string {
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

func filterFeedbackByDate(feedback []feedbackModels.Feedback, since time.Time) []feedbackModels.Feedback {
	var filtered []feedbackModels.Feedback
	for _, f := range feedback {
		if f.CreatedAt.After(since) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// New operational metrics methods

type QRMetrics struct {
	ActiveCount int64
	TotalScans  int64
	ScansToday  int64
}

func (s *analyticsService) getQRCodeMetrics(ctx context.Context, restaurantID uuid.UUID) (*QRMetrics, error) {
	// Get all QR codes for restaurant
	qrCodes, err := s.qrCodeRepo.FindByRestaurantID(ctx, restaurantID)
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
		
		// Count today's scans (approximation based on last_scanned_at)
		if qr.LastScannedAt != nil && qr.LastScannedAt.After(todayStart) {
			// This is an approximation - ideally we'd have daily scan logs
			metrics.ScansToday++
		}
	}
	
	return metrics, nil
}

func (s *analyticsService) getDeviceMetrics(ctx context.Context, restaurantID uuid.UUID) (map[string]int64, error) {
	// Get feedback with device info
	feedback, err := s.feedbackRepo.FindByRestaurantIDForAnalytics(ctx, restaurantID, 1000)
	if err != nil {
		return nil, err
	}
	
	deviceBreakdown := make(map[string]int64)
	platformBreakdown := make(map[string]int64)
	
	for _, f := range feedback {
		// Count by platform
		if f.DeviceInfo.Platform != "" {
			platformBreakdown[f.DeviceInfo.Platform]++
		}
		
		// Count by browser
		if f.DeviceInfo.Browser != "" {
			deviceBreakdown[f.DeviceInfo.Browser]++
		}
	}
	
	// Combine results - you might want to separate these
	result := make(map[string]int64)
	for k, v := range platformBreakdown {
		result[k] = v
	}
	for k, v := range deviceBreakdown {
		result[k+" Browser"] = v
	}
	
	return result, nil
}

func (s *analyticsService) getAverageResponseTime(ctx context.Context, restaurantID uuid.UUID) (float64, error) {
	// Get feedback with QR code info
	feedback, err := s.feedbackRepo.FindByRestaurantIDForAnalytics(ctx, restaurantID, 500)
	if err != nil {
		return 0, err
	}
	
	var totalTime float64
	var count int
	
	for _, f := range feedback {
		// Get QR code last scanned time
		qrCode, err := s.qrCodeRepo.FindByID(ctx, f.QRCodeID)
		if err != nil || qrCode.LastScannedAt == nil {
			continue
		}
		
		// Calculate time between scan and submission
		responseTime := f.CreatedAt.Sub(*qrCode.LastScannedAt)
		if responseTime > 0 && responseTime < 24*time.Hour { // Reasonable bounds
			totalTime += responseTime.Minutes()
			count++
		}
	}
	
	if count == 0 {
		return 0, nil
	}
	
	return totalTime / float64(count), nil
}

func (s *analyticsService) getPeakUsageHours(ctx context.Context, restaurantID uuid.UUID) ([]int, error) {
	// Get feedback from last 7 days
	weekAgo := time.Now().AddDate(0, 0, -7)
	feedback, err := s.feedbackRepo.FindByRestaurantIDForAnalytics(ctx, restaurantID, 1000)
	if err != nil {
		return nil, err
	}
	
	// Filter to last week
	recentFeedback := filterFeedbackByDate(feedback, weekAgo)
	
	// Count by hour
	hourCounts := make(map[int]int64)
	for _, f := range recentFeedback {
		hour := f.CreatedAt.Hour()
		hourCounts[hour]++
	}
	
	// Find top 3 hours
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
	
	// Return top 3 hours
	var peakHours []int
	for i, h := range hours {
		if i >= 3 {
			break
		}
		peakHours = append(peakHours, h.hour)
	}
	
	return peakHours, nil
}

func (s *analyticsService) getQRCodePerformance(ctx context.Context, restaurantID uuid.UUID) ([]analyticsModels.QRCodePerformance, error) {
	// Get restaurant name first
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	
	// Get all QR codes for restaurant
	qrCodes, err := s.qrCodeRepo.FindByRestaurantID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	
	var performance []analyticsModels.QRCodePerformance
	
	for _, qr := range qrCodes {
		// Count feedback for this QR code
		feedbackCount, err := s.feedbackRepo.CountByQRCodeID(ctx, qr.ID)
		if err != nil {
			logger.Error("Failed to count feedback for QR code", err, logrus.Fields{
				"qr_code_id": qr.ID,
			})
			feedbackCount = 0
		}
		
		// Calculate conversion rate with safety caps
		conversionRate := 0.0
		if qr.ScansCount > 0 && feedbackCount > 0 {
			rate := (float64(feedbackCount) / float64(qr.ScansCount)) * 100
			// Cap at 100% to avoid unrealistic rates due to historical data issues
			if rate > 100 {
				conversionRate = 100.0
			} else {
				conversionRate = rate
			}
		}
		
		perf := analyticsModels.QRCodePerformance{
			ID:             qr.ID,
			Label:          qr.Label,
			RestaurantID:   qr.RestaurantID,
			RestaurantName: restaurant.Name,
			ScansCount:     int64(qr.ScansCount),
			FeedbackCount:  feedbackCount,
			ConversionRate: conversionRate,
		}
		
		// Add location if available
		if qr.Location != nil {
			perf.Location = *qr.Location
		}
		
		performance = append(performance, perf)
	}
	
	// Sort by scan count (most used first)
	sort.Slice(performance, func(i, j int) bool {
		return performance[i].ScansCount > performance[j].ScansCount
	})
	
	return performance, nil
}