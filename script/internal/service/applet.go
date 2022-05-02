package service

import (
	"context"
	"project/model/cache"
	"time"
)

func (s *Service) WechatTokenGet(ctx context.Context) (string, error) {
	return s.redis.Get(ctx, cache.WechatTokenKey).Result()
}

func (s *Service) WechatTokenTTL(ctx context.Context) (time.Duration, error) {
	return s.redis.TTL(ctx, cache.WechatTokenKey).Result()
}

func (s *Service) WechatTokenSet(ctx context.Context, tk string, ttl time.Duration) error {
	return s.redis.Set(ctx, cache.WechatTokenKey, tk, ttl).Err()
}
