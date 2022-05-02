package service

import (
	"context"
)

func (s *Service) UserCount(ctx context.Context, begin, end string) (count int64, err error) {
	err = s.mysql.WithContext(ctx).
		Raw(`SELECT COUNT(*) FROM user WHERE create_time>? and create_time<?`, begin, end).
		Scan(&count).Error
	return
}

func (s *Service) UserAvatarUpdate(ctx context.Context, id int, path string) error {
	return s.mysql.WithContext(ctx).Where("id", id).Update("avatar_url", path).Error
}
