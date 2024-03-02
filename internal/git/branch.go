package git

import (
	"strings"

	"github.com/scottjr632/sequoia/internal/cli"
)

func CreateBranch(branch string) error {
	return cli.ExecuteCommandInTerminal("git", "branch", branch)
}

func CreateAndCheckoutBranch(branch string) error {
	err := CreateBranch(branch)
	if err != nil {
		return err
	}
	return cli.ExecuteCommandInTerminal("git", "checkout", branch)
}

// CreateBranchWithCommit creates a new branch and commits a message to it
// the branch name is a function of the commit message
func CreateBranchWithCommitAndCheckout(message string) error {
	branchName := FormatBranchNameFromCommit(message)
	return CreateAndCheckoutBranch(branchName)
}

func FormatBranchNameFromCommit(message string) string {
	san := strings.ReplaceAll(message, " ", "-")
	san = strings.ReplaceAll(san, "[", "_")
	san = strings.ReplaceAll(san, "]", "_")
	san = strings.ReplaceAll(san, "{", "_")
	san = strings.ReplaceAll(san, "}", "_")
	return san
}

func CheckoutBranch(branch string) (string, error) {
	return runGit("checkout", branch)
}

func GetAllBranchNames() ([]string, error) {
	result := make([]string, 0)
	out, err := runGit("branch", "--format", "%(refname:short)")
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(out, "\n") {
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

func GetCurrentBranchName() (string, error) {
	out, err := runGit("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(out), "\n", ""), nil
}

func PushCurrentBranchToRemoteIfNotExists() error {
	branch, err := GetCurrentBranchName()
	if err != nil {
		return err
	}
	exists, err := GetDoesBranchExistInRemote(branch)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = runGit("push", "origin", branch)
	return err
}

func GetDoesBranchExistInRemote(branch string) (bool, error) {
	out, err := runGit("ls-remote", "--heads", "origin", branch)
	if err != nil {
		return false, err
	}
	return out != "", nil
}

func GitPushForce(branchName string) error {
	_, err := runGit("push", "--force", "origin", branchName)
	return err
}

func PromptToPatch() error {
	return cli.ExecuteCommandInTerminal("git", "add", "-p")
}

func PromptToAddAll() error {
	return cli.ExecuteCommandInTerminal("git", "add", "-A")
}

func DeleteBranch(branchName string) error {
	return cli.ExecuteCommandInTerminal("git", "branch", "-d", branchName)
}

func DeleteBranchForce(branchName string) error {
	return cli.ExecuteCommandInTerminal("git", "branch", "-D", branchName)
}
