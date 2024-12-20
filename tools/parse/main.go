package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/cloudspannerecosystem/memefish/tools/util/poslang"
	"github.com/k0kubun/pp/v3"
)

var usage = heredoc.Doc(`
	Usage of tools/parse.go

	A testing tool for parsing Spanner SQL.

	Example:

	  $ go run ./tools/parse/main.go "SELECT 1 AS x"
	        Show the parse result of "SELECT 1 AS X".
	  $ go run ./tools/parse/main.go -mode expr "(SELECT 1) + 2"
	        Parse "(SELECT 1) + 2" on the expression mode.
	  $ go run ./tools/parse/main.go -pos "Query.end" "SELECT 1 AS x"
	        Evaluate the POS expression "Query.end" on "SELECT 1 AS x"
	  $ go run ./tools/parse/main.go -pos "As.end" -dig "Query.Results.0" "SELECT 1 AS x"
	        Evaluate the POS expression "As.end" on "1 AS x" of "SELECT 1 AS x"

	Options:
`)

var (
	mode    = flag.String("mode", "statement", `parsing mode (one of "statement", "query", "expr", "ddl", "dml")`)
	logging = flag.Bool("logging", false, "enable log (default: false)")
	dig     = flag.String("dig", "", "digging the result node before printing")
	pos     = flag.String("pos", "", "POS expression")
)

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
		return
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
	case "type":
		node, err = p.ParseType()
	case "ddl":
		node, err = p.ParseDDL()
	case "dml":
		node, err = p.ParseDML()
	default:
		log.Fatalf("unknown mode: %s", *mode)
	}
	logf("finish parsing successfully")

	if err != nil {
		fmt.Println("--- Error")
		fmt.Print(err)
		fmt.Println()
	}

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
	pprinter := pp.New()
	pprinter.SetOmitEmpty(true)
	_, _ = pprinter.Println(node)
	fmt.Println()

	fmt.Println("--- SQL")
	fmt.Println(node.SQL())

	if *pos != "" {
		fmt.Println("--- POS")

		expr, err := poslang.Parse(*pos)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(expr.EvalPos(node))
	}
}

func logf(msg string, params ...interface{}) {
	if *logging {
		log.Printf(msg, params...)
	}
}
