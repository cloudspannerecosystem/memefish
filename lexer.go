package memefish

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"github.com/cloudspannerecosystem/memefish/char"
	"github.com/cloudspannerecosystem/memefish/token"
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

	l.nextToken(false)
	return
}

func (l *Lexer) nextToken(noPanic bool) {
	l.lastTokenKind = l.Token.Kind
	l.Token = token.Token{}

	// Skips spaces and comments.
	var space string
	for {
		i := l.pos
		l.skipSpaces()
		space = l.Buffer[i:l.pos]

		i = l.pos
		hasError := l.skipComment(noPanic)

		if l.pos == i {
			break
		}
		l.Token.Comments = append(l.Token.Comments, token.TokenComment{
			Space: space,
			Raw:   l.Buffer[i:l.pos],
			Pos:   token.Pos(i),
			End:   token.Pos(l.pos),
		})

		if hasError {
			l.Token.Pos = token.Pos(l.pos)
			l.Token.End = token.Pos(l.pos)
			l.Token.Kind = token.TokenBad
			return
		}
	}

	l.Token.Space = space

	// Reads the next token.
	l.Token.Pos = token.Pos(l.pos)
	i := l.pos
	if l.dotIdent {
		l.consumeFieldToken(noPanic)
		l.dotIdent = false
	} else {
		l.consumeToken(noPanic)
	}
	l.Token.Raw = l.Buffer[i:l.pos]
	l.Token.End = token.Pos(l.pos)
}

func (l *Lexer) consumeToken(noPanic bool) {
	if l.eof() {
		l.Token.Kind = token.TokenEOF
		return
	}

	switch l.peek(0) {
	case '(', ')', '{', '}', ';', ',', '[', ']', '~', '*', '/', '&', '^', '%', ':',
		// Belows are not yet used in Spanner.
		'?', '\\', '$':
		l.Token.Kind = token.TokenKind([]byte{l.skip()})
		return
	case '.':
		nextDotIdent := isNextDotIdent(l.lastTokenKind)
		if !nextDotIdent && l.peekOk(1) && char.IsDigit(l.peek(1)) {
			l.consumeNumber(noPanic)
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
	case '+':
		switch {
		// KW_ADD_ASSIGN in ZetaSQL
		case l.peekIs(1, '='):
			l.skipN(2)
			l.Token.Kind = "+="
		default:
			l.skip()
			l.Token.Kind = "+"
		}
		return
	case '-':
		switch {
		// KW_SUB_ASSIGN in ZetaSQL
		case l.peekIs(1, '='):
			l.skipN(2)
			l.Token.Kind = "-="
		// KW_LAMBDA_ARROW in ZetaSQL
		case l.peekIs(1, '>'):
			l.skipN(2)
			l.Token.Kind = "->"
		default:
			l.skip()
			l.Token.Kind = "-"
		}
		return
	case '=':
		switch {
		case l.peekIs(1, '>'):
			l.skipN(2)
			l.Token.Kind = "=>"
		default:
			l.skip()
			l.Token.Kind = "="
		}
		return
	case '|':
		switch {
		case l.peekIs(1, '>'):
			l.skipN(2)
			l.Token.Kind = "|>"
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
		l.skip()
		l.Token.Kind = "!"
		return
	case '@':
		// KW_DOUBLE_AT is not yet used in Cloud Spanner, but used in BigQuery.
		if l.peekIs(1, '@') {
			l.skipN(2)
			l.Token.Kind = "@@"
			return
		}
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

		var hasError bool
		l.Token.AsString, hasError = l.consumeQuotedContent("`", false, true, "identifier", noPanic)
		if hasError {
			l.Token.Kind = token.TokenBad
		}
		return
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.consumeNumber(noPanic)
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
					l.consumeRawBytes(noPanic)
				case bytes:
					l.consumeBytes(noPanic)
				case raw:
					l.consumeRawString(noPanic)
				default:
					l.consumeString(noPanic)
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

	if noPanic {
		l.skip()
		l.Token.Kind = token.TokenBad
		return
	}

	panic(l.errorf("illegal input character: %q", l.peek(0)))
}

func (l *Lexer) consumeFieldToken(noPanic bool) {
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

	l.consumeToken(noPanic)
}

func (l *Lexer) consumeNumber(noPanic bool) {
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
			if !l.peekOk(i) || !char.IsDigit(l.peek(i)) {
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
		if noPanic {
			l.Token.Kind = token.TokenBad
			return
		}

		l.panicf("number literal cannot follow identifier without any spaces")
	}
}

func (l *Lexer) consumeRawBytes(noPanic bool) {
	l.Token.Kind = token.TokenBytes

	var hasError bool
	l.Token.AsString, hasError = l.consumeQuotedContent(l.peekDelimiter(), true, false, "raw bytes literal", noPanic)
	if hasError {
		l.Token.Kind = token.TokenBad
	}
}

func (l *Lexer) consumeBytes(noPanic bool) {
	l.Token.Kind = token.TokenBytes

	var hasError bool
	l.Token.AsString, hasError = l.consumeQuotedContent(l.peekDelimiter(), false, false, "bytes literal", noPanic)
	if hasError {
		l.Token.Kind = token.TokenBad
	}
}

func (l *Lexer) consumeRawString(noPanic bool) {
	l.Token.Kind = token.TokenString

	var hasError bool
	l.Token.AsString, hasError = l.consumeQuotedContent(l.peekDelimiter(), true, true, "raw string literal", noPanic)
	if hasError {
		l.Token.Kind = token.TokenBad
	}
}

func (l *Lexer) consumeString(noPanic bool) {
	l.Token.Kind = token.TokenString

	var hasError bool
	l.Token.AsString, hasError = l.consumeQuotedContent(l.peekDelimiter(), false, true, "string literal", noPanic)
	if hasError {
		l.Token.Kind = token.TokenBad
	}
}

func (l *Lexer) peekDelimiter() string {
	i := 0
	c := l.peek(i)
	if c != '"' && c != '\'' {
		// This error is unreachable
		panic(fmt.Sprintf("BUG: invalid delimiter: %v", c))
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

func (l *Lexer) consumeQuotedContent(q string, raw, unicode bool, name string, noPanic bool) (string, bool) {
	// https://cloud.google.com/spanner/docs/lexical#string-and-bytes-literals

	if len(q) == 3 {
		name = "triple-quoted " + name
	}

	i := len(q)
	var content []byte
	hasError := false

	for l.peekOk(i) {
		if l.slice(i, i+len(q)) == q {
			if len(content) == 0 && name == "identifier" {
				if noPanic {
					hasError = true
				} else {
					l.panicfAtPosition(token.Pos(l.pos), token.Pos(l.pos+i+len(q)), "invalid empty identifier")
				}
			}
			l.skipN(i + len(q))

			if hasError {
				return "", true
			}
			return string(content), false
		}

		c := l.peek(i)
		if c == '\\' {
			i++
			if !l.peekOk(i) {
				if noPanic {
					hasError = true
					continue
				}
				l.panicfAtPosition(token.Pos(l.pos+i-1), token.Pos(l.pos+i), "invalid escape sequence: \\<eof>")
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
				for j := 0; j < 2; j++ {
					if !l.peekOk(i+j) || !char.IsHexDigit(l.peek(i+j)) {
						if noPanic {
							hasError = true
							continue
						}
						l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+j+1), "invalid escape sequence: hex escape sequence must be follwed by 2 hex digits")
					}
				}
				u, err := strconv.ParseUint(l.slice(i, i+2), 16, 8)
				if err != nil {
					if noPanic {
						hasError = true
						continue
					}
					l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+2), "invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			case 'u', 'U':
				if !unicode {
					if noPanic {
						hasError = true
						continue
					}
					l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i), "invalid escape sequence: \\%c is not allowed in %s", c, name)
				}
				size := 4
				if c == 'U' {
					size = 8
				}
				for j := 0; j < size; j++ {
					if !l.peekOk(i+j) || !char.IsHexDigit(l.peek(i+j)) {
						if noPanic {
							hasError = true
							continue
						}
						l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+j+1), "invalid escape sequence: \\%c must be followed by %d hex digits", c, size)
					}
				}
				u, err := strconv.ParseUint(l.slice(i, i+size), 16, 32)
				if err != nil {
					if noPanic {
						hasError = true
						continue
					}
					l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+size), "invalid escape sequence: %v", err)
				}
				if 0xD800 <= u && u <= 0xDFFF || 0x10FFFF < u {
					if noPanic {
						hasError = true
						continue
					}
					l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+size), "invalid escape sequence: invalid code point: U+%04X", u)
				}
				var buf [utf8.MaxRune]byte
				n := utf8.EncodeRune(buf[:], rune(u))
				content = append(content, buf[:n]...)
				i += size
			case '0', '1', '2', '3':
				for j := 0; j < 2; j++ {
					if !l.peekOk(i+j) || !char.IsOctalDigit(l.peek(i+j)) {
						if noPanic {
							hasError = true
							continue
						}
						l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+j+1), "invalid escape sequence: octal escape sequence must be follwed by 3 octal digits")
					}
				}
				u, err := strconv.ParseUint(l.slice(i-1, i+2), 8, 8)
				if err != nil {
					if noPanic {
						hasError = true
						continue
					}
					l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i+2), "invalid escape sequence: %v", err)
				}
				content = append(content, byte(u))
				i += 2
			default:
				if noPanic {
					hasError = true
					continue
				}
				l.panicfAtPosition(token.Pos(l.pos+i-2), token.Pos(l.pos+i), "invalid escape sequence: \\%c", c)
			}

			continue
		}

		if c == '\n' && len(q) != 3 {
			if noPanic {
				hasError = true
				i++
				continue
			}
			l.panicfAtPosition(token.Pos(l.pos), token.Pos(l.pos+i), "unclosed %s: newline appears in non triple-quoted", name)
		}

		content = append(content, c)
		i++
	}

	if noPanic {
		l.skipN(i)
		return "", true
	}

	panic(l.errorfAtPosition(token.Pos(l.pos), token.Pos(l.pos+i), "unclosed %s", name))
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

func (l *Lexer) skipComment(noPanic bool) bool {
	r, _ := utf8.DecodeRuneInString(l.Buffer[l.pos:])
	switch {
	case r == '#' || r == '/' && l.peekIs(1, '/') || r == '-' && l.peekIs(1, '-'):
		return l.skipCommentUntil("\n", false, noPanic)
	case r == '/' && l.peekIs(1, '*'):
		return l.skipCommentUntil("*/", true, noPanic)
	default:
		return false
	}
}

func (l *Lexer) skipCommentUntil(end string, mustEnd bool, noPanic bool) bool {
	pos := token.Pos(l.pos)
	for !l.eof() {
		if l.slice(0, len(end)) == end {
			l.skipN(len(end))
			return false
		}
		l.skip()
	}
	if mustEnd {
		if noPanic {
			return true
		}
		l.panicfAtPosition(pos, token.Pos(l.pos), "unclosed comment")
	}

	return false
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

func (l *Lexer) errorfAtPosition(pos, end token.Pos, msg string, param ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, param...),
		Position: l.Position(pos, end),
	}
}

func (l *Lexer) panicf(msg string, param ...interface{}) {
	panic(l.errorf(msg, param...))
}

func (l *Lexer) panicfAtPosition(pos, end token.Pos, msg string, param ...interface{}) {
	panic(l.errorfAtPosition(pos, end, msg, param...))
}
