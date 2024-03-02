package engine

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/config"
	"github.com/scottjr632/sequoia/internal/git"
)

type ErrNoStack struct {
	BranchName string
}

func (e ErrNoStack) Error() string {
	return "no stack found for branch " + e.BranchName
}

const (
	CommitTypeAll   = "submit all changes (--all)"
	CommitTypePatch = "some changes (--patch)"
	CommitTypeAbort = "Abort changes"
)

type CommitOptions struct {
	AutoStage bool
}

func promptIfUnstagedFiles() error {
	if err := git.EnsureStagedFiles(); err != nil {
		if git.IsNoStagedFilesError(err) {
			color.Yellow("no staged files found")
		}
		prompts := promptui.Select{
			Label: "No staged files found. Would you like to add any?",
			Items: []string{CommitTypeAll, CommitTypePatch, CommitTypeAbort},
		}

		_, result, err := prompts.Run()
		if err != nil {
			return err
		}

		switch result {
		case CommitTypeAll:
			err := git.PromptToAddAll()
			if err != nil {
				return err
			}
		case CommitTypePatch:
			err := git.PromptToPatch()
			if err != nil {
				return err
			}
		case CommitTypeAbort:
			return fmt.Errorf("aborted")
		}
	}
	return nil
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

	if err = promptIfUnstagedFiles(); err != nil {
		return err
	}

	c, err := config.LoadPersistenConfig()
	if err != nil {
		return err
	}

	if c.BranchPrefix != "" && !strings.HasSuffix(c.BranchPrefix, "/") {
		c.BranchPrefix = c.BranchPrefix + "/"
	}

	// create a new stack for the new branch
	t := time.Now().Format("01-02")
	branchName := c.BranchPrefix + t + "-" + git.FormatBranchNameFromCommit(message)

	err = git.CreateAndCheckoutBranch(branchName)
	if err != nil {
		log.Println("failed to checkout branch", branchName, err)
		return err
	}

	err = git.Commit(message)
	if err != nil {
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
		color.Red("Failed to save stack. Please try again.", err)
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

	if err = promptIfUnstagedFiles(); err != nil {
		return err
	}

	err = git.AmendCommit()
	if err != nil {
		return err
	}

	stack, err = GetStackForCurrentBranch()
	if err != nil {
		return err
	}
	return RestackChildren(stack)
}
