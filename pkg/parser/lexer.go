package parser

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Lexer struct {
	Buffer               []byte
	FilePath             string
	Offset, Line, Column int
	Token                Token
}

type Token struct {
	Kind     TokenKind
	Space    []byte
	Raw      []byte
	AsString []byte  // available for TokenIdent, TokenString and TokenBytes
	AsInt    int64   // available for TokenInt
	AsFloat  float64 // available for TokenFloat
	Loc      *Location
	EndLoc   *Location
}

type TokenKind string

const (
	TokenEOF    TokenKind = "eof"
	TokenIdent  TokenKind = "ident"
	TokenParam  TokenKind = "param"
	TokenInt    TokenKind = "int"
	TokenFloat  TokenKind = "float"
	TokenString TokenKind = "string"
	TokenBytes  TokenKind = "bytes"
)

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isHexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func isOctalDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func isIdentPart(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func toUpper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		return c - 'a' + 'A'
	}
	return c
}

func (l *Lexer) NextToken() {
	// Skips spaces.
	l.Token = Token{}
	i := l.Offset
	l.skipSpaces()
	l.Token.Space = l.Buffer[i:l.Offset]

	// Reads the next token.
	l.Token.Loc = l.loc()
	i = l.Offset
	l.nextToken()
	l.Token.Raw = l.Buffer[i:l.Offset]
	l.Token.EndLoc = l.loc()
}

func (l *Lexer) nextToken() {
	if l.eof() {
		l.Token.Kind = TokenEOF
		return
	}

	switch l.peek(0) {
	case '(', ')', '{', '}', ':', ',', '[', ']', '~', '*', '/', '&', '^', '|', '=':
		l.Token.Kind = TokenKind([]byte{l.next()})
		return
	case '+', '-':
		if l.peekOk(1) && (isDigit(l.peek(1)) || l.peekIs(1, '.')) {
			l.nextNumber()
		} else {
			l.Token.Kind = TokenKind([]byte{l.next()})
		}
		return
	case '.':
		if l.peekOk(1) && isDigit(l.peek(1)) {
			l.nextNumber()
		} else {
			l.Token.Kind = "."
		}
		return
	case '<':
		switch {
		case l.peekIs(1, '<'):
			l.nextN(2)
			l.Token.Kind = "<<"
		case l.peekIs(1, '='):
			l.nextN(2)
			l.Token.Kind = "<="
		case l.peekIs(1, '>'):
			l.nextN(2)
			l.Token.Kind = "<>"
		default:
			l.next()
			l.Token.Kind = "<"
		}
		return
	case '>':
		switch {
		case l.peekIs(1, '>'):
			l.nextN(2)
			l.Token.Kind = ">>"
		case l.peekIs(1, '='):
			l.nextN(2)
			l.Token.Kind = ">="
		default:
			l.next()
			l.Token.Kind = ">"
		}
		return
	case '!':
		if l.peekIs(1, '=') {
			l.nextN(2)
			l.Token.Kind = "!="
			return
		}
	case '@':
		i := 1
		for l.peekOk(i) {
			if !isIdentPart(l.peek(i)) {
				break
			}
			i++
		}
		if i > 1 {
			l.Token.Kind = TokenParam
			l.Token.AsString = l.Buffer[l.Offset+1 : l.Offset+i]
			return
		}
	case '`':
		l.Token.Kind = TokenIdent
		l.Token.AsString = l.nextQuotedContent("`", false, true, "identifier")
		return
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.nextNumber()
		return
	case 'B', 'b', 'R', 'r', '"', '\'':
		i, bytes, raw := 0, false, false
		for {
			switch {
			case !bytes && (l.peekIs(i, 'B') || l.peekIs(i, 'b')):
				i++
				bytes = true
			case !raw && (l.peekIs(i, 'R') || l.peekIs(i, 'r')):
				i++
				raw = true
			case l.peekIs(i, '"') || l.peekIs(i, '\''):
				l.nextN(i)
				switch {
				case bytes && raw:
					l.nextRawBytes()
				case bytes:
					l.nextBytes()
				case raw:
					l.nextRawString()
				default:
					l.nextString()
				}
				return
			default:
				break
			}
		}
	}
}

func (l *Lexer) nextNumber() {
	// https://cloud.google.com/spanner/docs/lexical#integer-literals
	// https://cloud.google.com/spanner/docs/lexical#floating-point-literals

	i := 0
	sign, base := "", 10

	switch {
	case l.peekIs(i, '+'):
		i++
	case l.peekIs(i, '-'):
		i++
		sign = "-"
	}

	if l.peekIs(i, '0') && (l.peekIs(i+1, 'x') || l.peekIs(i+1, 'X')) {
		i += 2
		base = 16
	}

	offset, int := i, true

peek:
	for l.peekOk(i) {
		c := l.peek(i)
		switch {
		case base == 10 && isDigit(c):
			i++
		case base == 16 && isHexDigit(c):
			i++
		case base == 10 && c == '.':
			i++
			int = false
		case base == 10 && (c == 'E' || c == 'e'):
			i++
			if l.peekIs(i, '+') || l.peekIs(i, '-') {
				i++
			}
			if !(l.peekOk(i) && isDigit(l.peek(i))) {
				l.errorf("invalid number literal")
			}
			int = false
		default:
			break peek
		}
	}

	var err error
	if int {
		l.Token.Kind = TokenInt
		l.Token.AsInt, err = strconv.ParseInt(sign+l.slice(offset, i), base, 64)
	} else {
		l.Token.Kind = TokenFloat
		l.Token.AsFloat, err = strconv.ParseFloat(sign+l.slice(offset, i), 64)
	}
	if err != nil {
		l.errorf("invalid number literal: %v", err)
	}

	l.nextN(i)
}

func (l *Lexer) nextRawBytes() {
	l.Token.Kind = TokenBytes
	l.Token.AsString = l.nextQuotedContent(l.peekDelimiter(), true, false, "raw bytes literal")
}

func (l *Lexer) nextBytes() {
	l.Token.Kind = TokenBytes
	l.Token.AsString = l.nextQuotedContent(l.peekDelimiter(), false, false, "bytes literal")
}

func (l *Lexer) nextRawString() {
	l.Token.Kind = TokenString
	l.Token.AsString = l.nextQuotedContent(l.peekDelimiter(), true, true, "raw string literal")
}

func (l *Lexer) nextString() {
	l.Token.Kind = TokenString
	l.Token.AsString = l.nextQuotedContent(l.peekDelimiter(), false, true, "raw string literal")
}

func (l *Lexer) peekDelimiter() string {
	i := 0
	c := l.peek(i)
	if c != '"' && c != '\'' {
		l.errorf("invalid delimiter: %v", c)
	}
	i++

	triple := true
	for l.peekOk(i) && i < 3 {
		if !l.peekIs(i, c) {
			triple = false
			break
		}
	}

	switch {
	case !triple && c == '"':
		return `"`
	case !triple && c == '\'':
		return `'`
	case triple && c == '"':
		return `"""`
	case triple && c == '\'':
		return `'''`
	}

	panic("unreachable")
}

func (l *Lexer) nextQuotedContent(q string, raw, unicode bool, name string) []byte {
	// https://cloud.google.com/spanner/docs/lexical#string-and-bytes-literals

	if len(q) == 3 {
		name = "triple-quoted " + name
	}

	i := len(q)
	var content []byte

	for l.peekOk(i) {
		if l.slice(i, i+len(q)) == q {
			l.nextN(i + len(q))
			return content
		}

		c := l.peek(i)
		if c == '\\' {
			i++
			if !l.peekOk(i) {
				l.errorf("invalid escape sequence: \\<EOF>")
			}

			c := l.peek(i)
			i++

			if raw {
				content = append(content, '\\', c)
				continue
			}

			switch c {
			case 'a':
				content = append(content, '\a')
			case 'b':
				content = append(content, '\b')
			case 'f':
				content = append(content, '\f')
			case 'n':
				content = append(content, '\n')
			case 'r':
				content = append(content, '\r')
			case 't':
				content = append(content, '\t')
			case 'v':
				content = append(content, '\v')
			case '\\', '?', '"', '\'', '`':
				content = append(content, c)
			case 'x', 'X':
				if !(l.peekOk(i+1) && isHexDigit(l.peek(i)) && isHexDigit(l.peek(i+1))) {
					l.errorf("invalid escape sequence: hex escape sequence must be follwed by 2 hex digits")
				}
				u, err := strconv.ParseUint(l.slice(i, i+2), 16, 8)
				if err != nil {
					l.errorf("invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			case 'u', 'U':
				if !unicode {
					l.errorf("invalid escape sequence: \\%c is not allowed in %s", c, name)
				}
				size := 4
				if c == 'U' {
					size = 8
				}
				for j := 0; j < size; j++ {
					if !(l.peekOk(i+j) && isHexDigit(l.peek(i+j))) {
						l.errorf("invalid escape sequence: \\%c must be followed by %d hex digits", c, size)
					}
				}
				u, err := strconv.ParseUint(l.slice(i, i+size), 16, 32)
				if err != nil {
					l.errorf("invalid escape sequence: %v", err)
				}
				if 0xD800 <= u && u <= 0xDFFF || 0x10FFFF < u {
					l.errorf("invalid escape sequence: invalid code point: %04x", u)
				}
				var buf [utf8.MaxRune]byte
				n := utf8.EncodeRune(buf[:], rune(u))
				content = append(content, buf[:n]...)
				i += size
			case '0', '1', '2', '3':
				if l.peekOk(i+2) && isOctalDigit(l.peek(i+1)) && isOctalDigit(l.peek(i+2)) {
					l.errorf("invalid escape sequence: octal escape sequence must be follwed by 3 octal digits")
				}
				u, err := strconv.ParseUint(l.slice(i-1, i+2), 8, 8)
				if err != nil {
					l.errorf("invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			default:
				l.errorf("invalid escape sequence: \\%c", c)
			}

			continue
		}

		content = append(content, c)
		i++
	}

	l.errorf("unclosed %s", name)
	panic("unreachable")
}

func (l *Lexer) skipSpaces() {
	for !l.eof() {
		r, size := utf8.DecodeRune(l.Buffer[l.Offset:])
		switch {
		case unicode.IsSpace(r):
			l.nextN(size)
		case r == '#' || r == '/' && l.peekIs(1, '/') || r == '-' && l.peekIs(1, '-'):
			l.skipComment("\n")
		case r == '/' && l.peekIs(1, '*'):
			l.skipComment("*/")
		default:
			return
		}
	}
}

func (l *Lexer) skipComment(end string) {
	for !l.eof() {
		if l.consume(end) {
			return
		}
		l.next()
	}
}

func (l *Lexer) peek(i int) byte {
	return l.Buffer[l.Offset+i]
}

func (l *Lexer) peekOk(i int) bool {
	return l.Offset+i < len(l.Buffer)
}

func (l *Lexer) peekIs(i int, c byte) bool {
	return l.Offset+i < len(l.Buffer) && l.Buffer[l.Offset+i] == c
}

func (l *Lexer) next() byte {
	c := l.Buffer[l.Offset]
	l.Offset++
	if c == '\n' {
		l.Line += 1
		l.Column = 0
	} else {
		l.Column += 1
	}
	return c
}

func (l *Lexer) nextN(n int) {
	for i := 0; i < n; i++ {
		l.next()
	}
}

func (l *Lexer) consume(s string) bool {
	n := 0
	for i, c := range []byte(s) {
		if !l.peekIs(i, c) && !l.peekIs(i, toUpper(c)) {
			return false
		}
		n++
	}
	for i := 0; i < n; i++ {
		l.next()
	}
	return true
}

func (l *Lexer) slice(start, end int) string {
	return string(l.Buffer[l.Offset+start : l.Offset+end])
}

func (l *Lexer) eof() bool {
	return l.Offset >= len(l.Buffer)
}

func (l *Lexer) errorf(msg string, param ...interface{}) {
	err := &Error{
		Message: fmt.Sprintf(msg, param...),
		Loc:     l.loc(),
	}
	panic(err)
}

func (l *Lexer) loc() *Location {
	return &Location{
		FilePath: l.FilePath,
		Offset:   l.Offset,
		Line:     l.Line + 1,
		Column:   l.Column + 1,
	}
}
