package memefish

import (
	"fmt"
	"strings"

	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/char"
	"github.com/cloudspannerecosystem/memefish/token"
)

type Parser struct {
	*Lexer

	errors []*Error
}

// ParseStatement parses a SQL statement.
func (p *Parser) ParseStatement() (ast.Statement, error) {
	p.nextToken()
	stmt := p.parseStatement()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return stmt, err
	}

	return stmt, nil
}

// ParseStatements parses SQL statements list separated by semi-colon.
func (p *Parser) ParseStatements() ([]ast.Statement, error) {
	p.nextToken()
	stmts := parseStatements(p, p.parseStatement)
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return stmts, err
	}

	return stmts, nil
}

// ParseQuery parses a query statement.
func (p *Parser) ParseQuery() (*ast.QueryStatement, error) {
	p.nextToken()
	stmt := p.parseQueryStatement()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return stmt, err
	}

	return stmt, nil
}

// ParseExpr parses a SQL expression.
func (p *Parser) ParseExpr() (ast.Expr, error) {
	p.nextToken()
	expr := p.parseExpr()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return expr, err
	}

	return expr, nil
}

// ParseType parses a type name.
func (p *Parser) ParseType() (ast.Type, error) {
	p.nextToken()
	t := p.parseType()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return t, err
	}

	return t, nil
}

// ParseDDL parses a CREATE/ALTER/DROP statement.
func (p *Parser) ParseDDL() (ast.DDL, error) {
	p.nextToken()
	ddl := p.parseDDL()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return ddl, err
	}

	return ddl, nil
}

// ParseDDLs parses CREATE/ALTER/DROP statements list separated by semi-colon.
func (p *Parser) ParseDDLs() ([]ast.DDL, error) {
	p.nextToken()
	ddls := parseStatements(p, p.parseDDL)
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return ddls, err
	}

	return ddls, nil
}

// ParseDML parses a INSERT/DELETE/UPDATE statement.
func (p *Parser) ParseDML() (ast.DML, error) {
	p.nextToken()
	dml := p.parseDML()
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return dml, err
	}

	return dml, nil
}

// ParseDMLs parses INSERT/DELETE/UPDATE statements list separated by semi-colon.
func (p *Parser) ParseDMLs() ([]ast.DML, error) {
	p.nextToken()
	dmls := parseStatements(p, p.parseDML)
	if p.Token.Kind != token.TokenEOF {
		p.errors = append(p.errors, p.errorfAtToken(&p.Token, "expected token: <eof>, but: %s", p.Token.Kind))
	}

	if len(p.errors) > 0 {
		// Reset the errors and allow processing to continue
		err := MultiError(p.errors)
		p.errors = nil

		return dmls, err
	}

	return dmls, nil
}

func (p *Parser) parseStatement() (stmt ast.Statement) {
	l := p.Clone()
	defer func() {
		// Panic on tryParseHint()
		if r := recover(); r != nil {
			stmt = &ast.BadStatement{BadNode: p.handleParseStatementError(r, l)}
		}
	}()

	hint := p.tryParseHint()
	return p.parseStatementInternal(hint)
}

func (p *Parser) parseStatementInternal(hint *ast.Hint) (stmt ast.Statement) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			stmt = &ast.BadStatement{
				Hint:    hint,
				BadNode: p.handleParseStatementError(r, l),
			}
		}
	}()

	switch {
	case p.Token.Kind == "SELECT" || p.Token.Kind == "WITH" || p.Token.Kind == "(" || p.Token.Kind == "FROM":
		return p.parseQueryStatementInternal(hint)
	case p.Token.IsKeywordLike("INSERT") || p.Token.IsKeywordLike("DELETE") || p.Token.IsKeywordLike("UPDATE"):
		return p.parseDMLInternal(hint)
	case hint != nil:
		panic(p.errorfAtPosition(hint.Pos(), p.Token.End, "statement hint is only permitted before query or DML, but got: %s", p.Token.Raw))
	case p.Token.Kind == "CREATE" || p.Token.IsKeywordLike("ALTER") || p.Token.IsKeywordLike("DROP") ||
		p.Token.IsKeywordLike("RENAME") || p.Token.IsKeywordLike("GRANT") || p.Token.IsKeywordLike("REVOKE") ||
		p.Token.IsKeywordLike("ANALYZE"):
		return p.parseDDL()
	case p.Token.IsKeywordLike("CALL"):
		return p.parseOtherStatement()
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseOtherStatement() ast.Statement {
	switch {
	case p.Token.IsKeywordLike("CALL"):
		return p.parseCall()
	}

	panic(p.errorfAtToken(&p.Token, "unexpected token: %s", p.Token.Kind))
}

func (p *Parser) parseCall() *ast.Call {
	pos := p.expectKeywordLike("CALL").Pos
	name := p.parsePath()
	p.expect("(")

	var args []ast.TVFArg
	if p.Token.Kind != ")" {
		args = parseCommaSeparatedList(p, p.parseTVFArg)
	}

	rparen := p.expect(")").Pos

	return &ast.Call{
		Call:   pos,
		Rparen: rparen,
		Name:   name,
		Args:   args,
	}
}

func parseStatements[T ast.Node](p *Parser, doParse func() T) []T {
	var nodes []T
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ";" {
			p.nextToken()
			continue
		}

		nodes = append(nodes, doParse())

		if p.Token.Kind != ";" {
			break
		}
	}
	return nodes
}

// ================================================================================
//
// SELECT
//
// ================================================================================

func (p *Parser) parseQueryStatement() (stmt *ast.QueryStatement) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			// When parsing is failed on tryParseHint, the result of these methods are discarded
			// because they are concrete structs and we cannot fill them with *ast.BadNode.
			stmt = &ast.QueryStatement{
				Query: &ast.BadQueryExpr{BadNode: p.handleParseStatementError(r, l)},
			}
		}
	}()

	hint := p.tryParseHint()
	return p.parseQueryStatementInternal(hint)
}

func (p *Parser) parseQueryStatementInternal(hint *ast.Hint) (stmt *ast.QueryStatement) {
	// Can be a *ast.BadQueryExpr and won't panic
	query := p.parseQueryExpr()

	return &ast.QueryStatement{
		Hint:  hint,
		Query: query,
	}
}

func (p *Parser) parsePipeOperator() ast.PipeOperator {
	pos := p.expect("|>").Pos
	switch p.Token.Kind {
	case "SELECT":
		p.nextToken()

		allOrDistinct := p.tryParseAllOrDistinct()
		as := p.tryParseSelectAs()
		results := p.parseSelectResults()

		return &ast.PipeSelect{
			Pipe:          pos,
			AllOrDistinct: allOrDistinct,
			As:            as,
			Results:       results,
		}
	case "WHERE":
		p.nextToken()
		expr := p.parseExpr()
		return &ast.PipeWhere{
			Pipe: pos,
			Expr: expr,
		}
	default:
		panic(p.errorfAtToken(&p.Token, "expected pipe operator name, but: %q", p.Token.AsString))
	}
}

// parsePipeOperators parses pipe operators, which can be empty.
func (p *Parser) parsePipeOperators() []ast.PipeOperator {
	var pipeOps []ast.PipeOperator
	for p.Token.Kind == "|>" {
		pipeOps = append(pipeOps, p.parsePipeOperator())
	}
	return pipeOps
}

// parseQuery consumes ORDER BY, LIMIT, and pipe operators.
func (p *Parser) parseQuery() *ast.Query {
	with := p.tryParseWith()

	if p.Token.Kind == "WITH" {
		panic(p.errorfAtToken(&p.Token, "expect query expression, but unexpected WITH"))
	}

	query := p.parseQueryExpr()

	// If nested query expression is *ast.Query(with suffix), merge parsed WITH into it to avoid deep nest.
	if q, ok := query.(*ast.Query); ok {
		return &ast.Query{
			With:          with,
			Query:         q.Query,
			OrderBy:       q.OrderBy,
			Limit:         q.Limit,
			ForUpdate:     q.ForUpdate,
			PipeOperators: q.PipeOperators,
		}
	}

	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()
	forUpdate := p.tryParseForUpdate()
	pipeOps := p.parsePipeOperators()

	return &ast.Query{
		With:          with,
		Query:         query,
		OrderBy:       orderBy,
		Limit:         limit,
		ForUpdate:     forUpdate,
		PipeOperators: pipeOps,
	}
}

func (p *Parser) tryParseForUpdate() *ast.ForUpdate {
	if p.Token.Kind != "FOR" {
		return nil
	}

	forPos := p.expect("FOR").Pos
	update := p.expectKeywordLike("UPDATE").Pos

	return &ast.ForUpdate{
		For:    forPos,
		Update: update,
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
	key := p.parsePath()
	p.expect("=")
	value := p.parseExpr()
	return &ast.HintRecord{
		Key:   key,
		Value: value,
	}
}

func (p *Parser) tryParseWith() *ast.With {
	if p.Token.Kind != "WITH" {
		return nil
	}
	pos := p.Token.Pos
	p.nextToken()
	ctes := parseCommaSeparatedList(p, p.parseCTE)

	return &ast.With{
		With: pos,
		CTEs: ctes,
	}
}

func (p *Parser) parseCTE() *ast.CTE {
	name := p.parseIdent()
	p.expect("AS")
	p.expect("(")
	query := p.parseQueryExpr()
	rparen := p.expect(")").Pos
	return &ast.CTE{
		Name:      name,
		QueryExpr: query,
		Rparen:    rparen,
	}
}

func (p *Parser) parseQueryExpr() (query ast.QueryExpr) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			query = p.handleParseQueryExprError(false, r, l)
		}
	}()

	// If WITH is appeared, it is treated as an outer node than compound query.
	if p.Token.Kind == "WITH" {
		return p.parseQuery()
	}

	query = p.parseSimpleQueryExpr()

	// If the query is directly followed by ORDER BY, LIMIT or pipe operators, it won't be a compound query
	switch p.Token.Kind {
	case "ORDER", "LIMIT", "FOR", "|>":
		return p.parseQueryExprSuffix(query)
	}

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

		allOrDistinct := p.parseAllOrDistinct()

		right := p.parseSimpleQueryExpr()
		if c, ok := query.(*ast.CompoundQuery); ok {
			if c.Op != op || c.AllOrDistinct != allOrDistinct {
				p.panicfAtToken(&opTok, "all set operator at the same level must be the same, or wrap (...)")
			}
			c.Queries = append(c.Queries, right)
		} else {
			query = &ast.CompoundQuery{
				Op:            op,
				AllOrDistinct: allOrDistinct,
				Queries:       []ast.QueryExpr{query, right},
			}
		}
	}

	// LIMIT, ORDER BY, pipe operators can be placed after a compound query.
	return p.parseQueryExprSuffix(query)
}

func (p *Parser) parseFromQuery() *ast.FromQuery {
	if p.Token.Kind != "FROM" {
		panic(p.errorfAtToken(&p.Token, "expected 'FROM', got %s", p.Token.AsString))
	}
	from := p.tryParseFrom()

	// Although it can be parsed, it is better to reject invalid GoogleSQL queries.
	switch p.Token.Kind {
	case "ORDER":
		panic(p.errorfAtToken(&p.Token, "syntax error: ORDER BY not supported after FROM query; Consider using pipe operator `|> ORDER BY` or parentheses around the FROM query."))
	case "LIMIT":
		panic(p.errorfAtToken(&p.Token, "syntax error: LIMIT not supported after FROM query; Consider using pipe operator `|> LIMIT` or parentheses around the FROM query."))
	}

	return &ast.FromQuery{From: from}
}

// parseSimpleQueryExpr parses simple QueryExpr, which can be wrapped in Query or CompoundQuery.
func (p *Parser) parseSimpleQueryExpr() (query ast.QueryExpr) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			query = p.handleParseQueryExprError(true, r, l)
		}
	}()

	switch p.Token.Kind {
	// FROM and SELECT are the most primitive query form
	case "FROM":
		return p.parseFromQuery()
	case "SELECT":
		return p.parseSelect()
	case "(": // Query with paren
		lparen := p.expect("(").Pos
		q := p.parseQueryExpr()
		rparen := p.expect(")").Pos
		return &ast.SubQuery{
			Lparen: lparen,
			Rparen: rparen,
			Query:  q,
		}
	default:
		panic(p.errorfAtToken(&p.Token, `expected beginning of simple query "(", SELECT, FROM, but: %q`, p.Token.AsString))
	}
}

func (p *Parser) parseSelect() *ast.Select {
	sel := p.expect("SELECT").Pos
	allOrDistinct := p.tryParseAllOrDistinct()
	selectAs := p.tryParseSelectAs()
	results := p.parseSelectResults()
	from := p.tryParseFrom()
	where := p.tryParseWhere()
	groupBy := p.tryParseGroupBy()
	having := p.tryParseHaving()

	return &ast.Select{
		Select:        sel,
		AllOrDistinct: allOrDistinct,
		As:            selectAs,
		Results:       results,
		From:          from,
		Where:         where,
		GroupBy:       groupBy,
		Having:        having,
	}
}

func (p *Parser) tryParseSelectAs() ast.SelectAs {
	if p.Token.Kind != "AS" {
		return nil
	}
	pos := p.expect("AS").Pos

	switch {
	case p.Token.Kind == "STRUCT":
		structPos := p.expect("STRUCT").Pos
		return &ast.AsStruct{
			As:     pos,
			Struct: structPos,
		}
	case p.Token.IsKeywordLike("VALUE"):
		valuePos := p.expectKeywordLike("VALUE").Pos
		return &ast.AsValue{
			As:    pos,
			Value: valuePos,
		}
	default:
		namedType := p.parseNamedType()
		return &ast.AsTypeName{
			As:       pos,
			TypeName: namedType,
		}
	}
}

func (p *Parser) parseSelectResults() []ast.SelectItem {
	results := []ast.SelectItem{p.parseSelectItem()}
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
		if p.Token.Kind == token.TokenEOF || p.Token.Kind == "FROM" {
			break
		}
		results = append(results, p.parseSelectItem())
	}
	return results
}

// lookaheadStarModifierExcept is needed to distinct "* EXCEPT (columns)" and "* EXCEPT {ALL|DISTINCT}".
func (p *Parser) lookaheadStarModifierExcept() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if p.Token.Kind != "EXCEPT" {
		return false
	}
	p.nextToken()
	return p.Token.Kind == "("
}

func (p *Parser) tryParseStarModifierExcept() *ast.StarModifierExcept {
	if !p.lookaheadStarModifierExcept() {
		return nil
	}

	pos := p.expect("EXCEPT").Pos
	p.expect("(")
	columns := parseCommaSeparatedList(p, p.parseIdent)
	rparen := p.expect(")").Pos

	return &ast.StarModifierExcept{
		Except:  pos,
		Rparen:  rparen,
		Columns: columns,
	}
}

func (p *Parser) parseStarModifierReplaceItem() *ast.StarModifierReplaceItem {
	expr := p.parseExpr()
	p.expect("AS")
	name := p.parseIdent()

	return &ast.StarModifierReplaceItem{
		Expr: expr,
		Name: name,
	}
}

func (p *Parser) tryParseStarModifierReplace() *ast.StarModifierReplace {
	if !p.Token.IsKeywordLike("REPLACE") {
		return nil
	}

	pos := p.expectKeywordLike("REPLACE").Pos
	p.expect("(")
	columns := parseCommaSeparatedList(p, p.parseStarModifierReplaceItem)
	rparen := p.expect(")").Pos

	return &ast.StarModifierReplace{
		Replace: pos,
		Rparen:  rparen,
		Columns: columns,
	}
}

func (p *Parser) parseSelectItem() ast.SelectItem {
	if p.Token.Kind == "*" {
		pos := p.expect("*").Pos
		except := p.tryParseStarModifierExcept()
		replace := p.tryParseStarModifierReplace()
		return &ast.Star{
			Star:    pos,
			Except:  except,
			Replace: replace,
		}
	}

	expr := p.parseExpr()
	if as := p.tryParseAsAlias(withOptionalAs); as != nil {
		return &ast.Alias{
			Expr: expr,
			As:   as,
		}
	}

	if p.Token.Kind == "." {
		p.nextToken()
		pos := p.expect("*").Pos
		except := p.tryParseStarModifierExcept()
		replace := p.tryParseStarModifierReplace()
		return &ast.DotStar{
			Star:    pos,
			Expr:    expr,
			Except:  except,
			Replace: replace,
		}
	}

	return &ast.ExprSelectItem{
		Expr: expr,
	}
}

type withAs bool

const (
	withRequiredAs withAs = true
	withOptionalAs withAs = false
)

func (p *Parser) tryParseAsAlias(requiredAs withAs) *ast.AsAlias {
	pos := p.Token.Pos

	if p.Token.Kind == "AS" {
		p.nextToken()
		id := p.parseIdent()
		return &ast.AsAlias{
			As:    pos,
			Alias: id,
		}
	}

	if requiredAs {
		return nil
	}

	if p.Token.Kind == token.TokenIdent {
		id := p.parseIdent()
		return &ast.AsAlias{
			As:    token.InvalidPos,
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
	hint := p.tryParseHint()
	p.expect("BY")
	exprs := parseCommaSeparatedList(p, p.parseExpr)

	return &ast.GroupBy{
		Group: pos,
		Hint:  hint,
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

// parseQueryExprSuffix wraps QueryExpr if it is followed by ORDER BY, LIMIT, and/or pipe operators.
// It must not be called to *ast.Query itself because *ast.Query already consumes its suffix.
func (p *Parser) parseQueryExprSuffix(e ast.QueryExpr) ast.QueryExpr {
	// Query already consumes suffixes and currently it should be a logic bug.
	if _, ok := e.(*ast.Query); ok {
		panic(p.errorfAtToken(&p.Token, "invalid state of repeated processing of suffix of ast.Query. It is suspected a bug."))
	}

	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()
	forUpdate := p.tryParseForUpdate()

	pipeOps := p.parsePipeOperators()

	if orderBy == nil && limit == nil && forUpdate == nil && len(pipeOps) == 0 {
		return e
	}
	return &ast.Query{
		Query:         e,
		OrderBy:       orderBy,
		Limit:         limit,
		ForUpdate:     forUpdate,
		PipeOperators: pipeOps,
	}
}

func (p *Parser) tryParseAllOrDistinct() ast.AllOrDistinct {
	switch p.Token.Kind {
	case "ALL":
		p.nextToken()
		return ast.AllOrDistinctAll
	case "DISTINCT":
		p.nextToken()
		return ast.AllOrDistinctDistinct
	default:
		// not specified
		return ""
	}
}

func (p *Parser) parseAllOrDistinct() ast.AllOrDistinct {
	if p.Token.Kind != "ALL" && p.Token.Kind != "DISTINCT" {
		p.panicfAtToken(&p.Token, "expected token: ALL, DISTINCT, but: %s", p.Token.Kind)
	}

	return p.tryParseAllOrDistinct()
}
func (p *Parser) tryParseOrderBy() *ast.OrderBy {
	if p.Token.Kind != "ORDER" {
		return nil
	}

	pos := p.expect("ORDER").Pos
	p.expect("BY")

	items := parseCommaSeparatedList(p, p.parseOrderByItem)

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
			case p.Token.IsKeywordLike("LOOKUP"):
				p.nextToken()
				method = ast.LookupJoinMethod
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

		var cond ast.JoinCondition
		if op != ast.CrossJoin && op != ast.CommaJoin {
			switch right.(type) {
			case *ast.PathTableExpr, *ast.Unnest:
				cond = p.tryParseJoinCondition()
			default:
				cond = p.parseJoinCondition()
			}
		}

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
		as := p.tryParseAsAlias(withOptionalAs)
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
		return p.parseUnnestSuffix(expr, unnest, rparen)
	}

	if p.Token.Kind == token.TokenIdent {
		ids := p.parseIdentOrPath()
		if p.Token.Kind == "(" {
			return p.parseTVFCallExpr(ids)
		}
		if len(ids) == 1 {
			return p.parseTableNameSuffix(ids[0])
		}
		return p.parsePathTableExprSuffix(&ast.Path{Idents: ids})
	}

	panic(p.errorfAtToken(&p.Token, "expected token: (, UNNEST, <ident>, but: %s", p.Token.Kind))
}

func (p *Parser) parseTVFCallExpr(ids []*ast.Ident) *ast.TVFCallExpr {
	p.expect("(")

	var args []ast.TVFArg
	if p.Token.Kind != ")" {
		for !p.lookaheadNamedArg() {
			args = append(args, p.parseTVFArg())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}

	var namedArgs []*ast.NamedArg
	if p.lookaheadNamedArg() {
		namedArgs = parseCommaSeparatedList(p, p.parseNamedArg)
	}

	rparen := p.expect(")").Pos
	hint := p.tryParseHint()
	sample := p.tryParseTableSample()

	return &ast.TVFCallExpr{
		Rparen:    rparen,
		Name:      &ast.Path{Idents: ids},
		Args:      args,
		NamedArgs: namedArgs,
		Hint:      hint,
		Sample:    sample,
	}
}

func (p *Parser) parseTVFArg() ast.TVFArg {
	pos := p.Token.Pos
	switch {
	case p.Token.IsKeywordLike("TABLE"):
		p.nextToken()
		path := p.parsePath()

		return &ast.TableArg{
			Table: pos,
			Name:  path,
		}
	case p.Token.IsKeywordLike("MODEL"):
		p.nextToken()
		path := p.parsePath()

		return &ast.ModelArg{
			Model: pos,
			Name:  path,
		}
	default:
		return p.parseExprArg()
	}
}

func (p *Parser) parseIdentOrPath() []*ast.Ident {
	ids := []*ast.Ident{p.parseIdent()}
	for p.Token.Kind == "." {
		p.nextToken()
		ids = append(ids, p.parseIdent())
	}
	return ids
}

func (p *Parser) parseUnnestSuffix(expr ast.Expr, unnest, rparen token.Pos) ast.TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias(withOptionalAs)
	withOffset := p.tryParseWithOffset()

	return p.parseTableExprSuffix(&ast.Unnest{
		Unnest:     unnest,
		Rparen:     rparen,
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
	as := p.tryParseAsAlias(withOptionalAs)

	return &ast.WithOffset{
		With:   with,
		Offset: offset,
		As:     as,
	}
}

func (p *Parser) parseTableNameSuffix(id *ast.Ident) ast.TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias(withOptionalAs)
	return p.parseTableExprSuffix(&ast.TableName{
		Table: id,
		Hint:  hint,
		As:    as,
	})
}

func (p *Parser) parsePathTableExprSuffix(id *ast.Path) ast.TableExpr {
	hint := p.tryParseHint()
	as := p.tryParseAsAlias(withOptionalAs)
	withOffset := p.tryParseWithOffset()
	return p.parseTableExprSuffix(&ast.PathTableExpr{
		Path:       id,
		Hint:       hint,
		As:         as,
		WithOffset: withOffset,
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

func (p *Parser) tryParseJoinCondition() ast.JoinCondition {
	if p.Token.Kind != "ON" && p.Token.Kind != "USING" {
		return nil
	}
	return p.parseJoinCondition()
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
	idents := parseCommaSeparatedList(p, p.parseIdent)
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
	case *ast.PathTableExpr:
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

func (p *Parser) parseExpr() (expr ast.Expr) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			expr = p.handleParseExprError(r, l)
		}
	}()

	return p.parseOr()
}

func (p *Parser) parseOr() ast.Expr {
	expr := p.parseAnd()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "OR":
			op = ast.OpOr
		default:
			return expr
		}
		p.nextToken()
		expr = &ast.BinaryExpr{
			Left:  expr,
			Op:    op,
			Right: p.parseAnd(),
		}
	}
}

func (p *Parser) parseAnd() ast.Expr {
	expr := p.parseNot()
	for {
		var op ast.BinaryOp
		switch p.Token.Kind {
		case "AND":
			op = ast.OpAnd
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
			lexer := p.Clone()
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
			index := p.parseIndexSpecifier()
			rbrack := p.expect("]").Pos
			expr = &ast.IndexExpr{
				Rbrack: rbrack,
				Expr:   expr,
				Index:  index,
			}
		default:
			return expr
		}
	}
}

func (p *Parser) parseIndexSpecifier() ast.SubscriptSpecifier {
	pos := p.Token.Pos
	switch {
	case p.Token.IsIdent("OFFSET"), p.Token.IsIdent("ORDINAL"),
		p.Token.IsIdent("SAFE_OFFSET"), p.Token.IsIdent("SAFE_ORDINAL"):
		var keyword ast.PositionKeyword
		switch {
		case p.Token.IsIdent("OFFSET"):
			keyword = ast.PositionKeywordOffset
		case p.Token.IsIdent("ORDINAL"):
			keyword = ast.PositionKeywordOrdinal
		case p.Token.IsIdent("SAFE_OFFSET"):
			keyword = ast.PositionKeywordSafeOffset
		case p.Token.IsIdent("SAFE_ORDINAL"):
			keyword = ast.PositionKeywordSafeOrdinal

		}
		p.nextToken()
		p.expect("(")
		expr := p.parseExpr()
		rparen := p.expect(")").Pos

		return &ast.SubscriptSpecifierKeyword{
			KeywordPos: pos,
			Keyword:    keyword,
			Rparen:     rparen,
			Expr:       expr,
		}
	default:
		expr := p.parseExpr()
		return &ast.ExprArg{Expr: expr}
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
	case "IF":
		return p.parseIfExpr()
	case "CAST":
		return p.parseCastExpr()
	case "EXISTS":
		return p.parseExistsSubQuery()
	case "EXTRACT":
		return p.parseExtractExpr()
	case "WITH":
		return p.parseWithExpr()
	case "ARRAY":
		return p.parseArrayLiteralOrSubQuery()
	case "STRUCT":
		return p.parseStructLiteral()
	case "[":
		return p.parseSimpleArrayLiteral()
	case "(":
		return p.parseParenExpr()
	case "NEW":
		return p.parseNewConstructors()
	// In parser level, it is a valid ast.Expr, but it is semantically valid only in ast.BracedConstructorFieldExpr.
	case "{":
		return p.parseBracedConstructor()
	case "INTERVAL":
		return p.parseIntervalLiteral()
	case token.TokenIdent:
		id := p.Token
		switch {
		case id.IsKeywordLike("SAFE_CAST"):
			return p.parseCastExpr()
		case id.IsKeywordLike("REPLACE_FIELDS"):
			return p.parseReplaceFieldsExpr()
		}

		if p.lookaheadCallExpr() {
			return p.parseCallLike()
		}

		p.nextToken()
		switch p.Token.Kind {
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
			if id.IsKeywordLike("JSON") {
				return p.parseJSONLiteral(id)
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

func (p *Parser) lookaheadCallExpr() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	for {
		if p.Token.Kind != token.TokenIdent {
			return false
		}

		p.nextToken()
		switch p.Token.Kind {
		case "(":
			return true
		case ".":
			p.nextToken()
		default:
			return false
		}
	}
}

// parseCallLike parses after identifier part of function call like structures.
func (p *Parser) parseCallLike() ast.Expr {
	id := p.Token
	path := p.parsePath()
	p.expect("(")
	if len(path.Idents) == 1 && id.IsIdent("COUNT") && p.Token.Kind == "*" {
		p.nextToken()
		rparen := p.expect(")").Pos
		return &ast.CountStarExpr{
			Count:  path.Pos(),
			Rparen: rparen,
		}
	}

	distinct := false
	if p.Token.Kind == "DISTINCT" {
		p.nextToken()
		distinct = true
	}

	var args []ast.Arg
	if p.Token.Kind != ")" {
		for p.Token.Kind != token.TokenEOF && !p.lookaheadNamedArg() {
			args = append(args, p.parseArg())
			if p.Token.Kind != "," {
				break
			}
			p.nextToken()
		}
	}

	// https://github.com/google/zetasql/blob/master/docs/functions-reference.md#named-arguments
	// You cannot specify positional arguments after named arguments.
	var namedArgs []*ast.NamedArg
	for {
		namedArg := p.tryParseNamedArg()
		if namedArg == nil {
			break
		}
		namedArgs = append(namedArgs, namedArg)
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}

	nullHandling := p.tryParseNullHandlingModifier()
	having := p.tryParseHavingModifier()
	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()

	rparen := p.expect(")").Pos
	hint := p.tryParseHint()

	return &ast.CallExpr{
		Rparen:       rparen,
		Func:         path,
		Distinct:     distinct,
		Args:         args,
		NamedArgs:    namedArgs,
		NullHandling: nullHandling,
		Having:       having,
		OrderBy:      orderBy,
		Limit:        limit,
		Hint:         hint,
	}
}

func (p *Parser) lookaheadNamedArg() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if p.Token.Kind != token.TokenIdent {
		return false
	}
	p.parseIdent()
	return p.Token.Kind == "=>"
}

func (p *Parser) parseNamedArg() *ast.NamedArg {
	name := p.parseIdent()
	p.expect("=>")
	value := p.parseExpr()

	return &ast.NamedArg{
		Name:  name,
		Value: value,
	}
}

func (p *Parser) tryParseNamedArg() *ast.NamedArg {
	if !p.lookaheadNamedArg() {
		return nil
	}

	return p.parseNamedArg()
}

func (p *Parser) lookaheadLambdaArg() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	switch p.Token.Kind {
	case "(":
		p.nextToken()
		if p.Token.Kind != token.TokenIdent {
			return false
		}
		p.nextToken()

		for p.Token.Kind != ")" {
			if p.Token.Kind != "," {
				return false
			}
			p.nextToken()

			if p.Token.Kind != token.TokenIdent {
				return false
			}
			p.nextToken()
		}
		p.nextToken()

		return p.Token.Kind == "->"
	case token.TokenIdent:
		p.nextToken()
		return p.Token.Kind == "->"
	default:
		return false
	}
}

func (p *Parser) parseLambdaArg() *ast.LambdaArg {
	lparen := token.InvalidPos
	var args []*ast.Ident
	if p.Token.Kind == "(" {
		lparen = p.expect("(").Pos
		args = parseCommaSeparatedList(p, p.parseIdent)
		p.expect(")")
	} else {
		args = []*ast.Ident{p.parseIdent()}
	}

	p.expect("->")
	expr := p.parseExpr()

	return &ast.LambdaArg{
		Lparen: lparen,
		Args:   args,
		Expr:   expr,
	}
}

func (p *Parser) parseArg() ast.Arg {
	if s := p.tryParseSequenceArg(); s != nil {
		return s
	}
	if p.lookaheadLambdaArg() {
		return p.parseLambdaArg()
	}
	return p.parseExprArg()
}

func (p *Parser) tryParseSequenceArg() *ast.SequenceArg {
	if !p.Token.IsKeywordLike("SEQUENCE") {
		return nil
	}

	pos := p.Token.Pos
	p.nextToken()
	e := p.parseExpr()
	return &ast.SequenceArg{
		Sequence: pos,
		Expr:     e,
	}
}

func (p *Parser) parseExprArg() *ast.ExprArg {
	e := p.parseExpr()
	return &ast.ExprArg{
		Expr: e,
	}
}

func (p *Parser) tryParseHavingModifier() ast.HavingModifier {
	if p.Token.Kind != "HAVING" {
		return nil
	}

	having := p.expect("HAVING").Pos

	switch {
	case p.Token.IsKeywordLike("MIN"):
		p.nextToken()
		expr := p.parseExpr()

		return &ast.HavingMin{
			Having: having,
			Expr:   expr,
		}
	case p.Token.IsKeywordLike("MAX"):
		p.nextToken()
		expr := p.parseExpr()

		return &ast.HavingMax{
			Having: having,
			Expr:   expr,
		}
	default:
		panic(p.errorfAtToken(&p.Token, `expect MIN or MAX, but: %v`, p.Token.Kind))
	}
}

func (p *Parser) tryParseNullHandlingModifier() ast.NullHandlingModifier {
	switch p.Token.Kind {
	case "IGNORE":
		ignore := p.expect("IGNORE").Pos
		nulls := p.expect("NULLS").Pos

		return &ast.IgnoreNulls{
			Ignore: ignore,
			Nulls:  nulls,
		}
	case "RESPECT":
		respect := p.expect("RESPECT").Pos
		nulls := p.expect("NULLS").Pos

		return &ast.RespectNulls{
			Respect: respect,
			Nulls:   nulls,
		}
	default:
		return nil
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

func (p *Parser) parseIfExpr() *ast.IfExpr {
	pos := p.expect("IF").Pos
	p.expect("(")
	expr := p.parseExpr()
	p.expect(",")
	trueResult := p.parseExpr()
	p.expect(",")
	elseResult := p.parseExpr()
	rparen := p.expect(")").Pos

	return &ast.IfExpr{
		If:         pos,
		Rparen:     rparen,
		Expr:       expr,
		TrueResult: trueResult,
		ElseResult: elseResult,
	}
}

func (p *Parser) parseReplaceFieldsArg() *ast.ReplaceFieldsArg {
	expr := p.parseExpr()
	p.expect("AS")
	field := p.parsePath()

	return &ast.ReplaceFieldsArg{
		Expr:  expr,
		Field: field,
	}
}

func (p *Parser) parseReplaceFieldsExpr() *ast.ReplaceFieldsExpr {
	replaceFields := p.expectKeywordLike("REPLACE_FIELDS").Pos
	p.expect("(")
	expr := p.parseExpr()
	p.expect(",")
	fields := parseCommaSeparatedList(p, p.parseReplaceFieldsArg)
	rparen := p.expect(")").Pos

	return &ast.ReplaceFieldsExpr{
		ReplaceFields: replaceFields,
		Rparen:        rparen,
		Expr:          expr,
		Fields:        fields,
	}
}

func (p *Parser) parseCastExpr() *ast.CastExpr {
	if p.Token.Kind != "CAST" && !p.Token.IsKeywordLike("SAFE_CAST") {
		panic(p.errorfAtToken(&p.Token, `expected CAST keyword or SAFE_CAST pseudo keyword, but: %v`, p.Token.Kind))
	}

	var cast token.Pos
	var safe bool
	if p.Token.Kind == "CAST" {
		cast = p.expect("CAST").Pos
	} else {
		cast = p.expectKeywordLike("SAFE_CAST").Pos
		safe = true
	}

	p.expect("(")
	e := p.parseExpr()
	p.expect("AS")
	t := p.parseType()
	rparen := p.expect(")").Pos
	return &ast.CastExpr{
		Cast:   cast,
		Rparen: rparen,
		Safe:   safe,
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

func (p *Parser) parseWithExprVar() *ast.WithExprVar {
	name := p.parseIdent()
	p.expect("AS")
	expr := p.parseExpr()

	return &ast.WithExprVar{
		Expr: expr,
		Name: name,
	}
}

func (p *Parser) lookaheadWithExprVar() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if p.Token.Kind != token.TokenIdent {
		return false
	}

	p.parseIdent()
	return p.Token.Kind == "AS"
}

func (p *Parser) parseWithExpr() *ast.WithExpr {
	with := p.expect("WITH").Pos
	p.expect("(")

	var vars []*ast.WithExprVar
	for p.lookaheadWithExprVar() {
		vars = append(vars, p.parseWithExprVar())
		p.expect(",")
	}

	expr := p.parseExpr()
	rparen := p.expect(")").Pos
	return &ast.WithExpr{
		With:   with,
		Rparen: rparen,
		Vars:   vars,
		Expr:   expr,
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
		p.panicfAtToken(&paren, "cannot parse (...) as expression, tuple struct literal or subquery")
	}

	p.expect(",")

	values := append([]ast.Expr{expr}, parseCommaSeparatedList(p, p.parseExpr)...)
	rparen := p.expect(")").Pos

	return &ast.TupleStructLiteral{
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
		Array:  token.InvalidPos,
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

func (p *Parser) parseStructLiteral() ast.Expr {
	// Note that tuple struct syntax is handled in parseParenExpr.

	pos := p.expect("STRUCT").Pos
	if p.Token.Kind == "<" || p.Token.Kind == "<>" {
		return p.parseTypedStructLiteral(pos)
	}

	if p.Token.Kind != "(" {
		p.panicfAtToken(&p.Token, "expected token: <, <>, ( but: %s", p.Token.Kind)
	}

	return p.parseTypelessStructLiteral(pos)
}

func (p *Parser) parseTypedStructLiteral(pos token.Pos) *ast.TypedStructLiteral {
	fields, _ := p.parseStructTypeFields()

	p.expect("(")
	var values []ast.Expr
	if p.Token.Kind != ")" {
		values = parseCommaSeparatedList(p, p.parseExpr)
	}
	rparen := p.expect(")").Pos

	return &ast.TypedStructLiteral{
		Struct: pos,
		Rparen: rparen,
		Fields: fields,
		Values: values,
	}
}

func (p *Parser) parseTypelessStructLiteral(pos token.Pos) *ast.TypelessStructLiteral {
	p.expect("(")
	var values []ast.TypelessStructLiteralArg
	if p.Token.Kind != ")" {
		values = parseCommaSeparatedList(p, p.parseTypelessStructLiteralArg)
	}
	rparen := p.expect(")").Pos

	return &ast.TypelessStructLiteral{
		Struct: pos,
		Rparen: rparen,
		Values: values,
	}
}

func (p *Parser) parseTypelessStructLiteralArg() ast.TypelessStructLiteralArg {
	e := p.parseExpr()
	as := p.tryParseAsAlias(withRequiredAs)
	if as != nil {
		return &ast.Alias{
			Expr: e,
			As:   as,
		}
	}
	return &ast.ExprArg{
		Expr: e,
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

func (p *Parser) parseJSONLiteral(id token.Token) *ast.JSONLiteral {
	s := p.parseStringLiteral()
	return &ast.JSONLiteral{
		JSON:  id.Pos,
		Value: s,
	}
}

var validDateTimePartNames = []ast.DateTimePart{
	ast.DateTimePartYear,
	ast.DateTimePartMonth,
	ast.DateTimePartDay,
	ast.DateTimePartDayOfWeek,
	ast.DateTimePartDayOfYear,
	ast.DateTimePartQuarter,
	ast.DateTimePartHour,
	ast.DateTimePartMinute,
	ast.DateTimePartSecond,
	ast.DateTimePartMillisecond,
	ast.DateTimePartMicrosecond,
	ast.DateTimePartNanosecond,
	ast.DateTimePartWeek,
	ast.DateTimePartISOYear,
	ast.DateTimePartISOWeek,
	ast.DateTimePartDate,
	ast.DateTimePartDateTime,
	ast.DateTimePartTime,
}

func (p *Parser) parseDateTimePart() (part ast.DateTimePart, pos, end token.Pos) {
	// Note: This function will support WEEKDAY(weekday) when it is released in Spanner

	ident := p.parseIdent()

	for _, dateTimePart := range validDateTimePartNames {
		if char.EqualFold(ident.Name, string(dateTimePart)) {
			return dateTimePart, ident.Pos(), ident.End()
		}
	}

	panic(p.errorfAtPosition(ident.Pos(), ident.End(), "invalid date time part: %s", ident.Name))
}

func (p *Parser) parseIntervalLiteral() ast.Expr {
	interval := p.expect("INTERVAL").Pos
	expr := p.parseExpr()

	// string parameter is not supported, so ast.Param will be caught as ast.IntValue.
	switch e := expr.(type) {
	case ast.IntValue:
		unit, _, end := p.parseDateTimePart()

		return &ast.IntervalLiteralSingle{
			Interval:        interval,
			Value:           e,
			DateTimePart:    unit,
			DateTimePartEnd: end,
		}
	case *ast.StringLiteral:
		starting, _, _ := p.parseDateTimePart()

		p.expect("TO")
		ending, _, endingEnd := p.parseDateTimePart()

		return &ast.IntervalLiteralRange{
			Interval:              interval,
			Value:                 e,
			StartingDateTimePart:  starting,
			EndingDateTimePart:    ending,
			EndingDateTimePartEnd: endingEnd,
		}
	default:
		panic(p.errorfAtToken(&p.Token, `expect int64_expression or datetime_parts_string, but: %v`, p.Token.Kind))
	}
}

func (p *Parser) lookaheadSubQuery() bool {
	lexer := p.Clone()
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

// ================================================================================
//
// Type
//
// ================================================================================

func (p *Parser) parseNamedType() *ast.NamedType {
	path := p.parseIdentOrPath()
	return &ast.NamedType{Path: path}
}

func (p *Parser) parseType() (t ast.Type) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			t = p.handleParseTypeError(r, l)
		}
	}()

	switch p.Token.Kind {
	case token.TokenIdent, "INTERVAL":
		if !p.lookaheadSimpleType() {
			return p.parseNamedType()
		}
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
	"FLOAT32",
	"FLOAT64",
	"DATE",
	"TIMESTAMP",
	"NUMERIC",
	"STRING",
	"BYTES",
	"JSON",
	"TOKENLIST",
	"INTERVAL",
}

func (p *Parser) parseSimpleType() *ast.SimpleType {
	// Only INTERVAL is a keyword.
	if p.Token.Kind == "INTERVAL" {
		pos := p.expect("INTERVAL").Pos

		return &ast.SimpleType{
			NamePos: pos,
			Name:    ast.IntervalTypeName,
		}
	}

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
	if p.Token.Kind != "<" && p.Token.Kind != "<>" {
		p.panicfAtToken(&p.Token, "expected token: <, <>, but: %s", p.Token.Kind)
	}
	fields, gt := p.parseStructTypeFields()
	return &ast.StructType{
		Struct: pos,
		Gt:     gt,
		Fields: fields,
	}
}

func (p *Parser) parseStructTypeFields() ([]*ast.StructField, token.Pos) {
	if p.Token.Kind == "<>" {
		gt := p.Token.Pos + 1
		p.nextToken()
		return nil, gt
	}

	var fields []*ast.StructField
	p.expect("<")
	if p.Token.Kind != ">" && p.Token.Kind != ">>" {
		fields = parseCommaSeparatedList(p, p.parseFieldType)
	}

	var gt token.Pos
	if p.Token.Kind == ">>" {
		p.Token.Kind = ">"
		p.Token.Raw = ">"
		gt = p.Token.Pos
		p.Token.Pos += 1
	} else {
		gt = p.expect(">").Pos
	}

	return fields, gt
}

func (p *Parser) parseNewConstructor(newPos token.Pos, namedType *ast.NamedType) *ast.NewConstructor {
	p.expect("(")

	// Args can be empty like `NEW pkg.TypeName ()`.
	var args []ast.NewConstructorArg
	if p.Token.Kind != ")" {
		args = parseCommaSeparatedList(p, p.parseNewConstructorArg)
	}

	rparen := p.expect(")").Pos
	return &ast.NewConstructor{
		New:    newPos,
		Type:   namedType,
		Args:   args,
		Rparen: rparen,
	}
}

func (p *Parser) parseNewConstructorArg() ast.NewConstructorArg {
	// Currently, this method's contents are the same as `parseTypelessStructLiteralArg`.
	// It exists as an individual method for future extensibility.

	e := p.parseExpr()
	as := p.tryParseAsAlias(withRequiredAs)
	if as != nil {
		return &ast.Alias{
			Expr: e,
			As:   as,
		}
	}
	return &ast.ExprArg{
		Expr: e,
	}
}

func (p *Parser) parseBracedNewConstructorField() *ast.BracedConstructorField {
	name := p.parseIdent()
	var fieldValue ast.BracedConstructorFieldValue
	switch p.Token.Kind {
	case ":":
		colon := p.expect(":").Pos
		expr := p.parseExpr()
		fieldValue = &ast.BracedConstructorFieldValueExpr{Colon: colon, Expr: expr}
	case "{":
		fieldValue = p.parseBracedConstructor()
	}
	return &ast.BracedConstructorField{Name: name, Value: fieldValue}
}

func (p *Parser) parseBracedConstructor() *ast.BracedConstructor {
	lbrace := p.expect("{").Pos

	// Braced constructor permits empty.
	var fields []*ast.BracedConstructorField
	for p.Token.Kind != "}" {
		if p.Token.Kind != token.TokenIdent {
			p.panicfAtToken(&p.Token, "expect <ident>, but %v", p.Token.Kind)
		}
		fields = append(fields, p.parseBracedNewConstructorField())

		// It is an optional comma.
		if p.Token.Kind == "," {
			p.nextToken()
		}
	}

	rbrace := p.expect("}").Pos

	return &ast.BracedConstructor{
		Lbrace: lbrace,
		Rbrace: rbrace,
		Fields: fields,
	}
}

func (p *Parser) parseBracedNewConstructor(newPos token.Pos, namedType *ast.NamedType) *ast.BracedNewConstructor {
	body := p.parseBracedConstructor()
	return &ast.BracedNewConstructor{
		New:  newPos,
		Type: namedType,
		Body: body,
	}
}

func (p *Parser) parseNewConstructors() ast.Expr {
	newPos := p.expect("NEW").Pos
	namedType := p.parseNamedType()

	switch p.Token.Kind {
	case "(":
		return p.parseNewConstructor(newPos, namedType)
	case "{":
		return p.parseBracedNewConstructor(newPos, namedType)
	default:
		p.panicfAtToken(&p.Token, `expect '{' or '(', but %v`, p.Token.Kind)
	}
	return nil
}

func (p *Parser) parseFieldType() *ast.StructField {
	lexer := p.Clone()
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
	return p.Token.Kind == token.TokenIdent || p.Token.Kind == "INTERVAL" || p.Token.Kind == "ARRAY" || p.Token.Kind == "STRUCT"
}

func (p *Parser) lookaheadSimpleType() bool {
	if p.Token.Kind == "INTERVAL" {
		return true
	}

	if p.Token.Kind != token.TokenIdent {
		return false
	}

	id := p.Token

	for _, name := range simpleTypes {
		if id.IsIdent(name) {
			return true
		}
	}
	return false
}

// ================================================================================
//
// DDL
//
// ================================================================================

func (p *Parser) parseDDL() (ddl ast.DDL) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			ddl = &ast.BadDDL{BadNode: p.handleParseStatementError(r, l)}
		}
	}()

	pos := p.Token.Pos
	switch {
	case p.Token.Kind == "CREATE":
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("SCHEMA"):
			return p.parseCreateSchema(pos, false)
		case p.Token.IsKeywordLike("DATABASE"):
			return p.parseCreateDatabase(pos)
		case p.Token.IsKeywordLike("LOCALITY"):
			return p.parseCreateLocalityGroup(pos)
		case p.Token.IsKeywordLike("PLACEMENT"):
			return p.parseCreatePlacement(pos)
		case p.Token.Kind == "PROTO":
			return p.parseCreateProtoBundle(pos)
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseCreateTable(pos)
		case p.Token.IsKeywordLike("SEQUENCE"):
			return p.parseCreateSequence(pos)
		case p.Token.IsKeywordLike("VIEW"):
			return p.parseCreateView(pos, false)
		case p.Token.IsKeywordLike("INDEX") || p.Token.IsKeywordLike("UNIQUE") || p.Token.IsKeywordLike("NULL_FILTERED"):
			return p.parseCreateIndex(pos)
		case p.Token.IsKeywordLike("VECTOR"):
			return p.parseCreateVectorIndex(pos)
		case p.Token.IsKeywordLike("SEARCH"):
			return p.parseCreateSearchIndex(pos)
		case p.Token.IsKeywordLike("ROLE"):
			return p.parseCreateRole(pos)
		case p.Token.IsKeywordLike("CHANGE"):
			return p.parseCreateChangeStream(pos)
		case p.Token.IsKeywordLike("MODEL"):
			return p.parseCreateModel(pos, false)
		case p.Token.IsKeywordLike("FUNCTION"):
			return p.parseCreateFunction(pos, false)
		case p.Token.Kind == "OR":
			p.expect("OR")
			p.expectKeywordLike("REPLACE")
			switch {
			case p.Token.IsKeywordLike("VIEW"):
				return p.parseCreateView(pos, true)
			case p.Token.IsKeywordLike("SCHEMA"):
				return p.parseCreateSchema(pos, true)
			case p.Token.IsKeywordLike("MODEL"):
				return p.parseCreateModel(pos, true)
			case p.Token.IsKeywordLike("PROPERTY"):
				return p.parseCreatePropertyGraph(pos, true)
			case p.Token.IsKeywordLike("FUNCTION"):
				return p.parseCreateFunction(pos, true)
			}
		case p.Token.IsKeywordLike("PROPERTY"):
			return p.parseCreatePropertyGraph(pos, false)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: DATABASE, TABLE, INDEX, UNIQUE, NULL_FILTERED, ROLE, CHANGE, OR but: %s", p.Token.AsString)
	case p.Token.IsKeywordLike("ALTER"):
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseAlterTable(pos)
		case p.Token.IsKeywordLike("DATABASE"):
			return p.parseAlterDatabase(pos)
		case p.Token.IsKeywordLike("LOCALITY"):
			return p.parseAlterLocalityGroup(pos)
		case p.Token.Kind == "PROTO":
			return p.parseAlterProtoBundle(pos)
		case p.Token.IsKeywordLike("INDEX"):
			return p.parseAlterIndex(pos)
		case p.Token.IsKeywordLike("SEARCH"):
			return p.parseAlterSearchIndex(pos)
		case p.Token.IsKeywordLike("SEQUENCE"):
			return p.parseAlterSequence(pos)
		case p.Token.IsKeywordLike("CHANGE"):
			return p.parseAlterChangeStream(pos)
		case p.Token.IsKeywordLike("STATISTICS"):
			return p.parseAlterStatistics(pos)
		case p.Token.IsKeywordLike("MODEL"):
			return p.parseAlterModel(pos)
		case p.Token.IsKeywordLike("VECTOR"):
			return p.parseAlterVectorIndex(pos)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: TABLE, CHANGE, but: %s", p.Token.AsString)
	case p.Token.IsKeywordLike("DROP"):
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("SCHEMA"):
			return p.parseDropSchema(pos)
		case p.Token.IsKeywordLike("LOCALITY"):
			return p.parseDropLocalityGroup(pos)
		case p.Token.Kind == "PROTO":
			return p.parseDropProtoBundle(pos)
		case p.Token.IsKeywordLike("TABLE"):
			return p.parseDropTable(pos)
		case p.Token.IsKeywordLike("INDEX"):
			return p.parseDropIndex(pos)
		case p.Token.IsKeywordLike("SEARCH"):
			return p.parseDropSearchIndex(pos)
		case p.Token.IsKeywordLike("VECTOR"):
			return p.parseDropVectorIndex(pos)
		case p.Token.IsKeywordLike("SEQUENCE"):
			return p.parseDropSequence(pos)
		case p.Token.IsKeywordLike("VIEW"):
			return p.parseDropView(pos)
		case p.Token.IsKeywordLike("ROLE"):
			return p.parseDropRole(pos)
		case p.Token.IsKeywordLike("CHANGE"):
			return p.parseDropChangeStream(pos)
		case p.Token.IsKeywordLike("MODEL"):
			return p.parseDropModel(pos)
		case p.Token.IsKeywordLike("PROPERTY"):
			return p.parseDropPropertyGraph(pos)
		case p.Token.IsKeywordLike("FUNCTION"):
			return p.parseDropFunction(pos)
		}
		p.panicfAtToken(&p.Token, "expected pseudo keyword: TABLE, INDEX, ROLE, CHANGE, MODEL, but: %s", p.Token.AsString)
	case p.Token.IsKeywordLike("RENAME"):
		p.nextToken()
		return p.parseRenameTable(pos)
	case p.Token.IsKeywordLike("GRANT"):
		p.nextToken()
		return p.parseGrant(pos)
	case p.Token.IsKeywordLike("REVOKE"):
		p.nextToken()
		return p.parseRevoke(pos)
	case p.Token.IsKeywordLike("ANALYZE"):
		return p.parseAnalyze()
	}

	if p.Token.Kind != token.TokenIdent {
		panic(p.errorfAtToken(&p.Token, "expected token: CREATE, <ident>, but: %s", p.Token.Kind))
	}

	panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: ALTER, DROP, but: %s", p.Token.AsString))
}

func (p *Parser) parseCreateSchema(pos token.Pos, orReplace bool) *ast.CreateSchema {
	p.expectKeywordLike("SCHEMA")
	ifNotExists := p.parseIfNotExists()
	name := p.parseIdent()
	return &ast.CreateSchema{
		Create:      pos,
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
	}
}

func (p *Parser) parseDropSchema(pos token.Pos) *ast.DropSchema {
	p.expectKeywordLike("SCHEMA")
	ifExists := p.parseIfExists()
	name := p.parseIdent()
	return &ast.DropSchema{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) parseCreateDatabase(pos token.Pos) *ast.CreateDatabase {
	p.expectKeywordLike("DATABASE")
	name := p.parseIdent()

	return &ast.CreateDatabase{
		Create: pos,
		Name:   name,
	}
}

func (p *Parser) parseAlterDatabase(pos token.Pos) *ast.AlterDatabase {
	p.expectKeywordLike("DATABASE")
	name := p.parseIdent()
	p.expect("SET")
	options := p.parseOptions()

	return &ast.AlterDatabase{
		Alter:   pos,
		Name:    name,
		Options: options,
	}
}

func (p *Parser) parseCreateLocalityGroup(pos token.Pos) *ast.CreateLocalityGroup {
	p.expectKeywordLike("LOCALITY")
	p.expect("GROUP")
	name := p.parseIdent()

	options := p.tryParseOptions()

	return &ast.CreateLocalityGroup{
		Create:  pos,
		Name:    name,
		Options: options,
	}
}

func (p *Parser) parseAlterLocalityGroup(pos token.Pos) *ast.AlterLocalityGroup {
	p.expectKeywordLike("LOCALITY")
	p.expect("GROUP")
	name := p.parseIdent()

	p.expect("SET")
	options := p.parseOptions()

	return &ast.AlterLocalityGroup{
		Alter:   pos,
		Name:    name,
		Options: options,
	}
}

func (p *Parser) parseDropLocalityGroup(pos token.Pos) *ast.DropLocalityGroup {
	p.expectKeywordLike("LOCALITY")
	p.expect("GROUP")
	name := p.parseIdent()

	return &ast.DropLocalityGroup{
		Drop: pos,
		Name: name,
	}
}

func (p *Parser) parseCreatePlacement(pos token.Pos) *ast.CreatePlacement {
	p.expectKeywordLike("PLACEMENT")
	name := p.parseIdent()
	options := p.parseOptions()

	return &ast.CreatePlacement{
		Create:  pos,
		Name:    name,
		Options: options,
	}
}

func (p *Parser) parseProtoBundleTypes() *ast.ProtoBundleTypes {
	lparen := p.expect("(").Pos

	var types []*ast.NamedType
	if p.Token.Kind != ")" {
		for {
			types = append(types, p.parseNamedType())

			if p.Token.Kind != "," {
				break
			}
			p.nextToken()

			// case of trailing comma
			if p.Token.Kind == ")" {
				break
			}
		}
	}
	rparen := p.expect(")").Pos

	return &ast.ProtoBundleTypes{
		Lparen: lparen,
		Rparen: rparen,
		Types:  types,
	}
}

func (p *Parser) parseCreateProtoBundle(pos token.Pos) *ast.CreateProtoBundle {
	p.expect("PROTO")
	p.expectKeywordLike("BUNDLE")
	types := p.parseProtoBundleTypes()

	return &ast.CreateProtoBundle{
		Create: pos,
		Types:  types,
	}
}

func (p *Parser) parseAlterProtoBundle(pos token.Pos) *ast.AlterProtoBundle {
	p.expect("PROTO")
	bundle := p.expectKeywordLike("BUNDLE").Pos
	insert := p.tryParseAlterProtoBundleInsert()
	update := p.tryParseAlterProtoBundleUpdate()
	delete := p.tryParseAlterProtoBundleDelete()

	return &ast.AlterProtoBundle{
		Alter:  pos,
		Bundle: bundle,
		Insert: insert,
		Update: update,
		Delete: delete,
	}
}

func (p *Parser) tryParseAlterProtoBundleInsert() *ast.AlterProtoBundleInsert {
	if !p.Token.IsKeywordLike("INSERT") {
		return nil
	}

	pos := p.expectKeywordLike("INSERT").Pos
	types := p.parseProtoBundleTypes()

	return &ast.AlterProtoBundleInsert{
		Insert: pos,
		Types:  types,
	}
}

func (p *Parser) tryParseAlterProtoBundleUpdate() *ast.AlterProtoBundleUpdate {
	if !p.Token.IsKeywordLike("UPDATE") {
		return nil
	}

	pos := p.expectKeywordLike("UPDATE").Pos
	types := p.parseProtoBundleTypes()

	return &ast.AlterProtoBundleUpdate{
		Update: pos,
		Types:  types,
	}
}

func (p *Parser) tryParseAlterProtoBundleDelete() *ast.AlterProtoBundleDelete {
	if !p.Token.IsKeywordLike("DELETE") {
		return nil
	}

	pos := p.expectKeywordLike("DELETE").Pos
	types := p.parseProtoBundleTypes()

	return &ast.AlterProtoBundleDelete{
		Delete: pos,
		Types:  types,
	}
}

func (p *Parser) parseDropProtoBundle(pos token.Pos) *ast.DropProtoBundle {
	p.expect("PROTO")
	bundle := p.expectKeywordLike("BUNDLE").Pos

	return &ast.DropProtoBundle{
		Drop:   pos,
		Bundle: bundle,
	}
}
func (p *Parser) parseCreateTable(pos token.Pos) *ast.CreateTable {
	p.expectKeywordLike("TABLE")
	ifNotExists := p.parseIfNotExists()
	name := p.parsePath()

	// This loop allows parsing trailing comma intentionally.
	p.expect("(")
	var columns []*ast.ColumnDef
	var constraints []*ast.TableConstraint
	var synonyms []*ast.Synonym
	for p.Token.Kind != token.TokenEOF {
		if p.Token.Kind == ")" {
			break
		}
		switch {
		case p.Token.IsKeywordLike("CONSTRAINT"):
			constraints = append(constraints, p.parseConstraint())
		case p.Token.IsKeywordLike("FOREIGN"):
			fk := p.parseForeignKey()
			constraints = append(constraints, &ast.TableConstraint{
				ConstraintPos: token.InvalidPos,
				Constraint:    fk,
			})
		case p.Token.IsKeywordLike("CHECK"):
			c := p.parseCheck()
			constraints = append(constraints, &ast.TableConstraint{
				ConstraintPos: token.InvalidPos,
				Constraint:    c,
			})
		case p.Token.IsKeywordLike("SYNONYM"):
			synonym := p.parseSynonym()
			synonyms = append(synonyms, synonym)
		default:
			columns = append(columns, p.parseColumnDef())
		}
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}
	rparen := p.expect(")").Pos

	// PRIMARY KEY clause is now optional
	var keys []*ast.IndexKey
	primaryKeyRparen := token.InvalidPos
	if p.Token.IsKeywordLike("PRIMARY") {
		p.nextToken()
		p.expectKeywordLike("KEY")

		p.expect("(")
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
		primaryKeyRparen = p.expect(")").Pos
	}

	cluster := p.tryParseCluster()
	rdp := p.tryParseCreateRowDeletionPolicy()

	var options *ast.Options
	if p.Token.Kind == "," {
		p.nextToken()
		options = p.parseOptions()
	}

	return &ast.CreateTable{
		Create:            pos,
		Rparen:            rparen,
		PrimaryKeyRparen:  primaryKeyRparen,
		IfNotExists:       ifNotExists,
		Name:              name,
		Columns:           columns,
		TableConstraints:  constraints,
		Synonyms:          synonyms,
		PrimaryKeys:       keys,
		Cluster:           cluster,
		RowDeletionPolicy: rdp,
		Options:           options,
	}
}

func (p *Parser) parsePath() *ast.Path {
	path := p.parseIdentOrPath()
	return &ast.Path{Idents: path}
}

func (p *Parser) parseSequenceParams() []ast.SequenceParam {
	var params []ast.SequenceParam
	for {
		param := p.tryParseSequenceParam()
		if param == nil {
			return params
		}

		params = append(params, param)
	}
}

func (p *Parser) parseCreateSequence(pos token.Pos) *ast.CreateSequence {
	p.expectKeywordLike("SEQUENCE")
	ifNotExists := p.parseIfNotExists()
	name := p.parsePath()
	params := p.parseSequenceParams()
	options := p.parseOptions()

	return &ast.CreateSequence{
		Create:      pos,
		Name:        name,
		IfNotExists: ifNotExists,
		Params:      params,
		Options:     options,
	}
}

func (p *Parser) parseCreateView(pos token.Pos, orReplace bool) *ast.CreateView {
	p.expectKeywordLike("VIEW")

	name := p.parsePath()
	securityType := p.parseSqlSecurity()

	p.expect("AS")

	query := p.parseQueryExpr()

	return &ast.CreateView{
		Create:       pos,
		Name:         name,
		OrReplace:    orReplace,
		SecurityType: securityType,
		Query:        query,
	}
}

func (p *Parser) parseDropView(pos token.Pos) *ast.DropView {
	p.expectKeywordLike("VIEW")
	name := p.parsePath()

	return &ast.DropView{
		Drop: pos,
		Name: name,
	}
}

func (p *Parser) parseIdentityColumn() *ast.IdentityColumn {
	pos := p.expectKeywordLike("GENERATED").Pos
	p.expect("BY")
	p.expect("DEFAULT")
	p.expect("AS")
	identity := p.expectKeywordLike("IDENTITY").Pos

	var params []ast.SequenceParam
	rparen := token.InvalidPos
	if p.Token.Kind == "(" {
		p.nextToken()
		params = p.parseSequenceParams()
		rparen = p.expect(")").Pos
	}

	return &ast.IdentityColumn{
		Generated: pos,
		Identity:  identity,
		Rparen:    rparen,
		Params:    params,
	}
}

func (p *Parser) parseColumnDef() *ast.ColumnDef {
	name := p.parseIdent()
	t, notNull, null := p.parseTypeNotNull()

	var defaultSemantics ast.ColumnDefaultSemantics
	switch {
	case p.Token.Kind == "DEFAULT":
		defaultSemantics = p.parseColumnDefaultExpr()
	case p.Token.Kind == "AS":
		defaultSemantics = p.parseGeneratedColumnExpr()
	case p.Token.IsKeywordLike("GENERATED"):
		defaultSemantics = p.parseIdentityColumn()
	case p.Token.IsKeywordLike("AUTO_INCREMENT"):
		pos := p.expectKeywordLike("AUTO_INCREMENT").Pos
		defaultSemantics = &ast.AutoIncrement{AutoIncrement: pos}
	}

	hiddenPos := token.InvalidPos
	if p.Token.IsKeywordLike("HIDDEN") {
		hiddenPos = p.expectKeywordLike("HIDDEN").Pos
	}

	key := token.InvalidPos
	if p.Token.IsKeywordLike("PRIMARY") {
		p.nextToken()
		key = p.expectKeywordLike("KEY").Pos
	}

	options := p.tryParseOptions()

	return &ast.ColumnDef{
		Null:             null,
		Key:              key,
		Name:             name,
		Type:             t,
		NotNull:          notNull,
		PrimaryKey:       !key.Invalid(),
		DefaultSemantics: defaultSemantics,
		Hidden:           hiddenPos,
		Options:          options,
	}
}

func (p *Parser) parseConstraint() *ast.TableConstraint {
	pos := p.expectKeywordLike("CONSTRAINT").Pos
	name := p.parseIdent()
	var c ast.Constraint
	switch {
	case p.Token.IsKeywordLike("FOREIGN"):
		c = p.parseForeignKey()
	case p.Token.IsKeywordLike("CHECK"):
		c = p.parseCheck()
	default:
		panic(p.errorfAtToken(&p.Token, "unknown constraint %s", p.Token.AsString))
	}
	return &ast.TableConstraint{
		ConstraintPos: pos,
		Name:          name,
		Constraint:    c,
	}
}

func (p *Parser) parseForeignKey() *ast.ForeignKey {
	pos := p.expectKeywordLike("FOREIGN").Pos
	p.expectKeywordLike("KEY")

	p.expect("(")
	columns := parseCommaSeparatedList(p, p.parseIdent)
	p.expect(")")
	p.expectKeywordLike("REFERENCES")
	refTable := p.parsePath()

	p.expect("(")
	refColumns := parseCommaSeparatedList(p, p.parseIdent)
	rparen := p.expect(")").Pos

	onDelete, onDeleteEnd := p.tryParseOnDeleteAction()

	enforced := token.InvalidPos
	var enforcement ast.Enforcement
	switch {
	case p.Token.Kind == "NOT":
		p.nextToken()
		enforced = p.expectKeywordLike("ENFORCED").Pos

		enforcement = ast.NotEnforced
	case p.Token.IsKeywordLike("ENFORCED"):
		enforced = p.expectKeywordLike("ENFORCED").Pos

		enforcement = ast.Enforced
	}

	return &ast.ForeignKey{
		Foreign:          pos,
		Rparen:           rparen,
		OnDeleteEnd:      onDeleteEnd,
		Enforced:         enforced,
		Columns:          columns,
		ReferenceTable:   refTable,
		ReferenceColumns: refColumns,
		OnDelete:         onDelete,
		Enforcement:      enforcement,
	}
}

func (p *Parser) parseCheck() *ast.Check {
	pos := p.expectKeywordLike("CHECK").Pos
	p.expect("(")
	expr := p.parseExpr()
	rparen := p.expect(")").Pos
	return &ast.Check{
		Check:  pos,
		Rparen: rparen,
		Expr:   expr,
	}
}

func (p *Parser) parseSynonym() *ast.Synonym {
	pos := p.expectKeywordLike("SYNONYM").Pos
	p.expect("(")
	name := p.parseIdent()
	rparen := p.expect(")").Pos

	return &ast.Synonym{
		Synonym: pos,
		Rparen:  rparen,
		Name:    name,
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

func (p *Parser) tryParseColumnDefaultExpr() *ast.ColumnDefaultExpr {
	if p.Token.Kind != "DEFAULT" {
		return nil
	}

	return p.parseColumnDefaultExpr()
}

func (p *Parser) parseColumnDefaultExpr() *ast.ColumnDefaultExpr {
	def := p.expect("DEFAULT").Pos
	p.expect("(")
	expr := p.parseExpr()
	rparen := p.expect(")").Pos

	return &ast.ColumnDefaultExpr{
		Default: def,
		Rparen:  rparen,
		Expr:    expr,
	}
}

func (p *Parser) parseGeneratedColumnExpr() *ast.GeneratedColumnExpr {
	pos := p.expect("AS").Pos
	p.expect("(")
	expr := p.parseExpr()
	rparen := p.expect(")").Pos

	stored := token.InvalidPos
	if p.Token.IsKeywordLike("STORED") {
		stored = p.expectKeywordLike("STORED").Pos
	}

	return &ast.GeneratedColumnExpr{
		As:     pos,
		Stored: stored,
		Rparen: rparen,
		Expr:   expr,
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
	lexer := p.Clone()
	pos := p.expect(",").Pos
	if !p.Token.IsKeywordLike("INTERLEAVE") {
		p.Lexer = lexer
		return nil
	}
	p.nextToken()
	p.expect("IN")

	var enforced bool
	if p.Token.IsKeywordLike("PARENT") {
		p.nextToken()
		enforced = true
	}

	name := p.parsePath()

	onDelete, onDeleteEnd := p.tryParseOnDeleteAction()

	return &ast.Cluster{
		Comma:       pos,
		OnDeleteEnd: onDeleteEnd,
		TableName:   name,
		Enforced:    enforced,
		OnDelete:    onDelete,
	}
}

func (p *Parser) tryParseCreateRowDeletionPolicy() *ast.CreateRowDeletionPolicy {
	if p.Token.Kind != "," {
		return nil
	}

	lexer := p.Clone()
	pos := p.expect(",").Pos
	if !p.Token.IsKeywordLike("ROW") {
		p.Lexer = lexer
		return nil
	}

	rdp := p.parseRowDeletionPolicy()
	return &ast.CreateRowDeletionPolicy{
		Comma:             pos,
		RowDeletionPolicy: rdp,
	}
}

func (p *Parser) parseRowDeletionPolicy() *ast.RowDeletionPolicy {
	pos := p.expectKeywordLike("ROW").Pos
	p.expectKeywordLike("DELETION")
	p.expectKeywordLike("POLICY")
	p.expect("(")
	p.expectKeywordLike("OLDER_THAN")
	p.expect("(")
	timestampColumn := p.parseIdent()
	p.expect(",")
	p.expect("INTERVAL")
	numDays := p.parseIntLiteral()
	p.expectKeywordLike("DAY")
	p.expect(")")
	rparen := p.expect(")").Pos
	return &ast.RowDeletionPolicy{
		Row:        pos,
		ColumnName: timestampColumn,
		NumDays:    numDays,
		Rparen:     rparen,
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

func (p *Parser) parseOptionsDef() *ast.OptionsDef {
	key := p.parseIdent()
	p.expect("=")
	value := p.parseExpr()
	return &ast.OptionsDef{
		Name:  key,
		Value: value,
	}
}

func (p *Parser) tryParseOptions() *ast.Options {
	if !p.Token.IsKeywordLike("OPTIONS") {
		return nil
	}
	return p.parseOptions()
}

func (p *Parser) parseOptions() *ast.Options {
	pos := p.expectKeywordLike("OPTIONS").Pos
	p.expect("(")
	optionsDefs := parseCommaSeparatedList(p, p.parseOptionsDef)
	rparen := p.expect(")").Pos

	return &ast.Options{
		Options: pos,
		Rparen:  rparen,
		Records: optionsDefs,
	}
}

func (p *Parser) parseCreateVectorIndex(pos token.Pos) *ast.CreateVectorIndex {
	p.expectKeywordLike("VECTOR")
	p.expectKeywordLike("INDEX")
	ifNotExists := p.parseIfNotExists()
	name := p.parseIdent()
	p.expect("ON")
	tableName := p.parseIdent()
	p.expect("(")
	columnName := p.parseIdent()
	p.expect(")")

	storing := p.tryParseStoring()
	where := p.tryParseWhere()
	options := p.parseOptions()

	return &ast.CreateVectorIndex{
		Create:      pos,
		IfNotExists: ifNotExists,
		Name:        name,
		TableName:   tableName,
		ColumnName:  columnName,
		Storing:     storing,
		Where:       where,
		Options:     options,
	}
}

func (p *Parser) parseCreateSearchIndex(pos token.Pos) *ast.CreateSearchIndex {
	p.expectKeywordLike("SEARCH")
	p.expectKeywordLike("INDEX")
	name := p.parsePath()

	p.expect("ON")
	tableName := p.parsePath()

	p.expect("(")
	columnName := parseCommaSeparatedList(p, p.parseIdent)
	rparen := p.expect(")").Pos

	storing := p.tryParseStoring()

	var partitionColumns []*ast.Ident
	if p.Token.Kind == "PARTITION" {
		p.nextToken()
		p.expect("BY")
		partitionColumns = parseCommaSeparatedList(p, p.parseIdent)
	}

	orderBy := p.tryParseOrderBy()
	where := p.tryParseWhere()
	interleave := p.tryParseInterleaveIn()

	options := p.tryParseOptions()

	return &ast.CreateSearchIndex{
		Create:           pos,
		Name:             name,
		TableName:        tableName,
		TokenListPart:    columnName,
		Rparen:           rparen,
		Storing:          storing,
		PartitionColumns: partitionColumns,
		OrderBy:          orderBy,
		Where:            where,
		Interleave:       interleave,
		Options:          options,
	}
}

func (p *Parser) parseDropSearchIndex(pos token.Pos) *ast.DropSearchIndex {
	p.expectKeywordLike("SEARCH")
	p.expectKeywordLike("INDEX")
	ifExists := p.parseIfExists()
	name := p.parsePath()
	return &ast.DropSearchIndex{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}

}

func (p *Parser) parseAlterSearchIndex(pos token.Pos) *ast.AlterSearchIndex {
	p.expectKeywordLike("SEARCH")
	p.expectKeywordLike("INDEX")

	name := p.parsePath()
	alteration := p.parseIndexAlteration()

	return &ast.AlterSearchIndex{
		Alter:           pos,
		Name:            name,
		IndexAlteration: alteration,
	}
}

func (p *Parser) parseIndexAlteration() ast.IndexAlteration {
	switch {
	case p.Token.IsKeywordLike("ADD"):
		return p.parseAddStoredColumn()
	case p.Token.IsKeywordLike("DROP"):
		return p.parseDropStoredColumn()
	default:
		panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: ADD, DROP, but: %s", p.Token.AsString))
	}
}

func (p *Parser) parseVectorIndexAlteration() ast.VectorIndexAlteration {
	switch {
	case p.Token.IsKeywordLike("ADD"):
		return p.parseAddStoredColumn()
	case p.Token.IsKeywordLike("DROP"):
		return p.parseDropStoredColumn()
	default:
		panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: ADD, DROP, but: %s", p.Token.AsString))
	}
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

	ifNotExists := p.parseIfNotExists()

	name := p.parsePath()

	p.expect("ON")
	tableName := p.parsePath()

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
	options := p.tryParseOptions()

	return &ast.CreateIndex{
		Create:       pos,
		Rparen:       rparen,
		Unique:       unique,
		NullFiltered: nullFiltered,
		IfNotExists:  ifNotExists,
		Name:         name,
		TableName:    tableName,
		Keys:         keys,
		Storing:      storing,
		InterleaveIn: interleaveIn,
		Options:      options,
	}
}

func (p *Parser) parseCreateChangeStream(pos token.Pos) *ast.CreateChangeStream {
	p.expectKeywordLike("CHANGE")
	p.expectKeywordLike("STREAM")
	name := p.parseIdent()
	cs := &ast.CreateChangeStream{
		Create: pos,
		Name:   name,
	}
	if p.Token.Kind == "FOR" {
		cs.For = p.parseChangeStreamFor()
	}
	if p.Token.IsKeywordLike("OPTIONS") {
		cs.Options = p.parseOptions()
	}
	return cs

}

func (p *Parser) parseAlterChangeStream(pos token.Pos) *ast.AlterChangeStream {
	p.expectKeywordLike("CHANGE")
	p.expectKeywordLike("STREAM")
	name := p.parseIdent()
	cs := &ast.AlterChangeStream{
		Alter: pos,
		Name:  name,
	}
	if p.Token.Kind == "SET" {
		setpos := p.Token.Pos
		p.nextToken()
		if p.Token.Kind == "FOR" {
			cs.ChangeStreamAlteration = &ast.ChangeStreamSetFor{
				Set: setpos,
				For: p.parseChangeStreamFor(),
			}
			return cs
		} else if p.Token.IsKeywordLike("OPTIONS") {
			cs.ChangeStreamAlteration = &ast.ChangeStreamSetOptions{
				Set:     setpos,
				Options: p.parseOptions(),
			}
			return cs
		}
	} else if p.Token.IsKeywordLike("DROP") {
		droppos := p.Token.Pos
		p.nextToken()
		p.expect("FOR")
		allpos := p.expect("ALL").Pos
		cs.ChangeStreamAlteration = &ast.ChangeStreamDropForAll{
			Drop: droppos,
			All:  allpos,
		}
		return cs
	} else {
		p.panicfAtToken(&p.Token, "expected SET FOR or DROP FOR ALL or SET OPTIONS")
	}
	return cs
}

func (p *Parser) parseDropChangeStream(pos token.Pos) *ast.DropChangeStream {
	p.expectKeywordLike("CHANGE")
	p.expectKeywordLike("STREAM")
	name := p.parseIdent()
	return &ast.DropChangeStream{
		Drop: pos,
		Name: name,
	}
}

func (p *Parser) parseFunctionParam() *ast.FunctionParam {
	name := p.parseIdent()
	t := p.parseFunctionDataType()

	var defaultExpr ast.Expr
	if p.Token.Kind == "DEFAULT" {
		p.nextToken()
		defaultExpr = p.parseExpr()
	}

	return &ast.FunctionParam{
		Name:        name,
		Type:        t,
		DefaultExpr: defaultExpr,
	}
}

func (p *Parser) parseCreateFunction(pos token.Pos, orReplace bool) *ast.CreateFunction {
	p.expectKeywordLike("FUNCTION")
	name := p.parsePath()

	p.expect("(")
	var params []*ast.FunctionParam
	if p.Token.Kind != ")" {
		params = parseCommaSeparatedList(p, p.parseFunctionParam)
	}
	p.expect(")")

	var returnType ast.SchemaType
	if p.Token.IsKeywordLike("RETURNS") {
		p.nextToken()
		returnType = p.parseFunctionDataType()
	}

	determinism := p.tryParseDeterminism()

	var language string
	if p.Token.IsKeywordLike("LANGUAGE") {
		p.nextToken()
		language = strings.ToUpper(p.parseIdent().Name)
	}

	var remote bool
	if p.Token.IsKeywordLike("REMOTE") {
		p.nextToken()
		remote = true
	}

	var sqlSecurity ast.SecurityType
	if p.Token.IsKeywordLike("SQL") {
		sqlSecurity = p.parseSqlSecurity()
	}

	var options *ast.Options
	var body ast.Expr
	asPos := token.InvalidPos
	rparenAsPos := token.InvalidPos
	if p.Token.Kind == "AS" {
		asPos = p.expect("AS").Pos
		p.expect("(")
		body = p.parseExpr()
		rparenAsPos = p.expect(")").Pos
	} else if p.Token.IsKeywordLike("OPTIONS") {
		options = p.parseOptions()
	} else {
		p.panicfAtToken(&p.Token, "expected AS or OPTIONS")
	}

	return &ast.CreateFunction{
		Create:      pos,
		OrReplace:   orReplace,
		Name:        name,
		Params:      params,
		ReturnType:  returnType,
		SqlSecurity: sqlSecurity,
		Determinism: determinism,
		Language:    language,
		Remote:      remote,
		Options:     options,
		Definition:  body,
		As:          asPos,
		RparenAs:    rparenAsPos,
	}
}

func (p *Parser) parseDropFunction(pos token.Pos) *ast.DropFunction {
	p.expectKeywordLike("FUNCTION")

	ifExists := p.parseIfExists()
	name := p.parsePath()

	return &ast.DropFunction{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) tryParseDeterminism() ast.Determinism {
	if p.Token.IsKeywordLike("DETERMINISTIC") {
		p.nextToken()
		return ast.DeterminismDeterministic
	}

	if p.Token.Kind == "NOT" {
		lexer := p.Clone()
		p.nextToken()
		if p.Token.IsKeywordLike("DETERMINISTIC") {
			p.nextToken() // DETERMINISTIC
			return ast.DeterminismNotDeterministic
		}
		p.Lexer = lexer
	}

	return ""
}

func (p *Parser) parseSqlSecurity() ast.SecurityType {
	p.expectKeywordLike("SQL")
	p.expectKeywordLike("SECURITY")

	id := p.expect(token.TokenIdent)
	switch {
	case id.IsKeywordLike("INVOKER"):
		return ast.SecurityTypeInvoker
	case id.IsKeywordLike("DEFINER"):
		return ast.SecurityTypeDefiner
	default:
		p.panicfAtToken(id, "expected identifier: INVOKER, DEFINER, but: %s", id.Raw)
	}

	return ""
}

func (p *Parser) parseChangeStreamFor() ast.ChangeStreamFor {
	pos := p.expect("FOR").Pos
	if p.Token.Kind == "ALL" {
		p.nextToken()
		return &ast.ChangeStreamForAll{
			For: pos,
			All: p.Token.Pos,
		}

	}
	cswt := &ast.ChangeStreamForTables{
		For: pos,
	}
	for {
		tname := p.parseIdent()
		forTable := ast.ChangeStreamForTable{
			TableName: tname,
			Rparen:    token.InvalidPos,
		}

		if p.Token.Kind == "(" {
			p.nextToken()
			if p.Token.Kind != ")" {
				forTable.Columns = parseCommaSeparatedList(p, p.parseIdent)
			}
			forTable.Rparen = p.expect(")").Pos
		}

		cswt.Tables = append(cswt.Tables, &forTable)
		if p.Token.Kind == "," {
			p.nextToken()
			continue
		}
		break
	}
	return cswt

}

func (p *Parser) tryParseStoring() *ast.Storing {
	if !p.Token.IsKeywordLike("STORING") {
		return nil
	}
	pos := p.expectKeywordLike("STORING").Pos

	p.expect("(")
	columns := parseCommaSeparatedList(p, p.parseIdent)

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
	name := p.parsePath()

	var alteration ast.TableAlteration
	switch {
	case p.Token.IsKeywordLike("ADD"):
		alteration = p.parseAlterTableAdd()
	case p.Token.IsKeywordLike("DROP"):
		alteration = p.parseAlterTableDrop()
	case p.Token.IsKeywordLike("RENAME"):
		alteration = p.parseAlterTableRename()
	case p.Token.IsKeywordLike("REPLACE"):
		alteration = p.parseAlterTableReplace()
	case p.Token.Kind == "SET":
		alteration = p.parseAlterTableSet()
	case p.Token.IsKeywordLike("ALTER"):
		alteration = p.parseAlterColumn()
	default:
		if p.Token.Kind == token.TokenIdent {
			p.panicfAtToken(&p.Token, "expected pseuso keyword: ADD, ALTER, DROP, but: %s", p.Token.AsString)
		} else {
			p.panicfAtToken(&p.Token, "expected token: SET, <ident>, but: %s", p.Token.Kind)
		}
	}

	return &ast.AlterTable{
		Alter:           pos,
		Name:            name,
		TableAlteration: alteration,
	}
}

func (p *Parser) parseAddSynonym(add token.Pos) *ast.AddSynonym {
	p.expectKeywordLike("SYNONYM")
	name := p.parseIdent()

	return &ast.AddSynonym{
		Add:  add,
		Name: name,
	}
}

func (p *Parser) parseAlterTableAdd() ast.TableAlteration {
	pos := p.expectKeywordLike("ADD").Pos

	var alteration ast.TableAlteration

	switch {
	case p.Token.IsKeywordLike("SYNONYM"):
		alteration = p.parseAddSynonym(pos)
	case p.Token.IsKeywordLike("COLUMN"):
		p.expectKeywordLike("COLUMN")
		ifNotExists := p.parseIfNotExists()
		column := p.parseColumnDef()
		alteration = &ast.AddColumn{
			Add:         pos,
			IfNotExists: ifNotExists,
			Column:      column,
		}
	case p.Token.IsKeywordLike("CONSTRAINT"):
		alteration = &ast.AddTableConstraint{
			Add:             pos,
			TableConstraint: p.parseConstraint(),
		}
	case p.Token.IsKeywordLike("FOREIGN"):
		fk := p.parseForeignKey()
		alteration = &ast.AddTableConstraint{
			Add: pos,
			TableConstraint: &ast.TableConstraint{
				ConstraintPos: token.InvalidPos,
				Constraint:    fk,
			},
		}
	case p.Token.IsKeywordLike("CHECK"):
		c := p.parseCheck()
		alteration = &ast.AddTableConstraint{
			Add: pos,
			TableConstraint: &ast.TableConstraint{
				ConstraintPos: token.InvalidPos,
				Constraint:    c,
			},
		}
	case p.Token.IsKeywordLike("ROW"):
		rdp := p.parseRowDeletionPolicy()

		alteration = &ast.AddRowDeletionPolicy{
			Add:               pos,
			RowDeletionPolicy: rdp,
		}
	default:
		p.panicfAtToken(&p.Token, "expected pseuso keyword: COLUMN, CONSTRAINT, FOREIGN, but: %s", p.Token.AsString)
	}

	return alteration
}

func (p *Parser) parseAlterTableDrop() ast.TableAlteration {
	pos := p.expectKeywordLike("DROP").Pos

	var alteration ast.TableAlteration

	switch {
	case p.Token.IsKeywordLike("SYNONYM"):
		p.expectKeywordLike("SYNONYM")
		name := p.parseIdent()
		alteration = &ast.DropSynonym{
			Drop: pos,
			Name: name,
		}
	case p.Token.IsKeywordLike("COLUMN"):
		p.expectKeywordLike("COLUMN")
		name := p.parseIdent()
		alteration = &ast.DropColumn{
			Drop: pos,
			Name: name,
		}
	case p.Token.IsKeywordLike("CONSTRAINT"):
		p.expectKeywordLike("CONSTRAINT")
		name := p.parseIdent()
		alteration = &ast.DropConstraint{
			Drop: pos,
			Name: name,
		}
	case p.Token.IsKeywordLike("ROW"):
		p.expectKeywordLike("ROW")
		p.expectKeywordLike("DELETION")
		policyPos := p.expectKeywordLike("POLICY").Pos
		alteration = &ast.DropRowDeletionPolicy{
			Drop:   pos,
			Policy: policyPos,
		}
	default:
		p.panicfAtToken(&p.Token, "expected pseuso keyword: COLUMN, CONSTRAINT, but: %s", p.Token.AsString)
	}

	return alteration
}

func (p *Parser) parseAlterTableRename() ast.TableAlteration {
	pos := p.expectKeywordLike("RENAME").Pos
	p.expect("TO")
	name := p.parseIdent()

	var addSynonym *ast.AddSynonym
	if p.Token.Kind == "," {
		p.nextToken()
		add := p.expectKeywordLike("ADD").Pos
		addSynonym = p.parseAddSynonym(add)
	}

	return &ast.RenameTo{
		Rename:     pos,
		Name:       name,
		AddSynonym: addSynonym,
	}
}

func (p *Parser) parseAlterTableReplace() ast.TableAlteration {
	pos := p.expectKeywordLike("REPLACE").Pos
	rdp := p.parseRowDeletionPolicy()

	return &ast.ReplaceRowDeletionPolicy{
		Replace:           pos,
		RowDeletionPolicy: rdp,
	}
}

func (p *Parser) parseAlterTableSet() ast.TableAlteration {
	pos := p.expect("SET").Pos
	switch {
	case p.Token.Kind == "ON":
		onDelete, onDeleteEnd := p.parseOnDeleteAction()

		return &ast.SetOnDelete{
			Set:         pos,
			OnDeleteEnd: onDeleteEnd,
			OnDelete:    onDelete,
		}
	case p.Token.IsKeywordLike("OPTIONS"):
		options := p.parseOptions()
		return &ast.AlterTableSetOptions{
			Set:     pos,
			Options: options,
		}
	case p.Token.IsKeywordLike("INTERLEAVE"):
		p.nextToken()
		p.expect("IN")

		var enforced bool
		if p.Token.IsKeywordLike("PARENT") {
			p.nextToken()
			enforced = true
		}

		tableName := p.parsePath()
		onDelete, onDeleteEnd := p.tryParseOnDeleteAction()

		return &ast.SetInterleaveIn{
			Set:         pos,
			Enforced:    enforced,
			TableName:   tableName,
			OnDeleteEnd: onDeleteEnd,
			OnDelete:    onDelete,
		}
	default:
		panic(p.errorfAtToken(&p.Token, "expected token: ON, OPTIONS, but: %s", p.Token.AsString))
	}
}

func (p *Parser) parseColumnAlteration() ast.ColumnAlteration {
	switch {
	case p.Token.Kind == "SET":
		set := p.expect("SET").Pos
		if p.Token.Kind == "DEFAULT" {
			defaultExpr := p.parseColumnDefaultExpr()
			return &ast.AlterColumnSetDefault{
				Set:         set,
				DefaultExpr: defaultExpr,
			}
		} else {
			options := p.parseOptions()
			return &ast.AlterColumnSetOptions{
				Set:     set,
				Options: options,
			}
		}
	case p.Token.IsKeywordLike("DROP"):
		drop := p.expectKeywordLike("DROP").Pos
		def := p.expect("DEFAULT").Pos
		return &ast.AlterColumnDropDefault{
			Drop:    drop,
			Default: def,
		}
	case p.Token.IsKeywordLike("ALTER"):
		alter := p.expectKeywordLike("ALTER").Pos
		p.expectKeywordLike("IDENTITY")
		alteration := p.parseIdentityAlteration()
		return &ast.AlterColumnAlterIdentity{
			Alter:      alter,
			Alteration: alteration,
		}
	default:
		t, notNull, null := p.parseTypeNotNull()
		defaultExpr := p.tryParseColumnDefaultExpr()
		return &ast.AlterColumnType{
			Type:        t,
			Null:        null,
			NotNull:     notNull,
			DefaultExpr: defaultExpr,
		}
	}

}

func (p *Parser) parseSkipRange() *ast.SkipRange {
	pos := p.expectKeywordLike("SKIP").Pos
	p.expect("RANGE")
	min := p.parseIntLiteral()
	p.expect(",")
	max := p.parseIntLiteral()
	return &ast.SkipRange{
		Skip: pos,
		Min:  min,
		Max:  max,
	}
}

func (p *Parser) parseNoSkipRange() *ast.NoSkipRange {
	no := p.expect("NO").Pos
	p.expectKeywordLike("SKIP")
	rangePos := p.expect("RANGE").Pos

	return &ast.NoSkipRange{
		No:    no,
		Range: rangePos,
	}

}
func (p *Parser) parseIdentityAlteration() ast.IdentityAlteration {
	pos := p.Token.Pos

	switch {
	case p.Token.Kind == "SET":
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("SKIP"):
			skipRange := p.parseSkipRange()
			return &ast.SetSkipRange{
				Set:       pos,
				SkipRange: skipRange,
			}
		case p.Token.Kind == "NO":
			noSkipRange := p.parseNoSkipRange()
			return &ast.SetNoSkipRange{
				Set:         pos,
				NoSkipRange: noSkipRange,
			}
		default:
			panic(p.errorfAtToken(&p.Token, "expect SKIP or NO, but: %s", p.Token.AsString))
		}
	case p.Token.IsKeywordLike("RESTART"):
		return p.parseRestartCounterWith()
	default:
		panic(p.errorfAtToken(&p.Token, "expected token: SET, RESTART, but: %s", p.Token.AsString))
	}
}
func (p *Parser) parseAlterColumn() ast.TableAlteration {
	pos := p.expectKeywordLike("ALTER").Pos
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	alteration := p.parseColumnAlteration()
	return &ast.AlterColumn{
		Alter:      pos,
		Name:       name,
		Alteration: alteration,
	}
}

func (p *Parser) parseAlterIndex(pos token.Pos) *ast.AlterIndex {
	p.expectKeywordLike("INDEX")

	name := p.parsePath()
	alteration := p.parseIndexAlteration()

	return &ast.AlterIndex{
		Alter:           pos,
		Name:            name,
		IndexAlteration: alteration,
	}
}

func (p *Parser) parseRestartCounterWith() *ast.RestartCounterWith {
	restart := p.expectKeywordLike("RESTART").Pos
	p.expectKeywordLike("COUNTER")
	p.expect("WITH")

	counter := p.parseIntLiteral()

	return &ast.RestartCounterWith{
		Restart: restart,
		Counter: counter,
	}
}

func (p *Parser) tryParseSequenceParam() ast.SequenceParam {
	if !p.Token.IsKeywordLike("BIT_REVERSED_POSITIVE") && !p.Token.IsKeywordLike("SKIP") && !p.Token.IsKeywordLike("START") {
		return nil
	}

	return p.parseSequenceParam()
}

func (p *Parser) parseSequenceParam() ast.SequenceParam {
	pos := p.Token.Pos

	switch {
	case p.Token.IsKeywordLike("SKIP"):
		p.nextToken()
		p.expect("RANGE")

		min := p.parseIntLiteral()
		p.expect(",")
		max := p.parseIntLiteral()

		return &ast.SkipRange{
			Skip: pos,
			Min:  min,
			Max:  max,
		}
	case p.Token.IsKeywordLike("BIT_REVERSED_POSITIVE"):
		p.nextToken()

		return &ast.BitReversedPositive{BitReversedPositive: pos}
	case p.Token.IsKeywordLike("START"):
		p.nextToken()
		p.expectKeywordLike("COUNTER")
		p.expect("WITH")

		counter := p.parseIntLiteral()

		return &ast.StartCounterWith{
			Start:   pos,
			Counter: counter,
		}
	default:
		panic(p.errorfAtToken(&p.Token, `expect "BIT_REVERSED_POSITIVE", "SKIP", "START", but %q`, p.Token.AsString))
	}
}

func (p *Parser) parseAlterSequence(pos token.Pos) *ast.AlterSequence {
	p.expectKeywordLike("SEQUENCE")
	name := p.parsePath()

	var options *ast.Options
	if p.Token.Kind == "SET" {
		p.nextToken()
		options = p.parseOptions()
	}

	var skipRange *ast.SkipRange
	if p.Token.IsKeywordLike("SKIP") {
		skipRange = p.parseSkipRange()
	}

	var noSkipRange *ast.NoSkipRange
	if p.Token.Kind == "NO" {
		noSkipRange = p.parseNoSkipRange()
	}

	var restartCounterWith *ast.RestartCounterWith
	if p.Token.IsKeywordLike("RESTART") {
		restartCounterWith = p.parseRestartCounterWith()
	}

	return &ast.AlterSequence{
		Alter:              pos,
		Name:               name,
		Options:            options,
		RestartCounterWith: restartCounterWith,
		SkipRange:          skipRange,
		NoSkipRange:        noSkipRange,
	}
}

func (p *Parser) parseAddStoredColumn() *ast.AddStoredColumn {
	pos := p.expectKeywordLike("ADD").Pos
	p.expectKeywordLike("STORED")
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	return &ast.AddStoredColumn{
		Add:  pos,
		Name: name,
	}
}

func (p *Parser) parseDropStoredColumn() *ast.DropStoredColumn {
	pos := p.expectKeywordLike("DROP").Pos
	p.expectKeywordLike("STORED")
	p.expectKeywordLike("COLUMN")

	name := p.parseIdent()

	return &ast.DropStoredColumn{
		Drop: pos,
		Name: name,
	}
}
func (p *Parser) parseDropTable(pos token.Pos) *ast.DropTable {
	p.expectKeywordLike("TABLE")
	ifExists := p.parseIfExists()
	name := p.parsePath()
	return &ast.DropTable{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) parseDropIndex(pos token.Pos) *ast.DropIndex {
	p.expectKeywordLike("INDEX")
	ifExists := p.parseIfExists()
	name := p.parsePath()
	return &ast.DropIndex{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) parseAlterVectorIndex(pos token.Pos) *ast.AlterVectorIndex {
	p.expectKeywordLike("VECTOR")
	p.expectKeywordLike("INDEX")
	name := p.parsePath()
	alteration := p.parseVectorIndexAlteration()
	return &ast.AlterVectorIndex{
		Alter:      pos,
		Name:       name,
		Alteration: alteration,
	}
}

func (p *Parser) parseDropVectorIndex(pos token.Pos) *ast.DropVectorIndex {
	p.expectKeywordLike("VECTOR")
	p.expectKeywordLike("INDEX")
	ifExists := p.parseIfExists()
	name := p.parseIdent()
	return &ast.DropVectorIndex{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) parseDropSequence(pos token.Pos) *ast.DropSequence {
	p.expectKeywordLike("SEQUENCE")
	ifExists := p.parseIfExists()
	name := p.parsePath()
	return &ast.DropSequence{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

func (p *Parser) parseCreateRole(pos token.Pos) *ast.CreateRole {
	p.expectKeywordLike("ROLE")
	name := p.parseIdent()
	return &ast.CreateRole{
		Create: pos,
		Name:   name,
	}

}

func (p *Parser) parseDropRole(pos token.Pos) *ast.DropRole {
	p.expectKeywordLike("ROLE")
	name := p.parseIdent()
	return &ast.DropRole{
		Drop: pos,
		Name: name,
	}

}

func (p *Parser) parseGrant(pos token.Pos) *ast.Grant {
	privilege := p.parsePrivilege()
	p.expect("TO")
	p.expectKeywordLike("ROLE")
	roles := parseCommaSeparatedList(p, p.parseIdent)
	return &ast.Grant{
		Grant:     pos,
		Privilege: privilege,
		Roles:     roles,
	}
}

func (p *Parser) parseRevoke(pos token.Pos) *ast.Revoke {
	privilege := p.parsePrivilege()
	p.expect("FROM")
	p.expectKeywordLike("ROLE")
	roles := parseCommaSeparatedList(p, p.parseIdent)
	return &ast.Revoke{
		Revoke:    pos,
		Privilege: privilege,
		Roles:     roles,
	}
}

func (p *Parser) parsePrivilege() ast.Privilege {
	if s := p.tryParseSelectPrivilegeOnView(); s != nil {
		return s
	}
	if e := p.tryParseExecutePrivilegeOnTableFunction(); e != nil {
		return e
	}
	if r := p.tryParseRolePrivilege(); r != nil {
		return r
	}
	if c := p.tryParseSelectPrivilegeOnChangeStream(); c != nil {
		return c
	}
	return p.parsePrivilegeOnTable()
}

func (p *Parser) tryParseSelectPrivilegeOnView() *ast.SelectPrivilegeOnView {
	if p.Token.Kind != "SELECT" {
		return nil
	}
	lexer := p.Clone()
	pos := p.expect("SELECT").Pos
	if p.Token.Kind != "ON" {
		p.Lexer = lexer
		return nil
	}
	p.expect("ON")
	if !p.Token.IsKeywordLike("VIEW") {
		p.Lexer = lexer
		return nil
	}
	p.expectKeywordLike("VIEW")
	names := parseCommaSeparatedList(p, p.parsePath)
	return &ast.SelectPrivilegeOnView{
		Select: pos,
		Names:  names,
	}
}

func (p *Parser) tryParseExecutePrivilegeOnTableFunction() *ast.ExecutePrivilegeOnTableFunction {
	if !p.Token.IsKeywordLike("EXECUTE") {
		return nil
	}
	pos := p.expectKeywordLike("EXECUTE").Pos
	p.expect("ON")
	p.expectKeywordLike("TABLE")
	p.expectKeywordLike("FUNCTION")
	names := parseCommaSeparatedList(p, p.parsePath)
	return &ast.ExecutePrivilegeOnTableFunction{
		Execute: pos,
		Names:   names,
	}
}

func (p *Parser) tryParseRolePrivilege() *ast.RolePrivilege {
	if !p.Token.IsKeywordLike("ROLE") {
		return nil
	}
	pos := p.expectKeywordLike("ROLE").Pos
	names := parseCommaSeparatedList(p, p.parseIdent)
	return &ast.RolePrivilege{
		Role:  pos,
		Names: names,
	}
}

func (p *Parser) tryParseSelectPrivilegeOnChangeStream() *ast.SelectPrivilegeOnChangeStream {
	if p.Token.Kind != "SELECT" {
		return nil
	}
	lexer := p.Clone()
	pos := p.expect("SELECT").Pos
	if p.Token.Kind != "ON" {
		p.Lexer = lexer
		return nil
	}
	p.expect("ON")
	if !p.Token.IsKeywordLike("CHANGE") {
		p.Lexer = lexer
		return nil
	}
	p.expectKeywordLike("CHANGE")
	p.expectKeywordLike("STREAM")
	names := parseCommaSeparatedList(p, p.parsePath)

	return &ast.SelectPrivilegeOnChangeStream{
		Select: pos,
		Names:  names,
	}
}

func (p *Parser) parsePrivilegeOnTable() *ast.PrivilegeOnTable {
	privileges := parseCommaSeparatedList(p, p.parseTablePrivilege)
	p.expect("ON")
	p.expectKeywordLike("TABLE")
	names := parseCommaSeparatedList(p, p.parsePath)
	return &ast.PrivilegeOnTable{
		Privileges: privileges,
		Names:      names,
	}
}

func (p *Parser) parseTablePrivilege() ast.TablePrivilege {
	pos := p.Token.Pos
	switch {
	case p.Token.Kind == "SELECT":
		p.nextToken()
		columns, rparen := p.tryParseTablePrivilegeColumns()
		return &ast.SelectPrivilege{
			Select:  pos,
			Rparen:  rparen,
			Columns: columns,
		}
	case p.Token.IsKeywordLike("INSERT"):
		p.nextToken()
		columns, rparen := p.tryParseTablePrivilegeColumns()
		return &ast.InsertPrivilege{
			Insert:  pos,
			Rparen:  rparen,
			Columns: columns,
		}
	case p.Token.IsKeywordLike("UPDATE"):
		p.nextToken()
		columns, rparen := p.tryParseTablePrivilegeColumns()
		return &ast.UpdatePrivilege{
			Update:  pos,
			Rparen:  rparen,
			Columns: columns,
		}
	case p.Token.IsKeywordLike("DELETE"):
		p.nextToken()
		return &ast.DeletePrivilege{
			Delete: pos,
		}
	}
	if p.Token.Kind != token.TokenIdent {
		panic(p.errorfAtToken(&p.Token, "expected pseudo keyword: INSERT, UPDATE, DELETE, but: %s", p.Token.AsString))
	} else {
		panic(p.errorfAtToken(&p.Token, "expected token: SELECT, <ident>, but: %s", p.Token.Kind))
	}
}

func (p *Parser) tryParseTablePrivilegeColumns() ([]*ast.Ident, token.Pos) {
	if p.Token.Kind != "(" {
		return nil, token.InvalidPos
	}
	p.nextToken()
	columns := parseCommaSeparatedList(p, p.parseIdent)
	rparen := p.expect(")").Pos
	return columns, rparen
}

// begin CREATE PROPERTY GRAPH

func (p *Parser) parseCreatePropertyGraph(pos token.Pos, orReplace bool) *ast.CreatePropertyGraph {
	p.expectKeywordLike("PROPERTY")
	p.expectKeywordLike("GRAPH")

	ifNotExists := p.parseIfNotExists()
	name := p.parseIdent()
	content := p.parsePropertyGraphContent()

	return &ast.CreatePropertyGraph{
		Create:      pos,
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
		Content:     content,
	}
}

func (p *Parser) parsePropertyGraphContent() *ast.PropertyGraphContent {
	node := p.expectKeywordLike("NODE").Pos
	p.expectKeywordLike("TABLES")

	nodeTables := &ast.PropertyGraphNodeTables{
		Node:   node,
		Tables: p.parsePropertyGraphElementList(),
	}

	var edgeTables *ast.PropertyGraphEdgeTables
	if p.Token.IsKeywordLike("EDGE") {
		edge := p.expectKeywordLike("EDGE").Pos
		p.expectKeywordLike("TABLES")
		tables := p.parsePropertyGraphElementList()
		edgeTables = &ast.PropertyGraphEdgeTables{
			Edge:   edge,
			Tables: tables,
		}
	}

	return &ast.PropertyGraphContent{
		NodeTables: nodeTables,
		EdgeTables: edgeTables,
	}
}

func (p *Parser) parsePropertyGraphElementList() *ast.PropertyGraphElementList {
	lparen := p.expect("(").Pos
	elements := parseCommaSeparatedList(p, p.parsePropertyGraphElement)
	rparen := p.expect(")").Pos

	return &ast.PropertyGraphElementList{
		Lparen:   lparen,
		Rparen:   rparen,
		Elements: elements,
	}

}

func (p *Parser) parsePropertyGraphElement() *ast.PropertyGraphElement {
	name := p.parseIdent()

	var alias *ast.Ident
	if p.Token.Kind == "AS" {
		p.nextToken()
		alias = p.parseIdent()
	}

	keys := p.tryParsePropertyGraphElementKeys()
	properties := p.tryParsePropertyGraphLabelsOrProperties()
	dynamicLabel := p.tryParsePropertyGraphDynamicLabel()
	dynamicProperties := p.tryParsePropertyGraphDynamicProperties()

	return &ast.PropertyGraphElement{
		Name:              name,
		Alias:             alias,
		Keys:              keys,
		Properties:        properties,
		DynamicLabel:      dynamicLabel,
		DynamicProperties: dynamicProperties,
	}
}

// parsePropertyGraphLabelAndPropertiesList parses consecutive ast.PropertyGraphLabelAndProperties,
// and returns *ast.PropertyGraphLabelAndPropertiesList.
func (p *Parser) parsePropertyGraphLabelAndPropertiesList() *ast.PropertyGraphLabelAndPropertiesList {
	// list can be empty
	var list []*ast.PropertyGraphLabelAndProperties
	for p.Token.Kind == "DEFAULT" || p.Token.IsKeywordLike("LABEL") {
		elemLabel := p.parsePropertyGraphElementLabel()
		properties := p.tryParsePropertyGraphElementProperties()

		list = append(list, &ast.PropertyGraphLabelAndProperties{
			Label:      elemLabel,
			Properties: properties,
		})
	}

	return &ast.PropertyGraphLabelAndPropertiesList{
		LabelAndProperties: list,
	}
}

func (p *Parser) parsePropertyGraphElementLabel() ast.PropertyGraphElementLabel {
	if p.Token.Kind == "DEFAULT" {
		def := p.expect("DEFAULT").Pos
		label := p.expectKeywordLike("LABEL").Pos

		return &ast.PropertyGraphElementLabelDefaultLabel{
			Default: def,
			Label:   label,
		}
	}

	label := p.expectKeywordLike("LABEL").Pos
	name := p.parseIdent()

	return &ast.PropertyGraphElementLabelLabelName{
		Label: label,
		Name:  name,
	}
}

func (p *Parser) tryParsePropertyGraphElementProperties() ast.PropertyGraphElementProperties {
	if p.Token.Kind != "NO" && !p.Token.IsKeywordLike("PROPERTIES") {
		return nil
	}

	return p.parsePropertyGraphElementProperties()
}

func (p *Parser) parsePropertyGraphElementProperties() ast.PropertyGraphElementProperties {
	if p.Token.Kind != "NO" && !p.Token.IsKeywordLike("PROPERTIES") {
		p.panicfAtToken(&p.Token, `expect "NO" or "PROPERTIES", but %v`, p.Token.Kind)
	}

	if p.Token.Kind == "NO" {
		no := p.expect("NO").Pos
		properties := p.expectKeywordLike("PROPERTIES").Pos

		return &ast.PropertyGraphNoProperties{
			No:         no,
			Properties: properties,
		}
	}

	properties := p.expectKeywordLike("PROPERTIES")
	if p.Token.IsKeywordLike("ARE") || p.Token.Kind == "ALL" {
		// "ARE" is optional
		if p.Token.IsKeywordLike("ARE") {
			p.nextToken()
		}

		p.expect("ALL")
		columns := p.expectKeywordLike("COLUMNS").Pos

		var exceptColumns *ast.PropertyGraphColumnNameList
		if p.Token.Kind == "EXCEPT" {
			p.nextToken()
			exceptColumns = p.parsePropertyGraphColumnNameList()
		}

		return &ast.PropertyGraphPropertiesAre{
			Properties:    properties.Pos,
			Columns:       columns,
			ExceptColumns: exceptColumns,
		}
	}

	p.expect("(")
	list := parseCommaSeparatedList(p, p.parsePropertyGraphDerivedProperty)
	rparen := p.expect(")").Pos

	return &ast.PropertyGraphDerivedPropertyList{
		Rparen:            rparen,
		Properties:        properties.Pos,
		DerivedProperties: list,
	}
}

func (p *Parser) parsePropertyGraphDerivedProperty() *ast.PropertyGraphDerivedProperty {
	expr := p.parseExpr()

	var name *ast.Ident
	if p.Token.Kind == "AS" {
		p.nextToken()
		name = p.parseIdent()
	}

	return &ast.PropertyGraphDerivedProperty{
		Expr:  expr,
		Alias: name,
	}
}

func (p *Parser) tryParsePropertyGraphLabelsOrProperties() ast.PropertyGraphLabelsOrProperties {
	if p.Token.IsKeywordLike("LABEL") || p.Token.Kind == "DEFAULT" {
		return p.parsePropertyGraphLabelAndPropertiesList()
	}

	if properties := p.tryParsePropertyGraphElementProperties(); properties != nil {
		return &ast.PropertyGraphSingleProperties{
			Properties: properties,
		}
	}
	return nil
}

func (p *Parser) tryParsePropertyGraphElementKeys() ast.PropertyGraphElementKeys {
	if !p.Token.IsKeywordLike("KEY") && !p.Token.IsKeywordLike("SOURCE") {
		return nil
	}

	// element_key
	var elementKey *ast.PropertyGraphElementKey
	if p.Token.IsKeywordLike("KEY") {
		key := p.expectKeywordLike("KEY").Pos
		keyColumns := p.parsePropertyGraphColumnNameList()

		elementKey = &ast.PropertyGraphElementKey{
			Key:  key,
			Keys: keyColumns,
		}

		// if SOURCE KEY doesn't follow, it is node_element_key.
		if !p.Token.IsKeywordLike("SOURCE") {
			return &ast.PropertyGraphNodeElementKey{
				Key: elementKey,
			}
		}

	}

	// the rest of edge_element_keys

	// source_key
	source := p.expectKeywordLike("SOURCE").Pos
	p.expectKeywordLike("KEY")
	sourceColumns := p.parsePropertyGraphColumnNameList()
	p.expectKeywordLike("REFERENCES")
	sourceReference := p.parseIdent()
	sourceReferenceColumns := p.tryParsePropertyGraphColumnNameList()

	// destination_key
	destination := p.expectKeywordLike("DESTINATION").Pos
	p.expectKeywordLike("KEY")
	destinationColumns := p.parsePropertyGraphColumnNameList()
	p.expectKeywordLike("REFERENCES")
	destinationReference := p.parseIdent()
	destinationReferenceColumns := p.tryParsePropertyGraphColumnNameList()

	return &ast.PropertyGraphEdgeElementKeys{
		Element: elementKey,
		Source: &ast.PropertyGraphSourceKey{
			Source:           source,
			Keys:             sourceColumns,
			ElementReference: sourceReference,
			ReferenceColumns: sourceReferenceColumns,
		},
		Destination: &ast.PropertyGraphDestinationKey{
			Destination:      destination,
			Keys:             destinationColumns,
			ElementReference: destinationReference,
			ReferenceColumns: destinationReferenceColumns,
		},
	}
}

func (p *Parser) parsePropertyGraphColumnNameList() *ast.PropertyGraphColumnNameList {
	lparen := p.expect("(").Pos
	list := parseCommaSeparatedList(p, p.parseIdent)
	rparen := p.expect(")").Pos

	return &ast.PropertyGraphColumnNameList{
		Lparen:         lparen,
		Rparen:         rparen,
		ColumnNameList: list,
	}
}

func (p *Parser) tryParsePropertyGraphColumnNameList() *ast.PropertyGraphColumnNameList {
	if p.Token.Kind != "(" {
		return nil
	}
	return p.parsePropertyGraphColumnNameList()
}

func (p *Parser) lookaheadPropertyGraphDynamicLabel() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if !p.Token.IsKeywordLike("DYNAMIC") {
		return false
	}
	p.nextToken()

	return p.Token.IsKeywordLike("LABEL")
}

func (p *Parser) tryParsePropertyGraphDynamicLabel() *ast.PropertyGraphDynamicLabel {
	if !p.lookaheadPropertyGraphDynamicLabel() {
		return nil
	}

	dynamic := p.expectKeywordLike("DYNAMIC").Pos
	p.expectKeywordLike("LABEL")
	p.expect("(")
	name := p.parseIdent()
	rparen := p.expect(")").Pos

	return &ast.PropertyGraphDynamicLabel{
		Dynamic:    dynamic,
		Rparen:     rparen,
		ColumnName: name,
	}
}

func (p *Parser) lookaheadPropertyGraphDynamicProperties() bool {
	lexer := p.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if !p.Token.IsKeywordLike("DYNAMIC") {
		return false
	}
	p.nextToken()

	return p.Token.IsKeywordLike("PROPERTIES")
}

func (p *Parser) tryParsePropertyGraphDynamicProperties() *ast.PropertyGraphDynamicProperties {
	if !p.lookaheadPropertyGraphDynamicProperties() {
		return nil
	}

	dynamic := p.expectKeywordLike("DYNAMIC").Pos
	p.expectKeywordLike("PROPERTIES")
	p.expect("(")
	name := p.parseIdent()
	rparen := p.expect(")").Pos

	return &ast.PropertyGraphDynamicProperties{
		Dynamic:    dynamic,
		Rparen:     rparen,
		ColumnName: name,
	}
}

// end CREATE PROPERTY GRAPH

func (p *Parser) parseDropPropertyGraph(pos token.Pos) *ast.DropPropertyGraph {
	p.expectKeywordLike("PROPERTY")
	p.expectKeywordLike("GRAPH")
	ifExists := p.parseIfExists()
	name := p.parseIdent()

	return &ast.DropPropertyGraph{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
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
		end := p.expect(">").Pos

		var namedArgs []*ast.NamedArg
		rparen := token.InvalidPos
		if p.Token.Kind == "(" {
			p.nextToken()
			namedArgs = parseCommaSeparatedList(p, p.parseNamedArg)
			rparen = p.expect(")").Pos
		}

		return &ast.ArraySchemaType{
			Array:     pos,
			Gt:        end,
			Item:      t,
			NamedArgs: namedArgs,
			Rparen:    rparen,
		}
	case "STRUCT":
		return p.parseStructType()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ARRAY, <ident>, but: %s", p.Token.Kind))
}

func (p *Parser) parseAlterStatistics(pos token.Pos) *ast.AlterStatistics {
	p.expectKeywordLike("STATISTICS")
	name := p.parseIdent()
	p.expect("SET")
	options := p.parseOptions()

	return &ast.AlterStatistics{
		Alter:   pos,
		Name:    name,
		Options: options,
	}
}

func (p *Parser) parseAnalyze() *ast.Analyze {
	pos := p.expectKeywordLike("ANALYZE").Pos

	return &ast.Analyze{
		Analyze: pos,
	}
}

func (p *Parser) tryParseCreateModelColumn() *ast.CreateModelColumn {
	name := p.parseIdent()
	dataType := p.parseSchemaType()
	options := p.tryParseOptions()

	return &ast.CreateModelColumn{
		Name:     name,
		DataType: dataType,
		Options:  options,
	}
}

func (p *Parser) parseModelColumns(keyword string) (token.Pos, token.Pos, []*ast.CreateModelColumn) {
	keywordPos := p.expectKeywordLike(keyword).Pos
	p.expect("(")
	columns := []*ast.CreateModelColumn{p.tryParseCreateModelColumn()}

	for p.Token.Kind == "," {
		p.nextToken()
		if p.Token.Kind == ")" {
			break // allow trailing comma
		}
		columns = append(columns, p.tryParseCreateModelColumn())
	}

	rparen := p.expect(")").Pos

	return keywordPos, rparen, columns
}

func (p *Parser) tryParseCreateModelInputOutput() *ast.CreateModelInputOutput {
	if !p.Token.IsKeywordLike("INPUT") {
		return nil
	}

	pos, _, inputColumns := p.parseModelColumns("INPUT")
	_, rparen, outputColumns := p.parseModelColumns("OUTPUT")

	return &ast.CreateModelInputOutput{
		Input:         pos,
		Rparen:        rparen,
		InputColumns:  inputColumns,
		OutputColumns: outputColumns,
	}
}

func (p *Parser) parseCreateModel(pos token.Pos, orReplace bool) *ast.CreateModel {
	p.expectKeywordLike("MODEL")
	name := p.parseIdent()
	ifNotExists := p.parseIfNotExists()
	inputOutput := p.tryParseCreateModelInputOutput()
	remote := p.expectKeywordLike("REMOTE").Pos
	options := p.tryParseOptions()

	return &ast.CreateModel{
		Create:      pos,
		OrReplace:   orReplace,
		IfNotExists: ifNotExists,
		Name:        name,
		InputOutput: inputOutput,
		Remote:      remote,
		Options:     options,
	}

}

func (p *Parser) parseAlterModel(pos token.Pos) *ast.AlterModel {
	p.expectKeywordLike("MODEL")
	ifExists := p.parseIfExists()
	name := p.parseIdent()
	p.expect("SET")
	options := p.parseOptions()

	return &ast.AlterModel{
		Alter:    pos,
		IfExists: ifExists,
		Name:     name,
		Options:  options,
	}
}

func (p *Parser) parseDropModel(pos token.Pos) *ast.DropModel {
	p.expectKeywordLike("MODEL")
	ifExists := p.parseIfExists()
	name := p.parseIdent()

	return &ast.DropModel{
		Drop:     pos,
		IfExists: ifExists,
		Name:     name,
	}
}

// Usually, you don't need to update scalarSchemaTypes because NamedType can handle types.
// It is maintained for compatibility.
var scalarSchemaTypes = []string{
	"BOOL",
	"INT64",
	"FLOAT32",
	"FLOAT64",
	"DATE",
	"TIMESTAMP",
	"NUMERIC",
	"JSON",
	"TOKENLIST",
}

var sizedSchemaTypes = []string{
	"STRING",
	"BYTES",
}

// parseScalarSchemaType parses a scalar type for column definition.
// - STRING and BYTES: Size specification is mandatory. Always parsed as SizedSchemaType.
// - Other types: Size specification is denied. Parsed as ScalarSchemaType or NamedType.
func (p *Parser) parseScalarSchemaType() ast.SchemaType {
	if !p.lookaheadSimpleType() {
		return p.parseNamedType()
	}

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

	panic(p.errorfAtToken(id, "expect ident: %s, %s, but: %s", strings.Join(scalarSchemaTypes, ", "), strings.Join(sizedSchemaTypes, ", "), id.AsString))
}

func (p *Parser) parseFunctionDataType() ast.SchemaType {
	switch p.Token.Kind {
	case token.TokenIdent, "INTERVAL":
		return p.parseScalarFunctionSchemaType()
	case "ARRAY":
		return p.parseArrayFunctionSchemaType()
	case "STRUCT":
		return p.parseStructType()
	}

	panic(p.errorfAtToken(&p.Token, "expected token: ARRAY, STRUCT, INTERVAL, <ident>, but: %s", p.Token.Kind))
}

// parseScalarFunctionSchemaType parses a scalar type for function parameter or return type.
// - STRING and BYTES: Size specification is optional.
//   - If size is specified, parsed as SizedSchemaType.
//   - If size is omitted, parsed as ScalarSchemaType.
//
// - Other types: Size specification is denied. Parsed as ScalarSchemaType.
func (p *Parser) parseScalarFunctionSchemaType() ast.SchemaType {
	if !p.lookaheadSimpleType() {
		return p.parseNamedType()
	}

	if p.Token.Kind == "INTERVAL" {
		pos := p.expect("INTERVAL").Pos
		return &ast.ScalarSchemaType{
			NamePos: pos,
			Name:    ast.IntervalTypeName,
		}
	}

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
			if p.Token.Kind == "(" {
				p.nextToken()
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
			} else {
				return &ast.ScalarSchemaType{
					NamePos: pos,
					Name:    ast.ScalarTypeName(name),
				}
			}
		}
	}

	panic(p.errorfAtToken(id, "expect ident: %s, %s, but: %s", strings.Join(scalarSchemaTypes, ", "), strings.Join(sizedSchemaTypes, ", "), id.AsString))
}

func (p *Parser) parseArrayFunctionSchemaType() *ast.ArraySchemaType {
	pos := p.expect("ARRAY").Pos
	p.expect("<")
	t := p.parseFunctionDataType()
	end := p.expect(">").Pos

	var namedArgs []*ast.NamedArg
	rparen := token.InvalidPos
	if p.Token.Kind == "(" {
		p.nextToken()
		namedArgs = parseCommaSeparatedList(p, p.parseNamedArg)
		rparen = p.expect(")").Pos
	}

	return &ast.ArraySchemaType{
		Array:     pos,
		Gt:        end,
		Item:      t,
		NamedArgs: namedArgs,
		Rparen:    rparen,
	}
}

func (p *Parser) parseIfNotExists() bool {
	if p.Token.Kind == "IF" {
		p.nextToken()
		p.expect("NOT")
		p.expect("EXISTS")
		return true
	}
	return false
}

func (p *Parser) parseIfExists() bool {
	if p.Token.Kind == "IF" {
		p.nextToken()
		p.expect("EXISTS")
		return true
	}
	return false
}

// ================================================================================
//
// DML
//
// ================================================================================

func (p *Parser) parseDML() (dml ast.DML) {
	l := p.Clone()
	defer func() {
		// Panic on tryParseHint()
		if r := recover(); r != nil {
			dml = &ast.BadDML{BadNode: p.handleParseStatementError(r, l)}
		}
	}()

	hint := p.tryParseHint()
	return p.parseDMLInternal(hint)
}

func (p *Parser) parseDMLInternal(hint *ast.Hint) (dml ast.DML) {
	l := p.Clone()
	defer func() {
		if r := recover(); r != nil {
			dml = &ast.BadDML{
				Hint:    hint,
				BadNode: p.handleParseStatementError(r, l),
			}
		}
	}()

	id := p.expect(token.TokenIdent)
	pos := id.Pos
	switch {
	case id.IsKeywordLike("INSERT"):
		return p.parseInsert(pos, hint)
	case id.IsKeywordLike("DELETE"):
		return p.parseDelete(pos, hint)
	case id.IsKeywordLike("UPDATE"):
		return p.parseUpdate(pos, hint)
	}

	panic(p.errorfAtToken(id, "expect pseudo keyword: INSERT, DELETE,  UPDATE but: %s", id.AsString))
}

func (p *Parser) tryParseWithAction() *ast.WithAction {
	if p.Token.Kind != "WITH" {
		return nil
	}

	with := p.expect("WITH").Pos
	action := p.expectKeywordLike("ACTION").Pos
	alias := p.tryParseAsAlias(withRequiredAs)

	return &ast.WithAction{
		With:   with,
		Action: action,
		Alias:  alias,
	}
}

func (p *Parser) tryParseThenReturn() *ast.ThenReturn {
	if p.Token.Kind != "THEN" {
		return nil
	}

	then := p.expect("THEN").Pos
	p.expectKeywordLike("RETURN")
	withAction := p.tryParseWithAction()
	items := parseCommaSeparatedList(p, p.parseSelectItem)

	return &ast.ThenReturn{
		Then:       then,
		WithAction: withAction,
		Items:      items,
	}
}

func (p *Parser) parseInsert(pos token.Pos, hint *ast.Hint) *ast.Insert {
	var insertOrType ast.InsertOrType
	if p.Token.Kind == "OR" {
		p.nextToken()
		switch {
		case p.Token.IsKeywordLike("UPDATE"):
			insertOrType = ast.InsertOrTypeUpdate
		case p.Token.Kind == "IGNORE":
			insertOrType = ast.InsertOrTypeIgnore
		default:
			p.panicfAtToken(&p.Token, "expected pseudo keyword: UPDATE, IGNORE, but: %s", p.Token.AsString)
		}
		p.nextToken()
	}

	if p.Token.Kind == "INTO" {
		p.nextToken()
	}

	name := p.parsePath()
	tableHint := p.tryParseHint()

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

	thenReturn := p.tryParseThenReturn()

	return &ast.Insert{
		Insert:       pos,
		Hint:         hint,
		InsertOrType: insertOrType,
		TableName:    name,
		TableHint:    tableHint,
		Columns:      columns,
		Input:        input,
		ThenReturn:   thenReturn,
	}
}

func (p *Parser) parseValuesInput() *ast.ValuesInput {
	pos := p.expectKeywordLike("VALUES").Pos

	rows := parseCommaSeparatedList(p, p.parseValuesRow)

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

func (p *Parser) parseDelete(pos token.Pos, hint *ast.Hint) *ast.Delete {
	if p.Token.Kind == "FROM" {
		p.nextToken()
	}

	name := p.parsePath()
	tableHint := p.tryParseHint()
	as := p.tryParseAsAlias(withOptionalAs)
	where := p.parseWhere()
	thenReturn := p.tryParseThenReturn()

	return &ast.Delete{
		Delete:     pos,
		Hint:       hint,
		TableName:  name,
		TableHint:  tableHint,
		As:         as,
		Where:      where,
		ThenReturn: thenReturn,
	}
}

func (p *Parser) parseUpdate(pos token.Pos, hint *ast.Hint) *ast.Update {
	name := p.parsePath()
	tableHint := p.tryParseHint()
	as := p.tryParseAsAlias(withOptionalAs)

	p.expect("SET")

	items := parseCommaSeparatedList(p, p.parseUpdateItem)

	where := p.parseWhere()
	thenReturn := p.tryParseThenReturn()

	return &ast.Update{
		Update:     pos,
		Hint:       hint,
		TableName:  name,
		TableHint:  tableHint,
		As:         as,
		Updates:    items,
		Where:      where,
		ThenReturn: thenReturn,
	}
}

func (p *Parser) parseUpdateItem() *ast.UpdateItem {
	path := p.parseIdentOrPath()
	p.expect("=")
	defaultExpr := p.parseDefaultExpr()

	return &ast.UpdateItem{
		Path:        path,
		DefaultExpr: defaultExpr,
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

// ================================================================================
//
// Error Handlers
//
// ================================================================================

func (p *Parser) handleError(r any, l *Lexer) {
	e, ok := r.(*Error)
	if !ok {
		panic(r)
	}

	p.errors = append(p.errors, e)
	p.Lexer = l
}

func (p *Parser) handleParseStatementError(r any, l *Lexer) *ast.BadNode {
	p.handleError(r, l)

	var tokens []*token.Token
	pos := p.Token.Pos
	end := p.Token.Pos
skip:
	for p.Token.Kind != token.TokenEOF {
		switch p.Token.Kind {
		case ";":
			break skip
		}
		end = p.Token.End
		tokens = append(tokens, p.Token.Clone())
		p.Lexer.nextToken(true)
	}

	return &ast.BadNode{
		NodePos: pos,
		NodeEnd: end,
		Tokens:  tokens,
	}
}

func (p *Parser) handleParseQueryExprError(simple bool, r any, l *Lexer) *ast.BadQueryExpr {
	p.handleError(r, l)

	var tokens []*token.Token
	pos := p.Token.Pos
	end := p.Token.Pos
	nesting := 0
skip:
	for p.Token.Kind != token.TokenEOF {
		switch p.Token.Kind {
		case ";":
			break skip
		case "(":
			nesting += 1
		case ")":
			if nesting == 0 {
				break skip
			}
			nesting -= 1
		case "UNION", "INTERSECT", "EXCEPT":
			if simple && nesting == 0 {
				break skip
			}
		}
		end = p.Token.End
		tokens = append(tokens, p.Token.Clone())
		p.Lexer.nextToken(true)
	}

	return &ast.BadQueryExpr{
		BadNode: &ast.BadNode{
			NodePos: pos,
			NodeEnd: end,
			Tokens:  tokens,
		},
	}
}

func (p *Parser) handleParseExprError(r any, l *Lexer) *ast.BadExpr {
	p.handleError(r, l)

	var tokens []*token.Token
	pos := p.Token.Pos
	end := p.Token.Pos
	nesting := 0
skip:
	for p.Token.Kind != token.TokenEOF {
		switch p.Token.Kind {
		case ";":
			break skip
		case "(", "[", "CASE", "WHEN":
			nesting += 1
		case ")", "]", "}", "END", "THEN":
			if nesting == 0 {
				break skip
			}
			nesting -= 1
		case ",", "AS", "FROM", "GROUP", "HAVING", "ORDER", "LIMIT", "OFFSET", "AT", "UNION", "INTERSECT", "EXCEPT":
			if nesting == 0 {
				break skip
			}
		}
		end = p.Token.End
		tokens = append(tokens, p.Token.Clone())
		p.Lexer.nextToken(true)
	}

	return &ast.BadExpr{
		BadNode: &ast.BadNode{
			NodePos: pos,
			NodeEnd: end,
			Tokens:  tokens,
		},
	}
}

func (p *Parser) handleParseTypeError(r any, l *Lexer) *ast.BadType {
	p.handleError(r, l)

	var tokens []*token.Token
	pos := p.Token.Pos
	end := p.Token.Pos
	nesting := 0
skip:
	for p.Token.Kind != token.TokenEOF {
		switch p.Token.Kind {
		case ";", ")":
			break skip
		case "<":
			nesting += 1
		case ">":
			if nesting == 0 {
				break skip
			}
			nesting -= 1
		case ">>":
			if nesting == 0 {
				break skip
			}
			if nesting == 1 {
				p.Token.Kind = ">"
				p.Token.Pos += 1
				break skip
			}
			nesting -= 2
		case ",":
			if nesting == 0 {
				break skip
			}
		}
		tokens = append(tokens, p.Token.Clone())
		end = p.Token.End
		p.Lexer.nextToken(true)
	}

	return &ast.BadType{
		BadNode: &ast.BadNode{
			NodePos: pos,
			NodeEnd: end,
			Tokens:  tokens,
		},
	}
}

// ================================================================================
//
// Utilities
//
// ================================================================================

// parseCommaSeparatedList parses a comma separated list of nodes parsed by `doParse`.
//
// `doParse` should be a reference to a method of `Parser`. That is, this function should always be used on a single line, e.g.:
//
//	columns := parseCommaSeparatedList(p, p.parseIdent)
//
// TODO: create a linter for this.
func parseCommaSeparatedList[T ast.Node](p *Parser, doParse func() T) []T {
	nodes := []T{doParse()}
	for p.Token.Kind == "," {
		p.nextToken()
		nodes = append(nodes, doParse())
	}
	return nodes
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

func (p *Parser) parseRenameTableTo() *ast.RenameTableTo {
	old := p.parseIdent()
	p.expect("TO")
	new := p.parseIdent()

	return &ast.RenameTableTo{
		Old: old,
		New: new,
	}
}

func (p *Parser) parseRenameTable(pos token.Pos) *ast.RenameTable {
	p.expectKeywordLike("TABLE")
	tos := parseCommaSeparatedList(p, p.parseRenameTableTo)

	return &ast.RenameTable{
		Rename: pos,
		Tos:    tos,
	}

}

func (p *Parser) nextToken() {
	p.Lexer.nextToken(false)
}
