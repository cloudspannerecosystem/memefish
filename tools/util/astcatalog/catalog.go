package astcatalog

import (
	"go/token"
)

// Catalog is a catalog of AST types.
type Catalog struct {
	Structs    map[NodeStructType]*NodeStructDef
	Interfaces map[NodeInterfaceType]*NodeInterfaceDef
	Consts     map[ConstType]*ConstDef
}

// NodeStructDef is a definition of node structs in ast/ast.go.
type NodeStructDef struct {
	SourcePos  token.Pos
	Name       string
	Doc        string
	Tmpl       string
	Pos, End   string
	Fields     []*FieldDef
	Implements []NodeInterfaceType
}

// FieldDef is a field definition of node structs in ast/ast.go.
type FieldDef struct {
	Name    string
	Type    Type
	Comment string
}

// NodeInterfaceDef is a definition of node interfaces in ast/ast.go.
type NodeInterfaceDef struct {
	SourcePos   token.Pos
	Name        string
	Implemented []NodeStructType
}

// ConstDef is a definition of const types in ast/ast_const.go
type ConstDef struct {
	SourcePos token.Pos
	Name      string
	Values    []*ConstValueDef
}

// ConstValueDef is a value definition of const types in ast/ast_const.go.
type ConstValueDef struct {
	Name  string
	Value string
}

// Type represents types used in Catalog.
type Type interface {
	isType()
}

func (SliceType) isType()         {}
func (PointerType) isType()       {}
func (NodeStructType) isType()    {}
func (NodeInterfaceType) isType() {}
func (PrimitiveType) isType()     {}
func (ConstType) isType()         {}

// SliceType is a slice type.
type SliceType struct {
	Type Type
}

// PointerType is a pointer type.
type PointerType struct {
	Type Type
}

// NodeStructType is a type name of node structs defined in ast/ast.go.
type NodeStructType string

// NodeInterfaceType is a type name of node interfaces defined in ast/ast.go.
type NodeInterfaceType string

// ConstType is a type name of const types defined in ast/ast_const.go.
type ConstType string

// PrimitiveType is a type name which is neither a node pointer, a node interface, nor a const types.
type PrimitiveType string

// PrimitiveType values.
const (
	BoolType       PrimitiveType = "bool"
	IntType        PrimitiveType = "int"
	StringType     PrimitiveType = "string"
	TokenPosType   PrimitiveType = "token.Pos"
	TokenTokenType PrimitiveType = "token.Token"
)
