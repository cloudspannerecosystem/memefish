package analyzer

import (
	"fmt"

	"github.com/cloudspannerecosystem/memefish/pkg/ast"
	"github.com/cloudspannerecosystem/memefish/pkg/char"
)

type NameList []*Name

func (list NameList) Lookup(target string) *Name {
	var founds []*Name
	for _, name := range list {
		if char.EqualFold(name.Text, target) {
			founds = append(founds, name)
		}
	}

	switch len(founds) {
	case 0:
		return nil
	case 1:
		return founds[0]
	default:
		return makeAmbiguousName(target, founds)
	}
}

func (list NameList) ToType() *StructType {
	fields := make([]*StructField, len(list))
	for i, name := range list {
		fields[i] = &StructField{
			Name: name.Text,
			Type: name.Type,
		}
	}
	return &StructType{Fields: fields}
}

func (list NameList) toNameEnv() NameEnv {
	env := NameEnv{}
	for _, name := range list {
		err := env.Insert(name)
		if err != nil {
			panic(fmt.Sprintf("BUG: unexpected error: %v", err))
		}
	}
	return env
}

func (list NameList) toTableInfo() *TableInfo {
	return &TableInfo{
		List: list,
		Env:  list.toNameEnv(),
	}
}

func (list NameList) toNameScope(next *NameScope) *NameScope {
	return &NameScope{
		List: list,
		Env:  list.toNameEnv(),
		Next: next,
	}
}

// for ast.DotStar
func makeNameListFromType(t Type, node ast.Node) NameList {
	parent := &Name{
		Kind: ColumnName,
		Type: t,
		Node: node,
	}
	return NameList(parent.Children())
}
