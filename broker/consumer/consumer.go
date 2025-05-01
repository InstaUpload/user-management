package consumer

import (
	"sync"

	"github.com/IBM/sarama"
)

type Receiver struct {
	Email interface {
		Listen(*sync.WaitGroup)
	}
}

func NewReceiver(c sarama.Consumer) Receiver {
	emailReceiver := NewEmailReceiver(c)
	return Receiver{
		Email: &emailReceiver,
	}
}

func (r *Receiver) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	var innerWg sync.WaitGroup
	innerWg.Add(1)
	go r.Email.Listen(&innerWg)
	innerWg.Wait()
}
