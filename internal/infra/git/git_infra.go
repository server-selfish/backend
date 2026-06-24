package git_infra

import "github.com/go-git/go-git/v6"

type (
	GitInfra interface {
		PlainClone(path string, o *git.CloneOptions) (*git.Repository, error)
	}
	gitInfra struct{}
)

func NewGitInfra() GitInfra {
	return gitInfra{}
}

// PlainClone implements [GitInfra].
func (g gitInfra) PlainClone(path string, o *git.CloneOptions) (*git.Repository, error) {
	return git.PlainClone(path, o)
}
