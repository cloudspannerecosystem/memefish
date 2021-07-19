package parser

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/memefish/pkg/ast"
	"github.com/MakeNowJust/memefish/pkg/char"
	"github.com/MakeNowJust/memefish/pkg/token"
)

type Parser struct {
	*Lexer
}

// ParseStatement parses a SQL statement.
func (p *Parser) ParseStatement() (stmt ast.Statement, err error) {
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

	p.nextToken()
	stmt = p.parseStatement()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseStatements parses SQL statements list separated by semi-colon.
func (p *Parser) ParseStatements() (stmts []ast.Statement, err error) {
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

	p.nextToken()
	p.parseStatements(func() {
		stmts = append(stmts, p.parseStatement())
	})
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseQuery parses a SELECT query statement.
func (p *Parser) ParseQuery() (stmt *ast.QueryStatement, err error) {
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

	p.nextToken()
	stmt = p.parseQueryStatement()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseExpr parses a SQL expression.
func (p *Parser) ParseExpr() (expr ast.Expr, err error) {
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

	p.nextToken()
	expr = p.parseExpr()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDDL parses a CREATE/ALTER/DROP statement.
func (p *Parser) ParseDDL() (ddl ast.DDL, err error) {
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

	p.nextToken()
	ddl = p.parseDDL()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDDLs parses CREATE/ALTER/DROP statements list separated by semi-colon.
func (p *Parser) ParseDDLs() (ddls []ast.DDL, err error) {
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

	p.nextToken()
	p.parseStatements(func() {
		ddls = append(ddls, p.parseDDL())
	})
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDML parses a INSERT/DELETE/UPDATE statement.
func (p *Parser) ParseDML() (dml ast.DML, err error) {
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

	p.nextToken()
	dml = p.parseDML()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

// ParseDMLs parses INSERT/DELETE/UPDATE statements list separated by semi-colon.
func (p *Parser) ParseDMLs() (dmls []ast.DML, err error) {
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

	p.nextToken()
	p.parseStatements(func() {
		dmls = append(dmls, p.parseDML())
	})
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}

func (p *Parser) parseStatement() ast.Statement {
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
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ";" {
			p.nextToken()
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

func (p *Parser) parseQueryStatement() *ast.QueryStatement {
	hint := p.tryParseHint()
	query := p.parseQueryExpr()

	return &ast.QueryStatement{
		Hint:  hint,
		Query: query,
	}
}

func (p *Parser) tryParseHint() *ast.Hint {
	if p.Token.Kind != "@" {
		return nil
	}

	atmark := p.Token.Pos
	p.nextToken()
	p.expect("{")
	records := []*ast.HintRecord{p.parseHintRecord()}
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
		records = append(records, p.parseHintRecord())
	}
	rbrace := p.expect("}").Pos
	return &ast.Hint{
		Atmark:  atmark,
		Rbrace:  rbrace,
		Records: records,
	}
}

func (p *Parser) parseHintRecord() *ast.HintRecord {
	key := p.parseIdent()
	p.expect("=")
	value := p.parseExpr()
	return &ast.HintRecord{
		Key:   key,
		Value: value,
	}
}

func (p *Parser) parseQueryExpr() ast.QueryExpr {
	query := p.parseSimpleQueryExpr()

	for {
		var op ast.SetOp
		switch p.Token.Kind {
		case "UNION":
			op = ast.SetOpUnion
		case "INTERSECT":
			op = ast.SetOpIntersect
		case "EXCEPT":
			op = ast.SetOpExcept
		}

		if op == "" {
			break
		}

		opTok := p.Token
		p.nextToken()

		var distinct bool
		switch p.Token.Kind {
		case "ALL":
			distinct = false
		case "DISTINCT":
			distinct = true
		default:
			p.panicfAtToken(&p.Token, "expected token: ALL, DISTINCT, but: %s", p.Token.Kind)
		}
		p.nextToken()

		right := p.parseSimpleQueryExpr()
		if c, ok := query.(*ast.CompoundQuery); ok {
			if !(c.Op == op && c.Distinct == distinct) {
				p.panicfAtToken(&opTok, "all set operator at the same level must be the same, or wrap (...)")
			}
			c.Queries = append(c.Queries, right)
		} else {
			query = &ast.CompoundQuery{
				Op:       op,
				Distinct: distinct,
				Queries:  []ast.QueryExpr{query, right},
			}
		}
	}

	return p.parseQueryExprSuffix(query)
}

func (p *Parser) parseSimpleQueryExpr() ast.QueryExpr {
	if p.Token.Kind == "(" {
		lparen := p.expect("(").Pos
		query := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		return &ast.SubQuery{
			Lparen: lparen,
			Rparen: rparen,
			Query:  query,
		}
	}

	return p.parseSelect()
}

func (p *Parser) parseSelect() *ast.Select {
	sel := p.expect("SELECT").Pos
	var distinct bool
	if p.Token.Kind == "DISTINCT" {
		p.nextToken()
		distinct = true
	}
	var asStruct bool
	if p.Token.Kind == "AS" {
		p.nextToken()
		p.expect("STRUCT")
		asStruct = true
	}

	results := p.parseSelectResults()
	from := p.tryParseFrom()
	where := p.tryParseWhere()
	groupBy := p.tryParseGroupBy()
	having := p.tryParseHaving()

	return &ast.Select{
		Select:   sel,
		Distinct: distinct,
		AsStruct: asStruct,
		Results:  results,
		From:     from,
		Where:    where,
		GroupBy:  groupBy,
		Having:   having,
	}
}

func (p *Parser) parseSelectResults() []ast.SelectItem {
	results := []ast.SelectItem{p.parseSelectItem()}
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
		results = append(results, p.parseSelectItem())
	}
	return results
}

func (p *Parser) parseSelectItem() ast.SelectItem {
	if p.Token.Kind == "*" {
		pos := p.expect("*").Pos
		return &ast.Star{
			Star: pos,
		}
	}

	expr := p.parseExpr()
	if as := p.tryParseAsAlias(); as != nil {
		return &ast.Alias{
			Expr: expr,
			As:   as,
		}
	}

	if p.Token.Kind == "." {
		p.nextToken()
		pos := p.expect("*").Pos
		return &ast.DotStar{
			Star: pos,
			Expr: expr,
		}
	}

	return &ast.ExprSelectItem{
		Expr: expr,
	}
}

func (p *Parser) tryParseAsAlias() *ast.AsAlias {
	pos := p.Token.Pos

	if p.Token.Kind == "AS" {
		p.nextToken()
		id := p.parseIdent()
		return &ast.AsAlias{
			As:    token.InvalidPos,
			Alias: id,
		}
	}

	if p.Token.Kind == token.TokenIdent {
		id := p.parseIdent()
		return &ast.AsAlias{
			As:    pos,
			Alias: id,
		}
	}

	return nil
}

func (p *Parser) tryParseFrom() *ast.From {
	if p.Token.Kind != "FROM" {
		return nil
	}
	from := p.expect("FROM").Pos

	join := p.parseTableExpr(
		/* toplevel */ true,
	)

	return &ast.From{
		From:   from,
		Source: join,
	}
}

func (p *Parser) tryParseWhere() *ast.Where {
	if p.Token.Kind != "WHERE" {
		return nil
	}
	return p.parseWhere()
}

func (p *Parser) parseWhere() *ast.Where {
	pos := p.expect("WHERE").Pos
	expr := p.parseExpr()
	return &ast.Where{
		Where: pos,
		Expr:  expr,
	}
}

func (p *Parser) tryParseGroupBy() *ast.GroupBy {
	if p.Token.Kind != "GROUP" {
		return nil
	}
	pos := p.expect("GROUP").Pos
	p.expect("BY")
	exprs := []ast.Expr{p.parseExpr()}
	for p.Token.Kind == "," {
		p.nextToken()
		exprs = append(exprs, p.parseExpr())
	}

	return &ast.GroupBy{
		Group: pos,
		Exprs: exprs,
	}
}

func (p *Parser) tryParseHaving() *ast.Having {
	if p.Token.Kind != "HAVING" {
		return nil
	}
	pos := p.expect("HAVING").Pos
	expr := p.parseExpr()
	return &ast.Having{
		Having: pos,
		Expr:   expr,
	}
}

func (p *Parser) parseQueryExprSuffix(e ast.QueryExpr) ast.QueryExpr {
	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()

	switch e := e.(type) {
	case *ast.Select:
		e.OrderBy = orderBy
		e.Limit = limit
	case *ast.SubQuery:
		e.OrderBy = orderBy
		e.Limit = limit
	case *ast.CompoundQuery:
		e.OrderBy = orderBy
		e.Limit = limit
	}

	return e
}

func (p *Parser) tryParseOrderBy() *ast.OrderBy {
	if p.Token.Kind != "ORDER" {
		return nil
	}

	pos := p.expect("ORDER").Pos
	p.expect("BY")

	items := []*ast.OrderByItem{p.parseOrderByItem()}
	for p.Token.Kind == "," {
		p.nextToken()
		items = append(items, p.parseOrderByItem())
	}

	return &ast.OrderBy{
		Order: pos,
		Items: items,
	}
}

func (p *Parser) parseOrderByItem() *ast.OrderByItem {
	expr := p.parseExpr()
	collate := p.tryParseCollate()
	dir, dirPos := p.tryParseDirection()

	return &ast.OrderByItem{
		DirPos:  dirPos,
		Expr:    expr,
		Collate: collate,
		Dir:     dir,
	}
}

func (p *Parser) tryParseCollate() *ast.Collate {
	if p.Token.Kind != "COLLATE" {
		return nil
	}
	pos := p.expect("COLLATE").Pos
	value := p.parseStringValue()
	return &ast.Collate{
		Collate: pos,
		Value:   value,
	}
}

func (p *Parser) tryParseDirection() (ast.Direction, token.Pos) {
	var dir ast.Direction
	dirPos := token.InvalidPos
	switch p.Token.Kind {
	case "ASC":
		dirPos = p.expect("ASC").Pos
		dir = ast.DirectionAsc
	case "DESC":
		dirPos = p.expect("DESC").Pos
		dir = ast.DirectionDesc
	}

	return dir, dirPos
}

func (p *Parser) tryParseLimit() *ast.Limit {
	if p.Token.Kind != "LIMIT" {
		return nil
	}

	pos := p.expect("LIMIT").Pos
	count := p.parseIntValue()
	offset := p.tryParseOffset()

	return &ast.Limit{
		Limit:  pos,
		Count:  count,
		Offset: offset,
	}
}

func (p *Parser) tryParseOffset() *ast.Offset {
	if !p.Token.IsKeywordLike("OFFSET") {
		return nil
	}
	pos := p.expectKeywordLike("OFFSET").Pos
	value := p.parseIntValue()
	return &ast.Offset{
		Offset: pos,
		Value:  value,
	}
}

// ================================================================================
//
// JOIN
//
// ================================================================================

func (p *Parser) parseTableExpr(toplevel bool) ast.TableExpr {
	needJoin := !toplevel
	join := p.parseSimpleTableExpr()
	for {
		if needJoin {
			switch j := join.(type) {
			case *ast.Join:
				needJoin = false
			case *ast.SubQueryTableExpr:
				needJoin = j.As == nil
			}
		}

		op := ast.InnerJoin
		switch p.Token.Kind {
		case "INNER":
			p.nextToken()
			needJoin = true
		case "CROSS":
			p.nextToken()
			op = ast.CrossJoin
			needJoin = true
		case "FULL":
			p.nextToken()
			if p.Token.Kind == "OUTER" {
				p.nextToken()
			}
			op = ast.FullOuterJoin
			needJoin = true
		case "LEFT":
			p.nextToken()
			if p.Token.Kind == "OUTER" {
				p.nextToken()
			}
			op = ast.LeftOuterJoin
			needJoin = true
		case "RIGHT":
			p.nextToken()
			if p.Token.Kind == "OUTER" {
				p.nextToken()
			}
			op = ast.RightOuterJoin
			needJoin = true
		}

		if toplevel && p.Token.Kind == "," {
			op = ast.CommaJoin
		}

		var method ast.JoinMethod
		if op != ast.CommaJoin {
			switch {
			case p.Token.Kind == "HASH":
				p.nextToken()
				method = ast.HashJoinMethod
				needJoin = true
			case p.Token.IsKeywordLike("APPLY"):
				p.nextToken()
				method = ast.ApplyJoinMethod
				needJoin = true
			case p.Token.IsKeywordLike("LOOP"):
				p.nextToken()
				method = ast.LoopJoinMethod
				needJoin = true
			}
		}

		switch {
		case needJoin:
			p.expect("JOIN")
			needJoin = false
		case op == ast.CommaJoin || p.Token.Kind == "JOIN":
			p.nextToken()
		default:
			return join
		}

		hint := p.tryParseHint()
		right := p.parseSimpleTableExpr()

		if op == ast.CrossJoin || op == ast.CommaJoin {
			join = &ast.Join{
				Op:     op,
				Method: method,
				Hint:   hint,
				Left:   join,
				Right:  right,
			}
			continue
		}

		cond := p.parseJoinCondition()
		join = &ast.Join{
			Op:     op,
			Method: method,
			Hint:   hint,
			Left:   join,
			Right:  right,
			Cond:   cond,
		}
	}
}

func (p *Parser) parseSimpleTableExpr() ast.TableExpr {
	if p.lookaheadSubQuery() {
		lparen := p.expect("(").Pos
		query := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		as := p.tryParseAsAlias()
		return p.parseTableExprSuffix(&ast.SubQueryTableExpr{
			Lparen: lparen,
			Rparen: rparen,
			Query:  query,
			As:     as,
		})
	}

	if p.Token.Kind == "(" {
		lparen := p.expect("(").Pos
		join := p.parseTableExpr(
			/* toplevel */ false,
		)
		rparen := p.expect(")").Pos
		return p.parseTableExprSuffix(&ast.ParenTableExpr{
			Lparen: lparen,
			Rparen: rparen,
			Source: join,
		})
	}

	if p.Token.Kind == "UNNEST" {
		unnest := p.expect("UNNEST").Pos
		p.expect("(")
		expr := p.parseExpr()
		rparen := p.expect(")").Pos
		return p.parseUnnestSuffix(false, expr, unnest, rparen)
	}

	if p.Token.Kind == token.TokenIdent {
		ids := p.parseIdentOrPath()
		if len(ids) == 1 {
			return p.parseTableNameSuffix(ids[0])
		}
		return p.parseUnnestSuffix(true, &ast.Path{Idents: ids}, token.InvalidPos, token.InvalidPos)
	}

	panic(p.errorfAtToken(&p.Token, "expected token: (, UNNEST, <ident>, but: %s", p.Token.Kind))
}

func (p *Parser) parseIdentOrPath() []*ast.Ident {
	ids := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "." {
		p.nextToken()
		ids = append(ids, p.parseIdent())
	}
	return ids
}

func (p *Parser) parseUnnestSuffix(implicit bool, expr ast.Expr, unnest, rparen token.Pos) ast.TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias()
	withOffset := p.tryParseWithOffset()

	return p.parseTableExprSuffix(&ast.Unnest{
		Unnest:     unnest,
		Rparen:     rparen,
		Implicit:   implicit,
		Expr:       expr,
		Hint:       hint,
		As:         as,
		WithOffset: withOffset,
	})
}

func (p *Parser) tryParseWithOffset() *ast.WithOffset {
	if p.Token.Kind != "WITH" {
		return nil
	}

	with := p.expect("WITH").Pos
	offset := p.expectKeywordLike("OFFSET").Pos
	as := p.tryParseAsAlias()

	return &ast.WithOffset{
		With:   with,
		Offset: offset,
		As:     as,
	}
}

func (p *Parser) parseTableNameSuffix(id *ast.Ident) ast.TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias()
	return p.parseTableExprSuffix(&ast.TableName{
		Table: id,
		Hint:  hint,
		As:    as,
	})
}

func (p *Parser) parseJoinCondition() ast.JoinCondition {
	switch p.Token.Kind {
	case "ON":
		return p.parseOn()
	case "USING":
		return p.parseUsing()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ON, USING, but: %s", p.Token.Kind))
}

func (p *Parser) parseOn() *ast.On {
	pos := p.expect("ON").Pos
	expr := p.parseExpr()
	return &ast.On{
		On:   pos,
		Expr: expr,
	}
}

func (p *Parser) parseUsing() *ast.Using {
	using := p.expect("USING").Pos
	p.expect("(")
	idents := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.nextToken()
		idents = append(idents, p.parseIdent())
	}
	rparen := p.expect(")").Pos
	return &ast.Using{
		Using:  using,
		Rparen: rparen,
		Idents: idents,
	}
}

func (p *Parser) parseTableExprSuffix(join ast.TableExpr) ast.TableExpr {
	sample := p.tryParseTableSample()
	switch j := join.(type) {
	case *ast.Unnest:
		j.Sample = sample
	case *ast.TableName:
		j.Sample = sample
	case *ast.SubQueryTableExpr:
		j.Sample = sample
	case *ast.ParenTableExpr:
		j.Sample = sample
	default:
		panic(fmt.Sprintf("BUG: unexpected join: %#v", join))
	}
	return join
}

func (p *Parser) tryParseTableSample() *ast.TableSample {
	if p.Token.Kind != "TABLESAMPLE" {
		return nil
	}
	pos := p.expect("TABLESAMPLE").Pos

	id := p.expect(token.TokenIdent)
	var method ast.TableSampleMethod
	switch {
	case id.IsIdent("BERNOULLI"):
		method = ast.BernoulliSampleMethod
	case id.IsIdent("RESERVOIR"):
		method = ast.ReservoirSampleMethod
	default:
		p.panicfAtToken(id, "expected identifier: BERNOULLI, RESERVOIR, but: %s", id.Raw)
	}

	size := p.parseTableSampleSize()

	return &ast.TableSample{
		TableSample: pos,
		Method:      method,
		Size:        size,
	}
}

func (p *Parser) parseTableSampleSize() *ast.TableSampleSize {
	lparen := p.expect("(").Pos

	value := p.parseNumValue()
	var unit ast.TableSampleUnit
	switch {
	case p.Token.Kind == "ROWS":
		unit = ast.RowsTableSampleUnit
	case p.Token.IsKeywordLike("PERCENT"):
		unit = ast.PercentTableSampleUnit
	default:
		p.panicfAtToken(&p.Token, "expected token: ROWS, PERCENT, but: %s", p.Token.Kind)
	}
	p.nextToken()

	rparen := p.expect(")").Pos
	return &ast.TableSampleSize{
		Lparen: lparen,
		Rparen: rparen,
		Value:  value,
		Unit:   unit,
	}
}

// ================================================================================
//
// Expr
//
// ================================================================================

func (p *Parser) parseExpr() ast.Expr {
	return p.parseAndOr()
}

func (p *Parser) parseAndOr() ast.Expr {
	expr := p.parseNot()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "AND":
			op = ast.OpAnd
		case "OR":
			op = ast.OpOr
		default:
			return expr
		}
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseNot(),
		}
	}
}

func (p *Parser) parseNot() ast.Expr {
	if p.Token.Kind == "NOT" {
		pos := p.Token.Pos
		p.nextToken()
		return &ast.UnaryExpr{
			OpPos: pos,
			Op:    ast.OpNot,
			Expr:  p.parseNot(),
		}
	}

	return p.parseComparison()
}

func (p *Parser) parseComparison() ast.Expr {
	expr := p.parseBitOr()
	var op ast.BinaryOp
	switch p.Token.Kind {
	case "<":
		op = ast.OpLess
	case ">":
		op = ast.OpGreater
	case "<=":
		op = ast.OpLessEqual
	case ">=":
		op = ast.OpGreaterEqual
	case "=":
		op = ast.OpEqual
	case "!=", "<>":
		op = ast.OpNotEqual
	case "LIKE":
		op = ast.OpLike
	case "IN":
		p.nextToken()
		cond := p.parseInCondition()
		return &ast.InExpr{
			Left:  expr,
			Right: cond,
		}
	case "BETWEEN":
		p.nextToken()
		rightStart := p.parseBitOr()
		p.expect("AND")
		rightEnd := p.parseBitOr()
		return &ast.BetweenExpr{
			Left:       expr,
			RightStart: rightStart,
			RightEnd:   rightEnd,
		}
	case "NOT":
		p.nextToken()
		switch p.Token.Kind {
		case "LIKE":
			op = ast.OpNotLike
		case "IN":
			p.nextToken()
			cond := p.parseInCondition()
			return &ast.InExpr{
				Not:   true,
				Left:  expr,
				Right: cond,
			}
		case "BETWEEN":
			p.nextToken()
			rightStart := p.parseBitOr()
			p.expect("AND")
			rightEnd := p.parseBitOr()
			return &ast.BetweenExpr{
				Not:        true,
				Left:       expr,
				RightStart: rightStart,
				RightEnd:   rightEnd,
			}
		default:
			p.panicfAtToken(&p.Token, "expected token: LIKE, IN, but: %s", p.Token.Kind)
		}
	case "IS":
		p.nextToken()
		not := false
		if p.Token.Kind == "NOT" {
			p.nextToken()
			not = true
		}
		pos := p.Token.Pos
		switch p.Token.Kind {
		case "NULL":
			p.nextToken()
			return &ast.IsNullExpr{
				Null: pos,
				Left: expr,
				Not:  not,
			}
		case "TRUE":
			p.nextToken()
			return &ast.IsBoolExpr{
				RightPos: pos,
				Left:     expr,
				Not:      not,
				Right:    true,
			}
		case "FALSE":
			p.nextToken()
			return &ast.IsBoolExpr{
				RightPos: pos,
				Left:     expr,
				Not:      not,
				Right:    false,
			}
		default:
			p.panicfAtToken(&p.Token, "expected token: NULL, TRUE, FALSE, but: %s", p.Token.Kind)
		}
	default:
		return expr
	}
	p.nextToken()
	return &ast.BinaryExpr{
		Left:  expr,
		Op:    op,
		Right: p.parseBitOr(),
	}
}

func (p *Parser) parseInCondition() ast.InCondition {
	if p.lookaheadSubQuery() {
		lparen := p.expect("(").Pos
		query := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		return &ast.SubQueryInCondition{
			Lparen: lparen,
			Rparen: rparen,
			Query:  query,
		}
	}

	if p.Token.Kind == "(" {
		lparen := p.Token.Pos
		p.nextToken()
		exprs := []ast.Expr{p.parseExpr()}
		for p.Token.Kind != token.TokenEOF {
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
			exprs = append(exprs, p.parseExpr())
		}
		rparen := p.expect(")").Pos
		return &ast.ValuesInCondition{
			Lparen: lparen,
			Rparen: rparen,
			Exprs:  exprs,
		}
	}

	if p.Token.Kind == "UNNEST" {
		unnest := p.Token.Pos
		p.nextToken()
		p.expect("(")
		e := p.parseExpr()
		rparen := p.expect(")").Pos
		return &ast.UnnestInCondition{
			Unnest: unnest,
			Rparen: rparen,
			Expr:   e,
		}
	}

	panic(p.errorfAtToken(&p.Token, "expected token (, UNNEST, but: %s", p.Token.Kind))
}

func (p *Parser) parseBitOr() ast.Expr {
	expr := p.parseBitXor()
	for p.Token.Kind == "|" {
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    ast.OpBitOr,
			Right: p.parseBitXor(),
		}
	}
	return expr
}

func (p *Parser) parseBitXor() ast.Expr {
	expr := p.parseBitAnd()
	for p.Token.Kind == "^" {
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    ast.OpBitXor,
			Right: p.parseBitAnd(),
		}
	}
	return expr
}

func (p *Parser) parseBitAnd() ast.Expr {
	expr := p.parseBitShift()
	for p.Token.Kind == "&" {
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    ast.OpBitAnd,
			Right: p.parseBitShift(),
		}
	}
	return expr
}

func (p *Parser) parseBitShift() ast.Expr {
	expr := p.parseAddSub()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "<<":
			op = ast.OpBitLeftShift
		case ">>":
			op = ast.OpBitRightShift
		default:
			return expr
		}
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseAddSub(),
		}
	}
}

func (p *Parser) parseAddSub() ast.Expr {
	expr := p.parseMulDiv()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "+":
			op = ast.OpAdd
		case "-":
			op = ast.OpSub
		default:
			return expr
		}
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseMulDiv(),
		}
	}
}

func (p *Parser) parseMulDiv() ast.Expr {
	expr := p.parseUnary()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "*":
			op = ast.OpMul
		case "/":
			op = ast.OpDiv
		case "||":
			op = ast.OpConcat
		default:
			return expr
		}
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseUnary(),
		}
	}
}

func (p *Parser) parseUnary() ast.Expr {
	var op ast.UnaryOp
	switch p.Token.Kind {
	case "+":
		op = ast.OpPlus
	case "-":
		op = ast.OpMinus
	case "~":
		op = ast.OpBitNot
	default:
		return p.parseSelector()
	}
	pos := p.Token.Pos
	p.nextToken()

	e := p.parseUnary()
	if op != ast.OpBitNot {
		switch e := e.(type) {
		case *ast.IntLiteral:
			if e.Value[0] != '+' && e.Value[0] != '-' {
				e.ValuePos = pos
				e.Value = string(op) + e.Value
				return e
			}
		case *ast.FloatLiteral:
			if e.Value[0] != '+' && e.Value[0] != '-' {
				e.ValuePos = pos
				e.Value = string(op) + e.Value
				return e
			}
		}
	}

	return &ast.UnaryExpr{
		OpPos: pos,
		Op:    op,
		Expr:  e,
	}
}

func (p *Parser) parseSelector() ast.Expr {
	expr := p.parseLit()
	for {
		switch p.Token.Kind {
		case ".":
			lexer := p.Lexer.Clone()
			p.nextToken()
			if p.Token.Kind == "*" { // expr.* case
				p.Lexer = lexer
				return expr
			}

			ident := p.parseIdent()
			switch e := expr.(type) {
			case *ast.Ident:
				expr = &ast.Path{
					Idents: []*ast.Ident{e, ident},
				}
			case *ast.Path:
				e.Idents = append(e.Idents, ident)
			default:
				expr = &ast.SelectorExpr{
					Expr:  expr,
					Ident: ident,
				}
			}
		case "[":
			p.nextToken()
			id := p.expect(token.TokenIdent)
			ordinal := false
			if char.EqualFold(id.AsString, "ORDINAL") {
				ordinal = true
			} else if char.EqualFold(id.AsString, "OFFSET") {
				ordinal = false
			} else {
				p.panicfAtToken(id, "expected identifier: ORDINAL, OFFSET, but: %s", id.Raw)
			}
			p.expect("(")
			index := p.parseExpr()
			p.expect(")")
			rbrack := p.expect("]").Pos
			expr = &ast.IndexExpr{
				Rbrack:  rbrack,
				Ordinal: ordinal,
				Expr:    expr,
				Index:   index,
			}
		default:
			return expr
		}
	}
}

func (p *Parser) parseLit() ast.Expr {
	switch p.Token.Kind {
	case "NULL":
		return p.parseNullLiteral()
	case "TRUE", "FALSE":
		return p.parseBoolLiteral()
	case token.TokenInt:
		return p.parseIntLiteral()
	case token.TokenFloat:
		return p.parseFloatLiteral()
	case token.TokenString:
		return p.parseStringLiteral()
	case token.TokenBytes:
		return p.parseBytesLiteral()
	case token.TokenParam:
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
	case token.TokenIdent:
		id := p.Token
		p.nextToken()
		switch p.Token.Kind {
		case "(":
			return p.parseCall(id)
		case token.TokenString:
			if id.IsKeywordLike("DATE") {
				return p.parseDateLiteral(id)
			}
			if id.IsKeywordLike("TIMESTAMP") {
				return p.parseTimestampLiteral(id)
			}
			if id.IsKeywordLike("NUMERIC") {
				return p.parseNumericLiteral(id)
			}
		}
		return &ast.Ident{
			NamePos: id.Pos,
			NameEnd: id.End,
			Name:    id.AsString,
		}
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseCall(id token.Token) ast.Expr {
	p.expect("(")
	if id.IsIdent("COUNT") && p.Token.Kind == "*" {
		p.nextToken()
		rparen := p.expect(")").Pos
		return &ast.CountStarExpr{
			Count:  id.Pos,
			Rparen: rparen,
		}
	}

	fn := &ast.Ident{
		NamePos: id.Pos,
		NameEnd: id.End,
		Name:    id.AsString,
	}

	distinct := false
	if p.Token.Kind == "DISTINCT" {
		p.nextToken()
		distinct = true
	}

	var args []*ast.Arg
	if p.Token.Kind != ")" {
		for p.Token.Kind != token.TokenEOF {
			args = append(args, p.parseArg())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}
	rparen := p.expect(")").Pos
	return &ast.CallExpr{
		Rparen:   rparen,
		Func:     fn,
		Distinct: distinct,
		Args:     args,
	}
}

func (p *Parser) parseArg() *ast.Arg {
	if p.Token.Kind != "INTERVAL" {
		e := p.parseExpr()
		return &ast.Arg{
			Interval: token.InvalidPos,
			Expr:     e,
		}
	}

	pos := p.Token.Pos
	p.nextToken()
	e := p.parseExpr()
	unit := p.parseIdent()
	return &ast.Arg{
		Interval:     pos,
		Expr:         e,
		IntervalUnit: unit,
	}
}

func (p *Parser) parseCaseExpr() *ast.CaseExpr {
	pos := p.expect("CASE").Pos
	var expr ast.Expr
	if p.Token.Kind != "WHEN" {
		expr = p.parseExpr()
	}
	whens := []*ast.CaseWhen{p.parseCaseWhen()}
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind != "WHEN" {
			break
		}
		whens = append(whens, p.parseCaseWhen())
	}
	var els *ast.CaseElse
	if p.Token.Kind == "ELSE" {
		els = p.parseCaseElse()
	}
	end := p.expect("END").Pos
	return &ast.CaseExpr{
		Case:   pos,
		EndPos: end,
		Expr:   expr,
		Whens:  whens,
		Else:   els,
	}
}

func (p *Parser) parseCaseWhen() *ast.CaseWhen {
	pos := p.expect("WHEN").Pos
	cond := p.parseExpr()
	p.expect("THEN")
	then := p.parseExpr()
	return &ast.CaseWhen{
		When: pos,
		Cond: cond,
		Then: then,
	}
}

func (p *Parser) parseCaseElse() *ast.CaseElse {
	pos := p.expect("ELSE").Pos
	expr := p.parseExpr()
	return &ast.CaseElse{
		Else: pos,
		Expr: expr,
	}
}

func (p *Parser) parseCastExpr() *ast.CastExpr {
	cast := p.expect("CAST").Pos
	p.expect("(")
	e := p.parseExpr()
	p.expect("AS")
	t := p.parseType()
	rparen := p.expect(")").Pos
	return &ast.CastExpr{
		Cast:   cast,
		Rparen: rparen,
		Expr:   e,
		Type:   t,
	}
}

func (p *Parser) parseExistsSubQuery() *ast.ExistsSubQuery {
	exists := p.expect("EXISTS").Pos
	p.expect("(")
	query := p.parseQueryExpr()
	rparen := p.expect(")").Pos
	return &ast.ExistsSubQuery{
		Exists: exists,
		Rparen: rparen,
		Query:  query,
	}
}

func (p *Parser) parseExtractExpr() *ast.ExtractExpr {
	extract := p.expect("EXTRACT").Pos
	p.expect("(")
	part := p.parseIdent()
	p.expect("FROM")
	e := p.parseExpr()
	atTimeZone := p.tryParseAtTimeZone()
	rparen := p.expect(")").Pos
	return &ast.ExtractExpr{
		Extract:    extract,
		Rparen:     rparen,
		Part:       part,
		Expr:       e,
		AtTimeZone: atTimeZone,
	}
}

func (p *Parser) tryParseAtTimeZone() *ast.AtTimeZone {
	if p.Token.Kind != "AT" {
		return nil
	}

	pos := p.expect("AT").Pos
	p.expectKeywordLike("TIME")
	p.expectKeywordLike("ZONE")
	e := p.parseExpr()
	return &ast.AtTimeZone{
		At:   pos,
		Expr: e,
	}
}

func (p *Parser) parseParenExpr() ast.Expr {
	paren := p.Token

	if p.lookaheadSubQuery() {
		p.nextToken()
		query := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		return &ast.ScalarSubQuery{
			Lparen: paren.Pos,
			Rparen: rparen,
			Query:  query,
		}
	}

	p.nextToken()
	expr := p.parseExpr()

	if p.Token.Kind == ")" {
		rparen := p.Token.Pos
		p.nextToken()
		return &ast.ParenExpr{
			Lparen: paren.Pos,
			Rparen: rparen,
			Expr:   expr,
		}
	}

	if p.Token.Kind != "," {
		p.panicfAtToken(&paren, "cannot parse (...) as expression, struct literal or subquery")
	}

	values := []ast.Expr{expr}
	for p.Token.Kind == "," {
		p.nextToken()
		values = append(values, p.parseExpr())
	}
	rparen := p.expect(")").Pos
	return &ast.StructLiteral{
		Lparen: paren.Pos,
		Rparen: rparen,
		Values: values,
	}
}

func (p *Parser) parseArrayLiteralOrSubQuery() ast.Expr {
	pos := p.expect("ARRAY").Pos

	if p.Token.Kind == "(" {
		p.nextToken()
		query := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		return &ast.ArraySubQuery{
			Array:  pos,
			Rparen: rparen,
			Query:  query,
		}
	}

	var t ast.Type
	if p.Token.Kind == "<" {
		p.nextToken()
		t = p.parseType()
		p.expect(">")
	}

	values, lbrack, rbrack := p.parseArrayLiteralBody()
	return &ast.ArrayLiteral{
		Array:  pos,
		Lbrack: lbrack,
		Rbrack: rbrack,
		Type:   t,
		Values: values,
	}
}

func (p *Parser) parseSimpleArrayLiteral() *ast.ArrayLiteral {
	values, lbrack, rbrack := p.parseArrayLiteralBody()
	return &ast.ArrayLiteral{
		Lbrack: lbrack,
		Rbrack: rbrack,
		Values: values,
	}
}

func (p *Parser) parseArrayLiteralBody() (values []ast.Expr, lbrack, rbrack token.Pos) {
	lbrack = p.expect("[").Pos
	if p.Token.Kind != "]" {
		for p.Token.Kind != token.TokenEOF {
			values = append(values, p.parseExpr())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}
	rbrack = p.expect("]").Pos
	return
}

func (p *Parser) parseStructLiteral() *ast.StructLiteral {
	pos := p.expect("STRUCT").Pos
	fields, _ := p.parseStructTypeFields()
	lparen := p.expect("(").Pos
	var values []ast.Expr
	if p.Token.Kind != ")" {
		for p.Token.Kind != token.TokenEOF {
			values = append(values, p.parseExpr())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}
	rparen := p.expect(")").Pos
	return &ast.StructLiteral{
		Struct: pos,
		Lparen: lparen,
		Rparen: rparen,
		Fields: fields,
		Values: values,
	}
}

func (p *Parser) parseDateLiteral(id token.Token) *ast.DateLiteral {
	s := p.parseStringLiteral()
	return &ast.DateLiteral{
		Date:  id.Pos,
		Value: s,
	}
}

func (p *Parser) parseTimestampLiteral(id token.Token) *ast.TimestampLiteral {
	s := p.parseStringLiteral()
	return &ast.TimestampLiteral{
		Timestamp: id.Pos,
		Value:     s,
	}
}

func (p *Parser) parseNumericLiteral(id token.Token) *ast.NumericLiteral {
	s := p.parseStringLiteral()
	return &ast.NumericLiteral{
		Numeric: id.Pos,
		Value:   s,
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

	p.nextToken()
	// (SELECT ... indicates subquery.
	if p.Token.Kind == "SELECT" {
		return true
	}

	// ((...(SELECT maybe indicate subquery.
	nest := 0
	for p.Token.Kind == "(" {
		nest++
		p.nextToken()
	}
	if nest == 0 || p.Token.Kind != "SELECT" {
		return false
	}

	// ((...(SELECT ...)...) UNION indicates subquery.
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == "(" {
			nest++
		}
		if p.Token.Kind == ")" {
			nest--
		}

		if nest == 0 {
			break
		}
		p.nextToken()
	}
	if nest != 0 {
		return false
	}
	p.nextToken()
	switch p.Token.Kind {
	case "UNION", "INTERSECT", "EXCEPT", "ORDER", "LIMIT":
		return true
	}
	return false
}

// ParseType parses a type name.
func (p *Parser) ParseType() (typ ast.Type, err error) {
	defer func() {
		if r := recover(); r != nil {
			typ = nil
			if e, ok := r.(*Error); ok {
				err = e
			} else {
				panic(r)
			}
		}
	}()

	p.nextToken()
	typ = p.parseType()
	if p.Token.Kind != token.TokenEOF {
		p.panicfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind)
	}
	return
}
func (p *Parser) parseType() ast.Type {
	switch p.Token.Kind {
	case token.TokenIdent:
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
	"NUMERIC",
	"STRING",
	"BYTES",
}

func (p *Parser) parseSimpleType() *ast.SimpleType {
	id := p.expect(token.TokenIdent)
	for _, typeName := range simpleTypes {
		if id.IsIdent(typeName) {
			return &ast.SimpleType{
				NamePos: id.Pos,
				Name:    ast.ScalarTypeName(typeName),
			}
		}
	}

	panic(p.errorfAtToken(id, "expected identifier: %s, but: %s", strings.Join(simpleTypes, ", "), id.Raw))
}

func (p *Parser) parseArrayType() *ast.ArrayType {
	pos := p.expect("ARRAY").Pos
	p.expect("<")
	t := p.parseType()

	var gt token.Pos
	if p.Token.Kind == ">>" {
		p.Token.Kind = ">"
		p.Token.Raw = ">"
		gt = p.Token.Pos
		p.Token.Pos += 1
	} else {
		gt = p.expect(">").Pos
	}
	return &ast.ArrayType{
		Array: pos,
		Gt:    gt,
		Item:  t,
	}
}

func (p *Parser) parseStructType() *ast.StructType {
	pos := p.expect("STRUCT").Pos
	fields, gt := p.parseStructTypeFields()
	if fields == nil {
		p.panicfAtToken(&p.Token, "expected token: <, <>, but: %s", p.Token.Kind)
	}
	return &ast.StructType{
		Struct: pos,
		Gt:     gt,
		Fields: fields,
	}
}

func (p *Parser) parseStructTypeFields() (fields []*ast.StructField, gt token.Pos) {
	if p.Token.Kind != "<" && p.Token.Kind != "<>" {
		return
	}

	fields = make([]*ast.StructField, 0)
	if p.Token.Kind == "<>" {
		gt = p.expect("<>").Pos + 1
		return
	}

	p.expect("<")
	if p.Token.Kind != ">" && p.Token.Kind != ">>" {
		for p.Token.Kind != token.TokenEOF {
			fields = append(fields, p.parseFieldType())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}

	if p.Token.Kind == ">>" {
		p.Token.Kind = ">"
		p.Token.Raw = ">"
		gt = p.Token.Pos
		p.Token.Pos += 1
	} else {
		gt = p.expect(">").Pos
	}
	return
}

func (p *Parser) parseFieldType() *ast.StructField {
	lexer := p.Lexer.Clone()
	// Try to parse as "x INT64" case.
	if p.Token.Kind == token.TokenIdent {
		ident := p.parseIdent()
		if p.lookaheadType() {
			t := p.parseType()
			return &ast.StructField{
				Ident: ident,
				Type:  t,
			}
		}
	}

	p.Lexer = lexer
	return &ast.StructField{
		Type: p.parseType(),
	}
}

func (p *Parser) lookaheadType() bool {
	return p.Token.Kind == token.TokenIdent || p.Token.Kind == "ARRAY" || p.Token.Kind == "STRUCT"
}

// ================================================================================
//
// DDL
//
// ================================================================================

func (p *Parser) parseDDL() ast.DDL {
	pos := p.Token.Pos
	switch {
	case p.Token.Kind == "CREATE":
		p.nextToken()
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
		p.nextToken()
		return p.parseAlterTable(pos)
	case p.Token.IsKeywordLike("DROP"):
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseDropTable(pos)
		case p.Token.IsKeywordLike("INDEX"):
			return p.parseDropIndex(pos)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: TABLE, INDEX, but: %s", p.Token.AsString)
	}

	if p.Token.Kind != token.TokenIdent {
		panic(p.errorfAtToken(&p.Token, "expected token: CREATE, <ident>, but: %s", p.Token.Kind))
	}

	panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: ALTER, DROP, but: %s", p.Token.AsString))
}

func (p *Parser) parseCreateDatabase(pos token.Pos) *ast.CreateDatabase {
	p.expectKeywordLike("DATABASE")
	name := p.parseIdent()
	return &ast.CreateDatabase{
		Create: pos,
		Name:   name,
	}
}

func (p *Parser) parseCreateTable(pos token.Pos) *ast.CreateTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()

	// This loop allows parsing trailing comma intentionally.
	// TODO: is this allowed by Spanner really?
	p.expect("(")
	var columns []*ast.ColumnDef
	var foreignKeys []*ast.ForeignKey
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		switch {
		case p.Token.IsKeywordLike("CONSTRAINT"):
			foreignKeys = append(foreignKeys, p.parseConstraint())
		case p.Token.IsKeywordLike("FOREIGN"):
			foreignKeys = append(foreignKeys, p.parseForeignKey())
		default:
			columns = append(columns, p.parseColumnDef())
		}
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}
	p.expect(")")

	p.expectKeywordLike("PRIMARY")
	p.expectKeywordLike("KEY")

	p.expect("(")
	var keys []*ast.IndexKey
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		keys = append(keys, p.parseIndexKey())
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}
	rparen := p.expect(")").Pos

	cluster := p.tryParseCluster()

	return &ast.CreateTable{
		Create:      pos,
		Rparen:      rparen,
		Name:        name,
		Columns:     columns,
		ForeignKeys: foreignKeys,
		PrimaryKeys: keys,
		Cluster:     cluster,
	}
}

func (p *Parser) parseColumnDef() *ast.ColumnDef {
	name := p.parseIdent()
	t, notNull, null := p.parseTypeNotNull()
	generated := p.tryParseGeneratedColumnExpr()
	options := p.tryParseColumnDefOptions()

	return &ast.ColumnDef{
		Null:          null,
		Name:          name,
		Type:          t,
		NotNull:       notNull,
		GeneratedExpr: generated,
		Options:       options,
	}
}

func (p *Parser) parseConstraint() *ast.ForeignKey {
	pos := p.expectKeywordLike("CONSTRAINT").Pos
	name := p.parseIdent()
	fk := p.parseForeignKey()
	fk.Constraint = pos
	fk.Name = name
	return fk
}

func (p *Parser) parseForeignKey() *ast.ForeignKey {
	pos := p.expectKeywordLike("FOREIGN").Pos
	p.expectKeywordLike("KEY")

	p.expect("(")
	columns := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.nextToken()
		columns = append(columns, p.parseIdent())
	}
	p.expect(")")
	p.expectKeywordLike("REFERENCES")
	refTable := p.parseIdent()

	p.expect("(")
	refColumns := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.nextToken()
		refColumns = append(refColumns, p.parseIdent())
	}
	rparen := p.expect(")").End

	return &ast.ForeignKey{
		Foreign:          pos,
		Rparen:           rparen,
		Name:             nil,
		Columns:          columns,
		ReferenceTable:   refTable,
		ReferenceColumns: refColumns,
	}
}

func (p *Parser) parseTypeNotNull() (t ast.SchemaType, notNull bool, null token.Pos) {
	t = p.parseSchemaType()

	null = token.InvalidPos
	if p.Token.Kind == "NOT" {
		p.expect("NOT")
		null = p.expect("NULL").Pos
		notNull = true
	}
	return
}

func (p *Parser) tryParseGeneratedColumnExpr() *ast.GeneratedColumnExpr {
	if p.Token.Kind != "AS" {
		return nil
	}

	posAs := p.expect("AS").Pos
	p.expect("(")
	expr := p.parseExpr()
	p.expect(")")
	posEnd := p.expectKeywordLike("STORED").End

	return &ast.GeneratedColumnExpr{
		As:     posAs,
		Stored: posEnd,
		Expr:   expr,
	}
}

func (p *Parser) tryParseColumnDefOptions() *ast.ColumnDefOptions {
	if !p.Token.IsKeywordLike("OPTIONS") {
		return nil
	}

	return p.parseColumnDefOptions()
}

func (p *Parser) parseColumnDefOptions() *ast.ColumnDefOptions {
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
	p.nextToken()

	rparen := p.expect(")").End
	return &ast.ColumnDefOptions{
		Options:              pos,
		Rparen:               rparen,
		AllowCommitTimestamp: allowCommitTimestamp,
	}
}

func (p *Parser) parseIndexKey() *ast.IndexKey {
	name := p.parseIdent()
	dir, dirPos := p.tryParseDirection()

	return &ast.IndexKey{
		DirPos: dirPos,
		Name:   name,
		Dir:    dir,
	}
}

func (p *Parser) tryParseCluster() *ast.Cluster {
	if p.Token.Kind != "," {
		return nil
	}
	pos := p.expect(",").Pos

	p.expectKeywordLike("INTERLEAVE")
	p.expect("IN")
	p.expectKeywordLike("PARENT")
	name := p.parseIdent()

	onDelete, onDeleteEnd := p.tryParseOnDeleteAction()

	return &ast.Cluster{
		Comma:       pos,
		OnDeleteEnd: onDeleteEnd,
		TableName:   name,
		OnDelete:    onDelete,
	}
}

func (p *Parser) tryParseOnDeleteAction() (onDelete ast.OnDeleteAction, onDeleteEnd token.Pos) {
	onDeleteEnd = token.InvalidPos
	if p.Token.Kind != "ON" {
		return
	}

	onDelete, onDeleteEnd = p.parseOnDeleteAction()
	return
}

func (p *Parser) parseOnDeleteAction() (onDelete ast.OnDeleteAction, onDeleteEnd token.Pos) {
	p.expect("ON")
	p.expectKeywordLike("DELETE")
	switch p.Token.Kind {
	case token.TokenIdent:
		onDeleteEnd = p.expectKeywordLike("CASCADE").End
		onDelete = ast.OnDeleteCascade
	case "NO":
		p.nextToken()
		onDeleteEnd = p.expectKeywordLike("ACTION").End
		onDelete = ast.OnDeleteNoAction
	default:
		p.panicfAtToken(&p.Token, "expected token: NO, <ident>, but: %s", p.Token.Kind)
	}
	return
}

func (p *Parser) parseCreateIndex(pos token.Pos) *ast.CreateIndex {
	unique := false
	if p.Token.IsKeywordLike("UNIQUE") {
		p.nextToken()
		unique = true
	}

	nullFiltered := false
	if p.Token.IsKeywordLike("NULL_FILTERED") {
		p.nextToken()
		nullFiltered = true
	}

	p.expectKeywordLike("INDEX")
	name := p.parseIdent()

	p.expect("ON")
	tableName := p.parseIdent()

	p.expect("(")
	var keys []*ast.IndexKey
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		keys = append(keys, p.parseIndexKey())
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}
	rparen := p.expect(")").Pos

	storing := p.tryParseStoring()
	interleaveIn := p.tryParseInterleaveIn()

	return &ast.CreateIndex{
		Create:       pos,
		Rparen:       rparen,
		Unique:       unique,
		NullFiltered: nullFiltered,
		Name:         name,
		TableName:    tableName,
		Keys:         keys,
		Storing:      storing,
		InterleaveIn: interleaveIn,
	}
}

func (p *Parser) tryParseStoring() *ast.Storing {
	if !p.Token.IsKeywordLike("STORING") {
		return nil
	}
	pos := p.expectKeywordLike("STORING").Pos

	p.expect("(")
	columns := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "," {
		p.nextToken()
		columns = append(columns, p.parseIdent())
	}

	rparen := p.expect(")").Pos
	return &ast.Storing{
		Storing: pos,
		Rparen:  rparen,
		Columns: columns,
	}
}

func (p *Parser) tryParseInterleaveIn() *ast.InterleaveIn {
	if p.Token.Kind != "," {
		return nil
	}
	pos := p.expect(",").Pos
	p.expectKeywordLike("INTERLEAVE")
	p.expect("IN")
	name := p.parseIdent()

	return &ast.InterleaveIn{
		Comma:     pos,
		TableName: name,
	}
}

func (p *Parser) parseAlterTable(pos token.Pos) *ast.AlterTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()

	var alternation ast.TableAlternation
	switch {
	case p.Token.IsKeywordLike("ADD"):
		alternation = p.parseAlterTableAdd()
	case p.Token.IsKeywordLike("DROP"):
		alternation = p.parseAlterTableDrop()
	case p.Token.Kind == "SET":
		alternation = p.parseSetOnDelete()
	case p.Token.IsKeywordLike("ALTER"):
		alternation = p.parseAlterColumn()
	default:
		if p.Token.Kind == token.TokenIdent {
			p.panicfAtToken(&p.Token, "expected pseuso keyword: ADD, ALTER, DROP, but: %s", p.Token.AsString)
		} else {
			p.panicfAtToken(&p.Token, "expected token: SET, <ident>, but: %s", p.Token.Kind)
		}
	}

	return &ast.AlterTable{
		Alter:            pos,
		Name:             name,
		TableAlternation: alternation,
	}
}

func (p *Parser) parseAlterTableAdd() ast.TableAlternation {
	pos := p.expectKeywordLike("ADD").Pos

	var alternation ast.TableAlternation

	switch {
	case p.Token.IsKeywordLike("COLUMN"):
		p.expectKeywordLike("COLUMN")
		column := p.parseColumnDef()
		alternation = &ast.AddColumn{
			Add:    pos,
			Column: column,
		}
	case p.Token.IsKeywordLike("CONSTRAINT"):
		fk := p.parseConstraint()
		alternation = &ast.AddForeignKey{
			Add:        pos,
			ForeignKey: fk,
		}
	case p.Token.IsKeywordLike("FOREIGN"):
		fk := p.parseForeignKey()
		alternation = &ast.AddForeignKey{
			Add:        pos,
			ForeignKey: fk,
		}
	default:
		p.panicfAtToken(&p.Token, "expected pseuso keyword: COLUMN, CONSTRAINT, FOREIGN, but: %s", p.Token.AsString)
	}

	return alternation
}

func (p *Parser) parseAlterTableDrop() ast.TableAlternation {
	pos := p.expectKeywordLike("DROP").Pos

	var alternation ast.TableAlternation

	switch {
	case p.Token.IsKeywordLike("COLUMN"):
		p.expectKeywordLike("COLUMN")
		name := p.parseIdent()
		alternation = &ast.DropColumn{
			Drop: pos,
			Name: name,
		}
	case p.Token.IsKeywordLike("CONSTRAINT"):
		p.expectKeywordLike("CONSTRAINT")
		name := p.parseIdent()
		alternation = &ast.DropConstraint{
			Drop: pos,
			Name: name,
		}
	default:
		p.panicfAtToken(&p.Token, "expected pseuso keyword: COLUMN, CONSTRAINT, but: %s", p.Token.AsString)
	}

	return alternation
}

func (p *Parser) parseSetOnDelete() *ast.SetOnDelete {
	pos := p.expect("SET").Pos
	onDelete, onDeleteEnd := p.parseOnDeleteAction()

	return &ast.SetOnDelete{
		Set:         pos,
		OnDeleteEnd: onDeleteEnd,
		OnDelete:    onDelete,
	}
}

func (p *Parser) parseAlterColumn() ast.TableAlternation {
	pos := p.expectKeywordLike("ALTER").Pos
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	if p.Token.Kind == "SET" {
		p.nextToken()
		options := p.parseColumnDefOptions()
		return &ast.AlterColumnSet{
			Alter:   pos,
			Name:    name,
			Options: options,
		}
	}

	t, notNull, null := p.parseTypeNotNull()

	return &ast.AlterColumn{
		Alter:   pos,
		Null:    null,
		Name:    name,
		Type:    t,
		NotNull: notNull,
	}
}

func (p *Parser) parseDropTable(pos token.Pos) *ast.DropTable {
	p.expectKeywordLike("TABLE")
	name := p.parseIdent()
	return &ast.DropTable{
		Drop: pos,
		Name: name,
	}
}

func (p *Parser) parseDropIndex(pos token.Pos) *ast.DropIndex {
	p.expectKeywordLike("INDEX")
	name := p.parseIdent()
	return &ast.DropIndex{
		Drop: pos,
		Name: name,
	}
}

func (p *Parser) parseSchemaType() ast.SchemaType {
	switch p.Token.Kind {
	case token.TokenIdent:
		return p.parseScalarSchemaType()
	case "ARRAY":
		pos := p.expect("ARRAY").Pos
		p.expect("<")
		t := p.parseScalarSchemaType()
		end := p.expect(">").End
		return &ast.ArraySchemaType{
			Array: pos,
			Gt:    end,
			Item:  t,
		}
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ARRAY, <ident>, but: %s", p.Token.Kind))
}

var scalarSchemaTypes = []string{
	"BOOL",
	"INT64",
	"FLOAT64",
	"DATE",
	"TIMESTAMP",
	"NUMERIC",
}

var sizedSchemaTypes = []string{
	"STRING",
	"BYTES",
}

func (p *Parser) parseScalarSchemaType() ast.SchemaType {
	id := p.expect(token.TokenIdent)
	pos := id.Pos

	for _, name := range scalarSchemaTypes {
		if id.IsIdent(name) {
			return &ast.ScalarSchemaType{
				NamePos: pos,
				Name:    ast.ScalarTypeName(name),
			}
		}
	}

	for _, name := range sizedSchemaTypes {
		if id.IsIdent(name) {
			p.expect("(")
			max := false
			var size ast.IntValue
			if p.Token.IsIdent("MAX") {
				p.nextToken()
				max = true
			} else {
				size = p.parseIntValue()
			}
			rparen := p.expect(")").Pos
			return &ast.SizedSchemaType{
				NamePos: pos,
				Rparen:  rparen,
				Name:    ast.ScalarTypeName(name),
				Max:     max,
				Size:    size,
			}
		}
	}

	panic(p.errorfAtToken(id, "expect ident: %s, %s, but: %s", strings.Join(scalarSchemaTypes, ","), strings.Join(sizedSchemaTypes, ","), id.AsString))
}

// ================================================================================
//
// DML
//
// ================================================================================

func (p *Parser) parseDML() ast.DML {
	id := p.expect(token.TokenIdent)
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

func (p *Parser) parseInsert(pos token.Pos) *ast.Insert {
	if p.Token.Kind == "INTO" {
		p.nextToken()
	}

	name := p.parseIdent()

	p.expect("(")
	var columns []*ast.Ident
	if p.Token.Kind != ")" {
		for p.Token.Kind != token.TokenEOF {
			columns = append(columns, p.parseIdent())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}
	p.expect(")")

	var input ast.InsertInput
	if p.Token.IsKeywordLike("VALUES") {
		input = p.parseValuesInput()
	} else {
		input = p.parseSubQueryInput()
	}

	return &ast.Insert{
		Insert:    pos,
		TableName: name,
		Columns:   columns,
		Input:     input,
	}
}

func (p *Parser) parseValuesInput() *ast.ValuesInput {
	pos := p.expectKeywordLike("VALUES").Pos

	rows := []*ast.ValuesRow{p.parseValuesRow()}
	for p.Token.Kind == "," {
		p.nextToken()
		rows = append(rows, p.parseValuesRow())
	}

	return &ast.ValuesInput{
		Values: pos,
		Rows:   rows,
	}
}

func (p *Parser) parseValuesRow() *ast.ValuesRow {
	lparen := p.expect("(").Pos
	var exprs []*ast.DefaultExpr
	if p.Token.Kind != ")" {
		for p.Token.Kind != token.TokenEOF {
			exprs = append(exprs, p.parseDefaultExpr())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}
	rparen := p.expect(")").Pos

	return &ast.ValuesRow{
		Lparen: lparen,
		Rparen: rparen,
		Exprs:  exprs,
	}
}

func (p *Parser) parseDefaultExpr() *ast.DefaultExpr {
	if p.Token.Kind == "DEFAULT" {
		pos := p.expect("DEFAULT").Pos
		return &ast.DefaultExpr{
			DefaultPos: pos,
			Default:    true,
		}
	}

	expr := p.parseExpr()
	return &ast.DefaultExpr{
		DefaultPos: token.InvalidPos,
		Expr:       expr,
	}
}

func (p *Parser) parseSubQueryInput() *ast.SubQueryInput {
	query := p.parseQueryExpr()

	return &ast.SubQueryInput{
		Query: query,
	}
}

func (p *Parser) parseDelete(pos token.Pos) *ast.Delete {
	if p.Token.Kind == "FROM" {
		p.nextToken()
	}

	name := p.parseIdent()
	as := p.tryParseAsAlias()
	where := p.parseWhere()

	return &ast.Delete{
		Delete:    pos,
		TableName: name,
		As:        as,
		Where:     where,
	}
}

func (p *Parser) parseUpdate(pos token.Pos) *ast.Update {
	name := p.parseIdent()
	as := p.tryParseAsAlias()

	p.expect("SET")

	items := []*ast.UpdateItem{p.parseUpdateItem()}
	for p.Token.Kind == "," {
		p.nextToken()
		items = append(items, p.parseUpdateItem())
	}

	where := p.parseWhere()

	return &ast.Update{
		Update:    pos,
		TableName: name,
		As:        as,
		Updates:   items,
		Where:     where,
	}
}

func (p *Parser) parseUpdateItem() *ast.UpdateItem {
	path := p.parseIdentOrPath()
	p.expect("=")
	expr := p.parseExpr()

	return &ast.UpdateItem{
		Path: path,
		Expr: expr,
	}
}

// ================================================================================
//
// Primitives
//
// ================================================================================

func (p *Parser) parseIdent() *ast.Ident {
	id := p.expect(token.TokenIdent)
	return &ast.Ident{
		NamePos: id.Pos,
		NameEnd: id.End,
		Name:    id.AsString,
	}
}

func (p *Parser) parseParam() *ast.Param {
	param := p.expect(token.TokenParam)
	return &ast.Param{
		Atmark: param.Pos,
		Name:   param.AsString,
	}
}

func (p *Parser) parseNullLiteral() *ast.NullLiteral {
	tok := p.expect("NULL")
	return &ast.NullLiteral{
		Null: tok.Pos,
	}
}

func (p *Parser) parseBoolLiteral() *ast.BoolLiteral {
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
	p.nextToken()
	return &ast.BoolLiteral{
		ValuePos: pos,
		Value:    value,
	}
}

func (p *Parser) parseIntLiteral() *ast.IntLiteral {
	i := p.expect(token.TokenInt)
	return &ast.IntLiteral{
		ValuePos: i.Pos,
		ValueEnd: i.End,
		Base:     i.Base,
		Value:    i.Raw,
	}
}

func (p *Parser) parseFloatLiteral() *ast.FloatLiteral {
	f := p.expect(token.TokenFloat)
	return &ast.FloatLiteral{
		ValuePos: f.Pos,
		ValueEnd: f.End,
		Value:    f.Raw,
	}
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	s := p.expect(token.TokenString)
	return &ast.StringLiteral{
		ValuePos: s.Pos,
		ValueEnd: s.End,
		Value:    s.AsString,
	}
}

func (p *Parser) parseBytesLiteral() *ast.BytesLiteral {
	b := p.expect(token.TokenBytes)
	return &ast.BytesLiteral{
		ValuePos: b.Pos,
		ValueEnd: b.End,
		Value:    []byte(b.AsString),
	}
}

func (p *Parser) parseIntValue() ast.IntValue {
	switch p.Token.Kind {
	case token.TokenParam:
		return p.parseParam()
	case token.TokenInt:
		return p.parseIntLiteral()
	case "CAST":
		return p.parseCastIntValue()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <int>, CAST, but: %s", p.Token.Kind))
}

func (p *Parser) parseCastIntValue() *ast.CastIntValue {
	pos := p.expect("CAST").Pos
	p.expect("(")
	var v ast.IntValue
	switch p.Token.Kind {
	case token.TokenParam:
		v = p.parseParam()
	case token.TokenInt:
		v = p.parseIntLiteral()
	default:
		p.panicfAtToken(&p.Token, "expected token: <param>, <int>, but: %s", p.Token.Kind)
	}
	p.expect("AS")
	p.expectIdent("INT64")
	rparen := p.expect(")").Pos
	return &ast.CastIntValue{
		Cast:   pos,
		Rparen: rparen,
		Expr:   v,
	}
}

func (p *Parser) parseNumValue() ast.NumValue {
	switch p.Token.Kind {
	case token.TokenParam:
		return p.parseParam()
	case token.TokenInt:
		return p.parseIntLiteral()
	case token.TokenFloat:
		return p.parseFloatLiteral()
	case "CAST":
		return p.parseCastNumValue()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <int>, <float>, CAST, but: %s", p.Token.Kind))
}

func (p *Parser) parseCastNumValue() *ast.CastNumValue {
	pos := p.expect("CAST").Pos
	p.expect("(")
	var v ast.NumValue
	switch p.Token.Kind {
	case token.TokenParam:
		v = p.parseParam()
	case token.TokenInt:
		v = p.parseIntLiteral()
	case token.TokenFloat:
		v = p.parseFloatLiteral()
	default:
		p.panicfAtToken(&p.Token, "expected token: <param>, <int>, <float>, but: %s", p.Token.Kind)
	}
	p.expect("AS")
	id := p.expect(token.TokenIdent)
	var t ast.ScalarTypeName
	switch {
	case id.IsIdent("INT64"):
		t = ast.Int64TypeName
	case id.IsIdent("FLOAT64"):
		t = ast.Float64TypeName
	default:
		p.panicfAtToken(id, "expected identifier: INT64, FLOAT64, but: %s", id.Raw)
	}
	rparen := p.expect(")").Pos
	return &ast.CastNumValue{
		Cast:   pos,
		Rparen: rparen,
		Expr:   v,
		Type:   t,
	}
}

func (p *Parser) parseStringValue() ast.StringValue {
	switch p.Token.Kind {
	case token.TokenParam:
		return p.parseParam()
	case token.TokenString:
		return p.parseStringLiteral()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: <param>, <string>, but: %s", p.Token.Kind))
}

func (p *Parser) expect(kind token.TokenKind) *token.Token {
	if p.Token.Kind != kind {
		p.panicfAtToken(&p.Token, "expected token: %s, but: %s", kind, p.Token.Kind)
	}
	t := p.Token.Clone()
	p.nextToken()
	return t
}

func (p *Parser) expectIdent(s string) *token.Token {
	id := p.expect(token.TokenIdent)
	if !id.IsIdent(s) {
		p.panicfAtToken(id, "expected identifier: %s, but: %s", s, token.QuoteSQLIdent(id.AsString))
	}
	return id
}

func (p *Parser) expectKeywordLike(s string) *token.Token {
	id := p.expect(token.TokenIdent)
	if !id.IsKeywordLike(s) {
		if char.EqualFold(id.AsString, s) {
			p.panicfAtToken(id, "pseudo keyword %s cannot encloses with backquote", s)
		} else {
			p.panicfAtToken(id, "expected pseudo keyword: %s, but: %s", s, token.QuoteSQLIdent(id.AsString))
		}
	}
	return id
}

func (p *Parser) errorfAtToken(tok *token.Token, msg string, params ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(msg, params...),
		Position: p.Position(tok.Pos, tok.End),
	}
}

func (p *Parser) panicfAtToken(tok *token.Token, msg string, params ...interface{}) {
	panic(p.errorfAtToken(tok, msg, params...))
}
