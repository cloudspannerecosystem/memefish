package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type GroupByContext struct {
	Lists      map[parser.SelectItem]NameList
	ValidNames map[*Name]struct{}
	ValidExprs []parser.Expr
}

func (gbc *GroupByContext) AddValidName(name *Name) {
	name = name.Deref()
	gbc.ValidNames[name] = struct{}{}
}

func (gbc *GroupByContext) IsValidName(name *Name) bool {
	name = name.Deref()
	if _, ok := gbc.ValidNames[name]; ok {
		return true
	}
	return false
}
