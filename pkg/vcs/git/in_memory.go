package git

import (
	"fmt"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type InMemoryGit struct {
	repo *git.Repository
}

func (g *InMemoryGit) Clone(remoteURL string) error {

	if g.repo != nil {
		return nil
	}

	var (
		fs     = memfs.New()
		storer = memory.NewStorage()
	)

	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:               remoteURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Tags:              git.AllTags,
	})

	if err != nil {
		return err
	}

	g.repo = repo

	return nil
}

func (g *InMemoryGit) Checkout(versionish string) error {

	wt, err := g.repo.Worktree()

	if err != nil {
		return err
	}

	var (
		hash   = plumbing.NewHash(versionish)
		branch = plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/%s", versionish))
		tag    = plumbing.ReferenceName(fmt.Sprintf("refs/tags/%s", versionish))
	)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: branch,
	})

	if err != nil {

		if err == plumbing.ErrReferenceNotFound {

			err = wt.Checkout(&git.CheckoutOptions{
				Branch: tag,
			})

			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					if !hash.IsZero() {
						return wt.Checkout(&git.CheckoutOptions{
							Hash: hash,
						})
					}
				}
				return err
			}

		}
		return err
	}

	return nil

}

func (g *InMemoryGit) Open(file string) (billy.File, error) {
	wt, err := g.repo.Worktree()

	if err != nil {
		return nil, err
	}

	return wt.Filesystem.Open(file)
}

func (g *InMemoryGit) Cleanup() (err error) {
	*g = InMemoryGit{}
	return nil
}
