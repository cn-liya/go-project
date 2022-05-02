package entity

import (
	"encoding/base64"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
	"project/pkg/random"
	"time"
)

type AdminRole struct {
	ID         int               `json:"id"`
	Name       string            `json:"name"`
	Actions    JsonSlice[string] `json:"actions"`
	Menus      JsonSlice[string] `json:"menus"`
	CreateTime time.Time         `json:"create_time" gorm:"->"` // 只读
	UpdateTime time.Time         `json:"update_time" gorm:"->"` // 只读
}

func (*AdminRole) TableName() string {
	return "admin_role"
}

type AdminUser struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"-"`
	PwdExp     int64      `json:"-"`
	RoleID     int        `json:"role_id"`
	Status     int8       `json:"status"`
	CreateTime time.Time  `json:"create_time" gorm:"->"` // 只读
	UpdateTime time.Time  `json:"update_time" gorm:"->"` // 只读
	AdminRole  *AdminRole `json:"admin_role,omitempty" gorm:"foreignKey:ID;references:RoleID"`
}

func (*AdminUser) TableName() string {
	return "admin_user"
}

const (
	cost    = 4
	memory  = 1 << 12
	threads = 4
	keyLen  = 32
	saltLen = 13
	pwdTTL  = 7862400
)

func (a *AdminUser) BeforeSave(*gorm.DB) error {
	if a.Password != "" {
		salt := random.Bytes(saltLen)
		key := argon2.IDKey([]byte(a.Password), salt, cost, memory, threads, keyLen)
		a.Password = base64.URLEncoding.EncodeToString(append(key, salt...))
		a.PwdExp = time.Now().Unix() + pwdTTL
	} // 哈希密码和salt值合并存储到一个字段，前32字节为哈希密码，后13字节为salt值
	return nil
}

func (a *AdminUser) CheckPassword(input string) bool {
	b, _ := base64.URLEncoding.DecodeString(a.Password)
	if len(b) < keyLen {
		return false
	}
	key := argon2.IDKey([]byte(input), b[keyLen:], cost, memory, threads, keyLen)
	return string(key) == string(b[:keyLen])
}

func (a *AdminUser) IsSuper() bool { return a.RoleID == -1 }
