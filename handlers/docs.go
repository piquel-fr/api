package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/docs/render"
	"github.com/piquel-fr/api/utils"
	"github.com/piquel-fr/api/utils/errors"
	gh "github.com/piquel-fr/api/utils/github"
	"github.com/piquel-fr/api/utils/middleware"
)

func (h *Handler) CreateDocsHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", h.handleListDocs)
	handler.HandleFunc("POST /", h.handleNewDocs)
	handler.HandleFunc("GET /{documentation}", h.handleGetDocs)
	handler.HandleFunc("PUT /{documentation}", h.handleUpdateDocs)
	handler.HandleFunc("DELETE /{documentation}", h.handleDeleteDocs)
	handler.HandleFunc("GET /{documentation}/page/", h.handleGetDocsPage)

	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET", "POST"))
	handler.Handle("OPTIONS /{documentation}", middleware.CreateOptionsHandler("GET", "PUT", "DELETE"))
	handler.Handle("OPTIONS /{documentation}/page/", middleware.CreateOptionsHandler("GET"))

	return handler
}

func (h *Handler) handleListDocs(w http.ResponseWriter, r *http.Request) {
	requester, err := h.AuthService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	var instances []repository.DocsInstance

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

	if username := r.URL.Query().Get("user"); username != "" {
		user, err := h.AuthService.GetUserFromUsername(r.Context(), username)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		params := repository.ListUserDocsInstancesParams{
			OwnerId: user.ID,
			Limit:   int32(limit),
			Offset:  int32(offset),
		}

		instances, err = database.Queries.ListUserDocsInstances(r.Context(), params)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
	} else if r.URL.Query().Has("own") {
		if r.URL.Query().Has("count") {
			count, err := database.Queries.CountUserDocsInstances(r.Context(), requester.ID)
			if err != nil {
				errors.HandleError(w, r, err)
				return
			}

			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(strconv.Itoa(int(count))))
			return
		}

		params := repository.ListUserDocsInstancesParams{
			OwnerId: requester.ID,
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
		return
	} else {
		params := repository.ListDocsInstancesParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		}

		instances, err = database.Queries.ListDocsInstances(r.Context(), params)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
	}

	var returnedInstances []repository.DocsInstance
	for _, instance := range instances {
		if instance.Public {
			returnedInstances = append(returnedInstances, instance)
			continue
		}

		docsInstance := repository.DocsInstance(instance)
		authRequest := &auth.Request{
			User:      requester,
			Ressource: &docsInstance,
			Actions:   []string{"view"},
			Context:   r.Context(),
		}

		if err = h.AuthService.Authorize(authRequest); err != nil {
			continue
		}

		returnedInstances = append(returnedInstances, instance)
	}

	if r.URL.Query().Has("count") {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(len(returnedInstances))))
		return
	}

	data, err := json.Marshal(returnedInstances)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) handleNewDocs(w http.ResponseWriter, r *http.Request) {
	user, err := h.AuthService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	authRequest := &auth.Request{
		User:      user,
		Ressource: &repository.DocsInstance{},
		Actions:   []string{"create"},
		Context:   r.Context(),
	}

	if err := h.AuthService.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.AddDocsInstanceParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
	params.Root = utils.FormatLocalPathString(params.Root)

	if err := h.validateDocsInstance(params.Name, params.RepoOwner, params.RepoName, params.RepoRef, params.Root); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.OwnerId = user.ID
	if _, err = database.Queries.AddDocsInstance(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *Handler) handleGetDocs(w http.ResponseWriter, r *http.Request) {
	docsName := r.PathValue("documentation")
	config, err := database.Queries.GetDocsInstanceByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	docsConfig := repository.DocsInstance(config)

	user, err := h.AuthService.GetUserFromRequest(r)
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

	if err = h.AuthService.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(docsConfig)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) handleUpdateDocs(w http.ResponseWriter, r *http.Request) {
	docsName := r.PathValue("documentation")
	config, err := database.Queries.GetDocsInstanceByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	docsConfig := repository.DocsInstance(config)

	user, err := h.AuthService.GetUserFromRequest(r)
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

	if err = h.AuthService.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.UpdateDocsInstanceParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.Root = utils.FormatLocalPathString(params.Root)
	params.Name = strings.ToLower(params.Name)

	if err = h.validateDocsInstance(params.Name, params.RepoOwner, params.RepoName, params.RepoRef, params.Root); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.ID = docsConfig.ID
	err = database.Queries.UpdateDocsInstance(r.Context(), params)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) validateDocsInstance(name, owner, repo, ref, root string) error {
	// root cannot start with .
	if strings.HasPrefix(strings.Trim(root, "/"), ".") {
		return errors.NewError("root cannot start with \".\"", http.StatusBadRequest)
	}

	// repository must exist
	if !gh.RepositoryExists(owner, repo) {
		return errors.NewError(fmt.Sprintf("repository \"%s/%s\" does not exist", owner, repo), http.StatusBadRequest)
	}

	// root must exist
	if _, err := gh.GetRepositoryFile(owner, repo, ref, root); err != nil {
		return errors.NewError(fmt.Sprintf("file %s does not exist in %s/%s:%s", root, owner, repo, ref), http.StatusBadRequest)
	}

	// name cant have special characters
	if !utils.HasOnlyLettersAndNumbers(name) {
		return errors.NewError(fmt.Sprintf("name \"%s\" should only contain numbers or letter", name), http.StatusBadRequest)
	}
	return nil
}

func (h *Handler) handleDeleteDocs(w http.ResponseWriter, r *http.Request) {
	docsName := r.PathValue("documentation")
	config, err := database.Queries.GetDocsInstanceByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	docsConfig := repository.DocsInstance(config)

	user, err := h.AuthService.GetUserFromRequest(r)
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

	if err = h.AuthService.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	err = database.Queries.RemoveDocsInstance(r.Context(), docsConfig.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *Handler) handleGetDocsPage(w http.ResponseWriter, r *http.Request) {
	docsName := r.PathValue("documentation")
	page := r.URL.Path
	page = strings.Replace(page, docsName, "", 1)
	page = strings.Replace(page, "page", "", 1)
	page = utils.FormatLocalPathString(page)

	config, err := database.Queries.GetDocsInstanceByName(r.Context(), docsName)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}
	docsConfig := repository.DocsInstance(config)

	if !docsConfig.Public {
		user, err := h.AuthService.GetUserFromRequest(r)
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

		if err = h.AuthService.Authorize(authRequest); err != nil {
			errors.HandleError(w, r, err)
			return
		}
	}

	pathPrefix := r.URL.Query().Get("pathPrefix")
	utils.FormatLocalPathString(pathPrefix)

	renderConfig := render.RenderConfig{
		Instance:   &docsConfig,
		PathPrefix: pathPrefix,
	}

	html, err := h.DocsService.GetDocsInstancePage(page, &renderConfig)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(html)
}
