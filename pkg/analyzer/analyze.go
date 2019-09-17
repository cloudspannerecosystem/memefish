package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/MakeNowJust/memefish/pkg/token"
)

type Analyzer struct {
	File    *token.File
	Catalog *Catalog
	Params  map[string]interface{}

	Types     map[ast.Expr]*TypeInfo
	Tables    map[ast.TableExpr]*TableInfo
	NameLists map[ast.QueryExpr]NameList

	scope *NameScope
}

func (a *Analyzer) AnalyzeQueryStatement(q *ast.QueryStatement) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	a.analyzeQueryStatement(q)
	return
}

func (a *Analyzer) analyzeType(t ast.Type) Type {
	switch t := t.(type) {
	case *ast.SimpleType:
		return SimpleType(t.Name)
	case *ast.ArrayType:
		return &ArrayType{Item: a.analyzeType(t.Item)}
	case *ast.StructType:
		fields := make([]*StructField, len(t.Fields))
		for i, f := range t.Fields {
			var name string
			if f.Ident != nil {
				name = f.Ident.Name
			}
			fields[i] = &StructField{
				Name: name,
				Type: a.analyzeType(f.Type),
			}
		}
		return &StructType{Fields: fields}
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeIntValue(i ast.IntValue) int64 {
	switch i := i.(type) {
	case *ast.IntLiteral:
		v, err := strconv.ParseInt(i.Value, i.Base, 64)
		if err != nil {
			a.panicf(i, "error on parsing integer literal: %v", err)
		}
		return v
	case *ast.Param:
		v, ok := a.lookupParam(i.Name)
		if !ok {
			a.panicf(i, "unknown query parameter: %s", i.SQL())
		}
		iv, ok := v.(int64)
		if !ok {
			a.panicf(i, "invalid query parameter: %s", i.SQL())
		}
		return iv
	case *ast.CastIntValue:
		return a.analyzeIntValue(i.Expr)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeNumValue(n ast.NumValue) interface{} /* float64 | int64 */ { //nolint:unused
	switch n := n.(type) {
	case *ast.IntLiteral:
		v, err := strconv.ParseInt(n.Value, n.Base, 64)
		if err != nil {
			a.panicf(n, "error on parsing integer literal: %v", err)
		}
		return v
	case *ast.FloatLiteral:
		v, err := strconv.ParseFloat(n.Value, 64)
		if err != nil {
			a.panicf(n, "error on parsing integer literal: %v", err)
		}
		return v
	case *ast.Param:
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
	case *ast.CastNumValue:
		return a.analyzeNumValue(n.Expr)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeStringValue(s ast.StringValue) string {
	switch s := s.(type) {
	case *ast.StringLiteral:
		return s.Value
	case *ast.Param:
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

func (a *Analyzer) errorf(node ast.Node, msg string, params ...interface{}) *Error {
	var position *token.Position
	if node != nil {
		position = a.File.Position(node.Pos(), node.End())
	}

	return &Error{
		Message:  fmt.Sprintf(msg, params...),
		Position: position,
	}
}

func (a *Analyzer) panicf(node ast.Node, msg string, params ...interface{}) {
	panic(a.errorf(node, msg, params...))
}

func extractIdentFromExpr(e ast.Expr) *ast.Ident {
	switch e := e.(type) {
	case *ast.Ident:
		return e
	case *ast.Path:
		return e.Idents[len(e.Idents)-1]
	case *ast.SelectorExpr:
		return e.Ident
	case *ast.ParenExpr:
		return extractIdentFromExpr(e.Expr)
	}

	return nil
}
