package cmd

import (
	"log"

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

		return submitForParent(stack, stack.Parent)
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

	if err = git.PushCurrentBranchToRemoteIfNotExists(); err != nil {
		log.Println("Error pushing current branch to remote:", err)
		return err
	}
	if exists, err := gh.DoesPRExist(parentStack.Name, currentStack.Name); err != nil {
		log.Println("Error checking if PR exists:", err)
		return err
	} else if exists {
		log.Println("exists force pushing for current stack", currentStack.Name)
		return git.GitPushForce(currentStack.Name)
	}
	return gh.CreatePR(parentStack.Name, currentStack.Name)
}

func init() {
	prCommand.AddCommand(submitCmd)
}
