package engine

import "github.com/scottjr632/sequoia/utils/stack"

type BranchWithName struct {
	Name   string
	PRname string
}

func GetAllBranchNames() ([]BranchWithName, error) {
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

		branches = append(branches, BranchWithName{
			Name:   child.Name,
			PRname: child.PRName,
		})

		stack.PushMany(child.Children)
	}
	return branches, nil
}
