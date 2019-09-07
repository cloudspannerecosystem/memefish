package analyzer

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

func (t *TableScope) toNameScope(next *NameScope) *NameScope {
	scope := t.List.toNameScope(next)

	for name, ref := range t.Refs {
		scope.Env.insertRef(NewPath(name), ref)
	}

	return scope
}
