package main

import (
	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func main() {
	l := &parser.Lexer{
		File: parser.NewFile("[input]", "+12"),
	}

	for {
		l.NextToken()
		if l.Token.Kind == parser.TokenEOF {
			return
		}
		_, _ = pp.Println(l.Token)
	}
}
