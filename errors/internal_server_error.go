package errors

import (
	"net/http"
)

var ErrInternalServer = &InternalServerError{
	"There was an internal server error. Please try again",
}

type InternalServerError struct {
	string
}

func NewInternalServerError(s string) *InternalServerError {
	return &InternalServerError{s}
}

func (e *InternalServerError) Error() string {
	return e.string
}

func (e *InternalServerError) Code() int {
	return http.StatusInternalServerError
}
