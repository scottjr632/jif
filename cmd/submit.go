package cmd

import (
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var baseSubmitCmd = &cobra.Command{
	Use:     "submit",
	Short:   "Submit a new PR to GitHub",
	Aliases: []string{"s"},
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}

		return submitForParent(stack, stack.Parent)
	},
}
