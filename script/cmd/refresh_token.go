package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"project/script/internal/service"
	"syscall"
	"time"
)

var refreshTokenCmd = &cobra.Command{
	Use:   "refresh:token",
	Short: "刷新小程序AccessToken",
	Run: func(cmd *cobra.Command, args []string) {
		s := service.NewRefreshToken()
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
					s.WechatServerToken()
				}
			}
		}()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit
		close(stop)
		<-done
	},
}

func init() {
	rootCmd.AddCommand(refreshTokenCmd)
}
