package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/url"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	dialer *gomail.Dialer
	logger *slog.Logger
	From   string
}

func NewEmailSender(logger *slog.Logger, fromEmail, user, password string) *EmailSender {
	if fromEmail == "" {
		panic("email sender from_email cannot be empty")
	}

	if user == "" {
		panic("email sender username cannot be empty")
	}

	if password == "" {
		panic("email sender password cannot be empty")
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, user, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailSender{
		From:   fromEmail,
		dialer: d,
		logger: logger,
	}
}

func (s *EmailSender) Send(providerID string, _ int64, link string) error {
	s.logger.Info("send email", "provider_id", providerID, "link", link)
	message := gomail.NewMessage()

	message.SetAddressHeader("To", providerID, "")
	message.SetHeader("From", s.From)
	message.SetHeader("Subject", "Review Requested")

	u, err := url.Parse(link)
	if err != nil {
		return err
	}

	message.SetBody("text/html", fmt.Sprintf(`
		<html>
		<body>
			<h1>Hello! I am Review Bot</h1>
			<h3>You need to review the following request:</h3>
			<a href="%s">request</a>
		</body>
		</html>
		`, u.String()))

	return s.dialer.DialAndSend(message)
}
