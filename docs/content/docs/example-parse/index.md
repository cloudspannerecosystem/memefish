---
date: 2020-02-03 00:00:00 +0900
title: "Example: parse and unparse"
weight: 1
---

This example shows how to parse a Spanner SQL and unparse it.

 <!--more--> 

 ## Code

```go
package main

import (
	"fmt"
	"log"

	parser "github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/k0kubun/pp"
)

func main() {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: "SELECT * FROM customers",
	}
	p := &parser.Parser{
		Lexer: &parser.Lexer{File: file},
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
```

## Links

- [memefish - GoDoc](https://pkg.go.dev/github.com/cloudspannerecosystem/memefish)
- [ast - GoDoc](https://pkg.go.dev/github.com/cloudspannerecosystem/memefish/ast)
- [Query Syntax |  Cloud Spanner  |  Google Cloud](https://cloud.google.com/spanner/docs/query-syntax)