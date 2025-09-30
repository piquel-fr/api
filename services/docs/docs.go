package docs

import (
	"strings"

	"github.com/piquel-fr/api/services/docs/render"
	"github.com/piquel-fr/api/utils/errors"
	gh "github.com/piquel-fr/api/utils/github"
)

type DocsService interface {
	GetDocsInstancePage(route string, config *render.RenderConfig) ([]byte, error)
}

type realDocsService struct {
	renderer *render.Renderer
}

func NewRealDocsService() *realDocsService {
	return &realDocsService{renderer: render.NewRenderer()}
}

func (r *realDocsService) GetDocsInstancePage(route string, config *render.RenderConfig) ([]byte, error) {
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

	html, err := r.renderer.RenderPage(file, config)
	if err != nil {
		return nil, err
	}

	return html, err
}
