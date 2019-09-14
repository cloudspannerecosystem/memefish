package parser

import (
	"fmt"
	"strings"
)

type Parser struct {
	*Lexer
}

// ParseStatement parses a SQL statement.
func (p *Parser) ParseStatement() (stmt Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			stmt = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	stmt = p.parseStatement()
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseStatements parses SQL statements list separated by semi-colon.
func (p *Parser) ParseStatements() (stmts []Statement, err error) {
	defer func() {
		if r := recover(); r != nil {
			stmts = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	p.parseStatements(func() {
		stmts = append(stmts, p.parseStatement())
	})
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseQuery parses a SELECT query statement.
func (p *Parser) ParseQuery() (stmt *QueryStatement, err error) {
	defer func() {
		if r := recover(); r != nil {
			stmt = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	stmt = p.parseQueryStatement()
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

func (p *Parser) ParseExpr() (expr Expr, err error) {
	defer func() {
		if r := recover(); r != nil {
			expr = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	expr = p.parseExpr()
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDDL parses a CREATE/ALTER/DROP statement.
func (p *Parser) ParseDDL() (ddl DDL, err error) {
	defer func() {
		if r := recover(); r != nil {
			ddl = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	ddl = p.parseDDL()
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDDLs parses CREATE/ALTER/DROP statements list separated by semi-colon.
func (p *Parser) ParseDDLs() (ddls []DDL, err error) {
	defer func() {
		if r := recover(); r != nil {
			ddls = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	p.parseStatements(func() {
		ddls = append(ddls, p.parseDDL())
	})
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDML parses a INSERT/DELETE/UPDATE statement.
func (p *Parser) ParseDML() (dml DML, err error) {
	defer func() {
		if r := recover(); r != nil {
			dml = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	dml = p.parseDML()
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDMLs parses INSERT/DELETE/UPDATE statements list separated by semi-colon.
func (p *Parser) ParseDMLs() (dmls []DML, err error) {
	defer func() {
		if r := recover(); r != nil {
			dmls = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.NextToken()
	p.parseStatements(func() {
		dmls = append(dmls, p.parseDML())
	})
	if p.Token.Kind != TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

func (p *Parser) parseStatement() Statement {
	switch {
	case p.Token.Kind == "SELECT" || p.Token.Kind == "@" || p.Token.Kind == "(":
		return p.parseQueryStatement()
	case p.Token.Kind == "CREATE" || p.Token.IsKeywordLike("ALTER") || p.Token.IsKeywordLike("DROP"):
		return p.parseDDL()
	case p.Token.IsKeywordLike("INSERT") || p.Token.IsKeywordLike("DELETE") || p.Token.IsKeywordLike("UPDATE"):
		return p.parseDML()
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseStatements(doParse func()) {
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == ";" {
			p.NextToken()
			continue
		}

		doParse()

		if p.Token.Kind != ";" {
			break
		}
	}
}

// ================================================================================
//
// SELECT
//
// ================================================================================

func (p *Parser) parseQueryStatement() *QueryStatement {
	hint := p.tryParseHint()
	query := p.parseQueryExpr()

	return &QueryStatement{
		Hint:  hint,
		Query: query,
	}
}

func (p *Parser) tryParseHint() *Hint {
	if p.Token.Kind != "@" {
		return nil
	}

	pos := p.Token.Pos
	p.NextToken()
	p.expect("{")
	records := []*HintRecord{p.parseHintRecord()}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
		records = append(records, p.parseHintRecord())
	}
	end := p.expect("}").End
	return &Hint{
		pos:     pos,
		end:     end,
		Records: records,
	}
}

func (p *Parser) parseHintRecord() *HintRecord {
	key := p.parseIdent()
	p.expect("=")
	value := p.parseExpr()
	return &HintRecord{
		Key:   key,
		Value: value,
	}
}

func (p *Parser) parseQueryExpr() QueryExpr {
	query := p.parseSimpleQueryExpr()

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
		opTok := p.Token
		p.NextToken()

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

		right := p.parseSimpleQueryExpr()
		if c, ok := query.(*CompoundQuery); ok {
			if !(c.Op == op && c.Distinct == distinct) {
				p.panicfAtToken(&opTok, "all set operator at the same level must be the same, or wrap (...)")
			}
			c.Queries = append(c.Queries, right)
		} else {
			query = &CompoundQuery{
				Op:       op,
				Distinct: distinct,
				Queries:  []QueryExpr{query, right},
			}
		}
	}

	return p.parseQueryExprSuffix(query)
}

func (p *Parser) parseSimpleQueryExpr() QueryExpr {
	if p.Token.Kind == "(" {
		pos := p.expect("(").Pos
		query := p.parseQueryExpr()
		end := p.expect(")").End
		return &SubQuery{
			pos:   pos,
			end:   end,
			Query: query,
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

	results := p.parseSelectResults()
	from := p.tryParseFrom()
	where := p.tryParseWhere()
	groupBy := p.tryParseGroupBy()
	having := p.tryParseHaving()

	return &Select{
		pos:      pos,
		Distinct: distinct,
		AsStruct: asStruct,
		Results:  results,
		From:     from,
		Where:    where,
		GroupBy:  groupBy,
		Having:   having,
	}
}

func (p *Parser) parseSelectResults() []SelectItem {
	results := []SelectItem{p.parseSelectItem()}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
		results = append(results, p.parseSelectItem())
	}
	return results
}

func (p *Parser) parseSelectItem() SelectItem {
	if p.Token.Kind == "*" {
		star := p.expect("*")
		return &Star{
			pos: star.Pos,
		}
	}

	expr := p.parseExpr()
	if as := p.tryParseAsAlias(); as != nil {
		return &Alias{
			Expr: expr,
			As:   as,
		}
	}

	if p.Token.Kind == "." {
		p.NextToken()
		end := p.expect("*").End
		return &DotStar{
			end:  end,
			Expr: expr,
		}
	}

	return &ExprSelectItem{
		Expr: expr,
	}
}

func (p *Parser) tryParseAsAlias() *AsAlias {
	pos := p.Token.Pos

	if p.Token.Kind == "AS" {
		p.NextToken()
		id := p.parseIdent()
		return &AsAlias{
			pos:   pos,
			Alias: id,
		}
	}

	if p.Token.Kind == TokenIdent {
		id := p.parseIdent()
		return &AsAlias{
			pos:   pos,
			Alias: id,
		}
	}

	return nil
}

func (p *Parser) tryParseFrom() *From {
	if p.Token.Kind != "FROM" {
		return nil
	}
	pos := p.expect("FROM").Pos

	j := p.parseTableExpr(true, false)

	return &From{
		pos:    pos,
		Source: j,
	}
}

func (p *Parser) tryParseWhere() *Where {
	if p.Token.Kind != "WHERE" {
		return nil
	}
	return p.parseWhere()
}

func (p *Parser) parseWhere() *Where {
	pos := p.expect("WHERE").Pos
	expr := p.parseExpr()
	return &Where{
		pos:  pos,
		Expr: expr,
	}
}

func (p *Parser) tryParseGroupBy() *GroupBy {
	if p.Token.Kind != "GROUP" {
		return nil
	}
	pos := p.expect("GROUP").Pos
	p.expect("BY")
	exprs := []Expr{p.parseExpr()}
	for p.Token.Kind == "," {
		p.NextToken()
		exprs = append(exprs, p.parseExpr())
	}

	return &GroupBy{
		pos:   pos,
		Exprs: exprs,
	}
}

func (p *Parser) tryParseHaving() *Having {
	if p.Token.Kind != "HAVING" {
		return nil
	}
	pos := p.expect("HAVING").Pos
	expr := p.parseExpr()
	return &Having{
		pos:  pos,
		Expr: expr,
	}
}

func (p *Parser) parseQueryExprSuffix(e QueryExpr) QueryExpr {
	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()
	e.setQueryExprSuffix(orderBy, limit)
	return e
}

func (p *Parser) tryParseOrderBy() *OrderBy {
	if p.Token.Kind != "ORDER" {
		return nil
	}

	pos := p.expect("ORDER").Pos
	p.expect("BY")

	items := []*OrderByItem{p.parseOrderByItem()}
	for p.Token.Kind == "," {
		p.NextToken()
		items = append(items, p.parseOrderByItem())
	}

	return &OrderBy{
		pos:   pos,
		Items: items,
	}
}

func (p *Parser) parseOrderByItem() *OrderByItem {
	expr := p.parseExpr()

	var end Pos
	collate := p.tryParseCollate()
	if collate != nil {
		end = collate.End()
	}

	var dir Direction
	switch p.Token.Kind {
	case "ASC":
		end = p.expect("ASC").End
		dir = DirectionAsc
	case "DESC":
		end = p.expect("DESC").End
		dir = DirectionDesc
	}

	return &OrderByItem{
		end:     end,
		Expr:    expr,
		Collate: collate,
		Dir:     dir,
	}
}

func (p *Parser) tryParseCollate() *Collate {
	if p.Token.Kind != "COLLATE" {
		return nil
	}
	pos := p.expect("COLLATE").Pos
	value := p.parseStringValue()
	return &Collate{
		pos:   pos,
		Value: value,
	}
}

func (p *Parser) tryParseLimit() *Limit {
	if p.Token.Kind != "LIMIT" {
		return nil
	}

	pos := p.expect("LIMIT").Pos
	count := p.parseIntValue()
	offset := p.tryParseOffset()

	return &Limit{
		pos:    pos,
		Count:  count,
		Offset: offset,
	}
}

func (p *Parser) tryParseOffset() *Offset {
	if !p.Token.IsKeywordLike("OFFSET") {
		return nil
	}
	pos := p.expectKeywordLike("OFFSET").Pos
	value := p.parseIntValue()
	return &Offset{
		pos:   pos,
		Value: value,
	}
}

// ================================================================================
//
// JOIN
//
// ================================================================================

func (p *Parser) parseTableExpr(toplevel, needOp bool) TableExpr {
	j := p.parseSimpleTableExpr()
	for {
		needOp = needOp && j.isSimpleTableExpr()

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
		if toplevel && p.Token.Kind == "," {
			op = CommaJoin
		}

		var method JoinMethod
		if op != CommaJoin {
			switch {
			case p.Token.Kind == "HASH":
				p.NextToken()
				method = HashJoinMethod
				needOp = true
			case p.Token.IsKeywordLike("APPLY"):
				p.NextToken()
				method = ApplyJoinMethod
				needOp = true
			case p.Token.IsKeywordLike("LOOP"):
				p.NextToken()
				method = LoopJoinMethod
				needOp = true
			}
		}

		switch {
		case needOp:
			p.expect("JOIN")
			needOp = false
		case op == CommaJoin || p.Token.Kind == "JOIN":
			p.NextToken()
		default:
			return j
		}

		hint := p.tryParseHint()
		right := p.parseSimpleTableExpr()

		if op == CrossJoin || op == CommaJoin {
			j = &Join{
				Op:     op,
				Method: method,
				Hint:   hint,
				Left:   j,
				Right:  right,
			}
			continue
		}

		cond := p.parseJoinCondition()
		j = &Join{
			Op:     op,
			Method: method,
			Hint:   hint,
			Left:   j,
			Right:  right,
			Cond:   cond,
		}
	}
}

func (p *Parser) parseSimpleTableExpr() TableExpr {
	if p.lookaheadSubQuery() {
		pos := p.expect("(").Pos
		query := p.parseQueryExpr()
		end := p.expect(")").End
		as := p.tryParseAsAlias()
		if as != nil {
			end = as.End()
		}
		return p.parseTableExprSuffix(&SubQueryTableExpr{
			pos:   pos,
			end:   end,
			Query: query,
			As:    as,
		})
	}

	if p.Token.Kind == "(" {
		pos := p.expect("(").Pos
		j := p.parseTableExpr(false, true)
		end := p.expect(")").End
		return p.parseTableExprSuffix(&ParenTableExpr{
			pos:    pos,
			end:    end,
			Source: j,
		})
	}

	if p.Token.Kind == "UNNEST" {
		pos := p.expect("UNNEST").Pos
		p.expect("(")
		expr := p.parseExpr()
		end := p.expect(")").End
		return p.parseUnnestSuffix(false, expr, pos, end)
	}

	if p.Token.Kind == TokenIdent {
		expr := p.parseIdentOrPath()
		if id, ok := expr.(*Ident); ok {
			return p.parseTableNameSuffix(id)
		}
		return p.parseUnnestSuffix(true, expr, expr.Pos(), expr.End())
	}

	panic(p.errorfAtToken(&p.Token, "expected token: (, UNNEST, <ident>, but: %s", p.Token.Kind))
}

func (p *Parser) parseIdentOrPath() Expr {
	idents := []*Ident{p.parseIdent()}
	for p.Token.Kind == "." {
		p.NextToken()
		idents = append(idents, p.parseIdent())
	}
	if len(idents) == 1 {
		return idents[0]
	}
	return &Path{
		Idents: idents,
	}
}

func (p *Parser) parseUnnestSuffix(implicit bool, expr Expr, pos Pos, end Pos) TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias()
	withOffset := p.tryParseWithOffset()

	if withOffset != nil {
		end = withOffset.End()
	}
	if as != nil {
		end = as.End()
	}
	if hint != nil {
		end = hint.End()
	}

	return p.parseTableExprSuffix(&Unnest{
		pos:        pos,
		end:        end,
		Implicit:   implicit,
		Expr:       expr,
		Hint:       hint,
		As:         as,
		WithOffset: withOffset,
	})
}

func (p *Parser) tryParseWithOffset() *WithOffset {
	if p.Token.Kind != "WITH" {
		return nil
	}
	pos := p.expect("WITH").Pos
	end := p.expectKeywordLike("OFFSET").End

	as := p.tryParseAsAlias()
	if as != nil {
		end = as.End()
	}

	return &WithOffset{
		pos: pos,
		end: end,
		As:  as,
	}
}

func (p *Parser) parseTableNameSuffix(id *Ident) TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias()
	return p.parseTableExprSuffix(&TableName{
		Table: id,
		Hint:  hint,
		As:    as,
	})
}

func (p *Parser) parseJoinCondition() JoinCondition {
	switch p.Token.Kind {
	case "ON":
		return p.parseOn()
	case "USING":
		return p.parseUsing()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ON, USING, but: %s", p.Token.Kind))
}

func (p *Parser) parseOn() *On {
	pos := p.expect("ON").Pos
	expr := p.parseExpr()
	return &On{
		pos:  pos,
		Expr: expr,
	}
}

func (p *Parser) parseUsing() *Using {
	pos := p.expect("USING").Pos
	p.expect("(")
	idents := []*Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.NextToken()
		idents = append(idents, p.parseIdent())
	}
	end := p.expect(")").End
	return &Using{
		pos:    pos,
		end:    end,
		Idents: idents,
	}
}

func (p *Parser) parseTableExprSuffix(j TableExpr) TableExpr {
	sample := p.tryParseTableSample()
	if sample != nil {
		j.setSample(sample)
	}
	return j
}

func (p *Parser) tryParseTableSample() *TableSample {
	if p.Token.Kind != "TABLESAMPLE" {
		return nil
	}
	pos := p.expect("TABLESAMPLE").Pos

	id := p.expect(TokenIdent)
	var method TableSampleMethod
	switch {
	case id.IsIdent("BERNOULLI"):
		method = BernoulliSampleMethod
	case id.IsIdent("RESERVOIR"):
		method = ReservoirSampleMethod
	default:
		p.panicfAtToken(id, "expected identifier: BERNOULLI, RESERVOIR, but: %s", id.Raw)
	}

	size := p.parseTableSampleSize()

	return &TableSample{
		pos:    pos,
		Method: method,
		Size:   size,
	}
}

func (p *Parser) parseTableSampleSize() *TableSampleSize {
	pos := p.expect("(").Pos

	value := p.parseNumValue()
	var unit TableSampleUnit
	switch {
	case p.Token.Kind == "ROWS":
		unit = RowsTableSampleUnit
	case p.Token.IsKeywordLike("PERCENT"):
		unit = PercentTableSampleUnit
	default:
		p.panicfAtToken(&p.Token, "expected token: ROWS, PERCENT, but: %s", p.Token.Kind)
	}
	p.NextToken()

	end := p.expect(")").End
	return &TableSampleSize{
		pos:   pos,
		end:   end,
		Value: value,
		Unit:  unit,
	}
}

// ================================================================================
//
// Expr
//
// ================================================================================

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
		p.NextToken()
		cond := p.parseInCondition()
		return &InExpr{
			Left:  expr,
			Right: cond,
		}
	case "BETWEEN":
		p.NextToken()
		rightStart := p.parseBitOr()
		p.expect("AND")
		rightEnd := p.parseBitOr()
		return &BetweenExpr{
			Left:       expr,
			RightStart: rightStart,
			RightEnd:   rightEnd,
		}
	case "NOT":
		p.NextToken()
		switch p.Token.Kind {
		case "LIKE":
			op = OpNotLike
		case "IN":
			p.NextToken()
			cond := p.parseInCondition()
			return &InExpr{
				Not:   true,
				Left:  expr,
				Right: cond,
			}
		case "BETWEEN":
			p.NextToken()
			rightStart := p.parseBitOr()
			p.expect("AND")
			rightEnd := p.parseBitOr()
			return &BetweenExpr{
				Not:        true,
				Left:       expr,
				RightStart: rightStart,
				RightEnd:   rightEnd,
			}
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
			p.NextToken()
			return &IsNullExpr{Left: expr, Not: not}
		case "TRUE":
			p.NextToken()
			return &IsBoolExpr{Left: expr, Not: not, Right: true}
		case "FALSE":
			p.NextToken()
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

func (p *Parser) parseInCondition() InCondition {
	if p.lookaheadSubQuery() {
		pos := p.expect("(").Pos
		query := p.parseQueryExpr()
		end := p.expect(")").End
		return &SubQueryInCondition{
			pos:   pos,
			end:   end,
			Query: query,
		}
	}

	if p.Token.Kind == "(" {
		pos := p.Token.Pos
		p.NextToken()
		exprs := []Expr{p.parseExpr()}
		for p.Token.Kind != TokenEOF {
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
			exprs = append(exprs, p.parseExpr())
		}
		end := p.expect(")").End
		return &ValuesInCondition{
			pos:   pos,
			end:   end,
			Exprs: exprs,
		}
	}

	if p.Token.Kind == "UNNEST" {
		pos := p.Token.Pos
		p.NextToken()
		p.expect("(")
		e := p.parseExpr()
		end := p.expect(")").End
		return &UnnestInCondition{
			pos:  pos,
			end:  end,
			Expr: e,
		}
	}

	panic(p.errorfAtToken(&p.Token, "expected token (, UNNEST, but: %s", p.Token.Kind))
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
	if op != OpBitNot {
		switch e := e.(type) {
		case *IntLiteral:
			if e.Value[0] != '+' && e.Value[0] != '-' {
				e.pos = pos
				e.Value = string(op) + e.Value
				return e
			}
		case *FloatLiteral:
			if e.Value[0] != '+' && e.Value[0] != '-' {
				e.pos = pos
				e.Value = string(op) + e.Value
				return e
			}
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
			lexer := p.Lexer.Clone()
			p.NextToken()
			if p.Token.Kind == "*" { // expr.* case
				p.Lexer = lexer
				return expr
			}

			member := p.parseIdent()
			switch e := expr.(type) {
			case *Ident:
				expr = &Path{
					Idents: []*Ident{e, member},
				}
			case *Path:
				e.Idents = append(e.Idents, member)
			default:
				expr = &SelectorExpr{
					Expr:   expr,
					Member: member,
				}
			}
		case "[":
			p.NextToken()
			id := p.expect(TokenIdent)
			ordinal := false
			if strings.EqualFold(id.AsString, "ORDINAL") {
				ordinal = true
			} else if strings.EqualFold(id.AsString, "OFFSET") {
				ordinal = false
			} else {
				p.panicfAtToken(id, "expected identifier: ORDINAL, OFFSET, but: %s", id.Raw)
			}
			p.expect("(")
			index := p.parseExpr()
			p.expect(")")
			end := p.expect("]").End
			expr = &IndexExpr{
				end:     end,
				Ordinal: ordinal,
				Expr:    expr,
				Index:   index,
			}
		default:
			return expr
		}
	}
}

func (p *Parser) parseLit() Expr {
	switch p.Token.Kind {
	case "NULL":
		return p.parseNullLiteral()
	case "TRUE", "FALSE":
		return p.parseBoolLiteral()
	case TokenInt:
		return p.parseIntLiteral()
	case TokenFloat:
		return p.parseFloatLiteral()
	case TokenString:
		return p.parseStringLiteral()
	case TokenBytes:
		return p.parseBytesLiteral()
	case TokenParam:
		return p.parseParam()
	case "CASE":
		return p.parseCaseExpr()
	case "CAST":
		return p.parseCastExpr()
	case "EXISTS":
		return p.parseExistsSubQuery()
	case "EXTRACT":
		return p.parseExtractExpr()
	case "ARRAY":
		return p.parseArrayLiteralOrSubQuery()
	case "STRUCT":
		return p.parseStructLiteral()
	case "[":
		return p.parseSimpleArrayLiteral()
	case "(":
		return p.parseParenExpr()
	case TokenIdent:
		id := p.Token
		p.NextToken()
		switch p.Token.Kind {
		case "(":
			return p.parseCall(id)
		case TokenString:
			if id.IsKeywordLike("DATE") {
				return p.parseDateLiteral(id)
			}
			if id.IsKeywordLike("TIMESTAMP") {
				return p.parseTimestampLiteral(id)
			}
		}
		return &Ident{
			pos:  id.Pos,
			end:  id.End,
			Name: id.AsString,
		}
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseCall(id Token) Expr {
	p.expect("(")
	if id.IsIdent("COUNT") && p.Token.Kind == "*" {
		p.NextToken()
		end := p.expect(")").End
		return &CountStarExpr{
			pos: id.Pos,
			end: end,
		}
	}

	fn := &Ident{
		pos:  id.Pos,
		end:  id.End,
		Name: id.AsString,
	}

	distinct := false
	if p.Token.Kind == "DISTINCT" {
		p.NextToken()
		distinct = true
	}

	var args []*Arg
	if p.Token.Kind != ")" {
		for p.Token.Kind != TokenEOF {
			args = append(args, p.parseArg())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}
	end := p.expect(")").End
	return &CallExpr{
		end:      end,
		Func:     fn,
		Distinct: distinct,
		Args:     args,
	}
}

func (p *Parser) parseArg() *Arg {
	if p.Token.Kind != "INTERVAL" {
		e := p.parseExpr()
		return &Arg{
			pos:  e.Pos(),
			Expr: e,
		}
	}

	pos := p.Token.Pos
	p.NextToken()
	e := p.parseExpr()
	unit := p.parseIdent()
	return &Arg{
		pos:          pos,
		Expr:         e,
		IntervalUnit: unit,
	}
}

func (p *Parser) parseCaseExpr() *CaseExpr {
	pos := p.expect("CASE").Pos
	var expr Expr
	if p.Token.Kind != "WHEN" {
		expr = p.parseExpr()
	}
	whens := []*CaseWhen{p.parseCaseWhen()}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "WHEN" {
			break
		}
		whens = append(whens, p.parseCaseWhen())
	}
	var els *CaseElse
	if p.Token.Kind == "ELSE" {
		els = p.parseCaseElse()
	}
	end := p.expect("END").End
	return &CaseExpr{
		pos:   pos,
		end:   end,
		Expr:  expr,
		Whens: whens,
		Else:  els,
	}
}

func (p *Parser) parseCaseWhen() *CaseWhen {
	pos := p.expect("WHEN").Pos
	cond := p.parseExpr()
	p.expect("THEN")
	then := p.parseExpr()
	return &CaseWhen{
		pos:  pos,
		Cond: cond,
		Then: then,
	}
}

func (p *Parser) parseCaseElse() *CaseElse {
	pos := p.expect("ELSE").Pos
	expr := p.parseExpr()
	return &CaseElse{
		pos:  pos,
		Expr: expr,
	}
}

func (p *Parser) parseCastExpr() *CastExpr {
	pos := p.expect("CAST").Pos
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

func (p *Parser) parseExistsSubQuery() *ExistsSubQuery {
	pos := p.expect("EXISTS").Pos
	p.expect("(")
	query := p.parseQueryExpr()
	end := p.expect(")").End
	return &ExistsSubQuery{
		pos:   pos,
		end:   end,
		Query: query,
	}
}

func (p *Parser) parseExtractExpr() *ExtractExpr {
	pos := p.expect("EXTRACT").Pos
	p.expect("(")
	part := p.parseIdent()
	p.expect("FROM")
	e := p.parseExpr()
	var atTimeZone *AtTimeZone
	if p.Token.Kind == "AT" {
		atTimeZone = p.parseAtTimeZone()
	}
	end := p.expect(")").End
	return &ExtractExpr{
		pos:        pos,
		end:        end,
		Part:       part,
		Expr:       e,
		AtTimeZone: atTimeZone,
	}
}

func (p *Parser) parseAtTimeZone() *AtTimeZone {
	pos := p.expect("AT").Pos
	p.expectKeywordLike("TIME")
	p.expectKeywordLike("ZONE")
	e := p.parseExpr()
	return &AtTimeZone{
		pos:  pos,
		Expr: e,
	}
}

func (p *Parser) parseParenExpr() Expr {
	paren := p.Token

	if p.lookaheadSubQuery() {
		p.NextToken()
		query := p.parseQueryExpr()
		end := p.expect(")").End
		return &ScalarSubQuery{
			pos:   paren.Pos,
			end:   end,
			Query: query,
		}
	}

	p.NextToken()
	expr := p.parseExpr()

	if p.Token.Kind == ")" {
		end := p.Token.End
		p.NextToken()
		return &ParenExpr{
			pos:  paren.Pos,
			end:  end,
			Expr: expr,
		}
	}

	if p.Token.Kind != "," {
		p.panicfAtToken(&paren, "cannot parse (...) as expression, struct literal or subquery")
	}

	values := []Expr{expr}
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
		values = append(values, p.parseExpr())
	}
	end := p.expect(")").End
	return &StructLiteral{
		pos:    paren.Pos,
		end:    end,
		Values: values,
	}
}

func (p *Parser) parseArrayLiteralOrSubQuery() Expr {
	pos := p.expect("ARRAY").Pos

	if p.Token.Kind == "(" {
		p.NextToken()
		query := p.parseQueryExpr()
		end := p.expect(")").End
		return &ArraySubQuery{
			pos:   pos,
			end:   end,
			Query: query,
		}
	}

	var t Type
	if p.Token.Kind == "<" {
		p.NextToken()
		t = p.parseType()
		p.expect(">")
	}

	values, _, end := p.parseArrayLiteralBody()
	return &ArrayLiteral{
		pos:    pos,
		end:    end,
		Type:   t,
		Values: values,
	}
}

func (p *Parser) parseSimpleArrayLiteral() *ArrayLiteral {
	values, pos, end := p.parseArrayLiteralBody()
	return &ArrayLiteral{
		pos:    pos,
		end:    end,
		Values: values,
	}
}

func (p *Parser) parseArrayLiteralBody() (values []Expr, pos, end Pos) {
	pos = p.expect("[").Pos
	if p.Token.Kind != "]" {
		for p.Token.Kind != TokenEOF {
			values = append(values, p.parseExpr())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}
	end = p.expect("]").End
	return
}

func (p *Parser) parseStructLiteral() *StructLiteral {
	pos := p.expect("STRUCT").Pos
	fields, _ := p.parseStructTypeFields(false)
	p.expect("(")
	var values []Expr
	if p.Token.Kind != ")" {
		for p.Token.Kind != TokenEOF {
			values = append(values, p.parseExpr())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}
	end := p.expect(")").End
	return &StructLiteral{
		pos:    pos,
		end:    end,
		Fields: fields,
		Values: values,
	}
}

func (p *Parser) parseDateLiteral(id Token) *DateLiteral {
	s := p.expect(TokenString)
	return &DateLiteral{
		pos:   id.Pos,
		end:   s.End,
		Value: s.AsString,
	}
}

func (p *Parser) parseTimestampLiteral(id Token) *TimestampLiteral {
	s := p.expect(TokenString)
	return &TimestampLiteral{
		pos:   id.Pos,
		end:   s.End,
		Value: s.AsString,
	}
}

func (p *Parser) lookaheadSubQuery() bool {
	lexer := p.Lexer.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if p.Token.Kind != "(" {
		return false
	}

	p.NextToken()
	// (SELECT ... indicates subquery.
	if p.Token.Kind == "SELECT" {
		return true
	}

	// ((...(SELECT maybe indicate subquery.
	nest := 0
	for p.Token.Kind == "(" {
		nest++
		p.NextToken()
	}
	if nest == 0 || p.Token.Kind != "SELECT" {
		return false
	}

	// ((...(SELECT ...)...) UNION indicates subquery.
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == "(" {
			nest++
		}
		if p.Token.Kind == ")" {
			nest--
		}

		if nest == 0 {
			break
		}
		p.NextToken()
	}
	if nest != 0 {
		return false
	}
	p.NextToken()
	switch p.Token.Kind {
	case "UNION", "INTERSECT", "EXCEPT", "ORDER", "LIMIT":
		return true
	}
	return false
}

func (p *Parser) parseType() Type {
	switch p.Token.Kind {
	case TokenIdent:
		return p.parseSimpleType()
	case "ARRAY":
		return p.parseArrayType()
	case "STRUCT":
		return p.parseStructType()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <ident>, ARRAY, STRUCT, but: %s", p.Token.Kind))
}

var simpleTypes = []string{
	"BOOL",
	"INT64",
	"FLOAT64",
	"DATE",
	"TIMESTAMP",
	"STRING",
	"BYTES",
}

func (p *Parser) parseSimpleType() *SimpleType {
	id := p.expect(TokenIdent)
	for _, typeName := range simpleTypes {
		if id.IsIdent(typeName) {
			return &SimpleType{
				pos:  id.Pos,
				Name: ScalarTypeName(typeName),
			}
		}
	}

	panic(p.errorfAtToken(id, "expected identifier: %s, but: %s", strings.Join(simpleTypes, ", "), id.Raw))
}

func (p *Parser) parseArrayType() *ArrayType {
	pos := p.expect("ARRAY").Pos
	p.expect("<")
	t := p.parseType()

	var end Pos
	if p.Token.Kind == ">>" {
		p.Token.Kind = ">"
		p.Token.Raw = ">"
		p.Token.Pos += 1
		end = p.Token.Pos
	} else {
		end = p.expect(">").End
	}
	return &ArrayType{
		pos:  pos,
		end:  end,
		Item: t,
	}
}

func (p *Parser) parseStructType() *StructType {
	pos := p.expect("STRUCT").Pos
	fields, end := p.parseStructTypeFields(true)
	return &StructType{
		pos:    pos,
		end:    end,
		Fields: fields,
	}
}

func (p *Parser) parseStructTypeFields(inType bool) (fields []*FieldType, end Pos) {
	if p.Token.Kind != "<" && p.Token.Kind != "<>" {
		if inType {
			p.panicfAtToken(&p.Token, "expected token: <, <>, but: %s", p.Token.Kind)
		}
		return
	}

	fields = make([]*FieldType, 0)
	if p.Token.Kind == "<>" {
		end = p.expect("<>").End
		return
	}

	p.expect("<")
	if p.Token.Kind != ">" && p.Token.Kind != ">>" {
		for p.Token.Kind != TokenEOF {
			fields = append(fields, p.parseFieldType())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}

	if p.Token.Kind == ">>" {
		p.Token.Kind = ">"
		p.Token.Raw = ">"
		p.Token.Pos += 1
		end = p.Token.Pos
	} else {
		end = p.expect(">").End
	}
	return
}

func (p *Parser) parseFieldType() *FieldType {
	lexer := p.Lexer.Clone()
	// Try to parse as "x INT64" case.
	if p.Token.Kind == TokenIdent {
		member := p.parseIdent()
		if p.lookaheadType() {
			t := p.parseType()
			return &FieldType{
				Member: member,
				Type:   t,
			}
		}
	}
	p.Lexer = lexer
	return &FieldType{
		Type: p.parseType(),
	}
}

func (p *Parser) lookaheadType() bool {
	return p.Token.Kind == TokenIdent || p.Token.Kind == "ARRAY" || p.Token.Kind == "STRUCT"
}

// ================================================================================
//
// DDL
//
// ================================================================================

func (p *Parser) parseDDL() DDL {
	pos := p.Token.Pos
	switch {
	case p.Token.Kind == "CREATE":
		p.NextToken()
		switch {
		case p.Token.IsKeywordLike("DATABASE"):
			return p.parseCreateDatabase(pos)
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseCreateTable(pos)
		case p.Token.IsKeywordLike("INDEX") || p.Token.IsKeywordLike("UNIQUE") || p.Token.IsKeywordLike("NULL_FILTERED"):
			return p.parseCreateIndex(pos)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: DATABASE, TABLE, INDEX, UNIQUE, NULL_FILTERED, but: %s", p.Token.AsString)
	case p.Token.IsKeywordLike("ALTER"):
		p.NextToken()
		return p.parseAlterTable(pos)
	case p.Token.IsKeywordLike("DROP"):
		p.NextToken()
		switch {
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseDropTable(pos)
		case p.Token.IsKeywordLike("INDEX"):
			return p.parseDropIndex(pos)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: TABLE, INDEX, but: %s", p.Token.AsString)
	}

	if p.Token.Kind != TokenIdent {
		panic(p.errorfAtToken(&p.Token, "expected token: CREATE, <ident>, but: %s", p.Token.Kind))
	}

	panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: ALTER, DROP, but: %s", p.Token.AsString))
}

func (p *Parser) parseCreateDatabase(pos Pos) *CreateDatabase {
	p.expectKeywordLike("DATABASE")
	name := p.parseIdent()
	return &CreateDatabase{
		pos:  pos,
		Name: name,
	}
}

func (p *Parser) parseCreateTable(pos Pos) *CreateTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()

	// This loop allows parsing trailing comma intentionally.
	// TODO: is this allowed by Spanner really?
	p.expect("(")
	var columns []*ColumnDef
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		columns = append(columns, p.parseColumnDef())
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	p.expect(")")

	p.expectKeywordLike("PRIMARY")
	p.expectKeywordLike("KEY")

	p.expect("(")
	var keys []*IndexKey
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		keys = append(keys, p.parseIndexKey())
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	end := p.expect(")").End

	cluster := p.tryParseCluster()
	if cluster != nil {
		end = cluster.End()
	}

	return &CreateTable{
		pos:         pos,
		end:         end,
		Name:        name,
		Columns:     columns,
		PrimaryKeys: keys,
		Cluster:     cluster,
	}
}

func (p *Parser) parseColumnDef() *ColumnDef {
	name := p.parseIdent()
	t := p.parseSchemaType()
	end := t.End()

	notNull := false
	if p.Token.Kind == "NOT" {
		p.expect("NOT")
		end = p.expect("NULL").End
		notNull = true
	}

	options := p.tryParseColumnDefOptions()
	if options != nil {
		end = options.End()
	}

	return &ColumnDef{
		end:     end,
		Name:    name,
		Type:    t,
		NotNull: notNull,
		Options: options,
	}
}

func (p *Parser) tryParseColumnDefOptions() *ColumnDefOptions {
	if !p.Token.IsKeywordLike("OPTIONS") {
		return nil
	}

	return p.parseColumnDefOptions()
}

func (p *Parser) parseColumnDefOptions() *ColumnDefOptions {
	pos := p.expectKeywordLike("OPTIONS").Pos

	p.expect("(")
	p.expectIdent("allow_commit_timestamp")
	p.expect("=")

	var allowCommitTimestamp bool
	switch p.Token.Kind {
	case "TRUE":
		allowCommitTimestamp = true
	case "NULL":
		allowCommitTimestamp = false
	default:
		p.panicfAtToken(&p.Token, "expected token: TRUE, NULL, but: %s", p.Token.Kind)
	}
	p.NextToken()

	end := p.expect(")").End
	return &ColumnDefOptions{
		pos:                  pos,
		end:                  end,
		AllowCommitTimestamp: allowCommitTimestamp,
	}
}

func (p *Parser) parseIndexKey() *IndexKey {
	name := p.parseIdent()
	end := name.End()

	var dir Direction
	switch p.Token.Kind {
	case "ASC":
		end = p.expect("ASC").End
		dir = DirectionAsc
	case "DESC":
		end = p.expect("DESC").End
		dir = DirectionDesc
	}

	return &IndexKey{
		end:  end,
		Name: name,
		Dir:  dir,
	}
}

func (p *Parser) tryParseCluster() *Cluster {
	if p.Token.Kind != "," {
		return nil
	}
	pos := p.expect(",").Pos

	p.expectKeywordLike("INTERLEAVE")
	p.expect("IN")
	p.expectKeywordLike("PARENT")
	name := p.parseIdent()
	end := name.End()

	var onDelete OnDeleteAction
	if p.Token.Kind == "ON" {
		p.expect("ON")
		p.expectKeywordLike("DELETE")
		switch p.Token.Kind {
		case TokenIdent:
			end = p.expectKeywordLike("CASCADE").End
			onDelete = OnDeleteCascade
		case "NO":
			p.NextToken()
			end = p.expectKeywordLike("ACTION").End
			onDelete = OnDeleteNoAction
		default:
			p.panicfAtToken(&p.Token, "expected token: NO, <ident>, but: %s", p.Token.Kind)
		}
	}

	return &Cluster{
		pos:       pos,
		end:       end,
		TableName: name,
		OnDelete:  onDelete,
	}
}

func (p *Parser) parseCreateIndex(pos Pos) *CreateIndex {
	unique := false
	if p.Token.IsKeywordLike("UNIQUE") {
		p.NextToken()
		unique = true
	}

	nullFiltered := false
	if p.Token.IsKeywordLike("NULL_FILTERED") {
		p.NextToken()
		nullFiltered = true
	}

	p.expectKeywordLike("INDEX")
	name := p.parseIdent()

	p.expect("ON")
	tableName := p.parseIdent()

	p.expect("(")
	var keys []*IndexKey
	for p.Token.Kind != TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		keys = append(keys, p.parseIndexKey())
		if p.Token.Kind != "," {
			break
		}
		p.NextToken()
	}
	end := p.expect(")").End

	storing := p.tryParseStoring()
	if storing != nil {
		end = storing.End()
	}

	interleaveIn := p.tryParseInterleaveIn()
	if interleaveIn != nil {
		end = interleaveIn.End()
	}

	return &CreateIndex{
		pos:          pos,
		end:          end,
		Unique:       unique,
		NullFiltered: nullFiltered,
		Name:         name,
		TableName:    tableName,
		Keys:         keys,
		Storing:      storing,
		InterleaveIn: interleaveIn,
	}
}

func (p *Parser) tryParseStoring() *Storing {
	if !p.Token.IsKeywordLike("STORING") {
		return nil
	}
	pos := p.expectKeywordLike("STORING").Pos

	p.expect("(")
	columns := []*Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.NextToken()
		columns = append(columns, p.parseIdent())
	}

	end := p.expect(")").End
	return &Storing{
		pos:     pos,
		end:     end,
		Columns: columns,
	}
}

func (p *Parser) tryParseInterleaveIn() *InterleaveIn {
	if p.Token.Kind != "," {
		return nil
	}
	pos := p.expect(",").Pos
	p.expectKeywordLike("INTERLEAVE")
	p.expect("IN")
	name := p.parseIdent()

	return &InterleaveIn{
		pos:       pos,
		TableName: name,
	}
}

func (p *Parser) parseAlterTable(pos Pos) *AlterTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()

	var alternation TableAlternation
	switch {
	case p.Token.IsKeywordLike("ADD"):
		alternation = p.parseAddColumn()
	case p.Token.IsKeywordLike("DROP"):
		alternation = p.parseDropColumn()
	case p.Token.Kind == "SET":
		alternation = p.parseSetOnDelete()
	case p.Token.IsKeywordLike("ALTER"):
		alternation = p.parseAlterColumn()
	default:
		if p.Token.Kind == TokenIdent {
			p.panicfAtToken(&p.Token, "expected pseuso keyword: ADD, ALTER, DROP, but: %s", p.Token.AsString)
		} else {
			p.panicfAtToken(&p.Token, "expected token: SET, <ident>, but: %s", p.Token.Kind)
		}
	}

	return &AlterTable{
		pos:              pos,
		Name:             name,
		TableAlternation: alternation,
	}
}

func (p *Parser) parseAddColumn() *AddColumn {
	pos := p.expectKeywordLike("ADD").Pos
	p.expectKeywordLike("COLUMN")

	column := p.parseColumnDef()

	return &AddColumn{
		pos:    pos,
		Column: column,
	}
}

func (p *Parser) parseDropColumn() *DropColumn {
	pos := p.expectKeywordLike("DROP").Pos
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	return &DropColumn{
		pos:  pos,
		Name: name,
	}
}

func (p *Parser) parseSetOnDelete() *SetOnDelete {
	pos := p.expect("SET").Pos
	p.expect("ON")
	p.expectKeywordLike("DELETE")

	var onDelete OnDeleteAction
	var end Pos
	switch p.Token.Kind {
	case TokenIdent:
		end = p.expectKeywordLike("CASCADE").End
		onDelete = OnDeleteCascade
	case "NO":
		p.NextToken()
		end = p.expectKeywordLike("ACTION").End
		onDelete = OnDeleteNoAction
	default:
		p.panicfAtToken(&p.Token, "expected token: NO, <ident>, but: %s", p.Token.Kind)
	}

	return &SetOnDelete{
		pos:      pos,
		end:      end,
		OnDelete: onDelete,
	}
}

func (p *Parser) parseAlterColumn() TableAlternation {
	pos := p.expectKeywordLike("ALTER").Pos
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	if p.Token.Kind == "SET" {
		p.NextToken()
		options := p.parseColumnDefOptions()
		return &AlterColumnSet{
			pos:     pos,
			Name:    name,
			Options: options,
		}
	}

	t := p.parseSchemaType()

	end := name.End()
	notNull := false
	if p.Token.Kind == "NOT" {
		p.expect("NOT")
		end = p.expect("NULL").End
		notNull = true
	}

	return &AlterColumn{
		pos:     pos,
		end:     end,
		Name:    name,
		Type:    t,
		NotNull: notNull,
	}
}

func (p *Parser) parseDropTable(pos Pos) *DropTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()
	return &DropTable{
		pos:  pos,
		Name: name,
	}
}

func (p *Parser) parseDropIndex(pos Pos) *DropIndex {
	p.expectKeywordLike("INDEX")
	name := p.parseIdent()
	return &DropIndex{
		pos:  pos,
		Name: name,
	}
}

func (p *Parser) parseSchemaType() SchemaType {
	switch p.Token.Kind {
	case TokenIdent:
		return p.parseScalarSchemaType()
	case "ARRAY":
		pos := p.expect("ARRAY").Pos
		p.expect("<")
		t := p.parseScalarSchemaType()
		end := p.expect(">").End
		return &ArraySchemaType{
			pos:  pos,
			end:  end,
			Item: t,
		}
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ARRAY, <ident>, but: %s", p.Token.Kind))
}

func (p *Parser) parseScalarSchemaType() SchemaType {
	id := p.expect(TokenIdent)
	pos := id.Pos
	switch {
	case id.IsIdent("BOOL"):
		return &ScalarSchemaType{
			pos:  pos,
			Name: BoolTypeName,
		}
	case id.IsIdent("INT64"):
		return &ScalarSchemaType{
			pos:  pos,
			Name: Int64TypeName,
		}
	case id.IsIdent("FLOAT64"):
		return &ScalarSchemaType{
			pos:  pos,
			Name: Float64TypeName,
		}
	case id.IsIdent("DATE"):
		return &ScalarSchemaType{
			pos:  pos,
			Name: DateTypeName,
		}
	case id.IsIdent("TIMESTAMP"):
		return &ScalarSchemaType{
			pos:  pos,
			Name: TimestampTypeName,
		}
	case id.IsIdent("STRING"):
		p.expect("(")
		max := false
		var size IntValue
		if p.Token.IsIdent("MAX") {
			p.NextToken()
			max = true
		} else {
			size = p.parseIntValue()
		}
		end := p.expect(")").End
		return &SizedSchemaType{
			pos:  pos,
			end:  end,
			Name: StringTypeName,
			Max:  max,
			Size: size,
		}
	case id.IsIdent("BYTES"):
		p.expect("(")
		max := false
		var size IntValue
		if p.Token.IsIdent("MAX") {
			p.NextToken()
			max = true
		} else {
			size = p.parseIntValue()
		}
		end := p.expect(")").End
		return &SizedSchemaType{
			pos:  pos,
			end:  end,
			Name: BytesTypeName,
			Max:  max,
			Size: size,
		}
	}

	panic(p.errorfAtToken(id, "expect ident: BOOL, INT64, FLOAT64, DATE, TIMESTAMP, STRING, BYTES, but: %s", id.AsString))
}

// ================================================================================
//
// DML
//
// ================================================================================

func (p *Parser) parseDML() DML {
	id := p.expect(TokenIdent)
	pos := id.Pos
	switch {
	case id.IsKeywordLike("INSERT"):
		return p.parseInsert(pos)
	case id.IsKeywordLike("DELETE"):
		return p.parseDelete(pos)
	case id.IsKeywordLike("UPDATE"):
		return p.parseUpdate(pos)
	}

	panic(p.errorfAtToken(id, "expect pseudo keyword: INSERT, DELETE,  UPDATE but: %s", id.AsString))
}

func (p *Parser) parseInsert(pos Pos) *Insert {
	if p.Token.Kind == "INTO" {
		p.NextToken()
	}

	name := p.parseIdent()

	p.expect("(")
	var columns []*Ident
	if p.Token.Kind != ")" {
		for p.Token.Kind != TokenEOF {
			columns = append(columns, p.parseIdent())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}
	p.expect(")")

	var input InsertInput
	if p.Token.IsKeywordLike("VALUES") {
		input = p.parseValuesInput()
	} else {
		input = p.parseSubQueryInput()
	}

	return &Insert{
		pos:       pos,
		TableName: name,
		Columns:   columns,
		Input:     input,
	}
}

func (p *Parser) parseValuesInput() *ValuesInput {
	pos := p.expectKeywordLike("VALUES").Pos

	rows := []*ValuesRow{p.parseValuesRow()}
	for p.Token.Kind == "," {
		p.NextToken()
		rows = append(rows, p.parseValuesRow())
	}

	return &ValuesInput{
		pos:  pos,
		Rows: rows,
	}
}

func (p *Parser) parseValuesRow() *ValuesRow {
	pos := p.expect("(").Pos
	var values []*DefaultExpr
	if p.Token.Kind != ")" {
		for p.Token.Kind != TokenEOF {
			values = append(values, p.parseDefaultExpr())
			if p.Token.Kind != "," {
				break
			}
			p.NextToken()
		}
	}
	end := p.expect(")").End

	return &ValuesRow{
		pos:    pos,
		end:    end,
		Values: values,
	}
}

func (p *Parser) parseDefaultExpr() *DefaultExpr {
	if p.Token.Kind == "DEFAULT" {
		pos := p.expect("DEFAULT").Pos
		return &DefaultExpr{
			pos:     pos,
			Default: true,
		}
	}

	expr := p.parseExpr()
	return &DefaultExpr{
		pos:  expr.Pos(),
		Expr: expr,
	}
}

func (p *Parser) parseSubQueryInput() *SubQueryInput {
	query := p.parseQueryExpr()

	return &SubQueryInput{
		Query: query,
	}
}

func (p *Parser) parseDelete(pos Pos) *Delete {
	if p.Token.Kind == "FROM" {
		p.NextToken()
	}

	name := p.parseIdent()
	as := p.tryParseAsAlias()
	where := p.parseWhere()

	return &Delete{
		pos:       pos,
		TableName: name,
		As:        as,
		Where:     where,
	}
}

func (p *Parser) parseUpdate(pos Pos) *Update {
	name := p.parseIdent()
	as := p.tryParseAsAlias()

	p.expect("SET")

	items := []*UpdateItem{p.parseUpdateItem()}
	for p.Token.Kind == "," {
		p.NextToken()
		items = append(items, p.parseUpdateItem())
	}

	where := p.parseWhere()

	return &Update{
		pos:         pos,
		TableName:   name,
		As:          as,
		UpdateItems: items,
		Where:       where,
	}
}

func (p *Parser) parseUpdateItem() *UpdateItem {
	path := []*Ident{p.parseIdent()}
	for p.Token.Kind == "." {
		p.NextToken()
		path = append(path, p.parseIdent())
	}

	p.expect("=")
	expr := p.parseExpr()

	return &UpdateItem{
		Path: path,
		Expr: expr,
	}
}

// ================================================================================
//
// Primitives
//
// ================================================================================

func (p *Parser) parseIdent() *Ident {
	id := p.expect(TokenIdent)
	return &Ident{
		pos:  id.Pos,
		end:  id.End,
		Name: id.AsString,
	}
}

func (p *Parser) parseParam() *Param {
	param := p.expect(TokenParam)
	return &Param{
		pos:  param.Pos,
		Name: param.AsString,
	}
}

func (p *Parser) parseNullLiteral() *NullLiteral {
	tok := p.expect("NULL")
	return &NullLiteral{
		pos: tok.Pos,
	}
}

func (p *Parser) parseBoolLiteral() *BoolLiteral {
	var value bool
	pos := p.Token.Pos
	switch p.Token.Kind {
	case "TRUE":
		value = true
	case "FALSE":
		value = false
	default:
		p.panicfAtToken(&p.Token, "expected token: TRUE, FALSE, but: %s", p.Token.Kind)
	}
	p.NextToken()
	return &BoolLiteral{
		pos:   pos,
		Value: value,
	}
}

func (p *Parser) parseIntLiteral() *IntLiteral {
	i := p.expect(TokenInt)
	return &IntLiteral{
		pos:   i.Pos,
		end:   i.End,
		Base:  i.Base,
		Value: i.Raw,
	}
}

func (p *Parser) parseFloatLiteral() *FloatLiteral {
	f := p.expect(TokenFloat)
	return &FloatLiteral{
		pos:   f.Pos,
		end:   f.End,
		Value: f.Raw,
	}
}

func (p *Parser) parseStringLiteral() *StringLiteral {
	s := p.expect(TokenString)
	return &StringLiteral{
		pos:   s.Pos,
		end:   s.End,
		Value: s.AsString,
	}
}

func (p *Parser) parseBytesLiteral() *BytesLiteral {
	b := p.expect(TokenBytes)
	return &BytesLiteral{
		pos:   b.Pos,
		end:   b.End,
		Value: []byte(b.AsString),
	}
}

func (p *Parser) parseIntValue() IntValue {
	switch p.Token.Kind {
	case TokenParam:
		return p.parseParam()
	case TokenInt:
		return p.parseIntLiteral()
	case "CAST":
		return p.parseCastIntValue()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <int>, CAST, but: %s", p.Token.Kind))
}

func (p *Parser) parseCastIntValue() *CastIntValue {
	pos := p.expect("CAST").Pos
	p.expect("(")
	var v IntValue
	switch p.Token.Kind {
	case TokenParam:
		v = p.parseParam()
	case TokenInt:
		v = p.parseIntLiteral()
	default:
		p.panicfAtToken(&p.Token, "expected token: <param>, <int>, but: %s", p.Token.Kind)
	}
	p.expect("AS")
	p.expectIdent("INT64")
	end := p.expect(")").End
	return &CastIntValue{
		pos:  pos,
		end:  end,
		Expr: v,
	}
}

func (p *Parser) parseNumValue() NumValue {
	switch p.Token.Kind {
	case TokenParam:
		return p.parseParam()
	case TokenInt:
		return p.parseIntLiteral()
	case TokenFloat:
		return p.parseFloatLiteral()
	case "CAST":
		return p.parseCastNumValue()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <int>, <float>, CAST, but: %s", p.Token.Kind))
}

func (p *Parser) parseCastNumValue() *CastNumValue {
	pos := p.expect("CAST").Pos
	p.expect("(")
	var v NumValue
	switch p.Token.Kind {
	case TokenParam:
		v = p.parseParam()
	case TokenInt:
		v = p.parseIntLiteral()
	case TokenFloat:
		v = p.parseFloatLiteral()
	default:
		p.panicfAtToken(&p.Token, "expected token: <param>, <int>, <float>, but: %s", p.Token.Kind)
	}
	p.expect("AS")
	id := p.expect(TokenIdent)
	var t ScalarTypeName
	switch {
	case id.IsIdent("INT64"):
		t = Int64TypeName
	case id.IsIdent("FLOAT64"):
		t = Float64TypeName
	default:
		p.panicfAtToken(id, "expected identifier: INT64, FLOAT64, but: %s", id.Raw)
	}
	end := p.expect(")").End
	return &CastNumValue{
		pos:  pos,
		end:  end,
		Expr: v,
		Type: t,
	}
}

func (p *Parser) parseStringValue() StringValue {
	switch p.Token.Kind {
	case TokenParam:
		return p.parseParam()
	case TokenString:
		return p.parseStringLiteral()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <string>, but: %s", p.Token.Kind))
}

func (p *Parser) expect(kind TokenKind) *Token {
	if p.Token.Kind != kind {
		p.panicfAtToken(&p.Token, "expected token: %s, but: %s", kind, p.Token.Kind)
	}
	t := p.Token.Clone()
	p.NextToken()
	return t
}

func (p *Parser) expectIdent(s string) *Token {
	id := p.expect(TokenIdent)
	if !id.IsIdent(s) {
		p.panicfAtToken(id, "expected identifier: %s, but: %s", s, QuoteSQLIdent(id.AsString))
	}
	return id
}

func (p *Parser) expectKeywordLike(s string) *Token {
	id := p.expect(TokenIdent)
	if !id.IsKeywordLike(s) {
		if strings.EqualFold(id.AsString, s) {
			p.panicfAtToken(id, "pseudo keyword %s cannot encloses with backquote", s)
		} else {
			p.panicfAtToken(id, "expected pseudo keyword: %s, but: %s", s, QuoteSQLIdent(id.AsString))
		}
	}
	return id
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
