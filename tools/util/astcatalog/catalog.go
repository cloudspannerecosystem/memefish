package astcatalog

import (
	"go/token"
)

type Catalog map[NodeStructType]*NodeDef

// NodeDef represents a node definition.
type NodeDef struct {
	SourcePos  token.Pos
	Name       string
	Doc        string
	Tmpl       string
	Pos, End   string
	Fields     []*FieldDef
	Implements []NodeInterfaceType
}

type FieldDef struct {
	Name    string
	Type    FieldType
	Comment string
}

type FieldType interface{}

type SliceType struct {
	Type FieldType
}

type PointerType struct {
	Type FieldType
}

type NodeStructType string

type NodeInterfaceType string

// PrimitiveType is a type name which is neither a node pointer nor a node interface.
type PrimitiveType string

const (
	BoolType       PrimitiveType = "bool"
	IntType        PrimitiveType = "int"
	StringType     PrimitiveType = "string"
	TokenPosType   PrimitiveType = "token.Pos"
	TokenTokenType PrimitiveType = "token.Token"
)
