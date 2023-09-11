package domain

import "context"

type Mail struct {
	Body    string `json:"body" from:"body"`
	To      string `json:"to" from:"to"`
	Subject string `json:"subject" from:"subject"`
}

type MailUsecase interface {
	Feedback(ctx context.Context, main Mail) error
}
type MailRepository interface {
	Feedback(ctx context.Context, main Mail) error
}
type EmailsConsumer interface {
	StartConsumer(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error
}
