package cmd

import (
	"github.com/spf13/cobra"
)

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create semver label on a repository",
	Run:   setupLabelCreateCmdFunc(labelCreate),
}

func init() {
	labelCmd.AddCommand(labelCreateCmd)
}

func setupLabelCreateCmdFunc(f func(messageHandler, semverService, *cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		msgHandler := newMessageHandler()
		f(msgHandler, getSemverLabelService(), cmd, args)
	}
}

func labelCreate(msgHandler messageHandler, semverService semverService, cmd *cobra.Command, args []string) {
	if err := semverService.CreateList(); err != nil {
		msgHandler.errorFatal(err)
	}

	msgHandler.success("semver labels created")
}
