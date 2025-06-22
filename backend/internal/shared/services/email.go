package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"github.com/lecritique/api/internal/shared/config"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
	SendTeamInviteEmail(ctx context.Context, email, token, companyName string) error
}

type emailService struct {
	config *config.Config
}

func NewEmailService(config *config.Config) EmailService {
	return &emailService{config: config}
}

func (s *emailService) SendVerificationEmail(ctx context.Context, email, token string) error {
	subject := "Verify Your Email - LeCritique"
	verificationURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", s.config.App.URL, token)
	
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Welcome to LeCritique!</h2>
		<p>Please click the link below to verify your email address:</p>
		<p><a href="%s">Verify Email</a></p>
		<p>If you didn't create an account, please ignore this email.</p>
		<p>This link will expire in 24 hours.</p>
	</body>
	</html>
	`, verificationURL)

	return s.sendEmail(email, subject, body)
}

func (s *emailService) SendPasswordResetEmail(ctx context.Context, email, token string) error {
	subject := "Reset Your Password - LeCritique"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.App.URL, token)
	
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Password Reset Request</h2>
		<p>You requested to reset your password. Click the link below to set a new password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>If you didn't request this, please ignore this email.</p>
		<p>This link will expire in 1 hour.</p>
	</body>
	</html>
	`, resetURL)

	return s.sendEmail(email, subject, body)
}

func (s *emailService) SendTeamInviteEmail(ctx context.Context, email, token, companyName string) error {
	subject := fmt.Sprintf("Team Invitation - %s", companyName)
	inviteURL := fmt.Sprintf("%s/team/accept-invite?token=%s", s.config.App.URL, token)
	
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>You've been invited to join %s</h2>
		<p>Click the link below to accept the invitation:</p>
		<p><a href="%s">Accept Invitation</a></p>
		<p>If you don't recognize this invitation, please ignore this email.</p>
		<p>This invitation will expire in 7 days.</p>
	</body>
	</html>
	`, companyName, inviteURL)

	return s.sendEmail(email, subject, body)
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// For development, just log the email instead of actually sending it
	if s.config.App.Env == "development" {
		log.Printf("=== EMAIL ===\nTo: %s\nSubject: %s\nBody: %s\n=============", to, subject, body)
		return nil
	}

	// Production email sending (requires SMTP configuration)
	if s.config.SMTP == nil {
		log.Printf("SMTP not configured, skipping email to %s", to)
		return nil
	}

	// Compose email
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, subject, body)

	// Send email via SMTP
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port),
		auth,
		s.config.SMTP.From,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	return nil
}