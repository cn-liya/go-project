package service

import (
	"github.com/nsqio/go-nsq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/pkg/db"
	"project/pkg/gnsq"
	"project/pkg/gredis"
)

type Service struct {
	mysql    *gorm.DB
	redis    *redis.Client
	producer *nsq.Producer
}

type Option func(*Service)

func NewMysql(cfg *db.Config) Option {
	return func(s *Service) {
		if s.mysql == nil {
			s.mysql = db.NewMysql(cfg)
		}
	}
}

func NewRedis(cfg *gredis.Config) Option {
	return func(s *Service) {
		if s.redis == nil {
			s.redis = gredis.NewClient(cfg)
		}
	}
}

func NewProducer(addr string) Option {
	return func(s *Service) {
		if s.producer == nil {
			s.producer = gnsq.NewProducer(addr)
		}
	}
}

func New(options ...Option) *Service {
	s := &Service{}
	for _, opt := range options {
		opt(s)
	}
	return s
}
