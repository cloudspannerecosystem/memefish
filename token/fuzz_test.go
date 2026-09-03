package token_test

import (
	"testing"
	"unicode/utf8"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/token"
)

// lexSingleToken lexes s and returns its sole token, failing if s does not
// lex as exactly one token followed by EOF.
func lexSingleToken(t *testing.T, s string) token.Token {
	t.Helper()

	l := &memefish.Lexer{File: &token.File{Buffer: s}}
	if err := l.NextToken(); err != nil {
		t.Fatalf("failed to lex %q: %v", s, err)
	}
	tok := l.Token
	if err := l.NextToken(); err != nil {
		t.Fatalf("failed to lex %q after the first token: %v", s, err)
	}
	if l.Token.Kind != token.TokenEOF {
		t.Fatalf("expected a single token in %q, but found trailing %s token", s, l.Token.Kind)
	}
	return tok
}

func FuzzQuoteSQLString(f *testing.F) {
	for _, seed := range []string{"", "hello", `foo "bar" 'baz'`, "\\n\\", "\n\r\t\a\b\f\v", "日本語", "a\x00b", "\U0001F600"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// A SQL STRING value is Unicode; invalid UTF-8 cannot round-trip.
		if !utf8.ValidString(s) {
			t.Skip()
		}

		q := token.QuoteSQLString(s)
		tok := lexSingleToken(t, q)
		if tok.Kind != token.TokenString {
			t.Fatalf("QuoteSQLString(%q) = %q lexed as %s, not a string literal", s, q, tok.Kind)
		}
		if tok.AsString != s {
			t.Fatalf("QuoteSQLString(%q) = %q lexed back as %q", s, q, tok.AsString)
		}
	})
}

func FuzzQuoteSQLBytes(f *testing.F) {
	for _, seed := range [][]byte{{}, []byte("hello"), []byte(`"'`), {0x00, 0xff, 0x7f}, []byte("\\x00"), []byte("日本語")} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, bs []byte) {
		q := token.QuoteSQLBytes(bs)
		tok := lexSingleToken(t, q)
		if tok.Kind != token.TokenBytes {
			t.Fatalf("QuoteSQLBytes(%q) = %q lexed as %s, not a bytes literal", bs, q, tok.Kind)
		}
		if tok.AsString != string(bs) {
			t.Fatalf("QuoteSQLBytes(%q) = %q lexed back as %q", bs, q, []byte(tok.AsString))
		}
	})
}

func FuzzQuoteSQLIdent(f *testing.F) {
	for _, seed := range []string{"", "foo", "SELECT", "foo bar", "1foo", "_foo", "foo`bar", "日本語", "a\nb"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		// An identifier is non-empty Unicode text; QuoteSQLIdent formats without
		// validating, so invalid identifiers cannot round-trip through the lexer.
		if s == "" || !utf8.ValidString(s) {
			t.Skip()
		}

		q := token.QuoteSQLIdent(s)
		tok := lexSingleToken(t, q)
		if tok.Kind != token.TokenIdent {
			t.Fatalf("QuoteSQLIdent(%q) = %q lexed as %s, not an identifier", s, q, tok.Kind)
		}
		if tok.AsString != s {
			t.Fatalf("QuoteSQLIdent(%q) = %q lexed back as %q", s, q, tok.AsString)
		}
	})
}
