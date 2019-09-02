package main

import (
	"flag"
	"log"
	"fmt"

	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/k0kubun/pp"
)

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		log.Fatal("usage: ./parse [SQL query]")
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

	fmt.Println("--- AST")
	pp.Println(stmt)
	fmt.Println()
	fmt.Println("--- SQL")
	fmt.Println(stmt.SQL())
}
