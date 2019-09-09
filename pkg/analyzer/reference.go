package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Reference struct {
	Kind ReferenceKind

	Name string // Name == "" means this reference is anonymous
	Type Type

	Parent   *Reference
	Children map[string]*Reference

	TableSchema  *TableSchema
	ColumnSchema *ColumnSchema
	Node         parser.Node
	Ident        *parser.Ident // Ident corresponding to Name
	Origin       []*Reference
}

func (r *Reference) GetNode(fallback parser.Node) parser.Node {
	if r.Node != nil {
		return r.Node
	}

	for _, o := range r.Origin {
		n := o.GetNode(nil)
		if n != nil {
			return n
		}
	}

	return fallback
}

func (r *Reference) GetIdent(fallback parser.Node) parser.Node {
	if r.Ident != nil {
		return r.Ident
	}

	for _, o := range r.Origin {
		n := o.GetIdent(nil)
		if n != nil {
			return n
		}
	}

	return fallback
}

func (r *Reference) GetName() string {
	if r.Name == "" {
		return "(unspecified)"
	}
	return r.Name
}

type ReferenceKind int

const (
	_ ReferenceKind = iota
	TableRef
	ColumnRef
)

const (
	_ ReferenceKind = -iota
	InvalidRefAmbiguous
	InvalidRefAggregate
)

func (k ReferenceKind) Invalid() bool {
	return k < 0
}
