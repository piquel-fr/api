package gh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v74/github"
	"github.com/piquel-fr/api/errors"
)

var Client *github.Client

func InitGithubWrapper() {
	Client = github.NewClient(nil)
}

func GetRepositoryFile(owner, repo, ref, route string) (string, error) {
	file, _, res, err := Client.Repositories.GetContents(context.Background(), "piquel-fr", "docs-test", route, &github.RepositoryContentGetOptions{Ref: "main"})
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("github call failed with %d", res.StatusCode)
	}

	if file == nil {
		return "", errors.NewError(fmt.Sprintf("%s is a directory", route), http.StatusNotFound)
	}

	data, err := file.GetContent()
	if err != nil {
		return "", err
	}

	return data, nil
}
