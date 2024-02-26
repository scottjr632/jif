package cmd

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between parent and this branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}
		if stack.IsTrunk {
			color.Yellow("You are on the trunk branch, there is no parent to compare to")
			return nil
		}
		parentStack := stack.GetParent()
		return cli.ExecuteCommandInTerminal("git", "diff", parentStack.Name)
	},
}

func init() {
	RootCmd.AddCommand(diffCmd)
}
