package handlers

import (
	"encoding/json"
	"net/http"

	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/middleware"
)

func CreateProfileHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("GET /", handleGetProfileQuery)
	handler.HandleFunc("GET /{user}", handleGetProfile)
	handler.HandleFunc("PUT /{user}", handleUpdateProfile)

	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET"))
	handler.Handle("OPTIONS /{user}", middleware.CreateOptionsHandler("GET", "PUT"))

	return handler
}

func handleGetProfile(w http.ResponseWriter, r *http.Request) {
	writeProfile(w, r, r.PathValue("user"))
}

func handleGetProfileQuery(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		id, err := auth.GetUserId(r)
		if err != nil {
			http.Error(w, "Please login or specify a username", http.StatusUnauthorized)
			return
		}
		profile, err := auth.GetProfileFromUserId(id)
		if err != nil {
			http.Error(w, "Please login or specify a username", http.StatusUnauthorized)
			return
		}
		username = profile.Username
	}

	writeProfile(w, r, username)
}

func writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	profile, err := auth.GetProfileFromUsername(username)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("user")

	profile, err := auth.GetProfileFromUsername(username)
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

	if err := auth.Authorize(request); err != nil {
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
