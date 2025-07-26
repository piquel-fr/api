package render

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/chroma/v2/formatters/html"
	gh "github.com/piquel-fr/api/services/github"
	"github.com/piquel-fr/api/types"
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

func RenderPage(data []byte, config *types.UserDocsConfig) ([]byte, error) {
	custom, err := renderCustom(data, config)
	if err != nil {
		return nil, err
	}

	doc := parseMarkdown(custom, config)
	doc = fixupAST(doc, config)
	html := renderHTML(doc, config)
	return addStyles(html, config), nil
}

func loadInclude(path string, config *types.UserDocsConfig) ([]byte, error) {
	file, err := gh.GetRepositoryFile(config.RepoOwner, config.RepoName, config.RepoRef, ".common/includes/"+path+".md")
	return file, err
}
