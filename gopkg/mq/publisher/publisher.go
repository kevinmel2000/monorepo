package publisher

import (
	"encoding/json"
	"sync"

	"github.com/lab46/monorepo/gopkg/env"
	"github.com/lab46/monorepo/gopkg/log"
	"github.com/lab46/monorepo/gopkg/mq"
	nsq "github.com/nsqio/go-nsq"
)

// Producer struct
type Producer struct {
	Name     string
	Address  string
	Attempts int
	Prefix   string
}

var (
	defaultProducer mq.Publisher

	producer  *nsq.Producer
	producers map[string]*nsq.Producer
	mtx       = sync.RWMutex{}
	prefix    string
)

// InitNsqPublisher publisher
func InitNsqPublisher(p Producer) {
	var err error
	config := nsq.NewConfig()
	config.MaxAttempts = uint16(p.Attempts)

	producer, err = nsq.NewProducer(p.Address, config)
	if err != nil {
		panic("Failed to create NSQ producer")
	}
	defaultProducer = producer
}

// InitFakePublisher to publish using simplensq
func InitFakePublisher() {
	e := env.GetCurrentServiceEnv()
	if env.ServiceEnv(e) == env.ProductionEnv {
		log.Fatal("should not use fake publisher in production environment")
		return
	}
	fake := newFakeProducer()
	defaultProducer = fake
}

// RegisterFakeSubscriber for publishing to fake subscriber
func RegisterFakeSubscriber(subs mq.Subscribers) {
	switch defaultProducer.(type) {
	case *FakeProducer:
		p := defaultProducer.(*FakeProducer)
		p.RegisterSubcribers(subs)
	default:
		return
	}
}

// AddPublisher to add new publisher
func AddPublisher(p Producer) {
	config := nsq.NewConfig()
	config.MaxAttempts = uint16(p.Attempts)

	pr, err := nsq.NewProducer(p.Address, config)
	if err != nil {
		panic("Failed to create NSQ producer")
	}
	mtx.Lock()
	producers[p.Name] = pr
	mtx.Unlock()
}

// Publish data to message queue.
func Publish(topic string, data interface{}) (err error) {
	var payload []byte
	// topic = prefix + topic
	payload, err = json.Marshal(data)
	if err != nil {
		return err
	}

	if defaultProducer == nil {
		return
	}
	// log.Debugf("[MQ] Publishing topic=%s payload=%s", topic, string(payload))
	return defaultProducer.Publish(topic, payload)
}

// PublishString will only publish string
func PublishString(topic string, data string) (err error) {
	log.Println("publishing string with topic ", topic)
	// topic = prefix + topic
	payload := []byte(data)

	if defaultProducer == nil {
		return
	}
	return defaultProducer.Publish(topic, payload)
}

// PublishWithoutPrefix will publish data to message queue. No prefix added in topic
func PublishWithoutPrefix(topic string, data interface{}) (err error) {
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		return
	}

	if defaultProducer == nil {
		return
	}
	// log.Debugf("[MQ] Publishing topic=%s payload=%s", topic, string(payload))
	return defaultProducer.Publish(topic, payload)
}

// PublishCustomProducer will publish to custom producer name
func PublishCustomProducer(producerName, topic string, data interface{}) (err error) {
	var payload []byte
	payload, err = json.Marshal(data)
	if err != nil {
		return
	}
	mtx.RLock()
	p := producers[producerName]
	mtx.RUnlock()
	if p == nil {
		return
	}
	// log.Debugf("[MQ] Publishing topic=%s payload=%s", topic, string(payload))
	return defaultProducer.Publish(topic, payload)
}
