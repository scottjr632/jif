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

func CreatePR(base, branch string, args ...string) error {
	cmds := []string{"pr", "create", "--base", base, "--head", branch}
	cmds = append(cmds, args...)
	return RunGHCmd(cmds...)
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "merged", "--json", "number,headRefName,baseRefName,title,state,url", "--author", "@me")
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "closed", "--json", "number,headRefName,baseRefName,title,state,url", "--author", "@me")
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "open", "--json", "number,headRefName,baseRefName,title,state,url", "--author", "@me")
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
	out, err := cli.ExecuteCmd("gh", "pr", "list", "--state", "draft", "--json", "number,headRefName,baseRefName,title,state,url", "--author", "@me")
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

type PRBody struct {
	Body string `json:"body"`
}

func GetBodyForPR(prNumber string) (string, error) {
	out, err := cli.ExecuteCmd("gh", "pr", "view", prNumber, "--json", "body")
	if err != nil {
		return "", err
	}
	var prBody PRBody
	err = json.Unmarshal([]byte(out), &prBody)
	if err != nil {
		return "", err
	}
	return prBody.Body, nil
}

func UpdateBodyForPR(prNumber string, body string) error {
	_, err := cli.ExecuteCmd("gh", "pr", "edit", prNumber, "--body", body)
	return err
}
