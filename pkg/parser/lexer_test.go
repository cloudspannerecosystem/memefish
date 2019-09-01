package parser

import (
	"fmt"
	"strings"
	"testing"
)

var symbols = []string{
	".",
	",",
	";",
	"(",
	")",
	"{",
	"}",
	"[",
	"]",
	"@",
	"~",
	"+",
	"-",
	"*",
	"/",
	"&",
	"^",
	"|",
	"=",
	"<",
	"<<",
	"<=",
	"<>",
	">",
	">>",
	">=",
	"!=",
}

var lexerTestCases = []struct {
	source string
	tokens []*Token
}{
	// TokenInt
	{"0", []*Token{{Kind: TokenInt, Raw: "0"}}},
	{"1", []*Token{{Kind: TokenInt, Raw: "1"}}},
	{"123", []*Token{{Kind: TokenInt, Raw: "123"}}},
	{"+123", []*Token{{Kind: TokenInt, Raw: "+123"}}},
	{"-123", []*Token{{Kind: TokenInt, Raw: "-123"}}},
	{"9223372036854775807", []*Token{{Kind: TokenInt, Raw: "9223372036854775807"}}},
	{"-9223372036854775808", []*Token{{Kind: TokenInt, Raw: "-9223372036854775808"}}},
	{"0123", []*Token{{Kind: TokenInt, Raw: "0123"}}},
	{"0xbeaf", []*Token{{Kind: TokenInt, Raw: "0xbeaf"}}},
	{"0XBEAF", []*Token{{Kind: TokenInt, Raw: "0XBEAF"}}},
	// TokenFloat
	{"1.2", []*Token{{Kind: TokenFloat, Raw: "1.2"}}},
	{"+1.2", []*Token{{Kind: TokenFloat, Raw: "+1.2"}}},
	{"-1.2", []*Token{{Kind: TokenFloat, Raw: "-1.2"}}},
	{".1", []*Token{{Kind: TokenFloat, Raw: ".1"}}},
	{"00.1", []*Token{{Kind: TokenFloat, Raw: "00.1"}}},
	{"1.", []*Token{{Kind: TokenFloat, Raw: "1."}}},
	{"1e1", []*Token{{Kind: TokenFloat, Raw: "1e1"}}},
	{"1E1", []*Token{{Kind: TokenFloat, Raw: "1E1"}}},
	{"1e+1", []*Token{{Kind: TokenFloat, Raw: "1e+1"}}},
	{"1e-1", []*Token{{Kind: TokenFloat, Raw: "1e-1"}}},
}

func nextToken(l *Lexer) (tok *Token, err error) {
	defer func() {
		if r := recover(); r != nil {
			tok = nil
			err = r.(error)
		}
	}()

	l.NextToken()
	tok = &l.Token
	return
}

func tokenEqual(t1, t2 *Token) bool {
	if t1.Kind != t2.Kind || t1.Raw != t2.Raw {
		return false
	}

	switch t1.Kind {
	case TokenParam, TokenIdent, TokenString, TokenBytes:
		return t1.AsString == t2.AsString
	}

	return true
}

func testLexer(t *testing.T, source string, tokens []*Token) {
	t.Helper()
	l := &Lexer{
		File: NewFile("[test]", source),
	}
	for _, t2 := range tokens {
		t1, err := nextToken(l)
		if err != nil {
			t.Errorf("error on lexer: %v", err)
			return
		}
		if !tokenEqual(t1, t2) {
			t.Errorf("%#v != %#v", t1, t2)
			return
		}
	}
	t1, err := nextToken(l)
	if err != nil {
		t.Errorf("error on lexer: %v", err)
		return
	}
	if t1.Kind != TokenEOF {
		t.Errorf("expected EOF, but: %#v", t1)
		return
	}
}

func TestLexer(t *testing.T) {
	for _, s := range keywords {
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
