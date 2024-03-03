package cmd

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var revsCmd = &cobra.Command{
	Use:   "revs",
	Short: "List all revisions in the current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}
		color.Green("Revisions for %s:", stack.Name)
		revsLen := len(stack.Revisions)
		if revsLen == 0 {
			color.Yellow("no revisions found")
			return nil
		}

		for i := revsLen - 1; i >= 0; i-- {
			color.White("%d. %s", i+1, stack.Revisions[i])
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(revsCmd)
}
