package analyzer

import (
	"github.com/MakeNowJust/memefish/pkg/parser"
)

type Path []string

func NewSingletonPath(name string) Path {
	if name == "" {
		return Path{}
	}
	return Path{name}
}

func (p Path) IsAnonymous() bool {
	return len(p) == 0
}

type NameScope struct {
	List NameList
	Env  *NameEnv
}

type NameEnv struct {
	Refs     map[string]*Reference
	Children map[string]*NameEnv
}

func NewNameEnv() *NameEnv {
	return &NameEnv{
		Refs:     make(map[string]*Reference),
		Children: make(map[string]*NameEnv),
	}
}

func (n *NameEnv) HasRef(p Path) bool {
	r, _ := n.LookupRef(p)
	return r != nil
}

func (n *NameEnv) LookupRef(p Path) (*Reference, Path) {
	if len(p) == 0 {
		return nil, nil
	}

	if child, ok := n.Children[p[0]]; ok {
		r, rp := child.LookupRef(p[1:])
		if r != nil {
			return r, rp
		}
	}

	if r, ok := n.Refs[p[0]]; ok {
		return r, p[1:]
	}

	return nil, nil
}

func (n *NameEnv) InsertRef(p Path, r *Reference) bool {
	if r.Kind == InvalidRefAmbiguous {
		panic("BUG: insert ambiguous reference")
	}

	r2, rp := n.LookupRef(p)

	if r2 != nil && len(rp) == 0 {
		return n.tryInsertConflictRef(p, r, r2)
	}

	n.insertRef(p, r)
	return true
}

func (n *NameEnv) tryInsertConflictRef(p Path, r1, r2 *Reference) bool {
	switch r2.Kind {
	case TableRef:
		switch r1.Kind {
		case TableRef:
			return false
		case ColumnRef:
			return true
		case InvalidRefAggregate:
			n.insertRef(p, r1)
			return true
		}
	case ColumnRef:
		switch r1.Kind {
		case TableRef:
			n.insertRef(p, r1)
			return true
		case ColumnRef:
			n.insertRef(p, &Reference{
				Kind:   InvalidRefAmbiguous,
				Name:   r2.Name,
				Origin: []*Reference{r2, r1},
			})
			return true
		case InvalidRefAggregate:
			n.insertRef(p, r1)
			return true
		}
	case InvalidRefAmbiguous:
		r2.Origin = append(r2.Origin, r1)
		return true
	}

	panic("BUG: unexpected reference conflict")
}

func (n *NameEnv) insertRef(p Path, r *Reference) {
	if len(p) == 0 {
		panic("BUG: invalid path")
	}

	if len(p) == 1 {
		n.Refs[p[0]] = r
		return
	}

	child, ok := n.Children[p[0]]
	if !ok {
		child = NewNameEnv()
		n.Children[p[0]] = child
	}

	child.insertRef(p[1:], r)
}

func (n *NameEnv) InsertAllNameList(list NameList) {
	for _, r := range list {
		p := NewSingletonPath(r.Name)
		if !p.IsAnonymous() {
			if !n.InsertRef(p, r) {
				panic("BUG: unexpected failed on NameList insertion")
			}
		}
	}
}

type NameList []*Reference

func (list NameList) deriveSimple(n parser.Node) NameList {
	newList := make(NameList, len(list))
	for i, r := range list {
		newList[i] = r.deriveSimple(n)
	}
	return newList
}

func (list NameList) derive(n parser.Node, kind ReferenceKind) NameList {
	newList := make(NameList, len(list))
	for i, r := range list {
		newList[i] = r.derive(n, kind)
	}
	return newList
}
