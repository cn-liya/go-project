package cache

import (
	"strconv"
)

const (
	WechatTokenKey = "wx:tk" // 微信access_token
)

const (
	keyUserToken = "utk:"  // +token
	keyUserInfo  = "user:" // +uid
)

func UserTokenKey(token string) string { return keyUserToken + token }

type UserToken struct {
	ID         int    `json:"i"`
	Openid     string `json:"o"`
	Unionid    string `json:"u"`
	SessionKey string `json:"s"`
}

func UserInfoKey(id int) string { return keyUserInfo + strconv.Itoa(id) } // entity.User
