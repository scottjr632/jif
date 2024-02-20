package git

type RebaseOptions struct {
	GoBackToPreviousBranch bool
}

func NewRebaseOptions() RebaseOptions {
	return RebaseOptions{
		GoBackToPreviousBranch: false,
	}
}

func RebaseBranchOnto(branchName, onto string, options RebaseOptions) (string, error) {
	if options.GoBackToPreviousBranch {
		previousBranch, err := GetCurrentBranchName()
		if err != nil {
			return "", err
		}
		defer CheckoutBranch(previousBranch)
	}

	_, err := CheckoutBranch(branchName)
	if err != nil {
		return "", err
	}

	return runGit("rebase", onto)
}
