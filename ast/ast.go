// Package ast provides AST nodes definitions.
//
// The definitions of ASTs are based on the following document.
//
//   - https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language
//   - https://cloud.google.com/spanner/docs/query-syntax
//
// Each Node's documentation describes its syntax (SQL representation) in a text/template
// fashion with thw following custom functions.
//
//   - sql node: Returns the SQL representation of node.
//   - sqlOpt node: Like sql node, but returns the empty string if node is nil.
//   - sqlJoin sep nodes: Concatenates the SQL representations of nodes with sep.
//   - sqlIdentQuote x: Quotes the given identifier string if needed.
//   - sqlStringQuote s: Returns the SQL quoted string of s.
//   - sqlBytesQuote bs: Returns the SQL quotes bytes of bs.
//   - tokenJoin toks: Concateates the string representations of tokens.
//   - isnil v: Checks whether v is nil or others.
//
// Each Node's documentation has pos and end information using the following EBNF.
//
//	PosChoice -> PosExpr ("||" PosExpr)*
//	PosExpr   -> PosAtom ("+" IntAtom)*
//	PosAtom   -> PosVar | NodeExpr "." ("pos" | "end")
//	NodeExpr  -> NodeAtom | "(" NodeAtom ("??" NodeAtom)* ")"
//	NodeAtom  -> NodeVar | NodeSliceVar "[" (IntAtom | "$") "]"
//	IntAtom   -> IntVal
//	           | "len" "(" StringVar ")"
//	           | "(" BoolVar "?" IntAtom ":" IntAtom ")"
//	IntVal    -> "0" | "1" | ...
//
//	(PosVar, NodeVar, NodeSliceVar, and BoolVar are derived by its struct definition.)
package ast

// NOTE: ast.go and ast_*.go are used for automatic generation, so these files are conventional.

// NOTE: This file defines AST nodes and they are used for automatic generation,
//       so this file is conventional.
//
// Conventions:
//
//   - Each node interface (except for Node) must have isXXX method (XXX is a name of the interface itself).
//   - `isXXX` methods must be defined after the interface definition
//     and the receiver must be the non-pointer node struct type.
//   - Each node struct must have pos and end comments.
//   - Each node struct must have template lines in its doc comment.
//   - The fields of each node must be ordered by the position.

//go:generate go run ../tools/gen-ast-pos/main.go -astfile ast.go -constfile ast_const.go -outfile pos.go

import (
	"github.com/cloudspannerecosystem/memefish/token"
)

// Node is base interface of Spanner SQL AST nodes.
type Node interface {
	Pos() token.Pos
	End() token.Pos

	// Convert AST node into SQL source string (a.k.a. Unparse).
	SQL() string
}

// Statement represents toplevel statement node of Spanner SQL.
type Statement interface {
	Node
	isStatement()
}

// The order of this list follows the official documentation:
//
// - https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language
// - https://cloud.google.com/spanner/docs/reference/standard-sql/dml-syntax

func (BadStatement) isStatement()        {}
func (BadDDL) isStatement()              {}
func (BadDML) isStatement()              {}
func (QueryStatement) isStatement()      {}
func (CreateSchema) isStatement()        {}
func (DropSchema) isStatement()          {}
func (CreateDatabase) isStatement()      {}
func (AlterDatabase) isStatement()       {}
func (CreateLocalityGroup) isStatement() {}
func (AlterLocalityGroup) isStatement()  {}
func (DropLocalityGroup) isStatement()   {}
func (CreatePlacement) isStatement()     {}
func (CreateProtoBundle) isStatement()   {}
func (AlterProtoBundle) isStatement()    {}
func (DropProtoBundle) isStatement()     {}
func (CreateTable) isStatement()         {}
func (AlterTable) isStatement()          {}
func (DropTable) isStatement()           {}
func (RenameTable) isStatement()         {}
func (CreateIndex) isStatement()         {}
func (AlterIndex) isStatement()          {}
func (DropIndex) isStatement()           {}
func (CreateSearchIndex) isStatement()   {}
func (DropSearchIndex) isStatement()     {}
func (AlterSearchIndex) isStatement()    {}
func (CreateView) isStatement()          {}
func (DropView) isStatement()            {}
func (CreateChangeStream) isStatement()  {}
func (AlterChangeStream) isStatement()   {}
func (DropChangeStream) isStatement()    {}
func (CreateRole) isStatement()          {}
func (DropRole) isStatement()            {}
func (Grant) isStatement()               {}
func (Revoke) isStatement()              {}
func (CreateSequence) isStatement()      {}
func (AlterSequence) isStatement()       {}
func (DropSequence) isStatement()        {}
func (AlterStatistics) isStatement()     {}
func (CreateModel) isStatement()         {}
func (AlterModel) isStatement()          {}
func (DropModel) isStatement()           {}
func (Analyze) isStatement()             {}
func (CreateFunction) isStatement()      {}
func (DropFunction) isStatement()        {}
func (CreateVectorIndex) isStatement()   {}
func (AlterVectorIndex) isStatement()    {}
func (DropVectorIndex) isStatement()     {}
func (CreatePropertyGraph) isStatement() {}
func (DropPropertyGraph) isStatement()   {}
func (Insert) isStatement()              {}
func (Delete) isStatement()              {}
func (Update) isStatement()              {}
func (Call) isStatement()                {}

// QueryExpr represents query expression, which can be body of QueryStatement or subqueries.
// Select and FromQuery are leaf QueryExpr and others wrap other QueryExpr.
type QueryExpr interface {
	Node
	isQueryExpr()
}

func (BadQueryExpr) isQueryExpr()  {}
func (Select) isQueryExpr()        {}
func (Query) isQueryExpr()         {}
func (FromQuery) isQueryExpr()     {}
func (SubQuery) isQueryExpr()      {}
func (CompoundQuery) isQueryExpr() {}

// PipeOperator represents pipe operator node which can be appeared in Query.
type PipeOperator interface {
	Node
	isPipeOperator()
}

func (PipeSelect) isPipeOperator() {}
func (PipeWhere) isPipeOperator()  {}

// SelectItem represents expression in SELECT clause result columns list.
type SelectItem interface {
	Node
	isSelectItem()
}

func (Star) isSelectItem()           {}
func (DotStar) isSelectItem()        {}
func (Alias) isSelectItem()          {}
func (ExprSelectItem) isSelectItem() {}

// SelectAs represents AS VALUE/STRUCT/typename clause in SELECT clause.
type SelectAs interface {
	Node
	isSelectAs()
}

func (AsStruct) isSelectAs()   {}
func (AsValue) isSelectAs()    {}
func (AsTypeName) isSelectAs() {}

// TableExpr represents JOIN operands.
type TableExpr interface {
	Node
	isTableExpr()
}

func (Unnest) isTableExpr()            {}
func (TableName) isTableExpr()         {}
func (PathTableExpr) isTableExpr()     {}
func (SubQueryTableExpr) isTableExpr() {}
func (ParenTableExpr) isTableExpr()    {}
func (Join) isTableExpr()              {}
func (TVFCallExpr) isTableExpr()       {}

// JoinCondition represents condition part of JOIN expression.
type JoinCondition interface {
	Node
	isJoinCondition()
}

func (On) isJoinCondition()    {}
func (Using) isJoinCondition() {}

// Expr repersents an expression in SQL.
type Expr interface {
	Node
	isExpr()
}

func (BadExpr) isExpr()               {}
func (BinaryExpr) isExpr()            {}
func (UnaryExpr) isExpr()             {}
func (InExpr) isExpr()                {}
func (IsNullExpr) isExpr()            {}
func (IsBoolExpr) isExpr()            {}
func (BetweenExpr) isExpr()           {}
func (SelectorExpr) isExpr()          {}
func (IndexExpr) isExpr()             {}
func (CallExpr) isExpr()              {}
func (CountStarExpr) isExpr()         {}
func (CastExpr) isExpr()              {}
func (ExtractExpr) isExpr()           {}
func (WithExpr) isExpr()              {}
func (ReplaceFieldsExpr) isExpr()     {}
func (CaseExpr) isExpr()              {}
func (IfExpr) isExpr()                {}
func (ParenExpr) isExpr()             {}
func (ScalarSubQuery) isExpr()        {}
func (ArraySubQuery) isExpr()         {}
func (ExistsSubQuery) isExpr()        {}
func (Param) isExpr()                 {}
func (Ident) isExpr()                 {}
func (Path) isExpr()                  {}
func (ArrayLiteral) isExpr()          {}
func (TupleStructLiteral) isExpr()    {}
func (TypelessStructLiteral) isExpr() {}
func (TypedStructLiteral) isExpr()    {}
func (NullLiteral) isExpr()           {}
func (BoolLiteral) isExpr()           {}
func (IntLiteral) isExpr()            {}
func (FloatLiteral) isExpr()          {}
func (StringLiteral) isExpr()         {}
func (BytesLiteral) isExpr()          {}
func (DateLiteral) isExpr()           {}
func (TimestampLiteral) isExpr()      {}
func (NumericLiteral) isExpr()        {}
func (JSONLiteral) isExpr()           {}
func (IntervalLiteralSingle) isExpr() {}
func (IntervalLiteralRange) isExpr()  {}
func (NewConstructor) isExpr()        {}
func (BracedNewConstructor) isExpr()  {}
func (BracedConstructor) isExpr()     {}

// SubscriptSpecifier represents specifier of subscript operators.
type SubscriptSpecifier interface {
	Node
	isSubscriptSpecifier()
}

func (ExprArg) isSubscriptSpecifier()                   {}
func (SubscriptSpecifierKeyword) isSubscriptSpecifier() {}

// Arg represents argument of function call.
type Arg interface {
	Node
	isArg()
}

func (ExprArg) isArg()     {}
func (SequenceArg) isArg() {}
func (LambdaArg) isArg()   {}

type TVFArg interface {
	Node
	isTVFArg()
}

func (ExprArg) isTVFArg()  {}
func (ModelArg) isTVFArg() {}
func (TableArg) isTVFArg() {}

// NullHandlingModifier represents IGNORE/RESPECT NULLS of aggregate function calls
type NullHandlingModifier interface {
	Node
	isNullHandlingModifier()
}

func (IgnoreNulls) isNullHandlingModifier()  {}
func (RespectNulls) isNullHandlingModifier() {}

// HavingModifier represents HAVING clause of aggregate function calls.
type HavingModifier interface {
	Node
	isHavingModifier()
}

func (HavingMax) isHavingModifier() {}
func (HavingMin) isHavingModifier() {}

// InCondition is right-side value of IN operator.
type InCondition interface {
	Node
	isInCondition()
}

func (UnnestInCondition) isInCondition()   {}
func (SubQueryInCondition) isInCondition() {}
func (ValuesInCondition) isInCondition()   {}

// TypelessStructLiteralArg represents an argument of typeless STRUCT literals.
type TypelessStructLiteralArg interface {
	Node
	isTypelessStructLiteralArg()
}

func (ExprArg) isTypelessStructLiteralArg() {}
func (Alias) isTypelessStructLiteralArg()   {}

// NewConstructorArg represents an argument of NEW constructors.
type NewConstructorArg interface {
	Node
	isNewConstructorArg()
}

func (ExprArg) isNewConstructorArg() {}
func (Alias) isNewConstructorArg()   {}

// Type represents type node.
type Type interface {
	Node
	isType()
}

func (BadType) isType()    {}
func (SimpleType) isType() {}
func (ArrayType) isType()  {}
func (StructType) isType() {}
func (NamedType) isType()  {}

// IntValue represents integer values in SQL.
type IntValue interface {
	Node
	isIntValue()
}

func (Param) isIntValue()        {}
func (IntLiteral) isIntValue()   {}
func (CastIntValue) isIntValue() {}

// NumValue represents number values in SQL.
type NumValue interface {
	Node
	isNumValue()
}

func (Param) isNumValue()        {}
func (IntLiteral) isNumValue()   {}
func (FloatLiteral) isNumValue() {}
func (CastNumValue) isNumValue() {}

// StringValue represents string value in SQL.
type StringValue interface {
	Node
	isStringValue()
}

func (Param) isStringValue()         {}
func (StringLiteral) isStringValue() {}

// DDL represents data definition language in SQL.
//
// https://cloud.google.com/spanner/docs/data-definition-language
type DDL interface {
	Statement
	isDDL()
}

// The order of this list follows the official documentation:
//
// - https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language

func (BadDDL) isDDL()              {}
func (CreateSchema) isDDL()        {}
func (DropSchema) isDDL()          {}
func (CreateDatabase) isDDL()      {}
func (AlterDatabase) isDDL()       {}
func (CreateLocalityGroup) isDDL() {}
func (AlterLocalityGroup) isDDL()  {}
func (DropLocalityGroup) isDDL()   {}
func (CreatePlacement) isDDL()     {}
func (CreateProtoBundle) isDDL()   {}
func (AlterProtoBundle) isDDL()    {}
func (DropProtoBundle) isDDL()     {}
func (CreateTable) isDDL()         {}
func (AlterTable) isDDL()          {}
func (DropTable) isDDL()           {}
func (RenameTable) isDDL()         {}
func (CreateIndex) isDDL()         {}
func (AlterIndex) isDDL()          {}
func (DropIndex) isDDL()           {}
func (CreateView) isDDL()          {}
func (CreateSearchIndex) isDDL()   {}
func (DropSearchIndex) isDDL()     {}
func (AlterSearchIndex) isDDL()    {}
func (DropView) isDDL()            {}
func (CreateChangeStream) isDDL()  {}
func (AlterChangeStream) isDDL()   {}
func (DropChangeStream) isDDL()    {}
func (CreateRole) isDDL()          {}
func (DropRole) isDDL()            {}
func (Grant) isDDL()               {}
func (Revoke) isDDL()              {}
func (CreateSequence) isDDL()      {}
func (AlterSequence) isDDL()       {}
func (DropSequence) isDDL()        {}
func (AlterStatistics) isDDL()     {}
func (CreateModel) isDDL()         {}
func (AlterModel) isDDL()          {}
func (DropModel) isDDL()           {}
func (Analyze) isDDL()             {}
func (CreateFunction) isDDL()      {}
func (DropFunction) isDDL()        {}
func (CreateVectorIndex) isDDL()   {}
func (AlterVectorIndex) isDDL()    {}
func (DropVectorIndex) isDDL()     {}
func (CreatePropertyGraph) isDDL() {}
func (DropPropertyGraph) isDDL()   {}

// Constraint represents table constraint of CONSTARINT clause.
type Constraint interface {
	Node
	isConstraint()
}

func (ForeignKey) isConstraint() {}
func (Check) isConstraint()      {}

// TableAlteration represents ALTER TABLE action.
type TableAlteration interface {
	Node
	isTableAlteration()
}

func (AddSynonym) isTableAlteration()               {}
func (DropSynonym) isTableAlteration()              {}
func (RenameTo) isTableAlteration()                 {}
func (AddColumn) isTableAlteration()                {}
func (AddTableConstraint) isTableAlteration()       {}
func (AddRowDeletionPolicy) isTableAlteration()     {}
func (DropColumn) isTableAlteration()               {}
func (DropConstraint) isTableAlteration()           {}
func (DropRowDeletionPolicy) isTableAlteration()    {}
func (ReplaceRowDeletionPolicy) isTableAlteration() {}
func (SetOnDelete) isTableAlteration()              {}
func (SetInterleaveIn) isTableAlteration()          {}
func (AlterColumn) isTableAlteration()              {}
func (AlterTableSetOptions) isTableAlteration()     {}

// ColumnDefaultSemantics is interface of DefaultExpr, GeneratedColumnExpr, IdentityColumn, AutoIncrement.
// They are change default value of column and mutually exclusive.
type ColumnDefaultSemantics interface {
	Node
	isColumnDefaultSemantics()
}

func (ColumnDefaultExpr) isColumnDefaultSemantics()   {}
func (GeneratedColumnExpr) isColumnDefaultSemantics() {}
func (IdentityColumn) isColumnDefaultSemantics()      {}
func (AutoIncrement) isColumnDefaultSemantics()       {}

type SequenceParam interface {
	Node
	isSequenceParam()
}

func (BitReversedPositive) isSequenceParam() {}
func (SkipRange) isSequenceParam()           {}
func (StartCounterWith) isSequenceParam()    {}

// ColumnAlteration represents ALTER COLUMN action in ALTER TABLE.
type ColumnAlteration interface {
	Node
	isColumnAlteration()
}

func (AlterColumnType) isColumnAlteration()          {}
func (AlterColumnSetOptions) isColumnAlteration()    {}
func (AlterColumnSetDefault) isColumnAlteration()    {}
func (AlterColumnDropDefault) isColumnAlteration()   {}
func (AlterColumnAlterIdentity) isColumnAlteration() {}

type IdentityAlteration interface {
	Node
	isIdentityAlteration()
}

func (RestartCounterWith) isIdentityAlteration() {}
func (SetSkipRange) isIdentityAlteration()       {}
func (SetNoSkipRange) isIdentityAlteration()     {}

// Privilege represents privileges specified by GRANT and REVOKE.
type Privilege interface {
	Node
	isPrivilege()
}

func (PrivilegeOnTable) isPrivilege()                {}
func (SelectPrivilegeOnChangeStream) isPrivilege()   {}
func (SelectPrivilegeOnView) isPrivilege()           {}
func (ExecutePrivilegeOnTableFunction) isPrivilege() {}
func (RolePrivilege) isPrivilege()                   {}

// TablePrivilege represents privileges on table.
type TablePrivilege interface {
	Node
	isTablePrivilege()
}

func (SelectPrivilege) isTablePrivilege() {}
func (InsertPrivilege) isTablePrivilege() {}
func (UpdatePrivilege) isTablePrivilege() {}
func (DeletePrivilege) isTablePrivilege() {}

// SchemaType represents types for schema.
type SchemaType interface {
	Node
	isSchemaType()
}

func (ScalarSchemaType) isSchemaType() {}
func (SizedSchemaType) isSchemaType()  {}
func (ArraySchemaType) isSchemaType()  {}
func (StructType) isSchemaType()       {}
func (NamedType) isSchemaType()        {}

// IndexAlteration represents ALTER INDEX action.
type IndexAlteration interface {
	Node
	isIndexAlteration()
}

func (AddStoredColumn) isIndexAlteration()  {}
func (DropStoredColumn) isIndexAlteration() {}

// VectorIndexAlteration represents ALTER VECTOR INDEX action.
// Note: Currently, it is same as IndexAlteration,
// but cloud-spanner-emulator/backend/schema/parser/ddl_parser.jjt implies their difference.
type VectorIndexAlteration interface {
	Node
	isVectorIndexAlteration()
}

func (AddStoredColumn) isVectorIndexAlteration()  {}
func (DropStoredColumn) isVectorIndexAlteration() {}

// DML represents data manipulation language in SQL.
//
// https://cloud.google.com/spanner/docs/data-definition-language
type DML interface {
	Statement
	isDML()
}

func (BadDML) isDML() {}
func (Insert) isDML() {}
func (Delete) isDML() {}
func (Update) isDML() {}

// InsertInput represents input values of INSERT statement.
type InsertInput interface {
	Node
	isInsertInput()
}

func (ValuesInput) isInsertInput()   {}
func (SubQueryInput) isInsertInput() {}

// ChangeStreamFor represents FOR clause in CREATE/ALTER CHANGE STREAM statement.
type ChangeStreamFor interface {
	Node
	isChangeStreamFor()
}

func (ChangeStreamForAll) isChangeStreamFor()    {}
func (ChangeStreamForTables) isChangeStreamFor() {}

// ChangeStreamAlteration represents ALTER CHANGE STREAM action.
type ChangeStreamAlteration interface {
	Node
	isChangeStreamAlteration()
}

func (ChangeStreamSetFor) isChangeStreamAlteration()     {}
func (ChangeStreamDropForAll) isChangeStreamAlteration() {}
func (ChangeStreamSetOptions) isChangeStreamAlteration() {}

// ================================================================================
//
// Bad Node
//
// ================================================================================

// BadNode is a placeholder node for a source code containing syntax errors.
//
//	{{.Tokens | tokenJoin}}
type BadNode struct {
	// pos = NodePos
	// end = NodeEnd

	NodePos, NodeEnd token.Pos

	Tokens []*token.Token
}

// BadStatement is a BadNode for Statement.
//
//	{{.Hint | sqlOpt}} {{.BadNode | sql}}
type BadStatement struct {
	// pos = (Hint ?? BadNode).pos
	// end = BadNode.end

	Hint    *Hint
	BadNode *BadNode
}

// BadQueryExpr is a BadNode for QueryExpr.
//
//	{{.BadNode | sql}}
type BadQueryExpr struct {
	// pos = BadNode.pos
	// end = BadNode.end

	Hint    *Hint
	BadNode *BadNode
}

// BadExpr is a BadNode for Expr.
//
//	{{.BadNode | sql}}
type BadExpr struct {
	// pos = BadNode.pos
	// end = BadNode.end

	BadNode *BadNode
}

// BadType is a BadNode for Type.
//
//	{{.BadNode | sql}}
type BadType struct {
	// pos = BadNode.pos
	// end = BadNode.end

	BadNode *BadNode
}

// BadDDL is a BadNode for DDL.
//
//	{{.BadNode | sql}}
type BadDDL struct {
	// pos = BadNode.pos
	// end = BadNode.end

	BadNode *BadNode
}

// BadDML is a BadNode for DML.
//
//	{{.Hint | sqlOpt}} {{.BadNode | sql}}
type BadDML struct {
	// pos = (Hint ?? BadNode).pos
	// end = BadNode.end

	Hint    *Hint // optional
	BadNode *BadNode
}

// ================================================================================
//
// SELECT
//
// ================================================================================

// QueryStatement is query statement node.
//
//	{{.Hint | sqlOpt}} {{.Query | sql}}
type QueryStatement struct {
	// pos = (Hint ?? Query).pos
	// end = Query.end

	Hint  *Hint // optional
	Query QueryExpr
}

// Query is query expression node with optional CTE, ORDER BY, LIMIT, and pipe operators.
// Usually, it is used as outermost QueryExpr in SubQuery and QueryStatement
//
//	{{.With | sqlOpt}}
//	{{.Query | sql}}
//	{{.OrderBy | sqlOpt}}
//	{{.Limit | sqlOpt}}
//	{{.PipeOperators | sqlJoin ", "}}
//
// https://cloud.google.com/spanner/docs/query-syntax
type Query struct {
	// pos = (With ?? Query).pos
	// end = (PipeOperators[$] ?? ForUpdate ?? Limit ?? OrderBy ?? Query).end

	With  *With
	Query QueryExpr

	OrderBy       *OrderBy   // optional
	Limit         *Limit     // optional
	ForUpdate     *ForUpdate // optional
	PipeOperators []PipeOperator
}

// ForUpdate is FOR UPDATE node.
//
//	FOR UPDATE
type ForUpdate struct {
	// pos = For
	// end = Update + 6

	For, Update token.Pos
}

// Hint is hint node.
//
//	@{{"{"}}{{.Records | sqlJoin ","}}{{"}"}}
type Hint struct {
	// pos = Atmark
	// end = Rbrace + 1

	Atmark token.Pos // position of "@"
	Rbrace token.Pos // position of "}"

	Records []*HintRecord // len(Records) > 0
}

// HintRecord is hint record node.
//
//	{{.Key | sql}}={{.Value | sql}}
type HintRecord struct {
	// pos = Key.pos
	// end = Value.end

	Key   *Path
	Value Expr
}

// With is with clause node.
//
//	WITH {{.CTEs | sqlJoin ","}}
type With struct {
	// pos = With
	// end = CTEs[$].end

	With token.Pos // position of "WITH" keyword

	CTEs []*CTE
}

// CTE is common table expression node.
//
//	{{.Name}} AS ({{.QueryExpr}})
type CTE struct {
	// pos = Name.pos
	// end = Rparen + 1

	Rparen token.Pos // position of ")"

	Name      *Ident
	QueryExpr QueryExpr
}

// Select is SELECT statement node.
//
//	SELECT
//	  {{.AllOrDistinct}}
//	  {{.As | sqlOpt}}
//	  {{.Results | sqlJoin ","}}
//	  {{.From | sqlOpt}}
//	  {{.Where | sqlOpt}}
//	  {{.GroupBy | sqlOpt}}
//	  {{.Having | sqlOpt}}
type Select struct {
	// pos = Select
	// end = (Having ?? GroupBy ?? Where ?? From ?? Results[$]).end

	Select token.Pos // position of "select" keyword

	AllOrDistinct AllOrDistinct // optional
	As            SelectAs      // optional
	Results       []SelectItem  // len(Results) > 0
	From          *From         // optional
	Where         *Where        // optional
	GroupBy       *GroupBy      // optional
	Having        *Having       // optional
}

// AsStruct represents AS STRUCT node in SELECT clause.
//
//	AS STRUCT
type AsStruct struct {
	// pos = As
	// end = Struct + 6

	As     token.Pos
	Struct token.Pos
}

// AsValue represents AS VALUE node in SELECT clause.
//
//	AS VALUE
type AsValue struct {
	// pos = As
	// end = Value + 5

	As    token.Pos
	Value token.Pos
}

// AsTypeName represents AS typename node in SELECT clause.
//
//	AS {{.TypeName | sql}}
type AsTypeName struct {
	// pos = As
	// end = TypeName.end

	As       token.Pos
	TypeName *NamedType
}

// FromQuery is FROM query expression node.
//
//	FROM {{.From | sql}}
type FromQuery struct {
	// pos = From.pos
	// end = From.end

	From *From
}

// CompoundQuery is query expression node compounded by set operators.
// Note: A single CompoundQuery can express query expressions compounded by the same set operator.
// If there are mixed Op or Distinct in query expression, CompoundQuery will be nested.
//
//	{{.Queries | sqlJoin (printf "%s %s" .Op .AllOrDistinct)}}
type CompoundQuery struct {
	// pos = Queries[0].pos
	// end = Queries[$].end

	Op            SetOp
	AllOrDistinct AllOrDistinct
	Queries       []QueryExpr // len(Queries) >= 2
}

// SubQuery is parenthesized query expression node.
// Note: subquery expression is expressed as a ParenTableExpr. Maybe better to rename as like ParenQueryExpr?
//
//	({{.Query | sql}})
type SubQuery struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query QueryExpr
}

// StarModifierExcept is EXCEPT node in Star and DotStar of SelectItem.
//
//	EXCEPT ({{Columns | sqlJoin ", "}})
type StarModifierExcept struct {
	// pos = Except
	// end = Rparen + 1

	Except token.Pos
	Rparen token.Pos

	Columns []*Ident
}

// StarModifierReplaceItem is a single item of StarModifierReplace.
//
//	{{.Expr | sql}} AS {{.Name | sql}}
type StarModifierReplaceItem struct {
	// pos = Expr.pos
	// end = Name.end

	Expr Expr
	Name *Ident
}

// StarModifierReplace is REPLACE node in Star and DotStar of SelectItem.
//
//	REPLACE ({{Columns | sqlJoin ", "}})
type StarModifierReplace struct {
	// pos = Replace
	// end = Rparen + 1

	Replace token.Pos
	Rparen  token.Pos

	Columns []*StarModifierReplaceItem
}

// Star is a single * in SELECT result columns list.
//
//	{{"*"}} {{.Except | sqlOpt}} {{.Replace | sqlOpt}}
//
// Note: The text/template notation escapes * to avoid normalize * to - by some formatters
// because the prefix * is ambiguous with bulletin list.
type Star struct {
	// pos = Star
	// end = (Replace ?? Except).end || Star + 1

	Star token.Pos // position of "*"

	Except  *StarModifierExcept  // optional
	Replace *StarModifierReplace // optional
}

// DotStar is expression with * in SELECT result columns list.
//
//	{{.Expr | sql}}.* {{.Except | sqlOpt}} {{.Replace | sqlOpt}}
type DotStar struct {
	// pos = Expr.pos
	// end = (Replace ?? Except).end || Star + 1

	Star token.Pos // position of "*"

	Expr    Expr
	Except  *StarModifierExcept  // optional
	Replace *StarModifierReplace // optional
}

// Alias is aliased expression by AS clause.
//
// Typically, this appears in SELECT result columns list, but this can appear in typeless STRUCT literals
// and NEW constructors.
//
//	{{.Expr | sql}} {{.As | sql}}
type Alias struct {
	// pos = Expr.pos
	// end = As.end

	Expr Expr
	As   *AsAlias
}

// AsAlias is AS clause node for general purpose.
//
// It is used in Alias node and some JoinExpr nodes.
//
//	{{if not .As.Invalid}}AS {{end}}{{.Alias | sql}}
type AsAlias struct {
	// pos = As || Alias.pos
	// end = Alias.end

	As token.Pos // position of "AS" keyword, optional

	Alias *Ident
}

// ExprSelectItem is Expr wrapper for SelectItem.
//
//	{{.Expr | sql}}
type ExprSelectItem struct {
	// pos = Expr.pos
	// end = Expr.end

	Expr Expr
}

// From is FROM clause node.
//
//	FROM {{.Source | sql}}
type From struct {
	// pos = From
	// end = Source.end

	From token.Pos // position of "FROM" keyword

	Source TableExpr
}

// Where is WHERE clause node.
//
//	WHERE {{.Expr | sql}}
type Where struct {
	// pos = Where
	// end = Expr.end

	Where token.Pos // position of "WHERE" keyword

	Expr Expr
}

// GroupBy is GROUP BY clause node.
//
//	GROUP {{.Hint | sqlOpt}} BY {{.Exprs | sqlJoin ","}}
type GroupBy struct {
	// pos = Group
	// end = Exprs[$].end

	Group token.Pos // position of "GROUP" keyword

	Hint  *Hint  // optional
	Exprs []Expr // len(Exprs) > 0
}

// Having is HAVING clause node.
//
//	HAVING {{.Expr | sql}}
type Having struct {
	// pos = Having
	// end = Expr.end

	Having token.Pos // position of "HAVING" keyword

	Expr Expr
}

// OrderBy is ORDER BY clause node.
//
//	ORDER BY {{.Items | sqlJoin ","}}
type OrderBy struct {
	// pos = Order
	// end = Items[$].end

	Order token.Pos // position of "ORDER" keyword

	Items []*OrderByItem // len(Items) > 0
}

// OrderByItem is expression node in ORDER BY clause list.
//
//	{{.Expr | sql}} {{.Collate | sqlOpt}} {{.Direction}}
type OrderByItem struct {
	// pos = Expr.pos
	// end = DirPos + len(Dir) || (Collate ?? Expr).end

	DirPos token.Pos // position of Dir

	Expr    Expr
	Collate *Collate  // optional
	Dir     Direction // optional
}

// Collate is COLLATE clause node in ORDER BY item.
//
//	COLLATE {{.Value | sql}}
type Collate struct {
	// pos = Collate
	// end = Value.end

	Collate token.Pos // position of "COLLATE" keyword

	Value StringValue
}

// Limit is LIMIT clause node.
//
//	LIMIT {{.Count | sql}} {{.Offset | sqlOpt}}
type Limit struct {
	// pos = Limit
	// end = (Offset ?? Count).end

	Limit token.Pos // position of "LIMIT" keyword

	Count  IntValue
	Offset *Offset // optional
}

// Offset is OFFSET clause node in LIMIT clause.
//
//	OFFSET {{.Value | sql}}
type Offset struct {
	// pos = Offset
	// end = Value.end

	Offset token.Pos // position of "OFFSET" keyword

	Value IntValue
}

// ================================================================================
//
// Pipe Operators
//
// Must be the same order in document.
// https://github.com/google/zetasql/blob/master/docs/pipe-syntax.md#pipe-operators
// TODO: The reference should be Spanner reference when it is supported.
//
// ================================================================================

// PipeSelect is SELECT pipe operator node.
//
//	|> SELECT {{.AllOrDistinct}} {{.As | sqlOpt}} {{.Results | sqlJoin ", "}}
type PipeSelect struct {
	// pos = Pipe
	// end = Results[$].end

	Pipe token.Pos // position of "|>"

	AllOrDistinct AllOrDistinct // optional
	As            SelectAs      // optional
	Results       []SelectItem  // len(Results) > 0
}

// PipeWhere is WHERE pipe operator node.
//
//	|> WHERE {{.Expr | sql}}
type PipeWhere struct {
	// pos = Pipe
	// end = Expr.end

	Pipe token.Pos // position of "|>"

	Expr Expr
}

// ================================================================================
//
// JOIN
//
// ================================================================================

// Unnest is UNNEST call in FROM clause.
//
//	UNNEST({{.Expr | sql}})
//	{{.Hint | sqlOpt}}
//	{{.As | sqlOpt}}
//	{{.WithOffset | sqlOpt}}
//	{{.Sample | sqlOpt}}
type Unnest struct {
	// pos = Unnest
	// end = (Sample ?? WithOffset ?? As ?? Hint).end || Rparen + 1 || Expr.end

	Unnest token.Pos // position of "UNNEST"
	Rparen token.Pos // position of ")"

	Expr       Expr         // Path or Ident when Implicit is true
	Hint       *Hint        // optional
	As         *AsAlias     // optional
	WithOffset *WithOffset  // optional
	Sample     *TableSample // optional
}

// WithOffset is WITH OFFSET clause node after UNNEST call.
//
//	WITH OFFSET {{.As | sqlOpt}}
type WithOffset struct {
	// pos = With
	// end = As.end || Offset + 6

	With, Offset token.Pos // position of "WITH" and "OFFSET" keywords

	As *AsAlias // optional
}

// TableName is table name node in FROM clause.
//
//	{{.Table | sql}} {{.Hint | sqlOpt}} {{.As | sqlOpt}} {{.Sample | sqlOpt}}
type TableName struct {
	// pos = Table.pos
	// end = (Sample ?? As ?? Hint ?? Table).end

	Table  *Ident
	Hint   *Hint        // optional
	As     *AsAlias     // optional
	Sample *TableSample // optional
}

// PathTableExpr is path expression node in FROM clause.
// Parser cannot distinguish between `implicit UNNEST` and tables in a named schema.
// It is the job of a later phase to determine this distinction.
//
//	{{.Path | sql}} {{.Hint | sqlOpt}} {{.As | sqlOpt}} {{.Sample | sqlOpt}}
type PathTableExpr struct {
	// pos = Path.pos
	// end = (Sample ?? WithOffset ?? As ?? Hint ?? Path).end

	Path       *Path
	Hint       *Hint        // optional
	As         *AsAlias     // optional
	WithOffset *WithOffset  // optional
	Sample     *TableSample // optional
}

// SubQueryTableExpr is subquery inside JOIN expression.
//
//	({{.Query | sql}}) {{.As | sqlOpt}} {{.Sample | sqlOpt}}
type SubQueryTableExpr struct {
	// pos = Lparen
	// end = (Sample ?? As).end || Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query  QueryExpr
	As     *AsAlias     // optional
	Sample *TableSample // optional
}

// ParenTableExpr is parenthesized JOIN expression.
//
//	({{.Source | sql}}) {{.Sample | sqlOpt}}
type ParenTableExpr struct {
	// pos = Lparen
	// end = Sample.end || Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Source TableExpr    // SubQueryJoinExpr (without As) or Join
	Sample *TableSample // optional
}

// Join is JOIN expression.
//
//	{{.Left | sql}}
//	  {{.Op}} {{.Method}} {{.Hint | sqlOpt}}
//	  {{.Right | sql}}
//	{{.Cond | sqlOpt}}
type Join struct {
	// pos = Left.pos
	// end = (Cond ?? Right).end

	Left   TableExpr
	Op     JoinOp
	Method JoinMethod
	Hint   *Hint // optional
	Right  TableExpr

	// nil when Op is CrossJoin
	// optional when Right is PathTableExpr or Unnest
	// otherwise it must be set.
	Cond JoinCondition
}

// On is ON condition of JOIN expression.
//
//	ON {{.Expr | sql}}
type On struct {
	// pos = On
	// end = Expr.end

	On token.Pos // position of "ON" keyword

	Expr Expr
}

// Using is Using condition of JOIN expression.
//
//	USING ({{Idents | sqlJoin ","}})
type Using struct {
	// pos = Using
	// end = Rparen + 1

	Using  token.Pos // position of "USING" keyword
	Rparen token.Pos // position of ")"

	Idents []*Ident // len(Idents) > 0
}

// TableSample is TABLESAMPLE clause node.
//
//	TABLESAMPLE {{.Method}} {{.Size | sql}}
type TableSample struct {
	// pos = TableSample
	// end = Size.end

	TableSample token.Pos // position of "TABLESAMPLE" keyword

	Method TableSampleMethod
	Size   *TableSampleSize
}

// TableSampleSize is size part of TABLESAMPLE clause.
//
//	({{.Value | sql}} {{.Unit}})
type TableSampleSize struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Value NumValue
	Unit  TableSampleUnit
}

// ================================================================================
//
// Expr
//
// ================================================================================

// BinaryExpr is binary operator expression node.
//
//	{{.Left | sql}} {{.Op}} {{.Right | sql}}
type BinaryExpr struct {
	// pos = Left.pos
	// end = Right.end

	Op BinaryOp

	Left, Right Expr
}

// UnaryExpr is unary operator expression node.
//
//	{{.Op}} {{.Expr | sql}}
type UnaryExpr struct {
	// pos = OpPos
	// end = Expr.end

	OpPos token.Pos // position of Op

	Op   UnaryOp
	Expr Expr
}

// InExpr is IN expression node.
//
//	{{.Left | sql}} {{if .Not}}NOT{{end}} IN {{.Right | sql}}
type InExpr struct {
	// pos = Left.pos
	// end = Right.end

	Not   bool
	Left  Expr
	Right InCondition
}

// UnnestInCondition is UNNEST call at IN condition.
//
//	UNNEST({{.Expr | sql}})
type UnnestInCondition struct {
	// pos = Unnest
	// end = Rparen + 1

	Unnest token.Pos
	Rparen token.Pos

	Expr Expr
}

// SubQueryInCondition is subquery at IN condition.
//
//	({{.Query | sql}})
type SubQueryInCondition struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query QueryExpr
}

// ValuesInCondition is parenthesized values at IN condition.
//
//	({{.Exprs | sqlJoin ","}})
type ValuesInCondition struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Exprs []Expr // len(Exprs) > 0
}

// IsNullExpr is IS NULL expression node.
//
//	{{.Left | sql}} IS {{if .Not}}NOT{{end}} NULL
type IsNullExpr struct {
	// pos = Left.pos
	// end = Null + 4

	Null token.Pos // position of "NULL"

	Not  bool
	Left Expr
}

// IsBoolExpr is IS TRUE/FALSE expression node.
//
//	{{.Left | sql}} IS {{if .Not}}NOT{{end}} {{if .Right}}TRUE{{else}}FALSE{{end}}
type IsBoolExpr struct {
	// pos = Left.pos
	// end = RightPos + (Right ? 4 : 5)

	RightPos token.Pos // position of Right

	Not   bool
	Left  Expr
	Right bool
}

// BetweenExpr is BETWEEN expression node.
//
//	{{.Left | sql}} {{if .Not}}NOT{{end}} BETWEEN {{.RightStart | sql}} AND {{.RightEnd | sql}}
type BetweenExpr struct {
	// pos = Left.pos
	// end = RightEnd.end

	Not                        bool
	Left, RightStart, RightEnd Expr
}

// SelectorExpr is struct field access expression node.
//
//	{{.Expr | sql}}.{{.Ident | sql}}
type SelectorExpr struct {
	// pos = Expr.pos
	// end = Ident.end

	Expr  Expr
	Ident *Ident
}

// IndexExpr is a subscript operator expression node.
// This node can be:
//   - array subscript operator
//   - struct subscript operator
//   - JSON subscript operator
//
// Note: The name IndexExpr is a historical reason, maybe better to rename to SubscriptExpr.
//
//	{{.Expr | sql}}[{{.Index | sql}}]
type IndexExpr struct {
	// pos = Expr.pos
	// end = Rbrack + 1

	Rbrack token.Pos // position of "]"

	Expr  Expr
	Index SubscriptSpecifier
}

// SubscriptSpecifierKeyword is subscript specifier with position keyword.
//
//	{{string(.Keyword)}}({{.Expr | sql}})
type SubscriptSpecifierKeyword struct {
	// pos = KeywordPos
	// end = Rparen + 1

	KeywordPos token.Pos // position of Keyword
	Rparen     token.Pos // position of ")"

	Keyword PositionKeyword
	Expr    Expr
}

// CallExpr is function call expression node.
// It can represent both regular function calls and aggregate function calls.
//
//	{{.Func | sql}}(
//		{{if .Distinct}}DISTINCT{{end}}
//		{{.Args | sqlJoin ", "}}
//		{{if len(.Args) > 0 && len(.NamedArgs) > 0}}, {{end}}
//		{{.NamedArgs | sqlJoin ", "}}
//		{{.NullHandling | sqlOpt}}
//		{{.Having | sqlOpt}}
//		{{.OrderBy | sqlOpt}}
//		{{.Limit | sqlOpt}}
//	)
//	{{.Hint | sqlOpt}}
type CallExpr struct {
	// pos = Func.pos
	// end = Hint.end || Rparen + 1

	Rparen token.Pos // position of ")"

	Func         *Path
	Distinct     bool
	Args         []Arg
	NamedArgs    []*NamedArg
	NullHandling NullHandlingModifier // optional
	Having       HavingModifier       // optional
	OrderBy      *OrderBy             // optional
	Limit        *Limit               // optional
	Hint         *Hint                // optional
}

// TVFCallExpr is table-valued function call expression node.
//
//	{{.Name | sql}}(
//		{{.Args | sqlJoin ", "}}
//		{{if len(.Args) > 0 && len(.NamedArgs) > 0}}, {{end}}
//		{{.NamedArgs | sqlJoin ", "}}
//	)
//	{{.Hint | sqlOpt}}
//	{{.Sample | sqlOpt}}
type TVFCallExpr struct {
	// pos = Name.pos
	// end = (Sample ?? Hint).end || Rparen + 1

	Rparen token.Pos // position of ")"

	Name      *Path
	Args      []TVFArg
	NamedArgs []*NamedArg
	Hint      *Hint        // optional
	Sample    *TableSample // optional
}

// ExprArg is argument of the generic function call.
//
//	{{.Expr | sql}}
type ExprArg struct {
	// pos = Expr.pos
	// end = Expr.end

	Expr Expr
}

// SequenceArg is argument of sequence function call.
//
//	SEQUENCE {{.Expr | sql}}
type SequenceArg struct {
	// pos = Sequence
	// end = Expr.end

	Sequence token.Pos // position of "SEQUENCE" keyword

	Expr Expr
}

// LambdaArg is lambda expression argument of the generic function call.
//
//	{{if .Lparen.Invalid}}{{.Args | sqlJoin ", "}}{{else}}({{.Args | sqlJoin ", "}}) -> {{.Expr | sql}}
//
// Note: Args won't be empty. If Lparen is not appeared, Args have exactly one element.
type LambdaArg struct {
	// pos = Lparen || Args[0].pos
	// end = Expr.end

	Lparen token.Pos // optional

	Args []*Ident // if Lparen.Invalid() then len(Args) = 1 else len(Args) > 0
	Expr Expr
}

// ModelArg is argument of model function call.
//
//	MODEL {{.Name | sql}}
type ModelArg struct {
	// pos = Model
	// end = Name.end

	Model token.Pos // position of "MODEL" keyword

	Name *Path
}

// TableArg is TABLE table_name argument of table valued function call.
//
//	TABLE {{.Name | sql}}
type TableArg struct {
	// pos = Table
	// end = Name.end

	Table token.Pos // position of "TABLE" keyword

	Name *Path
}

// NamedArg represents a name and value pair in named arguments
//
//	{{.Name | sql}} => {{.Value | sql}}
type NamedArg struct {
	// pos = Name.pos
	// end = Value.end

	Name  *Ident
	Value Expr
}

// IgnoreNulls represents IGNORE NULLS of aggregate function calls.
//
//	IGNORE NULLS
type IgnoreNulls struct {
	// pos = Ignore
	// end = Nulls + 5

	Ignore token.Pos
	Nulls  token.Pos
}

// RespectNulls represents RESPECT NULLS of aggregate function calls
//
//	RESPECT NULLS
type RespectNulls struct {
	// pos = Respect
	// end = Nulls + 5

	Respect token.Pos
	Nulls   token.Pos
}

// HavingMax represents HAVING MAX of aggregate function calls.
//
//	HAVING MAX {{Expr | sql}}
type HavingMax struct {
	// pos = Having
	// end = Expr.end

	Having token.Pos
	Expr   Expr
}

// HavingMin represents HAVING MIN of aggregate function calls.
//
//	HAVING MIN {{Expr | sql}}
type HavingMin struct {
	// pos = Having
	// end = Expr.end

	Having token.Pos
	Expr   Expr
}

// CountStarExpr is node just for COUNT(*).
//
//	COUNT(*)
type CountStarExpr struct {
	// pos = Count
	// end = Rparen + 1

	Count  token.Pos // position of "COUNT"
	Rparen token.Pos // position of ")"
}

// ExtractExpr is EXTRACT call expression node.
//
//	EXTRACT({{.Part | sql}} FROM {{.Expr | sql}} {{.AtTimeZone | sqlOpt}})
type ExtractExpr struct {
	// pos = Extract
	// end = Rparen + 1

	Extract token.Pos // position of "EXTRACT" keyword
	Rparen  token.Pos // position of ")"

	Part       *Ident
	Expr       Expr
	AtTimeZone *AtTimeZone // optional
}

// ReplaceFieldsArg is value AS field_path node in ReplaceFieldsExpr.
//
//	{{.Expr | sql}} AS {{.Field | sql}}
type ReplaceFieldsArg struct {
	// pos = Expr.pos
	// end = Field.end

	Expr  Expr
	Field *Path
}

// ReplaceFieldsExpr is REPLACE_FIELDS call expression node.
//
//	REPLACE_FIELDS({{.Expr.| sql}}, {{.Fields | sqlJoin ", "}})
type ReplaceFieldsExpr struct {
	// pos = ReplaceFields
	// end = Rparen + 1

	ReplaceFields token.Pos // position of "REPLACE_FIELDS" keyword
	Rparen        token.Pos // position of ")"

	Expr   Expr
	Fields []*ReplaceFieldsArg
}

// AtTimeZone is AT TIME ZONE clause in EXTRACT call.
//
//	AT TIME ZONE {{.Expr | sql}}
type AtTimeZone struct {
	// pos = At
	// end = Expr.end

	At token.Pos // position of "AT" keyword

	Expr Expr
}

// WithExprVar is "name AS expr" node in WITH expression.
//
//	{{.Name | sql}} AS {{.Expr | sql}}
type WithExprVar struct {
	// pos = Name.pos
	// end = Expr.end

	Name *Ident
	Expr Expr
}

// WithExpr is WITH expression node.
//
//	WITH({{.Vars | sqlJoin ", "}}, {{.Expr | sql}})
type WithExpr struct {
	// pos = With
	// end = Rparen + 1

	With   token.Pos // position of "WITH" keyword
	Rparen token.Pos // position of ")"

	Vars []*WithExprVar // len(Vars) > 0
	Expr Expr
}

// CastExpr is CAST/SAFE_CAST call expression node.
//
//	{{if .Safe}}SAFE_{{end}}CAST({{.Expr | sql}} AS {{.Type | sql}})
type CastExpr struct {
	// pos = Cast
	// end = Rparen + 1

	Cast   token.Pos // position of "CAST" keyword or "SAFE_CAST" pseudo keyword
	Rparen token.Pos // position of ")"

	Safe bool

	Expr Expr
	Type Type
}

// CaseExpr is CASE expression node.
//
//	CASE {{.Expr | sqlOpt}}
//	  {{.Whens | sqlJoin "\n"}}
//	  {{.Else | sqlOpt}}
//	END
type CaseExpr struct {
	// pos = Case
	// end = EndPos + 3

	Case, EndPos token.Pos // position of "CASE" and "END" keywords

	Expr  Expr        // optional
	Whens []*CaseWhen // len(Whens) > 0
	Else  *CaseElse   // optional
}

// CaseWhen is WHEN clause in CASE expression.
//
//	WHEN {{.Cond | sql}} THEN {{.Then | sql}}
type CaseWhen struct {
	// pos = When
	// end = Then.end

	When token.Pos // position of "WHEN" keyword

	Cond, Then Expr
}

// CaseElse is ELSE clause in CASE expression.
//
//	ELSE {{.Expr | sql}}
type CaseElse struct {
	// pos = Else
	// end = Expr.end

	Else token.Pos // position of "ELSE" keyword

	Expr Expr
}

// IfExpr is IF conditional expression.
// Because IF is SQL keyword, it can't be a normal CallExpr.
//
//	IF({{.Expr | sql}}, {{.TrueResult | sql}}, {{.ElseResult | sql}})
type IfExpr struct {
	// pos = If
	// end = Rparen + 1

	If     token.Pos // position of "IF" keyword
	Rparen token.Pos // position of ")"

	Expr       Expr
	TrueResult Expr
	ElseResult Expr
}

// ParenExpr is parenthesized expression node.
//
//	({{. | sql}})
type ParenExpr struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Expr Expr
}

// ScalarSubQuery is subquery in expression.
//
//	({{.Query | sql}})
type ScalarSubQuery struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query QueryExpr
}

// ArraySubQuery is subquery in ARRAY call.
//
//	ARRAY({{.Query | sql}})
type ArraySubQuery struct {
	// pos = Array
	// end = Rparen + 1

	Array  token.Pos // position of "ARRAY" keyword
	Rparen token.Pos // position of ")"

	Query QueryExpr
}

// ExistsSubQuery is subquery in EXISTS call.
//
//	EXISTS {{.Hint | sqlOpt}} ({{.Query | sql}})
type ExistsSubQuery struct {
	// pos = Exists
	// end = Rparen + 1

	Exists token.Pos // position of "EXISTS" keyword
	Rparen token.Pos // position of ")"

	Hint  *Hint
	Query QueryExpr
}

// ================================================================================
//
// Literal
//
// ================================================================================

// Param is Query parameter node.
//
//	@{{.Name}}
type Param struct {
	// pos = Atmark
	// end = Atmark + 1 + len(Name)

	Atmark token.Pos

	Name string
}

// Ident is identifier node.
//
//	{{.Name | sqlIdentQuote}}
type Ident struct {
	// pos = NamePos
	// end = NameEnd

	NamePos, NameEnd token.Pos // position of this name

	Name string
}

// Path is dot-chained identifier list.
// It can be simple name without dot.
//
//	{{.Idents | sqlJoin "."}}
type Path struct {
	// pos = Idents[0].pos
	// end = Idents[$].end

	Idents []*Ident // len(Idents) > 0
}

// ArrayLiteral is array literal node.
//
//	{{if .Array.Invalid | not}}ARRAY{{end}}{{if .Type}}<{{.Type | sql}}>{{end}}[{{.Values | sqlJoin ","}}]
type ArrayLiteral struct {
	// pos = Array || Lbrack
	// end = Rbrack + 1

	Array          token.Pos // position of "ARRAY" keyword, optional
	Lbrack, Rbrack token.Pos // position of "[" and "]"

	Type   Type // optional
	Values []Expr
}

// TupleStructLiteral is tuple syntax struct literal node.
//
//	({{.Values | sqlJoin ","}})
type TupleStructLiteral struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Values []Expr // len(Values) > 1
}

// TypedStructLiteral is typed struct literal node.
//
//	STRUCT<{{.Fields | sqlJoin ","}}>({{.Values | sqlJoin ","}})
type TypedStructLiteral struct {
	// pos = Struct
	// end = Rparen + 1

	Struct token.Pos // position of "STRUCT"
	Rparen token.Pos // position of ")"

	Fields []*StructField
	Values []Expr
}

// TypelessStructLiteral is typeless struct literal node.
//
//	STRUCT({{.Values | sqlJoin ","}})
type TypelessStructLiteral struct {
	// pos = Struct
	// end = Rparen + 1

	Struct token.Pos // position of "STRUCT"
	Rparen token.Pos // position of ")"

	Values []TypelessStructLiteralArg
}

// NullLiteral is just NULL literal.
//
//	NULL
type NullLiteral struct {
	// pos = Null
	// end = Null + 4

	Null token.Pos // position of "NULL"
}

// BoolLiteral is boolean literal node.
//
//	{{if .Value}}TRUE{{else}}FALSE{{end}}
type BoolLiteral struct {
	// pos = ValuePos
	// end = ValuePos + (Value ? 4 : 5)

	ValuePos token.Pos // position of this value

	Value bool
}

// IntLiteral is integer literal node.
//
//	{{.Value}}
type IntLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Base  int // 10 or 16
	Value string
}

// FloatLiteral is floating point number literal node.
//
//	{{.Value}}
type FloatLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value string
}

// StringLiteral is string literal node.
//
//	{{.Value | sqlStringQuote}}
type StringLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value string
}

// BytesLiteral is bytes literal node.
//
//	B{{.Value | sqlBytesQuote}}
type BytesLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value []byte
}

// DateLiteral is date literal node.
//
//	DATE {{.Value | sql}}
type DateLiteral struct {
	// pos = Date
	// end = Value.end

	Date token.Pos // position of "DATE"

	Value *StringLiteral
}

// TimestampLiteral is timestamp literal node.
//
//	TIMESTAMP {{.Value | sql}}
type TimestampLiteral struct {
	// pos = Timestamp
	// end = Value.end

	Timestamp token.Pos // position of "TIMESTAMP"

	Value *StringLiteral
}

// NumericLiteral is numeric literal node.
//
//	NUMERIC {{.Value | sql}}
type NumericLiteral struct {
	// pos = Numeric
	// end = Value.end

	Numeric token.Pos // position of "NUMERIC"

	Value *StringLiteral
}

// JSONLiteral is JSON literal node.
//
//	JSON {{.Value | sql}}
type JSONLiteral struct {
	// pos = JSON
	// end = Value.end

	JSON token.Pos // position of "JSON"

	Value *StringLiteral
}

// IntervalLiteralSingle represents an interval literal with a single datetime part.
//
//	INTERVAL {{.Value}} {{.DateTimePart | sql}}
type IntervalLiteralSingle struct {
	// pos = Interval
	// end = DateTimePartEnd

	Interval        token.Pos // position of "INTERVAL" keyword
	DateTimePartEnd token.Pos

	Value IntValue

	DateTimePart DateTimePart
}

// IntervalLiteralRange represents an interval literal with a datetime part range.
//
//	INTERVAL {{.Value}} {{.StartingDateTimePart | sql}} TO {{.EndingDateTimePart | sql}}
type IntervalLiteralRange struct {
	// pos = Interval
	// end = EndingDateTimePartEnd

	Interval              token.Pos // position of "INTERVAL" keyword
	EndingDateTimePartEnd token.Pos

	Value *StringLiteral

	StartingDateTimePart DateTimePart
	EndingDateTimePart   DateTimePart
}

// ================================================================================
//
// NEW constructors
//
// ================================================================================

// BracedConstructorFieldValue represents value part of fields in BracedNewConstructor.
type BracedConstructorFieldValue interface {
	Node
	isBracedConstructorFieldValue()
}

func (BracedConstructor) isBracedConstructorFieldValue()               {}
func (BracedConstructorFieldValueExpr) isBracedConstructorFieldValue() {}

// NewConstructor represents NEW operator which creates a protocol buffer using a parenthesized list of arguments.
//
//	NEW {{.Type | sql}} ({{.Args | sqlJoin ", "}})
type NewConstructor struct {
	// pos = New
	// end = Rparen + 1

	New  token.Pos
	Type *NamedType

	Args []NewConstructorArg

	Rparen token.Pos
}

// BracedNewConstructor represents NEW operator which creates a protocol buffer using a map constructor.
//
//	NEW {{.Type | sql}} {{"{"}}{{"}"}}
type BracedNewConstructor struct {
	// pos = New
	// end = Body.end

	New token.Pos

	Type *NamedType
	Body *BracedConstructor
}

// BracedConstructor represents a single map constructor which is used in BracedNewConstructor.
// Actually, it is a top level Expr in syntax, but it is not permitted semantically in other place.
//
//	{{"{"}}{{.Fields | sqlJoin ", "}}{{"}"}}
type BracedConstructor struct {
	// pos = Lbrace
	// end = Rbrace + 1

	Lbrace, Rbrace token.Pos

	Fields []*BracedConstructorField
}

// BracedConstructorField represents a single field in BracedConstructor.
//
//	{{.Name | sql}} {{.Value | sql}}
type BracedConstructorField struct {
	// pos = Name.pos
	// end = Value.end

	Name  *Ident
	Value BracedConstructorFieldValue
}

// BracedConstructorFieldValueExpr represents a field value node.
//
//	: {{.Expr | sql}}
type BracedConstructorFieldValueExpr struct {
	// pos = Colon
	// end = Expr.end

	Colon token.Pos

	Expr Expr
}

// ================================================================================
//
// Type
//
// ================================================================================

// SimpleType is type node having no parameter like INT64, STRING.
//
//	{{.Name}}
type SimpleType struct {
	// pos = NamePos
	// end = NamePos + len(Name)

	NamePos token.Pos // position of this name

	Name ScalarTypeName
}

// ArrayType is array type node.
//
//	ARRAY<{{.Item | sql}}>
type ArrayType struct {
	// pos = Array
	// end = Gt + 1

	Array token.Pos // position of "ARRAY" keyword
	Gt    token.Pos // position of ">"

	Item Type
}

// StructType is struct type node.
//
//	STRUCT<{{.Fields | sqlJoin ","}}>
type StructType struct {
	// pos = Struct
	// end = Gt + 1

	Struct token.Pos // position of "STRUCT"
	Gt     token.Pos // position of ">"

	Fields []*StructField
}

// StructField is field in struct type node.
//
//	{{if .Ident}}{{.Ident | sql}}{{end}} {{.Type | sql}}
type StructField struct {
	// pos = (Ident ?? Type).pos
	// end = Type.end

	Ident *Ident
	Type  Type
}

// NamedType is named type node.
// It is currently PROTO or ENUM.
// Name is full qualified name, but it can be len(Name) == 1 if it doesn't contain ".".
//
//	{{.Path | sqlJoin "."}}
type NamedType struct {
	// pos = Path[0].pos
	// end = Path[$].end

	Path []*Ident // len(Path) > 0
}

// ================================================================================
//
// Cast for Special Cases
//
// ================================================================================

// CastIntValue is cast call in integer value context.
//
//	CAST({{.Expr | sql}} AS INT64)
type CastIntValue struct {
	// pos = Cast
	// end = Rparen + 1

	Cast   token.Pos // position of "CAST" keyword
	Rparen token.Pos // position of ")"

	Expr IntValue // IntLit or Param
}

// CasrNumValue is cast call in number value context.
//
//	CAST({{.Expr | sql}} AS {{.Type}})
type CastNumValue struct {
	// pos = Cast
	// end = Rparen + 1

	Cast   token.Pos // position of "CAST" keyword
	Rparen token.Pos // position of ")"

	Expr NumValue       // IntLit, FloatLit or Param
	Type ScalarTypeName // Int64Type or Float64Type
}

// ================================================================================
//
// DDL
//
// ================================================================================

// Options is generic OPTIONS clause node without key and value checking.
//
//	OPTIONS ({{.Records | sqlJoin ","}})
type Options struct {
	// pos = Options
	// end = Rparen + 1

	Options token.Pos // position of "OPTIONS" keyword
	Rparen  token.Pos // position of ")"

	Records []*OptionsDef // len(Records) > 0
}

// OptionsDef is single option definition for DDL statements.
//
//	{{.Name | sql}} = {{.Value | sql}}
type OptionsDef struct {
	// pos = Name.pos
	// end = Value.end

	Name  *Ident
	Value Expr
}

// CreateSchema is CREATE SCHEMA statement node.
//
//	CREATE {{if .OrReplace}}OR REPLACE{{end}} SCHEMA {{.Name | sql}}
type CreateSchema struct {
	// pos = Create
	// end = Name.end

	Create      token.Pos // position of "CREATE" keyword
	OrReplace   bool
	IfNotExists bool
	Name        *Ident
}

// DropSchema is DROP SCHEMA statement node.
//
//	DROP SCHEMA{{if .IfExists}} IF EXISTS{{end}} {{.Name | sql}}
type DropSchema struct {
	// pos = Drop
	// end = Name.end

	Drop     token.Pos // position of "DROP" keyword
	IfExists bool
	Name     *Ident
}

// CreateDatabase is CREATE DATABASE statement node.
//
//	CREATE DATABASE {{.Name | sql}}
type CreateDatabase struct {
	// pos = Create
	// end = Name.end

	Create token.Pos // position of "CREATE" keyword

	Name *Ident
}

// AlterDatabase is ALTER DATABASE statement node.
//
//	ALTER DATABASE {{.Name | sql}} SET {{.Options | sql}}
type AlterDatabase struct {
	// pos = Alter
	// end = Options.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *Options
}

// CreateLocalityGroup is CREATE LOCALITY GROUP statement node.
//
//	CREATE LOCALITY GROUP {{.Name | sql}} {{.Options | sqlOpt}}
type CreateLocalityGroup struct {
	// pos = Create
	// end = (Options ?? Name).end

	Create token.Pos // position of "CREATE" keyword

	Name    *Ident
	Options *Options // optional
}

// AlterLocalityGroup is ALTER LOCALITY GROUP statement node.
//
//	ALTER LOCALITY GROUP {{.Name | sql}} SET {{.Options | sql}}
type AlterLocalityGroup struct {
	// pos = Alter
	// end = Options.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *Options
}

// DropLocalityGroup is DROP LOCALITY GROUP statement node.
//
//	DROP LOCALITY GROUP {{.Name | sql}}
type DropLocalityGroup struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// CreatePlacement is CREATE PLACEMENT statement node.
//
//	CREATE PLACEMENT {{.Name | sql}} {{.Options | sqlOpt}}
type CreatePlacement struct {
	// pos = Create
	// end = (Options ?? Name).end

	Create token.Pos // position of "CREATE" keyword

	Name    *Ident
	Options *Options // optional
}

// ================================================================================
//
// PROTO BUNDLE statements
//
// ================================================================================

// ProtoBundleTypes is parenthesized Protocol Buffers type names node IN CREATE/ALTER PROTO BUNDLE statement.
//
//	({{.Types | sqlJoin ", "}})
type ProtoBundleTypes struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos
	Types          []*NamedType
}

// CreateProtoBundle is CREATE PROTO BUNDLE statement node.
//
//	CREATE PROTO BUNDLE {{.Types, | sql}}
type CreateProtoBundle struct {
	// pos = Create
	// end = Types.end

	Create token.Pos // position of "CREATE" keyword

	Types *ProtoBundleTypes
}

// AlterProtoBundle is ALTER PROTO BUNDLE statement node.
//
//	ALTER PROTO BUNDLE {{.Insert | sqlOpt}} {{.Update | sqlOpt}} {{.Delete | sqlOpt}}
type AlterProtoBundle struct {
	// pos = Alter
	// end = (Delete ?? Update ?? Insert).end || Bundle + 6

	Alter  token.Pos // position of "ALTER" keyword
	Bundle token.Pos

	Insert *AlterProtoBundleInsert // optional
	Update *AlterProtoBundleUpdate // optional
	Delete *AlterProtoBundleDelete // optional
}

// AlterProtoBundleInsert is INSERT (proto_type_name, ...) node in ALTER PROTO BUNDLE.
//
//	INSERT {{.Types | sql}}
type AlterProtoBundleInsert struct {
	// pos = Insert
	// end = Types.end

	Insert token.Pos // position of "INSERT" pseudo keyword

	Types *ProtoBundleTypes
}

// AlterProtoBundleUpdate is UPDATE (proto_type_name, ...) node in ALTER PROTO BUNDLE.
//
//	UPDATE {{.Types | sql}}
type AlterProtoBundleUpdate struct {
	// pos = Update
	// end = Types.end

	Update token.Pos // position of "UPDATE" pseudo keyword

	Types *ProtoBundleTypes
}

// AlterProtoBundleDelete is DELETE (proto_type_name, ...) node in ALTER PROTO BUNDLE.
//
//	DELETE {{.Types | sql}}
type AlterProtoBundleDelete struct {
	// pos = Delete
	// end = Types.end

	Delete token.Pos // position of "DELETE" pseudo keyword

	Types *ProtoBundleTypes
}

// DropProtoBundle is DROP PROTO BUNDLE statement node.
//
//	DROP PROTO BUNDLE
type DropProtoBundle struct {
	// pos = Drop
	// end = Bundle + 6

	Drop   token.Pos // position of "DROP" pseudo keyword
	Bundle token.Pos // position of "BUNDLE" pseudo keyword
}

// CreateTable is CREATE TABLE statement node.
//
//	CREATE TABLE {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}} (
//	  {{.Columns | sqlJoin ","}}{{if and .Columns (or .TableConstrains .Synonym)}},{{end}}
//	  {{.TableConstraints | sqlJoin ","}}{{if and .TableConstraints .Synonym}},{{end}}
//	  {{.Synonym | sqlJoin ","}}
//	)
//	{{if .PrimaryKeys}}PRIMARY KEY ({{.PrimaryKeys | sqlJoin ","}}){{end}}
//	{{.Cluster | sqlOpt}}
//	{{.CreateRowDeletionPolicy | sqlOpt}}
//	{{if .Options}}, {{.Options | sqlOpt}}{{end}}
//
// Spanner SQL allows to mix `Columns` and `TableConstraints` and `Synonyms`,
// however they are separated in AST definition for historical reasons. If you want to get
// the original order of them, please sort them by their `Pos()`.
type CreateTable struct {
	// pos = Create
	// end = Options.end || RowDeletionPolicy.end || Cluster.end || PrimaryKeyRparen + 1 || Rparen + 1

	Create           token.Pos // position of "CREATE" keyword
	Rparen           token.Pos // position of ")" of end of column definitions
	PrimaryKeyRparen token.Pos // position of ")" of PRIMARY KEY clause, optional

	IfNotExists       bool
	Name              *Path
	Columns           []*ColumnDef
	TableConstraints  []*TableConstraint
	PrimaryKeys       []*IndexKey // when omitted, len(PrimaryKeys) = 0
	Synonyms          []*Synonym
	Cluster           *Cluster                 // optional
	RowDeletionPolicy *CreateRowDeletionPolicy // optional
	Options           *Options                 // optional
}

// Synonym is SYNONYM node in CREATE TABLE
//
//	SYNONYM ({.Name | sql})
type Synonym struct {
	// pos = Synonym
	// end = Rparen + 1

	Synonym token.Pos // position of "SYNONYM" pseudo keyword
	Rparen  token.Pos // position of ")"

	Name *Ident
}

// CreateSequence is CREATE SEQUENCE statement node.
//
//	CREATE SEQUENCE {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}} }}
//	{{.Params | sqlJoin " "}}
//	{{.Options | sql}}
type CreateSequence struct {
	// pos = Create
	// end = Options.end || Params[$].end || Name.end

	Create token.Pos // position of "CREATE" keyword

	Name        *Path
	IfNotExists bool
	Params      []SequenceParam // len(Params) >= 0
	Options     *Options        // optional
}

// SkipRange is SKIP RANGE node.
//
//	SKIP RANGE {{.Min | sql}}, {{.Max | sql}}
type SkipRange struct {
	// pos = Skip
	// end = Max.end

	Skip token.Pos // position of "SKIP" keyword

	Min, Max *IntLiteral
}

// StartCounterWith is START COUNTER WITH node.
//
//	START COUNTER WITH {{.Counter | sql}}
type StartCounterWith struct {
	// pos = Start
	// end = Counter.end

	Start token.Pos // position of "START" keyword

	Counter *IntLiteral
}

// BitReversedPositive is BIT_RESVERSED_POSITIVE node.
//
//	BIT_REVERSED_POSITIVE
type BitReversedPositive struct {
	// pos = BitReversedPositive
	// end = BitReversedPositive + 21

	BitReversedPositive token.Pos // position of "BIT_REVERSED_POSITIVE" keyword
}

// ColumnDef is column definition in CREATE TABLE and ALTER TABLE ADD COLUMN.
// Note: Some fields are not valid in ADD COLUMN.
//
//	{{.Name | sql}}
//	{{.Type | sql}} {{if .NotNull}}NOT NULL{{end}}
//	{{.DefaultSemantics | sqlOpt}}
//	{{if .Hidden.Invalid | not)}}HIDDEN{{end}}
//	{{if .PrimaryKey}}PRIMARY KEY{{end}}
//	{{.Options | sqlOpt}}
type ColumnDef struct {
	// pos = Name.pos
	// end = Options.end || Key + 3 || Hidden + 6 || DefaultSemantics.end || Null + 4 || Type.end

	Null token.Pos // position of "NULL"
	Key  token.Pos // position of "KEY" of PRIMARY KEY

	Name       *Ident
	Type       SchemaType
	NotNull    bool
	PrimaryKey bool

	DefaultSemantics ColumnDefaultSemantics // optional

	Hidden  token.Pos // InvalidPos if not hidden
	Options *Options  // optional
}

// ColumnDefaultExpr is a default value expression for the column.
//
//	DEFAULT ({{.Expr | sql}})
type ColumnDefaultExpr struct {
	// pos = Default
	// end = Rparen

	Default token.Pos // position of "DEFAULT" keyword
	Rparen  token.Pos // position of ")"

	Expr Expr
}

// GeneratedColumnExpr is generated column expression.
//
//	AS ({{.Expr | sql}}) {{if .IsStored}}STORED{{end}}
type GeneratedColumnExpr struct {
	// pos = As
	// end = Stored + 6 || Rparen + 1

	As     token.Pos // position of "AS" keyword
	Stored token.Pos // position of "STORED" keyword, optional
	Rparen token.Pos // position of ")"

	Expr Expr
}

// IdentityColumn is GENERATED BY DEFAULT AS IDENTITY node.
//
//	GENERATED BY DEFAULT AS IDENTITY {{if not(.Rparen.Invalid)}}({{.Params | sqlJoin " "}})
type IdentityColumn struct {
	// pos = Generated
	// end = Rparen + 1 || Identity + 8

	Generated token.Pos // position of "GENERATED" keyword
	Identity  token.Pos // position of "IDENTITY" keyword
	Rparen    token.Pos // position of ")", optional

	Params []SequenceParam //  if Rparen.Invalid() then len(Param) = 0 else len(Param) > 0
}

// AutoIncrement is AUTO_INCREMENT node.
//
//	AUTO_INCREMENT
type AutoIncrement struct {
	// pos = AutoIncrement
	// end = AutoIncrement + 14

	AutoIncrement token.Pos // position of "AUTO_INCREMENT"
}

// TableConstraint is table constraint in CREATE TABLE and ALTER TABLE.
//
//	{{if .Name}}CONSTRAINT {{.Name}}{{end}}{{.Constraint | sql}}
type TableConstraint struct {
	// pos = ConstraintPos || Constraint.pos
	// end = Constraint.end

	ConstraintPos token.Pos // position of "CONSTRAINT" keyword when Name presents

	Name       *Ident // optional
	Constraint Constraint
}

// ForeignKey is foreign key specifier in CREATE TABLE and ALTER TABLE.
//
//	FOREIGN KEY ({{.ColumnNames | sqlJoin ","}}) REFERENCES {{.ReferenceTable}} ({{.ReferenceColumns | sqlJoin ","}})
//	{{.OnDelete}} {{.Enforcement}}
type ForeignKey struct {
	// pos = Foreign
	// end = Enforced + 8 || OnDeleteEnd || Rparen + 1

	Foreign     token.Pos // position of "FOREIGN" keyword
	Rparen      token.Pos // position of ")" after reference columns
	OnDeleteEnd token.Pos // end position of ON DELETE clause
	Enforced    token.Pos // position of "ENFORCED", optional

	Columns          []*Ident
	ReferenceTable   *Path
	ReferenceColumns []*Ident       // len(ReferenceColumns) > 0
	OnDelete         OnDeleteAction // optional
	Enforcement      Enforcement    // optional
}

// Check is check constraint in CREATE TABLE and ALTER TABLE.
//
//	Check ({{.Expr}})
type Check struct {
	// pos = Check
	// end = Rparen + 1

	Check  token.Pos // position of "CHECK" keyword
	Rparen token.Pos // position of ")" after Expr

	Expr Expr
}

// IndexKey is index key specifier in CREATE TABLE and CREATE INDEX.
//
//	{{.Name | sql}} {{.Dir}}
type IndexKey struct {
	// pos = Name.pos
	// end = DirPos + len(Dir) || Name.end

	DirPos token.Pos // position of Dir

	Name *Ident
	Dir  Direction // optional
}

// Cluster is INTERLEAVE IN [PARENT] clause in CREATE TABLE.
//
//	, INTERLEAVE IN {{if .Enforced}}PARENT{{end}} {{.TableName | sql}} {{.OnDelete}}
type Cluster struct {
	// pos = Comma
	// end = OnDeleteEnd || TableName.end

	Comma       token.Pos // position of ","
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	TableName *Path
	Enforced  bool           // true when PARENT is present
	OnDelete  OnDeleteAction // optional
}

// CreateRowDeletionPolicy is ROW DELETION POLICY clause in CREATE TABLE.
//
//	, {{.RowDeletionPolicy}}
type CreateRowDeletionPolicy struct {
	// pos = Comma
	// end = RowDeletionPolicy.end

	Comma token.Pos // position of ","

	RowDeletionPolicy *RowDeletionPolicy
}

// RowDeletionPolicy is ROW DELETION POLICY clause.
//
//	ROW DELETION POLICY (OLDER_THAN({{.ColymnName | sql}}, INTERVAL {{.NumDays}} DAY))
type RowDeletionPolicy struct {
	// pos = Row
	// end = Rparen + 1

	Row    token.Pos // position of "ROW"
	Rparen token.Pos // position of ")"

	ColumnName *Ident
	NumDays    *IntLiteral
}

// CreateView is CREATE VIEW statement node.
//
//	CREATE {{if .OrReplace}}OR REPLACE{{end}} VIEW {{.Name | sql}}
//	SQL SECURITY {{.SecurityType}} AS
//	{{.Query | sql}}
type CreateView struct {
	// pos = Create
	// end = Query.end

	Create token.Pos

	Name         *Path
	OrReplace    bool
	SecurityType SecurityType
	Query        QueryExpr
}

// DropView is DROP VIEW statement node.
//
//	DROP VIEW {{.Name | sql}}
type DropView struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos

	Name *Path
}

// AlterTable is ALTER TABLE statement node.
//
//	ALTER TABLE {{.Name | sql}} {{.TableAlteration | sql}}
type AlterTable struct {
	// pos = Alter
	// end = TableAlteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name            *Path
	TableAlteration TableAlteration
}

// AlterIndex is ALTER INDEX statement node.
//
//	ALTER INDEX {{.Name | sql}} {{.IndexAlteration | sql}}
type AlterIndex struct {
	// pos = Alter
	// end = IndexAlteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name            *Path
	IndexAlteration IndexAlteration
}

// AlterSequence is ALTER SEQUENCE statement node.
//
//	ALTER SEQUENCE {{.Name | sql}}
//	{{if .Options}}SET {{.Options | sqlOpt}}{{end}}
//	{{.RestartCounterWith | sqlOpt}}
//	{{.SkipRange | sqlOpt}}
//	{{.NoSkipRange | sqlOpt}}
type AlterSequence struct {
	// pos = Alter
	// end = (NoSkipRange ?? SkipRange ?? RestartCounterWith ?? Options).end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Path
	Options *Options // optional

	RestartCounterWith *RestartCounterWith // optional
	SkipRange          *SkipRange          // optional
	NoSkipRange        *NoSkipRange        // optional
}

// AlterChangeStream is ALTER CHANGE STREAM statement node.
//
//	ALTER CHANGE STREAM {{.Name | sql}} {{.ChangeStreamAlteration | sql}}
type AlterChangeStream struct {
	// pos = Alter
	// end = ChangeStreamAlteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name                   *Ident
	ChangeStreamAlteration ChangeStreamAlteration
}

// AddSynonym is ADD SYNONYM node in ALTER TABLE.
//
//	ADD SYNONYM {{.Name | sql}}
type AddSynonym struct {
	// pos = Add
	// end = Name.end

	Add  token.Pos // position of "ADD" pseudo keyword
	Name *Ident
}

// DropSynonym is DROP SYNONYM node in ALTER TABLE.
//
//	DROP SYNONYM {{.Name | sql}}
type DropSynonym struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" pseudo keyword
	Name *Ident
}

// RenameTo is RENAME TO node in ALTER TABLE.
//
//	RENAME TO {{.Name | sql}}{{if .AddSynonym}}, {{.AddSynonym | sql}}{{end}}
type RenameTo struct {
	// pos = Rename
	// end = (AddSynonym ?? Name).end

	Rename token.Pos // position of "RENAME" pseudo keyword

	Name       *Ident
	AddSynonym *AddSynonym // optional
}

// AddColumn is ADD COLUMN clause in ALTER TABLE.
//
//	ADD COLUMN {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Column | sql}}
type AddColumn struct {
	// pos = Add
	// end = Column.end

	Add token.Pos // position of "ADD" keyword

	IfNotExists bool
	Column      *ColumnDef
}

// AddTableConstraint is ADD table_constraint clause in ALTER TABLE.
//
//	ADD {{.TableConstraint}}
type AddTableConstraint struct {
	// pos = Add
	// end = TableConstraint.end

	Add token.Pos // position of "ADD" keyword

	TableConstraint *TableConstraint
}

// AddRowDeletionPolicy is ADD ROW DELETION POLICY clause in ALTER TABLE.
//
//	ADD {{.RowDeletionPolicy | sql}}
type AddRowDeletionPolicy struct {
	// pos = Add
	// end = RowDeletionPolicy.end

	Add token.Pos // position of "ADD" keyword

	RowDeletionPolicy *RowDeletionPolicy
}

// DropColumn is DROP COLUMN clause in ALTER TABLE.
//
//	DROP COLUMN {{.Name | sql}}
type DropColumn struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of  "DROP" keyword

	Name *Ident
}

// DropConstraint is DROP CONSTRAINT clause in ALTER TABLE.
//
//	DROP CONSTRAINT {{.Name | sql}}
type DropConstraint struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of  "DROP" keyword

	Name *Ident
}

// DropRowDeletionPolicy is DROP ROW DELETION POLICY clause in ALTER TABLE.
//
//	DROP ROW DELETION POLICY
type DropRowDeletionPolicy struct {
	// pos = Drop
	// end = Policy + 6

	Drop   token.Pos // position of  "DROP" keyword
	Policy token.Pos // position of  "POLICY" keyword
}

// ReplaceRowDeletionPolicy is REPLACE ROW DELETION POLICY clause in ALTER TABLE.
//
//	REPLACE {{.RowDeletionPolicy}}
type ReplaceRowDeletionPolicy struct {
	// pos = Replace
	// end = RowDeletionPolicy.end

	Replace token.Pos // position of  "REPLACE" keyword

	RowDeletionPolicy *RowDeletionPolicy
}

// SetOnDelete is SET ON DELETE clause in ALTER TABLE.
//
//	SET ON DELETE {{.OnDelete}}
type SetOnDelete struct {
	// pos = Set
	// end = OnDeleteEnd

	Set         token.Pos // position of "SET" keyword
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	OnDelete OnDeleteAction
}

// SetInterleaveIn is SET INTERLEAVE IN clause in ALTER TABLE.
//
//	SET INTERLEAVE IN {{if .Enforced}}PARENT{{end}} {{.TableName | sql}} {{.OnDelete}}
type SetInterleaveIn struct {
	// pos = Set
	// end = OnDeleteEnd || TableName.end

	Set         token.Pos // position of "SET" keyword
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	TableName *Path
	Enforced  bool
	OnDelete  OnDeleteAction // optional
}

// AlterTableSetOptions is SET OPTIONS node in ALTER TABLE.
//
//	SET {{.Options | sql}}
type AlterTableSetOptions struct {
	// pos = Set
	// end = Options.end

	Set     token.Pos
	Options *Options
}

// AlterColumn is ALTER COLUMN clause in ALTER TABLE.
//
//	ALTER COLUMN {{.Name | sql}} {{.Alteration | sql}}
type AlterColumn struct {
	// pos = Alter
	// end = Alteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name       *Ident
	Alteration ColumnAlteration
}

// AlterColumnType is action to change the data type of the column in ALTER COLUMN.
//
//	{{.Type | sql}} {{if .NotNull}}NOT NULL{{end}} {{.DefaultExpr | sqlOpt}}
type AlterColumnType struct {
	// pos = Type.pos
	// end = DefaultExpr.end || Null + 4 || Type.end

	Type        SchemaType
	Null        token.Pos // position of "NULL" keyword, optional
	NotNull     bool
	DefaultExpr *ColumnDefaultExpr // optional
}

// AlterColumnSetOptions is SET OPTIONS node in ALTER COLUMN.
//
//	SET {{.Options | sql}}
type AlterColumnSetOptions struct {
	// pos = Set
	// end = Options.end

	Set     token.Pos
	Options *Options
}

// AlterColumnSetDefault is SET DEFAULT node in ALTER COLUMN.
//
//	SET {{.DefaultExpr | sql}}
type AlterColumnSetDefault struct {
	// pos = Set
	// end = DefaultExpr.end

	Set         token.Pos
	DefaultExpr *ColumnDefaultExpr
}

// AlterColumnDropDefault is DROP DEFAULT node in ALTER COLUMN
//
//	DROP DEFAULT
type AlterColumnDropDefault struct {
	// pos = Drop
	// end = Default + 7

	Drop    token.Pos
	Default token.Pos
}

// AlterColumnAlterIdentity is ALTER IDENTITY node in ALTER COLUMN
//
//	ALTER IDENTITY {{.Alteration | sql}}
type AlterColumnAlterIdentity struct {
	// pos = Alter
	// end = Alteration.end

	Alter      token.Pos
	Alteration IdentityAlteration
}

// RestartCounterWith is RESTART COUNTER WITH node.
//
//	RESTART COUNTER WITH {{.Counter | sql}}
type RestartCounterWith struct {
	// pos = Restart
	// end = Counter.end

	Restart token.Pos // position of "RESTART" keyword

	Counter *IntLiteral
}

// SetSkipRange is SET SKIP RANGE node
//
//	SET {{.SkipRange | sql}}
type SetSkipRange struct {
	// pos = Set
	// end = SkipRange.end

	Set       token.Pos
	SkipRange *SkipRange
}

// NoSkipRange is NO SKIP RANGE node.
//
//	NO SKIP RANGE
type NoSkipRange struct {
	// pos = No
	// end = Range + 5

	No, Range token.Pos
}

// SetNoSkipRange is SET NO SKIP RANGE node.
//
//	SET {{.NoSkipRange | sql}}
type SetNoSkipRange struct {
	// pos = Set
	// end = NoSkipRange.end

	Set         token.Pos
	NoSkipRange *NoSkipRange
}

// DropTable is DROP TABLE statement node.
//
//	DROP TABLE {{if .IfExists}}IF NOT EXISTS{{end}} {{.Name | sql}}
type DropTable struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	IfExists bool
	Name     *Path
}

// RenameTable is RENAME TABLE statement node.
//
//	RENAME TABLE {{.Tos | sqlJoin ", "}}
type RenameTable struct {
	// pos = Rename
	// end = Tos[$].end

	Rename token.Pos // position of "RENAME" pseudo keyword

	Tos []*RenameTableTo // len(Tos) > 0
}

// RenameTableTo is old TO new node in RENAME TABLE statement.
//
//	{{.Old | sql}} TO {{.New | sql}}
type RenameTableTo struct {
	// pos = Old.pos
	// end = New.end

	Old *Ident
	New *Ident
}

// CreateIndex is CREATE INDEX statement node.
//
//	CREATE
//	  {{if .Unique}}UNIQUE{{end}}
//	  {{if .NullFiltered}}NULL_FILTERED{{end}}
//	INDEX {{if .IfExists}}IF NOT EXISTS{{end}} {{.Name | sql}} ON {{.TableName | sql}} (
//	  {{.Keys | sqlJoin ","}}
//	)
//	{{.Storing | sqlOpt}}
//	{{.InterleaveIn | sqlOpt}}
//	{{.Options | sqlOpt}}
type CreateIndex struct {
	// pos = Create
	// end = (Options ?? InterleaveIn ?? Storing).end || Rparen + 1

	Create token.Pos // position of "CREATE" keyword
	Rparen token.Pos // position of ")"

	Unique       bool
	NullFiltered bool
	IfNotExists  bool
	Name         *Path
	TableName    *Path
	Keys         []*IndexKey
	Storing      *Storing      // optional
	InterleaveIn *InterleaveIn // optional
	Options      *Options      // optional
}

// CreateVectorIndex is CREATE VECTOR INDEX statement node.
//
//	CREATE VECTOR INDEX {if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}}
//	ON {{.TableName | sql}}({{.ColumnName | sql}})
//	{{.Storing | sqlOpt}}
//	{{if .Where}}WHERE {{.Where | sql}}{{end}}
//	{{.Options | sql}}
type CreateVectorIndex struct {
	// pos = Create
	// end = Options.end

	Create token.Pos // position of "CREATE" keyword

	IfNotExists bool // optional
	Name        *Ident
	TableName   *Ident
	ColumnName  *Ident
	Storing     *Storing // optional

	// It only allows `WHERE column_name IS NOT NULL` for now, but we still relax the condition
	// by reusing the `parseWhere` function for sake of it may be extended more conditions in the future.
	//
	// Reference: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#vector_index_statements
	Where   *Where // optional
	Options *Options
}

// AlterVectorIndex is ALTER VECTOR INDEX statement node.
//
//	ALTER VECTOR INDEX {{.Name | sql}} {{.Alteration | sql}}
type AlterVectorIndex struct {
	// pos = Alter
	// end = Alteration.end

	Alter      token.Pos
	Name       *Path
	Alteration VectorIndexAlteration
}

// CreateChangeStream is CREATE CHANGE STREAM statement node.
//
//	CREATE CHANGE STREAM {{.Name | sql}} {{.For | sqlOpt}} {{.Options | sqlOpt}}
type CreateChangeStream struct {
	// pos = Create
	// end = (Options ?? For ?? Name).end

	Create token.Pos // position of "CREATE" keyword

	Name    *Ident
	For     ChangeStreamFor // optional
	Options *Options        // optional
}

// ChangeStreamForAll is FOR ALL node in CREATE CHANGE STREAM
//
//	FOR ALL
type ChangeStreamForAll struct {
	// pos = For
	// end = All + 3

	For token.Pos // position of "FOR" keyword
	All token.Pos // position of "ALL" keyword
}

// ChangeStreamForTables is FOR tables node in CREATE CHANGE STREAM
//
//	FOR {{.Tables | sqlJoin ","}}
type ChangeStreamForTables struct {
	// pos = For
	// end = Tables[$].end

	For token.Pos // position of "FOR" keyword

	Tables []*ChangeStreamForTable
}

// ChangeStreamForTable table node in CREATE CHANGE STREAM SET FOR
//
//	{{.TableName | sql}}{{if not(.Rparen.Invalid)}}({{.Columns | sqlJoin ","}}){{end}}
type ChangeStreamForTable struct {
	// pos = TableName.pos
	// end = Rparen + 1 || TableName.end

	Rparen token.Pos // position of ")"

	TableName *Ident
	Columns   []*Ident
}

// ChangeStreamSetFor is SET FOR tables node in ALTER CHANGE STREAM
//
//	SET FOR {{.For | sql}}
type ChangeStreamSetFor struct {
	// pos = Set
	// end = For.end

	Set token.Pos // position of "SET" keyword

	For ChangeStreamFor
}

// ChangeStreamDropForAll is DROP FOR ALL node in ALTER CHANGE STREAM
//
//	DROP FOR ALL
type ChangeStreamDropForAll struct {
	// pos = Drop
	// end = All + 3

	Drop token.Pos // position of "DROP" keyword
	All  token.Pos // position of "ALL" keyword
}

// ChangeStreamSetOptions is SET OPTIONS node in ALTER CHANGE STREAM
//
//	SET {{.Options | sql}}
type ChangeStreamSetOptions struct {
	// pos = Set
	// end = Options.end

	Set token.Pos // position of "SET" keyword

	Options *Options
}

// Storing is STORING clause in CREATE INDEX.
//
//	STORING ({{.Columns | sqlJoin ","}})
type Storing struct {
	// pos = Storing
	// end = Rparen + 1

	Storing token.Pos // position of "STORING" keyword
	Rparen  token.Pos // position of ")"

	Columns []*Ident
}

// InterleaveIn is INTERLEAVE IN clause in CREATE INDEX.
//
//	, INTERLEAVE IN {{.TableName | sql}}
type InterleaveIn struct {
	// pos = Comma
	// end = TableName.end

	Comma token.Pos // position of ","

	TableName *Ident
}

// AddStoredColumn is ADD STORED COLUMN clause in ALTER INDEX.
//
//	ADD STORED COLUMN {{.Name | sql}}
type AddStoredColumn struct {
	// pos = Add
	// end = Name.end

	Add token.Pos // position of "ADD" keyword

	Name *Ident
}

// DropStoredColumn is DROP STORED COLUMN clause in ALTER INDEX.
//
//	DROP STORED COLUMN {{.Name | sql}}
type DropStoredColumn struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// DropIndex is DROP INDEX statement node.
//
//	DROP INDEX {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropIndex struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	IfExists bool
	Name     *Path
}

// DropVectorIndex is DROP VECTOR INDEX statement node.
//
//	DROP VECTOR INDEX {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropVectorIndex struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	IfExists bool
	Name     *Ident
}

// DropSequence is DROP SEQUENCE statement node.
//
//	DROP SEQUENCE {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropSequence struct {
	// pos = Drop
	// end = Name.end

	Drop     token.Pos
	IfExists bool
	Name     *Path
}

// CreateRole is CREATE ROLE statement node.
//
//	CREATE ROLE {{.Name | sql}}
type CreateRole struct {
	// pos = Create
	// end = Name.end

	Create token.Pos // position of "CREATE" keyword

	Name *Ident
}

// DropRole is DROP ROLE statement node.
//
//	DROP ROLE {{.Name | sql}}
type DropRole struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// DropChangeStream is DROP CHANGE STREAM  statement node.
//
//	DROP CHANGE STREAM {{.Name | sql}}
type DropChangeStream struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// Grant is GRANT statement node.
//
//	GRANT {{.Privilege | sql}} TO ROLE {{.Roles | sqlJoin ","}}
type Grant struct {
	// pos = Grant
	// end = Roles[$].end

	Grant token.Pos // position of "GRANT" keyword

	Privilege Privilege
	Roles     []*Ident
}

// Revoke is REVOKE statement node.
//
//	REVOKE {{.Privilege | sql}} FROM ROLE {{.Roles | sqlJoin ","}}
type Revoke struct {
	// pos = Revoke
	// end = Roles[$].end

	Revoke token.Pos // position of "REVOKE" keyword

	Privilege Privilege // len(Privileges) > 0
	Roles     []*Ident  // len(Roles) > 0
}

// PrivilegeOnTable is ON TABLE privilege node in GRANT and REVOKE.
//
//	{{.Privileges | sqlJoin ","}} ON TABLE {{.Names | sqlJoin ","}}
type PrivilegeOnTable struct {
	// pos = Privileges[0].pos
	// end = Names[$].end

	Privileges []TablePrivilege // len(Privileges) > 0
	Names      []*Path          // len(Names) > 0
}

// SelectPrivilege is SELECT ON TABLE privilege node in GRANT and REVOKE.
//
//	SELECT{{if .Columns}}({{.Columns | sqlJoin ","}}){{end}}
type SelectPrivilege struct {
	// pos = Select
	// end = Rparen + 1 || Select + 6

	Select token.Pos // position of "SELECT" keyword
	Rparen token.Pos // position of ")" when len(Columns) > 0

	Columns []*Ident
}

// InsertPrivilege is INSERT ON TABLE privilege node in GRANT and REVOKE.
//
//	INSERT{{if .Columns}}({{.Columns | sqlJoin ","}}){{end}}
type InsertPrivilege struct {
	// pos = Insert
	// end = Rparen + 1 || Insert + 6

	Insert token.Pos // position of "INSERT" keyword
	Rparen token.Pos // position of ")" when len(Columns) > 0

	Columns []*Ident
}

// UpdatePrivilege is UPDATE ON TABLE privilege node in GRANT and REVOKE.
//
//	UPDATE{{if .Columns}}({{.Columns | sqlJoin ","}}){{end}}
type UpdatePrivilege struct {
	// pos = Update
	// end = Rparen + 1 || Update + 6

	Update token.Pos // position of "UPDATE" keyword
	Rparen token.Pos // position of ")" when len(Columns) > 0

	Columns []*Ident
}

// DeletePrivilege is DELETE ON TABLE privilege node in GRANT and REVOKE.
//
//	DELETE
type DeletePrivilege struct {
	// pos = Delete
	// end = Delete + 6

	Delete token.Pos // position of "DELETE" keyword
}

// SelectPrivilegeOnChangeStream is SELECT ON CHANGE STREAM privilege node in GRANT and REVOKE.
//
//	SELECT ON CHANGE STREAM {{.Names | sqlJoin ", "}}
type SelectPrivilegeOnChangeStream struct {
	// pos = Select
	// end = Names[$].end

	Select token.Pos

	Names []*Path // len(Names) > 0
}

// SelectPrivilegeOnView is SELECT ON VIEW privilege node in GRANT and REVOKE.
//
//	SELECT ON VIEW {{.Names | sqlJoin ","}}
type SelectPrivilegeOnView struct {
	// pos = Select
	// end = Names[$].end

	Select token.Pos

	Names []*Path // len(Names) > 0
}

// ExecutePrivilegeOnTableFunction is EXECUTE ON TABLE FUNCTION privilege node in GRANT and REVOKE.
//
//	EXECUTE ON TABLE FUNCTION {{.Names | sqlJoin ","}}
type ExecutePrivilegeOnTableFunction struct {
	// pos = Execute
	// end = Names[$].end

	Execute token.Pos

	Names []*Path // len(Names) > 0
}

// RolePrivilege is ROLE privilege node in GRANT and REVOKE.
//
//	ROLE {{.Names | sqlJoin ","}}
type RolePrivilege struct {
	// pos = Role
	// end = Names[$].end

	Role token.Pos

	Names []*Ident // len(Names) > 0
}

// AlterStatistics is ALTER STATISTICS statement node.
//
//	ALTER STATISTICS {{.Name | sql}} SET {{.Options | sql}}
type AlterStatistics struct {
	// pos = Alter
	// end = Options.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *Options
}

// Analyze is ANALYZE statement node.
//
//	ANALYZE
type Analyze struct {
	// pos = Analyze
	// end = Analyze + 7

	Analyze token.Pos // position of "ANALYZE" keyword
}

// FunctionParam is parameter in CREATE FUNCTION.
//
//	{{.Name | sql}} {{.Type | sql}}{{if .DefaultExpr}} DEFAULT {{.DefaultExpr | sql}}{{end}}
type FunctionParam struct {
	// pos = Name.pos
	// end = DefaultExpr.end || Type.end

	Name *Ident
	Type SchemaType

	DefaultExpr Expr // optional
}

// CreateFunction is CREATE FUNCTION statement node.
//
//	CREATE {{if .OrReplace}}OR REPLACE{{end}} FUNCTION {{.Name | sql}}
//	({{.Params | sqlJoin ", "}}){{if .ReturnType}} RETURNS {{.ReturnType | sql}}{{end}}
//	{{if .SqlSecurity}}SQL SECURITY {{.SqlSecurity}}{{end}}
//	{{if .Determinism}}{{.Determinism}}{{end}}
//	{{if .Language}}LANGUAGE {{.Language}}{{end}}
//	{{if .Remote}}REMOTE{{end}}
//	{{if .Options}}OPTIONS ({{.Options | sqlJoin ", "}}){{end}}
//	{{if .Definition}}AS ({{.Definition | sql}}){{end}}
type CreateFunction struct {
	// pos = Create
	// end = RparenAs + 1 || Options.end

	Create   token.Pos
	As       token.Pos // position of "AS", optional
	RparenAs token.Pos // optional

	OrReplace  bool
	Name       *Path
	Params     []*FunctionParam
	ReturnType SchemaType

	SqlSecurity SecurityType // optional
	Determinism Determinism  // optional
	Language    string       // optional
	Remote      bool
	Options     *Options // optional

	Definition Expr // optional
}

// DropFunction is DROP FUNCTION statement node.
//
//	DROP FUNCTION {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropFunction struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos

	IfExists bool
	Name     *Path
}

// CreateModelColumn is a single column definition node in CREATE MODEL.
//
//	{{.Name | sql}} {{.DataType | sql}} {{.Options | sqlOpt}}
type CreateModelColumn struct {
	// pos = Name.pos
	// end = (Options ?? DataType).end

	Name     *Ident
	DataType SchemaType
	Options  *Options // optional
}

// CreateModelInputOutput is INPUT and OUTPUT column list node.
//
//	INPUT ({{.InputColumns | sqlJoin ", "}}) OUTPUT ({{.OutputColumns | sqlJoin ", "}})
type CreateModelInputOutput struct {
	// pos = Input
	// end = Rparen + 1

	Input  token.Pos
	Rparen token.Pos // position of the last ")"

	InputColumns  []*CreateModelColumn
	OutputColumns []*CreateModelColumn
}

// CreateModel is CREATE MODEL statement node.
//
//	CREATE {{if .OrReplace}}OR REPLACE{{end}} MODEL {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}}
//	{{.InputOutput | sqlOpt}}
//	REMOTE
//	{{.Options | sqlOpt}}
type CreateModel struct {
	// pos = Create
	// end = Options.end || Remote + 6

	Create token.Pos // position of "CREATE" keyword
	Remote token.Pos // position of "REMOTE" keyword

	OrReplace   bool
	IfNotExists bool
	Name        *Ident
	InputOutput *CreateModelInputOutput // optional
	Options     *Options                // optional
}

// AlterModel is ALTER MODEL statement node.
//
//	ALTER MODEL {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}} SET {{.Options | sql}}
type AlterModel struct {
	// pos = Alter
	// end = Options.end

	Alter token.Pos

	IfExists bool
	Name     *Ident
	Options  *Options
}

// DropModel is DROP MODEL statement node.
//
//	DROP MODEL {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropModel struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos

	IfExists bool
	Name     *Ident
}

// ================================================================================
//
// Types for Schema
//
// ================================================================================

// ScalarSchemaType is scalar type node in schema.
// It is used for types without size specification.
// STRING and BYTES can also be ScalarSchemaType if they do not have a size specification.
//
//	{{.Name}}
type ScalarSchemaType struct {
	// pos = NamePos
	// end = NamePos + len(Name)

	NamePos token.Pos // position of this name

	Name ScalarTypeName
}

// SizedSchemaType is sized type node in schema.
//
//	{{.Name}}({{if .Max}}MAX{{else}}{{.Size | sql}}{{end}})
type SizedSchemaType struct {
	// pos = NamePos
	// end = Rparen + 1

	NamePos token.Pos // position of this name
	Rparen  token.Pos // position of ")"

	Name ScalarTypeName // StringTypeName or BytesTypeName
	// either Max or Size must be set
	Max  bool
	Size IntValue
}

// ArraySchemaType is array type node in schema.
//
//	ARRAY<{{.Item | sql}}>{{if .NamedArgs}}({{.NamedArgs | sqlJoin ", "}}){{end}}
type ArraySchemaType struct {
	// pos = Array
	// end = Rparen + 1 || Gt + 1

	Array  token.Pos // position of "ARRAY" keyword
	Gt     token.Pos // position of ">"
	Rparen token.Pos // position of ")" when len(NamedArgs) > 0

	Item      SchemaType // ScalarSchemaType or SizedSchemaType or NamedType
	NamedArgs []*NamedArg
}

// ================================================================================
//
// Search Index DDL
//
// ================================================================================

// CreateSearchIndex represents CREATE SEARCH INDEX statement
//
//	CREATE SEARCH INDEX {{.Name | sql}}
//	ON {{.TableName | sql}}
//	({{.TokenListPart | sqlJoin ", "}})
//	{{.Storing | sqlOpt}}
//	{{if .PartitionColumns}}PARTITION BY {{.PartitionColumns  | sqlJoin ", "}}{{end}}
//	{{.OrderBy | sqlOpt}}
//	{{.Where | sqlOpt}}
//	{{.Interleave | sqlOpt}}
//	{{.Options | sqlOpt}}
type CreateSearchIndex struct {
	// pos = Create
	// end = (Options ?? Interleave ?? Where ?? OrderBy ?? PartitionColumns[$] ?? Storing).end || Rparen + 1

	Create token.Pos

	Name             *Path
	TableName        *Path
	TokenListPart    []*Ident
	Rparen           token.Pos     // position of ")" after TokenListPart
	Storing          *Storing      // optional
	PartitionColumns []*Ident      // optional
	OrderBy          *OrderBy      // optional
	Where            *Where        // optional
	Interleave       *InterleaveIn // optional
	Options          *Options      // optional
}

// DropSearchIndex represents DROP SEARCH INDEX statement.
//
//	DROP SEARCH INDEX{{if .IfExists}}IF EXISTS{{end}} {{Name | sql}}
type DropSearchIndex struct {
	// pos = Drop
	// end = Name.end

	Drop     token.Pos
	IfExists bool
	Name     *Path
}

// AlterSearchIndex represents ALTER SEARCH INDEX statement.
//
//	ALTER SEARCH INDEX {{.Name | sql}} {{.IndexAlteration | sql}}
type AlterSearchIndex struct {
	// pos = Alter
	// end = IndexAlteration.end

	Alter           token.Pos
	Name            *Path
	IndexAlteration IndexAlteration
}

// ================================================================================
//
// GQL schema statements
//
// ================================================================================

// PropertyGraphLabelsOrProperties represents labels with properties or a single properties of node or edge.
type PropertyGraphLabelsOrProperties interface {
	Node
	isPropertyGraphLabelsOrProperties()
}

func (PropertyGraphSingleProperties) isPropertyGraphLabelsOrProperties()       {}
func (PropertyGraphLabelAndPropertiesList) isPropertyGraphLabelsOrProperties() {}

// PropertyGraphElementLabel represents a element label definition.
type PropertyGraphElementLabel interface {
	Node
	isPropertyGraphElementLabel()
}

func (PropertyGraphElementLabelLabelName) isPropertyGraphElementLabel()    {}
func (PropertyGraphElementLabelDefaultLabel) isPropertyGraphElementLabel() {}

// PropertyGraphElementKeys represents PropertyGraphNodeElementKey or PropertyGraphEdgeElementKeys.
type PropertyGraphElementKeys interface {
	Node
	isPropertyGraphElementKeys()
}

func (PropertyGraphNodeElementKey) isPropertyGraphElementKeys()  {}
func (PropertyGraphEdgeElementKeys) isPropertyGraphElementKeys() {}

// PropertyGraphElementProperties represents a definition of properties.
// See https://cloud.google.com/spanner/docs/reference/standard-sql/graph-schema-statements#element_table_property_definition.
type PropertyGraphElementProperties interface {
	Node
	isPropertyGraphElementProperties()
}

func (PropertyGraphNoProperties) isPropertyGraphElementProperties()        {}
func (PropertyGraphPropertiesAre) isPropertyGraphElementProperties()       {}
func (PropertyGraphDerivedPropertyList) isPropertyGraphElementProperties() {}

// CreatePropertyGraph is CREATE PROPERTY GRAPH statement node.
//
//	CREATE {{if .OrReplace}}OR REPLACE{{end}} PROPERTY GRAPH
//	{{if .IfNotExists}}IF NOT EXISTS{{end}}
//	{{.Name | sql}}
//	{{.Content | sql}}
type CreatePropertyGraph struct {
	// pos = Create
	// end = Content.end

	Create      token.Pos // position of "CREATE" keyword
	OrReplace   bool
	IfNotExists bool
	Name        *Ident
	Content     *PropertyGraphContent
}

// PropertyGraphContent represents body of CREATE PROPERTY GRAPH statement.
//
//	NODE TABLES {{.NodeTables | sql}} {{.EdgeTables | sqlOpt}}
type PropertyGraphContent struct {
	// pos = NodeTables.pos
	// end = (EdgeTables ?? NodeTables).end

	NodeTables *PropertyGraphNodeTables
	EdgeTables *PropertyGraphEdgeTables //optional
}

// PropertyGraphNodeTables is NODE TABLES node in CREATE PROPERTY GRAPH statement.
//
//	NODE TABLES {{.Tables | sql}}
type PropertyGraphNodeTables struct {
	// pos = Node
	// end = Tables.end

	Node   token.Pos
	Tables *PropertyGraphElementList
}

// PropertyGraphEdgeTables is EDGE TABLES node in CREATE PROPERTY GRAPH statement.
//
//	EDGE TABLES {{.Tables | sql}}
type PropertyGraphEdgeTables struct {
	// pos = Edge
	// end = Tables.end

	Edge   token.Pos
	Tables *PropertyGraphElementList
}

// PropertyGraphElementList represents element list in NODE TABLES or EDGE TABLES.
//
//	({{.Elements | sqlJoin ", "}})
type PropertyGraphElementList struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos
	Elements       []*PropertyGraphElement
}

// PropertyGraphElement represents a single element in NODE TABLES or EDGE TABLES.
//
//	{{.Name | sql}} {{if .Alias | isnil | not)}}AS {{.Alias | sql}}{{end}}
//	{{.Keys | sqlOpt}} {{.Properties | sqlOpt}}
//	{{.DynamicLabel | sqlOpt}} {{.DynamicProperties | sqlOpt}}
type PropertyGraphElement struct {
	// pos = Name.pos
	// end = (DynamicProperties ?? DynamicLabel ?? Properties ?? Keys ?? Alias ?? Name).end

	Name              *Ident
	Alias             *Ident                          // optional
	Keys              PropertyGraphElementKeys        // optional
	Properties        PropertyGraphLabelsOrProperties // optional
	DynamicLabel      *PropertyGraphDynamicLabel      // optional
	DynamicProperties *PropertyGraphDynamicProperties // optional
}

// PropertyGraphSingleProperties is wrapper node for PropertyGraphElementProperties in PropertyGraphElement.
// It implements PropertyGraphLabelsOrProperties.
//
//	{{.Properties | sql}}
type PropertyGraphSingleProperties struct {
	// pos = Properties.pos
	// end = Properties.end

	Properties PropertyGraphElementProperties
}

// PropertyGraphLabelAndPropertiesList represents whitespace-separated list of PropertyGraphLabelAndProperties.
// It implements PropertyGraphLabelsOrProperties.
//
//	{{.LabelAndProperties | sqlJoin " "}}
type PropertyGraphLabelAndPropertiesList struct {
	// pos = LabelAndProperties[0].pos
	// end = LabelAndProperties[$].end

	LabelAndProperties []*PropertyGraphLabelAndProperties // len(LabelAndProperties) > 0
}

// PropertyGraphLabelAndProperties represents label and properties definition for a single label.
//
//	{{.Label | sql}} {{.Properties | sqlOpt}}
type PropertyGraphLabelAndProperties struct {
	// pos = Label.pos
	// end = (Properties ?? Label).end

	Label      PropertyGraphElementLabel
	Properties PropertyGraphElementProperties // optional
}

// PropertyGraphElementLabelLabelName represents LABEL label_name node.
//
//	LABEL {{.Name | sql}}
type PropertyGraphElementLabelLabelName struct {
	// pos = Label
	// end = Name.end

	Label token.Pos
	Name  *Ident
}

// PropertyGraphElementLabelDefaultLabel represents DEFAULT LABEL node.
//
//	DEFAULT LABEL
type PropertyGraphElementLabelDefaultLabel struct {
	// pos = Default
	// end = Label + 5

	Default token.Pos
	Label   token.Pos
}

// PropertyGraphNodeElementKey is a wrapper of PropertyGraphElementKey to implement PropertyGraphElementKeys
// without deeper AST hierarchy.
//
//	{{.Key | sql}}
type PropertyGraphNodeElementKey struct {
	// pos = Key.pos
	// end = Key.end

	Key *PropertyGraphElementKey
}

// PropertyGraphEdgeElementKeys represents PropertyGraphSourceKey and PropertyGraphDestinationKey with optional PropertyGraphElementKey.
//
//	{{.Element | sqlOpt}} {{.Source | sql}} {{.Destination | sql}}
type PropertyGraphEdgeElementKeys struct {
	// pos = (Element ?? Source).pos
	// end = Destination.end

	Element     *PropertyGraphElementKey // optional
	Source      *PropertyGraphSourceKey
	Destination *PropertyGraphDestinationKey
}

// PropertyGraphElementKey represents the key that identifies the node or edge element.
//
//	KEY {{.Keys | sql}}
type PropertyGraphElementKey struct {
	// pos = Key
	// end = Keys.end

	Key  token.Pos
	Keys *PropertyGraphColumnNameList
}

// PropertyGraphSourceKey represents the key for the source node of the edge.
//
//	SOURCE KEY {{.Keys | sql}}
//	REFERENCES {{.ElementReference | sql}} {{.ReferenceColumns | sqlOpt}}
type PropertyGraphSourceKey struct {
	// pos = Source
	// end = (ReferenceColumns ?? ElementReference).end

	Source           token.Pos
	Keys             *PropertyGraphColumnNameList
	ElementReference *Ident
	ReferenceColumns *PropertyGraphColumnNameList // optional
}

// PropertyGraphDestinationKey represents the key for the destination node of the edge.
//
//	DESTINATION KEY {{.Keys | sql}}
//	REFERENCES {{.ElementReference | sql}} {{.ReferenceColumns | sqlOpt}}
type PropertyGraphDestinationKey struct {
	// pos = Destination
	// end = (ReferenceColumns ?? ElementReference).end

	Destination      token.Pos
	Keys             *PropertyGraphColumnNameList
	ElementReference *Ident
	ReferenceColumns *PropertyGraphColumnNameList // optional
}

// PropertyGraphColumnNameList represents one or more columns to assign to a key.
//
//	({{.ColumnNameList | sqlJoin ", "}})
type PropertyGraphColumnNameList struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos

	ColumnNameList []*Ident
}

// PropertyGraphNoProperties represents the element doesn't have properties.
//
//	NO PROPERTIES
type PropertyGraphNoProperties struct {
	// pos = No
	// end = Properties + 10

	No, Properties token.Pos // position of "NO" and "PROPERTIES"
}

// PropertyGraphPropertiesAre defines which columns to include as element properties.
//
//	PROPERTIES ARE ALL COLUMNS{{if .ExceptColumns | isnil | not}} EXCEPT {{.ExceptColumns | sql}}{{end}}
type PropertyGraphPropertiesAre struct {
	// pos = Properties
	// end = ExceptColumns.end || Columns + 7

	Properties token.Pos // position of "PROPERTIES"
	Columns    token.Pos // position of "COLUMNS"

	ExceptColumns *PropertyGraphColumnNameList // optional
}

// PropertyGraphDerivedPropertyList represents a list of PropertyGraphDerivedProperty.
// NOTE: In current syntax reference, "(" and ")" are missing.
//
//	PROPERTIES ({{.DerivedProperties | sqlJoin ", "}})
type PropertyGraphDerivedPropertyList struct {
	// pos = Properties
	// end = Rparen

	Properties        token.Pos                       // position of "PROPERTIES"
	Rparen            token.Pos                       // position of ")"
	DerivedProperties []*PropertyGraphDerivedProperty // len(DerivedProperties) > 0
}

// PropertyGraphDerivedProperty represents an expression that defines a property and can optionally reference the input table columns.
//
//	{{.Expr | sql}} {{if .Alias}}AS {{.Alias | sql}}{{end}}
type PropertyGraphDerivedProperty struct {
	// pos = Expr.pos
	// end = (Alias ?? Expr).end

	Expr  Expr
	Alias *Ident // optional
}

// DynamicLabel represents DYNAMIC LABEL clause in CREATE PROPERTY GRAPH statement.
//
//	DYNAMIC LABEL ({{.ColumnName | sql}})
type PropertyGraphDynamicLabel struct {
	// pos = Dynamic
	// end = Rparen + 1

	Dynamic token.Pos // position of "DYNAMIC"
	Rparen  token.Pos // position of ")"

	ColumnName *Ident
}

// DynamicProperties represents DYNAMIC PROPERTIES clause in CREATE PROPERTY GRAPH statement.
//
//	DYNAMIC PROPERTIES ({{.ColumnName | sql}})
type PropertyGraphDynamicProperties struct {
	// pos = Dynamic
	// end = Rparen + 1

	Dynamic token.Pos // position of "DYNAMIC"
	Rparen  token.Pos // position of ")"

	ColumnName *Ident
}

// DropPropertyGraph is DROP PROPERTY GRAPH statement node.
//
//	DROP PROPERTY GRAPH {{if .IfExists}}IF EXISTS{{end}} {{.Name | sql}}
type DropPropertyGraph struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos

	IfExists bool
	Name     *Ident
}

// ================================================================================
//
// DML
//
// ================================================================================

// WithAction is WITH ACTION clause in ThenReturn.
//
//	WITH ACTION {{.Alias | sqlOpt}}
type WithAction struct {
	// pos = With
	// end = Alias.end || Action + 6

	With   token.Pos // position of "WITH" keyword
	Action token.Pos // position of "ACTION" keyword

	Alias *AsAlias // optional
}

// ThenReturn is THEN RETURN clause in DML.
//
//	THEN RETURN {{.WithAction | sqlOpt}} {{.Items | sqlJoin ", "}}
type ThenReturn struct {
	// pos = Then
	// end = Items[$].end

	Then token.Pos // position of "THEN" keyword

	WithAction *WithAction // optional
	Items      []SelectItem
}

// Insert is INSERT statement node.
//
//	{{.Hint | sqlOpt}}
//	INSERT {{if .InsertOrType}}OR .InsertOrType{{end}}INTO {{.TableName | sql}}{{.TableHint | sqlOpt}} ({{.Columns | sqlJoin ","}}) {{.Input | sql}}
//	{{.ThenReturn | sqlOpt}}
type Insert struct {
	// pos = Hint.pos || Insert
	// end = (ThenReturn ?? Input).end

	Insert token.Pos // position of "INSERT" keyword

	InsertOrType InsertOrType

	Hint       *Hint // optional
	TableName  *Path
	TableHint  *Hint // optional
	Columns    []*Ident
	Input      InsertInput
	ThenReturn *ThenReturn // optional
}

// ValuesInput is VALUES clause in INSERT.
//
//	VALUES {{.Rows | sqlJoin ","}}
type ValuesInput struct {
	// pos = Values
	// end = Rows[$].end

	Values token.Pos // position of "VALUES" keyword

	Rows []*ValuesRow
}

// ValuesRow is row values of VALUES clause.
//
//	({{.Exprs | sqlJoin ","}})
type ValuesRow struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Exprs []*DefaultExpr
}

// DefaultExpr is DEFAULT or Expr.
//
//	{{if .Default}}DEFAULT{{else}}{{.Expr | sql}}{{end}}
type DefaultExpr struct {
	// pos = DefaultPos || Expr.pos
	// end = DefaultPos + 7 || Expr.end

	DefaultPos token.Pos // position of "DEFAULT"

	Default bool
	Expr    Expr
}

// SubQueryInput is query clause in INSERT.
//
//	{{.Query | sql}}
type SubQueryInput struct {
	// pos = Query.pos
	// end = Query.end

	Query QueryExpr
}

// Delete is DELETE statement.
//
//	{{.Hint | sqlOpt}}
//	DELETE FROM {{.TableName | sql}}{{.TableHint | sqlOpt}} {{.As | sqlOpt}} {{.Where | sql}}
//	{{.ThenReturn | sqlOpt}}
type Delete struct {
	// pos = Hint.pos || Delete
	// end = (ThenReturn ?? Where).end

	Delete token.Pos // position of "DELETE" keyword

	Hint       *Hint // optional
	TableName  *Path
	TableHint  *Hint    // optional
	As         *AsAlias // optional
	Where      *Where
	ThenReturn *ThenReturn // optional
}

// Update is UPDATE statement.
//
//	{{.Hint | sqlOpt}}
//	UPDATE {{.TableName | sql}}{{.TableHint | sqlOpt}} {{.As | sqlOpt}}
//	SET {{.Updates | sqlJoin ","}} {{.Where | sql}}
//	{{.ThenReturn | sqlOpt}}
type Update struct {
	// pos = Hint.pos || Update
	// end = (ThenReturn ?? Where).end

	Update token.Pos // position of "UPDATE" keyword

	Hint       *Hint // optional
	TableName  *Path
	TableHint  *Hint         // optional
	As         *AsAlias      // optional
	Updates    []*UpdateItem // len(Updates) > 0
	Where      *Where
	ThenReturn *ThenReturn // optional
}

// UpdateItem is SET clause items in UPDATE.
//
//	{{.Path | sqlJoin "."}} = {{.DefaultExpr | sql}}
type UpdateItem struct {
	// pos = Path[0].pos
	// end = DefaultExpr.end

	Path        []*Ident // len(Path) > 0
	DefaultExpr *DefaultExpr
}

// ================================================================================
//
// Procedural language
//
// ================================================================================

// Call is CALL statement.
//
//	CALL {{.Name | sql}}({{.Args | sqlJoin ", "}})
type Call struct {
	// pos = Call
	// end = Rparen +1

	Call   token.Pos
	Rparen token.Pos

	Name *Path
	Args []TVFArg // len(Args) > 0
}
