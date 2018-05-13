package git

import (
	"testing"
	"path/filepath"
	"os"
)

func TestInMemoryGit_Clone(t *testing.T) {

	var (
		testHelpersRepo = "../../../test-helpers"
	)

	tests := []struct {
		remoteURL string
		err error
		expectedFiles []string
		nonExistentFiles []string
	}{
		{
			remoteURL: getAbsolutePath(testHelpersRepo),
			err: nil,
			expectedFiles: []string {
				"README.md",
				"LICENSE",
			},
			nonExistentFiles: []string {
				"this-file-doesnt-exist",
			},
		},
	}

	for _, test := range tests {
		func () {
			git := InMemoryGit{}
			defer git.Cleanup()
			err := git.Clone(test.remoteURL)

			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

			for _, file := range test.expectedFiles {
				_, err := git.Open(file)

				if err != nil {
					t.Fatalf("%s should exist", file)
				}
			}

			for _, file := range test.nonExistentFiles {
				_, err := git.Open(file)

				if err != os.ErrNotExist {
					t.Fatalf("%s should not exist", file)
				}
			}
		}()
	}

}

func TestInMemoryGit_Checkout(t *testing.T) {

	var (
		testHelpersRepo = "../../../test-helpers"
	)
	
	tests := []struct {
		remoteURL string
		versions []string
	}{
		{
			remoteURL: getAbsolutePath(testHelpersRepo),
			versions: []string {
				"master",
				"with-pom",
			},
		},
	}

	for _, test := range tests {
		func () {
			git := InMemoryGit{}
			defer git.Cleanup()

			_ = git.Clone(test.remoteURL)

			for _, version := range test.versions {
				err := git.Checkout(version)

				if err != nil {
					t.Fatalf("expected branch %s", version)
				}
			}

		}()
	}

}

func TestInMemoryGit_Open(t *testing.T) {

	var (
		testHelpersRepo = "../../../test-helpers"
	)

	tests := []struct {
		remoteURL string
		err error
		expectedFiles []string
	}{
		{
			remoteURL: getAbsolutePath(testHelpersRepo),
			err: nil,
			expectedFiles: []string {
				"README.md",
				"LICENSE",
			},
		},
	}

	for _, test := range tests {
		func () {
			git := InMemoryGit{}
			defer git.Cleanup()
			err := git.Clone(test.remoteURL)

			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

			for _, file := range test.expectedFiles {
				_, err := git.Open(file)

				if err != nil {
					t.Fatalf("%s should exist", file)
				}
			}
		}()
	}

}

func getAbsolutePath(relative string) string {
	abs, err := filepath.Abs(relative)

	if err != nil {
		return ""
	}

	return abs
}
