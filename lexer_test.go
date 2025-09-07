package memefish

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	. "github.com/cloudspannerecosystem/memefish/token"
)

// Keep same order https://github.com/google/zetasql/blob/master/zetasql/parser/flex_tokenizer.l
var symbols = []string{
	"(",
	"[",
	"{",
	")",
	"]",
	"}",
	"*",
	",",
	"=",
	"+=",
	"-=",
	"!=",
	"<=",
	"<<",
	"=>",
	"->",
	"<",
	">",
	">=",
	"||",
	"|",
	"^",
	"&",
	"+",
	"-",
	"/",
	"~",
	"?",
	"!",
	"%",
	"|>",
	"@",
	"@@",
	".",
	":",
	"\\",
	";",
	"$",
	"<>", // <> is not a valid token in ZetaSQL, but it is a token in memefish
	">>", // >> is not a valid token in ZetaSQL, but it is a token in memefish.
}

var lexerTestCases = []struct {
	source string
	tokens []*Token
}{
	// Spaces
	{"  0", []*Token{{Kind: "<int>", Space: "  ", Raw: "0", Base: 10}}},
	// Comment
	{"# foo", nil},
	{"# foo\n0", []*Token{{Kind: "<int>", Space: "", Raw: "0", Base: 10, Comments: []TokenComment{{Space: "", Raw: "# foo\n", Pos: 0, End: 6}}}}},
	{"-- foo\n0", []*Token{{Kind: "<int>", Space: "", Raw: "0", Base: 10, Comments: []TokenComment{{Space: "", Raw: "-- foo\n", Pos: 0, End: 7}}}}},
	{"/* foo */ 0", []*Token{{Kind: "<int>", Space: " ", Raw: "0", Base: 10, Comments: []TokenComment{{Space: "", Raw: "/* foo */", Pos: 0, End: 9}}}}},
	{"-- aaa\n-- bbb\n/* foo */ 0", []*Token{{Kind: "<int>", Space: " ", Raw: "0", Base: 10, Comments: []TokenComment{
		{Space: "", Raw: "-- aaa\n", Pos: 0, End: 7},
		{Space: "", Raw: "-- bbb\n", Pos: 7, End: 14},
		{Space: "", Raw: "/* foo */", Pos: 14, End: 23},
	}}}},
	// TokenInt
	{"0", []*Token{{Kind: TokenInt, Raw: "0", Base: 10}}},
	{"1", []*Token{{Kind: TokenInt, Raw: "1", Base: 10}}},
	{"123", []*Token{{Kind: TokenInt, Raw: "123", Base: 10}}},
	{"+123", []*Token{{Kind: "+", Raw: "+"}, {Kind: TokenInt, Raw: "123", Base: 10}}},
	{"-123", []*Token{{Kind: "-", Raw: "-"}, {Kind: TokenInt, Raw: "123", Base: 10}}},
	{"9223372036854775807", []*Token{{Kind: TokenInt, Raw: "9223372036854775807", Base: 10}}},
	{"-9223372036854775808", []*Token{{Kind: "-", Raw: "-"}, {Kind: TokenInt, Raw: "9223372036854775808", Base: 10}}},
	{"0123", []*Token{{Kind: TokenInt, Raw: "0123", Base: 10}}}, // TODO: fix base
	{"0xbeaf", []*Token{{Kind: TokenInt, Raw: "0xbeaf", Base: 16}}},
	{"0XBEAF", []*Token{{Kind: TokenInt, Raw: "0XBEAF", Base: 16}}},
	// TokenFloat
	{"1.2", []*Token{{Kind: TokenFloat, Raw: "1.2"}}},
	{"+1.2", []*Token{{Kind: "+", Raw: "+"}, {Kind: TokenFloat, Raw: "1.2"}}},
	{"-1.2", []*Token{{Kind: "-", Raw: "-"}, {Kind: TokenFloat, Raw: "1.2"}}},
	{".1", []*Token{{Kind: TokenFloat, Raw: ".1"}}},
	{"00.1", []*Token{{Kind: TokenFloat, Raw: "00.1"}}},
	{"1.", []*Token{{Kind: TokenFloat, Raw: "1."}}},
	{"1e1", []*Token{{Kind: TokenFloat, Raw: "1e1"}}},
	{"1E1", []*Token{{Kind: TokenFloat, Raw: "1E1"}}},
	{"1e+1", []*Token{{Kind: TokenFloat, Raw: "1e+1"}}},
	{"1e-1", []*Token{{Kind: TokenFloat, Raw: "1e-1"}}},
	// TokenParam
	{"@foo", []*Token{{Kind: TokenParam, Raw: "@foo", AsString: "foo"}}},
	{"@foo.1", []*Token{{Kind: TokenParam, Raw: "@foo", AsString: "foo"}, {Kind: ".", Raw: "."}, {Kind: TokenIdent, Raw: "1", AsString: "1"}}},
	// TokenIdent
	{"foo", []*Token{{Kind: TokenIdent, Raw: "foo", AsString: "foo"}}},
	{"foo_bar", []*Token{{Kind: TokenIdent, Raw: "foo_bar", AsString: "foo_bar"}}},
	{"foo.1", []*Token{{Kind: TokenIdent, Raw: "foo", AsString: "foo"}, {Kind: ".", Raw: "."}, {Kind: TokenIdent, Raw: "1", AsString: "1"}}},
	{"foo.*", []*Token{{Kind: TokenIdent, Raw: "foo", AsString: "foo"}, {Kind: ".", Raw: "."}, {Kind: "*", Raw: "*"}}},
	{"`select`", []*Token{{Kind: TokenIdent, Raw: "`select`", AsString: "select"}}},
	{"`select`.1", []*Token{{Kind: TokenIdent, Raw: "`select`", AsString: "select"}, {Kind: ".", Raw: "."}, {Kind: TokenIdent, Raw: "1", AsString: "1"}}},
	{"].1", []*Token{{Kind: "]", Raw: "]"}, {Kind: ".", Raw: "."}, {Kind: TokenIdent, Raw: "1", AsString: "1"}}},
	{").1", []*Token{{Kind: ")", Raw: ")"}, {Kind: ".", Raw: "."}, {Kind: TokenIdent, Raw: "1", AsString: "1"}}},
	{"`foo\\u0031`", []*Token{{Kind: TokenIdent, Raw: "`foo\\u0031`", AsString: "foo1"}}},
	{"BR", []*Token{{Kind: TokenIdent, Raw: "BR", AsString: "BR"}}},
	{`R "foo"`, []*Token{{Kind: TokenIdent, Raw: "R", AsString: "R"}, {Kind: TokenString, Space: " ", Raw: `"foo"`, AsString: "foo"}}},
	// TokenString
	{`""`, []*Token{{Kind: TokenString, Raw: `""`, AsString: ""}}},
	{`''`, []*Token{{Kind: TokenString, Raw: `''`, AsString: ""}}},
	{`"foo"`, []*Token{{Kind: TokenString, Raw: `"foo"`, AsString: "foo"}}},
	{`'foo'`, []*Token{{Kind: TokenString, Raw: `'foo'`, AsString: "foo"}}},
	{`"""foo\nbar"""`, []*Token{{Kind: TokenString, Raw: `"""foo\nbar"""`, AsString: "foo\nbar"}}},
	{"'''foo\nbar'''", []*Token{{Kind: TokenString, Raw: "'''foo\nbar'''", AsString: "foo\nbar"}}},
	{`"\a\b\f\n\r\t\v\\\?\"\'"`, []*Token{{Kind: TokenString, Raw: `"\a\b\f\n\r\t\v\\\?\"\'"`, AsString: "\a\b\f\n\r\t\v\\?\"'"}}},
	{"'\\`'", []*Token{{Kind: TokenString, Raw: "'\\`'", AsString: "`"}}},
	{`"\061\x31\X31\u0031\U00000031"`, []*Token{{Kind: TokenString, Raw: `"\061\x31\X31\u0031\U00000031"`, AsString: "11111"}}},
	{`"\xff"`, []*Token{{Kind: TokenString, Raw: `"\xff"`, AsString: "\xff"}}},
	{`R"\\"`, []*Token{{Kind: TokenString, Raw: `R"\\"`, AsString: "\\\\"}}},
	{`R'\\'`, []*Token{{Kind: TokenString, Raw: `R'\\'`, AsString: "\\\\"}}},
	{`r"\\"`, []*Token{{Kind: TokenString, Raw: `r"\\"`, AsString: "\\\\"}}},
	{`r'\\'`, []*Token{{Kind: TokenString, Raw: `r'\\'`, AsString: "\\\\"}}},
	{`R"""\\"""`, []*Token{{Kind: TokenString, Raw: `R"""\\"""`, AsString: "\\\\"}}},
	{`R'''\\'''`, []*Token{{Kind: TokenString, Raw: `R'''\\'''`, AsString: "\\\\"}}},
	// ByteString
	{`B"foo"`, []*Token{{Kind: TokenBytes, Raw: `B"foo"`, AsString: "foo"}}},
	{`B'foo'`, []*Token{{Kind: TokenBytes, Raw: `B'foo'`, AsString: "foo"}}},
	{`b"foo"`, []*Token{{Kind: TokenBytes, Raw: `b"foo"`, AsString: "foo"}}},
	{`b'foo'`, []*Token{{Kind: TokenBytes, Raw: `b'foo'`, AsString: "foo"}}},
	{`B"""foo"""`, []*Token{{Kind: TokenBytes, Raw: `B"""foo"""`, AsString: "foo"}}},
	{`B'''foo'''`, []*Token{{Kind: TokenBytes, Raw: `B'''foo'''`, AsString: "foo"}}},
	{`B"\a\b\f\n\r\t\v\\\?\"\'"`, []*Token{{Kind: TokenBytes, Raw: `B"\a\b\f\n\r\t\v\\\?\"\'"`, AsString: "\a\b\f\n\r\t\v\\?\"'"}}},
	{`RB"foo"`, []*Token{{Kind: TokenBytes, Raw: `RB"foo"`, AsString: "foo"}}},
	{"RB'''foo\nbar'''", []*Token{{Kind: TokenBytes, Raw: "RB'''foo\nbar'''", AsString: "foo\nbar"}}},
	{`rb"foo"`, []*Token{{Kind: TokenBytes, Raw: `rb"foo"`, AsString: "foo"}}},
	{`BR"foo"`, []*Token{{Kind: TokenBytes, Raw: `BR"foo"`, AsString: "foo"}}},
}

var lexerWrongTestCase = []struct {
	source  string
	pos     Pos
	end     Pos
	message string
}{
	{"\b", 0, 0, "illegal input character: '\\b'"},
	{`"foo`, 0, 4, "unclosed string literal"},
	{`R"foo`, 1, 5, "unclosed raw string literal"},
	{"'foo\n", 0, 4, "unclosed string literal: newline appears in non triple-quoted"},
	{"R'foo\n", 1, 5, "unclosed raw string literal: newline appears in non triple-quoted"},
	{"R'foo\\", 5, 6, "invalid escape sequence: \\<eof>"},
	{`"\400"`, 1, 3, "invalid escape sequence: \\4"},
	{`"\3xx"`, 1, 4, "invalid escape sequence: octal escape sequence must be follwed by 3 octal digits"},
	{`"\xZZ"`, 1, 4, "invalid escape sequence: hex escape sequence must be follwed by 2 hex digits"},
	{`"\XZZ"`, 1, 4, "invalid escape sequence: hex escape sequence must be follwed by 2 hex digits"},
	{`B"\u0031"`, 2, 4, "invalid escape sequence: \\u is not allowed in bytes literal"},
	{`B"\U00000031"`, 2, 4, "invalid escape sequence: \\U is not allowed in bytes literal"},
	{`B"\U00000031"`, 2, 4, "invalid escape sequence: \\U is not allowed in bytes literal"},
	{`"\UFFFFFFFF"`, 1, 11, "invalid escape sequence: invalid code point: U+FFFFFFFF"},
	{"``", 0, 2, "invalid empty identifier"},
	{"1from", 1, 1, "number literal cannot follow identifier without any spaces"},
	{`'''0`, 0, 4, "unclosed triple-quoted string literal"},
	{`/*`, 0, 2, "unclosed comment"},
}

func testLexer(t *testing.T, source string, tokens []*Token) {
	t.Helper()
	l := &Lexer{
		File: &File{FilePath: "[test]", Buffer: source},
	}
	for _, t2 := range tokens {
		err := l.NextToken()
		if err != nil {
			t.Errorf("error on lexer: %v", err)
			return
		}
		opts := []cmp.Option{
			cmpopts.IgnoreFields(Token{}, "Pos", "End"),
		}
		if diff := cmp.Diff(&l.Token, t2, opts...); diff != "" {
			t.Errorf("(-got, +want)\n%s", diff)
		}
	}
	err := l.NextToken()
	if err != nil {
		t.Errorf("error on lexer: %v", err)
		return
	}
	if l.Token.Kind != TokenEOF {
		t.Errorf("expected EOF, but: %#v", &l.Token)
		return
	}
}

func TestLexer(t *testing.T) {
	for _, s := range Keywords {
		t.Run(fmt.Sprintf("keyword/%q", string(s)), func(t *testing.T) {
			testLexer(t, string(s), []*Token{{Kind: s, Raw: string(s)}})
		})
		l := strings.ToLower(string(s))
		t.Run(fmt.Sprintf("keyword/%q", l), func(t *testing.T) {
			testLexer(t, l, []*Token{{Kind: s, Raw: l}})
		})
	}

	for _, s := range symbols {
		t.Run(fmt.Sprintf("symbol/%q", s), func(t *testing.T) {
			testLexer(t, s, []*Token{{Kind: TokenKind(s), Raw: s}})
		})
	}

	for _, tc := range lexerTestCases {
		t.Run(fmt.Sprintf("testcase/%q", tc.source), func(t *testing.T) {
			testLexer(t, tc.source, tc.tokens)
		})
	}
}

func TestLexerWrong(t *testing.T) {
	for _, tc := range lexerWrongTestCase {
		t.Run(fmt.Sprintf("testcase/%q", tc.source), func(t *testing.T) {
			l := &Lexer{
				File: &File{FilePath: "[test]", Buffer: tc.source},
			}
			var err error
			for l.Token.Kind != TokenEOF {
				err = l.NextToken()
				if err != nil {
					break
				}
			}
			if err == nil {
				t.Errorf("unexpected EOF")
				return
			}
			if e, ok := err.(*Error); ok {
				if e.Message != tc.message {
					t.Errorf("expected error message: %q, but: %q", tc.message, e.Message)
				}
				if e.Position.Pos != tc.pos {
					t.Errorf("expected error position (pos): %v, but: %v", tc.pos, e.Position.Pos)
				}
				if e.Position.End != tc.end {
					t.Errorf("expected error position (end): %v, but: %v", tc.end, e.Position.End)
				}
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestLexerWrongNoError(t *testing.T) {
	for _, tc := range lexerWrongTestCase {
		t.Run(fmt.Sprintf("testcase/%q", tc.source), func(t *testing.T) {
			l := &Lexer{
				File: &File{FilePath: "[test]", Buffer: tc.source},
			}
			hasBad := false
			for l.Token.Kind != TokenEOF {
				l.nextToken(true)
				if l.Token.Kind == TokenBad {
					hasBad = true
				}
			}
			if !hasBad {
				t.Errorf("expected <bad>")
			}
		})
	}
}
