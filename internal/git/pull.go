package git

import "github.com/scottjr632/sequoia/internal/cli"

func FetchAndPullTrunk(trunkName string) error {
	// Fetch
	if err := cli.ExecuteCommandInTerminal("git", "fetch", "origin", trunkName); err != nil {
		return err
	}
	// Pull
	if err := cli.ExecuteCommandInTerminal("git", "pull", "origin", trunkName); err != nil {
		return err
	}
	return nil
}
