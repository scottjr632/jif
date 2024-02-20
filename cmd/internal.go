package cmd

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var internalCmd = &cobra.Command{
	Use:   "internal-only",
	Short: "This command is only for internal use",
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, err := engine.Load()
		if err != nil {
			return err
		}
		fmt.Println(stack)
		return nil
	},
}
