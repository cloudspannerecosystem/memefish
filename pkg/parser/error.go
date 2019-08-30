package parser

import (
	"fmt"
)

type Error struct {
	Message string
	Loc     *Location
}

func (e *Error) String() string {
	return e.Error()
}

func (e *Error) Error() string {
	return fmt.Sprintf("parse error: %s: %s", e.Loc, e.Message)
}
