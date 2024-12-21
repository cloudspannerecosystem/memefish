package memefish

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish/token"
)

func TestMultiError(t *testing.T) {
	err1 := &Error{
		Message: "error 1",
		Position: &token.Position{
			FilePath:  "foo",
			Pos:       0,
			End:       1,
			Line:      0,
			Column:    0,
			EndLine:   0,
			EndColumn: 1,
			Source:    "  1| a b\n   | ^",
		},
	}
	err2 := &Error{
		Message: "error 2",
		Position: &token.Position{
			FilePath:  "foo",
			Pos:       2,
			End:       3,
			Line:      0,
			Column:    2,
			EndLine:   0,
			EndColumn: 3,
			Source:    "  1| a b\n   |   ^",
		},
	}

	for _, testCase := range []struct {
		list      MultiError
		error     string
		fullError string
	}{
		{
			MultiError{},
			"(0 errors)",
			"",
		},
		{
			MultiError{err1},
			heredoc.Doc(`
				syntax error: foo:1:1: error 1
				  1| a b
				   | ^
			`),
			heredoc.Doc(`
				syntax error: foo:1:1: error 1
				  1| a b
				   | ^
			`),
		},
		{
			MultiError{err1, err2},
			heredoc.Doc(`
				syntax error: foo:1:1: error 1
				  1| a b
				   | ^
				(and 1 other error)
			`),
			heredoc.Doc(`
				syntax error: foo:1:1: error 1
				  1| a b
				   | ^
				syntax error: foo:1:3: error 2
				  1| a b
				   |   ^
			`),
		},
	} {
		if testCase.list.Error() != strings.TrimRight(testCase.error, "\n") {
			t.Errorf("error on MultiError.Error():\n%s", testCase.list.Error())
		}

		if testCase.list.FullError() != testCase.fullError {
			t.Errorf("error on MultiError.FullError():\n%s", testCase.list.FullError())
		}
	}
}
