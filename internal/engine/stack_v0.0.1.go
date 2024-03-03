package engine

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/scottjr632/sequoia/internal/git"
)

const (
	version         = "0.0.1"
	stackBinaryName = "stack"
)

var __stacks Stacks
var stacksOnce sync.Once

type StackID = uint64
type VersionID = uint64

type PRStatusType = string

const (
	PRStatusNone   PRStatusType = "none"
	PRStatusOpen   PRStatusType = "open"
	PRStatusClosed PRStatusType = "closed"
	PRStatusMerged PRStatusType = "merged"
	PRStatusDraft  PRStatusType = "draft"

	PRStackCommentIdentifier = "[//]: # (BEGIN JIFFY FOOTER)"
)

var engineFullPath = enginePath + stackBinaryName + "_" + version

type Stack struct {
	ID        StackID
	Name      string
	IsDirty   bool
	IsTrunk   bool
	Sha       string
	Parent    StackID
	PRStatus  PRStatusType
	PRNumber  string
	PRLink    string
	PRName    string
	Children  []StackID
	Revisions []string
}

type Stacks = []*Stack

func getStacks() Stacks {
	stacksOnce.Do(func() {
		localStack, err := Load()
		if err != nil {
			panic(err)
		}
		sort.SliceStable(localStack, func(i, j int) bool {
			return localStack[i].ID < localStack[j].ID
		})
		__stacks = localStack
	})
	return __stacks
}

func getNextID() StackID {
	localStacks := getStacks()
	if len(localStacks) == 0 {
		return 1
	}
	return localStacks[len(localStacks)-1].ID + 1
}

func newStack(id StackID, name string, isDirty bool, isTrunk bool, sha string, parentID StackID) *Stack {
	newStack := &Stack{
		ID:        id,
		Name:      name,
		IsDirty:   isDirty,
		IsTrunk:   isTrunk,
		Sha:       sha,
		Parent:    parentID,
		PRStatus:  PRStatusNone,
		Children:  make([]StackID, 0),
		Revisions: make([]string, 0),
	}
	__stacks = append(__stacks, newStack)
	return newStack
}

func NewStack(name string, isDirty bool, isTrunk bool, sha string, parentID StackID) *Stack {
	nextID := getNextID()
	newStack := newStack(nextID, name, isDirty, isTrunk, sha, parentID)
	return newStack
}

func newTrunk(name string, sha string) *Stack {
	trunk := newStack(1, name, false, true, sha, 0)
	return trunk
}

func (s *Stack) AddChild(childStackID StackID) {
	s.Children = append(s.Children, childStackID)
}

func (s *Stack) AddRevision(sha string) {
	s.Revisions = append(s.Revisions, sha)
}

func (s *Stack) PopRevision() string {
	if len(s.Revisions) == 0 {
		return ""
	}
	sha := s.Revisions[len(s.Revisions)-1]
	s.Revisions = s.Revisions[:len(s.Revisions)-1]
	return sha
}

func (s *Stack) GetParent() *Stack {
	if s.Parent == 0 {
		return nil
	}
	for _, stack := range getStacks() {
		if stack.ID == s.Parent {
			return stack
		}
	}
	return nil
}

type TrunkNotFoundError struct{}

func (e *TrunkNotFoundError) Error() string {
	return "Trunk not found"
}

func IsTrunkNotFoundError(err error) bool {
	_, ok := err.(*TrunkNotFoundError)
	return ok
}

func Save() error {
	file, err := os.OpenFile(engineFullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer file.Close()
	enc := gob.NewEncoder(file)
	stacks := getStacks()
	return enc.Encode(stacks)
}

func Load() (Stacks, error) {
	file, err := os.OpenFile(engineFullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	dec := gob.NewDecoder(file)
	var stacks []*Stack
	if err := dec.Decode(&stacks); err != nil {
		return nil, err
	}
	return stacks, nil
}

func initializeEngineFile(trunkName string, trunkSha string) (Stacks, error) {
	trunk := newTrunk(trunkName, trunkSha)
	file, err := os.OpenFile(engineFullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	enc := gob.NewEncoder(file)
	stacks := make(Stacks, 0)
	stacks = append(stacks, trunk)
	if err := enc.Encode(&stacks); err != nil {
		return nil, err
	}
	return stacks, nil
}

func GetTrunk() (*Stack, error) {
	for _, stack := range getStacks() {
		if stack.IsTrunk {
			return stack, nil
		}
	}
	return nil, &TrunkNotFoundError{}
}

func getStackForBranch(branchName string) (*Stack, error) {
	for _, stack := range getStacks() {
		if stack.Name == branchName {
			return stack, nil
		}
	}
	return nil, fmt.Errorf("stack not found for branch %s", branchName)
}

func GetStackForBranch(branchName string) (*Stack, error) {
	// start at the trunk and work our way up
	return getStackForBranch(branchName)
}

func GetStackForCurrentBranch() (*Stack, error) {
	currentBranch, err := git.GetCurrentBranchName()
	if err != nil {
		return nil, err
	}

	return GetStackForBranch(currentBranch)
}

func GetStackByID(stackID StackID) (*Stack, error) {
	for _, stack := range getStacks() {
		if stack.ID == stackID {
			return stack, nil
		}
	}
	return nil, fmt.Errorf("stack not found for ID %d", stackID)
}

func RemoveBranchFromStack(branchName string) error {
	stack, err := GetStackForBranch(branchName)
	if err != nil {
		return err
	}

	if stack.IsTrunk {
		return fmt.Errorf("cannot remove trunk")
	}

	parent := stack.GetParent()
	if parent == nil {
		return fmt.Errorf("parent not found")
	}

	for i, childID := range parent.Children {
		if childID == stack.ID {
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			break
		}
	}

	for _, childID := range stack.Children {
		child, err := GetStackByID(childID)
		if err != nil {
			return err
		}

		if _, err = git.RebaseBranchOnto(child.Name, parent.Name, git.RebaseOptions{GoBackToPreviousBranch: false}); err != nil {
			return err
		}
		child.Parent = parent.ID
		parent.AddChild(child.ID)
	}

	removeStackID(stack.ID)
	return Save()
}

func RestackChildren(stack *Stack) error {
	for _, childID := range stack.Children {
		child, err := GetStackByID(childID)
		if err != nil {
			return err
		}
		if child == nil {
			continue
		}

		_, err = git.RebaseBranchOnto(child.Name, stack.Name, git.RebaseOptions{GoBackToPreviousBranch: true})
		if err != nil {
			return err
		}
		RestackChildren(child)
	}
	_, err := git.CheckoutBranch(stack.Name)
	return err
}

func removeStackID(stackID StackID) {
	stacks := getStacks()
	for i, stack := range stacks {
		if stack.ID == stackID {
			stacks = append(stacks[:i], stacks[i+1:]...)
			break
		}
	}
	__stacks = stacks
}

func GetStackForCommentByStackID(stackID StackID) (string, error) {
	stack, err := GetStackByID(stackID)
	if err != nil {
		return "", err
	}
	return GetStackForCommentByStack(stack), nil
}

func GetStackForCommentByStack(stack *Stack) string {
	// we want to get all the parents up to the trunk and all the children that have been submited to GH

	stackNames := make([]string, 0)
	appendStackNameForParent(stack.GetParent(), &stackNames)
	stackWithLink := fmt.Sprintf("#%s", stack.PRNumber)
	stackNames = append(stackNames, "👉  "+stackWithLink)
	appendStackNameForChildren(stack, &stackNames)

	builder := PRStackCommentIdentifier + "\r\n"
	builder += "\r\n---\r\n"
	builder += "\r\n" + "Stacked dependencies on/for current PR:\r\n"

	for _, name := range stackNames {
		builder += "* " + name + "\r\n"
	}
	return builder
}

func GetStringWithoutStackComment(text string) string {
	index := strings.Index(text, PRStackCommentIdentifier)
	if index == -1 {
		return text
	}
	return text[:index]

}

func appendStackNameForParent(stack *Stack, stackNames *[]string) {
	if stack == nil || stack.IsTrunk || stack.PRNumber == "" {
		return
	}
	appendStackNameForParent(stack.GetParent(), stackNames)
	nameWithLink := fmt.Sprintf("#%s", stack.PRNumber)
	*stackNames = append(*stackNames, nameWithLink)
}

func appendStackNameForChildren(stack *Stack, stackNames *[]string) {
	if stack == nil || stack.PRNumber == "" {
		return
	}
	for i, childID := range stack.Children {
		child, err := GetStackByID(childID)
		if err != nil {
			continue
		}
		namekWithLink := fmt.Sprintf("#%s", child.PRNumber)
		*stackNames = append(*stackNames, strings.Repeat("   ", i)+namekWithLink)
		appendStackNameForChildren(child, stackNames)
	}
}

func GetStacksInStack(stack *Stack) []*Stack {
	stacks := make([]*Stack, 0)
	appendStackForParent(stack.GetParent(), &stacks)
	stacks = append(stacks, stack)
	appendStackForChildren(stack, &stacks)
	return stacks
}

func appendStackForParent(stack *Stack, stacks *[]*Stack) {
	if stack == nil || stack.IsTrunk {
		return
	}
	appendStackForParent(stack.GetParent(), stacks)
	*stacks = append(*stacks, stack)
}

func appendStackForChildren(stack *Stack, stacks *[]*Stack) {
	if stack == nil {
		return
	}
	for _, childID := range stack.Children {
		child, err := GetStackByID(childID)
		if err != nil {
			continue
		}
		*stacks = append(*stacks, child)
		appendStackForChildren(child, stacks)
	}
}
