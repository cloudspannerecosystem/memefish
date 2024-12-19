package token

import (
	"testing"
)

var quoteTestCases = []struct {
	input          string
	str, bytes, id string
}{
	{"foo", `"foo"`, `b"foo"`, "foo"},
	{"if", `"if"`, `b"if"`, "`if`"},
	{"\u0000", `"\x00"`, `b"\x00"`, "`\\x00`"},
	{"\U0010FFFF", `"\U0010ffff"`, `b"\xf4\x8f\xbf\xbf"`, "`\\U0010ffff`"},
	{"a\u2060b", `"a\u2060b"`, `b"a\xe2\x81\xa0b"`, "`a\\u2060b`"},
	{"\a\b\f\n\r\t\v\"'?\\", `"\x07\x08\x0c\n\r\t\x0b\"'?\\"`, `b"\x07\x08\x0c\x0a\x0d\x09\x0b\"'?\\"`, "`\\x07\\x08\\x0c\\n\\r\\t\\x0b\"'?\\\\`"},
	{"`", "\"`\"", "b\"`\"", "`\\``"},
}

func TestQuote(t *testing.T) {
	for _, tc := range quoteTestCases {
		s := QuoteSQLString(tc.input)
		if tc.str != s {
			t.Errorf("QuoteSQLString: %q (want) != %q (got)", tc.str, s)
		}

		b := QuoteSQLBytes([]byte(tc.input))
		if tc.bytes != b {
			t.Errorf("QuoteSQLBytes: %q (want) != %q (got)", tc.bytes, b)
		}

		id := QuoteSQLIdent(tc.input)
		if tc.id != id {
			t.Errorf("QuoteSQLIdent: %q (want) != %q (got)", tc.id, id)
		}
	}
}
