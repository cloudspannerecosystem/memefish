package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	*Lexer
}

func (p *Parser) ParseQuery() *QueryStatement {
	// TODO: recover
	p.NextToken()
	return p.parseQueryStatement()
}

func (p *Parser) ParseExpr() Expr {
	// TODO: recover
	p.NextToken()
	return p.parseExpr()
}

func (p *Parser) parseQueryStatement() *QueryStatement {
	hint := p.parseHint()
	expr := p.parseQueryExpr()

	return &QueryStatement{
		Hint: hint,
		Expr: expr,
	}
}

func (p *Parser) parseHint() *Hint {
	if p.Token.Kind != "@" {
		return nil
	}

	hint := &Hint{
		pos: p.Token.Pos,
		Map: make(map[string]Expr),
	}

	p.NextToken()
	p.expect("{")

	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == "}" {
			break
		}

		id := p.expect(TokenIdent)
		p.expect("=")
		val := p.parseExpr()
		hint.Map[id.AsString] = val

		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	hint.end = p.expect("}").End

	return hint
}

func (p *Parser) parseQueryExpr() QueryExpr {
	q := p.parseSimpleQueryExpr()

	var prevOp SetOp
	var prevDistinct bool
	hasPrev := false
setOp:
	for {
		var op SetOp
		switch p.Token.Kind {
		case "UNION":
			op = SetOpUnion
		case "INTERSECT":
			op = SetOpIntersect
		case "EXCEPT":
			op = SetOpExcept
		default:
			break setOp
		}
		p.NextToken()
		opTok := p.Token
		var distinct bool
		switch p.Token.Kind {
		case "ALL":
			distinct = false
		case "DISTINCT":
			distinct = true
		default:
			p.panicfAtToken(&p.Token, "expected token: ALL, DISTINCT, but: %s", p.Token.Kind)
		}
		p.NextToken()
		if hasPrev {
			if !(prevOp == op && prevDistinct == distinct) {
				p.panicfAtToken(&opTok, "all set operator at the same level must be the same, or wrap (...)")
			}
		} else {
			hasPrev = true
			prevOp = op
			prevDistinct = distinct
		}
		e := p.parseSimpleQueryExpr()
		q = &CompoundQuery{
			Op:       op,
			Distinct: distinct,
			Left:     q,
			Right:    e,
		}
	}

	// TODO: parse ORDER BY and LIMIT

	return q
}

func (p *Parser) parseSimpleQueryExpr() QueryExpr {
	if q := p.tryParseSubQuery(); q != nil {
		return q
	}
	// When p.Token.Kind is "(" even if p.tryParseSubQuery() returns nil,
	// but it should be parsed as SubQuery like (e.g. ((select 1)).)
	if p.Token.Kind == "(" {
		pos := p.Token.Pos
		p.NextToken()
		q := p.parseQueryExpr()
		end := p.expect(")").End
		return &SubQuery{
			pos:  pos,
			end:  end,
			Expr: q,
		}
	}
	return p.parseSelect()
}

func (p *Parser) parseSelect() *Select {
	pos := p.expect("SELECT").Pos
	var distinct bool
	if p.Token.Kind == "DISTINCT" {
		p.NextToken()
		distinct = true
	}
	var asStruct bool
	if p.Token.Kind == "AS" {
		p.NextToken()
		p.expect("STRUCT")
		asStruct = true
	}

	es := p.parseSelectExprList()
	end := es.End()

	var from FromItemList
	if p.Token.Kind == "FROM" {
		p.NextToken()
		from = p.parseFromItemList()
		end = from.End()
	}

	// TODO: parse WHERE, GROUP BY, HAVING, ORDER BY and LIMIT

	return &Select{
		pos:      pos,
		end:      end,
		Distinct: distinct,
		AsStruct: asStruct,
		List:     es,
		From:     from,
	}
}

func (p *Parser) parseSelectExprList() SelectExprList {
	es := SelectExprList{p.parseSelectExpr()}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
		es = append(es, p.parseSelectExpr())
	}
	return es
}

func (p *Parser) parseSelectExpr() *SelectExpr {
	if p.Token.Kind == "*" {
		e := &SelectExpr{
			pos:  p.Token.Pos,
			end:  p.Token.End,
			Star: true,
		}
		p.NextToken()
		return e
	}

	pos := p.Token.Pos
	e := p.parseExpr()
	switch p.Token.Kind {
	case TokenIdent:
		se := &SelectExpr{
			pos:  pos,
			end:  p.Token.End,
			Expr: e,
			As: &Ident{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: p.Token.AsString,
			},
		}
		p.NextToken()
		return se
	case "AS":
		p.NextToken()
		id := p.expect(TokenIdent)
		return &SelectExpr{
			pos:  pos,
			end:  id.End,
			Expr: e,
			As: &Ident{
				pos:  id.Pos,
				end:  id.End,
				Name: id.AsString,
			},
		}
	case ".":
		p.NextToken()
		end := p.expect("*").End
		return &SelectExpr{
			pos:  pos,
			end:  end,
			Star: true,
			Expr: e,
		}
	}

	return &SelectExpr{
		pos:  pos,
		end:  e.End(),
		Expr: e,
	}
}

func (p *Parser) parseFromItemList() FromItemList {
	from := FromItemList{p.parseFromItem()}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
		from = append(from, p.parseFromItem())
	}
	return from
}

func (p *Parser) parseFromItem() *FromItem {
	e := p.parseJoinExpr(false)
	end := e.End()
	var tableSample TableSampleMethod
	if p.Token.Kind == "TABLESAMPLE" {
		p.NextToken()
		id := p.expect(TokenIdent)
		end = id.End
		switch {
		case strings.EqualFold(id.AsString, "BERNOULLI"):
			tableSample = BernoulliSampleMethod
		case strings.EqualFold(id.AsString, "RESERVOIR"):
			tableSample = ReservoirSampleMethod
		default:
			p.panicfAtToken(&p.Token, "expected identifier: BERNOULLI, RESERVOIR, but: %s", p.Token.Raw)
		}
	}
	return &FromItem{
		end:         end,
		Expr:        e,
		TableSample: tableSample,
	}
}

func (p *Parser) parseJoinExpr(needOp bool) JoinExpr {
	j := p.parseSimpleJoinExpr()
	for {
		if _, ok := j.(*SubQueryJoinExpr); ok {
			needOp = false
		}

		op := InnerJoin
		switch p.Token.Kind {
		case "INNER":
			p.NextToken()
			needOp = true
		case "CROSS":
			p.NextToken()
			op = CrossJoin
			needOp = true
		case "FULL":
			p.NextToken()
			if p.Token.Kind == "OUTER" {
				p.NextToken()
			}
			op = FullOuterJoin
			needOp = true
		case "LEFT":
			p.NextToken()
			if p.Token.Kind == "OUTER" {
				p.NextToken()
			}
			op = LeftOuterJoin
			needOp = true
		case "RIGHT":
			p.NextToken()
			if p.Token.Kind == "OUTER" {
				p.NextToken()
			}
			op = RightOuterJoin
			needOp = true
		}

		var method JoinMethod
		if p.Token.Kind == "HASH" {
			p.NextToken()
			method = HashJoinMethod
			needOp = true
		} else if p.Token.Kind == TokenIdent {
			if strings.EqualFold(p.Token.Raw, "APPLY") {
				p.NextToken()
				method = ApplyJoinMethod
				needOp = true
			} else if strings.EqualFold(p.Token.Raw, "LOOP") {
				p.NextToken()
				method = LoopJoinMethod
				needOp = true
			}
		}
		if needOp {
			p.expect("JOIN")
			needOp = false
		} else if p.Token.Kind == "JOIN" {
			p.NextToken()
		} else {
			return j
		}

		hint := p.parseHint()
		right := p.parseSimpleJoinExpr()

		if op == CrossJoin {
			j = &Join{
				Op:     op,
				Method: method,
				Hint:   hint,
				Left:   j,
				Right:  right,
			}
			continue
		}

		var using IdentList
		var on Expr
		pos := p.Token.Pos
		var end Pos
		switch p.Token.Kind {
		case "USING":
			p.NextToken()
			using = make(IdentList, 0)
			p.expect("(")
			if p.Token.Kind != ")" {
				for p.Token.Kind != TokenEOF {
					id := p.expect(TokenIdent)
					using = append(using, &Ident{
						pos:  id.Pos,
						end:  id.End,
						Name: id.AsString,
					})
					if p.Token.Kind != "," {
						break
					}
					p.NextToken()
				}
			}
			end = p.expect(")").End
		case "ON":
			p.NextToken()
			on = p.parseExpr()
			end = on.End()
		}
		j = &Join{
			Op:     op,
			Method: method,
			Hint:   hint,
			Left:   j,
			Right:  right,
			Cond: &JoinCondition{
				pos:   pos,
				end:   end,
				On:    on,
				Using: using,
			},
		}
	}
}

func (p *Parser) parseSimpleJoinExpr() JoinExpr {
	if q := p.tryParseSubQuery(); q != nil {
		return p.parseSubQueryJoinExprSuffix(q)
	}

	if p.Token.Kind == "(" {
		pos := p.Token.Pos
		p.NextToken()
		j := p.parseJoinExpr(true)
		end := p.expect(")").End
		return &ParenJoinExpr{
			pos:  pos,
			end:  end,
			Expr: j,
		}
	}

	if p.Token.Kind == "UNNEST" {
		pos := p.Token.Pos
		p.NextToken()
		p.expect("(")
		e := p.parseExpr()
		end := p.expect(")").End
		return p.parseUnnestSuffix(e, pos, end)
	}

	e := p.parseExpr()
	if id, ok := e.(*Ident); ok {
		hint := p.parseHint()
		as := p.parseAs()
		end := id.End()
		if as != nil {
			end = as.End()
		} else if hint != nil {
			end = hint.End()
		}
		return &TableName{
			end:   end,
			Ident: id,
			Hint:  hint,
			As:    as,
		}
	}

	return p.parseUnnestSuffix(e, e.Pos(), e.End())
}

func (p *Parser) parseSubQueryJoinExprSuffix(q *SubQuery) *SubQueryJoinExpr {
	hint := p.parseHint()
	as := p.parseAs()
	end := q.End()
	if as != nil {
		end = as.End()
	} else if hint != nil {
		end = hint.End()
	}
	return &SubQueryJoinExpr{
		end:  end,
		Expr: q,
		Hint: hint,
		As:   as,
	}
}

func (p *Parser) parseUnnestSuffix(e Expr, pos Pos, end Pos) *Unnest {
	hint := p.parseHint()
	as := p.parseAs()
	if as != nil {
		end = as.End()
	} else if hint != nil {
		end = hint.End()
	}
	return &Unnest{
		pos:  pos,
		end:  end,
		Expr: e,
		Hint: hint,
		As:   as,
	}
}

func (p *Parser) parseAs() *Ident {
	if p.Token.Kind != "AS" {
		return nil
	}
	p.NextToken()
	id := p.expect(TokenIdent)
	return &Ident{
		pos:  id.Pos,
		end:  id.End,
		Name: id.AsString,
	}
}

func (p *Parser) parseExpr() Expr {
	return p.parseAndOr()
}

func (p *Parser) parseAndOr() Expr {
	expr := p.parseNot()
	for {
		var op BinaryOp
		switch p.Token.Kind {
		case "AND":
			op = OpAnd
		case "OR":
			op = OpOr
		default:
			return expr
		}
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseNot(),
		}
	}
}

func (p *Parser) parseNot() Expr {
	if p.Token.Kind == "NOT" {
		p.NextToken()
		return &UnaryExpr{
			Op:   OpNot,
			Expr: p.parseNot(),
		}
	}

	return p.parseComparison()
}

func (p *Parser) parseComparison() Expr {
	expr := p.parseBitOr()
	var op BinaryOp
	switch p.Token.Kind {
	case "<":
		op = OpLess
	case ">":
		op = OpGreater
	case "<=":
		op = OpLessEqual
	case ">=":
		op = OpGreaterEqual
	case "=":
		op = OpEqual
	case "!=", "<>":
		op = OpNotEqual
	case "LIKE":
		op = OpLike
	case "IN":
		panic("TODO: IN")
	case "BETWEEN":
		panic("TODO: BETWEEN")
	case "NOT":
		p.NextToken()
		switch p.Token.Kind {
		case "LIKE":
			op = OpNotLike
		case "IN":
			panic("TODO: NOT IN")
		case "BETWEEN":
			panic("TODO: NOT BETWEEN")
		default:
			p.panicfAtToken(&p.Token, "expected token: LIKE, IN, but: %s", p.Token.Kind)
		}
	case "IS":
		p.NextToken()
		not := false
		if p.Token.Kind == "NOT" {
			p.NextToken()
			not = true
		}
		switch p.Token.Kind {
		case "NULL":
			return &IsNullExpr{Left: expr, Not: not}
		case "TRUE":
			return &IsBoolExpr{Left: expr, Not: not, Right: true}
		case "FALSE":
			return &IsBoolExpr{Left: expr, Not: not, Right: false}
		default:
			p.panicfAtToken(&p.Token, "expected token: NULL, TRUE, FALSE, but: %s", p.Token.Kind)
		}
	default:
		return expr
	}
	p.NextToken()
	return &BinaryExpr{
		Left:  expr,
		Op:    op,
		Right: p.parseBitOr(),
	}
}

func (p *Parser) parseBitOr() Expr {
	expr := p.parseBitXor()
	for p.Token.Kind == "|" {
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    OpBitOr,
			Right: p.parseBitXor(),
		}
	}
	return expr
}

func (p *Parser) parseBitXor() Expr {
	expr := p.parseBitAnd()
	for p.Token.Kind == "^" {
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    OpBitXor,
			Right: p.parseBitAnd(),
		}
	}
	return expr
}

func (p *Parser) parseBitAnd() Expr {
	expr := p.parseBitShift()
	for p.Token.Kind == "&" {
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    OpBitAnd,
			Right: p.parseBitShift(),
		}
	}
	return expr
}

func (p *Parser) parseBitShift() Expr {
	expr := p.parseAddSub()
	for {
		var op BinaryOp
		switch p.Token.Kind {
		case "<<":
			op = OpBitLeftShift
		case ">>":
			op = OpBitRightShift
		default:
			return expr
		}
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseAddSub(),
		}
	}
}

func (p *Parser) parseAddSub() Expr {
	expr := p.parseMulDiv()
	for {
		var op BinaryOp
		switch p.Token.Kind {
		case "+":
			op = OpAdd
		case "-":
			op = OpSub
		default:
			return expr
		}
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseMulDiv(),
		}
	}
}

func (p *Parser) parseMulDiv() Expr {
	expr := p.parseUnary()
	for {
		var op BinaryOp
		switch p.Token.Kind {
		case "*":
			op = OpMul
		case "/":
			op = OpDiv
		default:
			return expr
		}
		p.NextToken()
		expr = &BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseUnary(),
		}
	}
}

func (p *Parser) parseUnary() Expr {
	var op UnaryOp
	switch p.Token.Kind {
	case "+":
		op = OpPlus
	case "-":
		op = OpMinus
	case "~":
		op = OpBitNot
	default:
		return p.parseSelector()
	}
	pos := p.Token.Pos
	p.NextToken()
	e := p.parseUnary()
	switch e := e.(type) {
	case *IntLit:
		if e.Value[0] != '+' && e.Value[0] != '-' {
			e.pos = pos
			e.Value = string(op) + e.Value
			return e
		}
	case *FloatLit:
		if e.Value[0] != '+' && e.Value[0] != '-' {
			e.pos = pos
			e.Value = string(op) + e.Value
			return e
		}
	}

	return &UnaryExpr{
		pos:  pos,
		Op:   op,
		Expr: e,
	}
}

func (p *Parser) parseSelector() Expr {
	expr := p.parseLit()
	for {
		switch p.Token.Kind {
		case ".":
			l := *p.Lexer
			p.NextToken()
			if p.Token.Kind == "*" { // expr.* case
				*p.Lexer = l
				return expr
			}
			id := p.expect(TokenIdent)
			expr = &SelectorExpr{
				Left: expr,
				Right: &Ident{
					pos:  id.Pos,
					end:  id.End,
					Name: id.AsString,
				},
			}
		case "[":
			p.NextToken()
			expr = &IndexExpr{
				Left:  expr,
				Right: p.parseExpr(),
			}
			expr.(*IndexExpr).end = p.expect("]").End
		default:
			return expr
		}
	}
}

func (p *Parser) parseLit() Expr {
	switch p.Token.Kind {
	case "NULL":
		lit := &NullLit{
			pos: p.Token.Pos,
		}
		p.NextToken()
		return lit
	case "TRUE", "FALSE":
		lit := &BoolLit{
			pos:   p.Token.Pos,
			Value: p.Token.Kind == "TRUE",
		}
		p.NextToken()
		return lit
	case TokenInt:
		lit := &IntLit{
			pos:   p.Token.Pos,
			end:   p.Token.End,
			Value: p.Token.Raw,
		}
		p.NextToken()
		return lit
	case TokenFloat:
		lit := &FloatLit{
			pos:   p.Token.Pos,
			end:   p.Token.End,
			Value: p.Token.Raw,
		}
		p.NextToken()
		return lit
	case TokenString:
		lit := &StringLit{
			pos:   p.Token.Pos,
			end:   p.Token.End,
			Value: p.Token.AsString,
		}
		p.NextToken()
		return lit
	case TokenBytes:
		lit := &BytesLit{
			pos:   p.Token.Pos,
			end:   p.Token.End,
			Value: []byte(p.Token.AsString),
		}
		p.NextToken()
		return lit
	case TokenParam:
		param := &Param{
			pos:  p.Token.Pos,
			Name: p.Token.AsString,
		}
		p.NextToken()
		return param
	}

	if p.Token.Kind == TokenIdent {
		id := p.Token
		p.NextToken()
		switch p.Token.Kind {
		case "(":
			p.NextToken()
			es, end := p.parseExprList(")")
			return &CallExpr{
				end: end,
				Func: &Ident{
					pos:  id.Pos,
					end:  id.End,
					Name: id.AsString,
				},
				Args: es,
			}
		case TokenString:
			// DATE "..."
			if strings.EqualFold(id.Raw, "DATE") {
				lit := &DateLit{
					pos:   id.Pos,
					end:   p.Token.End,
					Value: p.Token.AsString,
				}
				p.NextToken()
				return lit
			}
			// TIMESTAMP "..."
			if strings.EqualFold(id.Raw, "TIMESTAMP") {
				lit := &TimestampLit{
					pos:   id.Pos,
					end:   p.Token.End,
					Value: p.Token.AsString,
				}
				p.NextToken()
				return lit
			}
		}
		return &Ident{
			pos:  id.Pos,
			end:  id.End,
			Name: id.AsString,
		}
	}

	if p.Token.Kind == "ARRAY" {
		pos := p.Token.Pos
		p.NextToken()
		switch p.Token.Kind {
		case "(":
			parenPos := p.Token.Pos
			p.NextToken()
			q := &SubQuery{
				pos:  parenPos,
				Expr: p.parseQueryExpr(),
			}
			q.end = p.expect(")").End
			return &ArrayExpr{
				pos:  pos,
				Expr: q,
			}
		case "<":
			p.NextToken()
			t := p.parseType()
			p.expect(">")
			p.expect("[")
			es, end := p.parseExprList("]")
			return &ArrayLit{
				pos:    pos,
				end:    end,
				Type:   t,
				Values: es,
			}
		case "[":
			p.NextToken()
			es, end := p.parseExprList("]")
			return &ArrayLit{
				pos:    pos,
				end:    end,
				Values: es,
			}
		}

		p.panicfAtToken(&p.Token, "expected token: (, <, [, but: %s", p.Token.Kind)
	}

	if p.Token.Kind == "STRUCT" {
		pos := p.Token.Pos
		p.NextToken()
		switch p.Token.Kind {
		case "(":
			p.NextToken()
			es, end := p.parseExprList(")")
			return &StructLit{
				pos:    pos,
				end:    end,
				Values: es,
			}
		case "<>":
			p.NextToken()
			p.expect("(")
			es, end := p.parseExprList(")")
			return &StructLit{
				pos:    pos,
				end:    end,
				Fields: make([]*FieldSchema, 0),
				Values: es,
			}
		case "<":
			p.NextToken()
			fs, _ := p.parseFieldSchemas()
			p.expect("(")
			es, end := p.parseExprList(")")
			return &StructLit{
				pos:    pos,
				end:    end,
				Fields: fs,
				Values: es,
			}
		}

		p.panicfAtToken(&p.Token, "expected token: (, <, <>, but: %s", p.Token.Kind)
	}

	if p.Token.Kind == "CAST" {
		pos := p.Token.Pos
		p.NextToken()
		p.expect("(")
		e := p.parseExpr()
		p.expect("AS")
		t := p.parseType()
		end := p.expect(")").End
		return &CastExpr{
			pos:  pos,
			end:  end,
			Expr: e,
			Type: t,
		}
	}

	if p.Token.Kind == "(" {
		paren := p.Token
		if subQuery := p.tryParseSubQuery(); subQuery != nil {
			return subQuery
		}
		pos := p.Token.Pos
		p.NextToken()
		e := p.parseExpr()
		if p.Token.Kind == ")" {
			end := p.Token.End
			p.NextToken()
			return &ParenExpr{
				pos:  pos,
				end:  end,
				Expr: e,
			}
		}
		if p.Token.Kind != "," {
			p.panicfAtToken(&paren, "cannot parse (...) as expression, struct literal or subquery")
		}

		es := ExprList{e}
		for p.Token.Kind != TokenEOF {
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
			es = append(es, p.parseExpr())
		}
		end := p.expect(")").End
		return &StructLit{
			pos:    paren.Pos,
			end:    end,
			Values: es,
		}
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseExprList(end TokenKind) (ExprList, Pos) {
	var es ExprList
	if p.Token.Kind == end {
		return es, p.expect(end).End
	}

	for p.Token.Kind != TokenEOF {
		es = append(es, p.parseExpr())
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	return es, p.expect(end).End
}

func (p *Parser) tryParseSubQuery() *SubQuery {
	l := *p.Lexer
	pos := p.Token.Pos
	if p.Token.Kind != "(" {
		return nil
	}
	p.NextToken()
	if p.Token.Kind == "(" {
		p.NextToken()
		nest := 1
		for p.Token.Kind != TokenEOF {
			if p.Token.Kind == "(" {
				nest += 1
			}
			if p.Token.Kind == ")" {
				nest -= 1
			}
			p.NextToken()
			if nest == 0 {
				break
			}
		}
		if p.Token.Kind == "UNION" || p.Token.Kind == "INTERSECT" || p.Token.Kind == "EXCEPT" {
			*p.Lexer = l
			p.NextToken()
		} else {
			*p.Lexer = l
			return nil
		}
	} else if p.Token.Kind != "SELECT" {
		*p.Lexer = l
		return nil
	}
	q := p.parseQueryExpr()
	end := p.expect(")").End
	return &SubQuery{
		pos:  pos,
		end:  end,
		Expr: q,
	}
}

func (p *Parser) parseType() *Type {
	switch p.Token.Kind {
	case TokenIdent:
		if strings.EqualFold(p.Token.AsString, "BOOL") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: BoolType,
			}
			p.NextToken()
			return t
		}
		if strings.EqualFold(p.Token.AsString, "INT64") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: Int64Type,
			}
			p.NextToken()
			return t
		}
		if strings.EqualFold(p.Token.AsString, "FLOAT64") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: Float64Type,
			}
			p.NextToken()
			return t
		}
		if strings.EqualFold(p.Token.AsString, "DATE") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: DateType,
			}
			p.NextToken()
			return t
		}
		if strings.EqualFold(p.Token.AsString, "TIMESTAMP") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: TimestampType,
			}
			p.NextToken()
			return t
		}
		if strings.EqualFold(p.Token.AsString, "STRING") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: StringType,
			}
			p.NextToken()
			if p.Token.Kind == "(" {
				p.NextToken()
				var max bool
				var length IntValue
				switch p.Token.Kind {
				case TokenParam:
					length = &Param{
						pos:  p.Token.Pos,
						Name: p.Token.AsString,
					}
				case TokenInt:
					length = &IntLit{
						pos:   p.Token.Pos,
						end:   p.Token.End,
						Value: p.Token.Raw,
					}
				case TokenIdent:
					if !strings.EqualFold(p.Token.AsString, "MAX") { // TODO: is `Max` allowed?
						p.panicfAtToken(&p.Token, "expected identifier: MAX, but: %s", p.Token.Raw)
					}
					max = true
				}
				p.NextToken()
				t.MaxLength = max
				t.Length = length
				t.end = p.expect(")").End
			}
			return t
		}
		if strings.EqualFold(p.Token.AsString, "BYTES") {
			t := &Type{
				pos:  p.Token.Pos,
				end:  p.Token.End,
				Name: BytesType,
			}
			p.NextToken()
			if p.Token.Kind == "(" {
				p.NextToken()
				var max bool
				var length IntValue
				switch p.Token.Kind {
				case TokenParam:
					length = &Param{
						pos:  p.Token.Pos,
						Name: p.Token.AsString,
					}
				case TokenInt:
					length = &IntLit{
						pos:   p.Token.Pos,
						end:   p.Token.End,
						Value: p.Token.Raw,
					}
				case TokenIdent:
					if !strings.EqualFold(p.Token.AsString, "MAX") { // TODO: is `Max` allowed?
						p.panicfAtToken(&p.Token, "expected identifier: MAX, but: %s", p.Token.Raw)
					}
					max = true
				}
				p.NextToken()
				t.MaxLength = max
				t.Length = length
				t.end = p.expect(")").End
			}
			return t
		}

		p.panicfAtToken(&p.Token, "expected identifier: BOOL, INT64, FLOAT64, DATE, TIMESTAMP, STRING, BYTES, but: %s", p.Token.Raw)
	case "ARRAY":
		pos := p.Token.Pos
		p.NextToken()
		p.expect("<")
		t := p.parseType()
		end := p.expect(">").End
		return &Type{
			pos:   pos,
			end:   end,
			Name:  ArrayType,
			Value: t,
		}
	case "STRUCT":
		pos := p.Token.Pos
		p.NextToken()
		switch p.Token.Kind {
		case "<>":
			return &Type{
				pos:    pos,
				end:    p.Token.End,
				Name:   StructType,
				Fields: make([]*FieldSchema, 0),
			}
		case "<":
			p.NextToken()
			fs, end := p.parseFieldSchemas()
			return &Type{
				pos:    pos,
				end:    end,
				Name:   StructType,
				Fields: fs,
			}
		}
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <ident>, ARRAY, STRUCT, but: %s", p.Token.Kind))
}

func (p *Parser) parseFieldSchemas() ([]*FieldSchema, Pos) {
	fs := make([]*FieldSchema, 0)
	if p.Token.Kind == ">" {
		return fs, p.expect(">").End
	}
	for p.Token.Kind != TokenEOF {
		named := false
		l := *p.Lexer
		if p.Token.Kind == TokenIdent {
			id := p.Token
			p.NextToken()
			if p.Token.Kind == TokenIdent || p.Token.Kind == "ARRAY" || p.Token.Kind == "STRUCT" {
				t := p.parseType()
				fs = append(fs, &FieldSchema{
					Name: &Ident{
						pos:  id.Pos,
						end:  id.End,
						Name: id.AsString,
					},
					Type: t,
				})
				named = true
			} else {
				*p.Lexer = l
			}
		}
		if !named {
			t := p.parseType()
			fs = append(fs, &FieldSchema{
				Type: t,
			})
		}
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	return fs, p.expect(">").End
}

func (p *Parser) expect(kind TokenKind) *Token {
	if p.Token.Kind != kind {
		p.panicfAtToken(&p.Token, "expected token: %s, but: %s", kind, p.Token.Kind)
	}
	tok := p.Token
	p.NextToken()
	return &tok
}

func (p *Parser) errorfAtToken(tok *Token, msg string, params ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, params...),
		Position: p.Position(tok.Pos, tok.End),
	}
}

func (p *Parser) panicfAtToken(tok *Token, msg string, params ...interface{}) {
	panic(p.errorfAtToken(tok, msg, params...))
}
