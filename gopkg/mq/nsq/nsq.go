package nsq

import (
	"time"

	"github.com/lab46/monorepo/gopkg/log"
	"github.com/lab46/monorepo/mq"
	"github.com/nsqio/go-nsq"
)

var (
	err error
)

type subscriber struct {
	topic       string
	channel     string
	handler     nsq.Handler
	concurrent  int
	maxAttemps  uint16
	maxInFlight int
}

// MQ struct
type MQ struct {
	options     *Options
	consumers   []subscriber
	listenErrCh chan error
	birth       time.Time
}

// Options struct
type Options struct {
	ListenAddress  []string
	PublishAddress string
	Prefix         string
	MaxInFlight    int
	MaxAttempts    int
}

// New MQ
func New(o *Options) *MQ {
	m := &MQ{
		options:     o,
		listenErrCh: make(chan error),
		birth:       time.Now(),
	}

	// adding consumers to mq
	// for _, s := range subscribers {
	// 	m.consumers = append(m.consumers, s)
	// }
	return m
}

// RegisterSubcribers function
func (m *MQ) RegisterSubcribers(subs mq.Subscribers) error {
	for topicName, topic := range subs {
		for channelName, channel := range topic {
			m.consumers = append(m.consumers, subscriber{
				topic:      topicName,
				channel:    channelName,
				handler:    channel.Handler,
				concurrent: channel.Concurrent,
			})
		}
	}
	return nil
}

// Register consumer
func (m *MQ) Register(topic, channel string, handler nsq.Handler, concurrent int) {
	log.Debugf("Registering MQ Consumer : %s/%s concurrent=%d", topic, channel, concurrent)
	m.consumers = append(m.consumers, subscriber{
		topic:      topic,
		channel:    channel,
		handler:    handler,
		concurrent: concurrent,
	})
}

// RegisterWithoutPrefix for mq consumer
func (m *MQ) RegisterWithoutPrefix(topic, channel string, handler nsq.Handler, concurrents ...int) {
	var concurrent int
	if len(concurrents) > 0 && concurrents[0] > 0 {
		concurrent = concurrents[0]
	}
	log.Debugf("Registering MQ Consumer : %s/%s concurrent=%d", topic, channel, concurrent)
	m.consumers = append(m.consumers, subscriber{
		topic:      topic,
		channel:    channel,
		handler:    handler,
		concurrent: concurrent,
	})
}

// Run listener
func (m *MQ) Run() error {
	if m.options.PublishAddress == "" {
		m.options.PublishAddress = m.options.ListenAddress[0]
	}
	for _, consumer := range m.consumers {
		config := nsq.NewConfig()
		config.MaxAttempts = uint16(m.options.MaxAttempts)
		config.MaxInFlight = m.options.MaxInFlight
		// config.MaxAttempts = consumer.maxAttemps
		// config.MaxInFlight = consumer.maxInFlight
		q, err := nsq.NewConsumer(consumer.topic, consumer.channel, config)
		if err != nil {
			return err
		}

		if consumer.concurrent != 0 {
			q.AddConcurrentHandlers(consumer.handler, consumer.concurrent)
		} else {
			q.AddHandler(consumer.handler)
		}

		if len(m.options.ListenAddress) > 1 {
			err = q.ConnectToNSQLookupds(m.options.ListenAddress)
		} else {
			err = q.ConnectToNSQLookupd(m.options.ListenAddress[0])
		}

		if err != nil {
			return err
		}
	}
	return nil
}
