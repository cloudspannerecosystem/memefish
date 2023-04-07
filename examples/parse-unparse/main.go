package main

import (
	"fmt"
	"log"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/k0kubun/pp"
)

func main() {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: "SELECT * FROM customers",
	}
	p := &memefish.Parser{
		Lexer: &memefish.Lexer{File: file},
	}

	// Do parsing!
	stmt, err := p.ParseQuery()
	if err != nil {
		log.Fatal(err)
	}

	// Show AST.
	log.Print("AST")
	_, _ = pp.Println(stmt)

	// Unparse AST to SQL source string.
	log.Print("Unparse")
	fmt.Println(stmt.SQL())
}
