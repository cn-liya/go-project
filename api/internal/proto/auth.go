package proto

type UserToken struct {
	ID      int    `json:"i"`
	Openid  string `json:"o"`
	Unionid string `json:"u"`
}
