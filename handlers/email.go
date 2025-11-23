package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/utils/errors"
)

func (h *Handler) CreateEmailHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("POST /", h.handleAddAccount)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)

	// emails
	handler.HandleFunc("GET /{email}/get", h.handleListEmails)

	/**
	 * // accounts
	 * [x] GET / - list accounts
	 * [x] POST / - add account
	 * [x] DELETE {email} - delete account
	 * [ ] GET /{email} - get account info
	 *
	 * // emails
	 * [x] GET /{email}/get - list emails
	 * [ ] POST /{email} - send an email
	 *
	 * // OPTIONS handlers
	 */

	return handler
}

func (h *Handler) handleListEmails(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	requester, err := h.AuthService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid number %s specified for limit", limitStr), http.StatusBadRequest)
		return
	}

	if limit > 200 {
		limit = 200
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid number %s specified for offset", limitStr), http.StatusBadRequest)
		return
	}

	account, err := database.Queries.GetMailAccountByEmail(r.Context(), email)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	authRequest := &auth.Request{
		User:      requester,
		Ressource: &account,
		Actions:   []string{"fetch"},
		Context:   r.Context(),
	}

	if err = h.AuthService.Authorize(authRequest); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("count") {
		count, err := h.EmailService.CountEmailsForAccount(&account)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(int(count))))
		return
	}

	emails, err := h.EmailService.GetEmailsForAccount(&account, offset, limit)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	data, err := json.Marshal(emails)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *Handler) handleListAccounts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) handleAddAccount(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {}
