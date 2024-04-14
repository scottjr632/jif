package cmd

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var unamendCmd = &cobra.Command{
	Use:   "unamend",
	Short: "Undo the last amend",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}

		rev := stack.PopRevision()
		if rev == "" {
			color.Yellow("no revisions to unamend")
			return nil
		}

		if err = cli.ExecuteCommandInTerminal("git", "reset", "--soft", rev); err != nil {
			return err
		}

		if len(stack.Children) > 0 {
			if _, err = git.StageAll(); err != nil {
				return nil
			}

			if err = cli.ExecuteCommandInTerminal("git", "stash"); err != nil {
				return err
			}
			defer func(stack *engine.Stack) {
				git.CheckoutBranch(stack.Name)
				cli.ExecuteCommandInTerminal("git", "stash", "pop")
			}(stack)
		}
		if err = engine.RestackOntoParent(stack); err != nil {
			return err
		}
		return engine.Save()
	},
}

func init() {
	RootCmd.AddCommand(unamendCmd)
}
