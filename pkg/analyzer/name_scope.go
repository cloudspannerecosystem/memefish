package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Path []string

func NewPath(name string) Path {
	if name == "" {
		return Path{}
	}
	return Path{name}
}

func (p Path) IsAnonymous() bool {
	return len(p) == 0
}

type NameScope struct {
	List SelectList
	Env  *NameEnv
	Next *NameScope
}

func (scope *NameScope) LookupRef(p Path) (*Reference, Path) {
	r, rp := scope.Env.LookupRef(p)
	if r != nil {
		return r, rp
	}
	if scope.Next != nil {
		return scope.Next.LookupRef(p)
	}
	return nil, nil
}

func (scope *NameScope) toAggregateKeyScope(keys *NameEnv) *NameScope {
	list := make(SelectList, len(scope.List))
	for i, ref := range scope.List {
		if keys.HasRef(NewPath(ref.Name)) {
			list[i] = ref.deriveSimple(nil)
		} else {
			list[i] = ref.derive(nil, InvalidRefAggregate)
		}
	}

	env := newNameEnv()
	var walk func(p Path, e *NameEnv, f func(p Path, ref *Reference))
	walk = func(p Path, e *NameEnv, f func(p Path, ref *Reference)) {
		for name, ref := range e.Refs {
			f(append(p, name), ref)
		}
		for name, child := range e.Children {
			walk(append(p, name), child, f)
		}
	}

	walk(Path{}, keys, func(p Path, ref *Reference) {
		env.insertRef(p, ref.deriveSimple(nil))
	})

	walk(Path{}, scope.Env, func(p Path, ref *Reference) {
		env.insertRef(p, ref.derive(nil, InvalidRefAggregate))
	})

	return &NameScope{
		List: list,
		Env:  env,
		Next: scope.Next,
	}
}

type NameEnv struct {
	Refs     map[string]*Reference
	Children map[string]*NameEnv
}

func (env *NameEnv) HasRef(p Path) bool {
	ref, rp := env.LookupRef(p)
	return ref != nil && len(rp) == 0
}

func (env *NameEnv) LookupRef(p Path) (*Reference, Path) {
	if len(p) == 0 {
		return nil, nil
	}

	if child, ok := env.Children[p[0]]; ok {
		r, rp := child.LookupRef(p[1:])
		if r != nil {
			return r, rp
		}
	}

	if r, ok := env.Refs[p[0]]; ok {
		return r, p[1:]
	}

	return nil, nil
}

func newNameEnv() *NameEnv {
	return &NameEnv{
		Refs:     make(map[string]*Reference),
		Children: make(map[string]*NameEnv),
	}
}

func (env *NameEnv) insertRef(p Path, r *Reference) {
	if p.IsAnonymous() {
		return
	}

	oldRef, rp := env.LookupRef(p)

	if oldRef != nil && len(rp) == 0 {
		env.tryInsertConflictRef(p, oldRef, r)
		return
	}

	env.forceInsertRef(p, r)
	return
}

func (env *NameEnv) tryInsertConflictRef(p Path, oldRef, newRef *Reference) {
	// ## When insert reference to env
	//
	//   1. SelectList.toNameScope: insert ColumnRef/InvalidRefAmbiguous
	//   2. TableScope.toNameScope: insert TableRef
	//   3. NameScope.toAggregateKeyScope: insert InvalidRefAggregate
	//
	// ## Table for conflict resolution
	//
	// |   old      |   new      ||  result      |  situtation |  note
	// |------------|------------||--------------|-------------|---------...
	// |  table     |  table     || (unexpected) |  ???        |
	// |  table     |  column    || table        |  (1)        | keep old
	// |  table     |  aggregate || table        |  (3)        | keep old
	// |  table     |  ambiguous || ambiguous    |  (1)        | add old to new.Origin, and insert new
	// |  column    |  table     || table        |  (2)        | keep old
	// |  column    |  column    || ambiguous    |  (1)        | insert a new InsertRefAmbiguous with old and new Origin
	// |  column    |  aggregate || column       |  (3)        | keep old
	// |  column    |  ambiguous || ambiguous    |  (1)        | add old to new.Origin, and insert new
	// |  aggregate |  table     || (unexpected) |  ???        |
	// |  aggregate |  column    || (unexpected) |  ???        |
	// |  aggregate |  aggregate || (unexpected) |  ???        |
	// |  aggregate |  ambiguous || (unexpected) |  ???        |
	// |  ambiguous |  table     || amgiguous    |  (2)        | add new to old.Origin
	// |  ambiguous |  column    || ambiguous    |  (1)        | add new to old.Origin
	// |  ambiguous |  aggregate || aggregate    |  (3)        |
	// |  ambiguous |  ambiguous || ambiguous    |  (1)        | merge old.Origin and new.Origin

	if oldRef.Kind == InvalidRefAggregate || newRef.Kind == TableRef && oldRef.Kind == TableRef {
		panic("BUG: unexpected reference conflict")
	}

	if oldRef.Kind == InvalidRefAmbiguous || newRef.Kind == InvalidRefAggregate {
		env.forceInsertRef(p, newRef)
		return
	}

	if newRef.Kind == InvalidRefAggregate {
		// keep old
		return
	}

	if oldRef.Kind == InvalidRefAmbiguous && newRef.Kind == InvalidRefAmbiguous {
		// merge old.Origin and new.Origin
		oldRef.Origin = append(oldRef.Origin, newRef.Origin...)
		return
	}

	if oldRef.Kind == InvalidRefAmbiguous {
		// add new to old.Origin
		oldRef.Origin = append(oldRef.Origin, newRef)
		return
	}

	if newRef.Kind == InvalidRefAmbiguous {
		// add old to new.Origin, and insert new
		newRef.Origin = append(newRef.Origin, oldRef)
		env.forceInsertRef(p, newRef)
		return
	}

	if oldRef.Kind == ColumnRef || newRef.Kind == TableRef {
		env.forceInsertRef(p, newRef)
		return
	}

	if oldRef.Kind == TableRef || newRef.Kind == ColumnRef {
		// keep old
		return
	}

	if newRef.Kind == ColumnRef && oldRef.Kind == ColumnRef {
		env.forceInsertRef(p, &Reference{
			Kind:   InvalidRefAmbiguous,
			Name:   oldRef.Name,
			Origin: []*Reference{oldRef, newRef},
		})
		return
	}

	if newRef.Kind == ColumnRef && oldRef.Kind == TableRef {
		return
	}

	panic("BUG: unexpected reference conflict")
}

func (env *NameEnv) forceInsertRef(p Path, r *Reference) {
	if len(p) == 1 {
		env.Refs[p[0]] = r
		return
	}

	child, ok := env.Children[p[0]]
	if !ok {
		child = newNameEnv()
		env.Children[p[0]] = child
	}

	child.forceInsertRef(p[1:], r)
}

func (env *NameEnv) derive(n parser.Node, kind ReferenceKind) *NameEnv {
	newEnv := newNameEnv()
	for name, child := range env.Children {
		newEnv.Children[name] = child.derive(n, kind)
	}
	for name, ref := range env.Refs {
		newEnv.Refs[name] = ref.derive(n, kind)
	}
	return newEnv
}

func (env *NameEnv) deriveSimple(n parser.Node) *NameEnv {
	newEnv := newNameEnv()
	for name, child := range env.Children {
		newEnv.Children[name] = child.deriveSimple(n)
	}
	for name, ref := range env.Refs {
		newEnv.Refs[name] = ref.deriveSimple(n)
	}
	return newEnv
}
