package analyzer

import (
	"fmt"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Analyzer struct {
	File *parser.File

	Types       map[parser.Expr]*TypeInfo
	SelectLists map[parser.QueryExpr]SelectList

	scope          *NameScope
	aggregateScope *NameScope
}

func (a *Analyzer) AnalyzeQueryStatement(q *parser.QueryStatement) {
	// TODO: error handle
	a.analyzeQueryStatement(q)
}

func (a *Analyzer) analyzeType(t parser.Type) Type {
	switch t := t.(type) {
	case *parser.SimpleType:
		return SimpleType(t.Name)
	case *parser.ArrayType:
		return &ArrayType{Item: a.analyzeType(t.Item)}
	case *parser.StructType:
		fields := make([]*StructField, len(t.Fields))
		for i, f := range t.Fields {
			fields[i] = &StructField{
				Name: f.Member.Name,
				Type: a.analyzeType(f.Type),
			}
		}
		return &StructType{Fields: fields}
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) pushSelectListScope(list SelectList) {
	a.scope = list.toNameScope(a.scope)
}

func (a *Analyzer) pushTableScope(t *TableScope) {
	a.scope = t.toNameScope(a.scope)
}

func (a *Analyzer) popScope() {
	a.scope = a.scope.Next
}

func (a *Analyzer) lookupRef(p Path) (*Reference, Path) {
	if a.scope == nil {
		return nil, nil
	}
	return a.scope.LookupRef(p)
}

func (a *Analyzer) errorf(node parser.Node, msg string, params ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, params...),
		Position: a.File.Position(node.Pos(), node.End()),
	}
}

func (a *Analyzer) panicf(node parser.Node, msg string, params ...interface{}) {
	panic(a.errorf(node, msg, params...))
}

func extractIdentFromExpr(e parser.Expr) *parser.Ident {
	switch e := e.(type) {
	case *parser.Ident:
		return e
	case *parser.Path:
		return e.Idents[len(e.Idents)-1]
	case *parser.SelectorExpr:
		return e.Member
	case *parser.ParenExpr:
		return extractIdentFromExpr(e.Expr)
	}

	return nil
}

func extractNameFromExpr(e parser.Expr) string {
	id := extractIdentFromExpr(e)
	if id != nil {
		return id.Name
	}
	return ""
}
