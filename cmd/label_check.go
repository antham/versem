package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/antham/versem/github"

	"github.com/spf13/cobra"
)

var labelCheckCmd = &cobra.Command{
	Use:   "check [commitSha|pullRequestId]",
	Short: "Ensure a pull request contains a semver label",
	Run:   setupLabelCheckCmdFunc(labelCheck),
}

func init() {
	labelCmd.AddCommand(labelCheckCmd)
}

func setupLabelCheckCmdFunc(f func(messageHandler, semverService, *cobra.Command, []string)) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		msgHandler := newMessageHandler()
		f(msgHandler, getSemverLabelService(), cmd, args)
	}
}

func labelCheck(msgHandler messageHandler, semverService semverService, cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		msgHandler.errorFatalStr("provide a pull request number or full commit sha as first argument")
	}

	var version github.Version
	var err error
	var successMsg = "%s semver version found"

	switch {
	case regexp.MustCompile("[0-9a-f]{40}").MatchString(args[0]):
		version, err = semverService.GetFromCommit(args[0])
		if version == github.NORELEASE {
			successMsg += fmt.Sprintf(" or commit %s is not tied to a pull request", args[0])
		}
	case regexp.MustCompile("[0-9]+").MatchString(args[0]):
		n, _ := strconv.Atoi(args[0])
		version, err = semverService.GetFromPullRequest(n)
	default:
		err = fmt.Errorf("%s is not a valid number, nor a valid commit sha", args[0])
	}

	if err != nil {
		msgHandler.errorFatal(err)
	}

	msgHandler.success(successMsg, strings.ToLower(version.String()))
}
