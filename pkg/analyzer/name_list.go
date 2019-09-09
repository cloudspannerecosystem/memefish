package analyzer

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type NameList []*Name

func (list NameList) Lookup(target string) *Name {
	var founds []*Name
	for _, child := range list {
		if strings.EqualFold(child.Text, target) {
			founds = append(founds, child)
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

// for parser.StarPath
func makeNameListFromType(t Type, node parser.Node) []*Name {
	parent := &Name{
		Kind: ColumnName,
		Type: t,
		Node: node,
	}
	return parent.Children()
}
