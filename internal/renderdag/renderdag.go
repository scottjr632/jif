package renderdag

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/xlab/treeprint"
)

func addChild(parent treeprint.Tree, child *engine.Stack, currentBranchName string) {
	isCurrentBranch := child.Name == currentBranchName
	var branchName string
	if isCurrentBranch {
		branchName = fmt.Sprintf("(current) %s", child.Name)
	} else {
		branchName = child.Name
	}
	branch := parent.AddBranch(branchName)
	for _, c := range child.Children {
		childStack, err := engine.GetStackByID(c)
		if err != nil {
			continue
		}
		addChild(branch, childStack, currentBranchName)
	}
}

func RenderDag(trunk *engine.Stack) {
	tree := treeprint.New()

	currentBranch, err := git.GetCurrentBranchName()
	if err != nil {
		panic(err)
	}

	addChild(tree, trunk, currentBranch)

	fmt.Println(tree.String())
}
