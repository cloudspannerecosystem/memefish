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

func newSingletonNameList(path PathName, t Type, n parser.Node) *NameList {
	return &NameList{
		Columns: []*ColumnName{newColumnName(path, t, n)},
	}
}

func (n *NameList) concat(other *NameList) {
	n.Columns = append(n.Columns, other.Columns...)
}

func (n *NameList) derive(node parser.Node) *NameList {
	columns := make([]*ColumnName, len(n.Columns))
	for i, c := range n.Columns {
		columns[i] = c.derive(node)
	}
	return &NameList{Columns: columns}
}

type ColumnName struct {
	Path    PathName
	Invalid bool
	Type    Type

	Origin []*ColumnName
	Node   parser.Node
	Schema *ColumnSchema
}

func newColumnName(path PathName, t Type, n parser.Node) *ColumnName {
	return &ColumnName{
		Path: path,
		Type: t,

		Node: n,
	}
}

func (c *ColumnName) merge(d *ColumnName) bool {
	if c.Invalid || d.Invalid {
		panic("BUG: merge invalid ColumnName")
	}

	t, ok := MergeType(c.Type, d.Type)
	if !ok {
		return false
	}

	c.Type = t
	c.Origin = append(c.Origin, d)
	return true
}

func (c *ColumnName) derive(n parser.Node) *ColumnName {
	return &ColumnName{
		Path:    c.Path,
		Invalid: c.Invalid,
		Type:    c.Type,

		Origin: []*ColumnName{c},
		Node:   n,
	}
}

type TableName struct {
	Name string
	Type Type

	Node   parser.Node
	Schema *TableSchema
}
