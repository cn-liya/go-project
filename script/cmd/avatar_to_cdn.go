package cmd

import (
	"github.com/spf13/cobra"
	"project/model/queue"
	"project/pkg/gnsq"
	"project/pkg/process"
	"project/script/internal/handler"
	"project/script/internal/service"
)

var avatarToCdnCmd = &cobra.Command{
	Use:   "avatar:to:cdn",
	Short: "从临时链接获取用户头像转存到CND",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.Mysql(&cfg.Service.Mysql)))
		consumer := gnsq.NewConsumer(cfg.Service.Nsq.Consumer, queue.AvatarToCdnTP, queue.DefaultCH, 4, h.AvatarToCdn)
		//conn := rabbitmq.NewConnect(&cfg.Service.Rabbitmq)
		//consumer := rabbitmq.NewConsumer(conn, queue.AvatarToCdnQN, 4, h.AvatarToCdn)
		process.Notify()
		consumer.Stop()
	},
}
