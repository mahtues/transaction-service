package apperrors

import (
	"fmt"
)

var (
	UnknkownError = &AppError{message: "unknown error", cause: nil}
)

type AppError struct {
	message string
	cause   error
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Cause() error { return e.cause }

func Newf(format string, args ...interface{}) error {
	return &AppError{
		cause:   nil,
		message: fmt.Sprintf(format, args...),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &AppError{
		cause:   err,
		message: fmt.Sprintf(format, args...),
	}
}

func Cause(err error) *AppError {
	type causer interface {
		Cause() error
	}

	for err != nil {
		_, ok := err.(*AppError)
		if ok {
			break
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	apperr, ok := err.(*AppError)
	if !ok {
		return UnknkownError
	}

	return apperr
}
