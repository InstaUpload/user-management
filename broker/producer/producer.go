package producer

import (
	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/types"
)

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
