package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorNotAuthenticated    = NewError("user is not authenticated", http.StatusUnauthorized)
	ErrorForbidden           = NewError("you are not allowed to access this ressource", http.StatusForbidden)
	ErrorNotFound            = NewError("Not Found", http.StatusNotFound)
	ErrorInternalServerError = NewError("Internal Server Error", http.StatusInternalServerError)
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

func getError(err error) *Error {
	if err == nil {
		panic("nil error being handled")
	}

	switch err := err.(type) {
	case *Error:
		return err
	case *json.SyntaxError:
		return NewError("syntax error in json payload", http.StatusBadRequest)
	case *json.UnmarshalTypeError:
		return NewError("type error in json payload", http.StatusBadRequest)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrorNotFound
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		return ErrorNotAuthenticated
	}

	panic(err)
}

func HandleHumaError(api huma.API, ctx huma.Context, inErr error) {
	err := getError(inErr)
	huma.WriteErr(api, ctx, err.status, err.Error())
}

func HandleError(w http.ResponseWriter, r *http.Request, inErr error) {
	err := getError(inErr)
	http.Error(w, err.Error(), err.status)
}
