package memefish

import (
	"github.com/cloudspannerecosystem/memefish/ast"
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
	default:
		return nil
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
