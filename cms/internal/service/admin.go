package service

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"project/cms/internal/acl"
	"project/cms/internal/proto"
	"project/model"
	"strconv"
	"time"
)

const (
	keyAdminSSO   = "asso:" // +id [string] value = token
	keyAdminToken = "atk:"  // +token [string] value = JSON(proto.AdminToken)
	tokenTTL      = 2 * time.Hour
)

func ssoKey(id int) string {
	return keyAdminSSO + strconv.Itoa(id)
}

func (s *Service) TakeAdminByUsername(ctx context.Context, username string) (*acl.Admin, error) {
	var admin acl.Admin
	err := s.mysql.WithContext(ctx).Where("username = ?", username).Take(&admin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &admin, nil
}

func (s *Service) FindAdminByID(ctx context.Context, id int) (*acl.Admin, error) {
	var admin acl.Admin
	err := s.mysql.WithContext(ctx).Where("id = ?", id).Take(&admin).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &admin, nil
}

func (s *Service) SetAdminToken(ctx context.Context, admin *acl.Admin) (string, error) {
	sso := ssoKey(admin.ID)
	oldToken, err := s.redis.Get(ctx, sso).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(admin.Username))
	h.Write(uuid.NewV4().Bytes())
	newToken := base32.StdEncoding.EncodeToString(h.Sum(nil))
	b, _ := json.Marshal(&proto.AdminToken{
		ID:        admin.ID,
		Username:  admin.Username,
		Authority: admin.Authority,
	})
	tx := s.redis.TxPipeline()
	if oldToken != "" {
		tx.Del(ctx, keyAdminToken+oldToken)
	}
	tx.Set(ctx, keyAdminToken+newToken, b, tokenTTL)
	tx.Set(ctx, sso, newToken, tokenTTL)
	_, err = tx.Exec(ctx)
	return newToken, err
}

func (s *Service) GetAdminToken(ctx context.Context, token string) (*proto.AdminToken, error) {
	b, err := s.redis.Get(ctx, keyAdminToken+token).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var account proto.AdminToken
	if len(b) > 0 {
		_ = json.Unmarshal(b, &account)
	}
	return &account, nil
}

func (s *Service) DelAdminToken(ctx context.Context, token string) error {
	return s.redis.Del(ctx, keyAdminToken+token).Err()
}

func (s *Service) PaginateAdmin(ctx context.Context, p *proto.Pagination) (total int64, list []*acl.Admin, err error) {
	query := s.mysql.WithContext(ctx).Table((&acl.Admin{}).TableName()).Where("username <> ?", acl.SuperUser)
	err = query.Count(&total).Error
	if err != nil || total == 0 || p.Offset >= int(total) {
		return
	}
	err = query.Order("id").Limit(p.Limit).Offset(p.Offset).Find(&list).Error
	return
}

func (s *Service) CreateAdmin(ctx context.Context, data *acl.Admin) error {
	return s.mysql.WithContext(ctx).Create(data).Error
}

func (s *Service) UpdateAdmin(ctx context.Context, data *acl.Admin) error {
	opt := s.mysql.WithContext(ctx).Updates(data)
	if opt.Error != nil {
		return opt.Error
	}
	if opt.RowsAffected > 0 && (data.Status == model.StatusOff || data.Authority != nil) {
		token, err := s.redis.Get(ctx, ssoKey(data.ID)).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if token != "" {
			return s.redis.Del(ctx, keyAdminToken+token).Err()
		}
	}
	return nil
}
