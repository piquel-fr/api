package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorNotAuthenticated = NewError("user is not authenticated", http.StatusUnauthorized)
	ErrorForbidden        = NewError("you are not allowed to access this ressource", http.StatusForbidden)
	ErrorNotFound         = NewError("Not Found", http.StatusNotFound)
)

type Error struct {
	message string
	status  int
}

func NewError(message string, status int) *Error {
	return &Error{message: message, status: status}
}

func (e *Error) Error() string {
	return e.message
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		panic("nil error being handled")
	}

	switch err := err.(type) {
	case *Error:
		http.Error(w, err.Error(), err.status)
		return
	case *json.SyntaxError:
		http.Error(w, "syntax error in json payload", http.StatusBadRequest)
		return
	case *json.UnmarshalTypeError:
		http.Error(w, "type error in json payload", http.StatusBadRequest)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		HandleError(w, r, ErrorNotAuthenticated)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	panic(err)
}
