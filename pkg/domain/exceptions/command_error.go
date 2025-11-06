package exceptions

import (
	"strings"
)

type CommandError struct {
	Messages []string
}

func (e CommandError) Error() string {
	return "validation errors: " + strings.Join(e.Messages, "; ")
}
