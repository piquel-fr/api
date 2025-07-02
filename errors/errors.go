package errors

import (
	"encoding/json"
	"net/http"
)

var ErrorNotAuthenticated *Error = NewError("User is not authenticated!", http.StatusUnauthorized)
var ErrorForbidden *Error = NewError("You are not allowed to access this ressource!", http.StatusForbidden)

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
	switch err.(type) {
	case *Error:
		e := err.(*Error)
		http.Error(w, e.Error(), e.status)
	case *json.SyntaxError:
		http.Error(w, "syntax error in json payload", http.StatusBadRequest)
	case *json.UnmarshalTypeError:
		http.Error(w, "type error in json payload", http.StatusBadRequest)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		panic(err)
	}
}
