package docs

import (
	"github.com/piquel-fr/api/services/docs/render"
	gh "github.com/piquel-fr/api/services/github"
	"github.com/piquel-fr/api/types"
)

func InitDocumentation() error {
	return render.InitRenderer()
}

func GetDocumentaionPage(route string) ([]byte, error) {
	config := types.UserDocsConfig{
		HighlightStyleName: "tokyonight",
		FullPage:           false,
		UseTailwind:        true,
		Root:               "docs",
		RepoOwner:          "piquel-fr",
		RepoName:           "docs-test",
		RepoRef:            "main",
	}

	file, err := gh.GetRepositoryFile(config.RepoOwner, config.RepoName, config.RepoRef, route+".md")
	if err != nil {
		return nil, err
	}

	html, err := render.RenderPage(file, &config)
	if err != nil {
		return nil, err
	}

	return html, err
}
