package mailer

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/configs"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg configs.EmailConfig
}

func NewMailer(cfg configs.EmailConfig) *Mailer {
	return &Mailer{cfg: cfg}
}

func (mail *Mailer) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mail.cfg.MailFromAddress)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		mail.cfg.MailHost,
		mail.cfg.MailPort,
		mail.cfg.MailUsername,
		mail.cfg.MailPassword,
	)

	return d.DialAndSend(m)
}

func (mail *Mailer) GenerateRandom(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal membuat token verifikasi email", err)
	}

	return hex.EncodeToString(bytes), nil
}
