package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"project/model"
	"project/pkg/cache"
	"project/pkg/db"
	"project/pkg/logger"
	"project/pkg/wechat"
	"time"
)

type Service struct {
	mysql *gorm.DB
	redis *redis.Client
	//nsq    *nsq.Producer
	single *singleflight.Group
	wechat wechat.FullAPI
}

type Config struct {
	Mysql db.Mysql
	Redis cache.Redis
	Nsq   struct {
		Producer string
	}
	Wechat struct {
		Appid  string
		Secret string
	}
}

func New(cfg *Config) *Service {
	mysqlDB := db.NewMysqlDB(&cfg.Mysql)
	redisClient := cache.NewRedisClient(&cfg.Redis)
	//nsqProducer := mq.NewNsqProducer(cfg.Nsq.Producer)
	s := &Service{
		mysql: mysqlDB,
		redis: redisClient,
		//nsq:    nsqProducer,
		single: &singleflight.Group{},
	}
	httpClient := logger.NewHttpClient(8 * time.Second)
	s.wechat = wechat.NewFullAPI(cfg.Wechat.Appid, cfg.Wechat.Secret, httpClient, s.WechatToken)
	return s
}

func (s *Service) WechatToken(ctx context.Context) (string, error) {
	val, err, _ := s.single.Do("WechatToken", func() (any, error) {
		return s.redis.Get(ctx, model.KeyWechatToken).Result()
	})
	return val.(string), err
}
