package main

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func main() {
	l := &parser.Lexer{
		File: parser.NewFile(
			"[input]",
			/*
							heredoc.Doc(`
							select * from
				  		  (select 1 as A, 2 as  B union all select 1 as A, 5 as B) as X
							JOIN
							  (select 1 as A, 3 as C union all select 1 as A, 4 as C) as Y
							USING (A)
							LIMIT 1
							`),
			*/
			heredoc.Doc(`
				select ((select 1) LIMIT 1) + 1
			`),
		),
		// `((select (((select 1) union distinct select 1 union distinct select 1) + 1)))`),
	}
	p := &parser.Parser{
		Lexer: l,
	}

	e := p.ParseQuery()
	_, _ = pp.Println(e)
	_, _ = pp.Println(l.File.Position(e.Pos(), e.End()))
	_, _ = pp.Println(e.SQL())
}
