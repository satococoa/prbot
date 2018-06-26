package prbot

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// CreatePullRequest sends pull request to github
func CreatePullRequest(setting *Setting, prTitle, branchName, execLog string) (string, error) {
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
<pre>
` + execLog + `
</pre>
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
	return created.GetHTMLURL(), nil
}
