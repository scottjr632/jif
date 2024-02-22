package cmd

import (
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var trunkCmd = &cobra.Command{
	Use:   "trunk",
	Short: "Checkout trunk",
	RunE: func(cmd *cobra.Command, args []string) error {
		trunk, err := engine.ReadTrunkName()
		if err != nil {
			return err
		}
		return cli.ExecuteCmdToStdout("git", "checkout", trunk)
	},
}
