package git

import (
	"strings"

	"github.com/fatih/color"
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
	san = strings.ReplaceAll(san, "(", "_")
	san = strings.ReplaceAll(san, ")", "_")
	san = strings.ReplaceAll(san, ":", "-")
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

func GetLatestLocalCommitSha(branch string) (string, error) {
	out, err := runGit("rev-parse", branch)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(strings.TrimSpace(out), "\n", ""), nil
}

func PushCurrentBranchToRemoteIfNotExistsAndNeedsUpdate() error {
	branch, err := GetCurrentBranchName()
	if err != nil {
		return err
	}
	exists, err := GetDoesBranchExistInRemoteAndShouldUpdate(branch)
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

func GetDoesBranchExistInRemoteAndShouldUpdate(branch string) (bool, error) {
	exists, err := GetDoesBranchExistInRemote(branch)
	if err != nil {
		return false, err
	}

	if !exists {
		color.White("Branch %s does not exist in remote", branch)
		return false, nil
	}

	out, err := runGit("ls-remote", "origin", branch)
	if err != nil {
		return false, err
	}

	parts := strings.Split(out, "\t")
	if len(parts) != 2 {
		color.Red("Error parsing ls-remote output: %s", out)
		return false, err
	}

	remoteCommit := strings.ReplaceAll(strings.TrimSpace(parts[0]), "\n", "") // parts[0]
	localCommit, err := GetLatestLocalCommitSha(branch)
	if err != nil {
		color.Red("Error getting latest local commit sha: %s", err)
		return false, err
	}

	shouldUpdate := remoteCommit != localCommit
	if shouldUpdate {
		color.Yellow("Branch %s exists in remote and needs update", branch)
	} else {
		color.Yellow("Branch %s exists in remote and does not need update", branch)
	}

	return shouldUpdate, nil
}

func GitPushForce(branchName string) error {
	_, err := runGit("push", "--force-with-lease", "origin", branchName)
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
