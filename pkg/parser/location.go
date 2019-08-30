package parser

import (
	"fmt"
)

type Location struct {
	FilePath             string
	Offset, Line, Column int
}

func (loc *Location) String() string {
	return fmt.Sprintf("%s:%d:%d:", loc.FilePath, loc.Line, loc.Column)
}
