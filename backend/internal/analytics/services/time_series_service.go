package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grassmudhorses/vader-go/lexicon"
	"github.com/grassmudhorses/vader-go/sentitext"
	"kyooar/internal/analytics/models"
	"kyooar/internal/analytics/repositories"
	feedbackModels "kyooar/internal/feedback/models"
	feedbackRepos "kyooar/internal/feedback/repositories"
	organizationRepos "kyooar/internal/organization/repositories"
	"kyooar/internal/shared/logger"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type TimeSeriesService interface {
	CollectMetrics(ctx context.Context, organizationID uuid.UUID) error
	GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) (*models.TimeSeriesResponse, error)
	GetComparison(ctx context.Context, request models.ComparisonRequest) (*models.ComparisonResponse, error)
	CleanupOldMetrics(ctx context.Context, retentionDays int) error
}

type timeSeriesService struct {
	timeSeriesRepo repositories.TimeSeriesRepository
	feedbackRepo   feedbackRepos.FeedbackRepository
	organizationRepo organizationRepos.OrganizationRepository
	analyticsService AnalyticsService
	questionRepo   feedbackRepos.QuestionRepository
}

func NewTimeSeriesService(i *do.Injector) (TimeSeriesService, error) {
	return &timeSeriesService{
		timeSeriesRepo:   do.MustInvoke[repositories.TimeSeriesRepository](i),
		feedbackRepo:     do.MustInvoke[feedbackRepos.FeedbackRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
		analyticsService: do.MustInvoke[AnalyticsService](i),
		questionRepo:     do.MustInvoke[feedbackRepos.QuestionRepository](i),
	}, nil
}

func (s *timeSeriesService) CollectMetrics(ctx context.Context, organizationID uuid.UUID) error {
	now := time.Now()
	
	logger.Info("Starting metric collection", logrus.Fields{
		"organization_id": organizationID,
	})
	
	// Get organization to retrieve the correct account_id
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}
	
	accountID := organization.AccountID
	logger.Info("Found organization", logrus.Fields{
		"organization_id": organizationID,
		"account_id": accountID,
	})
	
	// Collect organization-level metrics
	orgMetrics, err := s.collectOrganizationMetrics(ctx, organizationID, accountID, now)
	if err != nil {
		logger.Error("Failed to collect organization metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
		return err
	}
	
	logger.Info("Collected metrics", logrus.Fields{
		"organization_id": organizationID,
		"metrics_count": len(orgMetrics),
	})
	
	// Delete all existing metrics for this organization to avoid duplicates
	if err := s.timeSeriesRepo.DeleteOrganizationMetrics(ctx, organizationID); err != nil {
		logger.Error("Failed to cleanup organization metrics", err, logrus.Fields{
			"organization_id": organizationID,
		})
	}
	
	// Save metrics in batch
	if err := s.timeSeriesRepo.CreateBatch(ctx, orgMetrics); err != nil {
		return fmt.Errorf("failed to save organization metrics: %w", err)
	}
	
	logger.Info("Saved metrics successfully", logrus.Fields{
		"organization_id": organizationID,
		"metrics_count": len(orgMetrics),
	})
	
	// Collect product-level metrics
	// This would iterate through all products in the organization
	// For now, we'll skip the implementation details
	
	return nil
}

func (s *timeSeriesService) collectOrganizationMetrics(ctx context.Context, organizationID, accountID uuid.UUID, timestamp time.Time) ([]models.TimeSeriesMetric, error) {
	var metrics []models.TimeSeriesMetric
	
	// Get all feedback for the organization (limit to recent 1000 for performance)
	feedbacks, err := s.feedbackRepo.FindByOrganizationIDForAnalytics(ctx, organizationID, 1000)
	if err != nil {
		return nil, err
	}
	
	// Build a map of question IDs to question details for metadata
	questionDetailsMap := make(map[uuid.UUID]*feedbackModels.Question)
	for _, feedback := range feedbacks {
		if feedback.Responses != nil {
			for _, response := range feedback.Responses {
				if response.QuestionID != uuid.Nil {
					if _, exists := questionDetailsMap[response.QuestionID]; !exists {
						// Get question details including labels
						question, err := s.questionRepo.GetQuestionByID(ctx, response.QuestionID)
						if err == nil {
							questionDetailsMap[response.QuestionID] = question
						}
					}
				}
			}
		}
	}
	
	// Group feedback by date, product, and question type for proper time series
	timeSeriesData := make(map[string]map[string]map[string]struct {
		responses []interface{}
		questionType string
		questionTexts []string
		productID string
		productName string
	})
	
	surveyResponsesByDate := make(map[string]int)
	
	for _, feedback := range feedbacks {
		feedbackDate := feedback.CreatedAt.Truncate(24 * time.Hour)
		dateKey := feedbackDate.Format("2006-01-02")
		
		surveyResponsesByDate[dateKey]++
		
		if feedback.Responses != nil {
			productKey := feedback.ProductID.String()
			productName := "Product"
			
			for _, response := range feedback.Responses {
				questionTypeKey := string(response.QuestionType)
				
				if timeSeriesData[dateKey] == nil {
					timeSeriesData[dateKey] = make(map[string]map[string]struct {
						responses []interface{}
						questionType string
						questionTexts []string
						productID string
						productName string
					})
				}
				
				if timeSeriesData[dateKey][productKey] == nil {
					timeSeriesData[dateKey][productKey] = make(map[string]struct {
						responses []interface{}
						questionType string
						questionTexts []string
						productID string
						productName string
					})
				}
				
				if _, exists := timeSeriesData[dateKey][productKey][questionTypeKey]; !exists {
					timeSeriesData[dateKey][productKey][questionTypeKey] = struct {
						responses []interface{}
						questionType string
						questionTexts []string
						productID string
						productName string
					}{
						responses: []interface{}{},
						questionType: string(response.QuestionType),
						questionTexts: []string{},
						productID: productKey,
						productName: productName,
					}
				}
				
				questionData := timeSeriesData[dateKey][productKey][questionTypeKey]
				questionData.responses = append(questionData.responses, response.Answer)
				
				questionTextExists := false
				for _, text := range questionData.questionTexts {
					if text == response.QuestionText {
						questionTextExists = true
						break
					}
				}
				if !questionTextExists {
					questionData.questionTexts = append(questionData.questionTexts, response.QuestionText)
				}
				
				timeSeriesData[dateKey][productKey][questionTypeKey] = questionData
			}
		}
	}
	
	for dateKey, dateData := range timeSeriesData {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		
		for _, productData := range dateData {
			for questionTypeKey, questionData := range productData {
				if len(questionData.responses) == 0 {
					continue
				}
				
				var value float64
				count := len(questionData.responses)
				
				switch questionData.questionType {
				case string(feedbackModels.QuestionTypeRating), string(feedbackModels.QuestionTypeScale):
					var total float64
					validCount := 0
					for _, resp := range questionData.responses {
						if val, ok := resp.(float64); ok {
							total += val
							validCount++
						} else if intVal, ok := resp.(int); ok {
							total += float64(intVal)
							validCount++
						}
					}
					if validCount > 0 {
						value = total / float64(validCount)
					}
				case string(feedbackModels.QuestionTypeYesNo):
					yesCount := 0
					for _, resp := range questionData.responses {
						if val, ok := resp.(bool); ok && val {
							yesCount++
						}
					}
					if count > 0 {
						value = (float64(yesCount) / float64(count)) * 100
					}
				case string(feedbackModels.QuestionTypeText):
					var totalSentiment float64
					validCount := 0
					
					for _, resp := range questionData.responses {
						if textResp, ok := resp.(string); ok && strings.TrimSpace(textResp) != "" {
							sentiment := s.analyzeSentiment(textResp)
							totalSentiment += sentiment
							validCount++
						}
					}
					
					if validCount > 0 {
						value = totalSentiment / float64(validCount)
					} else {
						value = 0.0
					}
				default:
					value = float64(count)
				}
				
				// Create question type metric
				metricType := questionTypeKey + "_questions"
				metricName := s.getMetricNameWithProduct(questionData.questionType, questionData.productName, questionData.questionTexts)
				
				var productID *uuid.UUID
				if pid, err := uuid.Parse(questionData.productID); err == nil {
					productID = &pid
				}
				
				logger.Info("Creating question type metric", logrus.Fields{
					"date": dateKey,
					"product": questionData.productName,
					"question_type": questionTypeKey,
					"metric_type": metricType,
					"value": value,
					"count": count,
				})
				
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
					Metadata:       s.createMetadataWithProduct(questionData.questionType, questionData.productName),
				})
				
			}
		}
	}
	
	// Create individual question metrics
	individualQuestionData := make(map[string]map[string]struct {
		responses []interface{}
		questionID string
		questionText string
		questionType string
		productID string
		productName string
	})
	
	// Group by date and question ID
	for _, feedback := range feedbacks {
		feedbackDate := feedback.CreatedAt.Truncate(24 * time.Hour)
		dateKey := feedbackDate.Format("2006-01-02")
		productKey := feedback.ProductID.String()
		productName := "Product"
		
		if feedback.Responses != nil {
			for _, response := range feedback.Responses {
				if response.QuestionID == uuid.Nil {
					continue
				}
				
				questionKey := response.QuestionID.String()
				
				if individualQuestionData[dateKey] == nil {
					individualQuestionData[dateKey] = make(map[string]struct {
						responses []interface{}
						questionID string
						questionText string
						questionType string
						productID string
						productName string
					})
				}
				
				if _, exists := individualQuestionData[dateKey][questionKey]; !exists {
					individualQuestionData[dateKey][questionKey] = struct {
						responses []interface{}
						questionID string
						questionText string
						questionType string
						productID string
						productName string
					}{
						responses: []interface{}{},
						questionID: questionKey,
						questionText: response.QuestionText,
						questionType: string(response.QuestionType),
						productID: productKey,
						productName: productName,
					}
				}
				
				qData := individualQuestionData[dateKey][questionKey]
				qData.responses = append(qData.responses, response.Answer)
				individualQuestionData[dateKey][questionKey] = qData
			}
		}
	}
	
	// Create metrics for individual questions
	for dateKey, questionMap := range individualQuestionData {
		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}
		
		for questionID, qData := range questionMap {
			if len(qData.responses) == 0 {
				continue
			}
			
			var questionValue float64
			questionCount := len(qData.responses)
			
			switch qData.questionType {
			case string(feedbackModels.QuestionTypeRating), string(feedbackModels.QuestionTypeScale):
				var total float64
				validCount := 0
				for _, resp := range qData.responses {
					if val, ok := resp.(float64); ok {
						total += val
						validCount++
					} else if intVal, ok := resp.(int); ok {
						total += float64(intVal)
						validCount++
					}
				}
				if validCount > 0 {
					questionValue = total / float64(validCount)
				}
			case string(feedbackModels.QuestionTypeYesNo):
				yesCount := 0
				for _, resp := range qData.responses {
					if val, ok := resp.(bool); ok && val {
						yesCount++
					}
				}
				if questionCount > 0 {
					questionValue = (float64(yesCount) / float64(questionCount)) * 100
				}
			case string(feedbackModels.QuestionTypeText):
				var totalSentiment float64
				validCount := 0
				
				for _, resp := range qData.responses {
					if textResp, ok := resp.(string); ok && strings.TrimSpace(textResp) != "" {
						sentiment := s.analyzeSentiment(textResp)
						totalSentiment += sentiment
						validCount++
					}
				}
				
				if validCount > 0 {
					questionValue = totalSentiment / float64(validCount)
				} else {
					questionValue = 0.0
				}
			default:
				questionValue = float64(questionCount)
			}
			
			var productID *uuid.UUID
			if pid, err := uuid.Parse(qData.productID); err == nil {
				productID = &pid
			}
			
			var questionUUID *uuid.UUID
			if qid, err := uuid.Parse(questionID); err == nil {
				questionUUID = &qid
			}
			
			// Get question details from our map
			var questionDetails *feedbackModels.Question
			if questionUUID != nil {
				questionDetails = questionDetailsMap[*questionUUID]
			}
			
			// Keep the question ID in metric type for frontend compatibility
			questionMetricType := "question_" + questionID
			// Use just the question text, without product info
			questionMetricName := qData.questionText
			
			logger.Info("Creating individual question metric", logrus.Fields{
				"date": dateKey,
				"product": qData.productName,
				"question_id": questionID,
				"question_text": qData.questionText,
				"metric_type": questionMetricType,
				"value": questionValue,
				"count": questionCount,
			})
			
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
				Metadata:       s.createMetadataWithQuestion(qData.questionType, qData.productName, qData.questionText, questionDetails),
			})
		}
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
	// Get raw metrics from repository
	metrics, err := s.timeSeriesRepo.GetTimeSeries(ctx, request)
	if err != nil {
		return nil, err
	}
	
	// Group metrics by type and process
	seriesMap := make(map[string]*models.TimeSeriesData)
	
	for _, metric := range metrics {
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
				Metadata:   metric.Metadata,
				Points: []models.TimeSeriesPoint{{
					Timestamp: metric.Timestamp,
					Value:     metric.Value,
					Count:     metric.Count,
				}},
			}
		}
	}
	
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
	
	return models.TimePeriodMetrics{
		StartDate:  startDate,
		EndDate:    endDate,
		Value:      total,
		Count:      count,
		Average:    average,
		Min:        min,
		Max:        max,
		DataPoints: dataPoints,
	}
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

func (s *timeSeriesService) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.timeSeriesRepo.DeleteOldMetrics(ctx, cutoffDate)
}