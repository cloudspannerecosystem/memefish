package token

import (
	"bytes"
	"fmt"
	"unicode"

	"github.com/cloudspannerecosystem/memefish/char"
)

// QuoteSQLString returns quoted string with SQL string escaping.
func QuoteSQLString(s string) string {
	quote := suitableQuote([]byte(s))

	var buf bytes.Buffer
	buf.WriteRune(quote)
	quoteSQLStringContent(quote, s, &buf)
	buf.WriteRune(quote)
	return buf.String()
}

func suitableQuote(b []byte) rune {
	var hasSingle, hasDouble bool
	for _, b := range b {
		switch b {
		case '\'':
			hasSingle = true
		case '"':
			hasDouble = true
		}
	}
	if !hasSingle && hasDouble {
		return '\''
	}
	return '"'
}

// QuoteSQLString returns quoted string with SQL bytes escaping.
func QuoteSQLBytes(bs []byte) string {
	quote := suitableQuote(bs)

	var buf bytes.Buffer
	buf.WriteString("b")
	buf.WriteRune(quote)
	for _, b := range bs {
		q := quoteSingleEscape(quote, rune(b))
		if q != "" {
			buf.WriteString(q)
			continue
		}

		// Note: char.IsPrint(' ') is false
		if b == ' ' || char.IsPrint(b) {
			buf.WriteByte(b)
			continue
		}
		fmt.Fprintf(&buf, `\x%02x`, uint64(b))
	}
	buf.WriteRune(quote)
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
	quoteSQLStringContent('`', s, &buf)
	buf.WriteByte('`')
	return buf.String()
}

func quoteSQLStringContent(quote rune, s string, buf *bytes.Buffer) {
	for _, r := range s {
		q := quoteSingleEscape(quote, r)
		if q != "" {
			buf.WriteString(q)
			continue
		}
		if unicode.IsPrint(r) {
			buf.WriteRune(r)
			continue
		}
		if r > 0xFFFF {
			fmt.Fprintf(buf, `\U%08x`, uint64(r))
		} else {
			fmt.Fprintf(buf, `\u%04x`, uint64(r))
		}
	}
}

func quoteSingleEscape(quote, r rune) string {
	if quote == r {
		return `\` + string(r)
	}

	switch r {
	case '\a':
		return `\a`
	case '\b':
		return `\b`
	case '\f':
		return `\f`
	case '\n':
		return `\n`
	case '\r':
		return `\r`
	case '\t':
		return `\t`
	case '\v':
		return `\v`
	case '?':
		return `\?`
	case '\\':
		return `\\`
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
