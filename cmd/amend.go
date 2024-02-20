package cmd

import (
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var amendCmd = &cobra.Command{
	Use:     "am",
	Aliases: []string{"amend"},
	Short:   "Amend the last commit",
	Long:    `Amend the last commit with the staged changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return engine.AmendCommit(engine.CommitOptions{
			AutoStage: true,
		})
	},
}
