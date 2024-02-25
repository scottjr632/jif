package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the stack with the remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		trunkName, err := engine.ReadTrunkName()
		if err != nil {
			return err
		}
		color.Yellow("Syncing the stack with the remote")
		mergedPRs, err := fetchAndPullTrunkWhileGettingMerged(trunkName)
		if err != nil {
			return err
		}
		color.Yellow("🫧 closing merged or closed PRs")
		closeMergedPRs(mergedPRs)
		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}
		return engine.RestackChildren(trunk)
	},
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

type MergedPRResult struct {
	prs []gh.PRState
	err error
}

func fetchAndPullTrunkWhileGettingMerged(trunkName string) ([]gh.PRState, error) {
	if _, err := git.CheckoutBranch(trunkName); err != nil {
		return nil, err
	}

	fpChan := make(chan error)
	mergedChan := make(chan MergedPRResult)
	closedChan := make(chan MergedPRResult)
	go func() {
		fpChan <- git.FetchAndPullTrunk(trunkName)
	}()
	go func() {
		prs, err := gh.GetMergedPRs()
		mergedChan <- MergedPRResult{prs, err}
	}()
	go func() {
		prs, err := gh.GetClosedPRs()
		closedChan <- MergedPRResult{prs, err}
	}()
	fpErr := <-fpChan
	mergedResult := <-mergedChan
	closedResult := <-closedChan
	if fpErr != nil {
		return nil, fpErr
	}
	combined := append(mergedResult.prs, closedResult.prs...)
	return combined, mergedResult.err
}

func doesExist(name string, names []string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func closeMergedPRs(prs []gh.PRState) error {
	stacksNames, err := engine.GetAllStackNames()
	if err != nil {
		return err
	}

	for _, pr := range prs {
		if !doesExist(pr.Branch, stacksNames) {
			continue
		}

		fmt.Printf("%s was merged\n", pr.Branch)
		prompt := promptui.Select{
			HideHelp:  true,
			IsVimMode: true,
			Label:     "Remove from local stack?",
			Items:     []string{"yes", "no"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		if result == "no" {
			color.WhiteString("skipping...")
			continue
		}

		color.Green("closing...")
		if err = engine.RemoveBranchFromStack(pr.Branch); err != nil {
			color.Red("error removing branch from stack: %s", err)
		}
	}
	return engine.Save()
}
