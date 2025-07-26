package docs

import (
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/docs/render"
	gh "github.com/piquel-fr/api/services/github"
)

func InitDocumentation() error {
	return render.InitRenderer()
}

func GetDocumentaionPage(route string, config *models.Documentation) ([]byte, error) {
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
