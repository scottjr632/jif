package cmd

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/utils/branches"
	"github.com/spf13/cobra"
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase [onto]",
	Short: "Rebase the current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		var onto string
		if len(args) > 0 {
			onto = args[0]
		} else {
			selectedBranch, err := branches.PromptForBranchesAndReturnSelection(true)
			if err != nil {
				return err
			}
			onto = selectedBranch
		}

		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			color.Red("Unable to get stack for the current branch: %s", err)
			return err
		}

		parentStack, err := engine.GetStackByID(stack.Parent)
		if err != nil {
			color.Red("Unable to get parent stack for the current branch: %s", err)
			return err
		}

		ontoStack, err := engine.GetStackForBranch(onto)
		if err != nil {
			color.Red("Unable to get stack for the onto branch: %s", err)
			return err
		}

		if _, err := cli.ExecuteCmd("git", "rebase", "--onto", ontoStack.Name, parentStack.Name, stack.Name); err != nil {
			err = cli.ExecuteCommandInTerminal("git", "add", ".")
			if err != nil {
				color.Red("Error rebasing: %s", err)
				return err
			}

			_, err = cli.ExecuteCmdWithEnviron("git", []string{"GIT_EDITOR=true"}, "rebase", "--continue")
			if err != nil {
				color.Red("Error rebasing: %s", err)
				return err
			}
		}

		return stack.Rebase(ontoStack)
	},
}

func init() {
	RootCmd.AddCommand(rebaseCmd)
}
