package handlers

import (
	"encoding/json"
	"net/http"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleBaseProfile(w http.ResponseWriter, r *http.Request) {
	// Get username from query params. Should look likes "GET api.piquel.fr/profile?profile=[username]
	username := r.URL.Query().Get("profile")
	if username == "" {
		var err error
		username, err = auth.GetUsername(r)
		if username == "" || err != nil {
			http.Error(w, "Please login or specify a username", http.StatusUnauthorized)
			return
		}
	}
	handleProfile(w, r, username)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]
	handleProfile(w, r, username)
}

func handleProfile(w http.ResponseWriter, r *http.Request, username string) {
	profile, err := users.GetProfile(username)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Properly redirect to cookied URL
			http.Redirect(w, r, "/", http.StatusNotFound)
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

func writeProfile(w http.ResponseWriter, r *http.Request, profile *types.UserProfile) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func updateProfile(w http.ResponseWriter, r *http.Request, profile *types.UserProfile) {
	w.WriteHeader(http.StatusOK)

	params := repository.UpdateUserParams{}
	if r.Header.Get("Content-Type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(params)
		if err != nil {
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

	w.WriteHeader(http.StatusOK)
}
