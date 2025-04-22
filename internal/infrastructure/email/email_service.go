package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type EmailService struct {
	host        string
	port        string
	senderEmail string
	senderName  string
	password    string
}

func NewEmialService(host, port, senderEmail, senderName, password string) *EmailService {
	return &EmailService{
		host:        host,
		port:        port,
		senderEmail: senderEmail,
		senderName:  senderName,
		password:    password,
	}
}

func (e *EmailService) SendEmail(to []string, subject, body string) error {
	auth := smtp.PlainAuth("", e.senderEmail, e.password, e.host)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", e.senderName, e.senderEmail)
	header["To"] = strings.Join(to, ",")
	header["Subject"] = subject
	header["NIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	message += "\r\n" + body

	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	return smtp.SendMail(addr, auth, e.senderEmail, to, []byte(message))
}

func (e *EmailService) SendResetPasswordEmail(to string, resetCode string) error {
	subject := "Reset Your Password"
	body := fmt.Sprintf(`
        <html>
            <body>
                <h1>Password Reset</h1>
                <p>You have requested to reset your password.</p>
                <p>Your 6-digit OTP is: <strong>%s</strong></p>
                <p>This code will expire in 1 hour.</p>
                <p>If you did not request a password reset, please ignore this email.</p>
            </body>
        </html>
    `, resetCode)

	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailService) SendVerificationEmail(to string, verificationCode string) error {
	subject := "Verify Your Email"
	body := fmt.Sprintf(`
        <html>
            <body>
                <h1>Email Verification</h1>
                <p>Thank you for registering. Please verify your email address to activate your account.</p>
                <p>Your 6-digit OTP is: <strong>%s</strong></p>
                <p>This code will expire in 24 hours.</p>
            </body>
        </html>
    `, verificationCode)

	return e.SendEmail([]string{to}, subject, body)
}
