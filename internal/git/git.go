package git

import (
	"fmt"

	"github.com/scottjr632/sequoia/internal/cli"
)

func runGit(cmds ...string) (result string, err error) {
	out, err := cli.ExecuteCmd("git", cmds...)
	if err != nil {
		fmt.Println(out)
	}
	return out, err
}

func runGitAsync(cmds ...string) <-chan cli.CmdResult {
	return cli.ExecuteCmdAsync("git", cmds...)
}
