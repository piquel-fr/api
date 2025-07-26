package render

import (
	"bytes"
	"io"
	"slices"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/piquel-fr/api/types"
	"github.com/piquel-fr/api/utils"
)

const tailwindBase = `
    h1 { font-size: 2em; }
    h2 { font-size: 1.5em; }
    h3 { font-size: 1.17em; }
    h4 { font-size: 1em; }
    h5 { font-size: 0.83em; }
    h6 { font-size: 0.67em; }
`

func parseMarkdown(md []byte, config *types.UserDocsConfig) ast.Node {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	return p.Parse(md)
}

func renderHTML(doc ast.Node, config *types.UserDocsConfig) []byte {
	htmlFlags := html.CommonFlags

	if config.FullPage {
		htmlFlags = htmlFlags | html.CompletePage
	}

	options := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: renderHook(config),
	}
	renderer := html.NewRenderer(options)

	return markdown.Render(doc, renderer)
}

func renderHook(config *types.UserDocsConfig) html.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		switch node := node.(type) {
		case *ast.CodeBlock:
			renderCodeBlock(w, node, entering, config)
			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	}
}

func addStyles(html []byte, config *types.UserDocsConfig) []byte {
	var styles []byte

	if config.UseTailwind {
		styles = append(styles, []byte(tailwindBase)...)
	}

	if styles == nil {
		return html
	}

	styles = slices.Concat([]byte("<style>\n"), styles, []byte("</style>\n"))
	return slices.Concat(html, styles)
}

func fixupAST(doc ast.Node, config *types.UserDocsConfig) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch node := node.(type) {
		case *ast.Link:
			fixupLink(node, entering, config)
		}

		return ast.GoToNext
	})
	return doc
}

func fixupLink(link *ast.Link, entering bool, config *types.UserDocsConfig) {
	if !entering {
		return
	}

	if bytes.HasPrefix(link.Destination, []byte("http")) {
		link.AdditionalAttributes = append(link.AdditionalAttributes, "target=\"_blank\"")
	} else {
		link.Destination = slices.Concat([]byte(config.Root), utils.FormatLocalPath(link.Destination, ".md"))
	}
}
