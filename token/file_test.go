package token

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
)

func stripMargin(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if len(line) > 0 && line[0] == '|' {
			lines[i] = line[1:]
		} else {
			lines[i] = line
		}
	}
	return strings.Join(lines, "\n")
}

var file = &File{
	FilePath: "test",
	Buffer: heredoc.Doc(`
		select 1 union all
		select 2
	`),
}

var positionTestCases = []struct {
	pos, end           Pos
	line, column       int
	endLine, endColumn int
	source             string
}{
	{
		pos: -1, end: -1,
		line: -1, column: -1,
		endLine: -1, endColumn: -1,
		source: "",
	},
	{
		pos: -1, end: 0,
		line: -1, column: -1,
		endLine: 0, endColumn: 0,
		source: "",
	},
	{
		pos: 0, end: -1,
		line: 0, column: 0,
		endLine: -1, endColumn: -1,
		source: "",
	},
	{
		pos: 0, end: 0,
		line: 0, column: 0,
		endLine: 0, endColumn: 0,
		source: stripMargin(heredoc.Doc(`
			|  1:  select 1 union all
			|      ^
		`)),
	},
	{
		pos: 0, end: 1,
		line: 0, column: 0,
		endLine: 0, endColumn: 1,
		source: stripMargin(heredoc.Doc(`
			|  1:  select 1 union all
			|      ^
		`)),
	},
	{
		pos: 0, end: 6,
		line: 0, column: 0,
		endLine: 0, endColumn: 6,
		source: stripMargin(heredoc.Doc(`
			|  1:  select 1 union all
			|      ^~~~~~
		`)),
	},
	{
		pos: 9, end: 18,
		line: 0, column: 9,
		endLine: 0, endColumn: 18,
		source: stripMargin(heredoc.Doc(`
			|  1:  select 1 union all
			|               ^~~~~~~~~
		`)),
	},
	{
		pos: 18, end: 19,
		line: 0, column: 18,
		endLine: 1, endColumn: 0,
		source: stripMargin(heredoc.Doc(`
			|  1:  select 1 union all
			|  2:  select 2
		`)),
	},
}

func TestPosition(t *testing.T) {
	for _, tc := range positionTestCases {
		pos := file.Position(tc.pos, tc.end)

		if tc.line != pos.Line {
			t.Errorf("Line: %d (want) != %d (got)", tc.line, pos.Line)
		}
		if tc.column != pos.Column {
			t.Errorf("Column: %d (want) != %d (got)", tc.column, pos.Column)
		}
		if tc.endLine != pos.EndLine {
			t.Errorf("EndLine: %d (want) != %d (got)", tc.endLine, pos.EndLine)
		}
		if tc.endColumn != pos.EndColumn {
			t.Errorf("EndColumn: %d (want) != %d (got)", tc.endColumn, pos.EndColumn)
		}
		if strings.TrimRight(tc.source, "\n") != pos.Source {
			t.Errorf("Source:\n-- want --\n%s\n-- got --\n%s\n", tc.source, pos.Source)
		}
	}
}
