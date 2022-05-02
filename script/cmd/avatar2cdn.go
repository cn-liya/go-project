package cmd

import (
	"github.com/spf13/cobra"
	"project/model/queue"
	"project/pkg/gnsq"
	"project/pkg/process"
	"project/script/internal/handler"
	"project/script/internal/service"
)

var avatar2cdnCmd = &cobra.Command{
	Use:   "avatar2cdn",
	Short: "从临时链接获取用户头像转存到CND",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.NewMysql(&cfg.Service.Mysql)))
		consumer := gnsq.NewConsumer(cfg.Service.Nsq.Consumer, queue.AvatarToCdnQN, queue.DefaultCSM, 4, h.Handle)
		process.Notify()
		consumer.Stop()
	},
}
