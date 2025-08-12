package constants

import "errors"

var (
	ErrFeedbackNotFound       = errors.New("feedback not found")
	ErrQuestionNotFound       = errors.New("question not found")
	ErrQuestionnaireNotFound  = errors.New("questionnaire not found")
	ErrInvalidProductID       = errors.New("invalid product ID")
	ErrInvalidOrganizationID  = errors.New("invalid organization ID")
	ErrInvalidFeedbackData    = errors.New("invalid feedback data")
	ErrUnauthorizedAccess     = errors.New("unauthorized access")
	ErrDuplicateFeedback      = errors.New("duplicate feedback submission")
	ErrInvalidQuestionType    = errors.New("invalid question type")
	ErrMissingRequiredField   = errors.New("missing required field")
)