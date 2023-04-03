package analyzer

import (
	"fmt"
	"strconv"

	"github.com/cloudspannerecosystem/memefish/pkg/ast"
	"github.com/cloudspannerecosystem/memefish/pkg/char"
)

type TypeInfo struct {
	Type  Type
	Name  *Name
	Value interface{}
}

func (a *Analyzer) analyzeExpr(e ast.Expr) *TypeInfo {
	var t *TypeInfo
	switch e := e.(type) {
	case *ast.BinaryExpr:
		t = a.analyzeBinaryExpr(e)
	case *ast.UnaryExpr:
		t = a.analyzeUnaryExpr(e)
	case *ast.InExpr:
		t = a.analyzeInExpr(e)
	case *ast.IsNullExpr:
		t = a.analyzeIsNullExpr(e)
	case *ast.IsBoolExpr:
		t = a.analyzeIsBoolExpr(e)
	case *ast.BetweenExpr:
		t = a.analyzeBetweenExpr(e)
	case *ast.SelectorExpr:
		t = a.analyzeSelectorExpr(e)
	case *ast.IndexExpr:
		t = a.analyzeIndexExpr(e)
	case *ast.CallExpr:
		t = a.analyzeCallExpr(e)
	case *ast.CountStarExpr:
		t = a.analyzeCountStarExpr(e)
	case *ast.CastExpr:
		t = a.analyzeCastExpr(e)
	case *ast.ExtractExpr:
		t = a.analyzeExtractExpr(e)
	case *ast.CaseExpr:
		t = a.analyzeCaseExpr(e)
	case *ast.ParenExpr:
		t = a.analyzeParenExpr(e)
	case *ast.ScalarSubQuery:
		t = a.analyzeScalarSubQuery(e)
	case *ast.ArraySubQuery:
		t = a.analyzeArraySubQuery(e)
	case *ast.ExistsSubQuery:
		t = a.analyzeExistsSubQuery(e)
	case *ast.Ident:
		t = a.analyzeIdent(e)
	case *ast.Path:
		t = a.analyzePath(e)
	case *ast.Param:
		t = a.analyzeParam(e)
	case *ast.ArrayLiteral:
		t = a.analyzeArrayLiteral(e)
	case *ast.StructLiteral:
		t = a.analyzeStructLiteral(e)
	case *ast.NullLiteral:
		t = a.analyzeNullLiteral(e)
	case *ast.BoolLiteral:
		t = a.analyzeBoolLiteral(e)
	case *ast.IntLiteral:
		t = a.analyzeIntLiteral(e)
	case *ast.FloatLiteral:
		t = a.analyzeFloatLiteral(e)
	case *ast.StringLiteral:
		t = a.analyzeStringLiteral(e)
	case *ast.BytesLiteral:
		t = a.analyzeBytesLiteral(e)
	case *ast.DateLiteral:
		t = a.analyzeDateLiteral(e)
	case *ast.TimestampLiteral:
		t = a.analyzeTimestampLiteral(e)
	case *ast.NumericLiteral:
		t = a.analyzeNumericLiteral(e)
	default:
		panic(fmt.Sprintf("BUG: unreachable: %t", e))
	}

	if a.Types == nil {
		a.Types = make(map[ast.Expr]*TypeInfo)
	}
	a.Types[e] = t
	return t
}

func (a *Analyzer) analyzeBinaryExpr(e *ast.BinaryExpr) *TypeInfo {
	lt := a.analyzeExpr(e.Left)
	rt := a.analyzeExpr(e.Right)

	switch e.Op {
	case ast.OpAnd, ast.OpOr:
		if TypeCoerce(lt.Type, BoolType) && TypeCoerce(rt.Type, BoolType) {
			return &TypeInfo{
				Type: BoolType,
			}
		}
		a.panicf(e, "operator %s requires two BOOL, but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	case ast.OpEqual, ast.OpNotEqual, ast.OpLess, ast.OpGreater, ast.OpLessEqual, ast.OpGreaterEqual:
		if TypeCoerce(lt.Type, rt.Type) || TypeCoerce(rt.Type, lt.Type) {
			return &TypeInfo{
				Type: BoolType,
			}
		}
		a.panicf(e, "operator %s requires two compatible types, but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	case ast.OpLike, ast.OpNotLike:
		if TypeCoerce(lt.Type, StringType) && TypeCoerce(rt.Type, StringType) {
			return &TypeInfo{
				Type: BoolType,
			}
		}
		if TypeCoerce(lt.Type, BytesType) && TypeCoerce(rt.Type, BytesType) {
			return &TypeInfo{
				Type: BoolType,
			}
		}
		a.panicf(e, "operator %s requires two STRING/BYTES, but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	case ast.OpBitAnd, ast.OpBitXor, ast.OpBitOr:
		if TypeCoerce(lt.Type, Int64Type) && TypeCoerce(rt.Type, Int64Type) {
			return &TypeInfo{
				Type: Int64Type,
			}
		}
		if TypeCoerce(lt.Type, BytesType) && TypeCoerce(rt.Type, BytesType) {
			return &TypeInfo{
				Type: BytesType,
			}
		}
		a.panicf(e, "operator %s requires two INT64/BYTES, but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	case ast.OpBitLeftShift, ast.OpBitRightShift:
		if TypeCoerce(lt.Type, Int64Type) && TypeCoerce(rt.Type, Int64Type) {
			return &TypeInfo{
				Type: Int64Type,
			}
		}
		if TypeCoerce(lt.Type, BytesType) && TypeCoerce(rt.Type, Int64Type) {
			return &TypeInfo{
				Type: BytesType,
			}
		}
		a.panicf(e, "operator %s requires (INT64, INT64) or (BYTES, INT64), but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	case ast.OpAdd, ast.OpSub, ast.OpMul, ast.OpDiv:
		if TypeCoerce(lt.Type, Int64Type) && TypeCoerce(rt.Type, Int64Type) {
			return &TypeInfo{
				Type: Int64Type,
			}
		}
		if TypeCoerce(lt.Type, Float64Type) && TypeCoerce(rt.Type, Float64Type) {
			return &TypeInfo{
				Type: Float64Type,
			}
		}
		a.panicf(e, "operator %s requires two INT64/FLOAT64, but: %s, %s", e.Op, TypeString(lt.Type), TypeString(rt.Type))
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeUnaryExpr(e *ast.UnaryExpr) *TypeInfo {
	t := a.analyzeExpr(e.Expr)

	switch e.Op {
	case ast.OpBitNot:
		if TypeCoerce(t.Type, Int64Type) {
			return &TypeInfo{
				Type: Int64Type,
			}
		}
		if TypeCoerce(t.Type, BytesType) {
			return &TypeInfo{
				Type: BytesType,
			}
		}
		a.panicf(e, "operator %s requires INT64/BYTES, but: %s", e.Op, TypeString(t.Type))
	case ast.OpPlus, ast.OpMinus:
		if TypeCoerce(t.Type, Int64Type) {
			return &TypeInfo{
				Type: Int64Type,
			}
		}
		if TypeCoerce(t.Type, Float64Type) {
			return &TypeInfo{
				Type: Float64Type,
			}
		}
		a.panicf(e, "operator %s requires INT64/FLOAT64, but: %s", e.Op, TypeString(t.Type))
	case ast.OpNot:
		if TypeCoerce(t.Type, BoolType) {
			return &TypeInfo{
				Type: BoolType,
			}
		}
		a.panicf(e, "operator NOT requires BOOL, but: %s", TypeString(t.Type))
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeInExpr(e *ast.InExpr) *TypeInfo {
	lt := a.analyzeExpr(e.Left)
	rt := a.analyzeInCondition(e.Right)

	if !(TypeCoerce(lt.Type, rt.Type) || TypeCoerce(rt.Type, lt.Type)) {
		a.panicf(e, "operator IN requires incompatible type, but: %s, %s", TypeString(lt.Type), TypeString(rt.Type))
	}

	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeInCondition(cond ast.InCondition) *TypeInfo {
	switch c := cond.(type) {
	case *ast.UnnestInCondition:
		return a.analyzeUnnestInCondition(c)
	case *ast.SubQueryInCondition:
		return a.analyzeSubQueryInCondition(c)
	case *ast.ValuesInCondition:
		return a.analyzeValuesInCondition(c)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeUnnestInCondition(cond *ast.UnnestInCondition) *TypeInfo {
	t := a.analyzeExpr(cond.Expr)
	tt, ok := TypeCastArray(t.Type)
	if !ok {
		a.panicf(cond, "UNNEST value must be ARRAY, but: %s", TypeString(t.Type))
	}

	return &TypeInfo{
		Type: tt.Item,
	}
}

func (a *Analyzer) analyzeSubQueryInCondition(cond *ast.SubQueryInCondition) *TypeInfo {
	list := a.analyzeQueryExpr(cond.Query)
	if len(list) != 1 {
		a.panicf(cond, "IN condition subquery must have just one column")
	}
	return &TypeInfo{
		Type: list[0].Type,
	}
}

func (a *Analyzer) analyzeValuesInCondition(cond *ast.ValuesInCondition) *TypeInfo {
	var t Type

	for _, e := range cond.Exprs {
		vt := a.analyzeExpr(e)
		t1, ok := MergeType(t, vt.Type)
		if !ok {
			panic(a.errorf(e, "%s is incompatible with %s", TypeString(t), TypeString(vt.Type)))
		}
		t = t1
	}

	return &TypeInfo{
		Type: t,
	}
}

func (a *Analyzer) analyzeIsNullExpr(e *ast.IsNullExpr) *TypeInfo {
	a.analyzeExpr(e.Left)
	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeIsBoolExpr(e *ast.IsBoolExpr) *TypeInfo {
	t := a.analyzeExpr(e.Left)
	if !TypeCoerce(t.Type, BoolType) {
		a.panicf(e, "operator IS TRUE/FALSE needs BOOL, but: %s", TypeString(t.Type))
	}
	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeBetweenExpr(e *ast.BetweenExpr) *TypeInfo {
	lt := a.analyzeExpr(e.Left)
	rst := a.analyzeExpr(e.RightStart)
	ret := a.analyzeExpr(e.RightEnd)

	if !(TypeCoerce(lt.Type, rst.Type) || TypeCoerce(rst.Type, lt.Type)) {
		a.panicf(e, "operator BETWEEN require two compatible types, but: %s, %s", TypeString(lt.Type), TypeString(rst.Type))
	}
	if !(TypeCoerce(lt.Type, ret.Type) || TypeCoerce(ret.Type, lt.Type)) {
		a.panicf(e, "operator BETWEEN require two compatible types, but: %s, %s", TypeString(lt.Type), TypeString(ret.Type))
	}

	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeCastExpr(e *ast.CastExpr) *TypeInfo {
	t1 := a.analyzeExpr(e.Expr)
	t2 := a.analyzeType(e.Type)
	if !TypeCast(t1.Type, t2) {
		a.panicf(e, "%s cannot cast to %s", TypeString(t1.Type), TypeString(t2))
	}

	return &TypeInfo{
		Type: t2,
	}
}

func (a *Analyzer) analyzeExtractExpr(e *ast.ExtractExpr) *TypeInfo {
	t := a.analyzeExpr(e.Expr)
	var resultType Type
	allowDate := false
	switch char.ToUpper(e.Part.Name) {
	case "DAYOFWEEK", "DAY", "DAYOFYEAR", "WEEK", "ISOWEEK", "MONTH", "QUARTER", "YEAR":
		allowDate = true
		resultType = Int64Type
	case "NANOSECOND", "MICROSECOND", "MILLISECOND", "SECOND", "MINUTE", "HOUR":
		resultType = Int64Type
	case "DATE":
		resultType = DateType
	default:
		a.panicf(e.Part, "unknown EXTRACT part name: %s", e.Part.SQL())
	}

	if !(TypeCoerce(t.Type, TimestampType) || allowDate && TypeCoerce(t.Type, DateType)) {
		allow := "TIMESTAMP"
		if allowDate {
			allow += "/DATE"
		}
		a.panicf(e.Part, "EXTRACT(%s FROM ...) requires %s, but: %s", char.ToUpper(e.Part.Name), allow, TypeString(t.Type))
	}

	// TODO: check e.AtTimeZone

	return &TypeInfo{
		Type: resultType,
	}
}

func (a *Analyzer) analyzeCaseExpr(e *ast.CaseExpr) *TypeInfo {
	var exprType Type
	if e.Expr == nil {
		exprType = BoolType
	} else {
		exprType = a.analyzeExpr(e.Expr).Type
	}

	var t Type

	for _, w := range e.Whens {
		ct := a.analyzeExpr(w.Cond)
		if !(TypeCoerce(ct.Type, exprType) || TypeCoerce(exprType, ct.Type)) {
			a.panicf(w.Cond, "WHEN clause condition requires %s, but: %s", TypeString(exprType), TypeString(ct.Type))
		}
		tt := a.analyzeExpr(w.Then)
		t1, ok := MergeType(t, tt.Type)
		if !ok {
			a.panicf(w.Then, "%s is incompatible with %s", TypeString(t), TypeString(tt.Type))
		}
		t = t1
	}

	if e.Else != nil {
		tt := a.analyzeExpr(e.Else.Expr)
		t1, ok := MergeType(t, tt.Type)
		if !ok {
			a.panicf(e.Else.Expr, "%s is incompatible with %s", TypeString(t), TypeString(tt.Type))
		}
		t = t1
	}

	return &TypeInfo{
		Type: t,
	}
}

func (a *Analyzer) analyzeSelectorExpr(e *ast.SelectorExpr) *TypeInfo {
	t := a.analyzeExpr(e.Expr)
	var names NameList
	if t.Name != nil {
		names = t.Name.Children()
	} else {
		names = makeNameListFromType(t.Type, e.Expr)
	}
	child := names.Lookup(e.Ident.Name)
	if child == nil {
		a.panicf(e.Ident, "unknown field: %s", e.Ident.SQL())
	}
	return &TypeInfo{
		Type: child.Type,
		Name: child,
	}
}

func (a *Analyzer) analyzeIndexExpr(e *ast.IndexExpr) *TypeInfo {
	et := a.analyzeExpr(e.Expr)
	it := a.analyzeExpr(e.Index)

	ett, ok := TypeCastArray(et.Type)
	if !ok {
		a.panicf(e.Expr, "element access using [] is not supported values of %s", TypeString(et.Type))
	}

	if !TypeCoerce(it.Type, Int64Type) {
		a.panicf(e.Expr, "array position in [] must be INT64, but: %s", TypeString(it.Type))
	}

	return &TypeInfo{
		Type: ett.Item,
	}
}

func (a *Analyzer) analyzeParenExpr(e *ast.ParenExpr) *TypeInfo {
	return a.analyzeExpr(e.Expr)
}

func (a *Analyzer) analyzeScalarSubQuery(e *ast.ScalarSubQuery) *TypeInfo {
	list := a.analyzeQueryExpr(e.Query)
	if len(list) != 1 {
		a.panicf(e, "scalar subquery must have just one column")
	}
	return &TypeInfo{
		Type: list[0].Type,
	}
}

func (a *Analyzer) analyzeArraySubQuery(e *ast.ArraySubQuery) *TypeInfo {
	list := a.analyzeQueryExpr(e.Query)
	if len(list) != 1 {
		a.panicf(e, "ARRAY subquery must have just one column")
	}
	return &TypeInfo{
		Type: &ArrayType{Item: list[0].Type},
	}
}

func (a *Analyzer) analyzeExistsSubQuery(e *ast.ExistsSubQuery) *TypeInfo {
	a.analyzeQueryExpr(e.Query)
	return &TypeInfo{
		Type: BoolType,
	}
}

func (a *Analyzer) analyzeIdent(e *ast.Ident) *TypeInfo {
	name, context := a.lookup(e.Name)
	if name == nil {
		a.panicf(e, "unknown name: %s", e.SQL())
	}
	if name.Ambiguous {
		a.panicf(e, "ambiguous name: %s", e.SQL())
	}
	if context != nil && !context.IsValidName(name) {
		a.panicf(e, "cannot use non-aggregate key: %s", e.SQL())
	}

	return &TypeInfo{
		Type: name.Type,
		Name: name,
	}
}

func (a *Analyzer) analyzePath(e *ast.Path) *TypeInfo {
	id0 := e.Idents[0]
	name, context := a.lookup(id0.Name)
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

	if context != nil && !context.IsValidName(name) {
		a.panicf(e, "cannot use non-aggregate key: %s", e.SQL())
	}

	return &TypeInfo{
		Type: name.Type,
		Name: name,
	}
}

func (a *Analyzer) analyzeParam(e *ast.Param) *TypeInfo {
	v, ok := a.lookupParam(e.Name)
	if !ok {
		a.panicf(e, "unknown query parameter: %s", e.SQL())
	}
	t, err := valueType(v)
	if err != nil {
		a.panicf(e, "invalid query parameter value: %v", err)
	}
	return &TypeInfo{
		Type:  t,
		Value: v,
	}
}

func (a *Analyzer) analyzeArrayLiteral(e *ast.ArrayLiteral) *TypeInfo {
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

func (a *Analyzer) analyzeArrayLiteralWithoutType(e *ast.ArrayLiteral) *TypeInfo {
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

func (a *Analyzer) analyzeStructLiteral(e *ast.StructLiteral) *TypeInfo {
	if e.Fields == nil {
		return a.analyzeStructLiteralWithoutType(e)
	}

	if len(e.Fields) != len(e.Values) {
		a.panicf(e, "STRUCT type has %d fields, but literal has %d values", len(e.Fields), len(e.Values))
	}

	fields := make([]*StructField, len(e.Fields))
	for i, f := range e.Fields {
		var name string
		if f.Ident != nil {
			name = f.Ident.Name
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

func (a *Analyzer) analyzeStructLiteralWithoutType(e *ast.StructLiteral) *TypeInfo {
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

func (a *Analyzer) analyzeNullLiteral(e *ast.NullLiteral) *TypeInfo {
	return &TypeInfo{}
}

func (a *Analyzer) analyzeBoolLiteral(e *ast.BoolLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  BoolType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeIntLiteral(e *ast.IntLiteral) *TypeInfo {
	v, err := strconv.ParseInt(e.Value, e.Base, 64)
	if err != nil {
		a.panicf(e, "error on parsing integer literal: %v", err)
	}
	return &TypeInfo{
		Type:  Int64Type,
		Value: v,
	}
}

func (a *Analyzer) analyzeFloatLiteral(e *ast.FloatLiteral) *TypeInfo {
	v, err := strconv.ParseFloat(e.Value, 64)
	if err != nil {
		a.panicf(e, "error on pasing floating point number literal: %v", err)
	}
	return &TypeInfo{
		Type:  Float64Type,
		Value: v,
	}
}

func (a *Analyzer) analyzeStringLiteral(e *ast.StringLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  StringType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeBytesLiteral(e *ast.BytesLiteral) *TypeInfo {
	return &TypeInfo{
		Type:  BytesType,
		Value: e.Value,
	}
}

func (a *Analyzer) analyzeDateLiteral(e *ast.DateLiteral) *TypeInfo {
	// TODO: check e.Value format
	return &TypeInfo{
		Type: DateType,
	}
}

func (a *Analyzer) analyzeTimestampLiteral(e *ast.TimestampLiteral) *TypeInfo {
	// TODO: check e.Value format
	return &TypeInfo{
		Type: TimestampType,
	}
}

func (a *Analyzer) analyzeNumericLiteral(e *ast.NumericLiteral) *TypeInfo {
	// TODO: check e.Value format
	return &TypeInfo{
		Type: NumericType,
	}
}
