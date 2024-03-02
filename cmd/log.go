package cmd

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/renderdag"
	"github.com/spf13/cobra"
)

type GHResult struct {
	PRs []gh.PRState
	Err error
}

var logCmd = &cobra.Command{
	Use:     "log",
	Short:   "Log stack",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}

		result := make(chan error)
		go func(trunk *engine.Stack) {
			result <- engine.SyncGithubWithLocal(trunk)
		}(trunk)

		renderdag.RenderDag(trunk)
		if err = <-result; err != nil {
			return err
		}

		return engine.Save()
	},
}

var logShortCmd = &cobra.Command{
	Use:     "log-short",
	Short:   "Log stack",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}

		renderdag.RenderDagShort(trunk)
		return engine.Save()
	},
}

func syncResults(trunk *engine.Stack, prResults *GHResult) error {
	return updateStack(trunk, prResults)
}

func updateStack(stack *engine.Stack, prResults *GHResult) error {
	if !stack.IsTrunk {
		pr := findPRStateForStack(prResults.PRs, stack)
		if pr == nil {
			stack.PRStatus = engine.PRStatusNone
		} else {
			stack.PRStatus = engine.PRStatusType(pr.State)
			stack.PRNumber = fmt.Sprint(pr.Number)
			stack.PRName = pr.Title
			stack.PRLink = pr.Link
		}
	}

	for _, child := range stack.Children {
		childStack, err := engine.GetStackByID(child)
		if err != nil {
			return err
		}
		if err = updateStack(childStack, prResults); err != nil {
			return err
		}
	}
	return nil
}

func findPRStateForStack(prs []gh.PRState, stack *engine.Stack) *gh.PRState {
	for _, pr := range prs {
		if pr.Branch == stack.Name {
			return &pr
		}
	}
	return nil
}
