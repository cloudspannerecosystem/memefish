package analyzer

import (
	"fmt"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Analyzer struct {
	File    *parser.File
	Catalog *Catalog

	Types     map[parser.Expr]*TypeInfo
	Tables    map[parser.TableExpr]*TableInfo
	NameLists map[parser.QueryExpr]NameList

	scope *NameScope
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

func (a *Analyzer) pushNameList(list NameList) {
	a.scope = list.toNameScope(a.scope)
}

func (a *Analyzer) pushTableInfo(ti *TableInfo) {
	a.scope = ti.toNameScope(a.scope)
}

func (a *Analyzer) popScope() {
	a.scope = a.scope.Next
}

func (a *Analyzer) lookup(target string) (*Name, *GroupByContext) {
	if a.scope == nil {
		return nil, nil
	}
	return a.scope.Lookup(target)
}

func (a *Analyzer) errorf(node parser.Node, msg string, params ...interface{}) *Error {
	var position *parser.Position
	if node != nil {
		position = a.File.Position(node.Pos(), node.End())
	}

	return &Error{
		Message:  fmt.Sprintf(msg, params...),
		Position: position,
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
