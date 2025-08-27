package gh

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v74/github"
	"github.com/piquel-fr/api/config"
	"github.com/piquel-fr/api/utils/errors"
)

var Client *github.Client

func InitGithubClient() {
	Client = github.NewClient(nil).WithAuthToken(config.Envs.GithubApiToken)
}

func GetRepositoryFile(owner, repo, ref, route string) ([]byte, error) {
	file, _, res, err := Client.Repositories.GetContents(context.Background(), owner, repo, route, &github.RepositoryContentGetOptions{Ref: ref})
	if res.StatusCode == http.StatusNotFound {
		return nil, errors.NewError(fmt.Sprintf("path %s does not exist in %s/%s:%s", route, owner, repo, ref), http.StatusNotFound)
	}
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
