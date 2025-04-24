package broker

import (
	"github.com/IBM/sarama"
)

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

type Sender struct {
	Email interface {
		SendVerificationEmail() error
	}
}

func NewSender(p sarama.SyncProducer) Sender {
	return Sender{
		Email: &EmailSender{
			producer: p,
		},
	}
}
