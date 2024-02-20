package cmd

import (
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/scottjr632/sequoia/internal/renderdag"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Log stack",
	RunE: func(cmd *cobra.Command, args []string) error {
		trunk, err := engine.GetTrunk()
		if err != nil {
			return err
		}
		renderdag.RenderDag(trunk)
		return nil
	},
}
