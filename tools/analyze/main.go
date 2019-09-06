package main

import (
	"flag"
	"log"
	"os"

	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/analyzer"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/olekukonko/tablewriter"
)

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		log.Fatal("usage: ./analyze [SQL query]")
	}

	query := flag.Arg(0)

	log.Printf("query: %q", query)

	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: parser.NewFile("[query]", query),
		},
	}

	log.Printf("start parsing")
	stmt, err := p.ParseQuery()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("finish parsing successfully")

	log.Printf("start analyzing")
	a := &analyzer.Analyzer{
		File: p.File,
	}
	a.AnalyzeQueryStatement(stmt)
	log.Printf("finish analyzing")

	list := a.NameLists[stmt.Query]
	if list == nil {
		log.Fatal("missing name list")
	}

	pp.Println(list)

	table := tablewriter.NewWriter(os.Stdout)
	var header []string
	for _, r := range list {
		header = append(header, r.GetName())
	}
	table.SetHeader(header)

	var types []string
	for _, r := range list {
		types = append(types, analyzer.TypeString(r.Type))
	}
	table.Append(types)
	table.Render()
}
