package cmd

import (
	"sync"

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

		if err = submitForChildren(stack, stack.GetParent()); err != nil {
			return err
		}

		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}
		if err = engine.SyncGithubWithLocal(trunk); err != nil {
			return err
		}
		addCommentsForStack(stack)
		return nil
	},
}

func addCommentsForStack(stack *engine.Stack) error {
	stacks := engine.GetStacksInStack(stack)
	wg := sync.WaitGroup{}
	for _, stack := range stacks {
		if stack.IsTrunk || stack.PRNumber == "" {
			continue
		}

		wg.Add(1)
		go func(stack *engine.Stack) {
			defer wg.Done()
			if err := getAndUpdateBodyForPR(stack); err != nil {
				color.Red("Error updating body for PR: %s", err)
			}
		}(stack)
	}
	wg.Wait()
	return nil
}

func getAndUpdateBodyForPR(stack *engine.Stack) error {
	// get the current body
	body, err := gh.GetBodyForPR(stack.PRNumber)
	if err != nil {
		return err
	}

	bodyWithoutComment := engine.GetStringWithoutStackComment(body)

	comment := engine.GetStackForCommentByStack(stack)
	bodyWithoutComment += "\r\n" + comment
	if err = gh.UpdateBodyForPR(stack.PRNumber, bodyWithoutComment); err != nil {
		return err
	}
	return nil
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

func submitForChildren(currentStack *engine.Stack, parentStack *engine.Stack) error {
	git.CheckoutBranch(currentStack.Name)
	exists, err := gh.DoesPRExist(parentStack.Name, currentStack.Name)
	if err != nil {
		color.Red("Error checking if PR exists: %s", err)
		return err
	}

	if !exists {
		return nil
	}

	color.Yellow("PR already exists for %s to %s, force pushing", currentStack.Name, parentStack.Name)
	if err = git.GitPushForce(currentStack.Name); err != nil {
		color.Red("Error force pushing: %s", err)
	}

	for _, child := range currentStack.Children {
		child, err := engine.GetStackByID(child)
		if err != nil {
			color.Red("Error getting child stack: %s", err)
			continue
		}

		_, err = git.RebaseBranchOnto(child.Name, currentStack.Name, git.RebaseOptions{GoBackToPreviousBranch: true})
		if err != nil {
			return err
		}
		submitForChildren(child, currentStack)
	}
	return nil
}

func init() {
	prCommand.AddCommand(submitCmd)
	prCommand.AddCommand(viewCmd)
}
