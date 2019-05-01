package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/* #nosec */
const githubToken = "GITHUB_TOKEN"
const githubOwner = "GITHUB_OWNER"
const githubRepository = "GITHUB_REPOSITORY"

var rootCmd = &cobra.Command{
	Use:   "versem",
	Short: "Semver manager",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		newMessageHandler().errorFatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig(newMessageHandler()))
}

func initConfig(msgHandler messageHandler) func() {
	return func() {
		viper.AutomaticEnv()
		for _, key := range []string{
			githubOwner,
			githubRepository,
			githubToken,
		} {
			if !viper.IsSet(key) {
				msgHandler.errorFatalStr("missing environment variable : %s", key)
			}
		}
	}
}
