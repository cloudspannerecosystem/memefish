package memefish_test

import (
	"testing"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

type parserEntryRoute struct {
	name   string
	valid  string
	parse  func(*memefish.Parser) ([]ast.Node, error)
	helper func(string) ([]ast.Node, error)
}

func publicParserEntryRoutes() []parserEntryRoute {
	return []parserEntryRoute{
		{
			name:  "ParseStatement",
			valid: "SELECT 1",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseStatement()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseStatement("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseStatements",
			valid: "SELECT 1",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				nodes, err := p.ParseStatements()
				return toASTNodes(nodes), err
			},
			helper: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseStatements("", input)
				return toASTNodes(nodes), err
			},
		},
		{
			name:  "ParseQuery",
			valid: "SELECT 1",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseQuery()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseQuery("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseGQLQuery",
			valid: "GRAPH FinGraph RETURN 1",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseGQLQuery()
				if node == nil {
					return nil, err
				}
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseExpr",
			valid: "1",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseExpr()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseExpr("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseType",
			valid: "INT64",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseType()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseType("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseSchemaType",
			valid: "STRING(MAX)",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseSchemaType()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseSchemaType("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseDDL",
			valid: "DROP TABLE t",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseDDL()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseDDL("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseDDLs",
			valid: "DROP TABLE t",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				nodes, err := p.ParseDDLs()
				return toASTNodes(nodes), err
			},
			helper: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseDDLs("", input)
				return toASTNodes(nodes), err
			},
		},
		{
			name:  "ParseDML",
			valid: "DELETE FROM t WHERE TRUE",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				node, err := p.ParseDML()
				return []ast.Node{node}, err
			},
			helper: func(input string) ([]ast.Node, error) {
				node, err := memefish.ParseDML("", input)
				return []ast.Node{node}, err
			},
		},
		{
			name:  "ParseDMLs",
			valid: "DELETE FROM t WHERE TRUE",
			parse: func(p *memefish.Parser) ([]ast.Node, error) {
				nodes, err := p.ParseDMLs()
				return toASTNodes(nodes), err
			},
			helper: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseDMLs("", input)
				return toASTNodes(nodes), err
			},
		},
	}
}

func TestPublicParserEntryFirstTokenRecovery(t *testing.T) {
	malformed := []struct {
		name  string
		input string
	}{
		{name: "illegal byte", input: "\x00"},
		{name: "unclosed string", input: `"foo`},
		{name: "unclosed comment", input: "/*"},
	}

	for _, route := range publicParserEntryRoutes() {
		for _, test := range malformed {
			t.Run(route.name+"/method/"+test.name, func(t *testing.T) {
				p := newEntryTestParser(test.input)
				nodes, err := parseEntryWithoutPanic(t, func() ([]ast.Node, error) {
					return route.parse(p)
				})
				if err == nil {
					t.Fatal("parse error = nil, want an error")
				}
				requireSafeEntryNodes(t, nodes)

				p.Lexer = newEntryTestLexer(route.valid)
				nodes, err = parseEntryWithoutPanic(t, func() ([]ast.Node, error) {
					return route.parse(p)
				})
				if err != nil {
					t.Fatalf("parser reuse error = %v, want nil", err)
				}
				requireSafeEntryNodes(t, nodes)
			})

			if route.helper != nil {
				t.Run(route.name+"/helper/"+test.name, func(t *testing.T) {
					nodes, err := parseEntryWithoutPanic(t, func() ([]ast.Node, error) {
						return route.helper(test.input)
					})
					if err == nil {
						t.Fatal("parse error = nil, want an error")
					}
					requireSafeEntryNodes(t, nodes)
				})
			}
		}
	}
}

func TestPublicParserEntryPostSemicolonRecovery(t *testing.T) {
	routes := []struct {
		name   string
		prefix string
		parse  func(string) ([]ast.Node, error)
	}{
		{
			name:   "ParseStatements",
			prefix: "SELECT 1; ",
			parse: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseStatements("", input)
				return toASTNodes(nodes), err
			},
		},
		{
			name:   "ParseDDLs",
			prefix: "DROP TABLE t; ",
			parse: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseDDLs("", input)
				return toASTNodes(nodes), err
			},
		},
		{
			name:   "ParseDMLs",
			prefix: "DELETE FROM t WHERE TRUE; ",
			parse: func(input string) ([]ast.Node, error) {
				nodes, err := memefish.ParseDMLs("", input)
				return toASTNodes(nodes), err
			},
		},
	}

	malformed := []struct {
		name  string
		input string
	}{
		{name: "illegal byte", input: "\x00"},
		{name: "unclosed string", input: `"foo`},
		{name: "unclosed comment", input: "/*"},
	}

	for _, route := range routes {
		for _, test := range malformed {
			t.Run(route.name+"/"+test.name, func(t *testing.T) {
				nodes, err := parseEntryWithoutPanic(t, func() ([]ast.Node, error) {
					return route.parse(route.prefix + test.input)
				})
				if err == nil {
					t.Fatal("parse error = nil, want an error")
				}
				if len(nodes) != 2 {
					t.Fatalf("len(nodes) = %d, want 2", len(nodes))
				}
				requireSafeEntryNodes(t, nodes)
			})
		}
	}
}

func newEntryTestParser(input string) *memefish.Parser {
	return &memefish.Parser{Lexer: newEntryTestLexer(input)}
}

func newEntryTestLexer(input string) *memefish.Lexer {
	return &memefish.Lexer{File: &token.File{Buffer: input}}
}

func toASTNodes[T ast.Node](nodes []T) []ast.Node {
	result := make([]ast.Node, len(nodes))
	for i, node := range nodes {
		result[i] = node
	}
	return result
}

func parseEntryWithoutPanic(t *testing.T, parse func() ([]ast.Node, error)) (nodes []ast.Node, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parse panicked: %v", r)
		}
	}()
	return parse()
}

func requireSafeEntryNodes(t *testing.T, nodes []ast.Node) {
	t.Helper()
	for i, node := range nodes {
		if node == nil {
			continue
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("node %d is unsafe: %v", i, r)
				}
			}()
			_ = node.SQL()
			_ = node.Pos()
			_ = node.End()
			ast.Inspect(node, func(ast.Node) bool { return true })
		}()
	}
}
