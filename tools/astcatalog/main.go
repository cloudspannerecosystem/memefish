package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish/tools/util/astcatalog"
	"github.com/k0kubun/pp/v3"
)

var usage = heredoc.Doc(`
	Usage of tools/astcatalog

	An utility to show the AST catalog.

	Example:

	  $ go run ./tools/astcatalog/main.go -astfile ast/ast.go -constfile ast/ast_const.go
	    Print the AST catalog of ast/ast.go and ast/ast_const.go.

	Flags:
`)

var (
	astfile   = flag.String("astfile", "ast/ast.go", "path to ast/ast.go")
	constfile = flag.String("constfile", "ast/ast_const.go", "path to ast/ast_const.go")
)

func main() {
	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}

	flag.Parse()

	catalog, err := astcatalog.Load(*astfile, *constfile)
	if err != nil {
		log.Fatalf("failed to load: %v", err)
	}
	pprinter := pp.New()
	pprinter.SetOmitEmpty(true)
	_, _ = pprinter.Println(catalog)
}
