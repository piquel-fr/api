package errors

import "errors"

var ErrorNotAuthenticated error = errors.New("User is not authenticated!")
