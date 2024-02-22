package gh

import (
	"encoding/json"
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

type PRState struct {
	State  string `json:"state"`
	Branch string `json:"headRefName"`
	Base   string `json:"baseRefName"`
	Title  string `json:"title"`
}

func GetMergedPRs() ([]PRState, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "merged", "--json", "number,headRefName,baseRefName,title,state", "--author", "@me")
	if err != nil {
		return nil, err
	}
	prs := []PRState{}
	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func GetClosedPRs() ([]PRState, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "closed", "--json", "number,headRefName,baseRefName,title,state", "--author", "@me")
	if err != nil {
		return nil, err
	}
	prs := []PRState{}
	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func GetOpenPRs() ([]PRState, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "open", "--json", "number,headRefName,baseRefName,title,state", "--author", "@me")
	if err != nil {
		return nil, err
	}
	prs := []PRState{}
	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

func GetDraftPRs() ([]PRState, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "draft", "--json", "number,headRefName,baseRefName,title,state", "--author", "@me")
	if err != nil {
		return nil, err
	}
	prs := []PRState{}
	err = json.Unmarshal([]byte(out), &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}
