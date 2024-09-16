package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
	"strings"
)

// https://cloud.google.com/spanner/docs/reference/standard-sql/graph-patterns

// GqlGraphPattern represents is the toplevel node of GQL graph patterns.
type GqlGraphPattern struct {
	// pos = GqlTopLevelPathPattern[0].pos
	// end = (WhereClause ?? PathPatternList[$]).end
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
	return sqlJoin(g.PathPatternList, ", ") + sqlOpt(" ", g.WhereClause, "")
}

// GqlTopLevelPathPattern is a PathPattern optionally prefixed by PathSearchPrefixOrPathMode.
type GqlTopLevelPathPattern struct {
	// pos = (PathSearchPrefixOrPathMode ?? PathPattern).pos
	// end = PathPattern.end
	PathSearchPrefixOrPathMode GqlPathSearchPrefixOrPathMode // optional
	PathPattern                *GqlPathPattern
}

func (g GqlTopLevelPathPattern) Pos() token.Pos {
	return firstValidPos(g.PathSearchPrefixOrPathMode, g.PathPattern)
}

func (g GqlTopLevelPathPattern) End() token.Pos {
	return g.PathPattern.End()
}

func (g GqlTopLevelPathPattern) SQL() string {
	return sqlOpt("", g.PathSearchPrefixOrPathMode, " ") + g.PathPattern.SQL()
}

// GqlPathSearchPrefixOrPathMode represents `{ path_search_prefix | path_mode }`
type GqlPathSearchPrefixOrPathMode interface {
	Node
	isGqlPathSearchPrefixOrPathMode()
}

// GqlEdgePattern represents edge pattern nodes.
//
//	edge_pattern:
//	 {
//	   full_edge_any |
//	   full_edge_left |
//	   full_edge_right |
//	   abbreviated_edge_any |
//	   abbreviated_edge_left |
//	   abbreviated_edge_right
//	 }
type GqlEdgePattern interface {
	GqlElementPattern
	isGqlEdgePattern()
}

// GqlFullEdgeAny is node representing`-[pattern_filler]-` .
type GqlFullEdgeAny struct {
	// pos = First.pos
	// end = Last.pos + 1
	First, Last   token.Pos
	PatternFiller *GqlPatternFiller
}

func (g GqlFullEdgeAny) Pos() token.Pos {
	return g.First
}

func (g GqlFullEdgeAny) End() token.Pos {
	return g.Last + 1
}

func (g GqlFullEdgeAny) SQL() string {
	return fmt.Sprintf("-[%v]-", g.PatternFiller.SQL())
}

func (g GqlFullEdgeAny) isGqlPathTerm()       {}
func (g GqlFullEdgeAny) isGqlElementPattern() {}
func (g GqlFullEdgeAny) isGqlEdgePattern()    {}

type GqlFullEdgeLeft struct {
	// pos = First
	// end = Last + 1
	First         token.Pos // position of "<"
	Last          token.Pos // position of the last "-"
	PatternFiller *GqlPatternFiller
}

func (g GqlFullEdgeLeft) Pos() token.Pos {
	return g.First
}

func (g GqlFullEdgeLeft) End() token.Pos {
	return g.Last + 1
}

func (g GqlFullEdgeLeft) SQL() string {
	return fmt.Sprintf("<-[%v]-", g.PatternFiller.SQL())
}

func (g GqlFullEdgeLeft) isGqlPathTerm() {}

func (g GqlFullEdgeLeft) isGqlElementPattern() {}

func (g GqlFullEdgeLeft) isGqlEdgePattern() {}

type GqlFullEdgeRight struct {
	// pos = First
	// end = Last + 1
	First         token.Pos // position of the first "-"
	Last          token.Pos // position of ">"
	PatternFiller *GqlPatternFiller
}

func (g GqlFullEdgeRight) Pos() token.Pos {
	return g.First
}

func (g GqlFullEdgeRight) End() token.Pos {
	return g.Last + 1
}

func (g GqlFullEdgeRight) SQL() string {
	return fmt.Sprintf("-[%v]->", g.PatternFiller.SQL())
}

func (g GqlFullEdgeRight) isGqlPathTerm() {}

func (g GqlFullEdgeRight) isGqlElementPattern() {}

func (g GqlFullEdgeRight) isGqlEdgePattern() {}

type GqlAbbreviatedEdgeAny struct {
	// pos = Hyphen
	// end = Hyphen +1
	Hyphen token.Pos // position of "-"
}

func (g GqlAbbreviatedEdgeAny) Pos() token.Pos {
	return g.Hyphen
}

func (g GqlAbbreviatedEdgeAny) End() token.Pos {
	return g.Hyphen + 1
}

func (g GqlAbbreviatedEdgeAny) SQL() string {
	return "-"
}

func (g GqlAbbreviatedEdgeAny) isGqlPathTerm() {}

func (g GqlAbbreviatedEdgeAny) isGqlElementPattern() {}

func (g GqlAbbreviatedEdgeAny) isGqlEdgePattern() {}

type GqlAbbreviatedEdgeLeft struct {
	// pos = First
	// end = Last + 1
	First token.Pos // position of "<"
	Last  token.Pos // position of "-"
}

func (g GqlAbbreviatedEdgeLeft) Pos() token.Pos {
	return g.First
}

func (g GqlAbbreviatedEdgeLeft) End() token.Pos {
	return g.Last + 1
}

func (g GqlAbbreviatedEdgeLeft) SQL() string {
	return "<-"
}

func (g GqlAbbreviatedEdgeLeft) isGqlPathTerm() {}

func (g GqlAbbreviatedEdgeLeft) isGqlElementPattern() {}

func (g GqlAbbreviatedEdgeLeft) isGqlEdgePattern() {}

type GqlAbbreviatedEdgeRight struct {
	// pos = First
	// end = Last + 1
	First token.Pos // position of "-"
	Last  token.Pos // position of ">"
}

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

type GqlQuantifiablePathTerm struct {
	Hint       *Hint // optional
	PathTerm   GqlPathTerm
	Quantifier GqlQuantifier // optional
}

func (g GqlQuantifiablePathTerm) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlQuantifiablePathTerm) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlQuantifiablePathTerm) SQL() string {
	var sql string
	sql += sqlOpt("", g.Hint, "")
	sql += g.PathTerm.SQL()
	if g.Quantifier != nil {
		sql += g.Quantifier.SQL()
	}
	return sql
}

type GqlPathPattern struct {
	PathTermList []*GqlQuantifiablePathTerm
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

func (g GqlWhereClause) isGqlPatternFillerFilter() {}

type GqlElementPattern interface {
	Node
	GqlPathTerm
	isGqlElementPattern()
}

type gqlPathModeEnum int

const (
	GqlPathModeWalk gqlPathModeEnum = iota
	GqlPathModeTrail
)

type GqlPathMode struct {
	Mode            gqlPathModeEnum
	ModeToken       *Ident
	PathOrPathToken *Ident
}

func (g *GqlPathMode) isGqlPathSearchPrefixOrPathMode() {
	//TODO implement me
	panic("implement me")
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

/*
	type GqlQuantifiedPathPrimary struct {
		PathPrimary
	}
*/
type GqlQuantifier interface {
	Node
	isGqlQuantifier()
}
type GqlFixedQuantifier struct {
	LBrace, RBrace token.Pos
	Bound          IntValue
}

func (g GqlFixedQuantifier) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFixedQuantifier) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlFixedQuantifier) SQL() string {
	return "{" + g.Bound.SQL() + "}"
}

func (g GqlFixedQuantifier) isGqlQuantifier() {
	//TODO implement me
	panic("implement me")
}

type GqlBoundedQuantifier struct {
	// pos = LBrace
	// end = RBrace + 1
	LBrace, RBrace token.Pos
	LowerBound     IntValue // optional
	UpperBound     IntValue
}

func (g *GqlBoundedQuantifier) Pos() token.Pos {
	return g.LBrace
}

func (g *GqlBoundedQuantifier) End() token.Pos {
	return g.RBrace + 1
}

func (g *GqlBoundedQuantifier) SQL() string {
	return fmt.Sprintf("{%v,%v}", sqlOpt("", g.LowerBound, ""), g.UpperBound.SQL())
}

func (g *GqlBoundedQuantifier) isGqlQuantifier() {
}

type GqlSubpathPattern struct {
	// pos = LParen
	// end = RParen + 1
	LParen, RParen token.Pos
	PathMode       *GqlPathMode
	PathPattern    *GqlPathPattern
	WhereClause    *Where
}

func (g GqlSubpathPattern) Pos() token.Pos {
	return g.LParen
}

func (g GqlSubpathPattern) End() token.Pos {
	return g.RParen + 1
}

func (g GqlSubpathPattern) SQL() string {
	return "(" +
		sqlOpt("", g.PathMode, " ") +
		g.PathPattern.SQL() +
		sqlOpt(" ", g.WhereClause, "") +
		")"
}

func (g GqlSubpathPattern) isGqlPathTerm() {
}

type GqlNodePattern struct {
	LParen, RParen token.Pos
	PatternFiller  *GqlPatternFiller
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
	sql += fmt.Sprintf("(%v)", g.PatternFiller.SQL())
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
/*
type EdgePattern interface {
	Node
	isEdgePattern()
}

*/

type GqlPatternFiller struct {
	// Hint is graph element hint which is a table hint.
	Hint                 *Hint
	GraphPatternVariable *Ident                 // optional
	IsLabelCondition     *GqlIsLabelCondition   // optional
	Filter               GqlPatternFillerFilter // optional
}

func (g GqlPatternFiller) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPatternFiller) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlPatternFiller) SQL() string {
	var sql string
	sql += sqlOpt("", g.Hint, "")
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

type GqlLabelOrExpression struct {
	Left, Right GqlLabelExpression
}

type GqlLabelParenExpression struct {
	LParen, RParen token.Pos
	LabelExpr      GqlLabelExpression
}

func (g GqlLabelParenExpression) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelParenExpression) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelParenExpression) SQL() string {
	return "(" + g.LabelExpr.SQL() + ")"
}

func (g GqlLabelParenExpression) isGqlLabelExpression() {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelOrExpression) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelOrExpression) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelOrExpression) SQL() string {
	return g.Left.SQL() + "|" + g.Right.SQL()
}

func (g GqlLabelOrExpression) isGqlLabelExpression() {
	//TODO implement me
	panic("implement me")
}

type GqlLabelAndExpression struct {
	Left, Right GqlLabelExpression
}

func (g GqlLabelAndExpression) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelAndExpression) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelAndExpression) SQL() string {
	return g.Left.SQL() + "&" + g.Right.SQL()
}

func (g GqlLabelAndExpression) isGqlLabelExpression() {
	//TODO implement me
	panic("implement me")
}

type GqlLabelNotExpression struct {
	NotPos          token.Pos
	LabelExpression GqlLabelExpression
}

func (g GqlLabelNotExpression) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelNotExpression) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLabelNotExpression) SQL() string {
	return "!" + g.LabelExpression.SQL()
}

func (g GqlLabelNotExpression) isGqlLabelExpression() {
	//TODO implement me
	panic("implement me")
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

func (g GqlLabelName) isGqlLabelExpression() {}

type GqlPatternFillerFilter interface {
	Node
	isGqlPatternFillerFilter()
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

func (g GqlPropertyFilters) isGqlPatternFillerFilter() {}

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

type gqlSearchPrefixEnum int

const (
	GqlPathSearchPrefixAll = iota
	GqlPathSearchPrefixAny
	GqlPathSearchPrefixAnyShortest
)

type GqlPathSearchPrefix struct {
	StartPos     token.Pos
	SearchPrefix gqlSearchPrefixEnum
}

func (p *GqlPathSearchPrefix) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (p *GqlPathSearchPrefix) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (p *GqlPathSearchPrefix) SQL() string {
	switch p.SearchPrefix {
	case GqlPathSearchPrefixAny:
		return "ANY"
	case GqlPathSearchPrefixAnyShortest:
		return "ANY SHORTEST"
	case GqlPathSearchPrefixAll:
		return "ALL"
	default:
		panic("invalid GqlPathSearchPrefix")
	}
}

func (p *GqlPathSearchPrefix) isGqlPathSearchPrefixOrPathMode() {}
