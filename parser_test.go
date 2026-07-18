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

	"github.com/k0kubun/pp/v3"
	"github.com/pmezard/go-difflib/difflib"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

var update = flag.Bool("update", false, "update result files")

type pathVisitor struct {
	f    func(path string, node ast.Node) bool
	path string
}

func (v *pathVisitor) Visit(node ast.Node) ast.Visitor {
	ok := v.f(v.path, node)
	if !ok {
		return nil
	}

	return &pathVisitor{
		f:    v.f,
		path: v.path,
	}
}

func (v *pathVisitor) VisitMany(nodes []ast.Node) ast.Visitor {
	return v
}

func (v *pathVisitor) Field(name string) ast.Visitor {
	return &pathVisitor{
		f:    v.f,
		path: fmt.Sprintf("%s.%s", v.path, name),
	}
}

func (v *pathVisitor) Index(index int) ast.Visitor {
	return &pathVisitor{
		f:    v.f,
		path: fmt.Sprintf("%s[%d]", v.path, index),
	}
}

func testParser(t *testing.T, inputPath, resultPath string, parse func(p *memefish.Parser) (ast.Node, error)) {
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
		bad := strings.HasPrefix(in.Name(), "!bad_")
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
			if !bad && node == nil {
				t.Fatal("parser returned a nil AST without an expected parse error")
			}

			pprinter := pp.New()
			pprinter.SetColoringEnabled(false)
			pprinter.SetOmitEmpty(true)

			var buf bytes.Buffer

			fmt.Fprintf(&buf, "--- %s\n", in.Name())
			fmt.Fprint(&buf, string(b))
			fmt.Fprintln(&buf)

			if err != nil {
				list, ok := err.(memefish.MultiError)
				if bad && ok {
					fmt.Fprintf(&buf, "--- Error\n%s\n\n", list.FullError())
				} else {
					t.Errorf("unexpected error: %v", err)
				}
			} else {
				if bad {
					t.Errorf("error is expected, but parsing succeeded")
				}
			}

			// pos/end test
			if !bad {
				ok := true

				ast.Walk(node, &pathVisitor{
					path: "",
					f: func(path string, node ast.Node) bool {
						if !ok {
							return false
						}

						if node.Pos().Invalid() {
							t.Errorf("pos must be valid, but got invalid on %v: %v", path, node.SQL())
							ok = false
							return false
						}

						if node.End().Invalid() {
							t.Errorf("end must be valid, but got invalid on %v: %v", path, node.SQL())
							ok = false
							return false
						}

						if node.End() <= node.Pos() {
							t.Errorf("pos must be smaller than end, but got pos: %v, end: %v on %v: %v", node.Pos(), node.End(), path, node.SQL())
							ok = false
							return false
						}

						// FIXME: The fields of `CreateTable` are not ordered by position for now,
						// so we skips to check the order of positions of `CreateTable` fields.
						_, isCreateTable := node.(*ast.CreateTable)

						lastEnd := token.InvalidPos
						ast.Walk(node, &pathVisitor{
							path: path,
							f: func(childPath string, child ast.Node) bool {
								if !ok {
									return false
								}

								if node == child {
									return true
								}

								if child.Pos() < lastEnd {
									if !isCreateTable {
										t.Errorf("pos must be larger or equal than end of last node pos %v, but got pos: %v on %v: %v", lastEnd, child.Pos(), childPath, child.SQL())
									} else {
										t.Logf("pos must be larger or equal than end of last node pos %v, but got pos: %v on %v: %v", lastEnd, child.Pos(), childPath, child.SQL())
									}
									ok = false
									return false
								}
								lastEnd = child.End()

								if child.Pos() < node.Pos() || node.End() < child.End() {
									t.Errorf("child position must be in node position [%v, %v], but got [%v, %v] on %v: %v", node.Pos(), node.End(), child.Pos(), child.End(), childPath, child.SQL())
									ok = false
								}

								return false
							},
						})

						return ok
					},
				})
			}

			fmt.Fprintf(&buf, "--- AST\n")
			_, _ = pprinter.Fprintln(&buf, node)
			fmt.Fprintln(&buf)

			fmt.Fprintf(&buf, "--- SQL\n")
			if node != nil {
				fmt.Fprintln(&buf, node.SQL())
			}

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
			if node == nil {
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
		"./testdata/input/gql",
		"./testdata/input/statement",
	}
	resultPath := "./testdata/result/statement"

	for _, inputPath := range inputPaths {
		testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
			return p.ParseStatement()
		})
	}
}

func TestParseGQLQuery(t *testing.T) {
	inputPath := "./testdata/input/gql"
	resultPath := "./testdata/result/gql"

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseGQLQuery()
	})
}

func TestParseGQLGraphPattern(t *testing.T) {
	inputPath := "./testdata/input/gql_graph_pattern"
	resultPath := "./testdata/result/gql_graph_pattern"

	testParser(t, inputPath, resultPath, func(p *memefish.Parser) (ast.Node, error) {
		return p.ParseGQLGraphPattern()
	})
}

func TestParseGQLGraphPatternMalformedReturnsBadNode(t *testing.T) {
	const input = "(a)-@"
	p := &memefish.Parser{
		Lexer: &memefish.Lexer{
			File: &token.File{Buffer: input},
		},
	}

	pattern, err := p.ParseGQLGraphPattern()
	if err == nil {
		t.Fatal("ParseGQLGraphPattern() error = nil, want error")
	}
	bad, ok := pattern.(*ast.BadGQLGraphPattern)
	if !ok {
		t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.BadGQLGraphPattern", pattern)
	}
	if got := bad.SQL(); got != input {
		t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, input)
	}
	if bad.Pos() != 0 || bad.End() != token.Pos(len(input)) {
		t.Errorf("ParseGQLGraphPattern() range = [%d, %d), want [0, %d)", bad.Pos(), bad.End(), len(input))
	}
	ast.Inspect(bad, func(ast.Node) bool { return true })
}

func TestParseGQLGraphPatternFirstTokenLexerErrorReturnsBadNode(t *testing.T) {
	const input = `"foo`
	p := &memefish.Parser{
		Lexer: &memefish.Lexer{
			File: &token.File{Buffer: input},
		},
	}

	pattern, err := p.ParseGQLGraphPattern()
	if err == nil {
		t.Fatal("ParseGQLGraphPattern() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "unclosed string literal") {
		t.Fatalf("ParseGQLGraphPattern() error = %v, want unclosed string literal", err)
	}
	bad, ok := pattern.(*ast.BadGQLGraphPattern)
	if !ok {
		t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.BadGQLGraphPattern", pattern)
	}
	if got := bad.SQL(); got != input {
		t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, input)
	}
	ast.Inspect(bad, func(ast.Node) bool { return true })

	p.Lexer = &memefish.Lexer{File: &token.File{Buffer: "(a)"}}
	pattern, err = p.ParseGQLGraphPattern()
	if err != nil {
		t.Fatalf("reused ParseGQLGraphPattern() error = %v, want nil", err)
	}
	if _, ok := pattern.(*ast.GQLGraphPattern); !ok {
		t.Fatalf("reused ParseGQLGraphPattern() pattern = %T, want *ast.GQLGraphPattern", pattern)
	}
}

func TestParseGQLGraphPatternHelper(t *testing.T) {
	const input = "(a)-[e]->(b)"
	pattern, err := memefish.ParseGQLGraphPattern("", input)
	if err != nil {
		t.Fatalf("ParseGQLGraphPattern() error = %v, want nil", err)
	}
	if _, ok := pattern.(*ast.GQLGraphPattern); !ok {
		t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.GQLGraphPattern", pattern)
	}
	if got := pattern.SQL(); got != input {
		t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, input)
	}
}

func TestParseGQLGraphPatternHelperMalformedReturnsBadNode(t *testing.T) {
	const input = "(a) trailing"
	pattern, err := memefish.ParseGQLGraphPattern("", input)
	if err == nil {
		t.Fatal("ParseGQLGraphPattern() error = nil, want error")
	}
	bad, ok := pattern.(*ast.BadGQLGraphPattern)
	if !ok {
		t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.BadGQLGraphPattern", pattern)
	}
	if got := bad.SQL(); got != input {
		t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, input)
	}
}

func TestParseGQLGraphPatternSupportedForms(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "path variable with search prefix", input: "p = ANY SHORTEST (a)->(b)", want: "p = ANY SHORTEST (a)->(b)"},
		{name: "prefix-like path variable", input: "shortest = (a)", want: "shortest = (a)"},
		{name: "path mode suffix is canonical", input: "p = trail paths (a)", want: "p = TRAIL PATHS (a)"},
		{name: "edge-only subpath", input: "(->)", want: "(->)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := memefish.ParseGQLGraphPattern("", tt.input)
			if err != nil {
				t.Fatalf("ParseGQLGraphPattern() error = %v, want nil", err)
			}
			if _, ok := pattern.(*ast.GQLGraphPattern); !ok {
				t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.GQLGraphPattern", pattern)
			}
			if got := pattern.SQL(); got != tt.want {
				t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseGQLGraphPatternRejectsEmptyProperties(t *testing.T) {
	const input = "(a {})"
	pattern, err := memefish.ParseGQLGraphPattern("", input)
	if err == nil {
		t.Fatal("ParseGQLGraphPattern() error = nil, want error")
	}
	bad, ok := pattern.(*ast.BadGQLGraphPattern)
	if !ok {
		t.Fatalf("ParseGQLGraphPattern() pattern = %T, want *ast.BadGQLGraphPattern", pattern)
	}
	if got := bad.SQL(); got != input {
		t.Errorf("ParseGQLGraphPattern() SQL() = %q, want %q", got, input)
	}
}
