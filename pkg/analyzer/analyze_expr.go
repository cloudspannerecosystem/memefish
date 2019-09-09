package analyzer

import (
	"fmt"
	"strconv"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type TypeInfo struct {
	Type  Type
	Name  *Name
	Value interface{}
}

func (a *Analyzer) analyzeExpr(e parser.Expr) *TypeInfo {
	var t *TypeInfo
	switch e := e.(type) {
	case *parser.CallExpr:
		t = a.analyzeCallExpr(e)
	case *parser.CountStarExpr:
		t = a.analyzeCountStarExpr(e)
	case *parser.ParenExpr:
		t = a.analyzeParenExpr(e)
	case *parser.ScalarSubQuery:
		t = a.analyzeScalarSubQuery(e)
	case *parser.ArraySubQuery:
		t = a.analyzeArraySubQuery(e)
	case *parser.ExistsSubQuery:
		t = a.analyzeExistsSubQuery(e)
	case *parser.Ident:
		t = a.analyzeIdent(e)
	case *parser.Path:
		t = a.analyzePath(e)
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
		panic(fmt.Sprintf("BUG: unreachable: %t", e))
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
	if len(list) != 1 {
		a.panicf(e, "scalar subquery must have just one column")
	}
	return &TypeInfo{
		Type: list[0].Type,
	}
}

func (a *Analyzer) analyzeArraySubQuery(e *parser.ArraySubQuery) *TypeInfo {
	list := a.analyzeQueryExpr(e.Query)
	if len(list) != 1 {
		a.panicf(e, "ARRAY subquery must have just one column")
	}
	return &TypeInfo{
		Type: &ArrayType{Item: list[0].Type},
	}
}

func (a *Analyzer) analyzeExistsSubQuery(e *parser.ExistsSubQuery) *TypeInfo {
	a.analyzeQueryExpr(e.Query)
	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeIdent(e *parser.Ident) *TypeInfo {
	name := a.lookup(e.Name)
	if name == nil {
		a.panicf(e, "unknown name: %s", e.SQL())
	}
	if name.Ambiguous {
		a.panicf(e, "ambiguous name: %s", e.SQL())
	}

	if a.scope.context != nil {
		if !a.scope.context.IsValidName(name) {
			a.panicf(e, "cannot use non-aggregate key: %s", e.SQL())
		}
	}

	return &TypeInfo{
		Type: name.Type,
		Name: name,
	}
}

func (a *Analyzer) analyzePath(e *parser.Path) *TypeInfo {
	id0 := e.Idents[0]
	name := a.lookup(id0.Name)
	if name == nil {
		a.panicf(e.Idents[0], "unknown name: %s", id0.SQL())
	}
	if name.Ambiguous {
		a.panicf(e, "ambiguous name: %s", id0.SQL())
	}

	for _, id := range e.Idents[1:] {
		child := name.LookupChild(id.Name)
		if child == nil {
			a.panicf(id, "unknown field: %s", id.SQL())
		}
		if child.Ambiguous {
			a.panicf(e, "ambiguous field: %s", id.SQL())
		}
		name = child
	}

	if a.scope.context != nil {
		if !a.scope.context.IsValidName(name) {
			a.panicf(e, "cannot use non-aggregate key: %s", e.SQL())
		}
	}

	return &TypeInfo{
		Type: name.Type,
		Name: name,
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
