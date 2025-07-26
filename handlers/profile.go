package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	repository "github.com/piquel-fr/api/database/generated"
	"github.com/piquel-fr/api/errors"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/database"
	"github.com/piquel-fr/api/services/users"
)

func HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]
	writeProfile(w, username)
}

func HandleGetProfileQuery(w http.ResponseWriter, r *http.Request) {
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

	writeProfile(w, username)
}

func writeProfile(w http.ResponseWriter, username string) {
	profile, err := users.GetProfileFromUsername(username)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, fmt.Sprintf("user %s does not exist", username), http.StatusNotFound)
			return
		}
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]

	profile, err := users.GetProfileFromUsername(username)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, fmt.Sprintf("user %s does not exist", username), http.StatusNotFound)
			return
		}
		panic(err)
	}

	request := &auth.Request{
		User:      profile.User,
		Ressource: profile,
		Actions:   []string{"update"},
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
