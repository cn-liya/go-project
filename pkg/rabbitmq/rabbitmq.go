package rabbitmq

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"project/pkg/core"
	"project/pkg/logger"
	"sync"
	"time"
)

type Config struct {
	Address       string
	Username      string
	Password      string
	Cert, Key, Ca string
}

// NewConnect 初始化单个连接 for consumer
func NewConnect(cfg *Config) *amqp.Connection {
	scheme := "amqp"
	var tlsConfig *tls.Config
	if cfg.Cert != "" && cfg.Key != "" && cfg.Ca != "" {
		certificate, err := tls.X509KeyPair([]byte(cfg.Cert), []byte(cfg.Key))
		if err != nil {
			log.Fatal(err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(cfg.Ca)) {
			log.Fatal("failed to parse root certificate")
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      pool,
		}
		scheme = "amqps"
	}
	dsn := scheme + "://" + cfg.Username + ":" + cfg.Password + "@" + cfg.Address
	conn, err := amqp.DialTLS(dsn, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// NewChannel 初始化单个连接并创建单个信道 for producer
func NewChannel(cfg *Config) *amqp.Channel {
	conn := NewConnect(cfg)
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	return ch
}

/*
消费者方法简单封装，队列已在mq管理后台或其他地方创建。
初始化一个*Consumer，指定连接，队列名，并发数，逻辑方法
handler方法返回nil会自动ack消息，返回error会reject并判断是否requeue
调用*Consumer的Stop方法，会停止获取消息，并阻塞等待正在处理消息的handler方法完成
*/

type Consumer struct {
	ch   *amqp.Channel
	stop context.CancelFunc
	wg   *sync.WaitGroup
}

type HandlerFunc func(*core.Context, *amqp.Delivery) error

func NewConsumer(conn *amqp.Connection, qname string, concurrency int, handler HandlerFunc) *Consumer {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	if err = ch.Qos(1, 0, false); err != nil {
		log.Fatal(err)
	}
	v2 := core.FuncName(handler)
	fn := func(msg *amqp.Delivery) (err error) {
		ctx := core.ContextWithVal(msg.MessageId, "", v2, "")
		begin := time.Now()
		defer func() {
			logger.FromContext(ctx).Trace("rabbitmq", map[string]any{
				"headers":   msg.Headers,
				"timestamp": msg.Timestamp.UnixNano(),
				"body":      logger.Compress(msg.Body),
			}, err, begin)
		}()
		defer core.RecoverE(ctx, &err)
		return handler(ctx, msg)
	}
	ctx, cancel := context.WithCancel(context.Background())
	consumer := &Consumer{
		ch:   ch,
		stop: cancel,
		wg:   &sync.WaitGroup{},
	}
	for i := 0; i < concurrency; i++ {
		consumer.wg.Add(1)
		go func() {
			msgs, err := ch.ConsumeWithContext(ctx, qname, "", false, false, false, false, nil)
			if err != nil {
				log.Fatal(err)
			}
			for d := range msgs {
				if fn(&d) == nil { // 处理成功确认
					_ = d.Ack(false)
				} else { // 处理失败拒绝，为避免无限requeue，仅当首次出错或队列有delivery-limit时requeue
					requeue := !d.Redelivered || (d.Headers != nil && d.Headers["x-delivery-count"] != nil)
					_ = d.Reject(requeue)
				}
			}
			consumer.wg.Done()
		}()
	}
	return consumer
}

func (c *Consumer) Stop() {
	c.stop()
	c.wg.Wait()
	_ = c.ch.Close()
}
