package email

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
)

// password reset related

// SendEmail sends a generic email using the provided configuration
func SendEmail(toEmail, subject, body string) error {
	cfg := config.GetConfig()
	from := cfg.Email.From
	username := cfg.Email.Username
	password := cfg.Email.Password
	smtpHost := cfg.Email.Host
	smtpPort := fmt.Sprintf("%d", cfg.Email.Port)

	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("email configuration is missing")
	}

	message := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent successfully to: %s", toEmail)
	return nil
}

// SendPasswordResetEmail sends a password reset email
func SendPasswordResetEmail(toEmail, resetLink string, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
        Hello,
        
        You have requested to reset your password. Please click the link below:
        %s

		You can get your token from here:
		%s
        
        This link will expire in 15 minutes.
        
        If you did not request this, please ignore this email.
        
        Best regards,
        Your App Team
    `, resetLink, resetToken)

	return SendEmail(toEmail, subject, body)
}

// varification related
func SendVerificationEmail(toEmail, token string) error {
	cfg := config.GetConfig()
	from := cfg.Email.From
	username := cfg.Email.Username
	password := cfg.Email.Password
	smtpHost := cfg.Email.Host
	smtpPort := fmt.Sprintf("%d", cfg.Email.Port)

	if from == "" || username == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("email configuration is missing")
	}

	verificationLink := fmt.Sprintf("%s?token=%s", cfg.Email.VerificationURL, token)
	message := []byte(fmt.Sprintf("Subject: Email Verification\n\nPlease verify your email by clicking the link: %s", verificationLink))

	auth := smtp.PlainAuth("", username, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{toEmail}, message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Verification email sent successfully")
	return nil
}
