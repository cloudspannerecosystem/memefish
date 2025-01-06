package main

import (
	"flag"
	"log"

	"github.com/cloudspannerecosystem/memefish/tools/util/astcatalog"
	"github.com/k0kubun/pp/v3"
)

var (
	infile = flag.String("infile", "", "input filename")
)

func main() {
	flag.Parse()

	catalog, err := astcatalog.Load(*infile)
	if err != nil {
		log.Fatalf("failed to load: %v", err)
	}

	pp.Print(catalog)
}
