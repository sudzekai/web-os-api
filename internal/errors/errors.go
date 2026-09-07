package errors

import (
	"fmt"
	"strings"
)

func MkError(format string, args ...any) error {
	return fmt.Errorf(
		"%s.",
		strings.TrimRight(fmt.Sprintf(format, args...), "."),
	)
}
