package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

func (h *Handler) CreateProfileHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", h.handleGetProfileQuery)
	handler.HandleFunc("GET /{user}", h.handleGetProfile)
	handler.HandleFunc("PUT /{user}", h.handleUpdateProfile)

	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{user}", middleware.CreateOptionsHandler("GET", "PUT"))

	return handler
}

func (h *Handler) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	h.writeProfile(w, r, r.PathValue("user"))
}

func (h *Handler) handleGetProfileQuery(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		id, err := h.AuthService.GetUserId(r)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
		profile, err := h.AuthService.GetProfileFromUserId(r.Context(), id)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
		username = profile.Username
	}

	h.writeProfile(w, r, username)
}

func (h *Handler) writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	profile, err := h.AuthService.GetProfileFromUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (h *Handler) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	profile, err := h.AuthService.GetProfileFromUsername(r.Context(), username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	request := &auth.Request{
		User:      profile.User,
		Ressource: profile,
		Actions:   []string{"update"},
		Context:   r.Context(),
	}

	if err := h.AuthService.Authorize(request); err != nil {
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

	params.ID = profile.ID

	if err := database.Queries.UpdateUser(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}
