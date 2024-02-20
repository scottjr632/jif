package cmd

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/cli"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Interact with git",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Unknown command, passing through to git")
		return cli.ExecuteCommandInTerminal("git", args...)
	},
}
