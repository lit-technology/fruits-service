package errors

type HttpError interface {
	Code() int
	Error() string
}
