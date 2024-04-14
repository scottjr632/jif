package git

func StageAll() (string, error) {
	return runGit("add", "-A")
}
