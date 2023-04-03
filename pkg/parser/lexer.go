package parser

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/cloudspannerecosystem/memefish/pkg/char"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
)

type Lexer struct {
	*token.File
	Token token.Token

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

	lastTokenKind token.TokenKind
	dotIdent      bool
}

func (l *Lexer) Clone() *Lexer {
	lex := *l
	return &lex
}

func isNextDotIdent(t token.TokenKind) bool {
	switch t {
	case token.TokenIdent, token.TokenParam, ")", "]":
		return true
	}
	return false
}

// NextToken reads a next token from source, then updates its Token field.
func (l *Lexer) NextToken() (err error) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*Error)
			if ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	l.nextToken()
	return
}

func (l *Lexer) nextToken() {
	l.lastTokenKind = l.Token.Kind
	l.Token = token.Token{}

	// Skips spaces and comments.
	var space string
	for {
		i := l.pos
		l.skipSpaces()
		space = l.Buffer[i:l.pos]

		i = l.pos
		l.skipComment()
		if l.pos == i {
			break
		}
		l.Token.Comments = append(l.Token.Comments, token.TokenComment{
			Space: space,
			Raw:   l.Buffer[i:l.pos],
			Pos:   token.Pos(i),
			End:   token.Pos(l.pos),
		})
	}

	l.Token.Space = space

	// Reads the next token.
	l.Token.Pos = token.Pos(l.pos)
	i := l.pos
	if l.dotIdent {
		l.consumeFieldToken()
		l.dotIdent = false
	} else {
		l.consumeToken()
	}
	l.Token.Raw = l.Buffer[i:l.pos]
	l.Token.End = token.Pos(l.pos)
}

func (l *Lexer) consumeToken() {
	if l.eof() {
		l.Token.Kind = token.TokenEOF
		return
	}

	switch l.peek(0) {
	case '(', ')', '{', '}', ';', ',', '[', ']', '~', '*', '/', '&', '^', '=', '+', '-':
		l.Token.Kind = token.TokenKind([]byte{l.skip()})
		return
	case '.':
		nextDotIdent := isNextDotIdent(l.lastTokenKind)
		if !nextDotIdent && l.peekOk(1) && char.IsDigit(l.peek(1)) {
			l.consumeNumber()
		} else {
			l.skip()
			l.Token.Kind = "."
			l.dotIdent = nextDotIdent
		}
		return
	case '<':
		switch {
		case l.peekIs(1, '<'):
			l.skipN(2)
			l.Token.Kind = "<<"
		case l.peekIs(1, '='):
			l.skipN(2)
			l.Token.Kind = "<="
		case l.peekIs(1, '>'):
			l.skipN(2)
			l.Token.Kind = "<>"
		default:
			l.skip()
			l.Token.Kind = "<"
		}
		return
	case '>':
		switch {
		case l.peekIs(1, '>'):
			l.skipN(2)
			l.Token.Kind = ">>"
		case l.peekIs(1, '='):
			l.skipN(2)
			l.Token.Kind = ">="
		default:
			l.skip()
			l.Token.Kind = ">"
		}
		return
	case '|':
		switch {
		case l.peekIs(1, '|'):
			l.skipN(2)
			l.Token.Kind = "||"
		default:
			l.skip()
			l.Token.Kind = "|"
		}
		return
	case '!':
		if l.peekIs(1, '=') {
			l.skipN(2)
			l.Token.Kind = "!="
			return
		}
	case '@':
		if l.peekOk(1) && char.IsIdentStart(l.peek(1)) {
			i := 1
			for l.peekOk(i) && char.IsIdentPart(l.peek(i)) {
				i++
			}
			l.Token.Kind = token.TokenParam
			l.Token.AsString = l.Buffer[l.pos+1 : l.pos+i]
			l.skipN(i)
			return
		}
		l.skip()
		l.Token.Kind = "@"
		return
	case '`':
		l.Token.Kind = token.TokenIdent
		l.Token.AsString = l.consumeQuotedContent("`", false, true, "identifier")
		return
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.consumeNumber()
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
				l.skipN(i)
				switch {
				case bytes && raw:
					l.consumeRawBytes()
				case bytes:
					l.consumeBytes()
				case raw:
					l.consumeRawString()
				default:
					l.consumeString()
				}
				return
			default:
				break loop
			}
		}
	}

	if char.IsIdentStart(l.peek(0)) {
		i := 0
		for l.peekOk(i) && char.IsIdentPart(l.peek(i)) {
			i++
		}
		s := l.slice(0, i)
		l.skipN(i)
		k := token.TokenKind(char.ToUpper(s))
		if _, ok := token.KeywordsMap[k]; ok {
			l.Token.Kind = k
		} else {
			l.Token.Kind = token.TokenIdent
			l.Token.AsString = s
		}
		return
	}

	panic(l.errorf("illegal input character: %q", l.peek(0)))
}

func (l *Lexer) consumeFieldToken() {
	if l.peekOk(0) && char.IsIdentPart(l.peek(0)) {
		i := 0
		for l.peekOk(i) && char.IsIdentPart(l.peek(i)) {
			i++
		}
		l.Token.Kind = token.TokenIdent
		l.Token.AsString = l.Buffer[l.pos : l.pos+i]
		l.skipN(i)
		return
	}

	l.consumeToken()
}

func (l *Lexer) consumeNumber() {
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
		case base == 10 && char.IsDigit(c):
			i++
			continue
		case base == 16 && char.IsHexDigit(c):
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
			if !(l.peekOk(i) && char.IsDigit(l.peek(i))) {
				i = rollback
				break
			}
			exp = true
			int = false
			continue
		}
		break
	}

	l.skipN(i)
	if int {
		l.Token.Kind = token.TokenInt
		l.Token.Base = base
	} else {
		l.Token.Kind = token.TokenFloat
	}

	if l.peekOk(0) && char.IsIdentPart(l.peek(0)) {
		l.panicf("number literal cannot follow identifier without any spaces")
	}
}

func (l *Lexer) consumeRawBytes() {
	l.Token.Kind = token.TokenBytes
	l.Token.AsString = l.consumeQuotedContent(l.peekDelimiter(), true, false, "raw bytes literal")
}

func (l *Lexer) consumeBytes() {
	l.Token.Kind = token.TokenBytes
	l.Token.AsString = l.consumeQuotedContent(l.peekDelimiter(), false, false, "bytes literal")
}

func (l *Lexer) consumeRawString() {
	l.Token.Kind = token.TokenString
	l.Token.AsString = l.consumeQuotedContent(l.peekDelimiter(), true, true, "raw string literal")
}

func (l *Lexer) consumeString() {
	l.Token.Kind = token.TokenString
	l.Token.AsString = l.consumeQuotedContent(l.peekDelimiter(), false, true, "string literal")
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

func (l *Lexer) consumeQuotedContent(q string, raw, unicode bool, name string) string {
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
			l.skipN(i + len(q))
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
				if !(l.peekOk(i+1) && char.IsHexDigit(l.peek(i)) && char.IsHexDigit(l.peek(i+1))) {
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
					if !(l.peekOk(i+j) && char.IsHexDigit(l.peek(i+j))) {
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
				if !(l.peekOk(i+1) && char.IsOctalDigit(l.peek(i)) && char.IsOctalDigit(l.peek(i+1))) {
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
			l.skipN(size)
		default:
			return
		}
	}
}

func (l *Lexer) skipComment() {
	r, _ := utf8.DecodeRuneInString(l.Buffer[l.pos:])
	switch {
	case r == '#' || r == '/' && l.peekIs(1, '/') || r == '-' && l.peekIs(1, '-'):
		l.skipCommentUntil("\n", false)
	case r == '/' && l.peekIs(1, '*'):
		l.skipCommentUntil("*/", true)
	default:
		return
	}
}

func (l *Lexer) skipCommentUntil(end string, mustEnd bool) {
	for !l.eof() {
		if l.slice(0, len(end)) == end {
			l.skipN(len(end))
			return
		}
		l.skip()
	}
	if mustEnd {
		// TODO: improve error position
		l.panicf("unclosed comment")
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

func (l *Lexer) skip() byte {
	c := l.Buffer[l.pos]
	l.pos++
	return c
}

func (l *Lexer) skipN(n int) {
	l.pos += n
}

func (l *Lexer) slice(start, end int) string {
	if len(l.Buffer) < l.pos+end {
		end = len(l.Buffer) - l.pos
	}
	return string(l.Buffer[l.pos+start : l.pos+end])
}

func (l *Lexer) eof() bool {
	return l.pos >= len(l.Buffer)
}

func (l *Lexer) errorf(msg string, param ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, param...),
		Position: l.Position(token.Pos(l.pos), token.Pos(l.pos)),
	}
}

func (l *Lexer) panicf(msg string, param ...interface{}) {
	panic(l.errorf(msg, param...))
}
