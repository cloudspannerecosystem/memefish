package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	. "github.com/MakeNowJust/memefish/pkg/token"
)

type Lexer struct {
	*File
	Token Token

	pos int

	// ZetaSQL's lexer comment:
	//
	// > After "." we allow more things, including all keywords and all
	// > integers, to be returned as identifiers. This state is initiated when we
	// > recognize an identifier followed by a ".". It is also initiated after a
	// > closing parenthesis, square bracket, or "?" (positional parameter) followed
	// > by a ".", to handle cases like foo[3].array. See the "." rule and the
	// > <DOT_IDENTIFIER>{generalized_identifier}.
	// > https://github.com/google/zetasql/blob/e269a26107e9b6c5a43a72d3a323faa19dd4599b/zetasql/parser/flex_tokenizer.l#L28-L33
	//
	// For implementing this, it should keep the lastTokenKind on lexing and have dotIdent flag.
	// But they are internal state, so let them private.

	lastTokenKind TokenKind
	dotIdent      bool
}

func (l *Lexer) Clone() *Lexer {
	lex := *l
	return &lex
}

func isDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

func isHexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

func isOctalDigit(c byte) bool {
	return '0' <= c && c <= '7'
}

func isIdentStart(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isIdentPart(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_'
}

func isNextDotIdent(t TokenKind) bool {
	switch t {
	case TokenIdent, TokenParam, ")", "]":
		return true
	}
	return false
}

func (l *Lexer) NextToken() {
	l.lastTokenKind = l.Token.Kind
	l.Token = Token{}

	// Skips spaces.
	i := l.pos
	l.skipSpaces()
	l.Token.Space = l.Buffer[i:l.pos]

	// Reads the next token.
	l.Token.Pos = Pos(l.pos)
	i = l.pos
	if l.dotIdent {
		l.nextFieldToken()
		l.dotIdent = false
	} else {
		l.nextToken()
	}
	l.Token.Raw = l.Buffer[i:l.pos]
	l.Token.End = Pos(l.pos)
}

func (l *Lexer) nextToken() {
	if l.eof() {
		l.Token.Kind = TokenEOF
		return
	}

	switch l.peek(0) {
	case '(', ')', '{', '}', ';', ',', '[', ']', '~', '*', '/', '&', '^', '|', '=':
		l.Token.Kind = TokenKind([]byte{l.next()})
		return
	case '+', '-':
		l.Token.Kind = TokenKind([]byte{l.next()})
		return
	case '.':
		nextDotIdent := isNextDotIdent(l.lastTokenKind)
		if !nextDotIdent && l.peekOk(1) && isDigit(l.peek(1)) {
			l.nextNumber()
		} else {
			l.next()
			l.Token.Kind = "."
			l.dotIdent = nextDotIdent
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
		if l.peekOk(1) && isIdentStart(l.peek(1)) {
			i := 1
			for l.peekOk(i) && isIdentPart(l.peek(i)) {
				i++
			}
			l.Token.Kind = TokenParam
			l.Token.AsString = l.Buffer[l.pos+1 : l.pos+i]
			l.nextN(i)
			return
		}
		l.next()
		l.Token.Kind = "@"
		return
	case '`':
		l.Token.Kind = TokenIdent
		l.Token.AsString = l.nextQuotedContent("`", false, true, "identifier")
		return
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.nextNumber()
		return
	case 'B', 'b', 'R', 'r', '"', '\'':
		bytes, raw := false, false
	loop:
		for i := 0; i < 3 && l.peekOk(i); i++ {
			switch {
			case !bytes && (l.peekIs(i, 'B') || l.peekIs(i, 'b')):
				bytes = true
			case !raw && (l.peekIs(i, 'R') || l.peekIs(i, 'r')):
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
				break loop
			}
		}
	}

	if isIdentStart(l.peek(0)) {
		i := 0
		for l.peekOk(i) && isIdentPart(l.peek(i)) {
			i++
		}
		s := l.slice(0, i)
		l.nextN(i)
		k := TokenKind(strings.ToUpper(s))
		if _, ok := KeywordsMap[k]; ok {
			l.Token.Kind = k
		} else {
			l.Token.Kind = TokenIdent
			l.Token.AsString = s
		}
		return
	}

	panic(l.errorf("illegal input character: %q", l.peek(0)))
}

func (l *Lexer) nextFieldToken() {
	if l.peekOk(0) && isIdentPart(l.peek(0)) {
		i := 0
		for l.peekOk(i) && isIdentPart(l.peek(i)) {
			i++
		}
		l.Token.Kind = TokenIdent
		l.Token.AsString = l.Buffer[l.pos : l.pos+i]
		l.nextN(i)
		return
	}

	l.nextToken()
}

func (l *Lexer) nextNumber() {
	// https://cloud.google.com/spanner/docs/lexical#integer-literals
	// https://cloud.google.com/spanner/docs/lexical#floating-point-literals

	i := 0
	base := 10

	if l.peekIs(i, '0') && (l.peekIs(i+1, 'x') || l.peekIs(i+1, 'X')) {
		i += 2
		base = 16
	}

	int, exp := true, false

	for l.peekOk(i) {
		c := l.peek(i)
		switch {
		case base == 10 && isDigit(c):
			i++
			continue
		case base == 16 && isHexDigit(c):
			i++
			continue
		case !exp && int && base == 10 && c == '.':
			i++
			int = false
			continue
		case !exp && base == 10 && (c == 'E' || c == 'e'):
			rollback := i
			i++
			if l.peekIs(i, '+') || l.peekIs(i, '-') {
				i++
			}
			if !(l.peekOk(i) && isDigit(l.peek(i))) {
				i = rollback
				break
			}
			exp = true
			int = false
			continue
		}
		break
	}

	l.nextN(i)
	if int {
		l.Token.Kind = TokenInt
		l.Token.Base = base
	} else {
		l.Token.Kind = TokenFloat
	}

	if l.peekOk(0) && isIdentPart(l.peek(0)) {
		l.panicf("number literal cannot follow identifier without any spaces")
	}
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
	l.Token.AsString = l.nextQuotedContent(l.peekDelimiter(), false, true, "string literal")
}

func (l *Lexer) peekDelimiter() string {
	i := 0
	c := l.peek(i)
	if c != '"' && c != '\'' {
		l.panicf("invalid delimiter: %v", c)
	}
	i++

	triple := true
	for i < 3 {
		if !l.peekIs(i, c) {
			triple = false
			break
		}
		i++
	}

	switch {
	case triple && c == '"':
		return `"""`
	case triple && c == '\'':
		return `'''`
	default:
		return string([]byte{c})
	}
}

func (l *Lexer) nextQuotedContent(q string, raw, unicode bool, name string) string {
	// https://cloud.google.com/spanner/docs/lexical#string-and-bytes-literals

	if len(q) == 3 {
		name = "triple-quoted " + name
	}

	i := len(q)
	var content []byte

	for l.peekOk(i) {
		if l.slice(i, i+len(q)) == q {
			if len(content) == 0 && name == "identifier" {
				l.panicf("invalid empty identifier")
			}
			l.nextN(i + len(q))
			return string(content)
		}

		c := l.peek(i)
		if c == '\\' {
			i++
			if !l.peekOk(i) {
				l.panicf("invalid escape sequence: \\<eof>")
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
					l.panicf("invalid escape sequence: hex escape sequence must be follwed by 2 hex digits")
				}
				u, err := strconv.ParseUint(l.slice(i, i+2), 16, 8)
				if err != nil {
					l.panicf("invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			case 'u', 'U':
				if !unicode {
					l.panicf("invalid escape sequence: \\%c is not allowed in %s", c, name)
				}
				size := 4
				if c == 'U' {
					size = 8
				}
				for j := 0; j < size; j++ {
					if !(l.peekOk(i+j) && isHexDigit(l.peek(i+j))) {
						l.panicf("invalid escape sequence: \\%c must be followed by %d hex digits", c, size)
					}
				}
				u, err := strconv.ParseUint(l.slice(i, i+size), 16, 32)
				if err != nil {
					l.panicf("invalid escape sequence: %v", err)
				}
				if 0xD800 <= u && u <= 0xDFFF || 0x10FFFF < u {
					l.panicf("invalid escape sequence: invalid code point: U+%04X", u)
				}
				var buf [utf8.MaxRune]byte
				n := utf8.EncodeRune(buf[:], rune(u))
				content = append(content, buf[:n]...)
				i += size
			case '0', '1', '2', '3':
				if !(l.peekOk(i+1) && isOctalDigit(l.peek(i)) && isOctalDigit(l.peek(i+1))) {
					l.panicf("invalid escape sequence: octal escape sequence must be follwed by 3 octal digits")
				}
				u, err := strconv.ParseUint(l.slice(i-1, i+2), 8, 8)
				if err != nil {
					l.panicf("invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			default:
				l.panicf("invalid escape sequence: \\%c", c)
			}

			continue
		}

		if c == '\n' && len(q) != 3 {
			l.panicf("unclosed %s: newline appears in non triple-quoted", name)
		}

		content = append(content, c)
		i++
	}

	panic(l.errorf("unclosed %s", name))
}

func (l *Lexer) skipSpaces() {
	for !l.eof() {
		r, size := utf8.DecodeRuneInString(l.Buffer[l.pos:])
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
		if l.slice(0, len(end)) == end {
			l.nextN(len(end))
			return
		}
		l.next()
	}
}

func (l *Lexer) peek(i int) byte {
	return l.Buffer[l.pos+i]
}

func (l *Lexer) peekOk(i int) bool {
	return l.pos+i < len(l.Buffer)
}

func (l *Lexer) peekIs(i int, c byte) bool {
	return l.pos+i < len(l.Buffer) && l.Buffer[l.pos+i] == c
}

func (l *Lexer) next() byte {
	c := l.Buffer[l.pos]
	l.pos++
	return c
}

func (l *Lexer) nextN(n int) {
	l.pos += n
}

func (l *Lexer) slice(start, end int) string {
	return string(l.Buffer[l.pos+start : l.pos+end])
}

func (l *Lexer) eof() bool {
	return l.pos >= len(l.Buffer)
}

func (l *Lexer) errorf(msg string, param ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, param...),
		Position: l.Position(Pos(l.pos), Pos(l.pos)),
	}
}

func (l *Lexer) panicf(msg string, param ...interface{}) {
	panic(l.errorf(msg, param...))
}
