package engine

import (
	"os"
)

type RepoState = string

const (
	RepoStateConflict RepoState = "conflict"
	RepoStateClean    RepoState = "clean"
	RepoStateUnknown  RepoState = "unknown"

	fileName = "status"
)

var fullFilePath = enginePath + fileName

func repoStateFromString(state string) RepoState {
	switch state {
	case "conflict":
		return RepoStateConflict
	case "clean":
		return RepoStateClean
	default:
		return RepoStateUnknown

	}
}

func GetState() (RepoState, error) {
	state, err := os.ReadFile(fullFilePath)
	if err != nil {
		return RepoStateUnknown, err
	}
	return repoStateFromString(string(state)), nil
}

func SetState(state RepoState) error {
	return os.WriteFile(fullFilePath, []byte(state), 0777)
}
