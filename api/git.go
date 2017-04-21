// Package api provides access to gitlab HTTP API, YAML config and GIT repo info
package api

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

// GitInfo holds git repo options
type GitInfo struct {
	CurrentBranch string
	LastCommit    string
}

// NewGitInfo read git options for path anf returns GitInfo object
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
