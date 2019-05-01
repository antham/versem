package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [commitSha|pullRequestId]",
	Short: "Ensure a pull request contains a semver label",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("check called")
	},
}

func init() {
	labelCmd.AddCommand(checkCmd)
}
