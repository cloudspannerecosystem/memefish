package analyzer

type NameScope struct {
	List NameList
	Env  NameEnv
	Next *NameScope

	context *GroupByContext
}

func (scope *NameScope) Lookup(target string) (*Name, *GroupByContext) {
	name := scope.Env.Lookup(target)
	if name != nil {
		return name, scope.context
	}
	if scope.Next != nil {
		return scope.Next.Lookup(target)
	}
	return nil, nil
}
