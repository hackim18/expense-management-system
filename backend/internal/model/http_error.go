package model

import (
	"go-expense-management-system/internal/messages"
	nethttp "net/http"
)

type HTTPError struct {
	Code int
	Msg  string
	Err  error
}

func (e *HTTPError) Error() string {
	return e.Msg
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Msg: message}
}

func (e *HTTPError) Status() int {
	return e.Code
}

func (e *HTTPError) Message() string {
	return e.Msg
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

var (
	ErrBadRequest          = NewHTTPError(nethttp.StatusBadRequest, messages.FailedInputFormat)
	ErrUnauthorized        = NewHTTPError(nethttp.StatusUnauthorized, messages.Unauthorized)
	ErrForbidden           = NewHTTPError(nethttp.StatusForbidden, messages.Forbidden)
	ErrNotFound            = NewHTTPError(nethttp.StatusNotFound, messages.StatusNotFound)
	ErrConflict            = NewHTTPError(nethttp.StatusConflict, messages.ConflictError)
	ErrInternalServerError = NewHTTPError(nethttp.StatusInternalServerError, messages.InternalServerError)
)
