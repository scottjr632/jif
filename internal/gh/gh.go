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
	Number int    `json:"number"`
	Link   string `json:"url"`
}

type PRResult struct {
	PRs []PRState
	Err error
}

func GetMergedPRs() ([]PRState, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "merged", "--json", "number,headRefName,baseRefName,title,state,url")
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "open", "--json", "number,headRefName,baseRefName,title,state,url")
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "draft", "--json", "number,headRefName,baseRefName,title,state,url")
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

func GetAllPRStats() ([]PRState, error) {
	mergedChan := make(chan PRResult)
	openChan := make(chan PRResult)
	draftChan := make(chan PRResult)

	go func() {
		prs, err := GetMergedPRs()
		mergedChan <- PRResult{PRs: prs, Err: err}
	}()

	go func() {
		prs, err := GetOpenPRs()
		openChan <- PRResult{PRs: prs, Err: err}
	}()

	go func() {
		prs, err := GetDraftPRs()
		draftChan <- PRResult{PRs: prs, Err: err}
	}()

	mergedResult := <-mergedChan
	// if mergedResult.Err != nil {
	// 	return nil, mergedResult.Err
	// }

	openResult := <-openChan
	// if openResult.Err != nil {
	// 	return nil, openResult.Err
	// }

	draftResult := <-draftChan
	// if draftResult.Err != nil {
	// 	return nil, draftResult.Err
	// }

	allPRs := append(mergedResult.PRs, openResult.PRs...)
	allPRs = append(allPRs, draftResult.PRs...)
	return allPRs, nil
}
