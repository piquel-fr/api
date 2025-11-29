package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

func (h *Handler) CreateEmailHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("PUT /", h.handleAddAccount)
	handler.HandleFunc("GET /{email}", h.handleAccountInfo)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)

	// OPTIONS handlers
	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET", "PUT"))
	handler.Handle("OPTIONS /{email}", middleware.CreateOptionsHandler("GET", "DELETE"))

	return handler
}

func (h *Handler) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	requester, err := h.AuthService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	var user repository.User
	if username := r.URL.Query().Get("user"); username != "" {
		user, err = h.AuthService.GetUserFromUsername(r.Context(), username)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
	} else {
		user = *requester
	}

	if err := h.AuthService.Authorize(&auth.Request{
		User:      requester,
		Ressource: &user,
		Context:   r.Context(),
		Actions:   []string{"list_emails"},
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("count") {
		count, err := h.EmailService.CountAccounts(r.Context(), user.ID)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(int(count))))
		return
	}

	accounts, err := h.EmailService.ListAccounts(r.Context(), user.ID)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	for _, account := range accounts {
		account.Username = ""
		account.Password = ""
	}

	data, err := json.Marshal(accounts)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) handleAddAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.AuthService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "please submit your creation request with the required json payload", http.StatusBadRequest)
		return
	}

	params := repository.AddEmailAccountParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params.OwnerId = user.ID
	if _, err = database.Queries.AddEmailAccount(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *Handler) handleAccountInfo(w http.ResponseWriter, r *http.Request)   {}
func (h *Handler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {}
