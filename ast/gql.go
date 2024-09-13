package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
	"strings"
)

type GqlGraphQuery struct {
	GraphClause                  *GqlGraphClause
	GqlMultiLinearQueryStatement *GqlMultiLinearQueryStatement
}

func (g GqlGraphQuery) Pos() token.Pos {
	return g.GraphClause.Pos()
}

func (g GqlGraphQuery) End() token.Pos {
	return g.GqlMultiLinearQueryStatement.End()
}

func (g GqlGraphQuery) isStatement() {}
func (s *GqlGraphQuery) SQL() string {
	return s.GraphClause.SQL() + "\n" + s.GqlMultiLinearQueryStatement.SQL()
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

func lastElem[T any](s []T) T {
	return s[len(s)-1]
}

func firstPos[T Node](s []T) token.Pos {
	return s[0].Pos()
}

func lastEnd[T Node](s []T) token.Pos {
	return lastElem(s).End()
}

func (s *GqlSimpleLinearQueryStatement) End() token.Pos {
	return lastEnd(s.PrimitiveQueryStatementList)
}

func sqlJoin[T Node](elems []T, sep string) string {
	var b strings.Builder
	for i, r := range elems {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(r.SQL())
	}
	return b.String()
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
	if g.Optional >= 0 {
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

type GqlReturnStatement struct {
	AllOrDistinct        GqlAllOrDistinctEnum
	Return               token.Pos
	ReturnItemList       []SelectItem
	GroupByClause        *GroupBy
	OrderByClause        *OrderBy
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
