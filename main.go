package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

	// clone or open
	clonePath := "/tmp/repo"
	repo, err := getRepository(setting, clonePath)
	if err != nil {
		log.Fatal(err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	// branch
	now := time.Now()
	branchName := strings.Replace(setting.title, " ", "_", -1) + "-" + now.Format("20060102")
	refName, err := createBranch(repo, branchName)
	if err != nil {
		log.Fatal(err)
	}
	worktree.Checkout(&git.CheckoutOptions{Branch: refName})

	// command
	execLog, err := execCommand(setting, clonePath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(execLog)

	status, err := worktree.Status()
	if err != nil {
		log.Fatal(err)
	}

	// not modified
	if status.IsClean() {
		log.Println("not modified")
		os.Exit(0)
	}

	// commit
	dateString := now.Format("2006-01-02 15:04:05")
	commitMessage := setting.title + " on " + dateString
	worktree.AddGlob(".")
	author := &object.Signature{
		Name:  setting.authorName,
		Email: setting.authorEmail,
		When:  now,
	}
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: author,
	})
	if err != nil {
		log.Fatal(err)
	}

	// push
	err = repo.Push(&git.PushOptions{Progress: os.Stdout})
	if err != nil {
		log.Fatal(err)
	}

	// Pull Request
	prURL, err := createPullRequest(setting, commitMessage, branchName, execLog)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(prURL)
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

func execCommand(setting *setting, clonePath string) (string, error) {
	err := os.Chdir(clonePath)
	if err != nil {
		return "", err
	}
	cmdStr := setting.command
	cmd := exec.Command("sh", "-c", cmdStr)
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return string(buffer.Bytes()), nil
}

func createBranch(repo *git.Repository, branchName string) (plumbing.ReferenceName, error) {
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

func createPullRequest(setting *setting, prTitle, branchName, execLog string) (string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: setting.accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	s := strings.Split(setting.repository, "/")
	owner, repo := s[0], s[1]
	base := "master"
	body := "Executed `" + setting.command + "`" + `

<details><summary>Output</summary>
` + execLog + `
</details>
`

	pr := &github.NewPullRequest{
		Title: &branchName,
		Head:  &branchName,
		Base:  &base,
		Body:  &body,
	}

	created, _, err := client.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		return "", err
	}
	return created.GetURL(), nil
}
