package handlers

import (
	"net/http"
	"strings"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/PiquelChips/piquel.fr/views/settings"
)

func HandleSettingsRedirect(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/settings/profile", http.StatusTemporaryRedirect)
}

func HandleProfileSettings(w http.ResponseWriter, r *http.Request) {
    settings.ProfileSettings(users.GetPageData(w, r), "", 200).Render(r.Context(), w)
}

func HandleProfileSettingsUpdate(w http.ResponseWriter, r *http.Request) {
    user_id := 20
    params := repository.UpdateUserParams{
        ID: int32(user_id),
        Username: strings.ReplaceAll(strings.ToLower(r.FormValue("username")), " ", ""),
        Name: r.FormValue("name"),
        Image: r.FormValue("image"),
    }
    database.Queries.UpdateUser(r.Context(), params)
    
    settings.ProfileSettings(users.GetPageData(w, r), "Successfully updated settings!", 200).Render(r.Context(), w)
}
