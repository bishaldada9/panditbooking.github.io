package email

import (
	"fmt"
	"log"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/bees/hindu-ritual-platform/pkg/configs"
)

type EmailService struct {
	config   *configs.EmailConfig
	mockMode bool
}

func NewEmailService(cfg *configs.EmailConfig, mockMode bool) *EmailService {
	return &EmailService{
		config:   cfg,
		mockMode: mockMode,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	if s.mockMode {
		log.Printf("[MOCK EMAIL] To: %s, Subject: %s, Body: %s", to, subject, body)
		return nil
	}

	if s.config.APIKey == "" {
		return fmt.Errorf("email API key is not configured")
	}

	from := mail.NewEmail(s.config.FromName, s.config.FromEmail)
	recipient := mail.NewEmail("", to)
	htmlContent := mail.NewContent("text/html", body)
	m := mail.NewV3MailInit(from, subject, recipient, htmlContent)

	request := sendgrid.GetRequest(s.config.APIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	response, err := sendgrid.API(request)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("email API returned status %d: %s", response.StatusCode, response.Body)
	}

	return nil
}

func (s *EmailService) SendVerificationEmail(to, token string) error {
	subject := "Verify Your Email - Hindu Ritual Platform"
	body := fmt.Sprintf(`
		<h2>Email Verification</h2>
		<p>Thank you for registering. Please verify your email address by clicking the link below:</p>
		<p><a href="https://platform.hinduritual.com/verify-email?token=%s">Verify Email</a></p>
		<p>This link will expire in 24 hours.</p>
		<p>If you did not create an account, please ignore this email.</p>
		<br>
		<p>Best regards,<br>Hindu Ritual Platform Team</p>
	`, token)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	subject := "Password Reset - Hindu Ritual Platform"
	body := fmt.Sprintf(`
		<h2>Password Reset Request</h2>
		<p>You have requested to reset your password. Click the link below to set a new password:</p>
		<p><a href="https://platform.hinduritual.com/reset-password?token=%s">Reset Password</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request a password reset, please ignore this email.</p>
		<br>
		<p>Best regards,<br>Hindu Ritual Platform Team</p>
	`, token)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendBookingConfirmation(to, name, bookingDetails string) error {
	subject := "Booking Confirmation - Hindu Ritual Platform"
	body := fmt.Sprintf(`
		<h2>Booking Confirmed!</h2>
		<p>Dear %s,</p>
		<p>Your booking has been confirmed. Here are the details:</p>
		<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">
			%s
		</div>
		<p>If you have any questions, please contact our support team.</p>
		<br>
		<p>Best regards,<br>Hindu Ritual Platform Team</p>
	`, name, bookingDetails)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendNotificationEmail(to, subject, message string) error {
	body := fmt.Sprintf(`
		<h2>Notification</h2>
		<p>%s</p>
		<br>
		<p>Best regards,<br>Hindu Ritual Platform Team</p>
	`, message)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendBookingReminder(to, name, details string) error {
	subject := "Booking Reminder - Hindu Ritual Platform"
	body := fmt.Sprintf(`
		<h2>Upcoming Booking Reminder</h2>
		<p>Dear %s,</p>
		<p>This is a reminder for your upcoming ritual booking:</p>
		<div style="background: #f9f9f9; padding: 15px; border-radius: 5px;">
			%s
		</div>
		<p>Please be prepared for the scheduled ritual.</p>
		<br>
		<p>Best regards,<br>Hindu Ritual Platform Team</p>
	`, name, details)
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendRaw(to, subject, body string) error {
	return s.SendEmail(to, subject, body)
}

func (s *EmailService) SendBatchEmails(recipients []string, subject, body string) error {
	for _, to := range recipients {
		if err := s.SendEmail(strings.TrimSpace(to), subject, body); err != nil {
			return fmt.Errorf("failed to send email to %s: %w", to, err)
		}
	}
	return nil
}
