package git_infra

import (
	"github.com/go-git/go-billy/v6"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/storage"
)

type (
	GitInfra interface {
		PlainClone(path string, o *git.CloneOptions) (*git.Repository, error)
		Clone(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error)
	}
	gitInfra struct{}
)

func NewGitInfra() GitInfra {
	return gitInfra{}
}

// Clone implements [GitInfra].
func (g gitInfra) Clone(s storage.Storer, worktree billy.Filesystem, o *git.CloneOptions) (*git.Repository, error) {
	return git.Clone(s, worktree, o)
}

// PlainClone implements [GitInfra].
func (g gitInfra) PlainClone(path string, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, o)
}
