package ast

import (
	"testing"
)

func TestStatement(t *testing.T) {
	Statement(&QueryStatement{}).isStatement()
	Statement(&CreateDatabase{}).isStatement()
	Statement(&CreateTable{}).isStatement()
	Statement(&CreateIndex{}).isStatement()
	Statement(&AlterTable{}).isStatement()
	Statement(&DropTable{}).isStatement()
	Statement(&DropIndex{}).isStatement()
	Statement(&Insert{}).isStatement()
	Statement(&Update{}).isStatement()
	Statement(&Delete{}).isStatement()
}

func TestSelectItem(t *testing.T) {
	SelectItem(&Star{}).isSelectItem()
	SelectItem(&DotStar{}).isSelectItem()
	SelectItem(&Alias{}).isSelectItem()
	SelectItem(&ExprSelectItem{}).isSelectItem()
}

func TestQueryExpr(t *testing.T) {
	QueryExpr(&Select{}).isQueryExpr()
	QueryExpr(&SubQuery{}).isQueryExpr()
	QueryExpr(&CompoundQuery{}).isQueryExpr()
}

func TestTableExpr(t *testing.T) {
	TableExpr(&Unnest{}).isTableExpr()
	TableExpr(&TableName{}).isTableExpr()
	TableExpr(&SubQueryTableExpr{}).isTableExpr()
	TableExpr(&ParenTableExpr{}).isTableExpr()
	TableExpr(&Join{}).isTableExpr()
}

func TestJoinCondition(t *testing.T) {
	JoinCondition(&On{}).isJoinCondition()
	JoinCondition(&Using{}).isJoinCondition()
}

func TestExpr(t *testing.T) {
	Expr(&BinaryExpr{}).isExpr()
	Expr(&UnaryExpr{}).isExpr()
	Expr(&InExpr{}).isExpr()
	Expr(&IsNullExpr{}).isExpr()
	Expr(&IsBoolExpr{}).isExpr()
	Expr(&BetweenExpr{}).isExpr()
	Expr(&SelectorExpr{}).isExpr()
	Expr(&IndexExpr{}).isExpr()
	Expr(&CallExpr{}).isExpr()
	Expr(&CountStarExpr{}).isExpr()
	Expr(&CastExpr{}).isExpr()
	Expr(&ExtractExpr{}).isExpr()
	Expr(&CaseExpr{}).isExpr()
	Expr(&ParenExpr{}).isExpr()
	Expr(&ScalarSubQuery{}).isExpr()
	Expr(&ArraySubQuery{}).isExpr()
	Expr(&ExistsSubQuery{}).isExpr()
	Expr(&Param{}).isExpr()
	Expr(&Ident{}).isExpr()
	Expr(&Path{}).isExpr()
	Expr(&ArrayLiteral{}).isExpr()
	Expr(&StructLiteral{}).isExpr()
	Expr(&NullLiteral{}).isExpr()
	Expr(&BoolLiteral{}).isExpr()
	Expr(&IntLiteral{}).isExpr()
	Expr(&FloatLiteral{}).isExpr()
	Expr(&StringLiteral{}).isExpr()
	Expr(&BytesLiteral{}).isExpr()
	Expr(&DateLiteral{}).isExpr()
	Expr(&TimestampLiteral{}).isExpr()
	Expr(&NumericLiteral{}).isExpr()
}

func TestInCondition(t *testing.T) {
	InCondition(&ValuesInCondition{}).isInCondition()
	InCondition(&UnnestInCondition{}).isInCondition()
	InCondition(&SubQueryInCondition{}).isInCondition()
}

func TestType(t *testing.T) {
	Type(&SimpleType{}).isType()
	Type(&ArrayType{}).isType()
	Type(&StructType{}).isType()
}

func TestIntValue(t *testing.T) {
	IntValue(&IntLiteral{}).isIntValue()
	IntValue(&Param{}).isIntValue()
	IntValue(&CastIntValue{}).isIntValue()
}

func TestNumValue(t *testing.T) {
	NumValue(&IntLiteral{}).isNumValue()
	NumValue(&FloatLiteral{}).isNumValue()
	NumValue(&Param{}).isNumValue()
	NumValue(&CastNumValue{}).isNumValue()
}

func TestStringValue(t *testing.T) {
	StringValue(&StringLiteral{}).isStringValue()
	StringValue(&Param{}).isStringValue()
}

func TestDDL(t *testing.T) {
	DDL(&CreateDatabase{}).isDDL()
	DDL(&CreateTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&AlterTable{}).isDDL()
	DDL(&DropTable{}).isDDL()
	DDL(&DropIndex{}).isDDL()
}

func TestTableAlternation(t *testing.T) {
	TableAlternation(&AddColumn{}).isTableAlternation()
	TableAlternation(&AddForeignKey{}).isTableAlternation()
	TableAlternation(&DropColumn{}).isTableAlternation()
	TableAlternation(&DropConstraint{}).isTableAlternation()
	TableAlternation(&SetOnDelete{}).isTableAlternation()
	TableAlternation(&AlterColumn{}).isTableAlternation()
	TableAlternation(&AlterColumnSet{}).isTableAlternation()
}

func TestSchemaType(t *testing.T) {
	SchemaType(&ScalarSchemaType{}).isSchemaType()
	SchemaType(&SizedSchemaType{}).isSchemaType()
	SchemaType(&ArraySchemaType{}).isSchemaType()
}

func TestDML(t *testing.T) {
	DML(&Insert{}).isDML()
	DML(&Delete{}).isDML()
	DML(&Update{}).isDML()
}

func TestInsertInput(t *testing.T) {
	InsertInput(&ValuesInput{}).isInsertInput()
	InsertInput(&SubQueryInput{}).isInsertInput()
}
