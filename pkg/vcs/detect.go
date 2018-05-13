package vcs

func DetectVcs(set VcsSet, vcs string) (Vcs, error) {
	switch vcs {

	case "GITHUB":
		return set.Git(), nil

	case "BITBUCKET-SERVER":
		return set.Git(), nil

	default:
		return nil, UnsupportedVcsType
	}
}
