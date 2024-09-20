package memefish_test

import (
	"github.com/cloudspannerecosystem/memefish"
	"github.com/google/go-cmp/cmp"
	"regexp"
	"testing"
)

func TestSeparateRawStatements(t *testing.T) {
	for _, test := range []struct {
		desc  string
		input string
		errRe *regexp.Regexp
		want  []string
	}{
		// SeparateRawStatements treats only lexical structures, so the test cases can be invalid statements.
		{desc: "empty input", input: "", want: []string{""}},
		{desc: "single statement ", input: `SELECT "123";`, want: []string{`SELECT "123"`}},
		{desc: "two statement", input: `SELECT "123"; SELECT "456";`, want: []string{`SELECT "123"`, `SELECT "456"`}},
		{desc: "second statement is empty", input: `SELECT 1; ;`, want: []string{`SELECT 1`, ``}},
		{desc: "two statement", input: "SELECT 1;\n SELECT 2;\n", want: []string{"SELECT 1", "SELECT 2"}},
		{desc: "single statement with line comment", input: `SELECT 1//
`, want: []string{"SELECT 1//\n"}},
		{desc: "semicolon in line comment", input: "SELECT 1 //;\n + 2", want: []string{"SELECT 1 //;\n + 2"}},
		{desc: "semicolon in multi-line comment", input: "SELECT 1 /*;\n*/ + 2", want: []string{"SELECT 1 /*;\n*/ + 2"}},
		{desc: "semicolon in double-quoted string", input: `SELECT "1;2;3";`, want: []string{`SELECT "1;2;3"`}},
		{desc: "semicolon in single-quoted string", input: `SELECT '1;2;3';`, want: []string{`SELECT '1;2;3'`}},
		{desc: "semicolon in back-quote", input: "SELECT `1;2;3`;", want: []string{"SELECT `1;2;3`"}},
		// $` may become a valid token in the future, but it's reasonable to check its current behavior.
		{desc: "unknown token", input: "SELECT $;", errRe: regexp.MustCompile(`illegal input character: '\$'`)},
	} {
		t.Run(test.desc, func(t *testing.T) {
			stmts, err := memefish.SeparateRawStatements("", test.input)
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
