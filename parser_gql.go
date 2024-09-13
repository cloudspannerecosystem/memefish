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

	var allOrDistinct ast.GqlAllOrDistinctEnum
	switch p.Token.Kind {
	case "ALL":
		p.nextToken()
		allOrDistinct = ast.GqlAllOrDistinctAll
	case "DISTINCT":
		p.nextToken()
		allOrDistinct = ast.GqlAllOrDistinctDistinct
	default:
		panic("unknown SetAllOrDistinct")
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

func (p *Parser) tryParseGqlPrimitiveQueryStatement() ast.GqlPrimitiveQueryStatement {

	switch {
	case p.Token.IsKeywordLike("RETURN"):
		return p.parseGqlReturnStatement()
	case p.Token.IsKeywordLike("LET"):
		return p.parseGqlLetStatement()
	case p.Token.IsKeywordLike("OPTIONAL") || p.Token.IsKeywordLike("MATCH"):
		return p.parseGqlMatchStatement()
	default:
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
	pattern := p.parseGqlGraphPattern()
	return &ast.GqlMatchStatement{Optional: optionalPos, Match: match, GraphPattern: pattern}
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
	return p.parseGqlLabelName()
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
func (p *Parser) parseGqlPatternFilter() *ast.GqlPatternFilter {
	var patternVar *ast.Ident
	if p.Token.Kind == token.TokenIdent {
		patternVar = p.parseIdent()
	}
	var isLabelCondition *ast.GqlIsLabelCondition
	if p.Token.Kind == ":" || p.Token.Kind == "IS" {
		isLabelCondition = p.parseGqlIsLabelCondition()
	}

	filter := p.tryParseGqlPatternFilterFilter()
	return &ast.GqlPatternFilter{
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
	if p.Token.Kind == "(" || p.Token.Kind == "-" || p.Token.Kind == "<" {
		return true
	}
	return false
}

func (p *Parser) parseGqlNodePattern() *ast.GqlNodePattern {
	lparen := p.expect("(").Pos
	filter := p.parseGqlPatternFilter()
	rparen := p.expect(")").Pos
	return &ast.GqlNodePattern{
		LParen:        lparen,
		RParen:        rparen,
		PatternFilter: filter,
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

func (p *Parser) tryParseGqlPathMode() *ast.GqlPathMode {
	if !p.Token.IsKeywordLike("TRAIL") && !p.Token.IsKeywordLike("WALK") {
		return nil
	}
	var pathMode ast.GqlPathModeEnum
	if p.Token.IsKeywordLike("TRAIL") {
		p.nextToken()
		pathMode = ast.GqlPathModeTrail
	} else if p.Token.IsKeywordLike("WALK") {
		p.nextToken()
		pathMode = ast.GqlPathModeWalk
	}
	if p.Token.IsKeywordLike("PATH") || p.Token.IsKeywordLike("PATHS") {
		p.nextToken()
	}
	return &ast.GqlPathMode{
		Mode: pathMode,
	}

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
	//TODO implement
	switch p.Token.Kind {
	case "<":
		p.nextToken()
		p.expect("-")
		if p.Token.Kind != "[" {
			return &ast.GqlAbbreviatedEdgeLeft{}
		}
		p.nextToken()
		patternFilter := p.parseGqlPatternFilter()
		p.expect("]")
		p.expect("-")
		return &ast.GqlFullEdgeLeft{PatternFilter: patternFilter}
	case "-":
		p.nextToken()
		switch p.Token.Kind {
		case ">":
			p.nextToken()
			return &ast.GqlAbbreviatedEdgeRight{}
		case "[":
			p.nextToken()
			patternFilter := p.parseGqlPatternFilter()
			p.expect("]")
			p.expect("-")
			if p.Token.Kind == ">" {
				p.nextToken()
				return &ast.GqlFullEdgeRight{PatternFilter: patternFilter}
			}
			return &ast.GqlFullEdgeAny{PatternFilter: patternFilter}
		default:
			return &ast.GqlAbbreviatedEdgeAny{}
		}
	default:
		panic(fmt.Sprintf("not implemented kind: %v %v", p.Token.Kind, p.Token.Raw))
	}
}

func (p *Parser) parseGqlPathPattern() *ast.GqlPathPattern {
	// TODO list
	list := []ast.GqlPathTerm{p.tryParseGqlPathTerm()}
	for p.Token.Kind != ")" && !p.Token.IsKeywordLike("WHERE") {
		term := p.tryParseGqlPathTerm()
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

func (p *Parser) tryParseGqlPatternFilterFilter() ast.GqlPatternFilterFilter {
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

func (p *Parser) parseGqlReturnStatement() *ast.GqlReturnStatement {
	ret := p.expectKeywordLike("RETURN").Pos
	var allOrDistinct ast.GqlAllOrDistinctEnum
	switch p.Token.Kind {
	case "ALL":
		p.nextToken()
		allOrDistinct = ast.GqlAllOrDistinctAll
	case "DISTINCT":
		p.nextToken()
		allOrDistinct = ast.GqlAllOrDistinctDistinct
	}
	returnItems := p.parseSelectResults()
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
