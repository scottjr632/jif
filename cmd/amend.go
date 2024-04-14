package cmd

import (
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var amendCmd = &cobra.Command{
	Use:     "am",
	Aliases: []string{"amend"},
	Short:   "Amend the last commit",
	Long:    `Amend the last commit with the staged changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stageAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}
		if stageAll {
			_, err = git.StageAll()
			if err != nil {
				return err
			}
		}

		return engine.AmendCommit(engine.CommitOptions{
			AutoStage: true,
		})
	},
}

func init() {
	RootCmd.AddCommand(amendCmd)

	amendCmd.Flags().BoolP("all", "a", false, "Stage all files")
}
