package errors

import (
	"net/http"
)

var ErrNotFound = &NotFoundError{
	"This object was not found",
}

type NotFoundError struct {
	string
}

func NewNotFoundError(s string) *NotFoundError {
	return &NotFoundError{s}
}

func (e *NotFoundError) Error() string {
	return e.string
}

func (e *NotFoundError) Code() int {
	return http.StatusMethodNotAllowed
}
