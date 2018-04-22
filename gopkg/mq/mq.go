package mq

import (
	"time"

	nsq "github.com/nsqio/go-nsq"
)

type (
	// Config struct
	Config struct {
		ListenAddress  []string `yaml:"listenaddress"`
		PublishAddress string   `yaml:"publishaddress"`
		MaxInFlight    int      `yaml:"maxinflight"`
		MaxAttempts    int      `yaml:"maxattempts"`
	}

	// MQ struct
	MQ interface {
		Register(string, string, func(*nsq.Message) error, int)
		Run() error
	}

	// Publisher interface
	Publisher interface {
		Publish(string, []byte) error
		Ping() error
	}

	// SubscriberOptions struct
	SubscriberOptions struct {
		Handler    nsq.Handler
		Concurrent int
	}

	// Consumer struct
	Consumer struct {
		topic       string
		channel     string
		handler     nsq.Handler
		concurrent  int
		maxAttemps  uint16
		maxInFlight int
	}

	// Subscribers of message queue
	Subscribers map[string]map[string]SubscriberOptions
)

// FakeID for fake nsq message
var FakeID nsq.MessageID

func init() {
	copy(FakeID[:], "fake")
}

// NewNSQLHandler to create handler of nsq
func NewNSQLHandler(h func(*nsq.Message) error) nsq.HandlerFunc {
	return nsq.HandlerFunc(h)
}

// Finish function
func Finish(msg *nsq.Message) {
	if msg == nil || nsq.MessageID((msg.ID)) == FakeID {
		return
	}
	msg.Finish()
}

// Requeue function
func Requeue(msg *nsq.Message, t time.Duration) {
	if msg == nil || nsq.MessageID((msg.ID)) == FakeID {
		return
	}
	msg.Requeue(t)
}

// RequeueWithoutBackoff function
func RequeueWithoutBackoff(msg *nsq.Message, t time.Duration) {
	if msg == nil || nsq.MessageID((msg.ID)) == FakeID {
		return
	}
	msg.RequeueWithoutBackoff(t)
}
