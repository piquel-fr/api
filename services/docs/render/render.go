package render

import (
	"bytes"
	"io"
	"slices"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/utils"
)

type RenderConfig struct {
	Instance   *models.DocsInstance
	PathPrefix string
}

func (r *Renderer) parseMarkdown(md []byte) ast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	return p.Parse(md)
}

func (r *Renderer) renderHTML(doc ast.Node) []byte {
	options := html.RendererOptions{
		Flags:          html.CommonFlags,
		RenderNodeHook: r.renderHook,
	}
	renderer := html.NewRenderer(options)

	return markdown.Render(doc, renderer)
}

func (r *Renderer) renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch node := node.(type) {
	case *ast.CodeBlock:
		r.renderCodeBlock(w, node)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func (r *Renderer) fixupAST(doc ast.Node, config *RenderConfig) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.Link:
			r.fixupLink(node, entering, config)
		}

		return ast.GoToNext
	})
	return doc
}

func (r *Renderer) fixupLink(link *ast.Link, entering bool, config *RenderConfig) {
	if !entering {
		return
	}

	if bytes.HasPrefix(link.Destination, []byte("http")) {
		link.AdditionalAttributes = append(link.AdditionalAttributes, "target=\"_blank\"")
	} else {
		link.Destination = slices.Concat([]byte(config.PathPrefix), utils.FormatLocalPath(link.Destination))
	}
}

var highlightStyle = styles.Get("tokyonight")

func (r *Renderer) renderCodeBlock(w io.Writer, codeBlock *ast.CodeBlock) error {
	lang := string(codeBlock.Info)
	source := string(codeBlock.Literal)
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	iterator, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}

	return r.htmlFormatter.Format(w, highlightStyle, iterator)
}
