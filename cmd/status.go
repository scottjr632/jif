package cmd

import (
	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Show the status",
	Aliases: []string{"st"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.ExecuteCommandInTerminal("git", "status")
	},
}
