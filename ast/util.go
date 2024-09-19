package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
	"strings"
)

// sqlOpt outputs:
//
//	when node != nil: left + node.SQL() + right
//	else            : empty string
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

func lastElem[T any](s []T) T {
	return s[len(s)-1]
}

func firstValidEnd(ns ...Node) token.Pos {
	for _, n := range ns {
		if n != nil && n.End() != token.InvalidPos {
			return n.End()
		}
	}
	return token.InvalidPos
}

func firstPos[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return s[0].Pos()
}

func lastEnd[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return lastElem(s).End()
}

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
