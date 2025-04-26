package consumer

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/InstaUpload/user-management/types"
)

type EmailReceiver struct {
	worker sarama.PartitionConsumer
}

func NewEmailReceiver(c sarama.Consumer) EmailReceiver {
	worker, err := c.ConsumePartition(types.VerificationKT, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Can not create Email Revicer err: %s", err.Error())
	}
	return EmailReceiver{
		worker: worker,
	}
}

func (e *EmailReceiver) ReceiveVerification() error {
	return errors.New("Email Receiver is unimplemented")
}

func (e *EmailReceiver) Listen(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		msg := <-e.worker.Messages()
		log.Printf("Received message: %s", string(msg.Value))
		time.Sleep(10 * time.Second)
		log.Printf("Working on email receiver")
		break
	}
}
