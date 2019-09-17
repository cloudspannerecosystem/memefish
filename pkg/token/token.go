package token

import (
	"github.com/MakeNowJust/memefish/pkg/char"
)

type Token struct {
	Kind     TokenKind
	Space    string // TODO: better comment support
	Raw      string
	AsString string // available for TokenIdent, TokenString and TokenBytes
	Base     int    // 10 or 16 on TokenInt
	Pos, End Pos
}

func (t *Token) IsIdent(s string) bool {
	return t.Kind == TokenIdent && char.EqualFold(t.AsString, s)
}

func (t *Token) IsKeywordLike(s string) bool {
	return t.Kind == TokenIdent && char.EqualFold(t.Raw, s)
}

func (t *Token) Clone() *Token {
	tok := *t
	return &tok
}

type TokenKind string

const (
	TokenEOF    TokenKind = "<eof>"
	TokenIdent  TokenKind = "<ident>"
	TokenParam  TokenKind = "<param>"
	TokenInt    TokenKind = "<int>"
	TokenFloat  TokenKind = "<float>"
	TokenString TokenKind = "<string>"
	TokenBytes  TokenKind = "<bytes>"
)
