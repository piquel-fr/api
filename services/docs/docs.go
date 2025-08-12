package docs

import (
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/docs/render"
	gh "github.com/piquel-fr/api/services/github"
)

func InitDocsInstance() error {
	return render.InitRenderer()
}

func GetDocsInstancePage(route string, config *models.DocsInstance) ([]byte, error) {
	if strings.HasPrefix(strings.Trim(route, "/"), ".") {
		return nil, errors.ErrorNotFound
	}

	if route == "/" {
		route = config.Root
	}

	file, err := gh.GetRepositoryFile(config.RepoOwner, config.RepoName, config.RepoRef, route)
	if err != nil {
		return nil, err
	}

	html, err := render.RenderPage(file, config)
	if err != nil {
		return nil, err
	}

	return html, err
}
