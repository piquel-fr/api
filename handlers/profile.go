package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/users"
	"github.com/piquel-fr/api/models"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleBaseProfile(w http.ResponseWriter, r *http.Request) {
	// Get username from query params. Should look likes "GET api.piquel.fr/profile?username=[username]

	username := r.URL.Query().Get("username")
	if username == "" {
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
	handleProfile(w, r, username)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]
	handleProfile(w, r, username)
}

func handleProfile(w http.ResponseWriter, r *http.Request, username string) {
	profile, err := users.GetProfileFromUsername(username)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Properly redirect to cookied URL
			http.Error(w, fmt.Sprintf("user %s does not exist", username), http.StatusNotFound)
			return
		}
		panic(err)
	}

	switch r.Method {
	case http.MethodGet:
		writeProfile(w, r, profile)
	case http.MethodPut:
		request := &auth.Request{
			User:      profile.User,
			Ressource: profile,
			Actions:   []string{"update"},
		}

		if err := auth.Authorize(request); err != nil {
			errors.HandleError(w, r, err)
			return
		}

		updateProfile(w, r, profile)
	}
}

func writeProfile(w http.ResponseWriter, r *http.Request, profile *models.UserProfile) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func updateProfile(w http.ResponseWriter, r *http.Request, profile *models.UserProfile) {
	w.WriteHeader(http.StatusOK)

	params := repository.UpdateUserParams{}
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			errors.HandleError(w, r, err)
			return
		}
	} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		if err := r.ParseForm(); err != nil {
			errors.HandleError(w, r, err)
			return
		}

		params.Name = r.FormValue("name")
		params.Username = r.FormValue("username")
		params.Image = r.FormValue("image")
	}

	params.ID = profile.ID

	if err := database.Queries.UpdateUser(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}
