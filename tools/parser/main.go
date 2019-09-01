package main

import (
	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func main() {
	l := &parser.Lexer{
		File: parser.NewFile("[input]", `CAST(STRUCT< >("foo", 42) AS STRUCT<STRING, INT64>)`),
	}
	p := &parser.Parser{
		Lexer: l,
	}

	e := p.ParseExpr()
	_, _ = pp.Println(e)
	_, _ = pp.Println(l.File.Position(e.Pos(), e.End()))
	_, _ = pp.Println(e.SQL())
}
