package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/types"
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
	user, err := database.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == pgx.ErrNoRows {
            // Properly redirect to cookied URL
			http.Redirect(w, r, "/", http.StatusNotFound)
			return
		}
	}

	profile := &types.UserProfile{User: user}

	group, err := database.Queries.GetGroupInfo(r.Context(), user.Group)
	if err != nil {
		panic(err)
	}

	profile.UserColor = group.Color
	profile.UserGroup = group.Displayname.String

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(profile)
}
