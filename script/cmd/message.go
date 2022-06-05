package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"project/script/internal/service"
	"syscall"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "nsq消费消息示例",
	Run: func(cmd *cobra.Command, args []string) {
		s := service.NewMessage()
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-quit
		s.Stop()
	},
}

func init() {
	rootCmd.AddCommand(messageCmd)
}
