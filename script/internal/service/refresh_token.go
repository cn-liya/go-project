package service

import (
	"encoding/hex"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"project/model"
	"project/pkg/cache"
	"project/pkg/logger"
	"project/pkg/wechat"
	"time"
)

type RefreshToken struct {
	redis  *redis.Client
	wechat wechat.BasicAPI
}

func NewRefreshToken() *RefreshToken {
	s := &RefreshToken{
		redis: cache.NewRedisClient(&config.Redis),
	}
	s.wechat = wechat.NewBasicAPI(config.Wechat.Appid, config.Wechat.Secret, logger.NewHttpClient(30*time.Second))
	return s
}

func (s *RefreshToken) WechatServerToken() {
	ctx, l := logger.NewCtxLog(hex.EncodeToString(uuid.NewV4().Bytes()), "RefreshToken", "WechatServerToken", "")
	ttl, err := s.redis.TTL(ctx, model.KeyWechatToken).Result()
	if err != nil {
		l.Error("redis.TTL error", model.KeyWechatToken, err)
		return
	}
	if ttl > 10*time.Minute {
		return
	}
	resp, err := s.wechat.GetAccessToken(ctx)
	if err != nil {
		l.Error("wechat.AccessToken error", nil, err)
		return
	}
	if resp.Errcode == 0 && resp.AccessToken != "" {
		err = s.redis.Set(ctx, model.KeyWechatToken, resp.AccessToken, time.Duration(resp.ExpiresIn)*time.Second).Err()
		if err != nil {
			l.Error("redis.Set error", model.KeyWechatToken, err)
		}
	} else {
		l.Warn("wechat.AccessToken fail", nil, resp)
	}
}
