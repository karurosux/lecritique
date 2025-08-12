package analyticsservice

import (
	"context"
	"encoding/json"
	"fmt"
	models "kyooar/internal/analytics/model"
	analyticsinterface "kyooar/internal/analytics/interface"
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
	"github.com/sirupsen/logrus"
)

type ResponseData struct {
	Responses     []any
	QuestionType  string
	QuestionTexts []string
	ProductID     string
	ProductName   string
}

type QuestionData struct {
	Responses    []any
	QuestionID   string
	QuestionText string
	QuestionType string
	ProductID    string
	ProductName  string
}

type TimeSeriesService struct {
	timeSeriesRepo      analyticsinterface.TimeSeriesRepository
	feedbackService     feedbackServices.FeedbackService
	organizationService organizationServices.OrganizationService
	analyticsService    analyticsinterface.AnalyticsService
	questionService     feedbackServices.QuestionService
}

func NewTimeSeriesService(
	timeSeriesRepo analyticsinterface.TimeSeriesRepository,
	feedbackService feedbackServices.FeedbackService,
	organizationService organizationServices.OrganizationService,
	analyticsService analyticsinterface.AnalyticsService,
	questionService feedbackServices.QuestionService,
) *TimeSeriesService {
	return &TimeSeriesService{
		timeSeriesRepo:      timeSeriesRepo,
		feedbackService:     feedbackService,
		organizationService: organizationService,
		analyticsService:    analyticsService,
		questionService:     questionService,
	}
}

func (s *TimeSeriesService) CollectMetrics(ctx context.Context, organizationID uuid.UUID) error {
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

	return nil
}

func (s *TimeSeriesService) collectOrganizationMetrics(ctx context.Context, organizationID, accountID uuid.UUID, timestamp time.Time) ([]models.TimeSeriesMetric, error) {
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

		value := s.processQuestionByType(questionData.QuestionType, questionData.Responses)
		count := len(questionData.Responses)

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

	for key, qData := range individualQuestionData {
		parts := strings.Split(key, "_")
		if len(parts) < 2 {
			continue
		}

		dateKey := parts[0]
		questionID := strings.Join(parts[1:], "_")

		date, err := time.Parse("2006-01-02", dateKey)
		if err != nil {
			continue
		}

		if len(qData.Responses) == 0 {
			continue
		}

		questionCount := len(qData.Responses)

		if qData.QuestionType == string(feedbackModels.QuestionTypeSingleChoice) {
			choiceMetrics := s.processSingleChoiceQuestion(
				qData.Responses, qData, questionID, accountID, organizationID, date, questionMap,
			)
			metrics = append(metrics, choiceMetrics...)
			continue
		}
		
		if qData.QuestionType == string(feedbackModels.QuestionTypeMultiChoice) {
			choiceMetrics := s.processMultiChoiceQuestion(
				qData.Responses, qData, questionID, accountID, organizationID, date, questionMap,
			)
			metrics = append(metrics, choiceMetrics...)
			continue
		}

		questionValue := s.processQuestionByType(qData.QuestionType, qData.Responses)

		var productID *uuid.UUID
		if pid, err := uuid.Parse(qData.ProductID); err == nil {
			productID = &pid
		}

		var questionUUID *uuid.UUID
		if qid, err := uuid.Parse(questionID); err == nil {
			questionUUID = &qid
		}

		questionMetricType := "question_" + questionID
		questionMetricName := fmt.Sprintf("%s - %s", qData.ProductName, qData.QuestionText)

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

func (s *TimeSeriesService) createMetadata(questionType string) *string {
	metadata := fmt.Sprintf(`{"question_type": "%s"}`, questionType)
	return &metadata
}

func (s *TimeSeriesService) createMetadataWithProduct(questionType string, productName string) *string {
	metadata := fmt.Sprintf(`{"question_type": "%s"}`, questionType)
	return &metadata
}

func (s *TimeSeriesService) createMetadataWithQuestion(questionType string, productName string, questionText string, question *feedbackModels.Question) *string {
	metadataMap := map[string]interface{}{
		"question_type": questionType,
		"question_text": questionText,
	}

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

	metadataBytes, err := json.Marshal(metadataMap)
	if err != nil {
			metadata := fmt.Sprintf(`{"question_type": "%s", "question_text": "%s"}`, questionType, questionText)
		return &metadata
	}

	metadataStr := string(metadataBytes)
	return &metadataStr
}

func (s *TimeSeriesService) createMetadataWithChoice(questionType string, productName string, questionText string, choiceOption string, question *feedbackModels.Question) *string {
	metadataMap := map[string]interface{}{
		"question_type": questionType,
		"question_text": questionText,
		"choice_option": choiceOption,
	}

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

	metadataBytes, err := json.Marshal(metadataMap)
	if err != nil {
			metadata := fmt.Sprintf(`{"question_type": "%s", "question_text": "%s", "choice_option": "%s"}`, questionType, questionText, choiceOption)
		return &metadata
	}

	metadataStr := string(metadataBytes)
	return &metadataStr
}

func (s *TimeSeriesService) analyzeSentiment(text string) float64 {
	if strings.TrimSpace(text) == "" {
		return 0.0
	}

	analyzer := sentitext.Parse(text, lexicon.DefaultLexicon)
	sentiment := sentitext.PolarityScore(analyzer)

	return sentiment.Compound
}

func (s *TimeSeriesService) getMetricName(questionType string, questionTexts []string) string {
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

func (s *TimeSeriesService) getMetricNameWithProduct(questionType string, productName string, questionTexts []string) string {
	return s.getMetricName(questionType, questionTexts)
}

func (s *TimeSeriesService) GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) (*models.TimeSeriesResponse, error) {
	metrics, err := s.timeSeriesRepo.GetTimeSeries(ctx, request)
	if err != nil {
		return nil, err
	}
	
	choiceMetrics, err := s.getChoiceMetricsForRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	
	allMetrics := append(metrics, choiceMetrics...)
	metrics = allMetrics

	seriesMap := s.groupMetricsWithChoiceHandling(metrics)

	var series []models.TimeSeriesData
	totalPoints := 0

	for _, data := range seriesMap {
		sort.Slice(data.Points, func(i, j int) bool {
			return data.Points[i].Timestamp.Before(data.Points[j].Timestamp)
		})

		data.Statistics = s.calculateStatistics(data.Points)

		series = append(series, *data)
		totalPoints += len(data.Points)
	}

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

func (s *TimeSeriesService) GetComparison(ctx context.Context, request models.ComparisonRequest) (*models.ComparisonResponse, error) {
	period1Metrics, period2Metrics, err := s.timeSeriesRepo.GetComparison(ctx, request)
	if err != nil {
		return nil, err
	}

	comparisons := s.compareMetrics(period1Metrics, period2Metrics, request)

	insights := s.generateInsights(comparisons)

	response := &models.ComparisonResponse{
		Request:     request,
		Comparisons: comparisons,
		Insights:    insights,
	}

	return response, nil
}

func (s *TimeSeriesService) compareMetrics(period1Metrics, period2Metrics []models.TimeSeriesMetric, request models.ComparisonRequest) []models.TimeSeriesComparison {
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

func (s *TimeSeriesService) groupMetricsByType(metrics []models.TimeSeriesMetric) map[string][]models.TimeSeriesMetric {
	grouped := make(map[string][]models.TimeSeriesMetric)

	for _, metric := range metrics {
		grouped[metric.MetricType] = append(grouped[metric.MetricType], metric)
	}

	return grouped
}

func (s *TimeSeriesService) aggregatePeriodMetrics(metrics []models.TimeSeriesMetric, startDate, endDate time.Time) models.TimePeriodMetrics {
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

	value := total
	if isYesNoQuestion {
		value = average
		fmt.Printf("DEBUG: Using average for yes/no question: %f\n", value)
	}

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

	if (isSingleChoiceQuestion || isMultiChoiceQuestion) && questionID != nil {
		choiceDistribution, mostPopular, topChoices := s.getChoiceDistribution(context.Background(), *questionID, startDate, endDate, isMultiChoiceQuestion)
		result.ChoiceDistribution = choiceDistribution
		result.MostPopularChoice = mostPopular
		result.TopChoices = topChoices
	}

	return result
}

func (s *TimeSeriesService) calculateStatistics(points []models.TimeSeriesPoint) models.TimeSeriesStats {
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

func (s *TimeSeriesService) calculateTrend(points []models.TimeSeriesPoint) (string, float64) {
	if len(points) < 2 {
		return models.TrendStable, 0.0
	}

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

	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	var direction string
	if math.Abs(slope) < 0.01 {
		direction = models.TrendStable
	} else if slope > 0 {
		direction = models.TrendImproving
	} else {
		direction = models.TrendDeclining
	}

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

func (s *TimeSeriesService) determineTrend(changePercent float64) string {
	if math.Abs(changePercent) < 5 {
		return models.TrendStable
	} else if changePercent > 0 {
		return models.TrendImproving
	} else {
		return models.TrendDeclining
	}
}

func (s *TimeSeriesService) generateInsights(comparisons []models.TimeSeriesComparison) []models.ComparisonInsight {
	var insights []models.ComparisonInsight

	for _, comp := range comparisons {
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

func (s *TimeSeriesService) getChoiceDistribution(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time, isMultiChoice bool) (map[string]int64, *models.ChoiceInfo, []models.ChoiceInfo) {
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

	for _, feedback := range feedbacks {
		if feedback.Responses != nil {
			for _, response := range feedback.Responses {
				if response.QuestionID == questionID {
					if isMultiChoice {
									switch answer := response.Answer.(type) {
						case []interface{}:
									for _, choice := range answer {
								if choiceStr, ok := choice.(string); ok && choiceStr != "" {
									choiceDistribution[choiceStr]++
								}
							}
						case string:
									if answer != "" {
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
								if answerStr, ok := response.Answer.(string); ok && answerStr != "" {
							choiceDistribution[answerStr]++
						}
					}
				}
			}
		}
	}

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

	var topChoices []models.ChoiceInfo
	type choicePair struct {
		choice string
		count  int64
	}

	var choices []choicePair
	for choice, count := range choiceDistribution {
		choices = append(choices, choicePair{choice: choice, count: count})
	}

	sort.Slice(choices, func(i, j int) bool {
		return choices[i].count > choices[j].count
	})

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

func (s *TimeSeriesService) processRatingScaleQuestion(responses []any) float64 {
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

func (s *TimeSeriesService) processYesNoQuestion(responses []any) float64 {
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

func (s *TimeSeriesService) processTextQuestion(responses []any) float64 {
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

func (s *TimeSeriesService) processFeedbackData(feedbacks []feedbackModels.Feedback) (
	map[string]*ResponseData,
	map[string]*QuestionData,
	map[string]int,
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

func (s *TimeSeriesService) processSingleChoiceQuestion(
	responses []any,
	questionData *QuestionData,
	questionID string,
	accountID, organizationID uuid.UUID,
	date time.Time,
	questionMap map[uuid.UUID]*feedbackModels.Question,
) []models.TimeSeriesMetric {
	var metrics []models.TimeSeriesMetric
	choiceDistribution := make(map[string]int)

	for _, resp := range responses {
		if choiceStr, ok := resp.(string); ok && strings.TrimSpace(choiceStr) != "" {
			cleanChoice := strings.Trim(strings.TrimSpace(choiceStr), `"'`)
			if cleanChoice != "" {
				choiceDistribution[cleanChoice]++
			}
		}
	}

	for choice, count := range choiceDistribution {
		var productID *uuid.UUID
		if pid, err := uuid.Parse(questionData.ProductID); err == nil {
			productID = &pid
		}

		var questionUUID *uuid.UUID
		if qid, err := uuid.Parse(questionID); err == nil {
			questionUUID = &qid
		}

		var question *feedbackModels.Question
		if questionUUID != nil {
			question = questionMap[*questionUUID]
		}

		choiceMetricType := fmt.Sprintf("question_%s_choice_%s", questionID, strings.ReplaceAll(strings.ToLower(choice), " ", "_"))
		choiceMetricName := fmt.Sprintf("%s - %s: %s", questionData.ProductName, questionData.QuestionText, choice)

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

func (s *TimeSeriesService) processMultiChoiceQuestion(
	responses []any,
	questionData *QuestionData,
	questionID string,
	accountID, organizationID uuid.UUID,
	date time.Time,
	questionMap map[uuid.UUID]*feedbackModels.Question,
) []models.TimeSeriesMetric {
	var metrics []models.TimeSeriesMetric
	choiceDistribution := make(map[string]int)

	for _, resp := range responses {
		var choices []string
		
		switch v := resp.(type) {
		case string:
			if strings.TrimSpace(v) != "" {
				var jsonArray []string
				if err := json.Unmarshal([]byte(v), &jsonArray); err == nil {
					choices = jsonArray
				} else {
					splitChoices := strings.Split(v, ",")
					for _, choice := range splitChoices {
						cleaned := strings.Trim(strings.TrimSpace(choice), `"'`)
						if cleaned != "" {
							choices = append(choices, cleaned)
						}
					}
				}
			}
		case []interface{}:
			for _, choice := range v {
				if choiceStr, ok := choice.(string); ok {
					cleaned := strings.Trim(strings.TrimSpace(choiceStr), `"'`)
					if cleaned != "" {
						choices = append(choices, cleaned)
					}
				}
			}
		case []string:
			for _, choice := range v {
				cleaned := strings.Trim(strings.TrimSpace(choice), `"'`)
				if cleaned != "" {
					choices = append(choices, cleaned)
				}
			}
		}

		for _, choice := range choices {
			if choice != "" {
				choiceDistribution[choice]++
			}
		}
	}

	for choice, count := range choiceDistribution {
		var productID *uuid.UUID
		if pid, err := uuid.Parse(questionData.ProductID); err == nil {
			productID = &pid
		}

		var questionUUID *uuid.UUID
		if qid, err := uuid.Parse(questionID); err == nil {
			questionUUID = &qid
		}

		var question *feedbackModels.Question
		if questionUUID != nil {
			question = questionMap[*questionUUID]
		}

		choiceMetricType := fmt.Sprintf("question_%s_choice_%s", questionID, strings.ReplaceAll(strings.ToLower(choice), " ", "_"))
		choiceMetricName := fmt.Sprintf("%s - %s: %s", questionData.ProductName, questionData.QuestionText, choice)

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

func (s *TimeSeriesService) processQuestionByType(questionType string, responses []any) float64 {
	switch questionType {
	case string(feedbackModels.QuestionTypeRating), string(feedbackModels.QuestionTypeScale):
		return s.processRatingScaleQuestion(responses)
	case string(feedbackModels.QuestionTypeYesNo):
		return s.processYesNoQuestion(responses)
	case string(feedbackModels.QuestionTypeText):
		return s.processTextQuestion(responses)
	case string(feedbackModels.QuestionTypeSingleChoice):
		return float64(len(responses))
	case string(feedbackModels.QuestionTypeMultiChoice):
		return float64(len(responses))
	default:
		return float64(len(responses))
	}
}

func (s *TimeSeriesService) CleanupOldMetrics(ctx context.Context, retentionDays int) error {
	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	return s.timeSeriesRepo.DeleteOldMetrics(ctx, cutoffDate)
}


func (s *TimeSeriesService) isSingleChoiceQuestion(ctx context.Context, questionID string) bool {
	_, err := uuid.Parse(questionID)
	if err != nil {
		return false
	}
	
	pattern := fmt.Sprintf("question_%s_choice_", questionID)
	hasChoiceMetrics := s.timeSeriesRepo.HasMetricsWithPattern(ctx, pattern)
	
	return hasChoiceMetrics
}

func (s *TimeSeriesService) getChoiceSpecificMetrics(ctx context.Context, questionID string) []string {
	pattern := fmt.Sprintf("question_%s_choice_", questionID)
	
	choiceMetrics := s.timeSeriesRepo.GetMetricTypesByPattern(ctx, pattern)
	
	return choiceMetrics
}

func (s *TimeSeriesService) getChoiceMetricsForRequest(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error) {
	var allChoiceMetrics []models.TimeSeriesMetric
	
	for _, metricType := range request.MetricTypes {
		if strings.HasPrefix(metricType, "question_") {
			questionID := strings.TrimPrefix(metricType, "question_")
			
			if s.isSingleChoiceQuestion(ctx, questionID) {
				choiceMetricTypes := s.getChoiceSpecificMetrics(ctx, questionID)
				
				if len(choiceMetricTypes) > 0 {
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

func (s *TimeSeriesService) groupMetricsWithChoiceHandling(metrics []models.TimeSeriesMetric) map[string]*models.TimeSeriesData {
	seriesMap := make(map[string]*models.TimeSeriesData)
	choiceMetricsMap := make(map[string][]models.TimeSeriesMetric)
	
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
	
	for _, metric := range metrics {
		if strings.Contains(metric.MetricType, "_choice_") {
			parts := strings.Split(metric.MetricType, "_choice_")
			if len(parts) == 2 {
				baseMetricType := parts[0]
				choiceMetricsMap[baseMetricType] = append(choiceMetricsMap[baseMetricType], metric)
				continue
			}
		}
		
		if strings.HasPrefix(metric.MetricType, "question_") && questionHasChoiceMetrics[metric.MetricType] {
			continue
		}
		
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
	
	for baseMetricType, choiceMetrics := range choiceMetricsMap {
		key := fmt.Sprintf("%s:%s", baseMetricType, choiceMetrics[0].ProductID)
		
		if _, exists := seriesMap[key]; !exists {
			var choiceMetadata *string
			if len(choiceMetrics) > 0 && choiceMetrics[0].Metadata != nil {
				if strings.Contains(*choiceMetrics[0].Metadata, `"question_type": "multi_choice"`) {
					choiceMetadata = s.createMultiChoiceMetadata()
				} else {
					choiceMetadata = s.createSingleChoiceMetadata()
				}
			} else {
				choiceMetadata = s.createSingleChoiceMetadata()
			}
			
			seriesMap[key] = &models.TimeSeriesData{
				MetricType: baseMetricType,
				MetricName: s.extractQuestionNameFromChoiceMetric(choiceMetrics[0]),
				ProductID:  choiceMetrics[0].ProductID,
				Metadata:   s.parseMetadata(choiceMetadata),
				Points:     []models.TimeSeriesPoint{},
			}
		}
		
		choiceSeriesMap := make(map[string][]models.TimeSeriesPoint)
		for _, metric := range choiceMetrics {
			choice := s.extractChoiceFromMetricType(metric.MetricType)
			choiceSeriesMap[choice] = append(choiceSeriesMap[choice], models.TimeSeriesPoint{
				Timestamp: metric.Timestamp,
				Value:     metric.Value,
				Count:     metric.Count,
			})
		}
		
		var choiceSeries []models.ChoiceSeriesData
		for choice, points := range choiceSeriesMap {
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

func (s *TimeSeriesService) extractQuestionNameFromChoiceMetric(metric models.TimeSeriesMetric) string {
	parts := strings.Split(metric.MetricName, ": ")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], ": ")
	}
	return metric.MetricName
}

func (s *TimeSeriesService) extractChoiceFromMetricType(metricType string) string {
	parts := strings.Split(metricType, "_choice_")
	if len(parts) == 2 {
		return strings.ReplaceAll(parts[1], "_", " ")
	}
	return metricType
}

func (s *TimeSeriesService) createSingleChoiceMetadata() *string {
	metadata := `{"question_type": "single_choice", "has_choice_series": true}`
	return &metadata
}

func (s *TimeSeriesService) createMultiChoiceMetadata() *string {
	metadata := `{"question_type": "multi_choice", "has_choice_series": true}`
	return &metadata
}

func (s *TimeSeriesService) parseMetadata(metadataStr *string) map[string]any {
	if metadataStr == nil {
		return nil
	}
	
	var metadata map[string]any
	err := json.Unmarshal([]byte(*metadataStr), &metadata)
	if err != nil {
			return nil
	}
	
	return metadata
}
