package gnsq

import (
	"github.com/nsqio/go-nsq"
	"log"
)

func NewProducer(addr string) *nsq.Producer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	producer.SetLogger(silent{}, nsq.LogLevelDebug)
	err = producer.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return producer
}

func NewConsumer(addr, topic, channel string, concurrent int, handler nsq.HandlerFunc) *nsq.Consumer {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = concurrent
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatal(err)
	}
	consumer.SetLogger(silent{}, nsq.LogLevelDebug)
	consumer.AddConcurrentHandlers(handler, concurrent)
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}

type silent struct{}

func (silent) Output(_ int, _ string) error {
	return nil
}
