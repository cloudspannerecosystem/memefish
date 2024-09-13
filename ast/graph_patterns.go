package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
	"strings"
)

type GqlGraphPattern struct {
	PathPatternList []*GqlTopLevelPathPattern
	WhereClause     *Where // optional
}

func (g GqlGraphPattern) Pos() token.Pos {
	return firstPos(g.PathPatternList)
}

func (g GqlGraphPattern) End() token.Pos {
	if g.WhereClause != nil {
		return g.WhereClause.End()
	}
	return lastEnd(g.PathPatternList)
}

func (g GqlGraphPattern) SQL() string {
	var sql string
	sql += g.PathPatternList[0].SQL()
	for _, pp := range g.PathPatternList[1:] {
		sql += ", " + pp.SQL()
	}
	if g.WhereClause != nil {
		sql += " " + g.WhereClause.SQL()
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

type GqlEdgePattern interface {
	GqlElementPattern
	isGqlEdgePattern()
}

type GqlFullEdgeAny struct {
	PatternFilter *GqlPatternFilter
}

func (g GqlFullEdgeAny) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeAny) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeAny) SQL() string {
	return fmt.Sprintf("-[%v]-", g.PatternFilter.SQL())
}

func (g GqlFullEdgeAny) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeAny) isGqlElementPattern() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeAny) isGqlEdgePattern() {
	//TODO implement me
	panic("implement me")
}

type GqlFullEdgeLeft struct {
	PatternFilter *GqlPatternFilter
}

func (g GqlFullEdgeLeft) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeLeft) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeLeft) SQL() string {
	return fmt.Sprintf("<-[%v]-", g.PatternFilter.SQL())
}

func (g GqlFullEdgeLeft) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeLeft) isGqlElementPattern() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeLeft) isGqlEdgePattern() {
	//TODO implement me
	panic("implement me")
}

type GqlFullEdgeRight struct {
	PatternFilter *GqlPatternFilter
}

func (g GqlFullEdgeRight) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeRight) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeRight) SQL() string {
	return fmt.Sprintf("-[%v]->", g.PatternFilter.SQL())
}

func (g GqlFullEdgeRight) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeRight) isGqlElementPattern() {
	//TODO implement me
	panic("implement me")
}

func (g GqlFullEdgeRight) isGqlEdgePattern() {
	//TODO implement me
	panic("implement me")
}

type GqlAbbreviatedEdgeAny struct{}

func (g GqlAbbreviatedEdgeAny) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeAny) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeAny) SQL() string {
	return "-"
}

func (g GqlAbbreviatedEdgeAny) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeAny) isGqlElementPattern() {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeAny) isGqlEdgePattern() {
	//TODO implement me
	panic("implement me")
}

type GqlAbbreviatedEdgeLeft struct{}

func (g GqlAbbreviatedEdgeLeft) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeLeft) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeLeft) SQL() string {
	return "<-"
}

func (g GqlAbbreviatedEdgeLeft) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeLeft) isGqlElementPattern() {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeLeft) isGqlEdgePattern() {
	//TODO implement me
	panic("implement me")
}

type GqlAbbreviatedEdgeRight struct{}

func (g GqlAbbreviatedEdgeRight) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeRight) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlAbbreviatedEdgeRight) SQL() string {
	return "->"
}

func (g GqlAbbreviatedEdgeRight) isGqlPathTerm() {}

func (g GqlAbbreviatedEdgeRight) isGqlElementPattern() {}

func (g GqlAbbreviatedEdgeRight) isGqlEdgePattern() {}

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

type GqlPathModeEnum int

const (
	GqlPathModeWalk GqlPathModeEnum = iota
	GqlPathModeTrail
)

type GqlPathMode struct {
	Mode GqlPathModeEnum
}

func (g *GqlPathMode) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlPathMode) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlPathMode) SQL() string {
	switch g.Mode {
	case GqlPathModeTrail:
		return "TRAIL"
	case GqlPathModeWalk:
		return "WALK"
	default:
		panic("UNKNOWN GqlPathMode")
	}
}

type GqlSubpathPattern struct {
	LParen, RParen token.Pos
	PathMode       *GqlPathMode
	PathPattern    *GqlPathPattern
	WhereClause    *Where
}

func (g GqlSubpathPattern) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlSubpathPattern) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlSubpathPattern) SQL() string {
	sql := "("
	if g.PathMode != nil {
		sql += g.PathMode.SQL() + " "
	}
	sql += g.PathPattern.SQL()
	if g.WhereClause != nil {
		sql += " " + g.WhereClause.SQL()
	}
	return sql + ")"
}

func (g GqlSubpathPattern) isGqlPathTerm() {
	//TODO implement me
	panic("implement me")
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
	sql += fmt.Sprintf("(%v)", g.PatternFilter.SQL())
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
	GraphPatternVariable *Ident                 // optional
	IsLabelCondition     *GqlIsLabelCondition   // optional
	Filter               GqlPatternFilterFilter // optional
}

func (g GqlPatternFilter) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPatternFilter) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPatternFilter) SQL() string {
	var sql string
	if g.GraphPatternVariable != nil {
		sql += g.GraphPatternVariable.SQL()
	}
	if g.IsLabelCondition != nil {
		sql += g.IsLabelCondition.SQL()
	}
	if g.Filter != nil {
		if sql == "" {
			sql = g.Filter.SQL()
		} else {
			sql += " " + g.Filter.SQL()
		}
	}
	return sql
}

type GqlIsLabelCondition struct {
	IsOrColon       token.Pos
	LabelExpression GqlLabelExpression
}

func (g GqlIsLabelCondition) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlIsLabelCondition) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlIsLabelCondition) SQL() string {
	return ":" + g.LabelExpression.SQL()
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

type GqlLabelName struct {
	IsPercent bool
	LabelName *Ident
}

func (g GqlLabelName) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelName) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelName) SQL() string {
	if g.IsPercent {
		return "%"
	}
	return g.LabelName.SQL()
}

func (g GqlLabelName) isGqlLabelExpression() {
	//TODO implement me
	panic("implement me")
}

type GqlPatternFilterFilter interface {
	Node
	isGqlPatternFilterFilter()
}

type GqlPropertyFilters struct {
	LBrace                 token.Pos
	PropertyFilterElemList []*GqlElementProperty
	RBrace                 token.Pos
}

func (g GqlPropertyFilters) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPropertyFilters) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPropertyFilters) SQL() string {
	var elemSqls []string
	for _, elem := range g.PropertyFilterElemList {
		elemSqls = append(elemSqls, elem.SQL())
	}

	return fmt.Sprintf("{%v}", strings.Join(elemSqls, ", "))
}

func (g GqlPropertyFilters) isGqlPatternFilterFilter() {}

type GqlElementProperty struct {
	ElementPropertyName  *Ident
	ElementPropertyValue Expr
}

func (g GqlElementProperty) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlElementProperty) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlElementProperty) SQL() string {
	var sql string
	sql += g.ElementPropertyName.SQL() + ": " + g.ElementPropertyValue.SQL()
	return sql
}
