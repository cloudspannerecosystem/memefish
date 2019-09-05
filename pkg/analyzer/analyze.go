package analyzer

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Analyzer struct {
	File *parser.File

	Types               map[parser.Expr]*TypeInfo
	NameLists           map[parser.QueryExpr]*NameList
	SelectItemNameLists map[parser.SelectItem]*NameList

	scope           *NameScope
	aggregateScope  *NameScope
	implicitAliasId int
}

type TypeInfo struct {
	Type           Type
	ResolvedTable  *TableName
	ResolvedColumn *ColumnName
	Scope          *NameScope
	Value          interface{}
}

func (a *Analyzer) AnalyzeQueryStatement(q *parser.QueryStatement) {
	// TODO: analyze q.Hint
	_ = a.analyzeQueryExpr(q.Query)
}

func (a *Analyzer) analyzeQueryExpr(q parser.QueryExpr) *NameList {
	var list *NameList
	switch q := q.(type) {
	case *parser.Select:
		list = a.analyzeSelect(q)
	case *parser.CompoundQuery:
		list = a.analyzeCompoundQuery(q)
	case *parser.SubQuery:
		list = a.analyzeSubQuery(q)
	}
	if a.NameLists == nil {
		a.NameLists = make(map[parser.QueryExpr]*NameList)
	}
	a.NameLists[q] = list
	return list
}

func (a *Analyzer) analyzeSelect(s *parser.Select) *NameList {
	switch {
	case s.From == nil:
		return a.analyzeSelectWithoutFrom(s)
	}

	panic("TODO")
}

func (a *Analyzer) analyzeSelectWithoutFrom(s *parser.Select) *NameList {
	if s.Where != nil {
		a.panicf(s.Where, "SELECT without FROM cannot have WHERE clause")
	}
	if s.GroupBy != nil {
		a.panicf(s.GroupBy, "SELECT without FROM cannot have GROUP BY clause")
	}
	if s.Having != nil {
		a.panicf(s.Having, "SELECT without FROM cannot have HAVING clause")
	}
	if s.OrderBy != nil {
		a.panicf(s.OrderBy, "SELECT without FROM cannot have ORDER BY clause")
	}

	var list NameList
	for _, item := range s.Results {
		itemList := a.analyzeSelectItem(item)
		list.concat(itemList)
	}

	return &list
}

func (a *Analyzer) analyzeCompoundQuery(q *parser.CompoundQuery) *NameList {
	list := a.analyzeQueryExpr(q.Queries[0])

	for _, query := range q.Queries[1:] {
		queryList := a.analyzeQueryExpr(query)

		if len(list.Columns) != len(queryList.Columns) {
			a.panicf(query, "queries in set operation have mismatched column count")
		}

		for i, c := range list.Columns {
			if c.merge(queryList.Columns[i]) {
				nodes := queryList.Columns[i].Nodes
				a.panicf(nodes[len(nodes)-1], "column %d has incompatible type", i+1)
			}
		}
	}

	a.pushNameListScope(list)
	a.analyzeOrderBy(q.OrderBy)
	a.analyzeLimit(q.Limit)
	a.popScope()

	return list
}

func (a *Analyzer) analyzeSubQuery(s *parser.SubQuery) *NameList {
	panic("TODO: implement")
}

func (a *Analyzer) analyzeSelectItem(s parser.SelectItem) *NameList {
	switch s := s.(type) {
	case *parser.Star:
		return a.analyzeStar(s)
	case *parser.StarPath:
		return a.analyzeStarPath(s)
	case *parser.Alias:
		return a.analyzeAlias(s)
	case *parser.ExprSelectItem:
		return a.analyzeExprSelectItem(s)
	}

	panic("unreachable")
}

func (a *Analyzer) analyzeStar(s *parser.Star) *NameList {
	if a.scope != nil || a.scope.List != nil {
		a.panicf(s, "SELECT * must have a FROM clause")
	}
	return a.scope.List.appendNodeToColumns(s)
}

func (a *Analyzer) analyzeStarPath(s *parser.StarPath) *NameList {
	t := a.AnalyzeExpr(s.Expr)
	list := convertTypeToNameList(t)
	if list != nil {
		a.panicf(s, "star expansion is not supported for type %s", t.String())
	}
	return list.appendNodeToColumns(s)
}

func convertTypeToNameList(t Type) *NameList {
	switch t := t.(type) {
	case *StructType:
		if len(t.Fields) == 0 {
			return nil
		}
		list := &NameList{
			Columns: make([]*ColumnName, len(t.Fields)),
		}
		for i, f := range t.Fields {
			list.Columns[i] = &ColumnName{
				Path: PathName{Name: f.Name},
				Type: f.Type,
			}
		}
		return list
	}

	return nil
}

func (a *Analyzer) analyzeAlias(s *parser.Alias) *NameList {
	t := a.AnalyzeExpr(s.Expr)
	return a.aliasColumn(s, s.As.Alias.Name, t)
}

func (a *Analyzer) analyzeExprSelectItem(s *parser.ExprSelectItem) *NameList {
	t := a.AnalyzeExpr(s.Expr)
	name := extractNameFromExpr(s.Expr)
	return a.aliasColumn(s, name, t)
}

func extractNameFromExpr(e parser.Expr) string {
	switch e := e.(type) {
	case *parser.Ident:
		return e.Name
	case *parser.Path:
		return e.Idents[len(e.Idents)-1].Name
	case *parser.SelectorExpr:
		return e.Member.Name
	}

	return ""
}

func (a *Analyzer) analyzeOrderBy(o *parser.OrderBy) {
	// TODO: implement
}

func (a *Analyzer) analyzeLimit(l *parser.Limit) {
	// TODO: implement
}

func (a *Analyzer) AnalyzeExpr(e parser.Expr) Type {
	var t *TypeInfo
	switch e := e.(type) {
	case *parser.NullLiteral:
		t = a.analyzeNullLiteral(e)
	case *parser.BoolLiteral:
		t = a.analyzeBoolLiteral(e)
	case *parser.IntLiteral:
		t = a.analyzeIntLiteral(e)
	case *parser.FloatLiteral:
		t = a.analyzeFloatLiteral(e)
	case *parser.StringLiteral:
		t = a.analyzeStringLiteral(e)
	case *parser.BytesLiteral:
		t = a.analyzeBytesLiteral(e)
	case *parser.DateLiteral:
		t = a.analyzeDateLiteral(e)
	case *parser.TimestampLiteral:
		t = a.analyzeTimestampLiteral(e)
	default:
		panic("TODO: implement")
	}

	if a.Types == nil {
		a.Types = make(map[parser.Expr]*TypeInfo)
	}
	a.Types[e] = t
	return t.Type
}

func (a *Analyzer) analyzeNullLiteral(e *parser.NullLiteral) *TypeInfo {
	return &TypeInfo{}
}

func (a *Analyzer) analyzeBoolLiteral(e *parser.BoolLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  BoolType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeIntLiteral(e *parser.IntLiteral) *TypeInfo {
	v, err := strconv.ParseInt(e.Value, e.Base, 64)
	if err != nil {
		a.panicf(e, "error on parsing integer literal: %v", err)
	}
	return &TypeInfo{
		Type:  Int64Type,
		Value: v,
	}
}

func (a *Analyzer) analyzeFloatLiteral(e *parser.FloatLiteral) *TypeInfo {
	v, err := strconv.ParseFloat(e.Value, 64)
	if err != nil {
		a.panicf(e, "error on pasing floating point number literal: %v", err)
	}
	return &TypeInfo{
		Type:  Float64Type,
		Value: v,
	}
}

func (a *Analyzer) analyzeStringLiteral(e *parser.StringLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  StringType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeBytesLiteral(e *parser.BytesLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  BytesType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeDateLiteral(e *parser.DateLiteral) *TypeInfo {
	// TODO: check e.Value format
	return &TypeInfo{
		Type: DateType,
	}
}

func (a *Analyzer) analyzeTimestampLiteral(e *parser.TimestampLiteral) *TypeInfo {
	// TODO: check e.Value format
	return &TypeInfo{
		Type: TimestampType,
	}
}

func (a *Analyzer) aliasColumn(s parser.SelectItem, name string, t Type) *NameList {
	var pathName PathName
	if name == "" {
		a.implicitAliasId++
		pathName = PathName{
			ImplicitAliasID: a.implicitAliasId,
		}
	} else {
		pathName = PathName{
			Name: name,
		}
	}
	return &NameList{
		Columns: []*ColumnName{
			&ColumnName{
				Path:  pathName,
				Type:  t,
				Nodes: []parser.SelectItem{s},
			},
		},
	}
}

func (a *Analyzer) pushNameListScope(list *NameList) {
	// TODO: implement
}

func (a *Analyzer) popScope() {
	// TODO: implement
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
