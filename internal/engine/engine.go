package engine

import (
	"fmt"
	"os"

	"github.com/scottjr632/sequoia/internal/gh"
	"github.com/scottjr632/sequoia/internal/git"
)

const (
	enginePath = "./.jf/"
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

func WriteTrunkName(trunkName string) error {
	file, err := os.Create(enginePath + "trunk")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(trunkName)
	return err
}

func ReadTrunkName() (string, error) {
	trunkName, err := os.ReadFile(enginePath + "trunk")
	if err != nil {
		return "", err
	}
	return string(trunkName), err
}

func GetAllStackNames() ([]string, error) {
	trunk, err := GetTrunk()
	if err != nil {
		return nil, err
	}
	stackNames := make([]string, 0)
	exploreAndAppendStack(trunk, &stackNames)
	return stackNames, nil
}

func exploreAndAppendStack(stack *Stack, stackNames *[]string) {
	*stackNames = append(*stackNames, stack.Name)
	for _, child := range stack.Children {
		childStack, err := GetStackByID(child)
		if err != nil {
			return
		}
		exploreAndAppendStack(childStack, stackNames)
	}
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

	if err := WriteTrunkName(trunkName); err != nil {
		return err
	}

	_, err = initializeEngineFile(trunkName, sha)
	if err != nil {
		return err
	}
	return nil
}

func syncResults(trunk *Stack, prs []gh.PRState) error {
	return updateStack(trunk, prs)
}

func updateStack(stack *Stack, prs []gh.PRState) error {
	if !stack.IsTrunk {
		pr := findPRStateForStack(prs, stack)
		if pr == nil {
			stack.PRStatus = PRStatusNone
		} else {
			stack.PRStatus = PRStatusType(pr.State)
			stack.PRNumber = fmt.Sprint(pr.Number)
			stack.PRName = pr.Title
			stack.PRLink = pr.Link
		}
	}

	for _, child := range stack.Children {
		childStack, err := GetStackByID(child)
		if err != nil {
			return err
		}
		if err = updateStack(childStack, prs); err != nil {
			return err
		}
	}
	return nil
}

func findPRStateForStack(prs []gh.PRState, stack *Stack) *gh.PRState {
	for _, pr := range prs {
		if pr.Branch == stack.Name {
			return &pr
		}
	}
	return nil
}

func SyncGithubWithLocal(trunk *Stack) error {
	prs, err := gh.GetAllPRStats()
	if err != nil {
		return err
	}
	if err = syncResults(trunk, prs); err != nil {
		return err
	}
	return Save()
}

func (s *Stack) ensureUniqueChildren(with *Stack) error {
	for _, child := range s.Children {
		if child == with.ID {
			return fmt.Errorf("duplicate child: %s", s.Name)
		}
	}
	return nil
}

func (s *Stack) removeChildStack(with *Stack) {
	for i, child := range s.Children {
		if child == with.ID {
			s.Children = append(s.Children[:i], s.Children[i+1:]...)
			break
		}
	}
}

func (s *Stack) Rebase(onto *Stack) error {
	if err := s.ensureUniqueChildren(onto); err != nil {
		return err
	}

	parentStack, err := GetStackByID(s.Parent)
	if err != nil {
		return err
	}

	parentStack.removeChildStack(s)
	onto.Children = append(onto.Children, s.ID)
	s.Parent = onto.ID
	return Save()
}
