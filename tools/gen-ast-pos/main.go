package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"unicode"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish/tools/util/astcatalog"
	"github.com/cloudspannerecosystem/memefish/tools/util/poslang"
)

var (
	usage = heredoc.Doc(`
		Usage of tools/gen-ast-pos.go

		A generator of ast/pos.go.

		Example:

		  $ go run ./tools/gen-ast-pos/main.go -infile ast/ast.go
		        Print the generated ast/pos.go to stdout.

		Flags:
	`)
	prologue = heredoc.Doc(`
		// Code generated by tools/gen-ast-pos; DO NOT EDIT.

		package ast

		import (
			"github.com/cloudspannerecosystem/memefish/token"
		)
	`)
)

var (
	astfile   = flag.String("astfile", "ast/ast.go", "path to ast/ast.go")
	constfile = flag.String("constfile", "ast/ast_const.go", "path to ast/ast_const.go")
	outfile   = flag.String("outfile", "", "output filename (if it is not specified, the result is printed to stdout.)")
)

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	catalog, err := astcatalog.Load(*astfile, *constfile)
	if err != nil {
		log.Fatal(err)
	}

	structs := make([]*astcatalog.NodeStructDef, 0, len(catalog.Structs))
	for _, structDef := range catalog.Structs {
		structs = append(structs, structDef)
	}
	sort.Slice(structs, func(i, j int) bool {
		return structs[i].SourcePos < structs[j].SourcePos
	})

	var buffer bytes.Buffer
	buffer.WriteString(prologue)

	for _, structDef := range structs {
		x := string(unicode.ToLower(rune(structDef.Name[0])))

		posExpr, err := poslang.Parse(structDef.Pos)
		if err != nil {
			log.Fatalf("error on parsing pos: %v", err)
		}

		endExpr, err := poslang.Parse(structDef.End)
		if err != nil {
			log.Fatalf("error on parsing pos: %v", err)
		}

		fmt.Fprintln(&buffer)
		fmt.Fprintf(&buffer, "func (%s *%s) Pos() token.Pos {\n", x, structDef.Name)
		fmt.Fprintf(&buffer, "\treturn %s\n", posExpr.PosExprToGo(x))
		fmt.Fprintf(&buffer, "}\n")
		fmt.Fprintln(&buffer)
		fmt.Fprintf(&buffer, "func (%s *%s) End() token.Pos {\n", x, structDef.Name)
		fmt.Fprintf(&buffer, "\treturn %s\n", endExpr.PosExprToGo(x))
		fmt.Fprintf(&buffer, "}\n")
	}

	if *outfile == "" {
		fmt.Print(buffer.String())
		return
	}

	err = os.WriteFile(*outfile, buffer.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
