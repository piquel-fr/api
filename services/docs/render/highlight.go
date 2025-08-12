package render

import (
	"fmt"
	"io"
	"net/http"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown/ast"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
)

func renderCodeBlock(w io.Writer, codeBlock *ast.CodeBlock, entering bool, config *models.DocsInstance) error {
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

	style, err := getHighlightStyle(config)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, style, iterator)
}

func getHighlightStyle(config *models.DocsInstance) (*chroma.Style, error) {
	styleName := config.HighlightStyle
	if styleName == "" {
		styleName = "tokyonight"
	}

	highlightStyle := styles.Get(styleName)
	if highlightStyle == nil {
		return nil, errors.NewError(fmt.Sprintf("Couldn't find the style %s", styleName), http.StatusBadRequest)
	}

	return highlightStyle, nil
}
