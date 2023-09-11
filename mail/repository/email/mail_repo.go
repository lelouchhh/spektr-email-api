package email

import (
	"context"
	"fmt"
	"net/smtp"
	"spektr-email-api/domain"
)

type MailRepository struct {
	From     string
	FromName string
	Pass     string
}

func (m MailRepository) Feedback(ctx context.Context, mail domain.Mail) error {
	fmt.Println(m)
	err := m.sendEmail(mail)
	if err != nil {
		return err
	}
	return nil
}

func NewMailRepository(from, fromName, pass string) domain.MailRepository {
	return &MailRepository{
		From:     from,
		FromName: fromName,
		Pass:     pass,
	}
}

func (m MailRepository) sendEmail(email domain.Mail) error {
	smtpHost := "smtp.yandex.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", m.From, m.Pass, smtpHost)

	message := fmt.Sprintf(
		"To: %s\r\n"+
			"From: %s<%s>\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		email.To,
		m.FromName,
		m.From,
		email.Subject,
		email.Body,
	)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, m.From, []string{email.To}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
