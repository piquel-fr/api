package handlers

import (
	"net/http"
)

func (h *Handler) CreateEmailHandler() http.Handler {
	handler := http.NewServeMux()

	// accounts
	handler.HandleFunc("GET /", h.handleListAccounts)
	handler.HandleFunc("PUT /", h.handleAddAccount)
	handler.HandleFunc("DELETE /{email}", h.handleRemoveAccount)
	handler.HandleFunc("GET /{email}", h.handleAccountInfo)

	// OPTIONS handlers

	return handler
}

func (h *Handler) handleListAccounts(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) handleAddAccount(w http.ResponseWriter, r *http.Request)    {}
func (h *Handler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleAccountInfo(w http.ResponseWriter, r *http.Request)   {}
