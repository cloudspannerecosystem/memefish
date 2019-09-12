package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/k0kubun/pp"
	"github.com/MakeNowJust/memefish/pkg/parser"
)

var t = flag.String("type", "query", "statement type")

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
			File: &parser.File{FilePath: "", Buffer: query},
		},
	}

	log.Printf("start parsing")
	var node parser.Node
	var err error
	switch *t {
	case "query":
		node, err = p.ParseQuery()
	case "expr":
		node, err = p.ParseExpr()
	case "ddl":
		node, err = p.ParseDDL()
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("finish parsing successfully")

	fmt.Println("--- AST")
	pp.Println(node)
	fmt.Println()
	fmt.Println("--- SQL")
	fmt.Println(node.SQL())
}
