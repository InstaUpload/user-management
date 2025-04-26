package broker

import (
	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/broker/consumer"
	"github.com/InstaUpload/user-management/broker/producer"
)

// Re-exported types.
type (
	Sender   producer.Sender
	Receiver consumer.Receiver
)

// Re-exported functions.
var (
	NewSender   = producer.NewSender
	NewReceiver = consumer.NewReceiver
)

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	return sarama.NewConsumer(brokers, config)
}
