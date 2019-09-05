package analyzer

import (
	"fmt"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Error struct {
	Message  string
	Position *parser.Position
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Error() string {
	// TODO: improve error message when valid e.Position.End is available.
	return fmt.Sprintf("syntax error: %s: %s", e.Position, e.Message)
}
