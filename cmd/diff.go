package cmd

import (
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
			cli.ExecuteCommandInTerminal("git", "diff")
			return nil
		}
		parentStack := stack.GetParent()
		return cli.ExecuteCommandInTerminal("git", "diff", parentStack.Name)
	},
}

func init() {
	RootCmd.AddCommand(diffCmd)
}
