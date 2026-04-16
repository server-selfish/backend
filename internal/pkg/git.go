package pkg

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/storer"
)

func ExtractRepoMetaData(repo *git.Repository) (string, string, string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", "", "", err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", "", "", err
	}
	var commitId, commitMsg, version string
	commitId = commit.Hash.String()[:7]
	commitMsg = commit.Message

	var matchedTag string

	tags, err := repo.Tags()
	if err != nil {
		return "", "", "", err
	}
	_ = tags.ForEach(func(refTag *plumbing.Reference) error {
		tagHash := refTag.Hash()

		tagObj, err := repo.TagObject(tagHash)
		if err == nil {
			tagHash = tagObj.Target
		}

		if tagHash == commit.Hash {
			matchedTag = refTag.Name().Short()
			return storer.ErrStop
		}

		return nil
	})
	if matchedTag != "" {
		version = matchedTag
	} else {
		version = commit.Hash.String()[:7]
	}
	return commitId, commitMsg, version, nil
}
