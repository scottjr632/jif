package engine

import (
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/scottjr632/sequoia/utils/stack"
)

type BranchWithName struct {
	Name          string
	PRname        string
	CommitMessage string
}

func GetAllBranchesWithNames() ([]BranchWithName, error) {
	trunk, err := GetTrunk()
	if err != nil {
		return nil, err
	}
	branches := make([]BranchWithName, 0)
	stack := stack.New(trunk.Children...)
	for !stack.IsEmpty() {
		childID, err := stack.Pop()
		if err != nil {
			return nil, err
		}

		child, err := GetStackByID(*childID)
		if err != nil {
			return nil, err
		}

		if child.CommitMessage == "" {
			commitMessage, err := git.GetLatestCommitMessage(child.Name)
			if err == nil {
				child.CommitMessage = commitMessage
			}
		}

		branches = append(branches, BranchWithName{
			Name:          child.Name,
			PRname:        child.PRName,
			CommitMessage: child.CommitMessage,
		})

		stack.PushMany(child.Children)
	}
	return branches, nil
}
