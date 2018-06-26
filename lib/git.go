package prbot

import (
	"log"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// FetchRepository clones and pull primary branch
func FetchRepository(setting *Setting, clonePath string) (*git.Repository, error) {
	cloneURL := "https://" + setting.accessToken + "@github.com/" + setting.repository
	_, err := os.Stat(clonePath)
	log.Printf("Clone: %s -> %s", cloneURL, clonePath)
	if err != nil {
		repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
			URL: cloneURL,
		})
		if err != nil {
			log.Printf(err.Error())
			return nil, err
		}
		return repo, nil
	}

	repo, err := git.PlainOpen(clonePath)
	if err != nil {
		return nil, err
	}
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	err = worktree.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err.Error() != "already up-to-date" {
		return nil, err
	}

	return repo, nil
}

// CreateBranch create new branch
func CreateBranch(repo *git.Repository, branchName string) (plumbing.ReferenceName, error) {
	refName := plumbing.ReferenceName("refs/heads/" + branchName)
	headRef, err := repo.Head()
	if err != nil {
		return "", err
	}
	ref := plumbing.NewHashReference(refName, headRef.Hash())
	err = repo.Storer.SetReference(ref)
	if err != nil {
		return "", err
	}
	return refName, nil
}

// Commit commit
func Commit(setting *Setting, repo *git.Repository, commitMessage string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	worktree.AddGlob(".")
	author := &object.Signature{
		Name:  setting.authorName,
		Email: setting.authorEmail,
		When:  time.Now(),
	}
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		return err
	}

	return nil
}
