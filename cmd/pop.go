package cmd

import (
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Remove the current branch but keep the changes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pop()
	},
}

func pop() error {
	stack, err := engine.GetStackForCurrentBranch()
	if err != nil {
		return err
	}
	parentSha, err := getParentShaForStack(stack)
	if err != nil {
		return err
	}
	if err = cli.ExecuteCommandInTerminal("git", "reset", "--soft", parentSha); err != nil {
		return err
	}

	if err = engine.RemoveBranchFromStack(stack.Name); err != nil {
		return err
	}

	if _, err = git.CheckoutBranch(stack.GetParent().Name); err != nil {
		return err
	}

	return git.DeleteBranchForce(stack.Name)
}

func getParentShaForStack(stack *engine.Stack) (string, error) {
	parent := stack.GetParent()
	out, err := git.CheckoutBranch(parent.Name)
	if err != nil {
		return out, err
	}
	defer git.CheckoutBranch(stack.Name)

	return git.GetCurrentBranchCommitSha()
}

func init() {
	RootCmd.AddCommand(popCmd)
}
