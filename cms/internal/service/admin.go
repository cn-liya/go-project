package service

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"project/cms/internal/proto"
	"project/model/cache"
	"project/model/entity"
	"project/pkg/logger"
	"project/pkg/random"
	"time"
)

func (s *Service) AdminTokenSet(ctx context.Context, data *cache.AdminToken) (string, error) {
	ssoKey := cache.AdminSSOKey(data.ID)
	oldToken, err := s.redis.Get(ctx, ssoKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	b, _ := json.Marshal(data)
	h := sha1.New()
	h.Write(b)
	h.Write(random.Bytes(10))
	newToken := base32.StdEncoding.EncodeToString(h.Sum(nil))
	tx := s.redis.TxPipeline()
	if oldToken != "" {
		tx.Del(ctx, cache.AdminTokenKey(oldToken))
	}
	tx.Set(ctx, cache.AdminTokenKey(newToken), b, time.Hour)
	tx.Set(ctx, ssoKey, newToken, time.Hour)
	_, err = tx.Exec(ctx)
	return newToken, err
}

func (s *Service) AdminTokenGet(ctx context.Context, token string) (*cache.AdminToken, error) {
	key := cache.AdminTokenKey(token)
	b, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var result cache.AdminToken
	if len(b) == 0 {
		return &result, nil
	}
	err = json.Unmarshal(b, &result)
	tx := s.redis.TxPipeline()
	tx.Expire(ctx, key, time.Hour)
	tx.Expire(ctx, cache.AdminSSOKey(result.ID), time.Hour)
	if _, err := tx.Exec(ctx); err != nil {
		logger.FromContext(ctx).Error("redis.TxPipeline.Exec error", nil, err)
	}
	return &result, err
}

func (s *Service) AdminUserLogout(ctx context.Context, id int) error {
	ssoKey := cache.AdminSSOKey(id)
	token, err := s.redis.Get(ctx, ssoKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if token != "" {
		return s.redis.Del(ctx, cache.AdminTokenKey(token), ssoKey).Err()
	}
	return nil
}

func (s *Service) AdminUserFindByName(ctx context.Context, name string) (*entity.AdminUser, error) {
	var result entity.AdminUser
	err := s.mysql.WithContext(ctx).Where("username=?", name).Preload("AdminRole").Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminUserFindByID(ctx context.Context, id int) (*entity.AdminUser, error) {
	var result entity.AdminUser
	err := s.mysql.WithContext(ctx).Where("id=?", id).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminUserUpdate(ctx context.Context, data *entity.AdminUser) error {
	opt := s.mysql.WithContext(ctx).Updates(data)
	if opt.Error != nil {
		return opt.Error
	}
	if opt.RowsAffected > 0 && (data.Status == entity.StatusOff || data.RoleID != 0) {
		return s.AdminUserLogout(ctx, data.ID) //禁用或变更角色强制退出登录
	}
	return nil
}

func (s *Service) AdminRoleList(ctx context.Context,
	p *proto.ListArgs) (total int64, list []*entity.AdminRole, err error) {
	query := s.mysql.WithContext(ctx).Model(&entity.AdminRole{})
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) AdminRoles(ctx context.Context) ([]*entity.AdminRole, error) {
	var result []*entity.AdminRole
	err := s.mysql.WithContext(ctx).Select("id", "name").
		Order("id").Find(&result).Error
	return result, err
}

func (s *Service) AdminRoleFindByID(ctx context.Context, id int) (*entity.AdminRole, error) {
	var result entity.AdminRole
	err := s.mysql.WithContext(ctx).Where("id=?", id).Take(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &result, nil
}

func (s *Service) AdminRoleCreate(ctx context.Context, data *entity.AdminRole) error {
	return s.mysql.WithContext(ctx).Create(data).Error
}

func (s *Service) AdminRoleUpdate(ctx context.Context, data *entity.AdminRole) error {
	return s.mysql.WithContext(ctx).Updates(data).Error
}

func (s *Service) AdminUserList(ctx context.Context,
	p *proto.AdminUserListArgs) (total int64, list []*entity.AdminUser, err error) {
	query := s.mysql.WithContext(ctx).Model(&entity.AdminUser{})
	if p.RoleID > 0 {
		query = query.Where("role_id=?", p.RoleID)
	}
	if p.Username != "" {
		query = query.Where("username LIKE ?", p.Username+"%")
	}
	if p.Status != 0 {
		query = query.Where("status=?", p.Status)
	}
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) AdminUserCreate(ctx context.Context, data *entity.AdminUser) (bool, error) {
	opt := s.mysql.WithContext(ctx).FirstOrCreate(data, "username=?", data.Username)
	return opt.RowsAffected > 0, opt.Error
}
