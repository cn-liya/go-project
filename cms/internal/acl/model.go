package acl

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const SuperUser = "admin"

type Admin struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Authority  Authority `json:"authority" gorm:"type:json"`
	Status     int8      `json:"status"`
	CreateTime time.Time `json:"create_time" gorm:"->"` // 只读
	UpdateTime time.Time `json:"update_time" gorm:"->"` // 只读
}

func (*Admin) TableName() string {
	return "admin"
}

func (a *Admin) BeforeSave(*gorm.DB) error {
	if a.Password != "" {
		a.Password = cryptoHash(a.Password)
	}
	return nil
}

type Authority map[string]int8

func (auth *Authority) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	return json.Unmarshal(b, auth) // receiver必须为指针
}

func (auth Authority) Value() (driver.Value, error) {
	if auth == nil {
		return []byte{'{', '}'}, nil
	}
	return json.Marshal(auth) // receiver不能为指针
}

func CheckPassword(input, crypt string) bool {
	return cryptoHash(input) == crypt
}

func cryptoHash(pwd string) string {
	h := sha256.Sum224([]byte(pwd))
	return base64.URLEncoding.EncodeToString(h[:24])
}
