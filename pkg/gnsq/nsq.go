package gnsq

import (
	"github.com/nsqio/go-nsq"
	"log"
	"project/pkg/core"
	"project/pkg/logger"
	"time"
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

type HandlerFunc func(*core.Context, *nsq.Message) error

func NewConsumer(addr, topic, channel string, concurrency int, handler HandlerFunc) *nsq.Consumer {
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = concurrency
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		log.Fatal(err)
	}
	consumer.SetLogger(silent{}, nsq.LogLevelDebug)
	v2 := core.FuncName(handler)
	var fn nsq.HandlerFunc = func(msg *nsq.Message) (err error) {
		ctx := core.ContextWithVal(string(msg.ID[:]), "", v2, "")
		begin := time.Now()
		defer func() {
			logger.FromContext(ctx).Trace("nsq", map[string]any{
				"attempts":  msg.Attempts,
				"timestamp": msg.Timestamp,
				"body":      logger.Compress(msg.Body),
			}, err, begin)
		}()
		defer core.RecoverE(ctx, &err)
		return handler(ctx, msg)
	}
	consumer.AddConcurrentHandlers(fn, concurrency)
	err = consumer.ConnectToNSQLookupd(addr)
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}

type silent struct{}

func (silent) Output(_ int, _ string) error { return nil }
