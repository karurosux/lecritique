package middleware

import (
	"github.com/labstack/echo/v4"
	feedbackinterface "kyooar/internal/feedback/interface"
)

type FeedbackMiddleware struct {
	feedbackService feedbackinterface.FeedbackService
}

func NewFeedbackMiddleware(feedbackService feedbackinterface.FeedbackService) *FeedbackMiddleware {
	return &FeedbackMiddleware{
		feedbackService: feedbackService,
	}
}

func (m *FeedbackMiddleware) ValidateFeedbackAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}