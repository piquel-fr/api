package handlers

import (
	"net/http"

	"github.com/PiquelChips/piquel.fr/services/users"
	"github.com/PiquelChips/piquel.fr/views"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	views.Home(users.GetPageData(w, r)).Render(r.Context(), w)
}
