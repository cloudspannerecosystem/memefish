package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/cloudspannerecosystem/memefish/tools/util/poslang"
	"github.com/k0kubun/pp"
)

var mode = flag.String("mode", "statement", "parsing mode")
var logging = flag.Bool("logging", false, "enable log")
var dig = flag.String("dig", "", "digging the result node before printing")
var pos = flag.String("pos", "", "POS expression")

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

	if *dig != "" {
		value := reflect.ValueOf(node)
		for _, name := range strings.Split(*dig, ".") {
			if value.Kind() == reflect.Interface {
				value = value.Elem()
			}

			if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
				index, err := strconv.Atoi(name)
				if err != nil {
					log.Fatal(err)
				}
				value = value.Index(index)
				continue
			}

			value = value.Elem().FieldByName(name)
		}

		if n, ok := value.Interface().(ast.Node); ok {
			node = n
		} else {
			log.Fatalf("invalid value: %#v", value)
		}
	}

	fmt.Println("--- AST")
	_, _ = pp.Println(node)
	fmt.Println()
	fmt.Println("--- SQL")
	fmt.Println(node.SQL())

	if *pos != "" {
		expr, err := poslang.Parse(*pos)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("--- Pos")
		fmt.Println(expr.EvalPos(node))
	}
}

func logf(msg string, params ...interface{}) {
	if *logging {
		log.Printf(msg, params...)
	}
}
