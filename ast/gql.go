package ast

import (
	"fmt"
	"github.com/cloudspannerecosystem/memefish/token"
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
	sql := s.GqlLinearQueryStatements[0].SQL()
	for _, r := range s.GqlLinearQueryStatements[1:] {
		sql += "\n" + r.SQL()
	}
	return sql
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

func (*GqlSimpleLinearQueryStatement) isGqlLinearQueryStatement() {}
func (s *GqlSimpleLinearQueryStatement) SQL() string {
	// TODO: list should not be empty
	if len(s.PrimitiveQueryStatementList) == 0 {
		return ""
	}
	sql := s.PrimitiveQueryStatementList[0].SQL()
	for _, r := range s.PrimitiveQueryStatementList[1:] {
		sql += "\n" + r.SQL()
	}
	return sql
}

type GqlCompositeLinearQueryStatement struct{}

func (*GqlCompositeLinearQueryStatement) isGqlLinearQueryStatement() {}

type GqlPrimitiveQueryStatement interface {
	Node
	isGqlPrimitiveQueryStatement()
}
type GqlReturnStatement struct {
	Return         token.Pos
	ReturnItemList []SelectItem
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
	sql := fmt.Sprintf("LET %v", s.LinearGraphVariableList[0].SQL())
	for _, v := range s.LinearGraphVariableList[1:] {
		sql += ", " + v.SQL()
	}
	return sql
}
func (*GqlLetStatement) isGqlPrimitiveQueryStatement() {}

func (*GqlReturnStatement) isGqlPrimitiveQueryStatement() {}

func (s *GqlReturnStatement) SQL() string {
	sql := "RETURN "
	sql += s.ReturnItemList[0].SQL()
	for _, r := range s.ReturnItemList[1:] {
		sql += ", " + r.SQL()
	}
	return sql
}
