package handlers

import (
	"net/http"
	"strings"

	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/docs"
)

func HandleDocs(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path

	if strings.Trim(page, "/") == "" {
		// TODO: get user configurated root
		page = "index.md"
	}

	config := &models.Documentation{
		HighlightStyle: "tokyonight",
		FullPage:       false,
		UseTailwind:    true,
		Root:           "docs",
		RepoOwner:      "piquel-fr",
		RepoName:       "docs-test",
		RepoRef:        "main",
	}

	html, err := docs.GetDocumentaionPage(page, config)
	if err != nil {
		errors.HandleError(w, r, err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func HandleNewDocs(w http.ResponseWriter, r *http.Request)      {}
func HandleUpdateDocs(w http.ResponseWriter, r *http.Request)   {}
func HandleTransferDocs(w http.ResponseWriter, r *http.Request) {}
func HandleDeleteDocs(w http.ResponseWriter, r *http.Request)   {}
