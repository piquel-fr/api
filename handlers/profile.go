package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PiquelChips/piquel.fr/errors"
	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/permissions"
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
		request := &permissions.Request{
			User:      profile.User,
			Ressource: profile,
			Actions:   []string{"update"},
		}

		if err := permissions.Authorize(request); err != nil {
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
}
