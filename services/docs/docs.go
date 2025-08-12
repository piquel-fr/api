package docs

import (
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/docs/render"
	gh "github.com/piquel-fr/api/services/github"
)

func InitDocsService() {
	render.InitRenderer()
}

func GetDocsInstancePage(route string, config *render.RenderConfig) ([]byte, error) {
	if strings.HasPrefix(strings.Trim(route, "/"), ".") {
		return nil, errors.ErrorNotFound
	}

	if route == "/" {
		route = config.Instance.Root
	}

	file, err := gh.GetRepositoryFile(
		config.Instance.RepoOwner, config.Instance.RepoName,
		config.Instance.RepoRef, route,
	)
	if err != nil {
		return nil, err
	}

	html, err := render.RenderPage(file, config)
	if err != nil {
		return nil, err
	}

	return html, err
}
