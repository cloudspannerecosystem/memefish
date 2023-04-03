---
date: 2020-02-03 00:00:00 +0900
title: "Example: analyze"
weight: 2
---

This example shows how to analyze a Spanner SQL after parsing.

 <!--more--> 

## Code

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudspannerecosystem/memefish/pkg/analyzer"
	"github.com/cloudspannerecosystem/memefish/pkg/parser"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
)

func main() {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: "SELECT * FROM singers",
	}
	p := &parser.Parser{
		Lexer: &parser.Lexer{File: file},
	}

	// Do parsing!
	stmt, err := p.ParseQuery()
	if err != nil {
		log.Fatal(err)
	}

	// Create table catalog.
	catalog := &analyzer.Catalog{
		Tables: map[string]*analyzer.TableSchema{
			"SINGERS": {
				Name: "Singers",
				Columns: []*analyzer.ColumnSchema{
					{Name: "SingerId", Type: analyzer.Int64Type},
					{Name: "FirstName", Type: analyzer.StringType},
					{Name: "LastName", Type: analyzer.StringType},
				},
			},
		},
	}

	// Create a new Analyzer instance.
	a := &analyzer.Analyzer{
		File:    file,
		Catalog: catalog,
	}

	// Analyze!
	err = a.AnalyzeQueryStatement(stmt)
	if err != nil {
		log.Fatal(err)
	}

	// Get first column information.
	columns := a.NameLists[stmt.Query]
	fmt.Printf("1st column name  : %s\n", columns[0].Text)
	fmt.Printf("1st column type  : %s\n", columns[0].Type)
	fmt.Printf("1st column schema: %#v\n", columns[0].Deref().ColumnSchema) // == catalog.Tables["SINGERS"].Columns[0]
}
```

## Links

- [analyzer - GoDoc](https://godoc.org/github.com/cloudspannerecosystem/memefish/pkg/analyzer)