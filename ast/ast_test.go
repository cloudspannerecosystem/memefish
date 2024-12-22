package ast

import (
	"testing"
)

func TestStatement(t *testing.T) {
	Statement(&QueryStatement{}).isStatement()
	Statement(&CreateDatabase{}).isStatement()
	Statement(&AlterDatabase{}).isStatement()
	Statement(&CreateTable{}).isStatement()
	Statement(&AlterTable{}).isStatement()
	Statement(&DropTable{}).isStatement()
	Statement(&CreateIndex{}).isStatement()
	Statement(&AlterIndex{}).isStatement()
	Statement(&DropIndex{}).isStatement()
	Statement(&CreateView{}).isStatement()
	Statement(&DropView{}).isStatement()
	Statement(&CreateChangeStream{}).isStatement()
	Statement(&AlterChangeStream{}).isStatement()
	Statement(&DropChangeStream{}).isStatement()
	Statement(&CreateRole{}).isStatement()
	Statement(&DropRole{}).isStatement()
	Statement(&Grant{}).isStatement()
	Statement(&Revoke{}).isStatement()
	Statement(&CreateSequence{}).isStatement()
	Statement(&AlterSequence{}).isStatement()
	Statement(&DropSequence{}).isStatement()
	Statement(&CreateVectorIndex{}).isStatement()
	Statement(&DropVectorIndex{}).isStatement()
	Statement(&Insert{}).isStatement()
	Statement(&Delete{}).isStatement()
	Statement(&Update{}).isStatement()
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
	Expr(&TypedStructLiteral{}).isExpr()
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
	DDL(&AlterDatabase{}).isDDL()
	DDL(&CreateTable{}).isDDL()
	DDL(&AlterTable{}).isDDL()
	DDL(&DropTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&AlterIndex{}).isDDL()
	DDL(&DropIndex{}).isDDL()
	DDL(&CreateSearchIndex{}).isDDL()
	DDL(&DropSearchIndex{}).isDDL()
	DDL(&AlterSearchIndex{}).isDDL()
	DDL(&CreateView{}).isDDL()
	DDL(&DropView{}).isDDL()
	DDL(&CreateChangeStream{}).isDDL()
	DDL(&AlterChangeStream{}).isDDL()
	DDL(&DropChangeStream{}).isDDL()
	DDL(&CreateRole{}).isDDL()
	DDL(&DropRole{}).isDDL()
	DDL(&Grant{}).isDDL()
	DDL(&Revoke{}).isDDL()
	DDL(&CreateSequence{}).isDDL()
	DDL(&AlterSequence{}).isDDL()
	DDL(&DropSequence{}).isDDL()
	DDL(&CreateVectorIndex{}).isDDL()
	DDL(&DropVectorIndex{}).isDDL()
	DDL(&CreatePropertyGraph{}).isDDL()
	DDL(&DropPropertyGraph{}).isDDL()
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

func TestPropertyGraphLabelsOrProperties(t *testing.T) {
	PropertyGraphLabelsOrProperties(&PropertyGraphSingleProperties{}).isPropertyGraphLabelsOrProperties()
	PropertyGraphLabelsOrProperties(&PropertyGraphLabelAndPropertiesList{}).isPropertyGraphLabelsOrProperties()
}

func TestPropertyGraphElementLabel(t *testing.T) {
	PropertyGraphElementLabel(&PropertyGraphElementLabelLabelName{}).isPropertyGraphElementLabel()
	PropertyGraphElementLabel(&PropertyGraphElementLabelDefaultLabel{}).isPropertyGraphElementLabel()
}

func TestPropertyGraphElementKeys(t *testing.T) {
	PropertyGraphElementKeys(&PropertyGraphNodeElementKey{}).isPropertyGraphElementKeys()
	PropertyGraphElementKeys(&PropertyGraphEdgeElementKeys{}).isPropertyGraphElementKeys()
}

func TestPropertyGraphElementProperties(t *testing.T) {
	PropertyGraphElementProperties(&PropertyGraphNoProperties{}).isPropertyGraphElementProperties()
	PropertyGraphElementProperties(&PropertyGraphPropertiesAre{}).isPropertyGraphElementProperties()
	PropertyGraphElementProperties(&PropertyGraphDerivedPropertyList{}).isPropertyGraphElementProperties()
}
