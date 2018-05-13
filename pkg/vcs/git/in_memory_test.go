package git

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInMemoryGit_Clone(t *testing.T) {
	tests := []struct {
		err error
	}{
		{
			err: nil,
		},
	}

	for _, test := range tests {
		func() {
			gitVcs := InMemoryGit{}
			defer gitVcs.Cleanup()

			dir, _, cleanup := createTestRepository(map[string][]byte{
				"file_1": []byte("contents"),
			})
			defer cleanup()

			err := gitVcs.Clone(dir)

			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

		}()
	}

}

func TestInMemoryGit_Checkout(t *testing.T) {

	tests := []struct {
		versions []string
	}{
		{
			versions: []string{
				"master",
			},
		},
	}

	for _, test := range tests {
		func() {
			gitVcs := InMemoryGit{}
			defer gitVcs.Cleanup()

			dir, h, cleanup := createTestRepository(map[string][]byte{
				"file_1": []byte("contents"),
			})
			defer cleanup()

			err := gitVcs.Clone(dir)

			if err != nil {
				panic(err)
			}

			// test if we can checkout created hash
			err = gitVcs.Checkout(h)
			if err != nil {
				t.Fatalf("expected hash %s", h)
			}

			for _, version := range test.versions {
				err := gitVcs.Checkout(version)

				if err != nil {
					t.Fatalf("expected branch %s", version)
				}
			}

		}()
	}

}

func TestInMemoryGit_Open(t *testing.T) {

	tests := []struct {
		err           error
		expectedFiles map[string][]byte
	}{
		{
			expectedFiles: map[string][]byte{
				"file_1": []byte("contents"),
				"file_2": []byte("some more\ncontent"),
			},
			err: nil,
		},
	}

	for _, test := range tests {
		func() {
			gitVcs := InMemoryGit{}
			defer gitVcs.Cleanup()

			dir, _, cleanup := createTestRepository(test.expectedFiles)
			defer cleanup()

			err := gitVcs.Clone(dir)

			if err != test.err {
				t.Fatalf("expected %v, got %v", test.err, err)
			}

			for file, contents := range test.expectedFiles {

				f, err := gitVcs.Open(file)

				if err != nil {

					t.Fatalf("%s should exist", file)
				}

				b, _ := ioutil.ReadAll(f)

				if string(b) != string(contents) {
					t.Fatalf("expected %s, got %s", string(contents), string(b))
				}

			}

		}()
	}

}

func createTestRepository(files map[string][]byte) (dest string, hash string, cleanup func()) {
	src, err := ioutil.TempDir("", "test")

	if err != nil {
		panic(err)
	}

	repo, err := git.PlainInit(src, false)

	if err != nil {
		panic(err)
	}
	_, err = repo.CreateRemote(&config.RemoteConfig{

		Name: "origin",
		URLs: []string{
			fmt.Sprintf("%s", src),
		},
	})

	if err != nil {
		panic(err)
	}

	wt, err := repo.Worktree()

	if err != nil {
		panic(err)
	}

	for name, contents := range files {

		err := ioutil.WriteFile(filepath.Join(src, name), contents, 0644)

		if err != nil {
			panic(err)
		}

		_, err = wt.Add(name)

		if err != nil {
			panic(err)
		}
	}

	// Commit the file
	h, err := wt.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	return src, h.String(), func() {
		os.RemoveAll(src)
	}
}
