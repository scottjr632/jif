package cmd

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var prCommand = &cobra.Command{
	Use:   "pr",
	Short: "Interact with GitHub pull requests",
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit a pull request",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}

		if err = submitForParent(stack, stack.Parent); err != nil {
			return err
		}

		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}
		return engine.SyncGithubWithLocal(trunk)
	},
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "View a pull request on GitHub",
	Aliases: []string{"open", "o", "v", "web"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return gh.RunGHCmd("pr", "view", "--web")
	},
}

func submitForParent(currentStack *engine.Stack, parentStackID engine.StackID) error {
	parentStack, err := engine.GetStackByID(parentStackID)
	if err != nil {
		return err
	}

	if parentStack.Parent != 0 && !parentStack.IsTrunk {
		if err = submitForParent(parentStack, parentStack.Parent); err != nil {
			return err
		}
	}

	git.CheckoutBranch(currentStack.Name)
	if err = git.PushCurrentBranchToRemoteIfNotExists(); err != nil {
		color.Red("Error pushing current branch to remote: %s", err)
		return err
	}
	if exists, err := gh.DoesPRExist(parentStack.Name, currentStack.Name); err != nil {
		color.Red("Error checking if PR exists: %s", err)
		return err
	} else if exists {
		color.Yellow("PR already exists for %s to %s, force pushing", currentStack.Name, parentStack.Name)
		return git.GitPushForce(currentStack.Name)
	}
	color.Green("Creating PR for %s to %s", currentStack.Name, parentStack.Name)
	return gh.CreatePR(parentStack.Name, currentStack.Name)
}

func init() {
	prCommand.AddCommand(submitCmd)
	prCommand.AddCommand(viewCmd)
}
