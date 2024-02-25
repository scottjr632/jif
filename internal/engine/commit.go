package engine

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/git"
)

type ErrNoStack struct {
	BranchName string
}

func (e ErrNoStack) Error() string {
	return "no stack found for branch " + e.BranchName
}

type CommitOptions struct {
	AutoStage bool
}

func CommitWithNewBranch(message string, options CommitOptions) error {
	stack, err := GetStackForCurrentBranch()
	if err != nil {
		return err
	}

	if stack == nil {
		currentBranch, _ := git.GetCurrentBranchName()
		return ErrNoStack{BranchName: currentBranch}
	}

	if err := git.EnsureStagedFiles(); err != nil {
		if git.IsNoStagedFilesError(err) {
			color.Yellow("no staged files found")
		}
		prompts := promptui.Select{
			Label: "No staged files found. Would you like to add any?",
			Items: []string{"submit all changes (--all)", "(some changes) --patch", "Abort changes"},
		}

		_, result, err := prompts.Run()
		if err != nil {
			return err
		}

		switch result {
		case "submit all changes (--all)":
			err := git.PromptToAddAll()
			if err != nil {
				return err
			}
		case "(some changes) --patch":
			err := git.PromptToPatch()
			if err != nil {
				return err
			}
		case "Abort changes":
			return fmt.Errorf("aborted")
		}
	}

	// create a new stack for the new branch
	branchName := git.FormatBranchNameFromCommit(message)
	_, err = git.CreateAndCheckoutBranch(branchName)
	if err != nil {
		log.Println("failed to checkout branch", branchName, err)
		return err
	}

	out, err := git.Commit(message)
	if err != nil {
		log.Println("failed to commit", message, err, string(out))
		return err
	}

	sha, err := git.GetCurrentBranchCommitSha()
	if err != nil {
		log.Println("failed to get sha", err)
		return err
	}

	newStack := NewStack(branchName, false, false, sha, stack.ID)
	stack.AddChild(newStack.ID)
	err = Save()
	if err != nil {
		log.Println("failed to save stack", err)
		return err
	}
	return nil
}

func AmendCommit(options CommitOptions) error {
	// ensure we're not amending trunk
	stack, err := GetStackForCurrentBranch()
	if err != nil {
		return err
	}
	if stack.IsTrunk {
		return fmt.Errorf("cannot amend trunk")
	}

	if err := git.EnsureStagedFiles(); err != nil {
		if git.IsNoStagedFilesError(err) {
			log.Println("no staged files found")
		}
		if options.AutoStage {
			log.Println("auto staging files")
			if _, err := git.StageAll(); err != nil {
				log.Println("failed to stage files", err)
				return err
			}
		} else {
			return err
		}
	}

	_, err = git.AmendCommit()
	if err != nil {
		return err
	}

	stack, err = GetStackForCurrentBranch()
	if err != nil {
		return err
	}
	return RestackChildren(stack)
}
