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
	// TODO: error handle
	a.analyzeQueryStatement(q)
}

func (a *Analyzer) analyzeQueryStatement(q *parser.QueryStatement) {
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

	if s.AsStruct {
		t := a.convertNameListToType(&list)
		return a.newSingletonNameList("", t, s)
	}

	return &list
}

func (a *Analyzer) analyzeCompoundQuery(q *parser.CompoundQuery) *NameList {
	list := a.analyzeQueryExpr(q.Queries[0]).derive(q)

	for _, query := range q.Queries[1:] {
		queryList := a.analyzeQueryExpr(query)

		if len(list.Columns) != len(queryList.Columns) {
			a.panicf(query, "queries in set operation have mismatched column count")
		}

		for i, c := range list.Columns {
			if !c.merge(queryList.Columns[i]) {
				a.panicf(
					queryList.Columns[i].Node,
					"%s is incompatible with %s (column %d)",
					TypeString(c.Type),
					TypeString(queryList.Columns[i].Type),
					i+1,
				)
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
	return a.scope.List.derive(s)
}

func (a *Analyzer) analyzeStarPath(s *parser.StarPath) *NameList {
	t := a.analyzeExpr(s.Expr)
	list := a.convertTypeToNameList(s, t.Type)
	if list == nil {
		a.panicf(s, "star expansion is not supported for type %s", TypeString(t.Type))
	}
	return list.derive(s)
}

func (a *Analyzer) convertTypeToNameList(n parser.Node, t Type) *NameList {
	switch t := t.(type) {
	case *StructType:
		if len(t.Fields) == 0 {
			return nil
		}
		list := &NameList{
			Columns: make([]*ColumnName, len(t.Fields)),
		}
		for i, f := range t.Fields {
			list.Columns[i] = a.newColumnName(f.Name, f.Type, n)
		}
		return list
	}

	return nil
}

func (a *Analyzer) convertNameListToType(list *NameList) Type {
	fields := make([]*StructField, len(list.Columns))
	for i, c := range list.Columns {
		fields[i] = &StructField{
			Name: c.Path.Name,
			Type: c.Type,
		}
	}
	return &StructType{Fields: fields}
}

func (a *Analyzer) analyzeAlias(s *parser.Alias) *NameList {
	t := a.analyzeExpr(s.Expr)
	return a.newSingletonNameList(s.As.Alias.Name, t.Type, s)
}

func (a *Analyzer) analyzeExprSelectItem(s *parser.ExprSelectItem) *NameList {
	t := a.analyzeExpr(s.Expr)
	name := extractNameFromExpr(s.Expr)
	return a.newSingletonNameList(name, t.Type, s)
}

func (a *Analyzer) newColumnName(name string, t Type, n parser.Node) *ColumnName {
	path := a.newPathName(name)
	return newColumnName(path, t, n)
}

func (a *Analyzer) newSingletonNameList(name string, t Type, n parser.Node) *NameList {
	path := a.newPathName(name)
	return newSingletonNameList(path, t, n)
}

func (a *Analyzer) newPathName(name string) PathName {
	if name == "" {
		a.implicitAliasId++
		return PathName{
			ImplicitAliasID: a.implicitAliasId,
		}
	} else {
		return PathName{
			Name: name,
		}
	}

}

func extractNameFromExpr(e parser.Expr) string {
	switch e := e.(type) {
	case *parser.Ident:
		return e.Name
	case *parser.Path:
		return e.Idents[len(e.Idents)-1].Name
	case *parser.SelectorExpr:
		return e.Member.Name
	case *parser.ParenExpr:
		return extractNameFromExpr(e.Expr)
	}

	return ""
}

func (a *Analyzer) analyzeOrderBy(o *parser.OrderBy) {
	// TODO: implement
}

func (a *Analyzer) analyzeLimit(l *parser.Limit) {
	// TODO: implement
}

func (a *Analyzer) analyzeExpr(e parser.Expr) *TypeInfo {
	var t *TypeInfo
	switch e := e.(type) {
	case *parser.ParenExpr:
		t = a.analyzeParenExpr(e)
	case *parser.ScalarSubQuery:
		t = a.analyzeScalarSubQuery(e)
	case *parser.ArraySubQuery:
		t = a.analyzeArraySubQuery(e)
	case *parser.ExistsSubQuery:
		t = a.analyzeExistsSubQuery(e)
	case *parser.ArrayLiteral:
		t = a.analyzeArrayLiteral(e)
	case *parser.StructLiteral:
		t = a.analyzeStructLiteral(e)
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
		panic(fmt.Sprintf("TODO: implement: %t", e))
	}

	if a.Types == nil {
		a.Types = make(map[parser.Expr]*TypeInfo)
	}
	a.Types[e] = t
	return t
}

func (a *Analyzer) analyzeParenExpr(e *parser.ParenExpr) *TypeInfo {
	return a.analyzeExpr(e.Expr)
}

func (a *Analyzer) analyzeScalarSubQuery(e *parser.ScalarSubQuery) *TypeInfo {
	list := a.analyzeQueryExpr(e.Query)
	if len(list.Columns) != 1 {
		a.panicf(e, "scalar subquery must have just one column")
	}
	return &TypeInfo{
		Type: list.Columns[0].Type,
	}
}

func (a *Analyzer) analyzeArraySubQuery(e *parser.ArraySubQuery) *TypeInfo {
	list := a.analyzeQueryExpr(e.Query)
	if len(list.Columns) != 1 {
		a.panicf(e, "ARRAY subquery must have just one column")
	}
	return &TypeInfo{
		Type: &ArrayType{Item: list.Columns[0].Type},
	}
}

func (a *Analyzer) analyzeExistsSubQuery(e *parser.ExistsSubQuery) *TypeInfo {
	a.analyzeQueryExpr(e.Query)
	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeArrayLiteral(e *parser.ArrayLiteral) *TypeInfo {
	if e.Type == nil {
		return a.analyzeArrayLiteralWithoutType(e)
	}

	t := a.analyzeType(e.Type)
	for _, v := range e.Values {
		vt := a.analyzeExpr(v)
		if !TypeCoerce(vt.Type, t) {
			a.panicf(v, "%s cannot coerce to %s", TypeString(vt.Type), TypeString(t))
		}
	}

	return &TypeInfo{
		Type: &ArrayType{Item: t},
	}
}

func (a *Analyzer) analyzeArrayLiteralWithoutType(e *parser.ArrayLiteral) *TypeInfo {
	var t Type

	for _, v := range e.Values {
		vt := a.analyzeExpr(v)
		t1, ok := MergeType(t, vt.Type)
		if !ok {
			panic(a.errorf(e, "%s is incompatible with %s", TypeString(t), TypeString(vt.Type)))
		}
		t = t1
	}

	return &TypeInfo{
		Type: &ArrayType{Item: t},
	}
}

func (a *Analyzer) analyzeStructLiteral(e *parser.StructLiteral) *TypeInfo {
	if e.Fields == nil {
		return a.analyzeStructLiteralWithoutType(e)
	}

	if len(e.Fields) != len(e.Values) {
		a.panicf(e, "STRUCT type has %d fields, but literal has %d values", len(e.Fields), len(e.Values))
	}

	fields := make([]*StructField, len(e.Fields))
	for i, f := range e.Fields {
		var name string
		if f.Member != nil {
			name = f.Member.Name
		}
		fields[i] = &StructField{
			Name: name,
			Type: a.analyzeType(f.Type),
		}
	}
	t := &StructType{Fields: fields}

	for i, v := range e.Values {
		vt := a.analyzeExpr(v)
		if !TypeCoerce(vt.Type, fields[i].Type) {
			a.panicf(v, "%s cannot coerce to %s", TypeString(vt.Type), TypeString(fields[i].Type))
		}
	}

	return &TypeInfo{
		Type: t,
	}
}

func (a *Analyzer) analyzeStructLiteralWithoutType(e *parser.StructLiteral) *TypeInfo {
	fields := make([]*StructField, len(e.Values))
	for i, v := range e.Values {
		t := a.analyzeExpr(v)
		fields[i] = &StructField{
			Type: t.Type,
		}
	}
	return &TypeInfo{
		Type: &StructType{Fields: fields},
	}
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

	panic("unreachable")
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
