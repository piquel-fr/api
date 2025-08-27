package render

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/chroma/v2/formatters/html"
	gh "github.com/piquel-fr/api/utils/github"
)

type Renderer struct {
	singleline, multiline *regexp.Regexp
	htmlFormatter         *html.Formatter
	gh                    *gh.GhWrapper
}

func NewRenderer(gh *gh.GhWrapper) *Renderer {
	singleline, err := regexp.Compile(`(?m)^{ *([a-z]+)(?: *\"(.*)\")? */}$`)
	if err != nil {
		panic(err)
	}

	multiline, err := regexp.Compile(`(?m)^{ *([a-z]+) *}\n?((?:.|\n)*?)\n?{/}$`)
	if err != nil {
		panic(err)
	}

	htmlFormatter := html.New()
	if htmlFormatter == nil {
		panic(fmt.Errorf("Error creating html formatter"))
	}

	return &Renderer{singleline: singleline, multiline: multiline, htmlFormatter: htmlFormatter}
}

func (r *Renderer) RenderPage(md []byte, config *RenderConfig) ([]byte, error) {
	md, err := r.renderCustom(md, config)
	if err != nil {
		return nil, err
	}

	ast := r.parseMarkdown(md)
	ast = r.fixupAST(ast, config)
	return r.renderHTML(ast), nil
}

func (r *Renderer) loadInclude(path string, config *RenderConfig) ([]byte, error) {
	file, err := r.gh.GetRepositoryFile(
		config.Instance.RepoOwner, config.Instance.RepoName,
		config.Instance.RepoRef, ".common/includes/"+path)
	return file, err
}
