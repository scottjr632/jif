package gh

import (
	"log"
	"strings"

	"github.com/scottjr632/sequoia/internal/cli"
)

func RunGHCmd(args ...string) error {
	return cli.ExecuteCommandInTerminal("gh", args...)
}

func RunAuthLogin() error {
	return RunGHCmd("auth", "login")
}

func CreatePR(base, branch string) error {
	log.Println("Creating PR for", branch, "to", base)
	return RunGHCmd("pr", "create", "--base", base, "--head", branch)
}

func CommentPR(prNumber, comment string) error {
	return RunGHCmd("pr", "comment", prNumber, "--body", comment)
}

func DoesPRExist(base, branch string) (bool, error) {
	log.Println("Checking if PR exists for", branch, "to", base)
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--base", base, "--head", branch)
	if err != nil {
		return false, err
	}
	doesExist := strings.Index(out, "no pull requests match") == -1 && out != ""
	return doesExist, nil
}
