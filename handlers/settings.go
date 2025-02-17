package handlers

import (
	"net/http"

	repository "github.com/PiquelChips/piquel.fr/database/generated"
	"github.com/PiquelChips/piquel.fr/services/database"
	"github.com/PiquelChips/piquel.fr/utils"
)

func HandleProfileSettingsUpdate(w http.ResponseWriter, r *http.Request) {
    user_id := 20
    params := repository.UpdateUserParams{
        ID: int32(user_id),
        Username: utils.FormatUsername(r.FormValue("username")),
        Name: r.FormValue("name"),
        Image: r.FormValue("image"),
    }
    database.Queries.UpdateUser(r.Context(), params)
    
    // Send udated profile back to user
}
