package memefish

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/ast"
	"github.com/cloudspannerecosystem/memefish/token"
)

// ================================================================================
//
// GRAPH
//
// ================================================================================

func (p *Parser) parseGqlQuery() *ast.GqlGraphQuery {
	graphClause := p.parseGqlGraphClause()
	multiQueryStatement := p.parseGqlMultiLinearQueryStatement()

	return &ast.GqlGraphQuery{GraphClause: graphClause, GqlMultiLinearQueryStatement: multiQueryStatement}
}

func (p *Parser) parseGqlSubquery() *ast.GqlSubQuery {
	lbrace := p.expect("{").Pos
	query := p.parseGqlQueryExpr()
	rbrace := p.expect("}").Pos

	return &ast.GqlSubQuery{LBrace: lbrace, RBrace: rbrace, Query: query}
}
func (p *Parser) parseGqlQueryExpr() *ast.GqlQueryExpr {
	var graphClause *ast.GqlGraphClause
	if p.Token.IsKeywordLike("GRAPH") {
		graphClause = p.parseGqlGraphClause()
	}
	multiQueryStatement := p.parseGqlMultiLinearQueryStatement()

	return &ast.GqlQueryExpr{GraphClause: graphClause, GqlMultiLinearQueryStatement: multiQueryStatement}
}

func (p *Parser) parseGqlGraphClause() *ast.GqlGraphClause {
	graphPos := p.expectKeywordLike("GRAPH").Pos
	graphName := p.parseIdent()
	return &ast.GqlGraphClause{
		GqlGraph:             graphPos,
		GqlPropertyGraphName: graphName,
	}
}

func (p *Parser) parseGqlLinearQueryStatement() ast.GqlLinearQueryStatement {
	stmt := p.parseGqlSimpleLinearQueryStatement()
	if p.Token.Kind != ("UNION") && p.Token.Kind != ("EXCEPT") && p.Token.Kind != ("DISTINCT") {
		return stmt
	}
	var tail []*ast.GqlSimpleLinearQueryStatementWithSetOperator
	for p.Token.Kind == ("UNION") || p.Token.Kind == ("EXCEPT") || p.Token.Kind == ("DISTINCT") {
		tail = append(tail, p.parseGqlSimpleLinearQueryStatementWithSetOperator())
	}
	return &ast.GqlCompositeLinearQueryStatement{HeadSimpleLinearQueryStatement: stmt, TailSimpleLinearQueryStatementList: tail}
}

func (p *Parser) parseGqlSimpleLinearQueryStatementWithSetOperator() *ast.GqlSimpleLinearQueryStatementWithSetOperator {
	pos := p.Token.Pos
	var op ast.GqlSetOperatorEnum
	switch p.Token.Kind {
	case "UNION":
		p.nextToken()
		op = ast.GqlSetOperatorUnion
	case "INTERSECT":
		p.nextToken()
		op = ast.GqlSetOperatorIntersect
	case "EXCEPT":
		p.nextToken()
		op = ast.GqlSetOperatorExcept
	default:
		panic("unknown SetOperatorEnum")
	}

	allOrDistinct := p.tryParseGqlAllOrDistinctEnum()
	if allOrDistinct == ast.GqlAllOrDistinctImplicitAll {
		p.panicfAtToken(&p.Token, `expect "ALL" or "DISTINCT, but: %v`, p.Token.Kind)
	}

	stmt := p.parseGqlSimpleLinearQueryStatement()
	return &ast.GqlSimpleLinearQueryStatementWithSetOperator{
		StartPos:      pos,
		SetOperator:   op,
		DistinctOrAll: allOrDistinct,
		Statement:     stmt,
	}
}
func (p *Parser) parseGqlSimpleLinearQueryStatement() *ast.GqlSimpleLinearQueryStatement {
	var stmts []ast.GqlPrimitiveQueryStatement
	for stmt := p.tryParseGqlPrimitiveQueryStatement(); stmt != nil; stmt = p.tryParseGqlPrimitiveQueryStatement() {
		stmts = append(stmts, stmt)
	}

	return &ast.GqlSimpleLinearQueryStatement{PrimitiveQueryStatementList: stmts}
}

func (p *Parser) tryParseGqlWithOffsetClause() *ast.GqlWithOffsetClause {
	if p.Token.Kind != "WITH" {
		return nil
	}
	with := p.expect("WITH").Pos
	offset := p.expectKeywordLike("OFFSET").Pos
	var offsetName *ast.Ident
	if p.Token.Kind == "AS" {
		p.nextToken()
		offsetName = p.parseIdent()
	}

	return &ast.GqlWithOffsetClause{
		With:       with,
		Offset:     offset,
		OffsetName: offsetName,
	}
}

func (p *Parser) parseGqlForStatement() *ast.GqlForStatement {
	forPos := p.expect("FOR").Pos
	name := p.parseIdent()
	p.expect("IN")
	expr := p.parseExpr()

	withOffset := p.tryParseGqlWithOffsetClause()

	return &ast.GqlForStatement{
		For:              forPos,
		ElementName:      name,
		ArrayExpression:  expr,
		WithOffsetClause: withOffset,
	}
}

func (p *Parser) tryParseGqlOrderBySpecificationList() []*ast.GqlOrderBySpecification {
	var list []*ast.GqlOrderBySpecification
	for {
		expr := p.parseExpr()
		var direction ast.GqlDirection
		directionPos := p.Token.Pos
		switch {
		case p.Token.Kind == "ASC":
			p.nextToken()
			direction = ast.GqlSortOrderAsc
		case p.Token.Kind == "DESC":
			p.nextToken()
			direction = ast.GqlSortOrderDesc
		case p.Token.Kind == "ASCENDING":
			p.nextToken()
			direction = ast.GqlSortOrderAscending
		case p.Token.Kind == "DESCENDING":
			p.nextToken()
			direction = ast.GqlSortOrderDescending
		default:
			directionPos = token.InvalidPos
		}
		list = append(list, &ast.GqlOrderBySpecification{Expr: expr, DirectionPos: directionPos, Direction: direction})
		if p.Token.Kind != "," {
			break
		}
	}
	return list
}

func (p *Parser) parseGqlOrderByStatement() *ast.GqlOrderByStatement {
	order := p.expect("ORDER").Pos
	p.expect("BY")
	spec := p.tryParseGqlOrderBySpecificationList()
	if spec == nil {
		p.panicfAtToken(&p.Token, "expect at least one order_by_specification")
	}

	return &ast.GqlOrderByStatement{
		Order:                    order,
		OrderBySpecificationList: spec,
	}
}

func (p *Parser) parseGqlLimitStatement() *ast.GqlLimitStatement {
	limit := p.expect("LIMIT").Pos
	expr := p.parseIntValue()

	return &ast.GqlLimitStatement{
		Limit: limit,
		Count: expr,
	}
}

func (p *Parser) parseGqlFilterStatement() *ast.GqlFilterStatement {
	filter := p.expectKeywordLike("FILTER").Pos
	where := token.InvalidPos
	if p.Token.Kind == "WHERE" {
		where = p.expect("WHERE").Pos
	}
	expr := p.parseExpr()

	return &ast.GqlFilterStatement{
		Filter: filter,
		Where:  where,
		Expr:   expr,
	}
}

func (p *Parser) parseGqlOffsetStatement() *ast.GqlOffsetStatement {
	var pos token.Pos
	var isSkip bool
	switch {
	case p.Token.IsKeywordLike("OFFSET"):
		pos = p.expectKeywordLike("OFFSET").Pos
	case p.Token.IsKeywordLike("SKIP"):
		pos = p.expectKeywordLike("SKIP").Pos
		isSkip = true
	default:
		p.panicfAtToken(&p.Token, `expected "OFFSET" or "SKIP", but: %v`, p.Token.Kind)
	}

	count := p.parseIntValue()

	return &ast.GqlOffsetStatement{
		Offset: pos,
		IsSkip: isSkip,
		Count:  count,
	}
}

func (p *Parser) parseGqlReturnItem() ast.GqlReturnItem {
	if p.Token.Kind == "*" {
		pos := p.expect("*").Pos
		return &ast.Star{
			Star: pos,
		}
	}

	expr := p.parseExpr()
	if p.Token.Kind == "AS" {
		if as := p.tryParseAsAlias(); as != nil {
			return &ast.Alias{
				Expr: expr,
				As:   as,
			}
		}
	}

	return &ast.ExprSelectItem{
		Expr: expr,
	}
}

func (p *Parser) tryParseAsAliasAsMandatory() *ast.AsAlias {
	if p.Token.Kind != "AS" {
		return nil
	}
	return p.tryParseAsAlias()
}

func (p *Parser) tryParseGqlPrimitiveQueryStatement() ast.GqlPrimitiveQueryStatement {
	switch {
	case p.Token.IsKeywordLike("RETURN"):
		return p.parseGqlReturnStatement()
	case p.Token.IsKeywordLike("LET"):
		return p.parseGqlLetStatement()
	case p.Token.IsKeywordLike("FILTER"):
		return p.parseGqlFilterStatement()
	case p.Token.Kind == "ORDER":
		return p.parseGqlOrderByStatement()
	case p.Token.Kind == "LIMIT":
		return p.parseGqlLimitStatement()
	case p.Token.IsKeywordLike("SKIP"), p.Token.IsKeywordLike("OFFSET"):
		return p.parseGqlOffsetStatement()
	case p.Token.Kind == "FOR":
		return p.parseGqlForStatement()
	case p.Token.IsKeywordLike("OPTIONAL") || p.Token.IsKeywordLike("MATCH"):
		return p.parseGqlMatchStatement()
	case p.Token.Kind == "WITH":
		return p.parseGqlWithStatement()
	default:
		return nil
	}
}

func (p *Parser) tryParseGqlPathModePathOrPaths() *ast.Ident {
	switch {
	case p.Token.IsKeywordLike("PATH"), p.Token.IsKeywordLike("PATHS"):
		return p.parseIdent()
	default:
		return nil
	}
}
func (p *Parser) tryParseGqlPathMode() *ast.GqlPathMode {
	switch {
	case p.Token.IsKeywordLike("WALK"):
		pathModeToken := p.parseIdent()
		pathOrPathsToken := p.tryParseGqlPathModePathOrPaths()
		return &ast.GqlPathMode{
			ModeToken:       pathModeToken,
			PathOrPathToken: pathOrPathsToken,
			Mode:            ast.GqlPathModeWalk,
		}

	case p.Token.IsKeywordLike("TRAIL"):
		pathModeToken := p.parseIdent()
		pathOrPathsToken := p.tryParseGqlPathModePathOrPaths()
		return &ast.GqlPathMode{
			ModeToken:       pathModeToken,
			PathOrPathToken: pathOrPathsToken,
			Mode:            ast.GqlPathModeTrail,
		}
	default:
		return nil
	}
}
func (p *Parser) tryParseGqlPathSearchPrefixOrPathMode() ast.GqlPathSearchPrefixOrPathMode {
	startPos := p.Token.Pos
	switch {
	case p.Token.Kind == "ALL":
		p.nextToken()
		return &ast.GqlPathSearchPrefix{
			StartPos:     startPos,
			SearchPrefix: ast.GqlPathSearchPrefixAll,
		}
	case p.Token.Kind == "ANY":
		p.nextToken()
		if p.Token.IsKeywordLike("SHORTEST") {
			p.nextToken()
			return &ast.GqlPathSearchPrefix{
				StartPos:     startPos,
				SearchPrefix: ast.GqlPathSearchPrefixAnyShortest,
			}
		} else {
			return &ast.GqlPathSearchPrefix{
				StartPos:     startPos,
				SearchPrefix: ast.GqlPathSearchPrefixAny,
			}
		}
	default:
		if pathMode := p.tryParseGqlPathMode(); pathMode != nil {
			return pathMode
		}
		return nil
	}

}

func (p *Parser) parseGqlMatchStatement() *ast.GqlMatchStatement {
	var optionalPos token.Pos
	if p.Token.IsKeywordLike("OPTIONAL") {
		optionalPos = p.Token.Pos
		p.nextToken()
	} else {
		optionalPos = token.InvalidPos
	}

	match := p.expectKeywordLike("MATCH").Pos
	hint := p.tryParseHint()
	var prefixOrMode = p.tryParseGqlPathSearchPrefixOrPathMode()
	pattern := p.parseGqlGraphPattern()
	return &ast.GqlMatchStatement{
		Optional:     optionalPos,
		Match:        match,
		MatchHint:    hint,
		PrefixOrMode: prefixOrMode,
		GraphPattern: pattern,
	}
}

func (p *Parser) parseGqlLabelName() *ast.GqlLabelName {
	if p.Token.Kind == "%" {
		p.nextToken()
		return &ast.GqlLabelName{IsPercent: true}
	}
	return &ast.GqlLabelName{
		IsPercent: false,
		LabelName: p.parseIdent(),
	}
}

func (p *Parser) parseGqlLabelExpression() ast.GqlLabelExpression {
	// TODO
	/*
		label_expression:
		  {
		    label_name
		    | or_expression
		    | and_expression
		    | not_expression
		  }
	*/
	var labelTerm ast.GqlLabelExpression
	switch {
	case p.Token.Kind == "!":
		notPos := p.expect("!").Pos
		expr := p.parseGqlLabelExpression()
		labelTerm = &ast.GqlLabelNotExpression{NotPos: notPos, LabelExpression: expr}
	case p.Token.Kind == "(":
		lparen := p.expect("(").Pos
		expr := p.parseGqlLabelExpression()
		rparen := p.expect(")").Pos
		labelTerm = &ast.GqlLabelParenExpression{
			LParen:    lparen,
			RParen:    rparen,
			LabelExpr: expr,
		}
	case p.Token.Kind == token.TokenIdent || p.Token.Kind == "%":
		labelTerm = p.parseGqlLabelName()
	default:
		p.panicfAtToken(&p.Token, `expected token: ",", "(", or "<ident>", but: %v`, p.Token.Kind)
		return nil
	}
	switch {
	case p.Token.Kind == "|":
		p.nextToken()
		right := p.parseGqlLabelExpression()
		return &ast.GqlLabelOrExpression{
			Left:  labelTerm,
			Right: right,
		}
	case p.Token.Kind == "&":
		p.nextToken()
		right := p.parseGqlLabelExpression()
		return &ast.GqlLabelAndExpression{
			Left:  labelTerm,
			Right: right,
		}
	}
	return labelTerm
}
func (p *Parser) parseGqlIsLabelCondition() *ast.GqlIsLabelCondition {

	pos := p.Token.Pos
	p.nextToken()
	labelExpr := p.parseGqlLabelExpression()
	return &ast.GqlIsLabelCondition{
		IsOrColon:       pos,
		LabelExpression: labelExpr,
	}
}
func (p *Parser) parseGqlPatternFiller() *ast.GqlPatternFiller {
	hint := p.tryParseHint()
	var patternVar *ast.Ident
	if p.Token.Kind == token.TokenIdent {
		patternVar = p.parseIdent()
	}
	var isLabelCondition *ast.GqlIsLabelCondition
	if p.Token.Kind == ":" || p.Token.Kind == "IS" {
		isLabelCondition = p.parseGqlIsLabelCondition()
	}

	filter := p.tryParseGqlPatternFillerFilter()
	return &ast.GqlPatternFiller{
		Hint:                 hint,
		GraphPatternVariable: patternVar,       // TODO
		IsLabelCondition:     isLabelCondition, // TODO
		Filter:               filter,           // TODO
	}
}

func (p *Parser) lookaheadGqlSubpathPattern() bool {
	lexer := p.Lexer.Clone()
	defer func() {
		p.Lexer = lexer
	}()

	if p.Token.Kind != "(" {
		return false
	}
	p.nextToken()
	if p.Token.Kind == "(" || p.Token.Kind == "-" || p.Token.Kind == "<" || p.Token.IsKeywordLike("WALK") || p.Token.IsKeywordLike("TRAIL") {
		return true
	}
	return false
}

func (p *Parser) parseGqlNodePattern() *ast.GqlNodePattern {
	lparen := p.expect("(").Pos
	filter := p.parseGqlPatternFiller()
	rparen := p.expect(")").Pos
	return &ast.GqlNodePattern{
		LParen:        lparen,
		RParen:        rparen,
		PatternFiller: filter,
	}
}

func (p *Parser) tryParseGqlQuantifier() ast.GqlQuantifier {
	if p.Token.Kind != "{" {
		return nil
	}
	lbrace := p.expect("{").Pos
	if p.Token.Kind == "," {
		upperBound := p.parseIntValue()
		rbrace := p.expect("}").Pos
		return &ast.GqlBoundedQuantifier{
			LBrace:     lbrace,
			RBrace:     rbrace,
			UpperBound: upperBound,
		}
	}
	bound := p.parseIntValue()
	if p.Token.Kind != "," {
		rbrace := p.expect("}").Pos
		return &ast.GqlFixedQuantifier{
			LBrace: lbrace,
			RBrace: rbrace,
			Bound:  bound,
		}
	}
	p.expect(",")
	upperBound := p.parseIntValue()
	rbrace := p.expect("}").Pos

	return &ast.GqlBoundedQuantifier{
		LBrace:     lbrace,
		RBrace:     rbrace,
		LowerBound: bound,
		UpperBound: upperBound,
	}
}

func (p *Parser) tryParseGqlQuantifiablePathTerm() *ast.GqlQuantifiablePathTerm {
	hint := p.tryParseHint()
	pt := p.tryParseGqlPathTerm()
	if pt == nil {
		return nil
	}
	q := p.tryParseGqlQuantifier()
	return &ast.GqlQuantifiablePathTerm{
		Hint:       hint,
		PathTerm:   pt,
		Quantifier: q,
	}
}

func (p *Parser) tryParseGqlPathTerm() ast.GqlPathTerm {
	// TODO
	/*
		path_term:
		  {
		    element_pattern
		    | subpath_pattern
		  }
	*/

	if p.lookaheadGqlSubpathPattern() {
		return p.parseGqlSubPathPattern()
	}
	if p.Token.Kind == "(" {
		return p.parseGqlNodePattern()
	}
	if p.Token.Kind == "-" || p.Token.Kind == "<" {
		return p.parseGqlEdgePattern()
	}
	return nil
}

func (p *Parser) parseGqlSubPathPattern() *ast.GqlSubpathPattern {
	lparen := p.expect("(").Pos
	pathMode := p.tryParseGqlPathMode()
	pattern := p.parseGqlPathPattern()
	where := p.tryParseWhere()
	rparen := p.expect(")").Pos

	return &ast.GqlSubpathPattern{
		LParen:      lparen,
		RParen:      rparen,
		PathMode:    pathMode,
		PathPattern: pattern,
		WhereClause: where,
	}
}

func (p *Parser) parseGqlEdgePattern() ast.GqlEdgePattern {
	// NOTE: Currently, this implementation allows whitespace in any place of "<-["  and "]->".
	// It is more relaxed syntax than spec.
	// It will be fixed when we are sure that the modification of lexer is safe.
	firstPos := p.Token.Pos
	switch p.Token.Kind {
	case "<":
		p.nextToken()
		firstHyphenPos := p.expect("-").Pos
		if p.Token.Kind != "[" {
			return &ast.GqlAbbreviatedEdgeLeft{
				First: firstPos,
				Last:  firstHyphenPos,
			}
		}
		p.nextToken()
		patternFilter := p.parseGqlPatternFiller()
		p.expect("]")
		lastPos := p.expect("-").Pos
		return &ast.GqlFullEdgeLeft{
			First:         firstPos,
			Last:          lastPos,
			PatternFiller: patternFilter,
		}
	case "-":
		p.nextToken()
		switch p.Token.Kind {
		case ">":
			lastPos := p.expect(">").Pos
			return &ast.GqlAbbreviatedEdgeRight{
				First: firstPos,
				Last:  lastPos,
			}
		case "[":
			p.nextToken()
			patternFiller := p.parseGqlPatternFiller()
			p.expect("]")
			lastHyphenPos := p.expect("-").Pos
			if p.Token.Kind == ">" {
				lastPos := p.Token.Pos
				p.nextToken()
				return &ast.GqlFullEdgeRight{
					First:         firstPos,
					Last:          lastPos,
					PatternFiller: patternFiller,
				}
			}
			return &ast.GqlFullEdgeAny{
				First:         firstPos,
				Last:          lastHyphenPos,
				PatternFiller: patternFiller,
			}
		default:
			return &ast.GqlAbbreviatedEdgeAny{Hyphen: firstPos}
		}
	default:
		panic(fmt.Sprintf("not implemented kind: %v %v", p.Token.Kind, p.Token.Raw))
	}
}

func (p *Parser) parseGqlPathPattern() *ast.GqlPathPattern {
	list := []*ast.GqlQuantifiablePathTerm{p.tryParseGqlQuantifiablePathTerm()}
	for p.Token.Kind != ")" && !p.Token.IsKeywordLike("WHERE") {
		term := p.tryParseGqlQuantifiablePathTerm()
		if term == nil {
			break
		}
		list = append(list, term)
	}
	return &ast.GqlPathPattern{
		PathTermList: list,
	}
}
func (p *Parser) parseGqlTopLevelPathPattern() *ast.GqlTopLevelPathPattern {
	pattern := p.parseGqlPathPattern()
	return &ast.GqlTopLevelPathPattern{
		PathSearchPrefixOrPathMode: nil, // TODO
		PathPattern:                pattern,
	}
}

func (p *Parser) parseGqlElementProperty() *ast.GqlElementProperty {
	ident := p.parseIdent()
	_ = p.expect(":")
	expr := p.parseExpr()
	return &ast.GqlElementProperty{
		ElementPropertyName:  ident,
		ElementPropertyValue: expr,
	}
}
func (p *Parser) parseGqlPropertyFilters() *ast.GqlPropertyFilters {
	lbrace := p.expect("{").Pos

	elementPropertyList := []*ast.GqlElementProperty{p.parseGqlElementProperty()}
	for p.Token.Kind == "," {
		p.nextToken()
		elementPropertyList = append(elementPropertyList, p.parseGqlElementProperty())
	}

	rbrace := p.expect("}").Pos
	return &ast.GqlPropertyFilters{
		LBrace:                 lbrace,
		PropertyFilterElemList: elementPropertyList,
		RBrace:                 rbrace,
	}
}

func (p *Parser) tryParseGqlPatternFillerFilter() ast.GqlPatternFillerFilter {
	switch {
	case p.Token.Kind == "WHERE":
		return p.parseWhere()
	case p.Token.Kind == "{":
		return p.parseGqlPropertyFilters()
	default:
		return nil
	}
}

func (p *Parser) parseGqlGraphPattern() *ast.GqlGraphPattern {
	patterns := []*ast.GqlTopLevelPathPattern{p.parseGqlTopLevelPathPattern()}
	for p.Token.Kind == "," {
		p.nextToken()
		patterns = append(patterns, p.parseGqlTopLevelPathPattern())
	}
	return &ast.GqlGraphPattern{
		PathPatternList: patterns,
		WhereClause:     p.tryParseWhere(), // TODO
	}
}

/*
***
GQL statements
***
*/

/*
***
GQL RETURN statement
***
*/
func (p *Parser) parseGqlReturnStatement() *ast.GqlReturnStatement {
	ret := p.expectKeywordLike("RETURN").Pos

	allOrDistinct := p.tryParseGqlAllOrDistinctEnum()
	returnItems := p.parseGqlReturnItemList()
	groupBy := p.tryParseGroupBy()
	orderBy := p.tryParseOrderBy()
	limit := p.tryParseLimit()
	offset := p.tryParseOffset()
	var limitAndOffset ast.GqlLimitAndOffsetClause
	switch {
	case limit != nil && offset != nil:
		limitAndOffset = &ast.GqlLimitWithOffsetClause{
			Limit:  limit,
			Offset: offset,
		}
	case limit != nil:
		limitAndOffset = &ast.GqlLimitClause{Limit: limit}
	case offset != nil:
		limitAndOffset = &ast.GqlOffsetClause{Offset: offset}
	default:
	}
	return &ast.GqlReturnStatement{Return: ret,
		AllOrDistinct:        allOrDistinct,
		GroupByClause:        groupBy,
		OrderByClause:        orderBy,
		LimitAndOffsetClause: limitAndOffset,
		ReturnItemList:       returnItems}
}

func (p *Parser) parseGqlReturnItemList() []ast.GqlReturnItem {
	var results []ast.GqlReturnItem
	for {
		results = append(results, p.parseGqlReturnItem())
		if p.Token.Kind != "," {
			break
		}
		p.nextToken()
	}
	return results
}

func (p *Parser) tryParseGqlAllOrDistinctEnum() ast.GqlAllOrDistinctEnum {
	switch p.Token.Kind {
	case "ALL":
		p.nextToken()
		return ast.GqlAllOrDistinctAll
	case "DISTINCT":
		p.nextToken()
		return ast.GqlAllOrDistinctDistinct
	default:
		return ast.GqlAllOrDistinctImplicitAll
	}
}

func (p *Parser) parseGqlWithStatement() *ast.GqlWithStatement {
	with := p.expect("WITH").Pos
	allOrDistinct := p.tryParseGqlAllOrDistinctEnum()
	returnItems := p.parseGqlReturnItemList()
	groupBy := p.tryParseGroupBy()

	return &ast.GqlWithStatement{
		With:           with,
		AllOrDistinct:  allOrDistinct,
		GroupByClause:  groupBy,
		ReturnItemList: returnItems}
}

func (p *Parser) parseGqlLinearGraphVariable() *ast.GqlLinearGraphVariable {
	name := p.parseIdent()
	p.expect("=")
	value := p.parseExpr()

	return &ast.GqlLinearGraphVariable{VariableName: name, Value: value}
}

func (p *Parser) parseGqlLetStatement() *ast.GqlLetStatement {
	let := p.expectKeywordLike("LET").Pos
	vars := []*ast.GqlLinearGraphVariable{p.parseGqlLinearGraphVariable()}

	for p.Token.Kind == "," {
		p.nextToken()
		vars = append(vars, p.parseGqlLinearGraphVariable())
	}

	return &ast.GqlLetStatement{Let: let, LinearGraphVariableList: vars}
}

func (p *Parser) parseGqlMultiLinearQueryStatement() *ast.GqlMultiLinearQueryStatement {
	items := []ast.GqlLinearQueryStatement{p.parseGqlLinearQueryStatement()}
	for p.Token.IsKeywordLike("NEXT") {
		p.expectKeywordLike("NEXT")
		items = append(items, p.parseGqlLinearQueryStatement())
	}

	return &ast.GqlMultiLinearQueryStatement{GqlLinearQueryStatements: items}
}
