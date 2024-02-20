package cmd

import (
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the stack with the remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		return git.FetchAndPullTrunk("main")
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}
