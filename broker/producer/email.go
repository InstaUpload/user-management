package producer

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/types"
)

type EmailSender struct {
	producer sarama.SyncProducer
}

func (e *EmailSender) SendVerification(message *types.SendVerificationKM) error {
	messageInBytes, err := json.Marshal(&message)
	key := []byte(types.MailVerificationKey)
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
	log.Printf("broker/producer/email.go|l:31 Message is stored in topic(%s)/partition(%d)/offset(%d)\n", types.EmailUserKT, partition, offset)
	return nil
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
	log.Printf("broker/producer/email.go|l:52 Message is stored in topic(%s)/partition(%d)/offset(%d)\n", types.EmailUserKT, partition, offset)
	return nil
}

func (e *EmailSender) SendEditorInvite(message *types.SendEditorRequestKM) error {
	// use json.Marshal() to convert the struct to json and send it to kafka topic.
	messageInBytes, err := json.Marshal(&message)
	key := []byte(types.MailEditorInviteKey)
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
	log.Printf("broker/producer/email.go|l:70 Message is stored in topic(%s)/partition(%d)/offset(%d)\n", types.EmailUserKT, partition, offset)
	return nil
}
