package renderdag

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/xlab/treeprint"
)

func findMaybeNameInPRs(prs []gh.PRState, name string) string {
	for _, pr := range prs {
		if pr.Branch == name {
			return fmt.Sprintf("%s (%s)", pr.Title, pr.State)
		}
	}
	return name
}

func addChild(parent treeprint.Tree, child *engine.Stack, currentBranchName string, prs []gh.PRState) {
	isCurrentBranch := child.Name == currentBranchName
	var branchName string
	nameToUse := findMaybeNameInPRs(prs, child.Name)
	if isCurrentBranch {
		branchName = fmt.Sprintf("(current) %s", nameToUse)
	} else {
		branchName = nameToUse
	}
	branch := parent.AddBranch(branchName)
	for _, c := range child.Children {
		childStack, err := engine.GetStackByID(c)
		if err != nil {
			continue
		}
		addChild(branch, childStack, currentBranchName, prs)
	}
}

func RenderDag(trunk *engine.Stack) {
	tree := treeprint.New()
	openBranches, _ := gh.GetOpenPRs()

	currentBranch, err := git.GetCurrentBranchName()
	if err != nil {
		panic(err)
	}

	addChild(tree, trunk, currentBranch, openBranches)

	fmt.Println(tree.String())
}
