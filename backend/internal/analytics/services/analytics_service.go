package services

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	analyticsModels "lecritique/internal/analytics/models"
	feedbackModels "lecritique/internal/feedback/models"
	feedbackRepos "lecritique/internal/feedback/repositories"
	menuRepos "lecritique/internal/menu/repositories"
	qrcodeRepos "lecritique/internal/qrcode/repositories"
	restaurantRepos "lecritique/internal/restaurant/repositories"
	"lecritique/internal/shared/logger"
	sharedModels "lecritique/internal/shared/models"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type AnalyticsService interface {
	GetDashboardMetrics(ctx context.Context, restaurantID uuid.UUID) (*analyticsModels.DashboardMetrics, error)
	GetDishInsights(ctx context.Context, dishID uuid.UUID) (*analyticsModels.DishInsights, error)
	GetRestaurantInsights(ctx context.Context, restaurantID uuid.UUID, period string) (*analyticsModels.RestaurantInsights, error)
	
	// Chart aggregation methods
	GetRestaurantChartData(ctx context.Context, restaurantID uuid.UUID, filters map[string]interface{}) (*analyticsModels.RestaurantChartData, error)
	GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*analyticsModels.ChartData, error)
}

type analyticsService struct {
	feedbackRepo   feedbackRepos.FeedbackRepository
	dishRepo       menuRepos.DishRepository
	qrCodeRepo     qrcodeRepos.QRCodeRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewAnalyticsService(i *do.Injector) (AnalyticsService, error) {
	return &analyticsService{
		feedbackRepo:   do.MustInvoke[feedbackRepos.FeedbackRepository](i),
		dishRepo:       do.MustInvoke[menuRepos.DishRepository](i),
		qrCodeRepo:     do.MustInvoke[qrcodeRepos.QRCodeRepository](i),
		restaurantRepo: do.MustInvoke[restaurantRepos.RestaurantRepository](i),
	}, nil
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
			LastScan:       qr.LastScannedAt,
			IsActive:       qr.IsActive,
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

// Chart aggregation implementations

func (s *analyticsService) GetRestaurantChartData(ctx context.Context, restaurantID uuid.UUID, filters map[string]interface{}) (*analyticsModels.RestaurantChartData, error) {
	// Build filter struct from map
	feedbackFilters := s.buildFeedbackFilters(filters)
	
	// Get filtered feedback data  
	feedback, err := s.feedbackRepo.FindByRestaurantIDWithFilters(ctx, restaurantID, sharedModels.PageRequest{
		Page: 1, Limit: 10000, // Large limit to get all relevant data
	}, feedbackFilters)
	if err != nil {
		return nil, err
	}
	
	// Get all dishes for this restaurant to map dish names
	dishes, err := s.dishRepo.FindByRestaurantID(ctx, restaurantID)
	if err != nil {
		logger.Error("Failed to get dishes", err, logrus.Fields{"restaurant_id": restaurantID})
	}
	
	// Create dish ID to name map
	dishMap := make(map[uuid.UUID]string)
	if dishes != nil {
		for _, dish := range dishes {
			dishMap[dish.ID] = dish.Name
		}
	}
	
	// Aggregate chart data by question across all dishes
	chartData := &analyticsModels.RestaurantChartData{
		RestaurantID: restaurantID,
		Charts:       []analyticsModels.ChartData{},
	}
	
	// Group responses by question_id and dish_id combination
	type questionDishKey struct {
		QuestionID uuid.UUID
		DishID     uuid.UUID
	}
	
	questionDishResponses := make(map[questionDishKey][]feedbackModels.Response)
	questionDishMeta := make(map[questionDishKey]struct {
		Text     string
		Type     string
		DishID   uuid.UUID
		DishName string
	})
	
	for _, f := range feedback.Data {
		dishName := dishMap[f.DishID]
		for _, response := range f.Responses {
			key := questionDishKey{
				QuestionID: response.QuestionID,
				DishID:     f.DishID,
			}
			questionDishResponses[key] = append(questionDishResponses[key], response)
			if questionDishMeta[key].Text == "" {
				questionDishMeta[key] = struct {
					Text     string
					Type     string
					DishID   uuid.UUID
					DishName string
				}{
					Text:     response.QuestionText,
					Type:     string(response.QuestionType),
					DishID:   f.DishID,
					DishName: dishName,
				}
			}
		}
	}
	
	// Generate chart data for each question-dish combination
	for key, responses := range questionDishResponses {
		meta := questionDishMeta[key]
		chart := s.aggregateQuestionResponses(key.QuestionID, meta.Text, meta.Type, responses)
		// Add dish information
		if meta.DishID != uuid.Nil {
			chart.DishID = &meta.DishID
			chart.DishName = meta.DishName
		}
		chartData.Charts = append(chartData.Charts, chart)
	}
	
	// Set summary information
	chartData.Summary.TotalResponses = int64(len(feedback.Data))
	chartData.Summary.FiltersApplied = filters
	
	// Set date range based on data
	if len(feedback.Data) > 0 {
		chartData.Summary.DateRange.Start = feedback.Data[len(feedback.Data)-1].CreatedAt
		chartData.Summary.DateRange.End = feedback.Data[0].CreatedAt
	}
	
	return chartData, nil
}

func (s *analyticsService) GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*analyticsModels.ChartData, error) {
	// Build filter struct
	feedbackFilters := s.buildFeedbackFilters(filters)
	
	// Get question details first to understand dish context
	// Note: We'd need to add a method to get question by ID, for now we'll work with what we have
	
	// For now, get all feedback that contains this question
	// This is not optimal - ideally we'd have a direct question-to-feedback lookup
	allFeedback, err := s.feedbackRepo.FindByRestaurantIDWithFilters(ctx, uuid.UUID{}, sharedModels.PageRequest{
		Page: 1, Limit: 10000,
	}, feedbackFilters)
	if err != nil {
		return nil, err
	}
	
	// Filter responses for this specific question
	var responses []feedbackModels.Response
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

// Helper methods

func (s *analyticsService) buildFeedbackFilters(filters map[string]interface{}) feedbackRepos.FeedbackFilter {
	feedbackFilters := feedbackRepos.FeedbackFilter{}
	
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
	
	if dishIDStr, ok := filters["dish_id"].(string); ok {
		if dishID, err := uuid.Parse(dishIDStr); err == nil {
			feedbackFilters.DishID = &dishID
		}
	}
	
	return feedbackFilters
}

func (s *analyticsService) aggregateQuestionResponses(questionID uuid.UUID, questionText, questionType string, responses []feedbackModels.Response) analyticsModels.ChartData {
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
		chartData.ChartType = "rating" // fallback
		chartData.Data = s.aggregateRatingResponses(responses)
	}
	
	return chartData
}

func (s *analyticsService) aggregateRatingResponses(responses []feedbackModels.Response) map[string]interface{} {
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
	
	// Calculate percentages
	percentages := make(map[string]float64)
	for key, count := range distribution {
		if total > 0 {
			percentages[key] = float64(count) / float64(total) * 100
		}
	}
	
	// Determine scale (common scales: 5, 10)
	scale := 5 // default
	if maxScale > 5 {
		scale = 10
	}
	
	return map[string]interface{}{
		"scale":        scale,
		"distribution": distribution,
		"average":      sum / float64(total),
		"total":        total,
		"percentages":  percentages,
	}
}

func (s *analyticsService) aggregateScaleResponses(responses []feedbackModels.Response) map[string]interface{} {
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
	
	// Calculate percentages
	percentages := make(map[string]float64)
	for key, count := range distribution {
		if total > 0 {
			percentages[key] = float64(count) / float64(total) * 100
		}
	}
	
	// Determine scale (common scales: 5, 10)
	scale := 5 // default
	if maxScale > 5 {
		scale = 10
	}
	
	return map[string]interface{}{
		"scale":        scale,
		"distribution": distribution,
		"average":      sum / float64(total),
		"total":        total,
		"percentages":  percentages,
		"is_scale":     true, // Flag to differentiate from rating
	}
}

func (s *analyticsService) aggregateChoiceResponses(responses []feedbackModels.Response, isMultiChoice bool) map[string]interface{} {
	options := make(map[string]int64)
	combinations := make(map[string]int64) // for multi-choice
	var total int64
	
	for _, response := range responses {
		total++
		
		switch v := response.Answer.(type) {
		case string:
			if isMultiChoice {
				// Handle multi-choice - could be comma-separated or JSON array
				choices := strings.Split(v, ",")
				for i, choice := range choices {
					choices[i] = strings.TrimSpace(choice)
					options[choices[i]]++
				}
				// Track combination
				sort.Strings(choices)
				combKey := strings.Join(choices, " + ")
				combinations[combKey]++
			} else {
				options[v]++
			}
		case []interface{}:
			// JSON array format for multi-choice
			if isMultiChoice {
				var choices []string
				for _, item := range v {
					if str, ok := item.(string); ok {
						choices = append(choices, str)
						options[str]++
					}
				}
				// Track combination
				sort.Strings(choices)
				combKey := strings.Join(choices, " + ")
				combinations[combKey]++
			}
		}
	}
	
	// Calculate percentages
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
	
	// Add combinations for multi-choice
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

func (s *analyticsService) aggregateYesNoResponses(responses []feedbackModels.Response) map[string]interface{} {
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
			// Handle string representations
			lower := strings.ToLower(v)
			if lower == "yes" || lower == "true" || lower == "1" {
				options["Yes"]++
			} else {
				options["No"]++
			}
		}
	}
	
	// Calculate percentages
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

func (s *analyticsService) aggregateTextResponses(responses []feedbackModels.Response) map[string]interface{} {
	var positive, neutral, negative int64
	var samples []string
	keywords := make(map[string]int)
	
	// Simple keyword-based sentiment analysis
	positiveWords := []string{"excellent", "great", "amazing", "love", "perfect", "fantastic", "wonderful"}
	negativeWords := []string{"terrible", "awful", "bad", "hate", "horrible", "disgusting", "worst"}
	
	for _, response := range responses {
		if text, ok := response.Answer.(string); ok && text != "" {
			lower := strings.ToLower(text)
			
			// Simple sentiment detection
			hasPositive := false
			hasNegative := false
			
			for _, word := range positiveWords {
				if strings.Contains(lower, word) {
					hasPositive = true
					keywords[word]++
					break
				}
			}
			
			for _, word := range negativeWords {
				if strings.Contains(lower, word) {
					hasNegative = true
					keywords[word]++
					break
				}
			}
			
			// Classify sentiment
			if hasPositive && !hasNegative {
				positive++
			} else if hasNegative && !hasPositive {
				negative++
			} else {
				neutral++
			}
			
			// Add to samples (limit to 5)
			if len(samples) < 5 {
				samples = append(samples, text)
			}
		}
	}
	
	// Get top keywords
	var topKeywords []string
	for word := range keywords {
		topKeywords = append(topKeywords, word)
		if len(topKeywords) >= 5 {
			break
		}
	}
	
	total := positive + neutral + negative
	
	return map[string]interface{}{
		"positive":  positive,
		"neutral":   neutral,
		"negative":  negative,
		"total":     total,
		"samples":   samples,
		"keywords":  topKeywords,
	}
}