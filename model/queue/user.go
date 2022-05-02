package queue

/*
CustomTP topic
CustomCH channel
CustomEX exchange
CustomQN queue-name
CustomCG consumer-group
Custom body-struct
*/

const DefaultCH = "default"

const AvatarToCdnTP = "avatar_to_cdn"

//const AvatarToCdnQN = "avatar_to_cdn"

type AvatarToCdn struct {
	ID        int    `json:"id"`
	Openid    string `json:"openid"`
	AvatarURL string `json:"avatar_url"`
}
