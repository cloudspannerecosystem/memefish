package main

import (
	"log"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func main() {
	l := &parser.Lexer{
		File: parser.NewFile(
			"[input]",
			heredoc.Doc(`
				select ARRAY_AGG(DISTINCT A)[OFFSET(0)] from
				  (select 1 as A, 2 as  B union all select 1 as A, 5 as B) as X
				join
				  (select 1 as A, 3 as C union all select 1 as A, 4 as C) as Y
				using (A)
				limit 1 offset 0
			`),
		),
	}
	p := &parser.Parser{
		Lexer: l,
	}

	e, err := p.ParseQuery()
	if err != nil {
		log.Fatal(err)
	}
	_, _ = pp.Println(e)
	_, _ = pp.Println(l.File.Position(e.Pos(), e.End()))
	_, _ = pp.Println(e.SQL())
}
