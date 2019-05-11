package errors

import (
	"net/http"
)

var ErrBadRequest = &BadRequestError{
	"The request was invalid",
}

type BadRequestError struct {
	string
}

func NewBadRequestError(s string) *BadRequestError {
	return &BadRequestError{s}
}

func (e *BadRequestError) Error() string {
	return e.string
}

func (e *BadRequestError) Code() int {
	return http.StatusBadRequest
}
