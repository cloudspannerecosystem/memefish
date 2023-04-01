package analyzer

import (
	"github.com/cloudspannerecosystem/memefish/pkg/ast"
)

type GroupByContext struct {
	Lists      map[ast.SelectItem]NameList
	ValidNames map[*Name]struct{}
	ValidExprs []ast.Expr
}

func (gbc *GroupByContext) AddValidName(name *Name) {
	name = name.Deref()
	if gbc.ValidNames == nil {
		gbc.ValidNames = make(map[*Name]struct{})
	}
	gbc.ValidNames[name] = struct{}{}
}

func (gbc *GroupByContext) IsValidName(name *Name) bool {
	name = name.Deref()
	if gbc.ValidNames == nil {
		return false
	}
	if _, ok := gbc.ValidNames[name]; ok {
		return true
	}
	return false
}
