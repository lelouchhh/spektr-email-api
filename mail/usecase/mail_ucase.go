package usecase

import (
	"context"
	"spektr-email-api/domain"
	"time"
)

func NewMailusecase(repo domain.MailRepository, timeout time.Duration) domain.MailUsecase {
	return &MailUsecase{
		mailRepo:       repo,
		contextTimeout: timeout,
	}
}

type MailUsecase struct {
	mailRepo       domain.MailRepository
	contextTimeout time.Duration
}

func (m MailUsecase) Feedback(ctx context.Context, mail domain.Mail) error {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	err := m.mailRepo.Feedback(ctx, mail)
	if err != nil {
		return err
	}
	return nil
}
