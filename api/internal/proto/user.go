package proto

type CodeArgs struct {
	Code string `json:"code" binding:"required"`
}

type LoginResp struct {
	Token   string `json:"token"`
	Openid  string `json:"openid"`
	Unionid string `json:"unionid"`
}

type WechatPhoneResp struct {
	PhoneNumber string `json:"phone_number"`
}

type UserProfileArgs struct {
	Nickname  string `json:"nickname" binding:"required"`
	AvatarURL string `json:"avatar_url" binding:"required"`
}

type UserInfoResp struct {
	PhoneNumber string `json:"phone_number"`
	Nickname    string `json:"nickname"`
	AvatarURL   string `json:"avatar_url"`
}
