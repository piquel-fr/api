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

func GetRepositoryFile(owner, repo, ref, route string) ([]byte, error) {
	file, _, res, err := Client.Repositories.GetContents(context.Background(), owner, repo, route, &github.RepositoryContentGetOptions{Ref: "main"})
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("github call failed with %d", res.StatusCode)
	}

	if file == nil {
		return nil, errors.NewError(fmt.Sprintf("%s is a directory", route), http.StatusNotFound)
	}

	data, err := file.GetContent()
	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

func RepositoryExists(owner, name string) bool {
	_, res, _ := Client.Repositories.Get(context.Background(), owner, name)
	if res.StatusCode == 200 {
		return true
	}
	return false
}
