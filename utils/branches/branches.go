package branches

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
)

func PromptForBranchesAndReturnSelection(includeTrunk bool) (string, error) {
	branchesWithNames, err := engine.GetAllBranchesWithNames()
	if err != nil {
		return "", err
	}

	if includeTrunk {
		trunkName, err := engine.ReadTrunkName()
		if err != nil {
			return "", err
		}
		branchesWithNames = append(branchesWithNames, engine.BranchWithName{Name: trunkName, PRname: ""})
	}

	items := make([]string, len(branchesWithNames))
	for i, branch := range branchesWithNames {
		if branch.PRname != "" {
			items[i] = branch.PRname
		} else {
			items[i] = branch.Name
		}
	}

	prompt := promptui.Select{
		Label:             "Select next branch",
		Items:             items,
		StartInSearchMode: true,
		Searcher: func(input string, index int) bool {
			return strings.Contains(items[index], input)
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	branchNameToCheckout := result
	for _, branch := range branchesWithNames {
		if branch.PRname == result {
			branchNameToCheckout = branch.Name
			break
		}
	}

	return branchNameToCheckout, nil
}
