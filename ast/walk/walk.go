package walk

import (
	"iter"

	"github.com/cloudspannerecosystem/memefish/ast"
)

//go:generate go run ../../tools/gen-ast-walk/main.go -astfile ../ast.go -constfile ../ast_const.go -outfile walk_internal.go

// Visitor is an interface for visiting AST nodes.
// If the result of Visit is nil, the node will not be traversed.
type Visitor interface {
	Visit(node ast.Node) Visitor
}

type stackItem struct {
	node    ast.Node
	visitor Visitor
}

// Walk traverses an AST in depth-first order.
func Walk(node ast.Node, v Visitor) {
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

type inspector func(ast.Node) bool

func (f inspector) Visit(node ast.Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order and calls f for each node.
func Inspect(node ast.Node, f func(ast.Node) bool) {
	Walk(node, inspector(f))
}

// Preorder returns an iterator that traverses an AST in depth-first preorder.
func Preorder(node ast.Node) iter.Seq[ast.Node] {
	return func(yield func(ast.Node) bool) {
		ok := true
		Inspect(node, func(n ast.Node) bool {
			ok = ok && yield(n)
			return ok
		})
	}
}
