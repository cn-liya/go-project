package service

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"project/model"
	"project/pkg/logger"
	"project/pkg/mq"
	"project/pkg/util/types"
)

type Message struct {
	consumer *nsq.Consumer
}

func NewMessage() *Message {
	s := &Message{}
	s.consumer = mq.NewNsqConsumer(config.Nsq.Consumer, model.TopicExample, "default", 4, s.handle)
	return s
}

func (s *Message) Stop() {
	s.consumer.Stop()
}

func (s *Message) handle(msg *nsq.Message) error {
	reqid := string(msg.ID[:])
	var data model.MsgExample
	_ = json.Unmarshal(msg.Body, &data)
	_, l := logger.NewCtxLog(reqid, "Message", "handle", types.Int2Str(msg.Timestamp))
	l.Info("msg.body", msg.Body, &data)
	return nil
}
