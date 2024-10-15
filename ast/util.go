package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
	"strings"
)

// Helper functions for SQL()

// sqlOpt outputs:
//
//	when node != nil: left + node.SQL() + right
//	else            : empty string
//
// This function corresponds to sqlOpt in ast.go
func sqlOpt[T interface {
	Node
	comparable
}](left string, node T, right string) string {
	var zero T
	if node == zero {
		return ""
	}
	return left + node.SQL() + right
}

// strOpt outputs:
//
//	when pred == true: s
//	else            : empty string
//
// This function corresponds to {{if pred}}s{{end}} in ast.go
func strOpt(pred bool, s string) string {
	if pred {
		return s
	}
	return ""
}

// sqlJoin outputs joined string of SQL() of all elems by sep.
// This function corresponds to sqlJoin in ast.go
func sqlJoin[T Node](elems []T, sep string) string {
	var b strings.Builder
	for i, r := range elems {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(r.SQL())
	}
	return b.String()
}

// Helper functions for Pos(), End()

// lastElem returns last element of slice s.
// This function corresponds to NodeSliceVar[$] in ast.go.
func lastElem[T any](s []T) T {
	return s[len(s)-1]
}

// firstValidEnd returns the first valid Pos() in argument.
// "valid" means the node is not nil and Pos().Invalid() is not true.
// This function corresponds to "(n0 ?? n1 ?? ...).End()"
func firstValidEnd(ns ...Node) token.Pos {
	for _, n := range ns {
		if n != nil && !n.End().Invalid() {
			return n.End()
		}
	}
	return token.InvalidPos
}

// firstPos returns the Pos() of the first node.
// If argument is an empty slice, this function returns token.InvalidPos.
// This function corresponds to NodeSliceVar[0].pos in ast.go.
func firstPos[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return s[0].Pos()
}

// lastEnd returns the End() of the last node.
// If argument is an empty slice, this function returns token.InvalidPos.
// This function corresponds to NodeSliceVar[$].end in ast.go.
func lastEnd[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return lastElem(s).End()
}
