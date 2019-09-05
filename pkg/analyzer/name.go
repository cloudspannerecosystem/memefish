package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type PathName struct {
	IsPath          bool
	Name            string
	ImplicitAliasID int
}

type NameScope struct {
	Next        *NameScope
	List        *NameList
	TableNames  map[string]*TableName
	ColumnNames map[PathName]*ColumnName
}

type NameList struct {
	Columns []*ColumnName
}

func (n *NameList) concat(other *NameList) {
	n.Columns = append(n.Columns, other.Columns...)
}

func (n *NameList) appendNodeToColumns(s parser.SelectItem) *NameList {
	list := &NameList{
		Columns: make([]*ColumnName, len(n.Columns)),
	}
	for i, c := range n.Columns {
		list.Columns[i] = c.appendNode(s)
	}
	return list
}

type ColumnName struct {
	Path    PathName
	Invalid bool
	Type    Type
	Nodes   []parser.SelectItem
}

func (c *ColumnName) merge(d *ColumnName) bool {
	if c.Type == nil {
		c.Type = d.Type
		return true
	}
	if d.Type == nil {
		return true
	}
	return c.Type.EqualTo(d.Type)
}

func (c *ColumnName) appendNode(s parser.SelectItem) *ColumnName {
	return &ColumnName{
		Path:    c.Path,
		Invalid: c.Invalid,
		Type:    c.Type,
		Nodes:   append(append([]parser.SelectItem{}, c.Nodes...), s),
	}
}

type TableName struct {
	Name string
	Type Type
}
