package memefish

import (
	"bytes"
	"fmt"

	"github.com/cloudspannerecosystem/memefish/token"
)

type ErrorList []*Error

func (list ErrorList) String() string {
	return list.String()
}

func (list ErrorList) Error() string {
	var message bytes.Buffer
	for i, err := range list {
		if i > 0 {
			fmt.Fprintln(&message)
		}
		fmt.Fprintln(&message, err.Error())
	}
	return message.String()
}

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
