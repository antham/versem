package cmd

import (
	"regexp"

	"github.com/spf13/cobra"
)

var releaseCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create semver release on a repository",
	Run:   setupReleaseCreateCmdFunc(releaseCreate),
}

func init() {
	releaseCmd.AddCommand(releaseCreateCmd)
}

func setupReleaseCreateCmdFunc(f func(messageHandler, semverService, releaseService, *cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		msgHandler := newMessageHandler()
		f(msgHandler, getSemverLabelService(), getReleaseService(), cmd, args)
	}
}

func releaseCreate(msgHandler messageHandler, semverService semverService, releaseService releaseService, cmd *cobra.Command, args []string) {
	if len(args) != 1 || !regexp.MustCompile("[0-9a-f]{7,}").MatchString(args[0]) {
		msgHandler.errorFatalStr("provide a commit sha as first argument")
	}

	version, err := semverService.GetFromCommit(args[0])
	if err != nil {
		msgHandler.errorFatal(err)
	}

	tag, err := releaseService.CreateNext(version, args[0])
	if err != nil {
		msgHandler.errorFatal(err)
	}

	msgHandler.success("tag %s created", tag)
}
