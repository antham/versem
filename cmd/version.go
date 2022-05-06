package cmd

import (
	"github.com/spf13/cobra"
)

var appVersion = ""

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run:   setupVersionCmdFunc(version),
}

func setupVersionCmdFunc(f func(messageHandler, *cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		msgHandler := newMessageHandler()
		f(msgHandler, cmd, args)
	}
}

func version(msgHandler messageHandler, cmd *cobra.Command, args []string) {
	msgHandler.success(appVersion)
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
