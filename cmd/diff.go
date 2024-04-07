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
		parentStack := stack.GetParent()
		if parentStack == nil {
			cli.ExecuteCommandInTerminal("git", "diff")
			return nil
		}

		return cli.ExecuteCommandInTerminal("git", "diff", parentStack.Name)
	},
}

func init() {
	RootCmd.AddCommand(diffCmd)
}
