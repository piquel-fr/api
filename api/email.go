package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/piquel-fr/api/database"
	"github.com/piquel-fr/api/database/repository"
	"github.com/piquel-fr/api/services/auth"
	"github.com/piquel-fr/api/services/email"
	"github.com/piquel-fr/api/utils/errors"
	"github.com/piquel-fr/api/utils/middleware"
)

type EmailHandler struct {
	authService  auth.AuthService
	emailService email.EmailService
}

func CreateEmailHandler(authService auth.AuthService, emailService email.EmailService) *EmailHandler {
	return &EmailHandler{authService, emailService}
}

func (h *EmailHandler) getName() string { return "email" }
func (h *EmailHandler) getSpec() Spec   { return nil }

func (h *EmailHandler) createHttpHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("PUT /", h.handleAddAccount)
	handler.HandleFunc("GET /{email}", h.handleAccountInfo)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)

	// sharing
	handler.HandleFunc("PUT /{email}/share", h.handleShareAccount)
	handler.HandleFunc("DELETE /{email}/share", h.handleRemoveAccountShare)

	// OPTIONS handlers
	handler.Handle("OPTIONS /", middleware.CreateOptionsHandler("GET", "PUT"))
	handler.Handle("OPTIONS /{email}", middleware.CreateOptionsHandler("GET", "DELETE"))
	handler.Handle("OPTIONS /{email}/share", middleware.CreateOptionsHandler("PUT", "DELETE"))

	return handler
}

func (h *EmailHandler) handleListAccounts(w http.ResponseWriter, r *http.Request) {
	requester, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	var user *repository.User
	if username := r.URL.Query().Get("user"); username != "" {
		user, err = h.authService.GetUserFromUsername(r.Context(), username)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}
	} else {
		user = requester
	}

	if err := h.authService.Authorize(&auth.Request{
		User:      requester,
		Ressource: user,
		Context:   r.Context(),
		Actions:   []string{auth.ActionListEmailAccounts},
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if r.URL.Query().Has("count") {
		count, err := h.emailService.CountAccounts(r.Context(), user.ID)
		if err != nil {
			errors.HandleError(w, r, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(strconv.Itoa(int(count))))
		return
	}

	accounts, err := h.emailService.ListAccounts(r.Context(), user.ID)
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

func (h *EmailHandler) handleAddAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetUserFromRequest(r)
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

func (h *EmailHandler) handleAccountInfo(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	accountInfo, err := h.emailService.GetAccountInfo(r.Context(), &account)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&auth.Request{
		User:      user,
		Ressource: &accountInfo,
		Actions:   []string{auth.ActionView},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	accountInfo.Username = ""
	accountInfo.Password = ""

	data, err := json.Marshal(accountInfo)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h *EmailHandler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&auth.Request{
		User:      user,
		Ressource: &account,
		Actions:   []string{auth.ActionDelete},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.emailService.RemoveAccount(r.Context(), account.ID); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleShareAccount(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&auth.Request{
		User:      user,
		Ressource: &account,
		Actions:   []string{auth.ActionShare},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	sharingUser, err := h.authService.GetUserFromUsername(r.Context(), r.URL.Query().Get("user"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.AddShareParams{
		UserId:     sharingUser.ID,
		Account:    account.ID,
		Permission: "",
	}

	if err := h.emailService.AddShare(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}

func (h *EmailHandler) handleRemoveAccountShare(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetUserFromRequest(r)
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	account, err := h.emailService.GetAccountByEmail(r.Context(), r.PathValue("email"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	if err := h.authService.Authorize(&auth.Request{
		User:      user,
		Ressource: &account,
		Actions:   []string{auth.ActionShare},
		Context:   r.Context(),
	}); err != nil {
		errors.HandleError(w, r, err)
		return
	}

	sharingUser, err := h.authService.GetUserFromUsername(r.Context(), r.URL.Query().Get("user"))
	if err != nil {
		errors.HandleError(w, r, err)
		return
	}

	params := repository.RemoveShareParams{
		UserId:  sharingUser.ID,
		Account: account.ID,
	}

	if err := h.emailService.RemoveShare(r.Context(), params); err != nil {
		errors.HandleError(w, r, err)
		return
	}
}
