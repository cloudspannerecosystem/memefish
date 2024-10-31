package memefish_test

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/k0kubun/pp"
	"github.com/pmezard/go-difflib/difflib"
)

var update = flag.Bool("update", false, "update result files")

func testParser(t *testing.T, inputPath, resultPath string, parse func(p *memefish.Parser) (ast.Node, error)) {
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

	inputs, err := os.ReadDir(inputPath)
	if err != nil {
		t.Fatalf("error on reading input path: %v", err)
	}

	for _, in := range inputs {
		in := in
		bad := strings.HasPrefix(in.Name(), "bad_")
		t.Run(in.Name(), func(t *testing.T) {
			t.Parallel()

			inputFilePath := filepath.Join(inputPath, in.Name())
			b, err := os.ReadFile(inputFilePath)
			if err != nil {
				t.Fatalf("error on reading input file: %v", err)
			}

			p := &memefish.Parser{
				Lexer: &memefish.Lexer{
					File: &token.File{FilePath: inputFilePath, Buffer: string(b)},
				},
			}

			node, err := parse(p)
			var buf bytes.Buffer

			fmt.Fprintf(&buf, "--- %s\n", in.Name())
			fmt.Fprint(&buf, string(b))
			fmt.Fprintln(&buf)

			if err != nil {
				if bad {
					fmt.Fprintln(&buf, "--- Error")
					fmt.Fprint(&buf, err)
					fmt.Fprintln(&buf)
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if bad {
					t.Errorf("error is expected, but parsing succeeded")
				}
			}

			fmt.Fprintf(&buf, "--- AST\n")
			_, _ = printer.Fprintln(&buf, node)
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- SQL\n")
			fmt.Fprintln(&buf, node.SQL())

			actual := buf.Bytes()

			if *update {
				t.Log("update " + in.Name() + ".txt")
				err = os.WriteFile(filepath.Join(resultPath, in.Name()+".txt"), buf.Bytes(), 0666)
				if err != nil {
					log.Fatalf("error on writing result file: %v", err)
				}
				return
			}

			expected, err := os.ReadFile(filepath.Join(resultPath, in.Name()+".txt"))
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
			p1 := &memefish.Parser{
				Lexer: &memefish.Lexer{
					File: &token.File{FilePath: in.Name() + " (SQL)", Buffer: s1},
				},
			}

			node1, _ := parse(p1)

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

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseQuery()
	})
}

func TestParseDDL(t *testing.T) {
	inputPath := "./testdata/input/ddl"
	resultPath := "./testdata/result/ddl"

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseDDL()
	})
}

func TestParseDML(t *testing.T) {
	inputPath := "./testdata/input/dml"
	resultPath := "./testdata/result/dml"

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseDML()
	})
}

func TestParseExpr(t *testing.T) {
	inputPath := "./testdata/input/expr"
	resultPath := "./testdata/result/expr"

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseExpr()
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
		testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
			return p.ParseStatement()
		})
	}
}
