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
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cloudspannerecosystem/memefish/token"
)

// node represents the subset interface of ast.Node.
//
// Package poslang is used by tools/gen-ast-pos and it is used for generating ast/pos.go.
// To avoid a mutual dependency, we introduce this type here.
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
	// Unparse returns the string representation of this expression.
	Unparse() string
}

// Typed expression types has a method Eval*.
// This method evaluates the expression on the given any value.

// PosExpr is a POS expression typed with token.Pos.
type PosExpr interface {
	Expr
	EvalPos(x any) token.Pos
	PosExprToGo(x string) string
}

// NodeExpr is a POS expression typed with ast.Node.
type NodeExpr interface {
	Expr
	EvalNode(x any) node
	NodeExprToGo(x string) string
}

// NodeSliceExpr is a POS expression typed with []ast.Node.
type NodeSliceExpr interface {
	Expr
	EvalNodeSlice(x any) []node
	NodeSliceExprToGo(x string) string
}

// IntExpr is a POS expression typed with int.
type IntExpr interface {
	Expr
	EvalInt(x any) int
	IntExprToGo(x string) string
}

// BoolExpr is a POS expression typed with bool.
type BoolExpr interface {
	Expr
	EvalBool(x any) bool
	BoolExprToGo(x string) string
}

// StringExpr is a POS expression typed with string.
type StringExpr interface {
	Expr
	EvalString(x any) string
	StringExprToGo(x string) string
}

// ===========================================
//
// AST Node Definitions
//
// ===========================================

// fieldOf returns the field value of x with some error checks.
func fieldOf(x any, name string) reflect.Value {
	value := reflect.ValueOf(x)
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !(value.IsValid() && value.Kind() == reflect.Struct) {
		panic("invalid context")
	}

	field := value.FieldByName(name)
	if !field.IsValid() {
		panic("invalid field: " + name)
	}

	return field
}

// Var represents an untyped variable in a POS expression.
//
// This type is determined by a context, i.e., on parsing.
type Var struct {
	Name string
}

func (v *Var) Unparse() string {
	return v.Name
}

func (v *Var) EvalPos(x any) token.Pos {
	field := fieldOf(x, v.Name)
	if !field.CanInt() {
		panic("expect int, but " + field.Kind().String())
	}

	return token.Pos(field.Int())
}

func (v *Var) PosExprToGo(x string) string {
	return fmt.Sprintf("%s.%s", x, v.Name)
}

func (v *Var) EvalNode(x any) node {
	field := fieldOf(x, v.Name)
	if !field.CanInterface() {
		panic("expect interface, but " + field.Kind().String())
	}

	if field.IsNil() {
		return nil
	}

	node, ok := field.Interface().(node)
	if !ok {
		panic("cannot convert the value to an AST node")
	}

	return node
}

func (v *Var) NodeExprToGo(x string) string {
	return fmt.Sprintf("wrapNode(%s.%s)", x, v.Name)
}

func (v *Var) EvalNodeSlice(x any) []node {
	field := fieldOf(x, v.Name)
	if field.Kind() != reflect.Slice {
		panic("expect slice, but " + field.Kind().String())
	}

	if field.IsNil() || field.Len() == 0 {
		return nil
	}

	n := field.Len()
	nodes := make([]node, n)
	for i := 0; i < n; i++ {
		item := field.Index(i)
		if !(item.IsValid() && item.CanInterface()) {
			panic("expect interface, but " + item.Kind().String())
		}

		var ok bool
		nodes[i], ok = item.Interface().(node)
		if !ok {
			panic("cannot convert the value to an AST node")
		}
	}

	return nodes
}

func (v *Var) NodeSliceExprToGo(x string) string {
	return fmt.Sprintf("%s.%s", x, v.Name)
}

func (v *Var) EvalBool(x any) bool {
	field := fieldOf(x, v.Name)
	if field.Kind() != reflect.Bool {
		panic("expect bool, but " + field.Kind().String())
	}

	return field.Bool()
}

func (v *Var) BoolExprToGo(x string) string {
	return fmt.Sprintf("%s.%s", x, v.Name)
}

func (v *Var) EvalString(x any) string {
	field := fieldOf(x, v.Name)
	if field.Kind() != reflect.String {
		panic("expect string, but " + field.Kind().String())
	}

	return field.String()
}

func (v *Var) StringExprToGo(x string) string {
	return fmt.Sprintf("%s.%s", x, v.Name)
}

// NodePos represents a "NodeExpr.pos" expression.
type NodePos struct {
	Expr NodeExpr
}

func (p *NodePos) Unparse() string {
	return fmt.Sprintf("%s.pos", p.Expr.Unparse())
}

func (p *NodePos) EvalPos(x any) token.Pos {
	node := p.Expr.EvalNode(x)
	if node == nil {
		return token.InvalidPos
	}
	return node.Pos()
}

func (p *NodePos) PosExprToGo(x string) string {
	return fmt.Sprintf("nodePos(%s)", p.Expr.NodeExprToGo(x))
}

// NodeEnd represents a "NodeExpr.end" expression.
type NodeEnd struct {
	Expr NodeExpr
}

func (e *NodeEnd) Unparse() string {
	return fmt.Sprintf("%s.end", e.Expr.Unparse())
}

func (e *NodeEnd) EvalPos(x any) token.Pos {
	node := e.Expr.EvalNode(x)
	if node == nil {
		return token.InvalidPos
	}
	return node.End()
}

func (e *NodeEnd) PosExprToGo(x string) string {
	return fmt.Sprintf("nodeEnd(%s)", e.Expr.NodeExprToGo(x))
}

// PosChoice represents a "PosExpr1 || PosExpr2 || ..." expression.
type PosChoice struct {
	Exprs []PosExpr
}

func (c *PosChoice) Unparse() string {
	ss := make([]string, 0, len(c.Exprs))
	for _, e := range c.Exprs {
		ss = append(ss, e.Unparse())
	}
	return strings.Join(ss, " || ")
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

func (c *PosChoice) PosExprToGo(x string) string {
	ss := make([]string, 0, len(c.Exprs))
	for _, e := range c.Exprs {
		ss = append(ss, e.PosExprToGo(x))
	}
	return fmt.Sprintf("posChoice(%s)", strings.Join(ss, ", "))
}

// PosAdd represents a "PosExpr + IntExpr" expression.
type PosAdd struct {
	Expr  PosExpr
	Value IntExpr
}

func (a *PosAdd) Unparse() string {
	return fmt.Sprintf("%s + %s", a.Expr.Unparse(), a.Value.Unparse())
}

func (a *PosAdd) EvalPos(x any) token.Pos {
	pos := a.Expr.EvalPos(x)
	if pos.Invalid() {
		return token.InvalidPos
	}

	value := a.Value.EvalInt(x)
	return token.Pos(int(pos) + value)
}

func (a *PosAdd) PosExprToGo(x string) string {
	return fmt.Sprintf("posAdd(%s, %s)", a.Expr.PosExprToGo(x), a.Value.IntExprToGo(x))
}

// NodeChoice represents a "(NodeExpr1 ?? NodeExpr2 ?? ...)" expression.
type NodeChoice struct {
	Exprs []NodeExpr
}

func (c *NodeChoice) Unparse() string {
	ss := make([]string, 0, len(c.Exprs))
	for _, e := range c.Exprs {
		ss = append(ss, e.Unparse())
	}
	return fmt.Sprintf("(%s)", strings.Join(ss, " ?? "))
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

func (c *NodeChoice) NodeExprToGo(x string) string {
	ss := make([]string, 0, len(c.Exprs))
	for _, e := range c.Exprs {
		ss = append(ss, e.NodeExprToGo(x))
	}
	return fmt.Sprintf("nodeChoice(%s)", strings.Join(ss, ", "))
}

// NodeSliceIndex represents a "NodeSliceExpr[IntExpr]" expression.
type NodeSliceIndex struct {
	Expr  NodeSliceExpr
	Index IntExpr
}

func (i *NodeSliceIndex) Unparse() string {
	return fmt.Sprintf("%s[%s]", i.Expr.Unparse(), i.Index.Unparse())
}

func (i *NodeSliceIndex) EvalNode(x any) node {
	nodes := i.Expr.EvalNodeSlice(x)
	if len(nodes) == 0 {
		return nil
	}

	index := i.Index.EvalInt(x)
	return nodes[index]
}

func (i *NodeSliceIndex) NodeExprToGo(x string) string {
	return fmt.Sprintf("nodeSliceIndex(%s, %s)", i.Expr.NodeSliceExprToGo(x), i.Index.IntExprToGo(x))
}

// NodeSliceLast represents a "NodeSliceExpr[$]" expression.
type NodeSliceLast struct {
	Expr NodeSliceExpr
}

func (l *NodeSliceLast) Unparse() string {
	return fmt.Sprintf("%s[$]", l.Expr.Unparse())
}

func (l *NodeSliceLast) EvalNode(x any) node {
	nodes := l.Expr.EvalNodeSlice(x)
	if len(nodes) == 0 {
		return nil
	}

	return nodes[len(nodes)-1]
}

func (l *NodeSliceLast) NodeExprToGo(x string) string {
	return fmt.Sprintf("nodeSliceLast(%s)", l.Expr.NodeSliceExprToGo(x))
}

// Len represents a "len(StringExpr)" expression.
type Len struct {
	Expr StringExpr
}

func (l *Len) Unparse() string {
	return fmt.Sprintf("len(%s)", l.Expr.Unparse())
}

func (l *Len) EvalInt(x any) int {
	return len(l.Expr.EvalString(x))
}

func (l *Len) IntExprToGo(x string) string {
	return fmt.Sprintf("len(%s)", l.Expr.StringExprToGo(x))
}

// IntLiteral represents an integer literal in a POS expression.
type IntLiteral struct {
	Value int
}

func (i *IntLiteral) Unparse() string {
	return strconv.Itoa(i.Value)
}

func (i *IntLiteral) EvalInt(x any) int {
	return i.Value
}

func (i *IntLiteral) IntExprToGo(x string) string {
	return strconv.Itoa(i.Value)
}

// IfThenElse represents a "(BoolExpr ? IntExpr1 : IntExpr2)" expression.
type IfThenElse struct {
	Cond       BoolExpr
	Then, Else IntExpr
}

func (i *IfThenElse) Unparse() string {
	return fmt.Sprintf("(%s ? %s : %s)", i.Cond.Unparse(), i.Then.Unparse(), i.Else.Unparse())
}

func (i *IfThenElse) EvalInt(x any) int {
	if i.Cond.EvalBool(x) {
		return i.Then.EvalInt(x)
	} else {
		return i.Else.EvalInt(x)
	}
}

func (i *IfThenElse) IntExprToGo(x string) string {
	return fmt.Sprintf("ifThenElse(%s, %s, %s)", i.Cond.BoolExprToGo(x), i.Then.IntExprToGo(x), i.Else.IntExprToGo(x))
}
