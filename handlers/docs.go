package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/models"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/docs"
	"github.com/piquel-fr/api/services/users"
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

	user, err := users.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	authRequest := &auth.Request{
		User:      user,
		Ressource: &docsConfig,
		Actions:   []string{"view"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	html, err := docs.GetDocumentaionPage(page, &docsConfig)
	if err != nil {
		errors.HandleError(w, r, err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}

func HandleNewDocs(w http.ResponseWriter, r *http.Request) {
	user, err := users.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	authRequest := &auth.Request{
		User:      user,
		Ressource: &models.Documentation{},
		Actions:   []string{"create"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.AddDocumentationParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.OwnerId = user.ID
	_, err = database.Queries.AddDocumentation(r.Context(), params)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func HandleUpdateDocs(w http.ResponseWriter, r *http.Request)   {}
func HandleTransferDocs(w http.ResponseWriter, r *http.Request) {}
func HandleDeleteDocs(w http.ResponseWriter, r *http.Request)   {}
