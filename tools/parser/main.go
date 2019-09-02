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
				/*
				select ARRAY_AGG(DISTINCT A)[OFFSET(0)] from
				  (select 1 as A, 2 as  B union all select 1 as A, 5 as B) as X
				join
				  (select 1 as A, 3 as C union all select 1 as A, 4 as C) as Y
				using (A)
				limit 1 offset 0
				*/
				SELECT DATE_ADD(DATE "2008-12-25", INTERVAL 5 DAY) as five_days_later
			`),
			/*
				heredoc.Doc(`
					-- select ((select 1) LIMIT 1) + 1  IN UNNEST(ARRAY(select 2 union all select 3)), 1 BETWEEN 0 AND 10
					-- select case 1 when 2 then 3 else 4 end
					select C.D A FROM B C
				`),
			*/
		),
		// `((select (((select 1) union distinct select 1 union distinct select 1) + 1)))`),
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
