package analyzer

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

func (a *Analyzer) analyzeSelectWithGroupBy(s *parser.Select) NameList {
	ti := a.analyzeFrom(s.From)

	a.pushTableInfo(ti)
	a.analyzeWhere(s.Where)

	var lists []NameList
	for _, item := range s.Results {
		itemList := a.analyzeSelectItem(item)
		lists = append(lists, itemList)
	}

	list, gbc := a.analyzeGroupBy(s.Results, s.GroupBy, lists)
	a.scope.context = gbc
	a.analyzeSelectResultsAfterGroupBy(s.Results, gbc)

	a.pushNameList(list)
	a.analyzeHaving(s.Having)
	a.analyzeOrderBy(s.OrderBy)
	a.analyzeLimit(s.Limit)
	a.popScope()
	a.popScope()

	if s.AsStruct {
		return NameList{makeNameListColumnName(list, s)}
	}

	return list
}

func (a *Analyzer) analyzeHaving(w *parser.Having) {
	// TODO: implement
}

func (a *Analyzer) analyzeGroupBy(results []parser.SelectItem, g *parser.GroupBy, lists []NameList) (NameList, *GroupByContext) {
	var list NameList
	for _, itemList := range lists {
		list = append(list, itemList...)
	}

	listsMap := make(map[parser.SelectItem]NameList)
	for i, item := range results {
		listsMap[item] = lists[i]
	}

	gbc := &GroupByContext{
		Lists: listsMap,
	}

	for _, expr := range g.Exprs {
		e := simplifyExpr(expr)
		switch e := e.(type) {
		case *parser.Ident:
			name := list.Lookup(e.Name)
			if name != nil {
				if name.Ambiguous {
					a.panicf(e, "ambiguous name: %s", e.SQL())
				}
				gbc.AddValidName(name)
				continue
			}
		case *parser.Path:
			name := list.Lookup(e.Idents[0].Name)
			if name != nil {
				a.panicf(e.Idents[1], "cannot access field of SELECT result column: %s", e.Idents[1].SQL())
			}
		case *parser.IntLiteral:
			v, err := strconv.ParseInt(e.Value, e.Base, 64)
			if err != nil {
				a.panicf(e, "error on parsing integer literal: %v", err)
			}
			if 1 <= v && int(v) <= len(list) {
				gbc.AddValidName(list[v-1])
				continue
			}
		}
		t := a.analyzeExpr(expr)
		if t.Name != nil {
			gbc.AddValidName(t.Name)
		}
		gbc.ValidExprs = append(gbc.ValidExprs, expr)
	}

	return list, gbc
}

func (a *Analyzer) analyzeSelectResultsAfterGroupBy(results []parser.SelectItem, gbc *GroupByContext) {
	for _, item := range results {
		list := gbc.Lists[item]

		hasValidName := false
		for _, name := range list {
			if gbc.IsValidName(name) {
				hasValidName = true
				break
			}
		}

		_, isStar := item.(*parser.Star)
		_, isStarPath := item.(*parser.StarPath)

		if hasValidName || isStar || isStarPath {
			for _, name := range list {
				if !gbc.IsValidName(name) {
					a.panicf(name.Node, "star-expansion contains invalid name after GROUP BY: %s", name.Quote())
				}
			}
			continue
		}

		switch item := item.(type) {
		case *parser.Alias:
			a.analyzeExprAfterGroupBy(item.Expr, gbc)
		case *parser.ExprSelectItem:
			a.analyzeExprAfterGroupBy(item.Expr, gbc)
		default:
			panic("BUG: unreachable")
		}

		for _, name := range list {
			gbc.AddValidName(name)
		}
	}
}

func (a *Analyzer) analyzeExprAfterGroupBy(expr parser.Expr, gbc *GroupByContext) {
	for _, validExpr := range gbc.ValidExprs {
		if isSameExprForGroupBy(expr, validExpr) {
			return
		}
	}

	a.analyzeExpr(expr)
}

func isSameExprForGroupBy(expr1, expr2 parser.Expr) bool {
	e1 := simplifyExpr(expr1)
	e2 := simplifyExpr(expr2)

	switch e1 := e1.(type) {
	case *parser.BinaryExpr:
		e2, ok := e2.(*parser.BinaryExpr)
		if !ok {
			return false
		}
		return e1.Op == e2.Op && isSameExprForGroupBy(e1.Left, e2.Left) && isSameExprForGroupBy(e1.Right, e2.Right)
	case *parser.UnaryExpr:
		e2, ok := e2.(*parser.UnaryExpr)
		if !ok {
			return false
		}
		return e1.Op == e2.Op && isSameExprForGroupBy(e1.Expr, e2.Expr)
	case *parser.Ident:
		e2, ok := e2.(*parser.Ident)
		if !ok {
			return false
		}
		return strings.EqualFold(e1.Name, e2.Name)
	case *parser.Param:
		e2, ok := e2.(*parser.Ident)
		if !ok {
			return false
		}
		return strings.EqualFold(e1.Name, e2.Name)
	case *parser.NullLiteral:
		_, ok := e2.(*parser.NullLiteral)
		return ok
	case *parser.BoolLiteral:
		e2, ok := e2.(*parser.BoolLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *parser.IntLiteral:
		e2, ok := e2.(*parser.IntLiteral)
		if !ok {
			return false
		}
		v1, err := strconv.ParseInt(e1.Value, e1.Base, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseInt(e2.Value, e2.Base, 64)
		if err != nil {
			return false
		}
		return v1 == v2
	case *parser.FloatLiteral:
		e2, ok := e2.(*parser.FloatLiteral)
		if !ok {
			return false
		}
		v1, err := strconv.ParseFloat(e1.Value, 64)
		if err != nil {
			return false
		}
		v2, err := strconv.ParseFloat(e2.Value, 64)
		if err != nil {
			return false
		}
		return v1 == v2
	case *parser.StringLiteral:
		e2, ok := e2.(*parser.StringLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *parser.BytesLiteral:
		e2, ok := e2.(*parser.BytesLiteral)
		if !ok {
			return false
		}
		return bytes.Equal(e1.Value, e2.Value)
	case *parser.DateLiteral:
		e2, ok := e2.(*parser.DateLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *parser.TimestampLiteral:
		e2, ok := e2.(*parser.TimestampLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	}

	return false
}

func simplifyExpr(e parser.Expr) parser.Expr {
	if e, ok := e.(*parser.ParenExpr); ok {
		return simplifyExpr(e.Expr)
	}
	return e
}
