package simplensq

import (
	"errors"
	"sync"
	"time"

	"github.com/lab46/monorepo/gopkg/log"
)

type (
	consumerHandler func(msg *Message)
	consumerChannel map[string]*Consumer
)

// Config struct
type Config struct {
	MaxAttempts uint64
	MaxInflight int
}

// NewConfig return address of config
func NewConfig() *Config {
	c := Config{}
	return &c
}

// Message struct
type Message struct {
	ID        string
	Body      []byte
	Timestamp int64
	Attempts  uint16

	topic string
}

// Requeue message
func (m *Message) Requeue(delay time.Duration) {
	time.Sleep(delay)
	m.Attempts++
	for _, channel := range listeners[m.topic] {
		channel.receiveMessage(m)
	}
}

// RequeueWithoutBackoff message
func (m *Message) RequeueWithoutBackoff(delay time.Duration) {
	time.Sleep(delay)
	m.Attempts++
	for _, channel := range listeners[m.topic] {
		channel.receiveMessage(m)
	}
}

// Finish message
func (m *Message) Finish() {
	return
}

var listeners map[string]consumerChannel

// Consumer struct
type Consumer struct {
	handler consumerHandler
	config  *Config
	running bool

	msgChan  chan *Message
	doneChan chan bool

	topic   string
	channel string
}

// NewConsumer to create new consumer
func NewConsumer(topic, channel string, config *Config) *Consumer {
	if c, ok := listeners[topic][channel]; ok {
		return c
	}

	c := &Consumer{
		config: config,
		// craete a buffered channel of 10
		msgChan: make(chan *Message, 10),
		topic:   topic,
		channel: channel,
	}
	cc := consumerChannel{}
	cc[channel] = c
	return c
}

// AddHandler to consumer
func (c *Consumer) AddHandler(h consumerHandler) {
	c.handler = h
	if !c.running {
		c.run()
	}
}

// AddConcurrentHandlers to handle message concurrently
func (c *Consumer) AddConcurrentHandlers(h consumerHandler, concurrency int) {
	c.handler = h
	if !c.running {
		c.run()
	}
}

func (c *Consumer) receiveMessage(msg *Message) {
	log.Debugf("send message %+v to topic %s", *msg, msg.topic)
	c.msgChan <- msg
}

// Run consumer
func (c *Consumer) run() {
	c.running = true
	for {
		select {
		case msg := <-c.msgChan:
			if msg == nil {
				continue
			}
			c.handler(msg)
		case <-c.doneChan:
			c.running = false
			return
		}
	}
}

func (c *Consumer) stop() {
	log.Debugf("consumer for topic %s and channel %s exiting", c.topic, c.channel)
	c.doneChan <- true
}

// ConnectToNSQDs function
func (c *Consumer) ConnectToNSQDs(listenerAddress []string) error {
	return nil
}

// ConnectToNSQD function
func (c *Consumer) ConnectToNSQD(listenerAddress string) error {
	return nil
}

// Producer struct
type Producer struct {
	guard sync.Mutex
}

// NewProducer return producer for queue
func NewProducer() *Producer {
	p := Producer{
		guard: sync.Mutex{},
	}
	return &p
}

// Ping producer
func (p *Producer) Ping() error {
	return nil
}

// Publish message
func (p *Producer) Publish(topic string, body []byte) error {
	m := Message{
		ID:        "someid",
		Body:      body,
		Timestamp: time.Now().UnixNano(),
		Attempts:  1,
		topic:     topic,
	}

	if _, ok := listeners[topic]; !ok {
		return errors.New("topic is not exist")
	}

	for _, channel := range listeners[topic] {
		channel.receiveMessage(&m)
	}
	return nil
}
