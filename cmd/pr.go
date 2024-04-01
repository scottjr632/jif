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
		currentBranchStack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}

		parentStack, err := engine.GetStackByID(currentBranchStack.Parent)
		if err != nil {
			return err
		}

		if err = submitForParents(parentStack); err != nil {
			return err
		}
		// submit for self
		if err = createPRForStack(currentBranchStack); err != nil {
			return err
		}

		for _, childID := range currentBranchStack.Children {
			childStack, err := engine.GetStackByID(childID)
			if err != nil {
				return err
			}

			// will do this for each child
			if err = submitForSelfAndChildrenIfPRExists(childStack); err != nil {
				return err
			}
		}

		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}
		if err = engine.SyncGithubWithLocal(trunk); err != nil {
			return err
		}
		addCommentsForStack(currentBranchStack)
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

func createPRForStack(currentStack *engine.Stack) error {
	parentStackID := currentStack.Parent
	parentStack, err := engine.GetStackByID(parentStackID)
	if err != nil {
		return err
	}

	git.CheckoutBranch(currentStack.Name)
	if err := git.PushCurrentBranchToRemoteIfNotExists(); err != nil {
		color.Red("Error pushing current branch to remote: %s", err)
		return err
	}
	if exists, err := gh.DoesPRExist(parentStack.Name, currentStack.Name); err != nil {
		color.Red("Error checking if PR exists: %s", err)
		return err
	} else if exists {
		color.Yellow("PR already exists for %s to %s, force pushing with lease", currentStack.Name, parentStack.Name)
		return git.GitPushForce(currentStack.Name)
	}
	color.Green("Creating PR for %s to %s", currentStack.Name, parentStack.Name)
	return gh.CreatePR(parentStack.Name, currentStack.Name)
}

func submitForParents(currentStack *engine.Stack) error {
	if currentStack.IsTrunk || currentStack.Parent == 0 {
		return nil
	}

	parentStackID := currentStack.Parent
	parentStack, err := engine.GetStackByID(parentStackID)
	if err != nil {
		return err
	}

	if parentStack.Parent != 0 && !parentStack.IsTrunk {
		if err = submitForParents(parentStack); err != nil {
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
		color.Yellow("PR already exists for %s to %s, force pushing with lease", currentStack.Name, parentStack.Name)
		return git.GitPushForce(currentStack.Name)
	}
	color.Green("Creating PR for %s to %s", currentStack.Name, parentStack.Name)
	return gh.CreatePR(parentStack.Name, currentStack.Name)
}

func submitForSelfAndChildrenIfPRExists(currentStack *engine.Stack) error {
	parentStack, err := engine.GetStackByID(currentStack.Parent)
	if err != nil {
		return err
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
	for _, childID := range currentStack.Children {
		childStack, err := engine.GetStackByID(childID)
		if err != nil {
			return err
		}
		if err = submitForSelfAndChildrenIfPRExists(childStack); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	prCommand.AddCommand(submitCmd)
	prCommand.AddCommand(viewCmd)
}
