package token

import (
	"bytes"
	"fmt"
	"strings"
)

// Pos is just source code position.
//
// Internally it is zero-origin offset of the buffer.
type Pos int

const InvalidPos Pos = -1

// Invalid returns whether p is invalid Pos or not.
func (p Pos) Invalid() bool {
	return p < 0
}

// Position is source code position with file path
// and source code around this position.
type Position struct {
	FilePath string
	Pos, End Pos

	// Line and Column are 0-origin.
	Line, Column       int
	EndLine, EndColumn int

	// Source is source code around this position with line number and cursor
	// for detailed error message.
	Source string
}

func (pos *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", pos.FilePath, pos.Line+1, pos.Column+1)
}

// File is input file with source code.
type File struct {
	FilePath string
	Buffer   string

	lines []Pos
}

// Position returns a new Position from pos and end on this File.
func (f *File) Position(pos, end Pos) *Position {
	line, column := f.ResolvePos(pos)
	endLine, endColumn := f.ResolvePos(end)

	// Calculate source coude around this position.
	var source bytes.Buffer
	switch {
	case pos.Invalid() || end.Invalid():
		break
	case line == endLine:
		lineBuffer := f.Buffer[f.lines[line] : f.lines[line+1]-1]
		count := endColumn - column - 1
		if count < 0 {
			count = 0
		}
		fmt.Fprintf(&source, "%3d:  %s\n", line+1, lineBuffer)
		fmt.Fprintf(&source, "      %s^%s\n", strings.Repeat(" ", column), strings.Repeat("~", count))
	case line < endLine:
		for l := line; l <= endLine; l++ {
			lineBuffer := f.Buffer[f.lines[l] : f.lines[l+1]-1]
			fmt.Fprintf(&source, "%3d:  %s\n", l+1, lineBuffer)
		}
	}

	return &Position{
		FilePath:  f.FilePath,
		Pos:       pos,
		End:       end,
		Line:      line,
		Column:    column,
		EndLine:   endLine,
		EndColumn: endColumn,
		Source:    source.String(),
	}
}

// ResolvePos returns line and column number from pos.
func (f *File) ResolvePos(pos Pos) (line int, column int) {
	line, column = -1, -1

	if pos.Invalid() {
		return
	}

	f.init()
	// TODO: for performance, use binary search instead
	for line = len(f.lines) - 1; line >= 0; line-- {
		linePos := f.lines[line]
		if linePos <= pos {
			column = int(pos - linePos)
			return
		}
	}

	return
}

// init initialize f.lines.
func (f *File) init() {
	if f.lines != nil {
		return
	}

	lines := []Pos{0}
	for i, line := range strings.Split(f.Buffer, "\n") {
		lines = append(lines, Pos(int(lines[i])+len(line)+1))
	}
	f.lines = lines
}
