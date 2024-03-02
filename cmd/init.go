package cmd

import (
	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize jiffy for the current repo at .",
	RunE: func(cmd *cobra.Command, args []string) error {
		branchNames, err := git.GetAllBranchNames()
		if err != nil {
			return err
		}
		result, err := promptForTrunk(branchNames)
		if err != nil {
			return err
		}
		return engine.InitEngine(result)
	},
}

func promptForTrunk(branchNames []string) (string, error) {
	prompt := promptui.Select{
		Label: "Select the trunk branch",
		Items: branchNames,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
