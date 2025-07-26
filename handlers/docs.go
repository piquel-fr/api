package handlers

import (
	"net/http"
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/docs"
)

func HandleDocs(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path

	if strings.Trim(page, "/") == "" {
		// TODO: get user configurated root
		page = "index"
	}

	html, err := docs.GetDocumentaionPage(page)
	if err != nil {
		errors.HandleError(w, r, err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}
