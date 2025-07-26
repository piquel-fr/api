package handlers

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/docs"
	"github.com/piquel-fr/api/utils"
)

func HandleDocs(w http.ResponseWriter, r *http.Request) {
	docsName := mux.Vars(r)["documentation"]
	page := r.URL.Path
	page = strings.Replace(page, "docs", "", 1)
	page = strings.Replace(page, docsName, "", 1)
	page = utils.FormatLocalPathString(page)

	config, err := database.Queries.GetDocumentationByName(r.Context(), docsName)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		panic(err)
	}

	docsConfig := models.Documentation(config)
	html, err := docs.GetDocumentaionPage(page, &docsConfig)
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
