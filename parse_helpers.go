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

// ParseStatement parses an input string containing a statement.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseStatement(filepath, s string) (ast.Statement, error) {
	return newParser(filepath, s).ParseStatement()
}

// ParseStatements parses an input string containing statements.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseStatements(filepath, s string) ([]ast.Statement, error) {
	return newParser(filepath, s).ParseStatements()
}

// ParseQuery parses an input string containing a query statement.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseQuery(filepath, s string) (*ast.QueryStatement, error) {
	return newParser(filepath, s).ParseQuery()
}

// ParseExpr parses an input string containing an expression.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseExpr(filepath, s string) (ast.Expr, error) {
	return newParser(filepath, s).ParseExpr()
}

// ParseType parses an input string containing a type.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseType(filepath, s string) (ast.Type, error) {
	return newParser(filepath, s).ParseType()
}

// ParseDDL parses an input string containing a DDL statement.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseDDL(filepath, s string) (ast.DDL, error) {
	return newParser(filepath, s).ParseDDL()
}

// ParseDDLs parses an input string containing DDL statements.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseDDLs(filepath, s string) ([]ast.DDL, error) {
	return newParser(filepath, s).ParseDDLs()
}

// ParseDML parses an input string containing a DML statement.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseDML(filepath, s string) (ast.DML, error) {
	return newParser(filepath, s).ParseDML()
}

// ParseDMLs parses an input string containing DML statements.
// This function won't panic but returns error if parser becomes error state.
// filepath can be empty, it is only used in error message.
func ParseDMLs(filepath, s string) ([]ast.DML, error) {
	return newParser(filepath, s).ParseDMLs()
}
