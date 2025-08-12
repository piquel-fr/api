package render

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/piquel-fr/api/models"
	gh "github.com/piquel-fr/api/services/github"
)

var (
	singleline    *regexp.Regexp
	multiline     *regexp.Regexp
	htmlFormatter *html.Formatter
)

func InitRenderer() error {
	var err error
	singleline, err = regexp.Compile(`(?m)^{ *([a-z]+)(?: *\"(.*)\")? */}$`)
	if err != nil {
		return err
	}

	multiline, err = regexp.Compile(`(?m)^{ *([a-z]+) *}\n?((?:.|\n)*?)\n?{/}$`)
	if err != nil {
		return err
	}

	htmlFormatter = html.New()
	if htmlFormatter == nil {
		return fmt.Errorf("Error creating html formatter")
	}

	return nil
}

func RenderPage(md []byte, config *models.DocsInstance) ([]byte, error) {
	md, err := renderCustom(md, config)
	if err != nil {
		return nil, err
	}

	ast := parseMarkdown(md)
	ast = fixupAST(ast, config)
	html := renderHTML(ast, config)
	return addStyles(html, config), nil
}

func loadInclude(path string, config *models.DocsInstance) ([]byte, error) {
	file, err := gh.GetRepositoryFile(config.RepoOwner, config.RepoName, config.RepoRef, ".common/includes/"+path)
	return file, err
}
