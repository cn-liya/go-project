package service

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"project/api/internal/proto"
	"project/model"
	"project/pkg/logger"
	"strconv"
	"time"
)

func userKey(id int) string {
	return model.KeyUserInfo + strconv.Itoa(id)
}

func (s *Service) WechatLogin(ctx context.Context, code string) (*model.User, error) {
	resp, err := s.wechat.JsCode2Session(ctx, code)
	if err != nil {
		return nil, err
	}
	data := &model.User{
		Openid:  resp.Openid,
		Unionid: resp.Unionid,
	}
	if data.Openid != "" {
		err = s.mysql.WithContext(ctx).FirstOrCreate(data, "openid = ?", resp.Openid).Error
		if data.ID > 0 {
			b, _ := json.Marshal(data)
			if err := s.redis.Set(ctx, userKey(data.ID), b, time.Hour); err != nil {
				logger.FromContext(ctx).Error("redis.Set error", nil, err)
			}
		}
	} else {
		logger.FromContext(ctx).Warn("wechat.JsCode2Session fail", code, resp)
	}
	return data, err
}

func (s *Service) SetUserToken(ctx context.Context, user *model.User) (string, error) {
	h := sha1.New()
	h.Write([]byte(user.Openid))
	h.Write(uuid.NewV4().Bytes())
	token := base32.StdEncoding.EncodeToString(h.Sum(nil))
	b, _ := json.Marshal(&proto.UserToken{
		ID:      user.ID,
		Openid:  user.Openid,
		Unionid: user.Unionid,
	})
	err := s.redis.Set(ctx, model.KeyUserToken+token, b, 2*time.Hour).Err()
	return token, err
}

func (s *Service) GetUserToken(ctx context.Context, token string) (*proto.UserToken, error) {
	b, err := s.redis.Get(ctx, model.KeyUserToken+token).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var account proto.UserToken
	if len(b) > 0 {
		_ = json.Unmarshal(b, &account)
	}
	return &account, nil
}

func (s *Service) QueryUser(ctx context.Context, id int) (*model.User, error) {
	key := userKey(id)
	b, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var res model.User
	if len(b) > 0 {
		err = json.Unmarshal(b, &res)
		return &res, err
	}

	err = s.mysql.WithContext(ctx).Where("id = ?", id).Take(&res).Error
	if err != nil {
		return nil, err
	}
	b, _ = json.Marshal(res)
	if err := s.redis.Set(ctx, key, b, time.Hour); err != nil {
		logger.FromContext(ctx).Error("redis.Set error", key, err)
	}
	return &res, nil
}

func (s *Service) WechatPhone(ctx context.Context, code string) (string, error) {
	resp, err := s.wechat.GetUserPhoneNumber(ctx, code)
	if err != nil {
		return "", err
	}
	if resp.PhoneInfo != nil {
		return resp.PhoneInfo.PhoneNumber, nil
	}
	logger.FromContext(ctx).Warn("wechat.GetUserPhoneNumber fail", code, resp)
	return "", nil
}

func (s *Service) UpdateUser(ctx context.Context, data *model.User) error {
	opt := s.mysql.WithContext(ctx).Updates(data) // gorm根据ID更新指定非零值字段
	if opt.Error != nil {
		return opt.Error
	}
	if opt.RowsAffected > 0 {
		return s.redis.Del(ctx, userKey(data.ID)).Err()
	}
	return nil
}
