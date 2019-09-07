package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type SelectList []*Reference

func (list SelectList) toType() Type {
	fields := make([]*StructField, len(list))
	for i, r := range list {
		fields[i] = &StructField{
			Name: r.Name,
			Type: r.Type,
		}
	}
	return &StructType{Fields: fields}
}

func (list SelectList) LookupRef(name string) *Reference {
	for _, ref := range list {
		if strings.EqualFold(ref.Name, name) {
			return ref
		}
	}
	return nil
}

func (list SelectList) toNameScope(next *NameScope) *NameScope {
	scope := &NameScope{
		List: list,
		Next: next,
	}

	scope.Env = newNameEnv()
	for _, ref := range list {
		scope.Env.insertRef(NewPath(ref.Name), ref)
	}

	return scope
}

func (list SelectList) deriveSimple(n parser.Node) SelectList {
	newList := make(SelectList, len(list))
	for i, r := range list {
		newList[i] = r.deriveSimple(n)
	}
	return newList
}

func (list SelectList) derive(n parser.Node, kind ReferenceKind) SelectList {
	newList := make(SelectList, len(list))
	for i, r := range list {
		newList[i] = r.derive(n, kind)
	}
	return newList
}
