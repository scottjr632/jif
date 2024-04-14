package cmd

import (
	"github.com/fatih/color"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/git"
	"github.com/spf13/cobra"
)

var prevCmd = &cobra.Command{
	Use:     "prev",
	Aliases: []string{"previous"},
	Short:   "Checkout the previous branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}
		parent, err := engine.GetStackByID(stack.Parent)
		if err != nil {
			return err
		}
		_, err = git.CheckoutBranch(parent.Name)
		return err
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Checkout the next branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.GetStackForCurrentBranch()
		if err != nil {
			return err
		}
		if len(stack.Children) == 0 {
			color.Yellow("No descendants found")
			return nil
		}

		if len(stack.Children) == 1 {
			child, err := engine.GetStackByID(stack.Children[0])
			if err != nil {
				return err
			}
			_, err = git.CheckoutBranch(child.Name)
			return err
		}

		childNames := make([]string, len(stack.Children))
		for i, childID := range stack.Children {
			child, err := engine.GetStackByID(childID)
			if err != nil {
				return err
			}
			childNames[i] = child.Name
		}

		prompt := promptui.Select{
			Label: "Select next branch",
			Items: childNames,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}

		_, err = git.CheckoutBranch(result)
		return err
	},
}

var checkoutCmd = &cobra.Command{
	Use:     "checkout",
	Aliases: []string{"co"},
	Short:   "Checkout a branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		branchesWithNames, err := engine.GetAllBranchesWithNames()
		if err != nil {
			return err
		}

		items := make([]string, len(branchesWithNames))
		for i, branch := range branchesWithNames {
			if branch.PRname != "" {
				items[i] = branch.PRname
			} else {
				if branch.CommitMessage == "" {
					items[i] = branch.Name
				} else {
					items[i] = branch.CommitMessage
				}
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
			return err
		}

		branchNameToCheckout := result
		for _, branch := range branchesWithNames {
			if branch.PRname == result || branch.CommitMessage == result {
				branchNameToCheckout = branch.Name
				break
			}
		}

		_, err = git.CheckoutBranch(branchNameToCheckout)
		return err
	},
}

func init() {
	RootCmd.AddCommand(checkoutCmd)
}
