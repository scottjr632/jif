package renderdag

import (
	"fmt"
	"github.com/fatih/color"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/xlab/treeprint"
)

func createName(stack *engine.Stack) string {
	if stack.IsTrunk {
		return stack.Name
	}

	var builder string
	if stack.PRName != "" {
		builder = stack.PRName
	} else {
		builder = stack.Name
	}

	if stack.PRStatus != engine.PRStatusNone {
		builder = fmt.Sprintf("%s (%s)", builder, stack.PRStatus)
	}

	if stack.PRNumber != "" {
		builder = fmt.Sprintf("%s #%s", builder, stack.PRNumber)
	}
	if stack.PRLink != "" {
		builder = fmt.Sprintf("%s\n  %s", builder, stack.PRLink)
	}
	return builder
}

func addChild(parent treeprint.Tree, child *engine.Stack, currentBranchName string) {
	isCurrentBranch := child.Name == currentBranchName
	var branchName string
	currentName := createName(child)
	if isCurrentBranch {
		currentString := color.GreenString("(current)")
		branchName = fmt.Sprintf("%s %s", currentString, currentName)
	} else {
		branchName = currentName
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
