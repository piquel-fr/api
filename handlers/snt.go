package handlers

import (
	"net/http"

	"github.com/PiquelChips/piquel.fr/views/snt"
)

func HandleLinus(w http.ResponseWriter, r *http.Request) {
    snt.Linus().Render(r.Context(), w)
}

func HandleLinux(w http.ResponseWriter, r *http.Request) {
    snt.Linux().Render(r.Context(), w)
}

func HandleMusset(w http.ResponseWriter, r *http.Request) {
    snt.Musset().Render(r.Context(), w)
}
