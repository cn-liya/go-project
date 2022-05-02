package gredis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"project/pkg/logger"
	"time"
)

// Client 继承*redis.Client的所有方法，原先使用*redis.Client的地方可以无缝切换。
type Client struct {
	*redis.Client
}

func NewPlusClient(cfg *Config) *Client { return &Client{NewClient(cfg)} }

// FetchXXX 入参：context, 缓存key, 结果指针, 缓存未命中的查询方法, 缓存时间。
// fn闭包函数内修改结果指针，返回error，闭包作用域内可以访问外层的ctx,查询参数,连接Client等。

func (c *Client) FetchJSON(ctx context.Context, key string, res any, fn func() error, ttl time.Duration) error {
	b, err := c.Get(ctx, key).Bytes()
	if err == redis.Nil {
		if err = fn(); err != nil {
			return err
		}
		b, _ = json.Marshal(res)
		if err = c.Set(ctx, key, b, ttl).Err(); err != nil {
			logger.FromContext(ctx).Error("redis.Set", key, err)
		}
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, res)
	return err
}

func (c *Client) FetchString(ctx context.Context, key string, res *string, fn func() error, ttl time.Duration) error {
	s, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		if err = fn(); err != nil {
			return err
		}
		if err = c.Set(ctx, key, *res, ttl).Err(); err != nil {
			logger.FromContext(ctx).Error("redis.Set", key, err)
		}
		return nil
	}
	if err != nil {
		return err
	}
	*res = s
	return nil
}

func (c *Client) FetchInt(ctx context.Context, key string, res *int, fn func() error, ttl time.Duration) error {
	i, err := c.Get(ctx, key).Int()
	if err == redis.Nil {
		if err = fn(); err != nil {
			return err
		}
		if err = c.Set(ctx, key, *res, ttl).Err(); err != nil {
			logger.FromContext(ctx).Error("redis.Set", key, err)
		}
		return nil
	}
	if err != nil {
		return err
	}
	*res = i
	return nil
}
