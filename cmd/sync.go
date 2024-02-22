package cmd

import (
	"log"

	"github.com/fatih/color"
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
		color.Green("Syncing the stack with the remote")
		mergedPRs, err := fetchAndPullTrunkWhileGettingMerged(trunkName)
		if err != nil {
			return err
		}
		log.Println("Merged PRs:", mergedPRs)
		color.Green("Closing merged PRs")
		closeMergedPRs(mergedPRs)
		return nil
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
	go func() {
		fpChan <- git.FetchAndPullTrunk(trunkName)
	}()
	go func() {
		prs, err := gh.GetMergedPRs()
		mergedChan <- MergedPRResult{prs, err}
	}()
	fpErr := <-fpChan
	mergedResult := <-mergedChan
	if fpErr != nil {
		return nil, fpErr
	}
	return mergedResult.prs, mergedResult.err
}

func closeMergedPRs(prs []gh.PRState) error {
	for _, pr := range prs {
		color.Green("Closing PR:", pr.Title)
		engine.RemoveBranchFromStack(pr.Branch)
	}
	return engine.Save()
}
