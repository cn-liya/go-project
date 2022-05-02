package service

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func (s *Service) ExampleLLen(ctx context.Context) (int64, error) {
	return s.redis.LLen(ctx, "example:list").Result()
}

func (s *Service) ExampleRPop(ctx context.Context) ([]byte, error) {
	b, err := s.redis.RPop(ctx, "example:list").Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	return b, nil
}
