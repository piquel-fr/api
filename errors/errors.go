package errors

import (
	"encoding/json"
	"errors"
	"net/http"

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

	e, ok := err.(*Error)
	if ok {
		http.Error(w, e.Error(), e.status)
		return
	}

	if errors.Is(err, pgx.ErrNoRows) {
		http.NotFound(w, r)
		return
	}

	if errors.Is(err, &json.SyntaxError{}) {
		http.Error(w, "syntax error in json payload", http.StatusBadRequest)
		return
	}

	if errors.Is(err, &json.UnmarshalTypeError{}) {
		http.Error(w, "type error in json payload", http.StatusBadRequest)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	panic(err)
}
