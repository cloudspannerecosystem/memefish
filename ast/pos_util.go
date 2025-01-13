package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
)

// This file contains utility types/functions for pos.go.

func nodePos(n Node) token.Pos {
	if n == nil {
		return token.InvalidPos
	}
	return n.Pos()
}

func nodeEnd(n Node) token.Pos {
	if n == nil {
		return token.InvalidPos
	}
	return n.End()
}

func posChoice(ps ...token.Pos) token.Pos {
	for _, p := range ps {
		if !p.Invalid() {
			return p
		}
	}
	return token.InvalidPos
}

func posAdd(p token.Pos, x int) token.Pos {
	if p.Invalid() {
		return token.InvalidPos
	}

	return token.Pos(int(p) + x)
}

func nodeChoice(ns ...Node) Node {
	for _, n := range ns {
		if n != nil {
			return n
		}
	}
	return nil
}

func nodeSliceIndex[T Node](ns []T, i int) Node {
	if len(ns) == 0 {
		return nil
	}

	return ns[i]
}

func nodeSliceLast[T Node](ns []T) Node {
	if len(ns) == 0 {
		return nil
	}

	return ns[len(ns)-1]
}

func ifThenElse(c bool, t int, e int) int {
	if c {
		return t
	} else {
		return e
	}
}
