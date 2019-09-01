package parser

import (
	"fmt"
	"strconv"
	"strings"
)

type prec int

const (
	precLit prec = iota
	precSelector
	precUnary
	precMulDiv
	precAddSub
	precBitShift
	precBitAnd
	precBitXor
	precBitOr
	precComparison
	precNot
	precAndOr
)

func exprPrec(e Expr) prec {
	switch e := e.(type) {
	case *CallExpr, *CastExpr, *CaseExpr, *SubQuery, *ParenExpr, *Param, *Ident, *ArrayLit, *StructLit, *NullLit, *BoolLit, *IntLit, *FloatLit, *StringLit, *BytesLit, *DateLit, *TimestampLit:
		return precLit
	case *IndexExpr, *SelectorExpr:
		return precSelector
	case *InExpr, *IsNullExpr, *IsBoolExpr, *BetweenExpr:
		return precComparison
	case *BinaryExpr:
		switch e.Op {
		case OpMul, OpDiv:
			return precMulDiv
		case OpAdd, OpSub:
			return precAddSub
		case OpBitLeftShift, OpBitRightShift:
			return precBitShift
		case OpBitAnd:
			return precBitAnd
		case OpBitXor:
			return precBitXor
		case OpBitOr:
			return precBitOr
		case OpEqual, OpNotEqual, OpLess, OpLessEqual, OpGreater, OpGreaterEqual, OpLike, OpNotLike:
			return precComparison
		case OpAnd, OpOr:
			return precAndOr
		}
	case *UnaryExpr:
		switch e.Op {
		case OpPlus, OpMinus, OpBitNot:
			return precUnary
		case OpNot:
			return precNot
		}
	}

	panic("exprPrec: unexpected")
}

func paren(p prec, e Expr) string {
	ep := exprPrec(e)
	if ep <= p {
		return e.SQL()
	} else {
		return "(" + e.SQL() + ")"
	}
}

func (e ExprList) SQL() string {
	var ss []string
	for _, i := range e {
		ss = append(ss, i.SQL())
	}
	return strings.Join(ss, ", ")
}

func (q *QueryStatement) SQL() string {
	if q.Hint == nil {
		return q.Expr.SQL()
	}

	return fmt.Sprintf("%s %s", q.Hint.SQL(), q.Expr.SQL())
}

func (h *Hint) SQL() string {
	var ss []string
	for name, expr := range h.Map {
		// TODO: escape name as same as identifier
		ss = append(ss, fmt.Sprintf("%s = %s", name, expr.SQL()))
	}
	return fmt.Sprintf("@{%s}", strings.Join(ss, ", "))
}

func (s *Select) SQL() string {
	sql := "SELECT "
	if s.Distinct {
		sql += "DISTINCT "
	}
	if s.AsStruct {
		sql += "AS STRUCT "
	}

	sql += s.List.SQL()
	if s.From != nil {
		sql += " FROM " + s.From.SQL()
	}
	if s.Where != nil {
		sql += " WHERE " + s.Where.SQL()
	}
	if s.GroupBy != nil {
		sql += " GROUP BY " + s.GroupBy.SQL()
	}
	if s.Having != nil {
		sql += " HAVING " + s.Having.SQL()
	}
	if s.OrderBy != nil {
		sql += " ORDER BY " + s.OrderBy.SQL()
	}
	if s.Limit != nil {
		sql += " LIMIT " + s.Limit.SQL()
	}

	return sql
}

func (c *CompoundQuery) SQL() string {
	op := string(c.Op)
	if c.Distinct {
		op += " DISTINCT"
	} else {
		op += " ALL"
	}
	sql := ""
	for i, e := range c.List {
		if i != 0 {
			sql += " " + op + " "
		}
		if s, ok := e.(*SubQueryExpr); ok && s.OrderBy == nil && s.Limit == nil {
			sql += s.SQL()
		} else if s, ok := e.(*Select); ok && s.OrderBy == nil && s.Limit == nil {
			sql += s.SQL()
		} else {
			sql += "(" + e.SQL() + ")"
		}
	}
	if c.OrderBy != nil {
		sql += " ORDER BY " + c.OrderBy.SQL()
	}
	if c.Limit != nil {
		sql += " LIMIT " + c.Limit.SQL()
	}
	return sql
}

func (s *SubQueryExpr) SQL() string {
	sql := s.Expr.SQL()
	if s.OrderBy != nil {
		sql += " ORDER BY " + s.OrderBy.SQL()
	}
	if s.Limit != nil {
		sql += " LIMIT " + s.Limit.SQL()
	}
	return sql
}

func (s *SelectExpr) SQL() string {
	if s.Expr == nil {
		return "*"
	}

	var sql string
	if s.Star {
		sql = paren(precSelector, s.Expr)
		sql += ".*"
	} else {
		sql = s.Expr.SQL()
	}

	if s.As != nil {
		sql += " AS " + s.As.SQL()
	}

	return sql
}

func (s SelectExprList) SQL() string {
	var ss []string
	for _, i := range s {
		ss = append(ss, i.SQL())
	}
	return strings.Join(ss, ", ")
}

func (f *FromItem) SQL() string {
	sql := f.Expr.SQL()
	if f.TableSample != "" {
		sql += " TABLESAMPLE " + string(f.TableSample)
	}
	return sql
}

func (f FromItemList) SQL() string {
	var ss []string
	for _, i := range f {
		ss = append(ss, i.SQL())
	}
	return strings.Join(ss, ", ")
}

func (t *TableName) SQL() string {
	sql := t.Ident.SQL()
	if t.Hint != nil {
		sql += " " + t.Hint.SQL()
	}
	if t.As != nil {
		sql += " AS " + t.As.SQL()
	}
	return sql
}

func (u *Unnest) SQL() string {
	sql := "UNNEST(" + u.Expr.SQL() + ")"
	if u.Hint != nil {
		sql += " " + u.Hint.SQL()
	}
	if u.As != nil {
		sql += " AS " + u.As.SQL()
	}
	if u.WithOffset != nil {
		sql += " WITH OFFSET " + u.WithOffset.SQL()
	}
	return sql
}

func (s *SubQueryJoinExpr) SQL() string {
	sql := s.Expr.SQL()
	if s.Hint != nil {
		sql += " " + s.Hint.SQL()
	}
	if s.As != nil {
		sql += " AS " + s.As.SQL()
	}
	return sql
}

func (p *ParenJoinExpr) SQL() string {
	return "(" + p.Expr.SQL() + ")"
}

func (j *Join) SQL() string {
	sql := j.Left.SQL()
	sql += " " + string(j.Op) + " JOIN "
	if j.Hint != nil {
		sql += " " + j.Hint.SQL() + " "
	}
	sql += j.Right.SQL()
	if j.Cond != nil {
		sql += " " + j.Cond.SQL()
	}
	return sql
}

func (j *JoinCondition) SQL() string {
	if j.On != nil {
		return "ON " + j.On.SQL()
	}

	return "USING (" + j.Using.SQL() + ")"
}

func (o *OrderExpr) SQL() string {
	sql := o.Expr.SQL()
	sql += " " + string(o.Dir)
	return sql
}

func (o OrderExprList) SQL() string {
	var ss []string
	for _, i := range o {
		ss = append(ss, i.SQL())
	}
	return strings.Join(ss, ", ")
}

func (l *Limit) SQL() string {
	sql := l.Count.SQL()
	if l.Offset != nil {
		sql += " OFFSET " + l.Offset.SQL()
	}
	return sql
}

func (b *BinaryExpr) SQL() string {
	p := exprPrec(b)
	sql := paren(p, b.Left)
	sql += " " + string(b.Op) + " "
	sql += paren(p, b.Right)
	return sql
}

func (u *UnaryExpr) SQL() string {
	p := exprPrec(u)
	if u.Op == OpNot {
		return "NOT " + paren(p, u.Expr)
	}
	return string(u.Op) + paren(p, u.Expr)
}

func (i *InExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	if i.Not {
		sql += " NOT"
	}
	sql += " IN "
	sql += i.Right.SQL()
	return sql
}

func (i *InCondition) SQL() string {
	if i.Unnest != nil {
		return "UNNEST(" + i.Unnest.SQL() + ")"
	}
	if i.SubQuery != nil {
		return i.SubQuery.SQL()
	}
	return "(" + i.Values.SQL() + ")"
}

func (i *IsNullExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	sql += " IS "
	if i.Not {
		sql += "NOT "
	}
	sql += "NULL"
	return sql
}

func (i *IsBoolExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	sql += " IS "
	if i.Not {
		sql += "NOT "
	}
	if i.Right {
		sql += "TRUE"
	} else {
		sql += "FALSE"
	}
	return sql
}

func (b *BetweenExpr) SQL() string {
	p := exprPrec(b)
	sql := paren(p, b.Left)
	if b.Not {
		sql += " NOT"
	}
	sql += " BETWEEN "
	sql += paren(p, b.RightStart)
	sql += " AND "
	sql += paren(p, b.RightEnd)
	return sql
}

func (s *SelectorExpr) SQL() string {
	p := exprPrec(s)
	sql := paren(p, s.Left)
	sql += "." + s.Right.SQL()
	return sql
}

func (i *IndexExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	sql += "[" + i.Right.SQL() + "]"
	return sql
}

func (c *CallExpr) SQL() string {
	return fmt.Sprintf("%s(%s)", c.Func.SQL(), c.Args.SQL())
}

func (c *CastExpr) SQL() string {
	return fmt.Sprintf("CAST(%s AS %s)", paren(precLit, c.Expr), c.Type.SQL())
}

func (c *CaseExpr) SQL() string {
	sql := "CASE"
	if c.Expr != nil {
		sql += " " + c.Expr.SQL()
	}
	for _, w := range c.When {
		sql += " WHEN" + w.SQL()
	}
	if c.Else != nil {
		sql += " ELSE " + c.Else.SQL()
	}
	return sql
}

func (w *When) SQL() string {
	return w.Cond.SQL() + " THEN " + w.Then.SQL()
}

func (s *SubQuery) SQL() string {
	return "(" + s.Expr.SQL() + ")"
}

func (p *ParenExpr) SQL() string {
	return "(" + p.Expr.SQL() + ")"
}

func (a *ArrayExpr) SQL() string {
	return "ARRAY" + a.Expr.SQL()
}

func (e *ExistsExpr) SQL() string {
	sql := "EXISTS"
	if e.Hint != nil {
		sql += " " + e.Hint.SQL() + " "
	}
	sql += e.Expr.SQL()
	return sql
}

func (p *Param) SQL() string {
	return "@" + p.Name
}

func (i *Ident) SQL() string {
	return i.Name // TODO: correct escape
}

func (i IdentList) SQL() string {
	var ss []string
	for _, s := range i {
		ss = append(ss, s.SQL())
	}
	return strings.Join(ss, ", ")
}

func (a *ArrayLit) SQL() string {
	sql := "ARRAY"
	if a.Type != nil {
		sql += "<" + a.Type.SQL() + ">"
	}
	sql += "[" + a.Values.SQL() + "]"
	return sql
}

func (s *StructLit) SQL() string {
	sql := "STRUCT"
	if s.Fields != nil {
		sql += "<"
		for i, f := range s.Fields {
			if i != 0 {
				sql += ", "
			}
			sql += f.SQL()
		}
		sql += ">"
	}
	sql += "(" + s.Values.SQL() + ")"
	return sql
}

func (*NullLit) SQL() string {
	return "NULL"
}

func (b *BoolLit) SQL() string {
	if b.Value {
		return "TRUE"
	} else {
		return "FALSE"
	}
}

func (i *IntLit) SQL() string {
	return i.Value
}

func (f *FloatLit) SQL() string {
	return f.Value
}

func (s *StringLit) SQL() string {
	return strconv.Quote(s.Value) // TODO: correct escape
}

func (b *BytesLit) SQL() string {
	return "B" + strconv.Quote(string(b.Value)) // TODO: correct escape
}

func (d *DateLit) SQL() string {
	return "DATE" + strconv.Quote(d.Value) // TODO: correct escape
}

func (t *TimestampLit) SQL() string {
	return "TIMESTAMP" + strconv.Quote(t.Value) // TODO: correct escape
}

func (t *Type) SQL() string {
	sql := string(t.Name)
	if t.Fields != nil {
		sql += "<"
		for i, f := range t.Fields {
			if i != 0 {
				sql += ", "
			}
			sql += f.SQL()
		}
		sql += ">"
	}
	if t.Value != nil {
		sql += "<" + t.Value.SQL() + ">"
	}
	return sql
}

func (t *FieldSchema) SQL() string {
	var sql string
	if t.Name != nil {
		sql += t.Name.SQL() + " "
	}
	sql += t.Type.SQL()
	return sql
}

func (c *CastIntValue) SQL() string {
	return "CAST(" + c.Expr.SQL() + " AS INT64)"
}
