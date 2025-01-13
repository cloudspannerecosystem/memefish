package ast

import (
	"iter"
)

//go:generate go run ../tools/gen-ast-walk/main.go -astfile ast.go -constfile ast_const.go -outfile walk_internal.go

// Visitor is an interface for visiting AST nodes.
// If the result of Visit is nil, the node will not be traversed.
type Visitor interface {
	Visit(node Node) Visitor
	VisitMany(nodes []Node) Visitor
	Field(name string) Visitor
	Index(index int) Visitor
}

type stackItem struct {
	node    Node
	nodes   []Node
	visitor Visitor
}

// Walk traverses an AST node in depth-first order.
func Walk(node Node, v Visitor) {
	stack := []*stackItem{{node: node, visitor: v}}
	walkMain(stack)
}

// Walk traverses AST nodes in depth-first order.
func WalkMany[T Node](nodes []T, v Visitor) {
	stack := []*stackItem{{nodes: wrapNodes(nodes), visitor: v}}
	walkMain(stack)
}

func walkMain(stack []*stackItem) {
	for len(stack) > 0 {
		last := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if last.node == nil && last.nodes == nil {
			continue
		}

		if last.nodes != nil {
			v := last.visitor.VisitMany(last.nodes)
			for i := len(last.nodes) - 1; i >= 0; i-- {
				stack = append(stack, &stackItem{node: last.nodes[i], visitor: v.Index(i)})
			}
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

func (f inspector) VisitMany(nodes []Node) Visitor {
	return f
}

func (f inspector) Field(name string) Visitor {
	return f
}

func (f inspector) Index(index int) Visitor {
	return f
}

// Inspect traverses an AST node in depth-first order and calls f for each node.
func Inspect(node Node, f func(Node) bool) {
	Walk(node, inspector(f))
}

// InspectMany traverses AST nodes in depth-first order and calls f for each node.
func InspectMany[T Node](nodes []T, f func(Node) bool) {
	WalkMany(nodes, inspector(f))
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

// PreorderMany returns an iterator that traverses AST nodes in depth-first preorder.
func PreorderMany[T Node](nodes []T) iter.Seq[Node] {
	return func(yield func(Node) bool) {
		ok := true
		InspectMany(nodes, func(n Node) bool {
			ok = ok && yield(n)
			return ok
		})
	}
}
