package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

func New(code int, msg string) *AppError { return &AppError{Code: code, Message: msg} }

var (
	ErrUnauthorized   = New(http.StatusUnauthorized, "unauthorized")
	ErrForbidden      = New(http.StatusForbidden, "forbidden")
	ErrNotFound       = New(http.StatusNotFound, "not found")
	ErrBadRequest     = New(http.StatusBadRequest, "bad request")
	ErrConflict       = New(http.StatusConflict, "conflict")
	ErrInternal       = New(http.StatusInternalServerError, "internal server error")
	ErrInvalidCreds   = New(http.StatusUnauthorized, "invalid credentials")
	ErrTenantNotFound = New(http.StatusNotFound, "tenant not found")
	ErrUserNotFound   = New(http.StatusNotFound, "user not found")
	ErrItemNotFound   = New(http.StatusNotFound, "item not found")
	ErrOrderNotFound  = New(http.StatusNotFound, "order not found")
	ErrTableNotFound  = New(http.StatusNotFound, "table not found")
)
