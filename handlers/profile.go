package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleProfileQuery(w http.ResponseWriter, r *http.Request) {
	// Get username from query params. Should look likes "GET api.piquel.fr/profile?[username]
	username := r.URL.Query().Get("profile")
	writeProfile(w, r, username)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]
	writeProfile(w, r, username)
}

func writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	profile, err := users.GetProfile(username)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Properly redirect to cookied URL
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
