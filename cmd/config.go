package cmd

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Manage user configuration settings",
}

var configDisplayCmd = &cobra.Command{
	Use:     "display",
	Aliases: []string{"d"},
	Short:   "Display the current configuration settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.LoadPersistenConfig()
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", c)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration setting for the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := config.LoadPersistenConfig()
		if err != nil {
			return err
		}

		branchPrefix, err := cmd.Flags().GetString("branch-prefix")
		if err != nil {
			return err
		}
		if branchPrefix != "" {
			if err := c.SetBranchPrefix(branchPrefix); err != nil {
				return err
			}
		}
		return nil
	},
}

func initFlags() {
	configSetCmd.Flags().StringP("branch-prefix", "b", "", "Set the branch prefix for the workspace")
}

func init() {
	initFlags()
	configCmd.AddCommand(configDisplayCmd)
	configCmd.AddCommand(configSetCmd)
	RootCmd.AddCommand(configCmd)
}
