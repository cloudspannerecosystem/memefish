package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/ast"
)

func isAggregateFuncName(name string) bool {
	switch strings.ToUpper(name) {
	case "ANY_VALUE", "ARRAY_AGG", "AVG",
		"BIT_AND", "BIT_OR", "BIT_XOR", "COUNT",
		"COUNTIF", "LOGICAL_AND", "LOGICAL_OR",
		"MAX", "MIN", "STRING_AGG", "SUM":
		return true
	}
	return false
}

func hasAggregateFunc(e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.BinaryExpr:
		return hasAggregateFunc(e.Left) || hasAggregateFunc(e.Right)
	case *ast.UnaryExpr:
		return hasAggregateFunc(e.Expr)
	case *ast.InExpr:
		if hasAggregateFunc(e.Left) {
			return true
		}
		switch r := e.Right.(type) {
		case *ast.UnnestInCondition:
			return hasAggregateFunc(r.Expr)
		case *ast.ValuesInCondition:
			for _, v := range r.Exprs {
				if hasAggregateFunc(v) {
					return true
				}
			}
			return false
		case *ast.SubQueryInCondition:
			// For example, `SELECT (SELECT SUM(x)) FROM table` is invalid
			// because subqeury is evaluated by other context.
			return false
		}
	case *ast.IsNullExpr:
		return hasAggregateFunc(e.Left)
	case *ast.IsBoolExpr:
		return hasAggregateFunc(e.Left)
	case *ast.BetweenExpr:
		return hasAggregateFunc(e.Left) || hasAggregateFunc(e.RightStart) || hasAggregateFunc(e.RightEnd)
	case *ast.SelectorExpr:
		return hasAggregateFunc(e.Expr)
	case *ast.IndexExpr:
		return hasAggregateFunc(e.Expr) || hasAggregateFunc(e.Index)
	case *ast.CallExpr:
		if isAggregateFuncName(e.Func.Name) {
			return true
		}
		for _, a := range e.Args {
			if hasAggregateFunc(a.Expr) {
				return true
			}
		}
		return false
	case *ast.CountStarExpr:
		return true
	case *ast.CastExpr:
		return hasAggregateFunc(e.Expr)
	case *ast.ExtractExpr:
		if e.AtTimeZone != nil && hasAggregateFunc(e.AtTimeZone.Expr) {
			return true
		}
		return hasAggregateFunc(e.Expr)
	case *ast.CaseExpr:
		if e.Expr != nil && hasAggregateFunc(e.Expr) {
			return true
		}
		for _, w := range e.Whens {
			if hasAggregateFunc(w.Cond) || hasAggregateFunc(w.Then) {
				return true
			}
		}
		return e.Else != nil && hasAggregateFunc(e.Else.Expr)
	case *ast.ParenExpr:
		return hasAggregateFunc(e.Expr)
	case *ast.ScalarSubQuery, *ast.ArraySubQuery, *ast.ExistsSubQuery,
		*ast.Param, *ast.Ident, *ast.Path:
		return false
	case *ast.ArrayLiteral:
		for _, v := range e.Values {
			if hasAggregateFunc(v) {
				return true
			}
		}
		return false
	case *ast.StructLiteral:
		for _, v := range e.Values {
			if hasAggregateFunc(v) {
				return true
			}
		}
		return false
	case *ast.NullLiteral, *ast.BoolLiteral, *ast.IntLiteral, *ast.FloatLiteral,
		*ast.StringLiteral, *ast.BytesLiteral, *ast.DateLiteral, *ast.TimestampLiteral:
		return false
	}

	panic("BUG: unreachable")
}

func hasAggregateFuncInSelectItem(s ast.SelectItem) bool {
	switch s := s.(type) {
	case *ast.Star:
		return false
	case *ast.DotStar:
		return hasAggregateFunc(s.Expr)
	case *ast.Alias:
		return hasAggregateFunc(s.Expr)
	case *ast.ExprSelectItem:
		return hasAggregateFunc(s.Expr)
	}

	panic("unreachable")
}
