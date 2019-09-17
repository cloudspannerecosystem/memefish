package token

import (
	"testing"
)

var quoteTestCases = []struct {
	input          string
	str, bytes, id string
}{
	{"foo", `"foo"`, `B"foo"`, "foo"},
	{"\u0000", `"\u0000"`, `B"\x00"`, "`\\u0000`"},
	{"\U0010FFFF", `"\U0010FFFF"`, `B"\xF4\x8F\xBF\xBF"`, "`\\U0010FFFF`"},
	{"a\u2060b", `"a\u2060b"`, `B"a\xE2\x81\xA0b"`, "`a\\u2060b`"},
	{"\a\b\f\n\r\t\v\"'?\\", `"\a\b\f\n\r\t\v\"\'\?\\"`, `B"\a\b\f\n\r\t\v\"\'\?\\"`, "`\\a\\b\\f\\n\\r\\t\\v\\\"\\'\\?\\\\`"},
	{"`", "\"\\`\"", "B\"\\`\"", "`\\``"},
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
