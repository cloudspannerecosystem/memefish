package analyzer

import (
	"strings"
)

type TableScope struct {
	List SelectList
	Refs map[string]*Reference
}

func newTableScope() *TableScope {
	return &TableScope{
		List: nil,
		Refs: make(map[string]*Reference),
	}
}

func (t *TableScope) LookupRef(name string) *Reference {
	if r, ok := t.Refs[strings.ToUpper(name)]; ok {
		return r
	}
	return t.List.LookupRef(name)
}

func (t *TableScope) toNameScope(next *NameScope) *NameScope {
	scope := t.List.toNameScope(next)

	for name, ref := range t.Refs {
		scope.Env.insertRef(NewPath(name), ref)
	}

	return scope
}
