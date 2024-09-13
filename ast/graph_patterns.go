package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
)

type GqlGraphPattern struct {
	PathPatternList     []*GqlTopLevelPathPattern
	OptionalWhereClause *GqlWhereClause // optional
}

func (g GqlGraphPattern) Pos() token.Pos {
	return firstPos(g.PathPatternList)
}

func (g GqlGraphPattern) End() token.Pos {
	if g.OptionalWhereClause != nil {
		return g.OptionalWhereClause.End()
	}
	return lastEnd(g.PathPatternList)
}

func (g GqlGraphPattern) SQL() string {
	var sql string
	if g.OptionalWhereClause != nil {
		sql += g.OptionalWhereClause.SQL() + " "
	}
	sql += g.PathPatternList[0].SQL()
	for _, pp := range g.PathPatternList[1:] {
		sql += ", " + pp.SQL()
	}
	return sql
}

type GqlTopLevelPathPattern struct {
	PathSearchPrefixOrPathMode GqlPathSearchPrefixOrPathMode // optional
	PathPattern                *GqlPathPattern
}

func (g GqlTopLevelPathPattern) Pos() token.Pos {
	if g.PathSearchPrefixOrPathMode != nil {
		return g.PathSearchPrefixOrPathMode.Pos()
	}
	return g.PathPattern.Pos()
}

func (g GqlTopLevelPathPattern) End() token.Pos {
	return g.PathPattern.End()
}

func (g GqlTopLevelPathPattern) SQL() string {
	var sql string
	if g.PathSearchPrefixOrPathMode != nil {
		sql += g.PathSearchPrefixOrPathMode.SQL() + " "
	}
	sql += g.PathPattern.SQL()
	return sql
}

// GqlPathSearchPrefixOrPathMode TODO
/*
{ path_search_prefix | path_mode }
*/
type GqlPathSearchPrefixOrPathMode interface {
	Node
	isGqlPathSearchPrefixOrPathMode()
}

type GqlPathPattern struct {
	PathTermList []GqlPathTerm
}

func (g GqlPathPattern) Pos() token.Pos {
	return firstPos(g.PathTermList)
}

func (g GqlPathPattern) End() token.Pos {
	return lastEnd(g.PathTermList)
}

func (g GqlPathPattern) SQL() string {
	var sql string
	sql += g.PathTermList[0].SQL()
	for _, elem := range g.PathTermList[1:] {
		sql += elem.SQL()
	}
	return sql
}

type GqlPathTerm interface {
	Node
	isGqlPathTerm()
}

type GqlWhereClause struct {
	Where          token.Pos
	BoolExpression Expr
}

func (g GqlWhereClause) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlWhereClause) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlWhereClause) SQL() string {
	//TODO implement me
	panic("implement me")
}

func (g GqlWhereClause) isGqlPatternFilterFilter() {}

type GqlElementPattern interface {
	Node
	GqlPathTerm
	isGqlElementPattern()
}

type GqlNodePattern struct {
	LParen, RParen token.Pos
	PatternFilter  *GqlPatternFilter
}

func (g GqlNodePattern) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlNodePattern) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlNodePattern) SQL() string {
	var sql string
	// TODO PatternFilter
	sql += fmt.Sprintf("( )")
	return sql
}

func (g GqlNodePattern) isGqlPathTerm() {}

// EdgePattern TODO
/*
edge_pattern:
  {
    full_edge_any |
    full_edge_left |
    full_edge_right |
    abbreviated_edge_any |
    abbreviated_edge_left |
    abbreviated_edge_right
  }
*/
type EdgePattern interface {
	Node
	isEdgePattern()
}

type GqlPatternFilter struct {
	GraphPatternVariable *GqlGraphPatternVariable // optional
	IsLabelCondition     *GqlIsLabelCondition     // optional
	Filter               *GqlPatternFilterFilter  // optional
}

type GqlGraphPatternVariable struct{}
type GqlIsLabelCondition struct {
	IsOrColon       token.Pos
	LabelExpression GqlLabelExpression
}

// GqlLabelExpression TODO
/*
label_expression:
  {
    label_name
    | or_expression
    | and_expression
    | not_expression
  }
*/
type GqlLabelExpression interface {
	Node
	isGqlLabelExpression()
}

type GqlPatternFilterFilter interface {
	Node
	isGqlPatternFilterFilter()
}

type GqlPropertyFilters struct {
	LBrace                 token.Pos
	PropertyFilterElemList []*GqlPropertyFilterElem
	RBrace                 token.Pos
}

type GqlPropertyFilterElem struct {
	Name  *Ident
	Value Expr
}

type GqlElementProperty struct {
	ElementPropertyName  *Ident
	ElementPropertyValue Expr
}
