package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/k0kubun/pp"
)

var mode = flag.String("mode", "statement", "parsing mode")

func init() {
	flag.Parse()
}

func main() {
	if flag.NArg() < 1 {
		log.Fatal("usage: ./parse [-mode statement|query|expr|ddl|dml] [SQL query]")
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
	switch *mode {
	case "statement":
		node, err = p.ParseStatement()
	case "query":
		node, err = p.ParseQuery()
	case "expr":
		node, err = p.ParseExpr()
	case "ddl":
		node, err = p.ParseDDL()
	case "dml":
		node, err = p.ParseDML()
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("finish parsing successfully")

	fmt.Println("--- AST")
	_, _ = pp.Println(node)
	fmt.Println()
	fmt.Println("--- SQL")
	fmt.Println(node.SQL())
}
