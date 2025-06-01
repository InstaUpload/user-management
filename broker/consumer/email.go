package consumer

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/mail"
	"github.com/InstaUpload/user-management/types"
	"github.com/InstaUpload/user-management/utils"
)

type EmailReceiver struct {
	worker     sarama.PartitionConsumer
	mailSender interface {
		SendWelcome(*types.SendWelcomeEmailKM)
		SendVerification(*types.SendVerificationKM)
		SendEditorInvite(*types.SendEditorRequestKM)
	}
}

func NewEmailReceiver(c sarama.Consumer) EmailReceiver {
	topics, err := c.Topics()
	log.Printf("broker/consumer/email.go| list of avaliable topics: %v", topics)
	worker, err := c.ConsumePartition(types.EmailUserKT, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Can not create Email Revicer err: %s", err.Error())
	}
	// Add mail.MailSender with types.MailConfig
	mailConfig := types.MailConfig{
		Host:        utils.GetEnvString("MAILHOST", "smtp.example.com"),
		Port:        utils.GetEnvInt("MAILPORT", 587),
		SenderEmail: utils.GetEnvString("MAILSENDEREMAIL", "gpt.sahaj28@gmail.com"),
		Password:    utils.GetEnvString("MAILSENDERPASSWORD", "yourpassword"),
	}
	mailSender := mail.NewMailSender(mailConfig)
	return EmailReceiver{
		worker:     worker,
		mailSender: mailSender,
	}
}

func (e *EmailReceiver) Listen(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := <-e.worker.Messages()
		if string(msg.Key) == types.MailWelcomeKey {
			// Decode string(msg.Value) to types.SendWelcomeEmailKM and store it in data variable.
			data := types.SendWelcomeEmailKM{}
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("broker/consumer/email.go| Error unmarshalling message value: %s", err.Error())
				continue
			}
			go e.mailSender.SendWelcome(&data)
		}
		if string(msg.Key) == types.MailVerificationKey {
			data := types.SendVerificationKM{}
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("broker/consumer/email.go| Error unmarshalling message value: %s", err.Error())
				continue
			}
			go e.mailSender.SendVerification(&data)
		}
		if string(msg.Key) == types.MailEditorInviteKey {
			data := types.SendEditorRequestKM{}
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("broker/consumer/email.go| Error unmarshalling message value: %s", err.Error())
				continue
			}
			go e.mailSender.SendEditorInvite(&data)
		}
	}
}
