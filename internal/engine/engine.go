package engine

import (
	"os"

	"github.com/scottjr632/sequoia/internal/git"
)

const (
	enginePath = "./.so/"
)

func DoesEngineExist() bool {
	_, err := os.Stat(enginePath)
	return !os.IsNotExist(err)
}

func writeGitIgnore() error {
	file, err := os.Create(enginePath + ".gitignore")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("/*")
	return err
}

func CreateEngine() error {
	if DoesEngineExist() {
		return nil
	}

	if err := os.Mkdir(enginePath, 0755); err != nil {
		return err
	}
	return writeGitIgnore()
}

func DeleteEngine() error {
	return os.RemoveAll(enginePath)
}

func InitEngine(trunkName string) error {
	if err := CreateEngine(); err != nil {
		return err
	}

	if _, err := git.CheckoutBranch(trunkName); err != nil {
		return err
	}

	sha, err := git.GetCurrentBranchCommitSha()
	if err != nil {
		return err
	}

	_, err = initializeEngineFile(trunkName, sha)
	if err != nil {
		return err
	}
	return nil
}
