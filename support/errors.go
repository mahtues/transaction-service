package support

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func PrintStackTrace(err error) string {
	if err == nil {
		return ""
	}

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
