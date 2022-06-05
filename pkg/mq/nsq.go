package mq

import (
	"github.com/nsqio/go-nsq"
	"log"
)

func NewNsqProducer(addr string) *nsq.Producer {
	producer, err := nsq.NewProducer(addr, nsq.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	err = producer.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return producer
}

func NewNsqConsumer(addr, topic, channel string, concurrent int, handler nsq.HandlerFunc) *nsq.Consumer {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = concurrent
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatal(err)
	}
	consumer.AddConcurrentHandlers(handler, concurrent)
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}
