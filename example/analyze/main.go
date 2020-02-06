package main

import (
	"fmt"
	"log"

	"github.com/MakeNowJust/memefish/pkg/analyzer"
	"github.com/MakeNowJust/memefish/pkg/parser"
	"github.com/MakeNowJust/memefish/pkg/token"
)

func main() {
	// Create a new Parser instance.
	file := &token.File{
		Buffer: "SELECT * FROM singers WHERE CAST(@param As STRUCT<foo STRING>).foo = FirstName",
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

	placeholder := analyzer.NewPlaceholder()
	params := map[string]interface{}{
		"PARAM": placeholder,
	}

	// Create a new Analyzer instance.
	a := &analyzer.Analyzer{
		File:    file,
		Catalog: catalog,
		Params:  params,
	}

	// Analyze!
	err = a.AnalyzeQueryStatement(stmt)
	if err != nil {
		log.Fatal(err)
	}

	// Get columns information.
	columns := a.NameLists[stmt.Query]
	for i, column := range columns {
		fmt.Printf("columns[%d] name  : %s\n", i, column.Text)
		fmt.Printf("columns[%d] type  : %s\n", i, column.Type)
		fmt.Printf("columns[%d] schema: %#v\n", i, column.Deref().ColumnSchema) // == catalog.Tables["SINGERS"].Columns[i]
		fmt.Println()
	}

	// Get inferred placeholder type.
	fmt.Printf("placeholder type : %s\n", placeholder.PlaceholderType)
}
