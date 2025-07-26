package docs

import (
	"github.com/piquel-fr/api/services/docs/render"
	gh "github.com/piquel-fr/api/services/github"
)

type UserDocsConfig struct{}

func InitDocumentation() error {
	return render.InitRenderer()
}

func GetDocumentaionPage(route string) ([]byte, error) {
	file, err := gh.GetRepositoryFile("piquel-fr", "docs-test", "main", route+".md")
	if err != nil {
		return nil, err
	}

	html, err := render.RenderPage(file, &render.RenderConfig{})
	if err != nil {
		return nil, err
	}

	return html, err
}
