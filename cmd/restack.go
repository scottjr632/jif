package cmd

import (
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var restackCmd = &cobra.Command{
	Use:   "restack",
	Short: "Restack the branches in the stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}
		return engine.RestackOntoParent(stack)
	},
}

func init() {
	RootCmd.AddCommand(restackCmd)
}
