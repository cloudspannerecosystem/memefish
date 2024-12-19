package memefish_test

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"

	"github.com/cloudspannerecosystem/memefish"
)

func ExampleParseStatement() {
	stmt, err := memefish.ParseStatement("path/to/file.sql", "SELECT * FROM foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(stmt.SQL())

	// Output:
	// SELECT * FROM foo
}

func ExampleParseStatements() {
	stmts, err := memefish.ParseStatements("path/to/file.sql", "SELECT 1; INSERT foo (x, y) VALUES (1, 2)")
	if err != nil {
		panic(err)
	}

	for _, stmt := range stmts {
		fmt.Printf("%s;\n", stmt.SQL())
	}

	// Output:
	// SELECT 1;
	// INSERT INTO foo (x, y) VALUES (1, 2);
}

func ExampleParseQuery() {
	stmt, err := memefish.ParseQuery("path/to/file.sql", "SELECT * FROM foo")
	if err != nil {
		panic(err)
	}

	fmt.Println(stmt.SQL())

	// Output:
	// SELECT * FROM foo
}

func ExampleParseExpr() {
	stmt, err := memefish.ParseExpr("", "a * b")
	if err != nil {
		panic(err)
	}

	fmt.Println(stmt.SQL())

	// Output:
	// a * b
}

func ExampleParseType() {
	stmt, err := memefish.ParseType("", "ARRAY<STRUCT<n INT64, s STRING>>")
	if err != nil {
		panic(err)
	}

	fmt.Println(stmt.SQL())

	// Output:
	// ARRAY<STRUCT<n INT64, s STRING>>
}
func ExampleParseDDL() {
	sql := heredoc.Doc(`
		CREATE TABLE foo (
			x int64,
			y int64,
		) PRIMARY KEY (x)
	`)

	ddl, err := memefish.ParseDDL("path/to/file.sql", sql)
	if err != nil {
		panic(err)
	}

	fmt.Println(ddl.SQL())

	// Output:
	// CREATE TABLE foo (
	//     x INT64,
	//     y INT64
	// ) PRIMARY KEY (x)
}

func ExampleParseDDLs() {
	sql := heredoc.Doc(`
		CREATE TABLE foo (x int64, y int64) PRIMARY KEY (x);

		CREATE TABLE bar (
			x int64, z int64,
		)
		PRIMARY KEY (x, z),
		INTERLEAVE IN PARENT foo;
	`)

	ddls, err := memefish.ParseDDLs("path/to/file.sql", sql)
	if err != nil {
		panic(err)
	}

	for _, ddl := range ddls {
		fmt.Printf("%s;\n", ddl.SQL())
	}

	// Output:
	// CREATE TABLE foo (
	//     x INT64,
	//     y INT64
	// ) PRIMARY KEY (x);
	// CREATE TABLE bar (
	//     x INT64,
	//     z INT64
	// ) PRIMARY KEY (x, z),
	// INTERLEAVE IN PARENT foo;
}

func ExampleParseDML() {
	sql := heredoc.Doc(`
		INSERT INTO foo (x, y)
		VALUES (1, 2),
		       (3, 4)
	`)

	dml, err := memefish.ParseDML("path/to/file.sql", sql)
	if err != nil {
		panic(err)
	}

	fmt.Println(dml.SQL())

	// Output:
	// INSERT INTO foo (x, y) VALUES (1, 2), (3, 4)
}

func ExampleParseDMLs() {
	sql := heredoc.Doc(`
		INSERT INTO foo (x, y) VALUES (1, 2), (3, 4);
		DELETE FROM foo WHERE foo.x = 1 AND foo.y = 2;
	`)

	dmls, err := memefish.ParseDMLs("path/to/file.sql", sql)
	if err != nil {
		panic(err)
	}

	for _, dml := range dmls {
		fmt.Printf("%s;\n", dml.SQL())
	}

	// Output:
	// INSERT INTO foo (x, y) VALUES (1, 2), (3, 4);
	// DELETE FROM foo WHERE foo.x = 1 AND foo.y = 2;
}
