package analyzer

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type TableInfo struct {
	List NameList
	Env  NameEnv
}

func (ti *TableInfo) toNameScope(next *NameScope) *NameScope {
	return &NameScope{
		List: ti.List,
		Env:  ti.Env,
		Next: next,
	}
}

func (a *Analyzer) analyzeFrom(f *parser.From) *TableInfo {
	return a.analyzeTableExpr(f.Source, &TableInfo{})
}

func (a *Analyzer) analyzeTableExpr(e parser.TableExpr, ti *TableInfo) *TableInfo {
	switch e := e.(type) {
	case *parser.TableName:
		return a.analyzeTableName(e, ti)
	case *parser.Unnest:
		return a.analyzeUnnest(e, ti)
	case *parser.SubQueryTableExpr:
		return a.analyzeSubQueryTableExpr(e, ti)
	case *parser.ParenTableExpr:
		return a.analyzeParenTableExpr(e, ti)
	case *parser.Join:
		return a.analyzeJoin(e, ti)
	}

	panic("BUG: unreachable")
}

func (a *Analyzer) analyzeTableName(e *parser.TableName, ti *TableInfo) *TableInfo {
	table, ok := a.lookupTable(e.Table.Name)
	if !ok {
		a.panicf(e, "unknown table: %s", e.Table.SQL())
	}

	var ident *parser.Ident
	if e.As != nil {
		ident = e.As.Alias
	}

	name := makeTableSchemaName(table, e, ident)
	list := NameList(name.Children())
	env := list.toNameEnv()
	err := env.Insert(name)
	if err != nil {
		panic(fmt.Sprintf("BUG: unexpected error: %v", err))
	}
	return &TableInfo{
		List: list,
		Env:  env,
	}
}

func (a *Analyzer) analyzeUnnest(e *parser.Unnest, ti *TableInfo) *TableInfo {
	a.pushTableInfo(ti)
	t := a.analyzeExpr(e.Expr)
	a.popScope()

	tt, ok := TypeCastArray(t.Type)
	if !ok {
		a.panicf(e, "UNNEST value must be ARRAY, but: %s", TypeString(t.Type))
	}

	var ident *parser.Ident
	if e.As != nil {
		ident = e.As.Alias
	} else if e.Implicit {
		ident = extractIdentFromExpr(e.Expr)
	}

	list := NameList{makeTableName("", tt.Item, e, ident)}
	result := list.toTableInfo()

	// TODO: check e.Hint

	// check WITH OFFSET clause
	if e.WithOffset != nil {
		result = a.mergeTableInfo(ti, a.analyzeWithOffset(e.WithOffset))
	}

	// TODO: check e.Sample

	return result
}

func (a *Analyzer) analyzeWithOffset(w *parser.WithOffset) *TableInfo {
	var ident *parser.Ident
	if w.As != nil {
		ident = w.As.Alias
	}

	list := NameList{makeTableName("offset", Int64Type, w, ident)}
	return list.toTableInfo()
}

func (a *Analyzer) analyzeSubQueryTableExpr(e *parser.SubQueryTableExpr, ti *TableInfo) *TableInfo {
	list := a.analyzeQueryExpr(e.Query)

	var ident *parser.Ident
	if e.As != nil {
		ident = e.As.Alias
	}

	name := makeNameListTableName(list, e, ident)

	env := list.toNameEnv()
	err := env.Insert(name)
	if err != nil {
		panic(fmt.Sprintf("BUG: unexpected error: %v", err))
	}
	return &TableInfo{
		List: name.Children(),
		Env:  env,
	}
}

func (a *Analyzer) analyzeParenTableExpr(e *parser.ParenTableExpr, ti *TableInfo) *TableInfo {
	return a.analyzeTableExpr(e.Source, &TableInfo{})
}

func (a *Analyzer) analyzeJoin(j *parser.Join, ti *TableInfo) *TableInfo {
	lti := a.analyzeTableExpr(j.Left, ti)
	rti := a.analyzeTableExpr(j.Right, a.mergeTableInfo(ti, lti))

	// TODO: check j.Method and j.Hint

	if j.Op == parser.CommaJoin || j.Op == parser.CrossJoin {
		if j.Cond != nil {
			a.panicf(j.Cond, "CROSS JOIN cannot have ON or USING clause")
		}
		return a.mergeTableInfo(lti, rti)
	}

	if j.Cond == nil {
		a.panicf(j, "%s must have ON or USING clause", j.Op)
	}

	var result *TableInfo

	switch cond := j.Cond.(type) {
	case *parser.On:
		result = a.mergeTableInfo(lti, rti)
		a.pushTableInfo(result)
		t := a.analyzeExpr(cond.Expr)
		a.popScope()
		if !TypeCoerce(t.Type, BoolType) {
			a.panicf(cond.Expr, "ON clause expression must be BOOL")
		}

	case *parser.Using:
		names := make(map[string]bool)
		for _, id := range cond.Idents {
			names[strings.ToUpper(id.Name)] = false
		}

		env := NameEnv{}
		for text, name := range lti.Env {
			if _, ok := names[text]; ok {
				continue
			}
			err := env.Insert(name)
			if err != nil {
				a.panicf(name.Ident, err.Error())
			}
		}
		for text, name := range rti.Env {
			if _, ok := names[text]; ok {
				continue
			}
			err := env.Insert(name)
			if err != nil {
				a.panicf(name.Ident, err.Error())
			}
		}

		var list NameList
		for _, id := range cond.Idents {
			text := strings.ToUpper(id.Name)
			if names[text] {
				continue
			}
			names[text] = true

			lname := lti.Env.Lookup(text)
			if lname == nil {
				a.panicf(id, "USING condition %s is not found in left-side", id.SQL())
			}
			rname := rti.Env.Lookup(text)
			if rname == nil {
				a.panicf(id, "USING condition %s is not found in right-side", id.SQL())
			}

			// TODO: check equality correctly
			if !(TypeCoerce(lname.Type, rname.Type) || TypeCoerce(rname.Type, lname.Type)) {
				a.panicf(
					id,
					"USING condition %s is incompatible type: %s and %s",
					id.SQL(), TypeString(lname.Type), TypeString(rname.Type),
				)
			}

			var name *Name
			switch j.Op {
			case parser.InnerJoin, parser.LeftOuterJoin:
				name = makeLeftJoinName(lname, rname)
			case parser.RightOuterJoin:
				name = makeRightJoinName(lname, rname)
			case parser.FullOuterJoin:
				var ok bool
				name, ok = makeFullJoinName(lname, rname)
				if !ok {
					a.panicf(
						id,
						"USING condition %s is incompatible type: %s and %s",
						id.SQL(), TypeString(lname.Type), TypeString(rname.Type),
					)
				}
			default:
				panic("BUG: unreachable")
			}
			env.InsertForce(name)
			list = append(list, name)
		}

		for _, name := range lti.List {
			if _, ok := names[strings.ToUpper(name.Text)]; ok {
				continue
			}
			list = append(list, name)
		}
		for _, name := range rti.List {
			if _, ok := names[strings.ToUpper(name.Text)]; ok {
				continue
			}
			list = append(list, name)
		}

		result = &TableInfo{
			List: list,
			Env:  env,
		}
	}

	return result
}

func (a *Analyzer) mergeTableInfo(ti1, ti2 *TableInfo) *TableInfo {
	var list NameList
	list = append(list, ti1.List...)
	list = append(list, ti2.List...)
	env := a.mergeNameEnv(ti1.Env, ti2.Env)
	return &TableInfo{
		List: list,
		Env:  env,
	}
}

func (a *Analyzer) mergeNameEnv(env1, env2 NameEnv) NameEnv {
	env := NameEnv{}
	for _, name := range env1 {
		err := env.Insert(name)
		if err != nil {
			a.panicf(name.Ident, err.Error())
		}
	}
	for _, name := range env2 {
		err := env.Insert(name)
		if err != nil {
			a.panicf(name.Ident, err.Error())
		}
	}
	return env
}
