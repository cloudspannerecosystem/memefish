package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/k0kubun/pp"
)

var mode = flag.String("mode", "statement", "parsing mode")
var logging = flag.Bool("logging", false, "enable log")

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("usage: ./parse [-mode statement|query|expr|ddl|dml] [-logging] <SQL query>")
	}

	query := flag.Arg(0)

	logf("query: %q", query)

	p := &memefish.Parser{
		Lexer: &memefish.Lexer{
			File: &token.File{FilePath: "", Buffer: query},
		},
	}

	logf("start parsing")
	var node ast.Node
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
	logf("finish parsing successfully")

	fmt.Println("--- AST")
	_, _ = pp.Println(node)
	fmt.Println()
	fmt.Println("--- SQL")
	fmt.Println(node.SQL())
}

func logf(msg string, params ...interface{}) {
	if *logging {
		log.Printf(msg, params...)
	}
}
