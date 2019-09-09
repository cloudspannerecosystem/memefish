package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Analyzer struct {
	File    *parser.File
	Catalog *Catalog
	Params  map[string]interface{}

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

func (a *Analyzer) analyzeIntValue(i parser.IntValue) int64 {
	switch i := i.(type) {
	case *parser.IntLiteral:
		v, err := strconv.ParseInt(i.Value, i.Base, 64)
		if err != nil {
			a.panicf(i, "error on parsing integer literal: %v", err)
		}
		return v
	case *parser.Param:
		v, ok := a.lookupParam(i.Name)
		if !ok {
			a.panicf(i, "unknown query parameter: %s", i.SQL())
		}
		iv, ok := v.(int64)
		if !ok {
			a.panicf(i, "invalid query parameter: %s", i.SQL())
		}
		return iv
	case *parser.CastIntValue:
		return a.analyzeIntValue(i.Expr)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeNumValue(n parser.NumValue) interface{} /* float64 | int64 */ {
	switch n := n.(type) {
	case *parser.IntLiteral:
		v, err := strconv.ParseInt(n.Value, n.Base, 64)
		if err != nil {
			a.panicf(n, "error on parsing integer literal: %v", err)
		}
		return v
	case *parser.FloatLiteral:
		v, err := strconv.ParseFloat(n.Value, 64)
		if err != nil {
			a.panicf(n, "error on parsing integer literal: %v", err)
		}
		return v
	case *parser.Param:
		v, ok := a.lookupParam(n.Name)
		if !ok {
			a.panicf(n, "unknown query parameter: %s", n.SQL())
		}
		iv, iok := v.(int64)
		fv, fok := v.(float64)
		if iok {
			return iv
		}
		if fok {
			return fv
		}
		a.panicf(n, "invalid query parameter: %s", n.SQL())
	case *parser.CastNumValue:
		return a.analyzeNumValue(n.Expr)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeStringValue(s parser.StringValue) string {
	switch s := s.(type) {
	case *parser.StringLiteral:
		return s.Value
	case *parser.Param:
		v, ok := a.lookupParam(s.Name)
		if !ok {
			a.panicf(s, "unknown query parameter: %s", s.SQL())
		}
		sv, ok := v.(string)
		if !ok {
			a.panicf(s, "invalid query parameter: %s", s.SQL())
		}
		return sv
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

func (a *Analyzer) lookupParam(target string) (interface{}, bool) {
	if a.Params == nil {
		return nil, false
	}
	p, ok := a.Params[strings.ToUpper(target)]
	return p, ok
}

func (a *Analyzer) lookupTable(target string) (*TableSchema, bool) {
	if a.Catalog == nil {
		return nil, false
	}
	table, ok := a.Catalog.Tables[strings.ToUpper(target)]
	return table, ok
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
