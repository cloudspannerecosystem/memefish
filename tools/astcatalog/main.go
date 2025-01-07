package main

import (
	"flag"
	"log"

	"github.com/cloudspannerecosystem/memefish/tools/util/astcatalog"
	"github.com/k0kubun/pp/v3"
)

var (
	astfile   = flag.String("astfile", "ast/ast.go", "path to ast/ast.go")
	constfile = flag.String("constfile", "ast/ast_const.go", "path to ast/ast_const.go")
)

func main() {
	flag.Parse()

	catalog, err := astcatalog.Load(*astfile, *constfile)
	if err != nil {
		log.Fatalf("failed to load: %v", err)
	}

	pp.Print(catalog)
}
