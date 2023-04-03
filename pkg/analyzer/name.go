package analyzer

import (
	"fmt"

	"github.com/cloudspannerecosystem/memefish/pkg/ast"
	"github.com/cloudspannerecosystem/memefish/pkg/token"
)

// Name represents typed names like table name and column name.
type Name struct {
	Kind NameKind

	Text string
	Type Type

	Parent *Name

	Origin    []*Name
	Ambiguous bool

	Ref *Name // for AliasName

	Node  ast.Node
	Ident *ast.Ident

	TableSchema  *TableSchema
	ColumnSchema *ColumnSchema

	// child names for STRUCT typed name.
	children []*Name
}

type NameKind int

const (
	_ NameKind = iota
	TableName
	ColumnName
	DerivedName
	AliasName
)

func (name *Name) Deref() *Name {
	for name.Kind == AliasName {
		name = name.Ref
	}
	return name
}

func (name *Name) Children() []*Name {
	if name.children != nil {
		return name.children
	}

	name = name.Deref()

	t, ok := TypeCastStruct(name.Type)
	if !ok {
		return nil
	}

	children := make([]*Name, len(t.Fields))
	for i, f := range t.Fields {
		children[i] = &Name{
			Kind:   ColumnName,
			Text:   f.Name,
			Type:   f.Type,
			Parent: name,
		}
	}

	name.children = children
	return children
}

func (name *Name) LookupChild(target string) *Name {
	children := name.Children()
	if children == nil {
		panic(fmt.Sprintf("BUG: cannot lookup child from non-STRUCT name: %#v", name))
	}
	return NameList(children).Lookup(target)
}

func (name *Name) Anonymous() bool {
	return name.Text == ""
}

func (name *Name) Quote() string {
	if name.Anonymous() {
		return "(anonymous)"
	}
	return token.QuoteSQLIdent(name.Text)
}

// for ast.Select with AS STRUCT
func makeNameListColumnName(list NameList, node ast.Node) *Name {
	return &Name{
		Kind:     ColumnName,
		Type:     list.ToType(),
		children: []*Name(list),
	}
}

// for ast.ExprSelectItem
func makeExprColumnName(t Type, expr ast.Expr, node ast.Node, ident *ast.Ident) *Name {
	if ident == nil {
		ident = extractIdentFromExpr(expr)
	}

	var text string
	if ident != nil {
		text = ident.Name
	}

	return &Name{
		Kind:  ColumnName,
		Text:  text,
		Type:  t,
		Node:  node,
		Ident: ident,
	}
}

// for ast.Alias
func makeAliasName(name *Name, node ast.Node, ident *ast.Ident) *Name {
	var text string
	if ident != nil {
		text = ident.Name
	}

	return &Name{
		Kind:  AliasName,
		Text:  text,
		Type:  name.Type,
		Ref:   name,
		Node:  node,
		Ident: ident,
	}
}

// for ast.TableName
func makeTableSchemaName(table *TableSchema, node ast.Node, ident *ast.Ident) *Name {
	text := table.Name
	if ident != nil {
		text = ident.Name
	}

	parent := &Name{
		Kind:        TableName,
		Text:        text,
		Type:        table.ToType(),
		Node:        node,
		Ident:       ident,
		TableSchema: table,
	}

	children := make([]*Name, len(table.Columns))
	for i, c := range table.Columns {
		children[i] = &Name{
			Kind:         ColumnName,
			Text:         c.Name,
			Type:         c.Type,
			Parent:       parent,
			ColumnSchema: c,
		}
	}

	parent.children = children
	return parent
}

// for ast.Unnest
func makeTableName(text string, t Type, node ast.Node, ident *ast.Ident) *Name {
	if ident != nil {
		text = ident.Name
	}

	return &Name{
		Kind:  TableName,
		Text:  text,
		Type:  t,
		Node:  node,
		Ident: ident,
	}
}

// for ast.SubQueryTableExpr
func makeNameListTableName(list NameList, node ast.Node, ident *ast.Ident) *Name {
	var text string
	if ident != nil {
		text = ident.Name
	}

	parent := &Name{
		Kind:  TableName,
		Text:  text,
		Type:  list.ToType(),
		Node:  node,
		Ident: ident,
	}

	children := make([]*Name, len(list))
	for i, name := range list {
		children[i] = &Name{
			Kind:   AliasName,
			Text:   name.Text,
			Type:   name.Type,
			Ref:    name,
			Parent: parent,
		}
	}

	parent.children = children
	return parent
}

// for ast.Join with InnerJoin and LeftOuterJoin
func makeLeftJoinName(left, right *Name) *Name {
	return &Name{
		Kind:   AliasName,
		Text:   left.Text,
		Type:   left.Type,
		Ref:    left,
		Origin: []*Name{left, right},
	}
}

// for ast.Join with RightOuterJoin
func makeRightJoinName(left, right *Name) *Name {
	return &Name{
		Kind:   AliasName,
		Text:   right.Text,
		Type:   right.Type,
		Ref:    right,
		Origin: []*Name{left, right},
	}
}

// for ast.Join with FullOuterJoin
func makeFullJoinName(left, right *Name) (*Name, bool) {
	t, ok := MergeType(left.Type, right.Type)
	if !ok {
		return nil, false
	}

	return &Name{
		Kind:   AliasName,
		Text:   left.Text,
		Type:   t,
		Ref:    left,
		Origin: []*Name{left, right},
	}, true
}

// for ast.CompoundQuery and ast.Join with FullOuterJoin
func makeCompoundQueryResultName(names []*Name, node ast.Node) (*Name, bool) {
	if len(names) <= 1 {
		panic(fmt.Sprintf("BUG: too few names: %#+v", names))
	}

	var t Type
	for _, name := range names {
		tt, ok := MergeType(t, name.Type)
		if !ok {
			return nil, false
		}

		t = tt
	}

	name0 := names[0]
	return &Name{
		Kind:   DerivedName,
		Text:   name0.Text,
		Type:   t,
		Origin: names,
	}, true
}

func makeAmbiguousName(text string, names []*Name) *Name {
	if len(names) <= 1 {
		panic(fmt.Sprintf("BUG: too few names: %#+v", names))
	}

	var origin []*Name

	for _, name := range names {
		if name.Kind == TableName {
			panic(fmt.Sprintf("BUG: invalid name: %#+v", name))
		}

		if name.Ambiguous {
			origin = append(origin, name.Origin...)
		} else {
			origin = append(origin, name)
		}
	}

	return &Name{
		Kind:      DerivedName,
		Text:      text,
		Ambiguous: true,
		Origin:    origin,
	}
}
