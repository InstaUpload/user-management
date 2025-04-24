package broker

import (
	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/types"
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
		SendVerification(*types.SendVerificationKM) error
		SendWelcome(*types.SendWelcomeEmailKM) error
	}
}

func NewSender(p sarama.SyncProducer) Sender {
	return Sender{
		Email: &EmailSender{
			producer: p,
		},
	}
}
