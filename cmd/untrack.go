package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var untrackCmd = &cobra.Command{
	Use:   "untrack",
	Short: "Untrack a branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName, err := git.GetCurrentBranchName()
		if err != nil {
			return err
		}
		prompt := promptui.Select{
			Label: fmt.Sprintf("Are you sure you want to untrack this branch? %s", branchName),
			Items: []string{"Yes", "No"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		if result == "No" {
			color.Green("Untrack aborted")
			return nil
		}

		if result == "Yes" {
			color.Green("Untracking branch")
			trunkBranch, err := engine.ReadTrunkName()
			if err != nil {
				return err
			}
			if err = engine.RemoveBranchFromStack(branchName); err != nil {
				return err
			}
			git.CheckoutBranch(trunkBranch)
			return engine.Save()
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(untrackCmd)
}
