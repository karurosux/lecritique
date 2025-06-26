package services

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/lecritique/api/internal/shared/config"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email, token string) error
	SendPasswordResetEmail(ctx context.Context, email, token string) error
	SendTeamInviteEmail(ctx context.Context, email, token, companyName string) error
	SendEmailChangeVerification(ctx context.Context, newEmail, token string) error
	SendDeactivationRequest(ctx context.Context, email string, deactivationDate string) error
	SendDeactivationCancelled(ctx context.Context, email string) error
	SendAccountDeactivated(ctx context.Context, email string) error
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
	// Use frontend URL for invitation acceptance
	frontendURL := s.config.App.FrontendURL
	if frontendURL == "" {
		// Default to localhost:5173 for development
		frontendURL = "http://localhost:5173"
	}
	inviteURL := fmt.Sprintf("%s/team/accept-invite?token=%s", frontendURL, token)
	
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

func (s *emailService) SendEmailChangeVerification(ctx context.Context, newEmail, token string) error {
	subject := "Confirm Your Email Change - LeCritique"
	confirmURL := fmt.Sprintf("%s/api/v1/auth/confirm-email-change?token=%s", s.config.App.URL, token)
	
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Confirm Your Email Change</h2>
		<p>You have requested to change your email address on LeCritique.</p>
		<p>Please click the link below to confirm this change:</p>
		<p><a href="%s">Confirm Email Change</a></p>
		<p>If you didn't request this change, please ignore this email and your email address will remain unchanged.</p>
		<p>This link will expire in 24 hours.</p>
	</body>
	</html>
	`, confirmURL)

	return s.sendEmail(newEmail, subject, body)
}

func (s *emailService) SendDeactivationRequest(ctx context.Context, email string, deactivationDate string) error {
	subject := "Account Deactivation Request - LeCritique"
	
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Account Deactivation Request</h2>
		<p>We've received your request to deactivate your account.</p>
		<p><strong>Your account will be permanently deactivated on %s.</strong></p>
		<p>If you change your mind, simply log in to your account before this date to cancel the deactivation.</p>
		<p>Once deactivated, all your data will be permanently deleted and cannot be recovered.</p>
		<p>If you didn't request this, please log in to your account immediately to cancel this request.</p>
	</body>
	</html>
	`, deactivationDate)

	return s.sendEmail(email, subject, body)
}

func (s *emailService) SendDeactivationCancelled(ctx context.Context, email string) error {
	subject := "Account Deactivation Cancelled - LeCritique"
	
	body := `
	<html>
	<body>
		<h2>Account Deactivation Cancelled</h2>
		<p>Your account deactivation request has been cancelled.</p>
		<p>Your account will remain active and you can continue using all features.</p>
		<p>Thank you for staying with us!</p>
	</body>
	</html>
	`

	return s.sendEmail(email, subject, body)
}

func (s *emailService) SendAccountDeactivated(ctx context.Context, email string) error {
	subject := "Account Deactivated - LeCritique"
	
	body := `
	<html>
	<body>
		<h2>Account Deactivated</h2>
		<p>Your account has been deactivated as requested.</p>
		<p>All your data has been permanently deleted.</p>
		<p>We're sorry to see you go. If you ever want to come back, you're always welcome to create a new account.</p>
		<p>Thank you for using LeCritique.</p>
	</body>
	</html>
	`

	return s.sendEmail(email, subject, body)
}

func (s *emailService) sendEmail(to, subject, body string) error {
	// For development, just log the email instead of actually sending it
	if s.config.App.Env == "development" {
		log.Printf("=== EMAIL ===\nTo: %s\nSubject: %s\nBody: %s\n=============", to, subject, body)
		
		// Extract and highlight any links for easy clicking in terminal
		if strings.Contains(body, "href=") {
			// Simple regex to find href links
			linkRegex := regexp.MustCompile(`href="([^"]+)"`)
			matches := linkRegex.FindAllStringSubmatch(body, -1)
			if len(matches) > 0 {
				log.Println("\n🔗 CLICKABLE LINKS:")
				for _, match := range matches {
					if len(match) > 1 {
						log.Printf("   👉 %s\n", match[1])
					}
				}
				log.Println("")
			}
		}
		
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