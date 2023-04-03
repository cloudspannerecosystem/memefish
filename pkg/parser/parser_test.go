package parser_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/cloudspannerecosystem/memefish/pkg/ast"
	"github.com/cloudspannerecosystem/memefish/pkg/parser"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
	"github.com/k0kubun/pp"
	"github.com/pmezard/go-difflib/difflib"
)

var update = flag.Bool("update", false, "update result files")

func testParser(t *testing.T, inputPath, resultPath string, parse func(p *parser.Parser) (ast.Node, error)) {
	printer := pp.New()
	printer.SetColoringEnabled(false)

	if *update {
		_, err := os.Stat(resultPath)
		if err == nil {
			err = os.RemoveAll(resultPath)
			if err != nil {
				log.Fatalf("error on remove result path: %v", err)
			}
		}
		err = os.MkdirAll(resultPath, 0777)
		if err != nil {
			log.Fatalf("error on create result path: %v", err)
		}
	}

	inputs, err := ioutil.ReadDir(inputPath)
	if err != nil {
		t.Fatalf("error on reading input path: %v", err)
	}

	for _, in := range inputs {
		in := in
		t.Run(in.Name(), func(t *testing.T) {
			t.Parallel()

			b, err := ioutil.ReadFile(filepath.Join(inputPath, in.Name()))
			if err != nil {
				t.Fatalf("error on reading input file: %v", err)
			}

			p := &parser.Parser{
				Lexer: &parser.Lexer{
					File: &token.File{FilePath: in.Name(), Buffer: string(b)},
				},
			}

			node, err := parse(p)
			if err != nil {
				log.Fatalf("error on parsing input file: %v", err)
			}

			var buf bytes.Buffer

			fmt.Fprintf(&buf, "--- %s\n", in.Name())
			fmt.Fprint(&buf, string(b))
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- AST\n")
			_, _ = printer.Fprintln(&buf, node)
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- SQL\n")
			fmt.Fprintln(&buf, node.SQL())

			actual := buf.Bytes()

			if *update {
				t.Log("update " + in.Name() + ".txt")
				err = ioutil.WriteFile(filepath.Join(resultPath, in.Name()+".txt"), buf.Bytes(), 0666)
				if err != nil {
					log.Fatalf("error on writing result file: %v", err)
				}
				return
			}

			expected, err := ioutil.ReadFile(filepath.Join(resultPath, in.Name()+".txt"))
			if err != nil {
				t.Fatalf("error on reading result file: %v", err)
			}

			if !bytes.Equal(actual, expected) {
				diff := difflib.UnifiedDiff{
					A:       difflib.SplitLines(string(expected)),
					B:       difflib.SplitLines(string(actual)),
					Context: 5,
				}
				d, err := difflib.GetUnifiedDiffString(diff)
				if err != nil {
					t.Fatalf("error on diff: %v", err)
				}
				t.Error(d)
				return
			}

			s1 := node.SQL()
			p1 := &parser.Parser{
				Lexer: &parser.Lexer{
					File: &token.File{FilePath: in.Name() + " (SQL)", Buffer: s1},
				},
			}

			node1, err := parse(p1)
			if err != nil {
				log.Fatalf("error on parsing unparsed SQL: %v", err)
			}

			s2 := node1.SQL()
			if s1 != s2 {
				t.Errorf("%q != %q", s1, s2)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	inputPath := "./testdata/input/query"
	resultPath := "./testdata/result/query"

	testParser(t, inputPath, resultPath, func(p *parser.Parser) (ast.Node, error) {
		return p.ParseQuery()
	})
}

func TestParseDDL(t *testing.T) {
	inputPath := "./testdata/input/ddl"
	resultPath := "./testdata/result/ddl"

	testParser(t, inputPath, resultPath, func(p *parser.Parser) (ast.Node, error) {
		return p.ParseDDL()
	})
}

func TestParseDML(t *testing.T) {
	inputPath := "./testdata/input/dml"
	resultPath := "./testdata/result/dml"

	testParser(t, inputPath, resultPath, func(p *parser.Parser) (ast.Node, error) {
		return p.ParseDML()
	})
}

func TestParseStatement(t *testing.T) {
	inputPaths := []string{
		"./testdata/input/query",
		"./testdata/input/ddl",
		"./testdata/input/dml",
	}
	resultPath := "./testdata/result/statement"

	for _, inputPath := range inputPaths {
		testParser(t, inputPath, resultPath, func(p *parser.Parser) (ast.Node, error) {
			return p.ParseStatement()
		})
	}
}

func ExampleParser_ParseStatements() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{Buffer: "SELECT 1; INSERT foo (x, y) VALUES (1, 2)"},
		},
	}

	stmts, err := p.ParseStatements()
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

func ExampleParser_ParseQuery() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{Buffer: "SELECT * FROM foo"},
		},
	}

	stmt, err := p.ParseQuery()
	if err != nil {
		panic(err)
	}

	fmt.Println(stmt.SQL())

	// Output:
	// SELECT * FROM foo
}

func ExampleParser_ParseDDL() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{
				Buffer: heredoc.Doc(`
					CREATE TABLE foo (
						x int64,
						y int64,
					) PRIMARY KEY (x)
				`),
			},
		},
	}

	ddl, err := p.ParseDDL()
	if err != nil {
		panic(err)
	}

	fmt.Println(ddl.SQL())

	// Output:
	// CREATE TABLE foo (x INT64, y INT64) PRIMARY KEY (x)
}

func ExampleParser_ParseDDLs() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{
				Buffer: heredoc.Doc(`
					CREATE TABLE foo (x int64, y int64) PRIMARY KEY (x);

					CREATE TABLE bar (
						x int64, z int64,
					)
					PRIMARY KEY (x, z),
					INTERLEAVE IN PARENT foo;
				`),
			},
		},
	}

	ddls, err := p.ParseDDLs()
	if err != nil {
		panic(err)
	}

	for _, ddl := range ddls {
		fmt.Printf("%s;\n", ddl.SQL())
	}

	// Output:
	// CREATE TABLE foo (x INT64, y INT64) PRIMARY KEY (x);
	// CREATE TABLE bar (x INT64, z INT64) PRIMARY KEY (x, z), INTERLEAVE IN PARENT foo;
}

func ExampleParser_ParseDML() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{
				Buffer: heredoc.Doc(`
					INSERT INTO foo (x, y)
					VALUES (1, 2),
					       (3, 4)
				`),
			},
		},
	}

	dml, err := p.ParseDML()
	if err != nil {
		panic(err)
	}

	fmt.Println(dml.SQL())

	// Output:
	// INSERT INTO foo (x, y) VALUES (1, 2), (3, 4)
}

func ExampleParser_ParseDMLs() {
	p := &parser.Parser{
		Lexer: &parser.Lexer{
			File: &token.File{
				Buffer: heredoc.Doc(`
					INSERT INTO foo (x, y) VALUES (1, 2), (3, 4);
					DELETE FROM foo WHERE foo.x = 1 AND foo.y = 2;
				`),
			},
		},
	}

	dmls, err := p.ParseDMLs()
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
