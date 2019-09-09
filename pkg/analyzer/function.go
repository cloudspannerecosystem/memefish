package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
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

func hasAggregateFunc(e parser.Expr) bool {
	switch e := e.(type) {
	case *parser.BinaryExpr:
		return hasAggregateFunc(e.Left) || hasAggregateFunc(e.Right)
	case *parser.UnaryExpr:
		return hasAggregateFunc(e.Expr)
	case *parser.InExpr:
		if hasAggregateFunc(e.Left) {
			return true
		}
		switch r := e.Right.(type) {
		case *parser.UnnestInCondition:
			return hasAggregateFunc(r.Expr)
		case *parser.ValuesInCondition:
			for _, v := range r.Exprs {
				if hasAggregateFunc(v) {
					return true
				}
			}
			return false
		case *parser.SubQueryInCondition:
			// For example, `SELECT (SELECT SUM(x)) FROM table` is invalid
			// because subqeury is evaluated by other context.
			return false
		}
	case *parser.IsNullExpr:
		return hasAggregateFunc(e.Left)
	case *parser.IsBoolExpr:
		return hasAggregateFunc(e.Left)
	case *parser.BetweenExpr:
		return hasAggregateFunc(e.Left) || hasAggregateFunc(e.RightStart) || hasAggregateFunc(e.RightEnd)
	case *parser.SelectorExpr:
		return hasAggregateFunc(e.Expr)
	case *parser.IndexExpr:
		return hasAggregateFunc(e.Expr) || hasAggregateFunc(e.Index)
	case *parser.CallExpr:
		if isAggregateFuncName(e.Func.Name) {
			return true
		}
		for _, a := range e.Args {
			if hasAggregateFunc(a.Expr) {
				return true
			}
		}
		return false
	case *parser.CountStarExpr:
		return true
	case *parser.CastExpr:
		return hasAggregateFunc(e.Expr)
	case *parser.ExtractExpr:
		if e.AtTimeZone != nil && hasAggregateFunc(e.AtTimeZone.Expr) {
			return true
		}
		return hasAggregateFunc(e.Expr)
	case *parser.CaseExpr:
		if e.Expr != nil && hasAggregateFunc(e.Expr) {
			return true
		}
		for _, w := range e.Whens {
			if hasAggregateFunc(w.Cond) || hasAggregateFunc(w.Then) {
				return true
			}
		}
		return e.Else != nil && hasAggregateFunc(e.Else.Expr)
	case *parser.ParenExpr:
		return hasAggregateFunc(e.Expr)
	case *parser.ScalarSubQuery, *parser.ArraySubQuery, *parser.ExistsSubQuery,
		*parser.Param, *parser.Ident, *parser.Path:
		return false
	case *parser.ArrayLiteral:
		for _, v := range e.Values {
			if hasAggregateFunc(v) {
				return true
			}
		}
		return false
	case *parser.StructLiteral:
		for _, v := range e.Values {
			if hasAggregateFunc(v) {
				return true
			}
		}
		return false
	case *parser.NullLiteral, *parser.BoolLiteral, *parser.IntLiteral, *parser.FloatLiteral,
		*parser.StringLiteral, *parser.BytesLiteral, *parser.DateLiteral, *parser.TimestampLiteral:
		return false
	}

	panic("BUG: unreachable")
}

func hasAggregateFuncInSelectItem(s parser.SelectItem) bool {
	switch s := s.(type) {
	case *parser.Star:
		return false
	case *parser.DotStar:
		return hasAggregateFunc(s.Expr)
	case *parser.Alias:
		return hasAggregateFunc(s.Expr)
	case *parser.ExprSelectItem:
		return hasAggregateFunc(s.Expr)
	}

	panic("unreachable")
}
