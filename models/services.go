package models

import "net/http"

type Service interface {
	GetHandler() http.Handler
}
