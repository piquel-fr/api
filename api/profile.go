package api

import (
	"encoding/json"
	"net/http"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type ProfileHandler struct {
	authService auth.AuthService
}

func CreateProfileHandler(authService auth.AuthService) *ProfileHandler {
	return &ProfileHandler{authService}
}

func (h *ProfileHandler) getName() string { return "profile" }
func (h *ProfileHandler) getSpec() Spec   { return nil }

func (h *ProfileHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", h.handleGetProfileQuery)
	handler.HandleFunc("GET /{user}", h.handleGetProfile)
	handler.HandleFunc("PUT /{user}", h.handleUpdateProfile)

	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{user}", middleware.CreateOptionsHandler("GET", "PUT"))

	return handler
}

func (h *ProfileHandler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	h.writeProfile(w, r, r.PathValue("user"))
}

func (h *ProfileHandler) handleGetProfileQuery(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		id, err := h.authService.GetUserId(r)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
		user, err := h.authService.GetUserFromUserId(r.Context(), id)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
		username = user.Username
	}

	h.writeProfile(w, r, username)
}

func (h *ProfileHandler) writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	user, err := h.authService.GetUserFromUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *ProfileHandler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	user, err := h.authService.GetUserFromUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &auth.Request{
		User:      user,
		Ressource: user,
		Actions:   []string{auth.ActionUpdate},
		Context:   r.Context(),
	}

	if err := h.authService.Authorize(request); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.UpdateUserParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.ID = user.ID

	if err := database.Queries.UpdateUser(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}
