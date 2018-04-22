package publisher

import (
	"sync"

	"github.com/lab46/monorepo/gopkg/log"
	"github.com/lab46/monorepo/gopkg/mq"
	nsq "github.com/nsqio/go-nsq"
)

// FakeProducer struct
type FakeProducer struct {
	consumers mq.Subscribers
	guard     sync.Mutex
}

func newFakeProducer() *FakeProducer {
	f := FakeProducer{
		guard:     sync.Mutex{},
		consumers: make(mq.Subscribers),
	}
	return &f
}

// Ping producer
func (p *FakeProducer) Ping() error {
	return nil
}

// Publish message
func (p *FakeProducer) Publish(topic string, body []byte) error {
	for _, channel := range p.consumers[topic] {
		msg := nsq.Message{
			ID:   mq.FakeID,
			Body: body,
		}
		err := channel.Handler.HandleMessage(&msg)
		if err != nil {
			log.Errors(err)
			return err
		}
	}
	return nil
}

// RegisterSubcribers function
func (p *FakeProducer) RegisterSubcribers(subs mq.Subscribers) error {
	for topicName, topic := range subs {
		for channelName, channel := range topic {
			if _, ok := p.consumers[topicName]; !ok {
				p.consumers[topicName] = make(map[string]mq.SubscriberOptions)
			}
			p.consumers[topicName][channelName] = channel
		}
	}
	return nil
}
