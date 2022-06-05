package model

type User struct {
	ID          int    `json:"id"`
	Openid      string `json:"openid"`
	Unionid     string `json:"unionid"`
	PhoneNumber string `json:"phone_number"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
}

func (*User) TableName() string {
	return "user"
}
