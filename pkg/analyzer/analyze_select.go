package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

func (a *Analyzer) analyzeQueryStatement(q *parser.QueryStatement) {
	// TODO: check q.Hint
	_ = a.analyzeQueryExpr(q.Query)
}

func (a *Analyzer) analyzeQueryExpr(q parser.QueryExpr) NameList {
	var list NameList
	switch q := q.(type) {
	case *parser.Select:
		list = a.analyzeSelect(q)
	case *parser.CompoundQuery:
		list = a.analyzeCompoundQuery(q)
	case *parser.SubQuery:
		list = a.analyzeSubQuery(q)
	}

	if a.NameLists == nil {
		a.NameLists = make(map[parser.QueryExpr]NameList)
	}
	a.NameLists[q] = list
	return list
}

func (a *Analyzer) analyzeSelect(s *parser.Select) NameList {
	switch {
	case s.From == nil:
		return a.analyzeSelectWithoutFrom(s)
	case s.GroupBy == nil:
		return a.analyzeSelectWithoutGroupBy(s)
	}

	return a.analyzeSelectWithGroupBy(s)
}

func (a *Analyzer) analyzeSelectWithoutFrom(s *parser.Select) NameList {
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

	a.pushTableInfo(&TableInfo{}) // prevent working SELECT * in subquery
	var list NameList
	for _, item := range s.Results {
		if hasAggregateFuncInSelectItem(item) {
			a.panicf(item, "SELECT without FROM cannot have aggregate function call")
		}

		itemList := a.analyzeSelectItem(item)
		list = append(list, itemList...)
	}

	a.analyzeLimit(s.Limit)
	a.popScope()

	if s.AsStruct {
		return NameList{makeNameListColumnName(list, s)}
	}

	return list
}

func (a *Analyzer) analyzeSelectWithoutGroupBy(s *parser.Select) NameList {
	if s.Having != nil {
		a.panicf(s.Having, "SELECT without GROUP BY cannot have HAVING clause")
	}

	ti := a.analyzeFrom(s.From)
	a.pushTableInfo(ti)
	a.analyzeWhere(s.Where)

	var lists []NameList
	for _, item := range s.Results {
		itemList := a.analyzeSelectItem(item)
		lists = append(lists, itemList)
	}

	var list NameList
	for _, itemList := range lists {
		list = append(list, itemList...)
	}

	listsMap := make(map[parser.SelectItem]NameList)
	hasAggregate := false

	for i, item := range s.Results {
		listsMap[item] = lists[i]
		if hasAggregateFuncInSelectItem(item) {
			hasAggregate = true
		}
	}

	gbc := &GroupByContext{
		Lists: listsMap,
	}

	if hasAggregate {
		a.analyzeSelectResultsAfterGroupBy(s.Results, gbc)
	}

	a.pushNameList(list)
	a.analyzeOrderBy(s.OrderBy)
	a.analyzeLimit(s.Limit)
	a.popScope()
	a.popScope()

	if s.AsStruct {
		return NameList{makeNameListColumnName(list, s)}
	}

	return list
}

func (a *Analyzer) analyzeCompoundQuery(q *parser.CompoundQuery) NameList {
	var lists []NameList

	for _, query := range q.Queries {
		lists = append(lists, a.analyzeQueryExpr(query))
	}

	for i, l := range lists {
		if len(l) != len(lists[0]) {
			a.panicf(q.Queries[i], "queries in set operation have mismatched column count")
		}
	}

	list := make(NameList, len(lists[0]))
	for i := 0; i < len(list); i++ {
		names := make([]*Name, len(lists))
		for j, l := range lists {
			names[j] = l[i]
		}

		name, ok := makeCompoundQueryResultName(names, q)
		if !ok {
			ts := make([]string, len(names))
			for j, name := range names {
				ts[j] = TypeString(name.Type)
			}
			a.panicf(q, "column %d of queries in set operation have incompatible type %s", i+1, strings.Join(ts, ","))
		}

		list[i] = name
	}

	a.pushNameList(list)
	a.analyzeOrderBy(q.OrderBy)
	a.analyzeLimit(q.Limit)
	a.popScope()

	return list
}

func (a *Analyzer) analyzeSubQuery(s *parser.SubQuery) NameList {
	list := a.analyzeQueryExpr(s.Query)

	a.pushNameList(list)
	a.analyzeOrderBy(s.OrderBy)
	a.analyzeLimit(s.Limit)
	a.popScope()

	return list
}

func (a *Analyzer) analyzeSelectItem(s parser.SelectItem) NameList {
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

func (a *Analyzer) analyzeStar(s *parser.Star) NameList {
	if a.scope == nil || a.scope.List == nil {
		a.panicf(s, "SELECT * must have a FROM clause")
	}
	return a.scope.List
}

func (a *Analyzer) analyzeStarPath(s *parser.StarPath) NameList {
	t := a.analyzeExpr(s.Expr)

	var list NameList
	if t.Name != nil {
		list = t.Name.Children()
	} else {
		list = makeNameListFromType(t.Type, s)
	}

	if list == nil {
		a.panicf(s, "star expansion is not supported for type %s", TypeString(t.Type))
	}

	return list
}

func (a *Analyzer) analyzeAlias(s *parser.Alias) NameList {
	t := a.analyzeExpr(s.Expr)
	if t.Name != nil {
		return NameList{makeAliasName(t.Name, s, s.As.Alias)}
	}
	return NameList{makeExprColumnName(t.Type, s.Expr, s, s.As.Alias)}
}

func (a *Analyzer) analyzeExprSelectItem(s *parser.ExprSelectItem) NameList {
	t := a.analyzeExpr(s.Expr)
	if t.Name != nil {
		return NameList{makeAliasName(t.Name, s, extractIdentFromExpr(s.Expr))}
	}
	return NameList{makeExprColumnName(t.Type, s.Expr, s, nil)}
}

func (a *Analyzer) analyzeWhere(w *parser.Where) {
	if w == nil {
		return
	}

	t := a.analyzeExpr(w.Expr)
	if !TypeCoerce(t.Type, BoolType) {
		a.panicf(w, "WHERE clause expression require BOOL, but: %s", TypeString(t.Type))
	}
}

func (a *Analyzer) analyzeOrderBy(o *parser.OrderBy) {
	if o == nil {
		return
	}

	for _, item := range o.Items {
		a.analyzeExpr(item.Expr)
		if item.Collate != nil {
			// TODO: check COLLATE value more
			a.analyzeStringValue(item.Collate.Value)
		}
	}
}

func (a *Analyzer) analyzeLimit(l *parser.Limit) {
	if l == nil {
		return
	}

	a.analyzeIntValue(l.Count)
	if l.Offset != nil {
		a.analyzeIntValue(l.Offset.Value)
	}
}
