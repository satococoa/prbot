package prbot

import (
	"errors"
	"os"
)

// Setting for execution
type Setting struct {
	repository  string
	accessToken string
	baseBranch  string
	command     string
	title       string
	authorName  string
	authorEmail string
}

// NewSetting returns *Setting
func NewSetting() (*Setting, error) {
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

	setting := &Setting{
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
