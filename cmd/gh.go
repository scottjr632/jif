package cmd

import (
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/spf13/cobra"
)

var ghCmd = &cobra.Command{
	Use:   "gh",
	Short: "Interact with GitHub",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.ExecuteCommandInTerminal("gh", args...)
	},
}
