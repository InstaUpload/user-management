package producer

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/types"
)

type EmailSender struct {
	producer sarama.SyncProducer
}

func (e *EmailSender) SendVerification(*types.SendVerificationKM) error {
	return errors.New("Kafka send Verification Email function is UnImplemented")
}

func (e *EmailSender) SendWelcome(message *types.SendWelcomeEmailKM) error {
	// use json.Marshal() to convert the struct to json and send it to kafka topic.
	messageInBytes, err := json.Marshal(&message)
	key := []byte(types.MailWelcomeKey)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: types.EmailUserKT,
		Value: sarama.StringEncoder(messageInBytes),
		Key:   sarama.StringEncoder(key),
	}
	partition, offset, err := e.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", types.EmailUserKT, partition, offset)
	return nil
}
