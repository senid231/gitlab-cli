package api

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

type GitInfo struct {
	CurrentBranch string
	LastCommit    string
}

func NewGitInfo(path *string) (*GitInfo, error) {
	info := &GitInfo{}
	repo, err := git.PlainOpen(*path)
	if err != nil {
		return info, err
	}
	ref, err := repo.Head()
	if err != nil {
		return info, err
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return info, err
	}
	info.LastCommit = commit.Message
	info.CurrentBranch = strings.TrimPrefix(ref.Name().String(), "refs/heads/")
	return info, nil
}
