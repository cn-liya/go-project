package queue

/*
SomeModelQN  队列名称： rabbitmq.name | nsq.topic | redis-stream.key
SomeModelCSM  消费者(组)： rabbitmq.consumer | nsq.channel | redis-stream.group
SomeModel  消息体结构
*/

const DefaultCSM = "default"

const AvatarToCdnQN = "avatar_to_cdn"

type AvatarToCdn struct {
	ID        int    `json:"id"`
	Openid    string `json:"openid"`
	AvatarURL string `json:"avatar_url"`
}
