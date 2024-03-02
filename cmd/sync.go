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
		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}

		color.Yellow("Syncing the stack with the remote...")

		errChan := make(chan error)
		go func() {
			errChan <- engine.SyncGithubWithLocal(trunk)
		}()

		mergedPRs, err := fetchAndPullTrunkWhileGettingMerged(trunkName)
		if err != nil {
			return err
		}

		// block for the sync from remote to finish, however, its not critical if it fails
		<-errChan

		fmt.Println("Checking merged or closed PRs...")
		closeMergedPRs(mergedPRs)

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
	// merged and closed PRs can overlap, so we need to dedupe them
	return dedupePRs(combined), mergedResult.err
}

func dedupePRs(prs []gh.PRState) []gh.PRState {
	seen := make(map[string]struct{})
	result := make([]gh.PRState, 0)
	for _, pr := range prs {
		if _, ok := seen[pr.Branch]; ok {
			continue
		}
		seen[pr.Branch] = struct{}{}
		result = append(result, pr)
	}
	return result
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

	fmt.Printf("\r\n\r\n")

	for _, pr := range prs {
		if !doesExist(pr.Branch, stacksNames) {
			continue
		}

		cyan := color.New(color.FgCyan).SprintFunc()
		label := fmt.Sprintf("%s was merged, remove it locally?", cyan(pr.Branch))
		prompt := promptui.Select{
			HideHelp:  true,
			IsVimMode: true,
			Label:     label,
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
		err = engine.RemoveBranchFromStack(pr.Branch)
		if err != nil {
			color.Red("error removing branch from stack: %s", err)
		}
		if err == nil {
			if err = git.DeleteBranchForce(pr.Branch); err != nil {
				color.Red("error deleting branch: %s", err)
			}
		}
	}
	return engine.Save()
}
