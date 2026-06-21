package httperr

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound       = &Error{Code: http.StatusNotFound, Message: "not found"}
	ErrConflict       = &Error{Code: http.StatusConflict, Message: "already exists"}
	ErrUnauthorized   = &Error{Code: http.StatusUnauthorized, Message: "unauthorized"}
	ErrForbidden      = &Error{Code: http.StatusForbidden, Message: "forbidden"}
	ErrValidation     = &Error{Code: http.StatusUnprocessableEntity, Message: "validation error"}
	ErrInternal       = &Error{Code: http.StatusInternalServerError, Message: "internal error"}
	ErrBadRequest     = &Error{Code: http.StatusBadRequest, Message: "bad request"}
	ErrDuplicateKey   = errors.New("duplicate key")
)

func New(code int, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func RespondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, err error) {
	var e *Error
	if errors.As(err, &e) {
		RespondJSON(w, e.Code, map[string]string{"error": e.Message})
		return
	}
	RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
}
