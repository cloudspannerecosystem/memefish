package memefish

import (
	"bytes"
	"fmt"

	"github.com/cloudspannerecosystem/memefish/token"
)

type Error struct {
	Message  string
	Position *token.Position
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Error() string {
	var message bytes.Buffer
	fmt.Fprintf(&message, "syntax error: %s: %s\n", e.Position, e.Message)
	if e.Position.Source != "" {
		fmt.Fprintln(&message)
		fmt.Fprint(&message, e.Position.Source)
	}
	return message.String()
}
