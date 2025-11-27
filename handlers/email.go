package handlers

import (
	"net/http"

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

func (h *Handler) handleListAccounts(w http.ResponseWriter, r *http.Request)  {}
func (h *Handler) handleAddAccount(w http.ResponseWriter, r *http.Request)    {}
func (h *Handler) handleAccountInfo(w http.ResponseWriter, r *http.Request)   {}
func (h *Handler) handleRemoveAccount(w http.ResponseWriter, r *http.Request) {}
