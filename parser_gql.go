package memefish

import (
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
	// TODO process ast.GqlCompositeQueryStatement
	return p.parseGqlSimpleLinearQueryStatement()
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

func (p *Parser) parseGqlPatternFilter() *ast.GqlPatternFilter {
	return &ast.GqlPatternFilter{
		GraphPatternVariable: nil, // TODO
		IsLabelCondition:     nil, // TODO
		Filter:               nil, // TODO
	}
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
func (p *Parser) parseGqlPathTerm() ast.GqlPathTerm {
	// TODO
	/*
		path_term:
		  {
		    element_pattern
		    | subpath_pattern
		  }
	*/

	return p.parseGqlNodePattern()
}

func (p *Parser) parseGqlPathPattern() *ast.GqlPathPattern {
	list := []ast.GqlPathTerm{p.parseGqlPathTerm()}
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
func (p *Parser) parseGqlGraphPattern() *ast.GqlGraphPattern {
	patterns := []*ast.GqlTopLevelPathPattern{p.parseGqlTopLevelPathPattern()}
	for p.Token.Kind == "," {
		p.nextToken()
		patterns = append(patterns, p.parseGqlTopLevelPathPattern())
	}
	return &ast.GqlGraphPattern{
		PathPatternList:     patterns,
		OptionalWhereClause: nil, // TODO
	}
}

func (p *Parser) parseGqlReturnStatement() *ast.GqlReturnStatement {
	ret := p.expectKeywordLike("RETURN").Pos
	returnItems := p.parseSelectResults()
	return &ast.GqlReturnStatement{Return: ret, ReturnItemList: returnItems}
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
	for false {
		items = append(items, p.parseGqlLinearQueryStatement())
	}

	return &ast.GqlMultiLinearQueryStatement{GqlLinearQueryStatements: items}
}
