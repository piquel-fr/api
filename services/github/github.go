package gh

import "github.com/google/go-github/v74/github"

var Client *github.Client

func InitGithubWrapper() {
	Client = github.NewClient(nil)
}
