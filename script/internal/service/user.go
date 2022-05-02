package service

import (
	"context"
	"project/model/entity"
)

func (s *Service) UserCount(ctx context.Context, begin, end string) (count int64, err error) {
	err = s.mysql.WithContext(ctx).Model(&entity.User{}).
		Where("create_at>? AND create_at<?", begin, end).Count(&count).Error
	return
}

func (s *Service) UserAvatarUpdate(ctx context.Context, id int, path string) error {
	return s.mysql.WithContext(ctx).Model(&entity.User{}).
		Where("id=?", id).Update("avatar_url", path).Error
}
