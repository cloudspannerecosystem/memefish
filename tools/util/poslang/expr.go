// Package poslang provides an implementation of the POS expression.
//
// POS is a small language to write a specification of Pos()/End() methods of ast.Node types.
// The syntax of POS is described in ast package documentation.
// See https://pkg.go.dev/github.com/cloudspannerecosystem/memefish/ast.
//
// This package provides:
//
//   - types for the POS expression
//   - a parser for the POS expression
//   - an interpreter of the POS expression
package poslang

import (
	"reflect"
	"strconv"

	"github.com/cloudspannerecosystem/memefish/token"
)

// node represents the subset interface of ast.Node.
//
// Package poslang is used by tools/gen-ast-pos and it is used for generating ast/pos.go.
// To avoid mutual dependency, we introduce this type here.
type node interface {
	Pos() token.Pos
	End() token.Pos
}

// ===========================================
//
// Expression Interfaces
//
// ===========================================

// Expr is an interface that represents a POS expression.
//
// This means an untyped expression and a base interface of all typed expressions.
type Expr interface {
	// PosExpr returns the string representation of this expression (i.e. unparse).
	PosExpr() string
}

// Typed expression types has a method EvalXXX.
// This method evaluates the expression on the given any value.

// PosExpr is a POS expression typed with token.Pos.
type PosExpr interface {
	Expr
	EvalPos(x any) token.Pos
}

// NodeExpr is a POS expression typed with ast.Node.
type NodeExpr interface {
	Expr
	EvalNode(x any) node
}

// NodeSliceExpr is a POS expression typed with []ast.Node.
type NodeSliceExpr interface {
	Expr
	EvalNodeSlice(x any) []node
}

// IntExpr is a POS expression typed with int.
type IntExpr interface {
	Expr
	EvalInt(x any) int
}

// BoolExpr is a POS expression typed with bool.
type BoolExpr interface {
	Expr
	EvalBool(x any) bool
}

// StringExpr is a POS expression typed with string.
type StringExpr interface {
	Expr
	EvalString(x any) string
}

// ===========================================
//
// AST Node Definitions
//
// ===========================================

// Var represents an untyped variable in a POS expression.
//
// This type is determined by a context, i.e., on parsing.
type Var struct {
	Name string
}

func (v *Var) PosExpr() string {
	return v.Name
}

func (v *Var) EvalPos(x any) token.Pos {
	value := reflect.ValueOf(x).Elem()
	field := value.FieldByName(v.Name)

	return token.Pos(field.Int())
}

func (v *Var) EvalNode(x any) node {
	value := reflect.ValueOf(x).Elem()
	field := value.FieldByName(v.Name)
	if field.IsNil() {
		return nil
	}

	node, _ := field.Interface().(node)

	return node
}

func (v *Var) EvalNodeSlice(x any) []node {
	value := reflect.ValueOf(x).Elem()
	field := value.FieldByName(v.Name)
	if field.IsNil() {
		return nil
	}

	n := field.Len()
	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		nodes[i], _ = field.Index(i).Interface().(node)
	}

	return nodes
}

func (v *Var) EvalBool(x any) bool {
	value := reflect.ValueOf(x).Elem()
	field := value.FieldByName(v.Name)

	return field.Bool()
}

func (v *Var) EvalString(x any) string {
	value := reflect.ValueOf(x).Elem()
	field := value.FieldByName(v.Name)

	return field.String()
}

// NodePos represents a "NodeExpr.pos" expression.
type NodePos struct {
	Expr NodeExpr
}

func (p *NodePos) PosExpr() string {
	return p.Expr.PosExpr() + ".pos"
}

func (p *NodePos) EvalPos(x any) token.Pos {
	return p.Expr.EvalNode(x).Pos()
}

// NodeEnd represents a "NodeExpr.end" expression.
type NodeEnd struct {
	Expr NodeExpr
}

func (e *NodeEnd) PosExpr() string {
	return e.Expr.PosExpr() + ".end"
}

func (p *NodeEnd) EvalPos(x any) token.Pos {
	node := p.Expr.EvalNode(x)
	if node == nil {
		return token.InvalidPos
	}
	return node.End()
}

// PosChoice represents a "PosExpr1 || PosExpr2 || ..." expression.
type PosChoice struct {
	Exprs []PosExpr
}

func (c *PosChoice) PosExpr() string {
	s := c.Exprs[0].PosExpr()
	for _, e := range c.Exprs[1:] {
		s += " || " + e.PosExpr()
	}
	return s
}

func (c *PosChoice) EvalPos(x any) token.Pos {
	for _, e := range c.Exprs {
		pos := e.EvalPos(x)
		if !pos.Invalid() {
			return pos
		}
	}

	return token.InvalidPos
}

// PosAdd represents a "PosExpr + IntExpr" expression.
type PosAdd struct {
	Expr  PosExpr
	Value IntExpr
}

func (a *PosAdd) PosExpr() string {
	return a.Expr.PosExpr() + " + " + a.Value.PosExpr()
}

func (a *PosAdd) EvalPos(x any) token.Pos {
	pos := a.Expr.EvalPos(x)
	if pos.Invalid() {
		return token.InvalidPos
	}

	value := a.Value.EvalInt(x)
	return token.Pos(int(pos) + value)
}

// NodeChoice represents a "(NodeExpr1 ?? NodeExpr2 ?? ...)" expression.
type NodeChoice struct {
	Exprs []NodeExpr
}

func (c *NodeChoice) PosExpr() string {
	s := "(" + c.Exprs[0].PosExpr()
	for _, e := range c.Exprs[1:] {
		s += " ?? " + e.PosExpr()
	}
	return s + ")"
}

func (c *NodeChoice) EvalNode(x any) node {
	for _, e := range c.Exprs {
		node := e.EvalNode(x)
		if node != nil {
			return node
		}
	}

	return nil
}

// NodeSliceIndex represents a "NodeSliceExpr[IntExpr]" expression.
type NodeSliceIndex struct {
	Expr  NodeSliceExpr
	Index IntExpr
}

func (i *NodeSliceIndex) PosExpr() string {
	return i.Expr.PosExpr() + "[" + i.Index.PosExpr() + "]"
}

func (i *NodeSliceIndex) EvalNode(x any) node {
	nodes := i.Expr.EvalNodeSlice(x)
	if nodes == nil {
		return nil
	}

	index := i.Index.EvalInt(x)
	return nodes[index]
}

// NodeSliceLast represents a "NodeSliceExpr[$]" expression.
type NodeSliceLast struct {
	Expr NodeSliceExpr
}

func (l *NodeSliceLast) PosExpr() string {
	return l.Expr.PosExpr() + "[$]"
}

func (l *NodeSliceLast) EvalNode(x any) node {
	nodes := l.Expr.EvalNodeSlice(x)
	if nodes == nil {
		return nil
	}

	return nodes[len(nodes)-1]
}

// Len represents a "len(StringExpr)" expression.
type Len struct {
	Expr StringExpr
}

func (l *Len) PosExpr() string {
	return "len(" + l.Expr.PosExpr() + ")"
}

func (l *Len) EvalInt(x any) int {
	return len(l.Expr.EvalString(x))
}

// IntLiteral represents an integer literal in a POS expression.
type IntLiteral struct {
	Value int
}

func (i *IntLiteral) PosExpr() string {
	return strconv.Itoa(i.Value)
}

func (i *IntLiteral) EvalInt(x any) int {
	return i.Value
}

// IfThenElse represents a "(BoolExpr ? IntExpr1 : IntExpr2)" expression.
type IfThenElse struct {
	Cond       BoolExpr
	Then, Else IntExpr
}

func (i *IfThenElse) PosExpr() string {
	return "(" + i.Cond.PosExpr() + " ? " + i.Then.PosExpr() + " : " + i.Else.PosExpr() + ")"
}

func (i *IfThenElse) EvalInt(x any) int {
	if i.Cond.EvalBool(x) {
		return i.Then.EvalInt(x)
	} else {
		return i.Else.EvalInt(x)
	}
}
