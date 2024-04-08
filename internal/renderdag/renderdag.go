package renderdag

import (
	"fmt"
	"github.com/fatih/color"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/xlab/treeprint"
)

func createNameShort(stack *engine.Stack) string {
	if stack.IsTrunk {
		return stack.Name
	}

	if stack.PRName != "" {
		return stack.PRName
	}
	return stack.Name
}

func createNameLong(stack *engine.Stack) string {
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

func addChild(parent treeprint.Tree, child *engine.Stack, currentBranchName string, createNameFn func(stack *engine.Stack) string) {
	isCurrentBranch := child.Name == currentBranchName
	var branchName string
	currentName := createNameFn(child)
	if isCurrentBranch {
		currentString := color.GreenString("(current)")
		branchName = fmt.Sprintf("%s %s", currentString, currentName)
	} else {
		branchName = currentName
	}

	if child.NeedsRestack {
		branchName = fmt.Sprintf("%s (needs restack)", branchName)
	}

	branch := parent.AddBranch(branchName)
	for _, c := range child.Children {
		childStack, err := engine.GetStackByID(c)
		if err != nil {
			continue
		}
		addChild(branch, childStack, currentBranchName, createNameFn)
	}
}

func RenderDag(trunk *engine.Stack) {
	tree := treeprint.New()

	currentBranch, err := git.GetCurrentBranchName()
	if err != nil {
		panic(err)
	}

	addChild(tree, trunk, currentBranch, createNameLong)

	fmt.Println(tree.String())
}

func RenderDagShort(trunk *engine.Stack) {
	tree := treeprint.New()

	currentBranch, err := git.GetCurrentBranchName()
	if err != nil {
		panic(err)
	}

	addChild(tree, trunk, currentBranch, createNameShort)

	fmt.Println(tree.String())
}
