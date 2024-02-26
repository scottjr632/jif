package cmd

import (
	"log"

	"github.com/manifoldco/promptui"
	"github.com/scottjr632/sequoia/internal/engine"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:     "commit",
	Short:   "Commit the current stack",
	Aliases: []string{"c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			log.Println(err)
			return err
		}

		if msg == "" {
			prompt := promptui.Prompt{
				Label: "Commit message",
			}
			res, err := prompt.Run()
			if err != nil {
				return err
			}
			msg = res
		}

		err = engine.CommitWithNewBranch(msg, engine.CommitOptions{AutoStage: true})
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	},
}

func init() {
	commitCmd.Flags().StringP("message", "m", "", "The message to use for the commit")
}
