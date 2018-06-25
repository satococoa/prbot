package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type setting struct {
	repository  string
	accessToken string
	baseBranch  string
	command     string
	title       string
	authorName  string
	authorEmail string
}

func main() {
	setting, err := getSetting()
	if err != nil {
		log.Fatal(err)
	}

	clonePath := "/tmp/repo"
	repo, err := getRepository(setting, clonePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("repository cloned into " + clonePath)
	log.Printf("repo: %+v\n", repo)

	err = execCommand(setting, clonePath)
	if err != nil {
		log.Fatal(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	status, err := worktree.Status()
	if err != nil {
		log.Fatal(err)
	}
	if status.IsClean() {
		os.Exit(0)
	}

	// commit
	dateString := time.Now().Format("2006-01-02 15:04:05")
	commitMessage := setting.title + " on " + dateString
	worktree.AddGlob(".")
	author := &object.Signature{
		Name:  setting.authorName,
		Email: setting.authorEmail,
	}
	hash, err := worktree.Commit(commitMessage, &git.CommitOptions{Author: author})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(hash)
}

func getSetting() (*setting, error) {
	repository, exists := os.LookupEnv("GITHUB_REPOSITORY")
	if !exists {
		return nil, errors.New("Please set GITHUB_REPOSITORY")
	}

	accessToken, exists := os.LookupEnv("GITHUB_ACCESS_TOKEN")
	if !exists {
		return nil, errors.New("Please set GITHUB_ACCESS_TOKEN")
	}

	baseBranch, exists := os.LookupEnv("BASE_BRANCH")
	if !exists {
		return nil, errors.New("Please set BASE_BRANCH")
	}

	command, exists := os.LookupEnv("COMMAND")
	if !exists {
		return nil, errors.New("Please set COMMAND")
	}

	title, exists := os.LookupEnv("TITLE")
	if !exists {
		return nil, errors.New("Please set TITLE")
	}

	authorName, exists := os.LookupEnv("AUTHOR_NAME")
	if !exists {
		return nil, errors.New("Please set AUTHOR_NAME")
	}

	authorEmail, exists := os.LookupEnv("AUTHOR_EMAIL")
	if !exists {
		return nil, errors.New("Please set AUTHOR_EMAIL")
	}

	setting := &setting{
		repository:  repository,
		accessToken: accessToken,
		baseBranch:  baseBranch,
		command:     command,
		title:       title,
		authorName:  authorName,
		authorEmail: authorEmail,
	}

	return setting, nil
}

func getRepository(setting *setting, clonePath string) (*git.Repository, error) {
	cloneURL := "https://" + setting.accessToken + "@github.com/" + setting.repository
	_, err := os.Stat(clonePath)
	if err != nil {
		repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
			URL:      cloneURL,
			Progress: os.Stdout,
		})
		if err != nil {
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

func execCommand(setting *setting, clonePath string) error {
	err := os.Chdir(clonePath)
	if err != nil {
		return err
	}
	args := strings.Split(setting.command, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
