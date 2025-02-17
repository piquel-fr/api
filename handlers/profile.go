package handlers

import (
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/types"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["profile"]

	user, err := database.Queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == pgx.ErrNoRows {
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

    // Return the profile data to the user
}
