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

	"github.com/cloudspannerecosystem/memefish/pkg/parser"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
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

- [ast - GoDoc](https://godoc.org/github.com/cloudspannerecosystem/memefish/pkg/ast)
- [parser - GoDoc](https://godoc.org/github.com/cloudspannerecosystem/memefish/pkg/parser)
- [Query Syntax |  Cloud Spanner  |  Google Cloud](https://cloud.google.com/spanner/docs/query-syntax)