package memefish

import (
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

func newParser(filepath, s string) *Parser {
	return &Parser{
		Lexer: &Lexer{
			File: &token.File{FilePath: filepath, Buffer: s},
		},
	}
}

func ParseStatement(filepath, s string) (ast.Statement, error) {
	return newParser(filepath, s).ParseStatement()
}

func ParseStatements(filepath, s string) ([]ast.Statement, error) {
	return newParser(filepath, s).ParseStatements()
}

func ParseQuery(filepath, s string) (*ast.QueryStatement, error) {
	return newParser(filepath, s).ParseQuery()
}

func ParseExpr(filepath, s string) (ast.Expr, error) {
	return newParser(filepath, s).ParseExpr()
}

func ParseType(filepath, s string) (ast.Type, error) {
	return newParser(filepath, s).ParseType()
}

func ParseDDL(filepath, s string) (ast.DDL, error) {
	return newParser(filepath, s).ParseDDL()
}

func ParseDDLs(filepath, s string) ([]ast.DDL, error) {
	return newParser(filepath, s).ParseDDLs()
}

func ParseDML(filepath, s string) (ast.DML, error) {
	return newParser(filepath, s).ParseDML()
}

func ParseDMLs(filepath, s string) ([]ast.DML, error) {
	return newParser(filepath, s).ParseDMLs()
}
