package memefish_test

import (
	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"testing"
)

func TestSplitRawStatements(t *testing.T) {
	for _, test := range []struct {
		desc  string
		input string
		errRe *regexp.Regexp
		want  []*memefish.RawStatement
	}{
		// SplitRawStatements treats only lexical structures, so the test cases can be invalid statements.
		{desc: "empty input", input: "", want: []*memefish.RawStatement{{Statement: ""}}},
		{desc: "single statement ends with semicolon", input: `SELECT "123";`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT "123"`, End: token.Pos(12)},
			}},
		{desc: "single statement ends with EOF", input: `SELECT "123"`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT "123"`, End: token.Pos(12)},
			}},
		{desc: "two statement ends with semicolon", input: `SELECT "123"; SELECT "456";`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT "123"`, End: token.Pos(12)},
				{Statement: `SELECT "456"`, Pos: token.Pos(14), End: token.Pos(26)},
			}},
		{desc: "two statement ends with EOF", input: `SELECT "123"; SELECT "456"`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT "123"`, End: token.Pos(12)},
				{Statement: `SELECT "456"`, Pos: token.Pos(14), End: token.Pos(26)},
			}},
		{desc: "second statement is empty", input: `SELECT 1; ;`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT 1`, End: token.Pos(8)},
				{Statement: ``, Pos: token.Pos(10), End: token.Pos(10)},
			}},
		{desc: "two statement with new lines", input: "SELECT 1;\n SELECT 2;\n",
			want: []*memefish.RawStatement{
				{Statement: "SELECT 1", End: token.Pos(8)},
				{Statement: "SELECT 2", Pos: token.Pos(11), End: token.Pos(19)},
			}},
		{desc: "single statement with line comment", input: `SELECT 1//
`, want: []*memefish.RawStatement{
			{Statement: "SELECT 1//\n", End: token.Pos(11)},
		}},
		{desc: "semicolon in line comment", input: "SELECT 1 //;\n + 2",
			want: []*memefish.RawStatement{
				{Statement: "SELECT 1 //;\n + 2", End: token.Pos(17)},
			}},
		{desc: "semicolon in multi-line comment", input: "SELECT 1 /*;\n*/ + 2",
			want: []*memefish.RawStatement{
				{Statement: "SELECT 1 /*;\n*/ + 2", End: token.Pos(19)},
			}},
		{desc: "semicolon in double-quoted string", input: `SELECT "1;2;3";`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT "1;2;3"`, End: token.Pos(14)},
			}},
		{desc: "semicolon in single-quoted string", input: `SELECT '1;2;3';`,
			want: []*memefish.RawStatement{
				{Statement: `SELECT '1;2;3'`, End: token.Pos(14)},
			}},
		{desc: "semicolon in back-quote", input: "SELECT `1;2;3`;",
			want: []*memefish.RawStatement{
				{Statement: "SELECT `1;2;3`", End: token.Pos(14)},
			}},
	} {
		t.Run(test.desc, func(t *testing.T) {
			stmts, err := memefish.SplitRawStatements("", test.input)
			if err != nil {
				if test.errRe == nil {
					t.Errorf("should success, but %v", err)
					return
				}
				if !test.errRe.MatchString(err.Error()) {
					t.Errorf("error message should match %q, but %q", test.errRe, err)
					return
				}
			}
			if err == nil && test.errRe != nil {
				t.Errorf("success, but should fail %q", test.errRe)
				return
			}
			if diff := cmp.Diff(stmts, test.want); diff != "" {
				t.Errorf("differs: %v", diff)
				return
			}
		})
	}
}
