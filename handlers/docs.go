package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
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
		errors.HandleError(w, r, err)
		return
	}
	docsConfig := models.Documentation(config)

	if !docsConfig.Public {
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
	}

	if r.URL.Query().Has("config") {
		data, err := json.Marshal(docsConfig)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
		return
	}

	html, err := docs.GetDocumentaionPage(page, &docsConfig)
	if err != nil {
		errors.HandleError(w, r, err)
		return
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

func HandleUpdateDocs(w http.ResponseWriter, r *http.Request) {
	docsName := mux.Vars(r)["documentation"]
	config, err := database.Queries.GetDocumentationByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
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
		Actions:   []string{"update"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.UpdateDocumentationParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.ID = docsConfig.ID
	err = database.Queries.UpdateDocumentation(r.Context(), params)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func HandleTransferDocs(w http.ResponseWriter, r *http.Request) {
	docsName := mux.Vars(r)["documentation"]
	destination := mux.Vars(r)["username"]

	config, err := database.Queries.GetDocumentationByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
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
		Actions:   []string{"transfer"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	destinationUser, err := database.Queries.GetUserByUsername(r.Context(), destination)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.TransferDocumentationParams{
		ID:      docsConfig.ID,
		OwnerId: destinationUser.ID,
	}

	err = database.Queries.TransferDocumentation(r.Context(), params)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func HandleDeleteDocs(w http.ResponseWriter, r *http.Request) {
	docsName := mux.Vars(r)["documentation"]
	config, err := database.Queries.GetDocumentationByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
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
		Actions:   []string{"delete"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	err = database.Queries.RemoveDocumentation(r.Context(), docsConfig.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func HandleListDocs(w http.ResponseWriter, r *http.Request) {
	username, ok := mux.Vars(r)["username"]
	if !ok {
		id, err := auth.GetUserId(r)
		if err != nil {
			http.Error(w, "Please login or specify a username", http.StatusUnauthorized)
			return
		}
		profile, err := users.GetProfileFromUserId(id)
		if err != nil {
			http.Error(w, "Please login or specify a username", http.StatusUnauthorized)
			return
		}
		username = profile.Username
	}

	user, err := users.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	requestedUser, err := database.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	authRequest := &auth.Request{
		User:      user,
		Ressource: &models.Documentation{OwnerId: requestedUser.ID},
		Actions:   []string{"list"},
		Context:   r.Context(),
	}

	if err = auth.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("count") {
		count, err := database.Queries.CountUserDocsInstances(r.Context(), requestedUser.ID)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(int(count))))
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid number %s specified for limit", limitStr), http.StatusBadRequest)
		return
	}

	if limit > 200 {
		limit = 200
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid number %s specified for offset", limitStr), http.StatusBadRequest)
		return
	}

	params := repository.ListUserDocsInstancesParams{
		OwnerId: requestedUser.ID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	}

	configs, err := database.Queries.ListUserDocsInstances(r.Context(), params)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(configs)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
