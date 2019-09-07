package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func (a *Analyzer) analyzeQueryStatement(q *parser.QueryStatement) {
	// TODO: analyze q.Hint
	_ = a.analyzeQueryExpr(q.Query)
}

func (a *Analyzer) analyzeQueryExpr(q parser.QueryExpr) SelectList {
	var list SelectList
	switch q := q.(type) {
	case *parser.Select:
		list = a.analyzeSelect(q)
	case *parser.CompoundQuery:
		list = a.analyzeCompoundQuery(q)
	case *parser.SubQuery:
		list = a.analyzeSubQuery(q)
	}

	if a.SelectLists == nil {
		a.SelectLists = make(map[parser.QueryExpr]SelectList)
	}
	a.SelectLists[q] = list
	return list
}

func (a *Analyzer) analyzeSelect(s *parser.Select) SelectList {
	switch {
	case s.From == nil:
		return a.analyzeSelectWithoutFrom(s)
	case s.GroupBy == nil:
		return a.analyzeSelectWithoutGroupBy(s)
	}

	panic("TODO")
}

func (a *Analyzer) analyzeSelectWithoutFrom(s *parser.Select) SelectList {
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

	var list SelectList
	for _, item := range s.Results {
		if hasAggregateFuncInSelectItem(item) {
			a.panicf(item, "SELECT without FROM cannot have aggregate function call")
		}

		itemList := a.analyzeSelectItem(item)
		list = append(list, itemList...)
	}

	a.analyzeLimit(s.Limit)

	if s.AsStruct {
		t := list.toType()
		return SelectList{newColumnReference("", t, s)}
	}

	return list
}

func (a *Analyzer) analyzeSelectWithoutGroupBy(s *parser.Select) SelectList {
	if s.Having != nil {
		a.panicf(s.Having, "SELECT without GROUP BY cannot have HAVING clause")
	}

	t := a.analyzeFrom(s.From)

	a.pushTableScope(t)
	a.analyzeWhere(s.Where)

	hasAgg := false
	for _, item := range s.Results {
		if hasAggregateFuncInSelectItem(item) {
			hasAgg = true
			break
		}
	}
	if hasAgg {
		return a.analyzeSelectWithoutGroupByAggregate(s)
	}

	var list SelectList
	for _, item := range s.Results {
		itemList := a.analyzeSelectItem(item)
		list = append(list, itemList...)
	}

	a.pushSelectListScope(list)
	a.analyzeOrderBy(s.OrderBy)
	a.analyzeLimit(s.Limit)
	a.popScope()
	a.popScope()

	return list
}

func (a *Analyzer) analyzeSelectWithoutGroupByAggregate(s *parser.Select) SelectList {
	panic("TODO: implement")
}

func (a *Analyzer) analyzeCompoundQuery(q *parser.CompoundQuery) SelectList {
	list := a.analyzeQueryExpr(q.Queries[0]).deriveSimple(q)

	for _, query := range q.Queries[1:] {
		queryList := a.analyzeQueryExpr(query)

		if len(list) != len(queryList) {
			a.panicf(query, "queries in set operation have mismatched column count")
		}

		for i, r := range list {
			if !r.merge(queryList[i]) {
				a.panicf(
					queryList[i].GetNode(q),
					"%s is incompatible with %s (column %d)",
					TypeString(r.Type),
					TypeString(queryList[i].Type),
					i+1,
				)
			}
		}
	}

	a.pushSelectListScope(list)
	a.analyzeOrderBy(q.OrderBy)
	a.analyzeLimit(q.Limit)
	a.popScope()

	return list
}

func (a *Analyzer) analyzeSubQuery(s *parser.SubQuery) SelectList {
	panic("TODO: implement")
}

func (a *Analyzer) analyzeSelectItem(s parser.SelectItem) SelectList {
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

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeStar(s *parser.Star) SelectList {
	if a.scope == nil || a.scope.List == nil {
		a.panicf(s, "SELECT * must have a FROM clause")
	}
	return a.scope.List.deriveSimple(s)
}

func (a *Analyzer) analyzeStarPath(s *parser.StarPath) SelectList {
	t := a.analyzeExpr(s.Expr)
	list := typeToSelectList(t.Type, s)
	if list == nil {
		a.panicf(s, "star expansion is not supported for type %s", TypeString(t.Type))
	}
	return list.deriveSimple(s)
}

func (a *Analyzer) analyzeAlias(s *parser.Alias) SelectList {
	t := a.analyzeExpr(s.Expr)
	return SelectList{newColumnReference(s.As.Alias.Name, t.Type, s)}
}

func (a *Analyzer) analyzeExprSelectItem(s *parser.ExprSelectItem) SelectList {
	t := a.analyzeExpr(s.Expr)
	name := extractNameFromExpr(s.Expr)
	return SelectList{newColumnReference(name, t.Type, s)}
}

func (a *Analyzer) analyzeWhere(w *parser.Where) {
	// TODO: implement
}

func (a *Analyzer) analyzeOrderBy(o *parser.OrderBy) {
	// TODO: implement
}

func (a *Analyzer) analyzeLimit(l *parser.Limit) {
	// TODO: implement
}
