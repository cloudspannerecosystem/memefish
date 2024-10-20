package ast

import (
	"testing"
)

func TestStatement(t *testing.T) {
	Statement(&QueryStatement{}).isStatement()
	Statement(&CreateDatabase{}).isStatement()
	Statement(&CreateTable{}).isStatement()
	Statement(&CreateView{}).isStatement()
	Statement(&CreateIndex{}).isStatement()
	Statement(&CreateSequence{}).isStatement()
	Statement(&CreateRole{}).isStatement()
	Statement(&AlterTable{}).isStatement()
	Statement(&AlterIndex{}).isStatement()
	Statement(&DropTable{}).isStatement()
	Statement(&DropIndex{}).isStatement()
	Statement(&DropVectorIndex{}).isStatement()
	Statement(&DropRole{}).isStatement()
	Statement(&Insert{}).isStatement()
	Statement(&Delete{}).isStatement()
	Statement(&Update{}).isStatement()
	Statement(&Grant{}).isStatement()
	Statement(&Revoke{}).isStatement()
}

func TestQueryExpr(t *testing.T) {
	QueryExpr(&Select{}).isQueryExpr()
	QueryExpr(&SubQuery{}).isQueryExpr()
	QueryExpr(&CompoundQuery{}).isQueryExpr()
}

func TestSelectItem(t *testing.T) {
	SelectItem(&Star{}).isSelectItem()
	SelectItem(&DotStar{}).isSelectItem()
	SelectItem(&Alias{}).isSelectItem()
	SelectItem(&ExprSelectItem{}).isSelectItem()
}

func TestSelectAs(t *testing.T) {
	SelectAs(&AsStruct{}).isSelectAs()
	SelectAs(&AsValue{}).isSelectAs()
	SelectAs(&AsTypeName{}).isSelectAs()
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

func TestArg(t *testing.T) {
	Arg(&IntervalArg{}).isArg()
	Arg(&ExprArg{}).isArg()
	Arg(&SequenceArg{}).isArg()
}

func TestInCondition(t *testing.T) {
	InCondition(&UnnestInCondition{}).isInCondition()
	InCondition(&SubQueryInCondition{}).isInCondition()
	InCondition(&ValuesInCondition{}).isInCondition()
}

func TestType(t *testing.T) {
	Type(&SimpleType{}).isType()
	Type(&ArrayType{}).isType()
	Type(&StructType{}).isType()
}

func TestIntValue(t *testing.T) {
	IntValue(&Param{}).isIntValue()
	IntValue(&IntLiteral{}).isIntValue()
	IntValue(&CastIntValue{}).isIntValue()
}

func TestNumValue(t *testing.T) {
	NumValue(&Param{}).isNumValue()
	NumValue(&IntLiteral{}).isNumValue()
	NumValue(&FloatLiteral{}).isNumValue()
	NumValue(&CastNumValue{}).isNumValue()
}

func TestStringValue(t *testing.T) {
	StringValue(&Param{}).isStringValue()
	StringValue(&StringLiteral{}).isStringValue()
}

func TestDDL(t *testing.T) {
	DDL(&CreateDatabase{}).isDDL()
	DDL(&CreateTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&CreateVectorIndex{}).isDDL()
	DDL(&CreateSequence{}).isDDL()
	DDL(&AlterSequence{}).isDDL()
	DDL(&DropSequence{}).isDDL()
	DDL(&CreateView{}).isDDL()
	DDL(&AlterTable{}).isDDL()
	DDL(&DropTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&AlterIndex{}).isDDL()
	DDL(&DropIndex{}).isDDL()
	DDL(&DropVectorIndex{}).isDDL()
	DDL(&CreateRole{}).isDDL()
	DDL(&DropRole{}).isDDL()
	DDL(&Grant{}).isDDL()
	DDL(&Revoke{}).isDDL()
	DDL(&CreateSearchIndex{}).isDDL()
	DDL(&DropSearchIndex{}).isDDL()
	DDL(&AlterSearchIndex{}).isDDL()
}

func TestConstraint(t *testing.T) {
	Constraint(&ForeignKey{}).isConstraint()
	Constraint(&Check{}).isConstraint()
}

func TestTableAlteration(t *testing.T) {
	TableAlteration(&AddColumn{}).isTableAlteration()
	TableAlteration(&AddTableConstraint{}).isTableAlteration()
	TableAlteration(&DropColumn{}).isTableAlteration()
	TableAlteration(&DropConstraint{}).isTableAlteration()
	TableAlteration(&SetOnDelete{}).isTableAlteration()
	TableAlteration(&AlterColumn{}).isTableAlteration()
	TableAlteration(&AlterColumnSet{}).isTableAlteration()
}

func TestPrivilege(t *testing.T) {
	Privilege(&PrivilegeOnTable{}).isPrivilege()
	Privilege(&SelectPrivilegeOnView{}).isPrivilege()
	Privilege(&ExecutePrivilegeOnTableFunction{}).isPrivilege()
	Privilege(&RolePrivilege{}).isPrivilege()
}

func TestTablePrivilege(t *testing.T) {
	TablePrivilege(&SelectPrivilege{}).isTablePrivilege()
	TablePrivilege(&InsertPrivilege{}).isTablePrivilege()
	TablePrivilege(&UpdatePrivilege{}).isTablePrivilege()
	TablePrivilege(&DeletePrivilege{}).isTablePrivilege()
}

func TestSchemaType(t *testing.T) {
	SchemaType(&ScalarSchemaType{}).isSchemaType()
	SchemaType(&SizedSchemaType{}).isSchemaType()
	SchemaType(&ArraySchemaType{}).isSchemaType()
}

func TestIndexAlteration(t *testing.T) {
	IndexAlteration(&AddStoredColumn{}).isIndexAlteration()
	IndexAlteration(&DropStoredColumn{}).isIndexAlteration()
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
