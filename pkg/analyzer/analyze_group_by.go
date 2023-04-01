package analyzer

import (
	"bytes"
	"strconv"

	"github.com/cloudspannerecosystem/memefish/pkg/ast"
	"github.com/cloudspannerecosystem/memefish/pkg/char"
)

func (a *Analyzer) analyzeSelectWithGroupBy(s *ast.Select) NameList {
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

func (a *Analyzer) analyzeHaving(h *ast.Having) {
	if h == nil {
		return
	}

	t := a.analyzeExpr(h.Expr)
	if !TypeCoerce(t.Type, BoolType) {
		a.panicf(h, "HAVING clause expression require BOOL, but: %s", TypeString(t.Type))
	}
}

func (a *Analyzer) analyzeGroupBy(results []ast.SelectItem, g *ast.GroupBy, lists []NameList) (NameList, *GroupByContext) {
	var list NameList
	for _, itemList := range lists {
		list = append(list, itemList...)
	}

	listsMap := make(map[ast.SelectItem]NameList)
	for i, item := range results {
		listsMap[item] = lists[i]
	}

	gbc := &GroupByContext{
		Lists: listsMap,
	}

	for _, expr := range g.Exprs {
		e := simplifyExpr(expr)
		switch e := e.(type) {
		case *ast.Ident:
			name := list.Lookup(e.Name)
			if name != nil {
				if name.Ambiguous {
					a.panicf(e, "ambiguous name: %s", e.SQL())
				}
				gbc.AddValidName(name)
				continue
			}
		case *ast.Path:
			name := list.Lookup(e.Idents[0].Name)
			if name != nil {
				a.panicf(e.Idents[1], "cannot access field of SELECT result column: %s", e.Idents[1].SQL())
			}
		case *ast.IntLiteral:
			v, err := strconv.ParseInt(e.Value, e.Base, 64)
			if err != nil {
				a.panicf(e, "error on parsing integer literal: %v", err)
			}
			if 1 <= v && int(v) <= len(list) {
				gbc.AddValidName(list[v-1])
				continue
			}
		case *ast.Param:
			v, ok := a.lookupParam(e.Name)
			if !ok {
				a.panicf(e, "unknown query parameter: %s", e.SQL())
			}
			iv, ok := v.(int64)
			if ok && 1 <= iv && int(iv) <= len(list) {
				gbc.AddValidName(list[iv-1])
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

func (a *Analyzer) analyzeSelectResultsAfterGroupBy(results []ast.SelectItem, gbc *GroupByContext) {
	for _, item := range results {
		list := gbc.Lists[item]

		hasValidName := false
		for _, name := range list {
			if gbc.IsValidName(name) {
				hasValidName = true
				break
			}
		}

		_, isStar := item.(*ast.Star)
		_, isDotStar := item.(*ast.DotStar)

		if hasValidName || isStar || isDotStar {
			for _, name := range list {
				if !gbc.IsValidName(name) {
					a.panicf(name.Node, "star-expansion contains invalid name after GROUP BY: %s", name.Quote())
				}
			}
			continue
		}

		switch item := item.(type) {
		case *ast.Alias:
			a.analyzeExprAfterGroupBy(item.Expr, gbc)
		case *ast.ExprSelectItem:
			a.analyzeExprAfterGroupBy(item.Expr, gbc)
		default:
			panic("BUG: unreachable")
		}

		for _, name := range list {
			gbc.AddValidName(name)
		}
	}
}

func (a *Analyzer) analyzeExprAfterGroupBy(expr ast.Expr, gbc *GroupByContext) {
	for _, validExpr := range gbc.ValidExprs {
		if isSameExprForGroupBy(expr, validExpr) {
			return
		}
	}

	// FIXME: it works, but it is very inefficient.
	a.analyzeExpr(expr)
}

func isSameExprForGroupBy(expr1, expr2 ast.Expr) bool {
	e1 := simplifyExpr(expr1)
	e2 := simplifyExpr(expr2)

	switch e1 := e1.(type) {
	case *ast.BinaryExpr:
		e2, ok := e2.(*ast.BinaryExpr)
		if !ok {
			return false
		}
		return e1.Op == e2.Op && isSameExprForGroupBy(e1.Left, e2.Left) && isSameExprForGroupBy(e1.Right, e2.Right)
	case *ast.UnaryExpr:
		e2, ok := e2.(*ast.UnaryExpr)
		if !ok {
			return false
		}
		return e1.Op == e2.Op && isSameExprForGroupBy(e1.Expr, e2.Expr)
	case *ast.Ident:
		e2, ok := e2.(*ast.Ident)
		if !ok {
			return false
		}
		return char.EqualFold(e1.Name, e2.Name)
	case *ast.Param:
		e2, ok := e2.(*ast.Ident)
		if !ok {
			return false
		}
		return char.EqualFold(e1.Name, e2.Name)
	case *ast.NullLiteral:
		_, ok := e2.(*ast.NullLiteral)
		return ok
	case *ast.BoolLiteral:
		e2, ok := e2.(*ast.BoolLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *ast.IntLiteral:
		e2, ok := e2.(*ast.IntLiteral)
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
	case *ast.FloatLiteral:
		e2, ok := e2.(*ast.FloatLiteral)
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
	case *ast.StringLiteral:
		e2, ok := e2.(*ast.StringLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *ast.BytesLiteral:
		e2, ok := e2.(*ast.BytesLiteral)
		if !ok {
			return false
		}
		return bytes.Equal(e1.Value, e2.Value)
	case *ast.DateLiteral:
		e2, ok := e2.(*ast.DateLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	case *ast.TimestampLiteral:
		e2, ok := e2.(*ast.TimestampLiteral)
		if !ok {
			return false
		}
		return e1.Value == e2.Value
	}

	// TODO: handle missing ASTs

	return false
}

func simplifyExpr(e ast.Expr) ast.Expr {
	if e, ok := e.(*ast.ParenExpr); ok {
		return simplifyExpr(e.Expr)
	}
	return e
}
