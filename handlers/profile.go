package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleProfileQuery(w http.ResponseWriter, r *http.Request) {
	// Get username from query params. Should look likes "GET api.piquel.fr/profile?[username]
	username := r.URL.Query().Get("profile")
	if username == "" {
		var err error
		username, err = auth.GetUsername(r)
		if username == "" || err != nil {
			http.Error(w, "You are not logged in", http.StatusUnauthorized)
			return
		}
	}
	writeProfile(w, r, username)
}

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]
	writeProfile(w, r, username)
}

func writeProfile(w http.ResponseWriter, r *http.Request, username string) {
	log.Printf("Fetching profile for %s!", username)
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
