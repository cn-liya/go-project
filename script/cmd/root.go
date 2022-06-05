package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"project/script/internal/service"
)

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Println("start ...")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		log.Println("stop ...")
	},
}

func init() {
	rootCmd.AddCommand(cronjobCmd)
	cobra.OnInitialize(func() {
		service.Initialize()
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
