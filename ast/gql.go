package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
)

// GqlGraphQuery is toplevel node of GRAPH query.
type GqlGraphQuery struct {
	// pos = (GraphClause ?? GqlMultiLinearQueryStatement).pos
	// end = GqlMultiLinearQueryStatement.pos
	GraphClause                  *GqlGraphClause
	GqlMultiLinearQueryStatement *GqlMultiLinearQueryStatement
}

func (g *GqlGraphQuery) Pos() token.Pos {
	return g.GraphClause.Pos()
}

func (g *GqlGraphQuery) End() token.Pos {
	return g.GqlMultiLinearQueryStatement.End()
}

func (g *GqlGraphQuery) isStatement() {}
func (s *GqlGraphQuery) SQL() string {
	return s.GraphClause.SQL() + "\n" + s.GqlMultiLinearQueryStatement.SQL()
}

// GqlQueryExpr is similar to GqlGraphQuery,
// but it is appeared in GQL subqueries and it can optionally have GRAPH clause
type GqlQueryExpr struct {
	// pos = (GraphClause ?? GqlMultiLinearQueryStatement).pos
	// end = GqlMultiLinearQueryStatement.pos
	GraphClause                  *GqlGraphClause // optional
	GqlMultiLinearQueryStatement *GqlMultiLinearQueryStatement
}

func (g *GqlQueryExpr) Pos() token.Pos {
	return firstValidPos(g.GraphClause, g.GqlMultiLinearQueryStatement)
}

func (g *GqlQueryExpr) End() token.Pos {
	return g.GqlMultiLinearQueryStatement.End()
}

func (g *GqlQueryExpr) SQL() string {
	return sqlOpt("", g.GraphClause, "\n") + g.GqlMultiLinearQueryStatement.SQL()
}

type GqlGraphClause struct {
	GqlGraph             token.Pos
	GqlPropertyGraphName *Ident
}

func (s *GqlGraphClause) Pos() token.Pos {
	return s.GqlGraph
}

func (s *GqlGraphClause) End() token.Pos {
	return s.GqlPropertyGraphName.End()
}

func (s *GqlGraphClause) SQL() string {
	return fmt.Sprintf("GRAPH %v", s.GqlPropertyGraphName.SQL())
}

type GqlNextStatement struct {
	GqlNext *token.Pos
}

func (s *GqlNextStatement) SQL() string {
	return "NEXT"
}

type GqlMultiLinearQueryStatement struct {
	GqlLinearQueryStatements []GqlLinearQueryStatement
}

func (s *GqlMultiLinearQueryStatement) Pos() token.Pos {
	return s.GqlLinearQueryStatements[0].Pos()
}

func (s *GqlMultiLinearQueryStatement) End() token.Pos {
	return s.GqlLinearQueryStatements[len(s.GqlLinearQueryStatements)-1].End()
}

func (s *GqlMultiLinearQueryStatement) SQL() string {
	return sqlJoin(s.GqlLinearQueryStatements, "\nNEXT\n")
}

type GqlLinearQueryStatement interface {
	Node
	isGqlLinearQueryStatement()
}

type GqlSimpleLinearQueryStatement struct {
	PrimitiveQueryStatementList []GqlPrimitiveQueryStatement
}

func (s *GqlSimpleLinearQueryStatement) Pos() token.Pos {
	return firstPos(s.PrimitiveQueryStatementList)
}

func (s *GqlSimpleLinearQueryStatement) End() token.Pos {
	return lastEnd(s.PrimitiveQueryStatementList)
}

func (*GqlSimpleLinearQueryStatement) isGqlLinearQueryStatement() {}

func (s *GqlSimpleLinearQueryStatement) SQL() string {
	return sqlJoin(s.PrimitiveQueryStatementList, "\n")
}

type GqlAllOrDistinctEnum int

const (
	GqlAllOrDistinctImplicitAll = iota
	GqlAllOrDistinctAll         = iota
	GqlAllOrDistinctDistinct
)

type GqlSetOperatorEnum int

const (
	GqlSetOperatorUnion = iota
	GqlSetOperatorIntersect
	GqlSetOperatorExcept
)

type GqlSimpleLinearQueryStatementWithSetOperator struct {
	StartPos      token.Pos
	SetOperator   GqlSetOperatorEnum
	DistinctOrAll GqlAllOrDistinctEnum
	Statement     *GqlSimpleLinearQueryStatement
}

type GqlCompositeLinearQueryStatement struct {
	HeadSimpleLinearQueryStatement     *GqlSimpleLinearQueryStatement
	TailSimpleLinearQueryStatementList []*GqlSimpleLinearQueryStatementWithSetOperator
}

func (s *GqlCompositeLinearQueryStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (s *GqlCompositeLinearQueryStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (s *GqlCompositeLinearQueryStatement) SQL() string {
	var sql string
	sql += s.HeadSimpleLinearQueryStatement.SQL()
	for _, s := range s.TailSimpleLinearQueryStatementList {
		sql += "\n"
		switch s.SetOperator {
		case GqlSetOperatorUnion:
			sql += "UNION"
		case GqlSetOperatorIntersect:
			sql += "INTERSECT"
		case GqlSetOperatorExcept:
			sql += "INTERSECT"
		}

		switch s.DistinctOrAll {
		case GqlAllOrDistinctAll:
			sql += " " + "ALL"
		case GqlAllOrDistinctDistinct:
			sql += " " + "DISTINCT"
		}
		sql += "\n" + s.Statement.SQL()
	}
	return sql
}

func (*GqlCompositeLinearQueryStatement) isGqlLinearQueryStatement() {}

type GqlPrimitiveQueryStatement interface {
	Node
	isGqlPrimitiveQueryStatement()
}

type GqlMatchStatement struct {
	Optional     token.Pos
	Match        token.Pos
	MatchHint    *Hint
	PrefixOrMode GqlPathSearchPrefixOrPathMode // optional
	GraphPattern *GqlGraphPattern
}

func (g GqlMatchStatement) Pos() token.Pos {
	return g.Match
}

func (g GqlMatchStatement) End() token.Pos {
	return g.GraphPattern.Pos()
}

func (g GqlMatchStatement) SQL() string {
	var sql string
	if g.Optional != token.InvalidPos {
		sql = "OPTIONAL MATCH"
	} else {
		sql = "MATCH"
	}
	if g.MatchHint != nil {
		sql += g.MatchHint.SQL()
	}
	if g.PrefixOrMode != nil {
		sql += " " + g.PrefixOrMode.SQL()
	}
	sql += " " + g.GraphPattern.SQL()
	return sql
}

func (g GqlMatchStatement) isGqlPrimitiveQueryStatement() {}

type GqlLimitAndOffsetClause interface {
	Node
	isGqlLimitAndOffsetClause()
}

// GqlLimitClause is wrapper of Limit for GQL
type GqlLimitClause struct {
	Limit *Limit
}

// GqlOffsetClause is wrapper of Offset for GQL
type GqlOffsetClause struct {
	Offset *Offset
}

func (g GqlOffsetClause) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOffsetClause) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOffsetClause) SQL() string {
	//TODO implement me
	panic("implement me")
}

func (g GqlOffsetClause) isGqlLimitAndOffsetClause() {
	//TODO implement me
	panic("implement me")
}

// GqlLimitClauseWithOffset
type GqlLimitWithOffsetClause struct {
	Limit  *Limit
	Offset *Offset
}

func (g GqlLimitWithOffsetClause) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitWithOffsetClause) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitWithOffsetClause) SQL() string {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitWithOffsetClause) isGqlLimitAndOffsetClause() {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitClause) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitClause) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitClause) SQL() string {
	//TODO implement me
	panic("implement me")
}

func (g GqlLimitClause) isGqlLimitAndOffsetClause() {
	//TODO implement me
	panic("implement me")
}

type GqlFilterStatement struct {
	Filter token.Pos
	Where  token.Pos
	Expr   Expr
}

func (g *GqlFilterStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlFilterStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlFilterStatement) SQL() string {
	if g.Where != token.InvalidPos {
		return "FILTER WHERE " + g.Expr.SQL()
	}
	return "FILTER " + g.Expr.SQL()
}

func (g *GqlFilterStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

type GqlForStatement struct {
	For              token.Pos
	ElementName      *Ident
	ArrayExpression  Expr
	WithOffsetClause *GqlWithOffsetClause
}

func (g GqlForStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlForStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlForStatement) SQL() string {
	return "FOR " + g.ElementName.SQL() + " IN " + g.ArrayExpression.SQL() + sqlOpt(" ", g.WithOffsetClause, "")
}

func (g GqlForStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

type GqlLimitStatement struct {
	Limit token.Pos
	Count IntValue
}

func (g *GqlLimitStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlLimitStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g *GqlLimitStatement) SQL() string {
	return "LIMIT " + g.Count.SQL()
}

func (g *GqlLimitStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

/* OFFSET/SKIP statement */

// GqlOffsetStatement represents OFFSET statement.
// It also represents SKIP statement as the synonym.
type GqlOffsetStatement struct {
	// pos = Offset.pos
	// end = Count.end
	Offset token.Pos
	IsSkip bool
	Count  IntValue
}

func (g GqlOffsetStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOffsetStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOffsetStatement) SQL() string {
	if g.IsSkip {
		return "SKIP " + g.Count.SQL()
	}
	return "OFFSET " + g.Count.SQL()
}

func (g GqlOffsetStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

type GqlOrderByStatement struct {
	Order                    token.Pos
	OrderBySpecificationList []*GqlOrderBySpecification
}

func (g GqlOrderByStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOrderByStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOrderByStatement) SQL() string {
	return "ORDER BY " + sqlJoin(g.OrderBySpecificationList, ", ")
}

func (g GqlOrderByStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

type GqlDirection string

const (
	GqlSortOrderUnspecified = ""
	GqlSortOrderAsc         = "ASC"
	GqlSortOrderAscending   = "ASCENDING"
	GqlSortOrderDesc        = "DESC"
	GqlSortOrderDescending  = "DESCENDING"
)

// TODO
type GqlOrderBySpecification struct {
	Expr Expr
	// TODO
	// CollationSpecification *GqlCollationSpecification
	DirectionPos token.Pos
	Direction    GqlDirection
}

func (g GqlOrderBySpecification) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOrderBySpecification) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlOrderBySpecification) SQL() string {
	sql := g.Expr.SQL()
	if g.Direction != GqlSortOrderUnspecified {
		sql += " " + string(g.Direction)
	}
	return sql
}

type GqlWithStatement struct {
	With           token.Pos
	AllOrDistinct  GqlAllOrDistinctEnum
	ReturnItemList []GqlReturnItem
	GroupByClause  *GroupBy // optional
}

func (g GqlWithStatement) Pos() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlWithStatement) End() token.Pos {
	//TODO implement me
	panic("implement me")
}

func (g GqlWithStatement) SQL() string {
	var allOrDistinctStr string
	switch g.AllOrDistinct {
	case GqlAllOrDistinctDistinct:
		allOrDistinctStr = "DISTINCT "
	case GqlAllOrDistinctAll:
		allOrDistinctStr = "ALL "
	case GqlAllOrDistinctImplicitAll:
		allOrDistinctStr = ""
	}
	return "WITH " + allOrDistinctStr + sqlJoin(g.ReturnItemList, ", ") + sqlOpt(" ", g.GroupByClause, "")
}

func (g GqlWithStatement) isGqlPrimitiveQueryStatement() {
	//TODO implement me
	panic("implement me")
}

type GqlWithOffsetClause struct {
	// pos = With.pos
	// end = OffsetName.end ?? Offset + 6
	With       token.Pos
	Offset     token.Pos
	OffsetName *Ident
}

func (g *GqlWithOffsetClause) Pos() token.Pos {
	return g.With
}

func (g *GqlWithOffsetClause) End() token.Pos {
	if g.OffsetName != nil {
		return g.OffsetName.End()
	}
	return g.Offset + 6
}

func (g *GqlWithOffsetClause) SQL() string {
	return "WITH OFFSET" + sqlOpt(" AS ", g.OffsetName, "")
}

// GqlReturnItem is similar to SelectItem,
// but it don't permit DotStar and AsAlias without AS.
type GqlReturnItem SelectItem

type GqlReturnStatement struct {
	AllOrDistinct  GqlAllOrDistinctEnum
	Return         token.Pos
	ReturnItemList []GqlReturnItem

	// Use GoogleSQL GroupBy because it is referenced in docs
	GroupByClause *GroupBy //optional

	// Use GoogleSQL OrderBy because it is referenced in docs
	OrderByClause *OrderBy //optional

	LimitAndOffsetClause GqlLimitAndOffsetClause
}

func (s *GqlReturnStatement) Pos() token.Pos {
	return s.Return
}

func (s *GqlReturnStatement) End() token.Pos {
	return s.ReturnItemList[len(s.ReturnItemList)-1].End()
}

type GqlLinearGraphVariable struct {
	VariableName *Ident
	Value        Expr
}

func (s *GqlLinearGraphVariable) Pos() token.Pos { return s.VariableName.Pos() }

func (s *GqlLinearGraphVariable) End() token.Pos { return s.Value.End() }

func (s *GqlLinearGraphVariable) SQL() string {
	return fmt.Sprintf("%v = %v", s.VariableName.SQL(), s.Value.SQL())
}

type GqlLetStatement struct {
	// pos = Let.pos
	// end = LinearGraphVariableList[$].end

	Let                     token.Pos
	LinearGraphVariableList []*GqlLinearGraphVariable
}

func (s *GqlLetStatement) Pos() token.Pos {
	return s.Let
}

func (s *GqlLetStatement) End() token.Pos {
	if len(s.LinearGraphVariableList) == 0 {
		return s.Let + 3
	}
	return s.LinearGraphVariableList[len(s.LinearGraphVariableList)].End()
}

func (s *GqlLetStatement) SQL() string {
	return "LET " + sqlJoin(s.LinearGraphVariableList, ", ")
}
func (*GqlLetStatement) isGqlPrimitiveQueryStatement() {}

func (*GqlReturnStatement) isGqlPrimitiveQueryStatement() {}

func (s *GqlReturnStatement) SQL() string {
	sql := "RETURN "
	switch s.AllOrDistinct {
	case GqlAllOrDistinctAll:
		sql += "ALL "
	case GqlAllOrDistinctDistinct:
		sql += "DISTINCT "
	case GqlAllOrDistinctImplicitAll:
		// empty
	}

	sql += sqlJoin(s.ReturnItemList, ", ")

	if s.GroupByClause != nil {
		sql += " " + s.GroupByClause.SQL()
	}
	if s.OrderByClause != nil {
		sql += " " + s.OrderByClause.SQL()
	}
	if s.LimitAndOffsetClause != nil {
		sql += " " + s.LimitAndOffsetClause.SQL()
	}
	return sql
}
