package handlers

import "net/http"

func (h *Handler) CreateEmailHandler() http.Handler {
	handler := http.NewServeMux()
	return handler
}
