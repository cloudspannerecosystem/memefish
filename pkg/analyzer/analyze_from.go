package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

func (a *Analyzer) analyzeFrom(f *parser.From) *TableScope {
	return a.analyzeTableExpr(f.Source, newTableScope())
}

func (a *Analyzer) analyzeTableExpr(e parser.TableExpr, ts *TableScope) *TableScope {
	switch e := e.(type) {
	case *parser.Unnest:
		return a.analyzeUnnest(e, ts)
	case *parser.SubQueryTableExpr:
		return a.analyzeSubQueryJoinExpr(e, ts)
	case *parser.ParenTableExpr:
		return a.analyzeParenJoinExpr(e, ts)
	case *parser.Join:
		return a.analyzeJoin(e, ts)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeUnnest(e *parser.Unnest, ts *TableScope) *TableScope {
	a.pushTableScope(ts)
	t := a.analyzeExpr(e.Expr)
	a.popScope()

	tt, ok := TypeCastArray(t.Type)
	if !ok {
		a.panicf(e, "UNNEST value must be ARRAY, but: %s", TypeString(t.Type))
	}

	var id *parser.Ident
	if e.As != nil {
		id = e.As.Alias
	} else if e.Implicit {
		id = extractIdentFromExpr(e.Expr)
	}

	var name string
	if id != nil {
		name = id.Name
	}

	refs := make(map[string]*Reference)
	tableRef := &Reference{
		Kind:  TableRef,
		Name:  name,
		Type:  tt.Item,
		Node:  e,
		Ident: id,
	}

	if name != "" {
		refs[name] = tableRef
	}

	list := typeToSelectList(tt.Item, e)
	if list == nil {
		list = SelectList{tableRef.derive(nil, ColumnRef)}
	}

	// TODO: check e.Hint

	// check WITH OFFSET clause
	if e.WithOffset != nil {
		offsetName := "offset"
		var offsetId *parser.Ident
		if e.WithOffset.As != nil {
			offsetId = e.WithOffset.As.Alias
			offsetName = offsetId.Name
		}
		if offsetName == name {
			n := parser.Node(e.WithOffset)
			if offsetId != nil {
				n = offsetId
			}
			a.panicf(n, "duplicate alias %s found", parser.QuoteSQLIdent(name))
		}
		offsetRef := &Reference{
			Kind:  TableRef,
			Name:  offsetName,
			Type:  Int64Type,
			Node:  e.WithOffset,
			Ident: offsetId,
		}
		refs[offsetName] = offsetRef
		list = append(list, offsetRef.derive(nil, ColumnRef))
	}

	// TODO: check e.Sample

	return &TableScope{
		Refs: refs,
		List: list,
	}
}

func (a *Analyzer) analyzeSubQueryJoinExpr(e *parser.SubQueryTableExpr, ts *TableScope) *TableScope {
	list := a.analyzeQueryExpr(e.Query)

	var id *parser.Ident
	if e.As != nil {
		id = e.As.Alias
	}

	var name string
	if id != nil {
		name = id.Name
	}

	refs := make(map[string]*Reference)
	if name != "" {
		refs[name] = &Reference{
			Kind:  TableRef,
			Name:  name,
			Type:  list.toType(),
			Node:  e,
			Ident: id,
		}
	}

	return &TableScope{
		Refs: refs,
		List: list,
	}
}

func (a *Analyzer) analyzeParenJoinExpr(e *parser.ParenTableExpr, ts *TableScope) *TableScope {
	return a.analyzeTableExpr(e.Source, newTableScope())
}

func (a *Analyzer) analyzeJoin(j *parser.Join, ts *TableScope) *TableScope {
	lts := a.analyzeTableExpr(j.Left, ts)
	rts := a.analyzeTableExpr(j.Right, a.merge(ts, lts, j.Left))

	// TODO: check j.Method and j.Hint

	if j.Op == parser.CommaJoin || j.Op == parser.CrossJoin {
		if j.Cond != nil {
			a.panicf(j.Cond, "CROSS JOIN cannot have ON or USING clause")
		}
		return a.merge(lts, rts, j)
	}

	if j.Cond == nil {
		a.panicf(j, "%s must have ON or USING clause", j.Op)
	}

	var result *TableScope

	switch cond := j.Cond.(type) {
	case *parser.On:
		result := a.merge(lts, rts, j)
		a.pushTableScope(result)
		t := a.analyzeExpr(cond.Expr)
		a.popScope()
		if !TypeCoerce(t.Type, BoolType) {
			a.panicf(cond.Expr, "ON clause expression must be BOOL")
		}
	case *parser.Using:
		names := make(map[string]bool)
		for _, id := range cond.Idents {
			names[id.Name] = false
		}
		refs := a.mergeRefs(lts.Refs, rts.Refs, names, j)

		var list SelectList
		for _, id := range cond.Idents {
			if names[id.Name] {
				continue
			}
			names[id.Name] = true

			lref := lts.List.LookupRef(id.Name)
			if lref == nil {
				a.panicf(id, "USING condition %s is not found in left-side", id.SQL())
			}
			rref := rts.List.LookupRef(id.Name)
			if rref == nil {
				a.panicf(id, "USING condition %s is not found in right-side", id.SQL())
			}
			if !(TypeCoerce(lref.Type, rref.Type) || TypeCoerce(rref.Type, lref.Type)) {
				a.panicf(id, "USING condition %s is incompatible type: %s and %s", id.SQL(), TypeString(lref.Type), TypeString(rref.Type))
			}
			switch j.Op {
			case parser.InnerJoin, parser.LeftOuterJoin:
				ref := lref.deriveSimple(nil)
				ref.Origin = append(ref.Origin, rref)
				refs[id.Name] = ref
			case parser.RightOuterJoin:
				ref := rref.deriveSimple(nil)
				ref.Origin = append(ref.Origin, lref)
				refs[id.Name] = ref
			case parser.FullOuterJoin:
				ref := lref.deriveSimple(nil)
				if !ref.merge(rref) {
					a.panicf(id, "USING condition %s is incompatible type: %s and %s", id.SQL(), TypeString(lref.Type), TypeString(rref.Type))
				}
				refs[id.Name] = ref
			default:
				panic("BUG: unreachable")
			}
			list = append(list, refs[id.Name])
		}

		for _, ref := range lts.List {
			if r, ok := refs[ref.Name]; ok && r.Kind == ColumnRef {
				continue
			}
			list = append(list, ref)
		}
		for _, ref := range rts.List {
			if r, ok := refs[ref.Name]; ok && r.Kind == ColumnRef {
				continue
			}
			list = append(list, ref)
		}

		result = &TableScope{
			Refs: refs,
			List: list,
		}
	}

	return result
}

func (a *Analyzer) merge(t, u *TableScope, n parser.Node) *TableScope {
	s := &TableScope{}

	s.List = append(s.List, t.List...)
	s.List = append(s.List, u.List...)
	s.Refs = a.mergeRefs(t.Refs, u.Refs, make(map[string]bool), n)
	return s
}

func (a *Analyzer) mergeRefs(tRefs, uRefs map[string]*Reference, names map[string]bool, n parser.Node) map[string]*Reference {
	refs := make(map[string]*Reference)

	for name, ref := range tRefs {
		if _, ok := names[name]; ok {
			continue
		}
		refs[name] = ref
	}

	for name, ref := range uRefs {
		if _, ok := names[name]; ok {
			continue
		}
		if tRef, ok := refs[name]; ok {
			if tRef.Kind == TableRef && ref.Kind == TableRef {
				a.panicf(ref.GetIdent(n), "duplicate alias %s found", parser.QuoteSQLIdent(name))
			}
			if tRef.Kind == TableRef && ref.Kind == ColumnRef {
				continue
			}
		}
		refs[name] = ref
	}

	return refs
}
