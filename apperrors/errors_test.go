package apperrors

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestAppErrors(t *testing.T) {
	n := 5

	t.Run("no AppError (aka unknown error)", func(t *testing.T) {
		expected := UnknkownError

		actual := Cause(errors.New("not app error"))

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("no error chain", func(t *testing.T) {
		expected := Newf("app error")
		actual := Cause(expected)
		if actual != expected {
			t.Errorf("expected: %v (%p), actual: %v (%p)", expected, expected, actual, actual)
		}
	})

	t.Run("no error chain with stack", func(t *testing.T) {
		expected := Newf("app error")
		err := errors.WithStack(expected)
		actual := Cause(err)
		if actual != expected {
			t.Errorf("expected: %v (%p), actual: %v (%p)", expected, expected, actual, actual)
		}
	})

	t.Run("no error chain with stack", func(t *testing.T) {
		expected := Newf("app error")
		err := errors.WithStack(expected)
		actual := Cause(err)
		if actual != expected {
			t.Errorf("expected: %v (%p), actual: %v (%p)", expected, expected, actual, actual)
		}
	})

	t.Run("error chain with app error", func(t *testing.T) {
		expected := Wrapf(errors.New("root cause"), "app error")

		err := errors.WithMessage(expected, "some message")
		for i := 0; i < n; i++ {
			err = errors.Wrapf(err, "i == %d", i)
		}

		actual := Cause(err)

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})

	t.Run("error chain without app error", func(t *testing.T) {
		err := errors.New("root cause")
		for i := 0; i < n; i++ {
			err = errors.Wrapf(err, "i == %d", i)
		}

		actual := Cause(err)

		if actual != UnknkownError {
			t.Errorf("expected: %v, actual: %v", UnknkownError, actual)
		}
	})

	t.Run("error chain with many app error", func(t *testing.T) {
		err := Wrapf(errors.New("root cause"), "app error #1")
		for i := 0; i < n; i++ {
			err = errors.Wrapf(err, "i == %d", i)
		}

		err = Wrapf(err, "app error #2")
		expected := err

		actual := Cause(err)

		for i := 0; i < n; i++ {
			err = errors.Wrapf(err, "i == %d", i)
		}

		if actual != expected {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func printStackTrace(err error) string {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	lines := []string{}

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			lines = append(lines, fmt.Sprintf("%+s:%d", f, f))
		}
	}

	return strings.Join(lines, "\n")
}
