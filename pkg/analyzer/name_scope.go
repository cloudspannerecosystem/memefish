package analyzer

type NameScope struct {
	List NameList
	Env  NameEnv
	Next *NameScope

	context *GroupByContext
}

func (scope *NameScope) Lookup(target string) *Name {
	name := scope.Env.Lookup(target)
	if name != nil {
		return name
	}
	if scope.Next != nil {
		return scope.Next.Lookup(target)
	}
	return nil
}
