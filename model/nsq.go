package model

// 定义队列的topic和数据结构

const (
	TopicExample = "example"
)

type MsgExample struct {
	UUID   string `json:"uuid"`
	Number int    `json:"number"`
}
