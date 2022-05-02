package cmd

import (
	"github.com/spf13/cobra"
	"project/pkg/process"
	"project/script/internal/handler"
	"project/script/internal/service"
	"time"
)

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh:token",
	Short: "刷新小程序AccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		h := handler.New(&cfg.Handler, service.New(service.NewRedis(&cfg.Service.Redis)))
		stop := make(chan struct{})
		done := make(chan struct{})
		go func() {
			tk := time.Tick(2 * time.Minute)
			for {
				select {
				case <-stop:
					close(done)
					return
				case <-tk:
					h.WechatServerToken()
				}
			}
		}()
		process.Notify()
		close(stop)
		<-done
	},
}
