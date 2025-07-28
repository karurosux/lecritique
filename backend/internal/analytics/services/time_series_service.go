package services

import (
	"context"
	"encoding/json"
	"fmt"
	"kyooar/internal/analytics/models"
	"kyooar/internal/analytics/repositories"
	"kyooar/internal/shared/logger"
	"math"
	"sort"
	"strings"
	"time"

	feedbackModels "kyooar/internal/feedback/models"
	feedbackServices "kyooar/internal/feedback/services"
	organizationServices "kyooar/internal/organization/services"

	"github.com/google/uuid"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

// ResponseData represents processed response data for metrics calculation
type ResponseData struct {
	Responses     []any
	QuestionType  string
	QuestionTexts []string
	ProductID     string
	ProductName   string
}

// QuestionData represents individual question data for metrics
type QuestionData struct {
	Responses    []any
	QuestionID   string
	QuestionText string
	QuestionType string
	ProductID    string
	ProductName  string
}

type TimeSeriesService interface {
	CollectMetrics(ctx context.Context, organizationID uuid.UUID) error
	GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) (*models.TimeSeriesResponse, error)
	GetComparison(ctx context.Context, request models.ComparisonRequest) (*models.ComparisonResponse, error)
	CleanupOldMetrics(ctx context.Context, retentionDays int) error
}

type timeSeriesService struct {
	timeSeriesRepo      repositories.TimeSeriesRepository
	feedbackService     feedbackServices.FeedbackService
	organizationService organizationServices.OrganizationService
	analyticsService    AnalyticsService
	questionService     feedbackServices.QuestionService
}

func NewTimeSeriesService(i *do.Injector) (TimeSeriesService, error) {
	return &timeSeriesService{
		timeSeriesRepo:      do.MustInvoke[repositories.TimeSeriesRepository](i),
		feedbackService:     do.MustInvoke[feedbackServices.FeedbackService](i),
		organizationService: do.MustInvoke[organizationServices.OrganizationService](i),
		analyticsService:    do.MustInvoke[AnalyticsService](i),
		questionService:     do.MustInvoke[feedbackServices.QuestionService](i),
	}, nil
}

func (s *timeSeriesService) CollectMetrics(ctx context.Context, organizationID uuid.UUID) error {
	now := time.Now()

	organization, err := s.organizationService.GetByIDForAnalytics(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}

	accountID := organization.AccountID

	orgMetrics, err := s.collectOrganizationMetrics(ctx, organizationID, accountID, now)
	if err != nil {
		logger.Error("Failed to collect organization metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return err
	}

	if err := s.timeSeriesRepo.DeleteOrganizationMetrics(ctx, organizationID); err != nil {
		logger.Error("Failed to cleanup organization metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
	}

	if err := s.timeSeriesRepo.CreateBatch(ctx, orgMetrics); err != nil {
		return fmt.Errorf("failed to save organization metrics: %w", err)
	}

	// Collect product-level metrics
	// This would iterate through all products in the organization
	// For now, we'll skip the implementation details
	return nil
}

func (s *timeSeriesService) collectOrganizationMetrics(ctx context.Context, organizationID, accountID uuid.UUID, timestamp time.Time) ([]models.TimeSeriesMetric, error) {
	var metrics []models.TimeSeriesMetric

	feedbacks, err := s.feedbackService.GetByOrganizationIDForAnalytics(ctx, organizationID, 1000)
	if err != nil {
		return nil, err
	}

	questionMap := make(map[uuid.UUID]*feedbackModels.Question)
	processedProducts := make(map[uuid.UUID]bool)
	for _, feedback := range feedbacks {
		if feedback.ProductID != uuid.Nil && !processedProducts[feedback.ProductID] {
			questions, err := s.questionService.GetQuestionsByProductIDForAnalytics(ctx, feedback.ProductID)
			if err == nil {
				for _, question := range questions {
					questionMap[question.ID] = question
				}
			}
			processedProducts[feedback.ProductID] = true
		}
	}

	timeSeriesData, individualQuestionData, surveyResponsesByDate := s.processFeedbackData(feedbacks)

	for key, questionData := range timeSeriesData {
		parts := strings.Split(key, "_")
		if len(parts) < 3 {
			continue
		}

		dateKey := parts[0]
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		if len(questionData.Responses) == 0 {
			continue
		}

		// Process question using centralized logic
		value := s.processQuestionByType(questionData.QuestionType, questionData.Responses)
		count := len(questionData.Responses)

		// Create question type metric
		metricType := questionData.QuestionType + "_questions"
		metricName := s.getMetricNameWithProduct(questionData.QuestionType, questionData.ProductName, questionData.QuestionTexts)

		var productID *uuid.UUID
		if pid, err := uuid.Parse(questionData.ProductID); err == nil {
			productID = &pid
		}

		metrics = append(metrics, models.TimeSeriesMetric{
			AccountID:      accountID,
			OrganizationID: organizationID,
			ProductID:      productID,
			MetricType:     metricType,
			MetricName:     metricName,
			Value:          value,
			Count:          int64(count),
			Timestamp:      date,
			Granularity:    models.GranularityDaily,
			Metadata:       s.createMetadataWithProduct(questionData.QuestionType, questionData.ProductName),
		})
	}

	// Process individual question metrics
	for key, qData := range individualQuestionData {
		parts := strings.Split(key, "_")
		if len(parts) < 2 {
			continue
		}

		dateKey := parts[0]
		questionID := strings.Join(parts[1:], "_") // Handle UUIDs with underscores

		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		if len(qData.Responses) == 0 {
			continue
		}

		questionCount := len(qData.Responses)

		// Handle single choice questions specially to create choice-specific metrics
		if qData.QuestionType == string(feedbackModels.QuestionTypeSingleChoice) {
			// Create choice-specific metrics
			choiceMetrics := s.processSingleChoiceQuestion(
				qData.Responses, qData, questionID, accountID, organizationID, date, questionMap,
			)
			metrics = append(metrics, choiceMetrics...)
			// Skip creating summary metric for single choice questions
			continue
		}

		// Process question using centralized logic for summary metric
		questionValue := s.processQuestionByType(qData.QuestionType, qData.Responses)

		var productID *uuid.UUID
		if pid, err := uuid.Parse(qData.ProductID); err == nil {
			productID = &pid
		}

		var questionUUID *uuid.UUID
		if qid, err := uuid.Parse(questionID); err == nil {
			questionUUID = &qid
		}

		// Keep the question ID in metric type for frontend compatibility
		questionMetricType := "question_" + questionID
		// Include product name with question text for clarity
		questionMetricName := fmt.Sprintf("%s - %s", qData.ProductName, qData.QuestionText)

		// Get the question from our map
		var question *feedbackModels.Question
		if questionUUID != nil {
			question = questionMap[*questionUUID]
		}

		metrics = append(metrics, models.TimeSeriesMetric{
			AccountID:      accountID,
			OrganizationID: organizationID,
			ProductID:      productID,
			QuestionID:     questionUUID,
			MetricType:     questionMetricType,
			MetricName:     questionMetricName,
			Value:          questionValue,
			Count:          int64(questionCount),
			Timestamp:      date,
			Granularity:    models.GranularityDaily,
			Metadata:       s.createMetadataWithQuestion(qData.QuestionType, qData.ProductName, qData.QuestionText, question),
		})
	}

	// Create survey response count metrics for each date
	for dateKey, responseCount := range surveyResponsesByDate {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		metrics = append(metrics, models.TimeSeriesMetric{
			AccountID:      accountID,
			OrganizationID: organizationID,
			MetricType:     "survey_responses",
			MetricName:     "Total Survey Responses",
			Value:          float64(responseCount),
			Count:          int64(responseCount),
			Timestamp:      date,
			Granularity:    models.GranularityDaily,
		})
	}

	return metrics, nil
}

func (s *timeSeriesService) createMetadata(questionType string) *string {
	metadata := fmt.Sprintf(`{"question_type": "%s"}`, questionType)
	return &metadata
}

func (s *timeSeriesService) createMetadataWithProduct(questionType string, productName string) *string {
	// Only include question type, not product info
	metadata := fmt.Sprintf(`{"question_type": "%s"}`, questionType)
	return &metadata
}

func (s *timeSeriesService) createMetadataWithQuestion(questionType string, productName string, questionText string, question *feedbackModels.Question) *string {
	// Create metadata with question type, text, and scale labels if available
	metadataMap := map[string]interface{}{
		"question_type": questionType,
		"question_text": questionText,
	}

	// Add scale labels if available
	if question != nil {
		if question.MinLabel != "" {
			metadataMap["min_label"] = question.MinLabel
		}
		if question.MaxLabel != "" {
			metadataMap["max_label"] = question.MaxLabel
		}
		if question.MinValue != nil {
			metadataMap["min_value"] = *question.MinValue
		}
		if question.MaxValue != nil {
			metadataMap["max_value"] = *question.MaxValue
		}
	}

	// Convert to JSON string
	metadataBytes, err := json.Marshal(metadataMap)
	if err != nil {
		// Fallback to simple format
		metadata := fmt.Sprintf(`{"question_type": "%s", "question_text": "%s"}`, questionType, questionText)
		return &metadata
	}

	metadataStr := string(metadataBytes)
	return &metadataStr
}

func (s *timeSeriesService) createMetadataWithChoice(questionType string, productName string, questionText string, choiceOption string, question *feedbackModels.Question) *string {
	// Create metadata with question type, text, and choice option information
	metadataMap := map[string]interface{}{
		"question_type": questionType,
		"question_text": questionText,
		"choice_option": choiceOption,
	}

	// Add scale labels if available (though not typically used for single choice)
	if question != nil {
		if question.MinLabel != "" {
			metadataMap["min_label"] = question.MinLabel
		}
		if question.MaxLabel != "" {
			metadataMap["max_label"] = question.MaxLabel
		}
		if question.MinValue != nil {
			metadataMap["min_value"] = *question.MinValue
		}
		if question.MaxValue != nil {
			metadataMap["max_value"] = *question.MaxValue
		}
	}

	// Convert to JSON string
	metadataBytes, err := json.Marshal(metadataMap)
	if err != nil {
		// Fallback to simple format
		metadata := fmt.Sprintf(`{"question_type": "%s", "question_text": "%s", "choice_option": "%s"}`, questionType, questionText, choiceOption)
		return &metadata
	}

	metadataStr := string(metadataBytes)
	return &metadataStr
}

func (s *timeSeriesService) analyzeSentiment(text string) float64 {
	if strings.TrimSpace(text) == "" {
		return 0.0
	}

	analyzer := sentitext.Parse(text, lexicon.DefaultLexicon)
	sentiment := sentitext.PolarityScore(analyzer)

	return sentiment.Compound
}

func (s *timeSeriesService) getMetricName(questionType string, questionTexts []string) string {
	switch questionType {
	case string(feedbackModels.QuestionTypeRating):
		return "Rating Questions"
	case string(feedbackModels.QuestionTypeScale):
		return "Scale Questions"
	case string(feedbackModels.QuestionTypeYesNo):
		return "Yes/No Questions"
	case string(feedbackModels.QuestionTypeText):
		return "Text Sentiment"
	case string(feedbackModels.QuestionTypeSingleChoice):
		return "Single Choice Questions"
	case string(feedbackModels.QuestionTypeMultiChoice):
		return "Multiple Choice Questions"
	default:
		return "Other Questions"
	}
}

func (s *timeSeriesService) getMetricNameWithProduct(questionType string, productName string, questionTexts []string) string {
	// Just return the base metric name without product info to keep it clean
	return s.getMetricName(questionType, questionTexts)
}

func (s *timeSeriesService) GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) (*models.TimeSeriesResponse, error) {
	// Get raw metrics from repository (including both summary and choice-specific metrics)
	metrics, err := s.timeSeriesRepo.GetTimeSeries(ctx, request)
	if err != nil {
		return nil, err
	}
	
	// Also get choice-specific metrics for single choice questions
	choiceMetrics, err := s.getChoiceMetricsForRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	
	// Combine both metric sets
	allMetrics := append(metrics, choiceMetrics...)
	metrics = allMetrics

	// Group and process metrics, handling single choice questions specially
	seriesMap := s.groupMetricsWithChoiceHandling(metrics)

	// Convert map to slice and calculate statistics
	var series []models.TimeSeriesData
	totalPoints := 0

	for _, data := range seriesMap {
		// Sort points by timestamp
		sort.Slice(data.Points, func(i, j int) bool {
			return data.Points[i].Timestamp.Before(data.Points[j].Timestamp)
		})

		// Calculate statistics
		data.Statistics = s.calculateStatistics(data.Points)

		series = append(series, *data)
		totalPoints += len(data.Points)
	}

	// Create response
	response := &models.TimeSeriesResponse{
		Request: request,
		Series:  series,
		Summary: models.TimeSeriesSummary{
			TotalDataPoints: totalPoints,
			DateRange: models.DateRange{
				Start: request.StartDate,
				End:   request.EndDate,
			},
			Granularity: request.Granularity,
		},
	}

	return response, nil
}

func (s *timeSeriesService) GetComparison(ctx context.Context, request models.ComparisonRequest) (*models.ComparisonResponse, error) {
	// Get metrics for both periods
	period1Metrics, period2Metrics, err := s.timeSeriesRepo.GetComparison(ctx, request)
	if err != nil {
		return nil, err
	}

	// Group and compare metrics
	comparisons := s.compareMetrics(period1Metrics, period2Metrics, request)

	// Generate insights
	insights := s.generateInsights(comparisons)

	response := &models.ComparisonResponse{
		Request:     request,
		Comparisons: comparisons,
		Insights:    insights,
	}

	return response, nil
}

func (s *timeSeriesService) compareMetrics(period1Metrics, period2Metrics []models.TimeSeriesMetric, request models.ComparisonRequest) []models.TimeSeriesComparison {
	// Group metrics by type
	period1Map := s.groupMetricsByType(period1Metrics)
	period2Map := s.groupMetricsByType(period2Metrics)

	var comparisons []models.TimeSeriesComparison

	for metricType, period1Data := range period1Map {
		period2Data, exists := period2Map[metricType]

		comparison := models.TimeSeriesComparison{
			MetricType: metricType,
			MetricName: period1Data[0].MetricName,
			Period1:    s.aggregatePeriodMetrics(period1Data, request.Period1Start, request.Period1End),
			Metadata:   period1Data[0].Metadata,
		}

		if exists {
			comparison.Period2 = s.aggregatePeriodMetrics(period2Data, request.Period2Start, request.Period2End)
			comparison.Change = comparison.Period2.Value - comparison.Period1.Value
			if comparison.Period1.Value != 0 {
				comparison.ChangePercent = (comparison.Change / comparison.Period1.Value) * 100
			}
			comparison.Trend = s.determineTrend(comparison.ChangePercent)
		} else {
			comparison.Period2 = models.TimePeriodMetrics{
				StartDate: request.Period2Start,
				EndDate:   request.Period2End,
			}
			comparison.Trend = models.TrendDeclining
		}

		comparisons = append(comparisons, comparison)
	}

	return comparisons
}

func (s *timeSeriesService) groupMetricsByType(metrics []models.TimeSeriesMetric) map[string][]models.TimeSeriesMetric {
	grouped := make(map[string][]models.TimeSeriesMetric)

	for _, metric := range metrics {
		grouped[metric.MetricType] = append(grouped[metric.MetricType], metric)
	}

	return grouped
}

func (s *timeSeriesService) aggregatePeriodMetrics(metrics []models.TimeSeriesMetric, startDate, endDate time.Time) models.TimePeriodMetrics {
	if len(metrics) == 0 {
		return models.TimePeriodMetrics{
			StartDate: startDate,
			EndDate:   endDate,
		}
	}

	var total float64
	var count int64
	min := math.MaxFloat64
	max := -math.MaxFloat64

	var dataPoints []models.TimeSeriesPoint

	isYesNoQuestion := false
	isSingleChoiceQuestion := false
	isMultiChoiceQuestion := false
	var questionID *uuid.UUID

	if len(metrics) > 0 && metrics[0].Metadata != nil {
		fmt.Printf("DEBUG: Metadata=%s, MetricType=%s\n", *metrics[0].Metadata, metrics[0].MetricType)
		if strings.Contains(*metrics[0].Metadata, `"question_type": "yes_no"`) {
			isYesNoQuestion = true
			fmt.Printf("DEBUG: Detected yes/no question\n")
		} else if strings.Contains(*metrics[0].Metadata, `"question_type": "single_choice"`) {
			isSingleChoiceQuestion = true
			questionID = metrics[0].QuestionID
			fmt.Printf("DEBUG: Detected single choice question\n")
		} else if strings.Contains(*metrics[0].Metadata, `"question_type": "multi_choice"`) {
			isMultiChoiceQuestion = true
			questionID = metrics[0].QuestionID
			fmt.Printf("DEBUG: Detected multi choice question\n")
		}
	}

	for _, metric := range metrics {
		total += metric.Value
		count += metric.Count

		if metric.Value < min {
			min = metric.Value
		}
		if metric.Value > max {
			max = metric.Value
		}

		dataPoints = append(dataPoints, models.TimeSeriesPoint{
			Timestamp: metric.Timestamp,
			Value:     metric.Value,
			Count:     metric.Count,
		})
	}

	average := total / float64(len(metrics))

	// For certain question types, use average as the value instead of sum
	value := total
	if isYesNoQuestion {
		value = average
		fmt.Printf("DEBUG: Using average for yes/no question: %f\n", value)
	}

	// For rating and scale questions, also use average as the main value
	if len(metrics) > 0 && metrics[0].Metadata != nil {
		if strings.Contains(*metrics[0].Metadata, `"question_type": "rating"`) {
			value = average
			fmt.Printf("DEBUG: Using average for rating question: %f\n", value)
		} else if strings.Contains(*metrics[0].Metadata, `"question_type": "scale"`) {
			value = average
			fmt.Printf("DEBUG: Using average for scale question: %f\n", value)
		}
	}

	result := models.TimePeriodMetrics{
		StartDate:  startDate,
		EndDate:    endDate,
		Value:      value,
		Count:      count,
		Average:    average,
		Min:        min,
		Max:        max,
		DataPoints: dataPoints,
	}

	// For choice questions, get choice distribution
	if (isSingleChoiceQuestion || isMultiChoiceQuestion) && questionID != nil {
		choiceDistribution, mostPopular, topChoices := s.getChoiceDistribution(context.Background(), *questionID, startDate, endDate, isMultiChoiceQuestion)
		result.ChoiceDistribution = choiceDistribution
		result.MostPopularChoice = mostPopular
		result.TopChoices = topChoices
	}

	return result
}

func (s *timeSeriesService) calculateStatistics(points []models.TimeSeriesPoint) models.TimeSeriesStats {
	if len(points) == 0 {
		return models.TimeSeriesStats{}
	}

	var total float64
	var count int64
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, point := range points {
		total += point.Value
		count += point.Count

		if point.Value < min {
			min = point.Value
		}
		if point.Value > max {
			max = point.Value
		}
	}

	average := total / float64(len(points))

	// Calculate trend
	trendDirection, trendStrength := s.calculateTrend(points)

	return models.TimeSeriesStats{
		Total:          total,
		Average:        average,
		Min:            min,
		Max:            max,
		Count:          count,
		TrendDirection: trendDirection,
		TrendStrength:  trendStrength,
	}
}

func (s *timeSeriesService) calculateTrend(points []models.TimeSeriesPoint) (string, float64) {
	if len(points) < 2 {
		return models.TrendStable, 0.0
	}

	// Simple linear regression to determine trend
	n := float64(len(points))
	var sumX, sumY, sumXY, sumX2 float64

	for i, point := range points {
		x := float64(i)
		y := point.Value

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate slope
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Determine trend direction
	var direction string
	if math.Abs(slope) < 0.01 {
		direction = models.TrendStable
	} else if slope > 0 {
		direction = models.TrendImproving
	} else {
		direction = models.TrendDeclining
	}

	// Calculate R-squared for trend strength
	meanY := sumY / n
	var ssTot, ssRes float64

	for i, point := range points {
		x := float64(i)
		y := point.Value
		yPred := (slope * x) + ((sumY - slope*sumX) / n)

		ssTot += math.Pow(y-meanY, 2)
		ssRes += math.Pow(y-yPred, 2)
	}

	rSquared := 1 - (ssRes / ssTot)
	if math.IsNaN(rSquared) {
		rSquared = 0
	}

	return direction, math.Abs(rSquared)
}

func (s *timeSeriesService) determineTrend(changePercent float64) string {
	if math.Abs(changePercent) < 5 {
		return models.TrendStable
	} else if changePercent > 0 {
		return models.TrendImproving
	} else {
		return models.TrendDeclining
	}
}

func (s *timeSeriesService) generateInsights(comparisons []models.TimeSeriesComparison) []models.ComparisonInsight {
	var insights []models.ComparisonInsight

	for _, comp := range comparisons {
		// Generate insights based on significant changes
		if math.Abs(comp.ChangePercent) > 20 {
			severity := "info"
			if math.Abs(comp.ChangePercent) > 50 {
				severity = "warning"
			}

			insightType := "significant_change"
			message := fmt.Sprintf("%s has changed by %.1f%% between periods", comp.MetricName, comp.ChangePercent)

			insight := models.ComparisonInsight{
				Type:       insightType,
				Severity:   severity,
				Message:    message,
				MetricType: comp.MetricType,
				Change:     comp.ChangePercent,
			}

			// Add recommendations based on metric type and change
			if comp.MetricType == models.MetricTypeAverageRating && comp.ChangePercent < -10 {
				insight.Recommendation = "Review recent feedback to identify areas of concern"
			} else if comp.MetricType == models.MetricTypeFeedbackCount && comp.ChangePercent < -20 {
				insight.Recommendation = "Consider increasing customer engagement efforts"
			} else if comp.MetricType == models.MetricTypeCompletionRate && comp.ChangePercent < -15 {
				insight.Recommendation = "Simplify the feedback process to improve completion rates"
			}

			insights = append(insights, insight)
		}
	}

	return insights
}

func (s *timeSeriesService) getChoiceDistribution(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time, isMultiChoice bool) (map[string]int64, *models.ChoiceInfo, []models.ChoiceInfo) {
	// Get all feedback responses for this question within the time period
	feedbacks, err := s.feedbackService.GetByQuestionInPeriod(ctx, questionID, startDate, endDate)
	if err != nil {
		logger.Error("Failed to get feedback for choice distribution", err, logrus.Fields{
			"question_id": questionID,
			"start_date":  startDate,
			"end_date":    endDate,
		})
		return nil, nil, nil
	}

	choiceDistribution := make(map[string]int64)

	// Count choices
	for _, feedback := range feedbacks {
		if feedback.Responses != nil {
			for _, response := range feedback.Responses {
				if response.QuestionID == questionID {
					if isMultiChoice {
						// For multiple choice, the answer might be an array or comma-separated string
						switch answer := response.Answer.(type) {
						case []interface{}:
							// Handle array of choices
							for _, choice := range answer {
								if choiceStr, ok := choice.(string); ok && choiceStr != "" {
									choiceDistribution[choiceStr]++
								}
							}
						case string:
							// Handle comma-separated string or single choice
							if answer != "" {
								// Try to split by comma in case it's multiple selections
								choices := strings.Split(answer, ",")
								for _, choice := range choices {
									trimmedChoice := strings.TrimSpace(choice)
									if trimmedChoice != "" {
										choiceDistribution[trimmedChoice]++
									}
								}
							}
						}
					} else {
						// Single choice - convert answer to string
						if answerStr, ok := response.Answer.(string); ok && answerStr != "" {
							choiceDistribution[answerStr]++
						}
					}
				}
			}
		}
	}

	// Find most popular choice
	var mostPopular *models.ChoiceInfo
	var maxCount int64 = 0

	for choice, count := range choiceDistribution {
		if count > maxCount {
			maxCount = count
			mostPopular = &models.ChoiceInfo{
				Choice: choice,
				Count:  count,
			}
		}
	}

	// Get top 3 choices sorted by count
	var topChoices []models.ChoiceInfo
	type choicePair struct {
		choice string
		count  int64
	}

	var choices []choicePair
	for choice, count := range choiceDistribution {
		choices = append(choices, choicePair{choice: choice, count: count})
	}

	// Sort by count (descending)
	sort.Slice(choices, func(i, j int) bool {
		return choices[i].count > choices[j].count
	})

	// Take top 3
	limit := 3
	if len(choices) < limit {
		limit = len(choices)
	}

	for i := 0; i < limit; i++ {
		topChoices = append(topChoices, models.ChoiceInfo{
			Choice: choices[i].choice,
			Count:  choices[i].count,
		})
	}

	return choiceDistribution, mostPopular, topChoices
}

// processRatingScaleQuestion processes rating and scale questions
func (s *timeSeriesService) processRatingScaleQuestion(responses []any) float64 {
	var total float64
	validCount := 0
	for _, resp := range responses {
		if val, ok := resp.(float64); ok {
			total += val
			validCount++
		} else if intVal, ok := resp.(int); ok {
			total += float64(intVal)
			validCount++
		}
	}
	if validCount > 0 {
		return total / float64(validCount)
	}
	return 0
}

// processYesNoQuestion processes yes/no questions
func (s *timeSeriesService) processYesNoQuestion(responses []any) float64 {
	yesCount := 0
	totalCount := len(responses)

	for _, resp := range responses {
		switch v := resp.(type) {
		case bool:
			if v {
				yesCount++
			}
		case string:
			if v == "true" || v == "yes" || v == "1" {
				yesCount++
			}
		case float64:
			if v == 1 {
				yesCount++
			}
		case int:
			if v == 1 {
				yesCount++
			}
		}
	}

	if totalCount > 0 {
		return (float64(yesCount) / float64(totalCount)) * 100
	}
	return 0
}

// processTextQuestion processes text questions using sentiment analysis
func (s *timeSeriesService) processTextQuestion(responses []any) float64 {
	var totalSentiment float64
	validCount := 0

	for _, resp := range responses {
		if textResp, ok := resp.(string); ok && strings.TrimSpace(textResp) != "" {
			sentiment := s.analyzeSentiment(textResp)
			totalSentiment += sentiment
			validCount++
		}
	}

	if validCount > 0 {
		return totalSentiment / float64(validCount)
	}
	return 0
}

// processFeedbackData processes feedback and organizes it by date and question
func (s *timeSeriesService) processFeedbackData(feedbacks []feedbackModels.Feedback) (
	map[string]*ResponseData, // aggregated by question type
	map[string]*QuestionData, // individual questions
	map[string]int, // survey responses by date
) {
	timeSeriesData := make(map[string]*ResponseData)
	individualQuestionData := make(map[string]*QuestionData)
	surveyResponsesByDate := make(map[string]int)

	for _, feedback := range feedbacks {
		feedbackDate := feedback.CreatedAt.Truncate(24 * time.Hour)
		dateKey := feedbackDate.Format("2006-01-02")

		surveyResponsesByDate[dateKey]++

		if feedback.Responses == nil {
			continue
		}

		productKey := feedback.ProductID.String()
		productName := feedback.Product.Name
		if productName == "" {
			productName = "Unknown Product"
		}

		for _, response := range feedback.Responses {
			// Process for aggregated metrics (by question type)
			questionTypeKey := fmt.Sprintf("%s_%s_%s", dateKey, productKey, string(response.QuestionType))

			if _, exists := timeSeriesData[questionTypeKey]; !exists {
				timeSeriesData[questionTypeKey] = &ResponseData{
					Responses:     []any{},
					QuestionType:  string(response.QuestionType),
					QuestionTexts: []string{},
					ProductID:     productKey,
					ProductName:   productName,
				}
			}

			timeSeriesData[questionTypeKey].Responses = append(timeSeriesData[questionTypeKey].Responses, response.Answer)

			// Add question text if not already present
			questionTextExists := false
			for _, text := range timeSeriesData[questionTypeKey].QuestionTexts {
				if text == response.QuestionText {
					questionTextExists = true
					break
				}
			}
			if !questionTextExists {
				timeSeriesData[questionTypeKey].QuestionTexts = append(timeSeriesData[questionTypeKey].QuestionTexts, response.QuestionText)
			}

			// Process for individual question metrics
			if response.QuestionID != uuid.Nil {
				individualQuestionKey := fmt.Sprintf("%s_%s", dateKey, response.QuestionID.String())

				if _, exists := individualQuestionData[individualQuestionKey]; !exists {
					individualQuestionData[individualQuestionKey] = &QuestionData{
						Responses:    []any{},
						QuestionID:   response.QuestionID.String(),
						QuestionText: response.QuestionText,
						QuestionType: string(response.QuestionType),
						ProductID:    productKey,
						ProductName:  productName,
					}
				}

				individualQuestionData[individualQuestionKey].Responses = append(individualQuestionData[individualQuestionKey].Responses, response.Answer)
			}
		}
	}

	return timeSeriesData, individualQuestionData, surveyResponsesByDate
}

// processSingleChoiceQuestion creates metrics for each choice option
func (s *timeSeriesService) processSingleChoiceQuestion(
	responses []any,
	questionData *QuestionData,
	questionID string,
	accountID, organizationID uuid.UUID,
	date time.Time,
	questionMap map[uuid.UUID]*feedbackModels.Question,
) []models.TimeSeriesMetric {
	var metrics []models.TimeSeriesMetric
	choiceDistribution := make(map[string]int)

	// Count choices
	for _, resp := range responses {
		if choiceStr, ok := resp.(string); ok && strings.TrimSpace(choiceStr) != "" {
			// Clean the choice string by removing quotes and trimming whitespace
			cleanChoice := strings.Trim(strings.TrimSpace(choiceStr), `"'`)
			if cleanChoice != "" {
				choiceDistribution[cleanChoice]++
			}
		}
	}

	// Create separate metrics for each choice option
	for choice, count := range choiceDistribution {
		var productID *uuid.UUID
		if pid, err := uuid.Parse(questionData.ProductID); err == nil {
			productID = &pid
		}

		var questionUUID *uuid.UUID
		if qid, err := uuid.Parse(questionID); err == nil {
			questionUUID = &qid
		}

		// Get the question from our map
		var question *feedbackModels.Question
		if questionUUID != nil {
			question = questionMap[*questionUUID]
		}

		// Create metric type and name for this specific choice option
		choiceMetricType := fmt.Sprintf("question_%s_choice_%s", questionID, strings.ReplaceAll(strings.ToLower(choice), " ", "_"))
		choiceMetricName := fmt.Sprintf("%s - %s: %s", questionData.ProductName, questionData.QuestionText, choice)

		// Create metadata with choice information
		choiceMetadata := s.createMetadataWithChoice(questionData.QuestionType, questionData.ProductName, questionData.QuestionText, choice, question)

		metrics = append(metrics, models.TimeSeriesMetric{
			AccountID:      accountID,
			OrganizationID: organizationID,
			ProductID:      productID,
			QuestionID:     questionUUID,
			MetricType:     choiceMetricType,
			MetricName:     choiceMetricName,
			Value:          float64(count),
			Count:          int64(count),
			Timestamp:      date,
			Granularity:    models.GranularityDaily,
			Metadata:       choiceMetadata,
		})
	}

	return metrics
}

// processQuestionByType processes a question based on its type and returns the value
func (s *timeSeriesService) processQuestionByType(questionType string, responses []any) float64 {
	switch questionType {
	case string(feedbackModels.QuestionTypeRating), string(feedbackModels.QuestionTypeScale):
		return s.processRatingScaleQuestion(responses)
	case string(feedbackModels.QuestionTypeYesNo):
		return s.processYesNoQuestion(responses)
	case string(feedbackModels.QuestionTypeText):
		return s.processTextQuestion(responses)
	case string(feedbackModels.QuestionTypeSingleChoice):
		return float64(len(responses)) // For summary metric, return total count
	default:
		return float64(len(responses))
	}
}

func (s *timeSeriesService) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.timeSeriesRepo.DeleteOldMetrics(ctx, cutoffDate)
}


// isSingleChoiceQuestion checks if a question is a single choice question
func (s *timeSeriesService) isSingleChoiceQuestion(ctx context.Context, questionID string) bool {
	// Parse question UUID to validate format
	_, err := uuid.Parse(questionID)
	if err != nil {
		return false
	}
	
	// Check if choice-specific metrics exist for this question by querying the repository
	// This is more reliable than trying to get the question type directly
	pattern := fmt.Sprintf("question_%s_choice_", questionID)
	hasChoiceMetrics := s.timeSeriesRepo.HasMetricsWithPattern(ctx, pattern)
	
	return hasChoiceMetrics
}

// getChoiceSpecificMetrics returns all choice-specific metric types for a question
func (s *timeSeriesService) getChoiceSpecificMetrics(ctx context.Context, questionID string) []string {
	// Create a pattern to match choice-specific metrics
	pattern := fmt.Sprintf("question_%s_choice_", questionID)
	
	// Get all metric types that match this pattern from the repository
	choiceMetrics := s.timeSeriesRepo.GetMetricTypesByPattern(ctx, pattern)
	
	return choiceMetrics
}

// getChoiceMetricsForRequest gets choice-specific metrics for single choice questions in the request
func (s *timeSeriesService) getChoiceMetricsForRequest(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error) {
	var allChoiceMetrics []models.TimeSeriesMetric
	
	for _, metricType := range request.MetricTypes {
		if strings.HasPrefix(metricType, "question_") {
			questionID := strings.TrimPrefix(metricType, "question_")
			
			if s.isSingleChoiceQuestion(ctx, questionID) {
				choiceMetricTypes := s.getChoiceSpecificMetrics(ctx, questionID)
				
				if len(choiceMetricTypes) > 0 {
					// Create a request for the choice-specific metrics
					choiceRequest := request
					choiceRequest.MetricTypes = choiceMetricTypes
					
					choiceMetrics, err := s.timeSeriesRepo.GetTimeSeries(ctx, choiceRequest)
					if err != nil {
						return nil, err
					}
					
					allChoiceMetrics = append(allChoiceMetrics, choiceMetrics...)
				}
			}
		}
	}
	
	return allChoiceMetrics, nil
}

// groupMetricsWithChoiceHandling groups metrics and handles single choice questions specially
func (s *timeSeriesService) groupMetricsWithChoiceHandling(metrics []models.TimeSeriesMetric) map[string]*models.TimeSeriesData {
	seriesMap := make(map[string]*models.TimeSeriesData)
	choiceMetricsMap := make(map[string][]models.TimeSeriesMetric) // question_id -> choice metrics
	
	// First pass: identify all questions that have choice metrics
	questionHasChoiceMetrics := make(map[string]bool)
	for _, metric := range metrics {
		if strings.Contains(metric.MetricType, "_choice_") {
			parts := strings.Split(metric.MetricType, "_choice_")
			if len(parts) == 2 {
				baseMetricType := parts[0]
				questionHasChoiceMetrics[baseMetricType] = true
			}
		}
	}
	
	// Second pass: process metrics
	for _, metric := range metrics {
		// Check if this is a choice-specific metric
		if strings.Contains(metric.MetricType, "_choice_") {
			// Extract the base question metric type
			parts := strings.Split(metric.MetricType, "_choice_")
			if len(parts) == 2 {
				baseMetricType := parts[0] // e.g., "question_c83a860a-5a94-4390-bf11-a9af77bbf5ec"
				choiceMetricsMap[baseMetricType] = append(choiceMetricsMap[baseMetricType], metric)
				continue // Don't add choice metrics to regular series
			}
		}
		
		// Skip summary metrics for questions that have choice metrics
		if strings.HasPrefix(metric.MetricType, "question_") && questionHasChoiceMetrics[metric.MetricType] {
			continue
		}
		
		// Handle regular metrics (non-choice-specific)
		key := fmt.Sprintf("%s:%s", metric.MetricType, metric.ProductID)
		
		if series, exists := seriesMap[key]; exists {
			series.Points = append(series.Points, models.TimeSeriesPoint{
				Timestamp: metric.Timestamp,
				Value:     metric.Value,
				Count:     metric.Count,
			})
		} else {
			seriesMap[key] = &models.TimeSeriesData{
				MetricType: metric.MetricType,
				MetricName: metric.MetricName,
				ProductID:  metric.ProductID,
				Metadata:   s.parseMetadata(metric.Metadata),
				Points: []models.TimeSeriesPoint{{
					Timestamp: metric.Timestamp,
					Value:     metric.Value,
					Count:     metric.Count,
				}},
			}
		}
	}
	
	// Process choice metrics and add them as choice_series to their base question series
	for baseMetricType, choiceMetrics := range choiceMetricsMap {
		key := fmt.Sprintf("%s:%s", baseMetricType, choiceMetrics[0].ProductID)
		
		// Create or get the base series
		if _, exists := seriesMap[key]; !exists {
			// Create base series for the question
			seriesMap[key] = &models.TimeSeriesData{
				MetricType: baseMetricType,
				MetricName: s.extractQuestionNameFromChoiceMetric(choiceMetrics[0]),
				ProductID:  choiceMetrics[0].ProductID,
				Metadata:   s.parseMetadata(s.createSingleChoiceMetadata()),
				Points:     []models.TimeSeriesPoint{},
			}
		}
		
		// Group choice metrics by choice option
		choiceSeriesMap := make(map[string][]models.TimeSeriesPoint)
		for _, metric := range choiceMetrics {
			choice := s.extractChoiceFromMetricType(metric.MetricType)
			choiceSeriesMap[choice] = append(choiceSeriesMap[choice], models.TimeSeriesPoint{
				Timestamp: metric.Timestamp,
				Value:     metric.Value,
				Count:     metric.Count,
			})
		}
		
		// Convert to ChoiceSeriesData
		var choiceSeries []models.ChoiceSeriesData
		for choice, points := range choiceSeriesMap {
			// Sort points by timestamp
			sort.Slice(points, func(i, j int) bool {
				return points[i].Timestamp.Before(points[j].Timestamp)
			})
			
			choiceSeries = append(choiceSeries, models.ChoiceSeriesData{
				Choice:     choice,
				Points:     points,
				Statistics: s.calculateStatistics(points),
			})
		}
		
		seriesMap[key].ChoiceSeries = choiceSeries
	}
	
	return seriesMap
}

// extractQuestionNameFromChoiceMetric extracts the question name from a choice metric
func (s *timeSeriesService) extractQuestionNameFromChoiceMetric(metric models.TimeSeriesMetric) string {
	// Extract the question part from the metric name
	// e.g., "City Walking Tour - What type of tour experience do you prefer?: Guided Tour"
	// should become "City Walking Tour - What type of tour experience do you prefer?"
	parts := strings.Split(metric.MetricName, ": ")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ": ")
	}
	return metric.MetricName
}

// extractChoiceFromMetricType extracts the choice option from a choice-specific metric type
func (s *timeSeriesService) extractChoiceFromMetricType(metricType string) string {
	// e.g., "question_c83a860a-5a94-4390-bf11-a9af77bbf5ec_choice_guided_tour" -> "guided_tour"
	parts := strings.Split(metricType, "_choice_")
	if len(parts) == 2 {
		return strings.ReplaceAll(parts[1], "_", " ")
	}
	return metricType
}

// createSingleChoiceMetadata creates metadata indicating this is a single choice question
func (s *timeSeriesService) createSingleChoiceMetadata() *string {
	metadata := `{"question_type": "single_choice", "has_choice_series": true}`
	return &metadata
}

// parseMetadata parses a JSON metadata string into a map
func (s *timeSeriesService) parseMetadata(metadataStr *string) map[string]any {
	if metadataStr == nil {
		return nil
	}
	
	var metadata map[string]any
	err := json.Unmarshal([]byte(*metadataStr), &metadata)
	if err != nil {
		// If parsing fails, return nil or empty map
		return nil
	}
	
	return metadata
}
