package poslang

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Parse parses source as a POS expression.
func Parse(source string) (expr PosExpr, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				expr = nil
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p := &parser{source: []rune(source), offset: 0}

	p.skip()
	expr = p.parsePosChoice()
	if !p.eof() {
		panic(fmt.Errorf("unexpected character: '%c'", p.current()))
	}

	return
}

type parser struct {
	source []rune
	offset int
}

func (p *parser) parsePosChoice() PosExpr {
	expr := p.parsePosExpr()

	if !p.consume("||") {
		return expr
	}

	exprs := []PosExpr{expr, p.parsePosExpr()}
	for p.consume("||") {
		exprs = append(exprs, p.parsePosExpr())
	}

	return &PosChoice{Exprs: exprs}
}

func (p *parser) parsePosExpr() PosExpr {
	expr := p.parsePosAtom()

	for p.consume("+") {
		value := p.parseIntAtom()
		expr = &PosAdd{Expr: expr, Value: value}
	}

	return expr
}

func (p *parser) parsePosAtom() PosExpr {
	expr := p.parseNodeExpr()

	if p.consume(".") {
		x := p.parseVar()
		if x.Name == "pos" {
			return &NodePos{Expr: expr}
		} else if x.Name == "end" {
			return &NodeEnd{Expr: expr}
		}
		panic(errors.New(`expected: "pos", "end"`))
	}

	if v, ok := expr.(*Var); ok {
		return v
	}

	panic(fmt.Errorf("invalid type expression: %s", expr.PosExpr()))
}

func (p *parser) parseNodeExpr() NodeExpr {
	if p.consume("(") {
		exprs := []NodeExpr{p.parseNodeAtom()}
		for p.consume("??") {
			exprs = append(exprs, p.parseNodeAtom())
		}
		p.expect(")")
		return &NodeChoice{Exprs: exprs}
	}

	return p.parseNodeAtom()
}

func (p *parser) parseNodeAtom() NodeExpr {
	x := p.parseVar()

	if !p.consume("[") {
		return x
	}

	if p.consume("$") {
		p.expect("]")
		return &NodeSliceLast{Expr: x}
	}

	index := p.parseIntAtom()
	p.expect("]")

	return &NodeSliceIndex{Expr: x, Index: index}
}

func (p *parser) parseIntAtom() IntExpr {
	if p.consume("(") {
		c := p.parseVar()

		// Here's `?` does not conflict with `??` in `parseNodeExpr`
		// because `parseIntAtom` is not called from problematic positions for now.
		p.expect("?")

		t := p.parseIntAtom()
		p.expect(":")
		e := p.parseIntAtom()
		p.expect(")")

		return &IfThenElse{Cond: c, Then: t, Else: e}
	}

	n := p.consumeIntLiteral()
	if n != nil {
		return n
	}

	x := p.consumeVar()
	if x == nil || x.Name != "len" {
		panic(errors.New(`expected: "(", <int>, "var"`))
	}

	p.expect("(")
	s := p.parseVar()
	p.expect(")")

	return &Len{Expr: s}
}

func (p *parser) parseVar() *Var {
	x := p.consumeVar()
	if x == nil {
		panic(errors.New("expected: var"))
	}
	return x
}

func (p *parser) consumeVar() *Var {
	if p.eof() || !unicode.IsLetter(p.current()) {
		return nil
	}

	start := p.offset
	for !p.eof() {
		c := p.current()
		if !(unicode.IsLetter(c) || c == '_' || unicode.IsDigit(c)) {
			break
		}
		p.next()
	}

	name := string(p.source[start:p.offset])
	p.skip()

	return &Var{Name: name}
}

func (p *parser) consumeIntLiteral() *IntLiteral {
	if p.eof() || !unicode.IsDigit(p.current()) {
		return nil
	}

	start := p.offset
	for !p.eof() {
		c := p.current()
		if !unicode.IsDigit(c) {
			break
		}
		p.next()
	}

	str := string(p.source[start:p.offset])
	value, err := strconv.Atoi(str)
	if err != nil {
		panic(errors.New("invalid integer"))
	}

	p.skip()

	return &IntLiteral{Value: value}
}

func (p *parser) consume(s string) bool {
	if p.offset+len(s) > len(p.source) {
		return false
	}

	part := string(p.source[p.offset : p.offset+len(s)])
	if part != s {
		return false
	}

	p.offset += len(s)
	p.skip()

	return true
}

func (p *parser) expect(s string) {
	if p.consume(s) {
		return
	}

	panic(fmt.Errorf("expected: \"%s\"", s))
}

func (p *parser) skip() {
	for !p.eof() && unicode.IsSpace(p.current()) {
		p.next()
	}
}

func (p *parser) current() rune {
	return p.source[p.offset]
}

func (p *parser) next() {
	p.offset += 1
}

func (p *parser) eof() bool {
	return len(p.source) <= p.offset
}
