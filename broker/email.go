package broker

import (
	"errors"

	"github.com/IBM/sarama"
)

type EmailSender struct {
	producer sarama.SyncProducer
}

func (e *EmailSender) SendVerificationEmail() error {
	return errors.New("UnImplemented")
}
