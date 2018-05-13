package vcs

import (
	"testing"
	"github.com/vastness-io/parser/pkg/vcs/git"
	"reflect"
)

func TestDetectVcs(t *testing.T) {
	tests := []struct{
		vcsSet VcsSet
		vcsType string
		expected Vcs
		err error
	} {
		{
			vcsSet:NewVcsSet(&git.InMemoryGit{}),
			vcsType: "GITHUB",
			expected: &git.InMemoryGit{},
			err: nil,
		},
		{
			vcsSet:NewVcsSet(&git.InMemoryGit{}),
			vcsType: "BITBUCKET-SERVER",
			expected: &git.InMemoryGit{},
			err: nil,
		},
		{
			vcsSet:NewVcsSet(&git.InMemoryGit{}),
			vcsType: "svn",
			expected: nil,
			err: UnsupportedVcsType,
		},
	}

	for _, test := range tests {

		v, err := DetectVcs(test.vcsSet, test.vcsType)

		if err != test.err {
			t.Fatalf("should equal")
		}

		if !reflect.DeepEqual(test.expected, v) {
			t.Fatalf("expected %v, got %v", test.expected, v)
		}
	}
}
