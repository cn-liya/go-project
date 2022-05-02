package service

import (
	"context"
	"github.com/nsqio/go-nsq"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"project/model/cache"
	"project/pkg/db"
	"project/pkg/gnsq"
	"project/pkg/gredis"
)

type Service struct {
	mysql    *gorm.DB
	redis    *gredis.Client
	producer *nsq.Producer
	//rabbit   *amqp.Channel
	single *singleflight.Group
}

type Config struct {
	Mysql db.Config
	Redis gredis.Config
	Nsq   struct {
		Producer string
	}
	//Rabbitmq rabbitmq.Config
}

func New(cfg *Config) *Service {
	return &Service{
		mysql:    db.NewMysql(&cfg.Mysql),
		redis:    gredis.NewPlusClient(&cfg.Redis),
		producer: gnsq.NewProducer(cfg.Nsq.Producer),
		//rabbit:   rabbitmq.NewChannel(&cfg.Rabbitmq),
		single: &singleflight.Group{},
	}
}

func (s *Service) WechatToken(ctx context.Context) (string, error) {
	val, err, _ := s.single.Do("WechatToken", func() (any, error) {
		return s.redis.Get(ctx, cache.WechatTokenKey).Result()
	})
	return val.(string), err
}
