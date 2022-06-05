package service

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"project/model"
	"project/pkg/logger"
	"sort"
	"time"
)

func (s *Service) GetCityBanners(ctx context.Context, city string) ([]*model.Banner, error) {
	if _, ok := model.Cities[city]; !ok {
		city = model.DefaultCity
	}
	key := model.KeyBanners + city
	val, err, _ := s.single.Do(key, func() (any, error) {
		b, err := s.redis.Get(ctx, key).Bytes()
		if err != nil && err != redis.Nil {
			return nil, err
		}
		var res []*model.Banner
		if len(b) > 0 {
			err = json.Unmarshal(b, &res)
		} else {
			err = s.mysql.WithContext(ctx).Where("city = ? AND status = ?", city, model.StatusOn).
				Find(&res).Error
			if err == nil {
				sort.Slice(res, func(i, j int) bool {
					return res[i].Sort < res[j].Sort
				})
				b, _ = json.Marshal(res)
				if err := s.redis.Set(ctx, key, b, time.Hour); err != nil {
					logger.FromContext(ctx).Error("redis.Set error", key, err)
				}
			}
		}
		return res, err
	})

	if err != nil {
		return nil, err
	}
	return val.([]*model.Banner), nil
}

//func (s *Service) PushMessage(_ context.Context, data *model.MsgExample) error {
//	b, _ := json.Marshal(data)
//	return s.nsq.Publish(model.TopicExample, b)
//}
