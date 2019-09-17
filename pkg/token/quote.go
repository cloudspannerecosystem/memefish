package token

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/MakeNowJust/memefish/pkg/char"
)

// QuoteSQLString returns quoted string with SQL string escaping.
func QuoteSQLString(s string) string {
	var buf bytes.Buffer
	buf.WriteByte('"')
	quoteSQLStringContent(s, &buf)
	buf.WriteByte('"')
	return buf.String()
}

// QuoteSQLString returns quoted string with SQL bytes escaping.
func QuoteSQLBytes(bs []byte) string {
	var buf bytes.Buffer
	buf.WriteString("B\"")
	for _, b := range bs {
		q := quoteSingleEscape(rune(b))
		if q != "" {
			buf.WriteString(q)
			continue
		}
		if char.IsPrint(b) {
			buf.WriteByte(b)
			continue
		}
		fmt.Fprintf(&buf, "\\x%02X", uint64(b))
	}
	buf.WriteRune('"')
	return buf.String()
}

// QuoteSQLString returns quoted string with SQL bytes escaping if needed,
// otherwise it returns the input string.
func QuoteSQLIdent(s string) string {
	if !needQuoteSQLIdent(s) {
		return s
	}

	var buf bytes.Buffer
	buf.WriteByte('`')
	quoteSQLStringContent(s, &buf)
	buf.WriteByte('`')
	return buf.String()
}

func quoteSQLStringContent(s string, buf *bytes.Buffer) {
	for _, r := range s {
		q := quoteSingleEscape(r)
		if q != "" {
			buf.WriteString(q)
			continue
		}
		if unicode.IsPrint(r) {
			buf.WriteRune(r)
			continue
		}
		if r > 0xFFFF {
			fmt.Fprintf(buf, "\\U%08X", uint64(r))
		} else {
			fmt.Fprintf(buf, "\\u%04X", uint64(r))
		}
	}
}

func quoteSingleEscape(r rune) string {
	switch r {
	case '\a':
		return "\\a"
	case '\b':
		return "\\b"
	case '\f':
		return "\\f"
	case '\n':
		return "\\n"
	case '\r':
		return "\\r"
	case '\t':
		return "\\t"
	case '\v':
		return "\\v"
	case '"':
		return "\\\""
	case '\'':
		return "\\'"
	case '`':
		return "\\`"
	case '?':
		return "\\?"
	case '\\':
		return "\\\\"
	}
	return ""
}

func needQuoteSQLIdent(s string) bool {
	// When s is keyword, it should be quoted.
	if IsKeyword(s) {
		return true
	}

	// Then, check s can be parsed as TokenIdent without backquoted.
	if !char.IsIdentStart(s[0]) {
		return true
	}
	for i := 0; i < len(s); i++ {
		if !char.IsIdentPart(s[i]) {
			return true
		}
	}

	// After passing all check, then s doesn't need to be quoted.
	return false
}
