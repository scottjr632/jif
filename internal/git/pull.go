package git

import (
	"github.com/fatih/color"
	"github.com/scottjr632/sequoia/internal/cli"
)

func FetchAndPullTrunk(trunkName string) error {
	color.White("Fast forwarding trunk...")
	if err := cli.ExecuteCmdToStdout("git", "fetch", "origin", trunkName); err != nil {
		return err
	}
	if err := cli.ExecuteCmdToStdout("git", "pull", "origin", trunkName); err != nil {
		return err
	}
	color.Green("Trunk updated!")
	return nil
}
