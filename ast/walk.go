package ast

import (
	"iter"
)

//go:generate go run ../tools/gen-ast-walk/main.go -astfile ast.go -constfile ast_const.go -outfile walk_internal.go

// Visitor is an interface for visiting AST nodes.
// If the result of Visit is nil, the node will not be traversed.
type Visitor interface {
	Visit(node Node) Visitor
}

type stackItem struct {
	node    Node
	visitor Visitor
}

// Walk traverses an AST in depth-first order.
func Walk(node Node, v Visitor) {
	var stack []*stackItem

	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if last.node == nil {
			continue
		}

		v := last.visitor.Visit(last.node)
		if v == nil {
			continue
		}

		stack = walkInternal(last.node, v, stack)
	}
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order and calls f for each node.
func Inspect(node Node, f func(Node) bool) {
	Walk(node, inspector(f))
}

// Preorder returns an iterator that traverses an AST in depth-first preorder.
func Preorder(node Node) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		ok := true
		Inspect(node, func(n Node) bool {
			ok = ok && yield(n)
			return ok
		})
	}
}
