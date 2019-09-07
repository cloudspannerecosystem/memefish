package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Reference struct {
	Kind ReferenceKind

	Name string // Name == "" means this reference is anonymous
	Type Type

	TableSchema  *TableSchema
	ColumnSchema *ColumnSchema
	Node         parser.Node
	Ident        *parser.Ident
	Origin       []*Reference
}

type ReferenceKind int

const (
	_ ReferenceKind = iota
	TableRef
	ColumnRef
	_ ReferenceKind = -iota
	InvalidRefAmbiguous
	InvalidRefAggregate
)

func (k ReferenceKind) Invalid() bool {
	return k < 0
}

func newColumnReference(name string, t Type, n parser.Node) *Reference {
	return &Reference{
		Kind: ColumnRef,

		Name: name,
		Type: t,

		Node: n,
	}
}

// merge merges two references.
//
// Try to merge r.Type and s.Type, then set r.Type as this and append s to r.Origin .
// If type merging is failed, it returns true, otherwise it returns false.
//
// NOTE: it is mutable method.
func (r *Reference) merge(s *Reference) bool {
	t, ok := MergeType(r.Type, s.Type)
	if !ok {
		return false
	}

	r.Type = t
	r.Origin = append(r.Origin, s)
	return true
}

func (r *Reference) deriveSimple(n parser.Node) *Reference {
	return r.derive(n, r.Kind)
}

func (r *Reference) derive(n parser.Node, kind ReferenceKind) *Reference {
	return &Reference{
		Kind: kind,

		Name: r.Name,
		Type: r.Type,

		Node:   n,
		Origin: []*Reference{r},
	}
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
