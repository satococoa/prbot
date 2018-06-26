package prbot

import (
	"log"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

// Execute will execute main routine
func Execute(setting *Setting) error {
	clonePath := "/tmp/repo"
	repo, err := FetchRepository(setting, clonePath)
	if err != nil {
		return err
	}
	log.Printf("cloned into %s", clonePath)

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	// branch
	now := time.Now()
	branchName := strings.Replace(setting.title, " ", "_", -1) + "-" + now.Format("20060102")
	refName, err := CreateBranch(repo, branchName)
	if err != nil {
		return err
	}
	worktree.Checkout(&git.CheckoutOptions{Branch: refName})

	// command
	execLog, err := ExecCommand(setting, clonePath)
	if err != nil {
		return err
	}

	status, err := worktree.Status()
	if err != nil {
		return err
	}

	// not modified
	if status.IsClean() {
		log.Printf("not modified")
		return nil
	}

	// commit
	dateString := now.Format("2006-01-02 15:04:05")
	commitMessage := setting.title + " on " + dateString
	err = Commit(setting, repo, commitMessage)
	if err != nil {
		return err
	}

	// push
	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	// Pull Request
	prURL, err := CreatePullRequest(setting, commitMessage, branchName, execLog)
	if err != nil {
		return err
	}

	log.Printf("Pull Request: %s", prURL)
	return nil
}
