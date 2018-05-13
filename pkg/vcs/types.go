package vcs

import (
	"errors"
	"gopkg.in/src-d/go-billy.v4"
)

var (
	UnsupportedVcsType = errors.New("unsupported vcs")
	FallbackError      = errors.New("FB")
)

type Vcs interface {
	Clone(remoteURL string) error
	Checkout(versionish string) error
	Open(string) (billy.File, error)
	Cleanup() error
}

type VcsSet interface {
	Git() Vcs
}

type vcsSet struct {
	git Vcs
}

func (vs *vcsSet) Git() Vcs {
	return vs.git
}

func NewVcsSet(git Vcs) VcsSet {
	return &vcsSet{
		git: git,
	}
}
