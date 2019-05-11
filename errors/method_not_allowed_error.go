package errors

import (
	"net/http"
)

var ErrMethodNotAllowed = &MethodNotAllowedError{
	"This HTTP method is not allowed",
}

type MethodNotAllowedError struct {
	string
}

func NewMethodNotAllowedError(s string) *MethodNotAllowedError {
	return &MethodNotAllowedError{s}
}

func (e *MethodNotAllowedError) Error() string {
	return e.string
}

func (e *MethodNotAllowedError) Code() int {
	return http.StatusMethodNotAllowed
}
