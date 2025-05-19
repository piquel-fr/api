package errors

import (
	"fmt"
	"net/http"
)

var ErrorNotAuthenticated error = NewError("User is not authenticated!", http.StatusUnauthorized)

type Error struct {
	status int
	error
}

func NewError(message string, status int) Error {
	return Error{error: fmt.Errorf(message), status: status}
}

func (e *Error) Handle(w http.ResponseWriter) {
	http.Error(w, e.Error(), e.status)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if internalError, ok := err.(Error); ok {
		internalError.Handle(w)
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	panic(err)
}
