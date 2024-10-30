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

// This file must contain only AST definitions.
// We use the following go:generate directive for generating pos.go. Thus, all AST definitions must have pos and end lines.
//go:generate go run ../tools/gen-ast-pos/main.go -infile ast.go -outfile pos.go

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

func (QueryStatement) isStatement()     {}
func (CreateSchema) isStatement()       {}
func (DropSchema) isStatement()         {}
func (CreateDatabase) isStatement()     {}
func (AlterDatabase) isStatement()      {}
func (CreateTable) isStatement()        {}
func (AlterTable) isStatement()         {}
func (DropTable) isStatement()          {}
func (RenameTable) isStatement()        {}
func (CreateIndex) isStatement()        {}
func (AlterIndex) isStatement()         {}
func (DropIndex) isStatement()          {}
func (CreateSearchIndex) isStatement()  {}
func (DropSearchIndex) isStatement()    {}
func (AlterSearchIndex) isStatement()   {}
func (CreateView) isStatement()         {}
func (DropView) isStatement()           {}
func (CreateChangeStream) isStatement() {}
func (AlterChangeStream) isStatement()  {}
func (DropChangeStream) isStatement()   {}
func (CreateRole) isStatement()         {}
func (DropRole) isStatement()           {}
func (Grant) isStatement()              {}
func (Revoke) isStatement()             {}
func (CreateSequence) isStatement()     {}
func (AlterSequence) isStatement()      {}
func (DropSequence) isStatement()       {}
func (AlterStatistics) isStatement()    {}
func (CreateVectorIndex) isStatement()  {}
func (DropVectorIndex) isStatement()    {}
func (Insert) isStatement()             {}
func (Delete) isStatement()             {}
func (Update) isStatement()             {}

// GRAPH query is top level statement which can be executed by ExecuteSQL API.
func (*GQLGraphQuery) isStatement() {}

// QueryExpr represents set operator operands.
type QueryExpr interface {
	Node
	isQueryExpr()
}

func (Select) isQueryExpr()        {}
func (SubQuery) isQueryExpr()      {}
func (CompoundQuery) isQueryExpr() {}

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
func (*GraphTableExpr) isTableExpr()   {}

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

func (BinaryExpr) isExpr()            {}
func (UnaryExpr) isExpr()             {}
func (InExpr) isExpr()                {}
func (IsNullExpr) isExpr()            {}
func (IsBoolExpr) isExpr()            {}
func (IsSourceExpr) isExpr()          {}
func (IsDestinationExpr) isExpr()     {}
func (BetweenExpr) isExpr()           {}
func (SelectorExpr) isExpr()          {}
func (IndexExpr) isExpr()             {}
func (CallExpr) isExpr()              {}
func (CountStarExpr) isExpr()         {}
func (CastExpr) isExpr()              {}
func (ExtractExpr) isExpr()           {}
func (CaseExpr) isExpr()              {}
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
func (NewConstructor) isExpr()        {}
func (BracedNewConstructor) isExpr()  {}
func (BracedConstructor) isExpr()     {}
func (ArrayGQLSubQuery) isExpr()      {}
func (ValueGQLSubQuery) isExpr()      {}
func (ExistsGQLSubQuery) isExpr()     {}

// Arg represents argument of function call.
type Arg interface {
	Node
	isArg()
}

func (ExprArg) isArg()     {}
func (IntervalArg) isArg() {}
func (SequenceArg) isArg() {}

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

func (UnnestInCondition) isInCondition()       {}
func (SubQueryInCondition) isInCondition()     {}
func (ValuesInCondition) isInCondition()       {}
func (*GQLSubQueryInCondition) isInCondition() {}

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

func (CreateSchema) isDDL()       {}
func (DropSchema) isDDL()         {}
func (CreateDatabase) isDDL()     {}
func (AlterDatabase) isDDL()      {}
func (CreateTable) isDDL()        {}
func (AlterTable) isDDL()         {}
func (DropTable) isDDL()          {}
func (RenameTable) isDDL()        {}
func (CreateIndex) isDDL()        {}
func (AlterIndex) isDDL()         {}
func (DropIndex) isDDL()          {}
func (CreateView) isDDL()         {}
func (CreateSearchIndex) isDDL()  {}
func (DropSearchIndex) isDDL()    {}
func (AlterSearchIndex) isDDL()   {}
func (DropView) isDDL()           {}
func (CreateChangeStream) isDDL() {}
func (AlterChangeStream) isDDL()  {}
func (DropChangeStream) isDDL()   {}
func (CreateRole) isDDL()         {}
func (DropRole) isDDL()           {}
func (Grant) isDDL()              {}
func (Revoke) isDDL()             {}
func (CreateSequence) isDDL()     {}
func (AlterSequence) isDDL()      {}
func (DropSequence) isDDL()       {}
func (AlterStatistics) isDDL()    {}
func (CreateVectorIndex) isDDL()  {}
func (DropVectorIndex) isDDL()    {}

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
func (AlterColumn) isTableAlteration()              {}

// ColumnAlteration represents ALTER COLUMN action in ALTER TABLE.
type ColumnAlteration interface {
	Node
	isColumnAlteration()
}

func (AlterColumnType) isColumnAlteration()        {}
func (AlterColumnSetOptions) isColumnAlteration()  {}
func (AlterColumnSetDefault) isColumnAlteration()  {}
func (AlterColumnDropDefault) isColumnAlteration() {}

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
func (NamedType) isSchemaType()        {}

// IndexAlteration represents ALTER INDEX action.
type IndexAlteration interface {
	Node
	isIndexAlteration()
}

func (AddStoredColumn) isIndexAlteration()  {}
func (DropStoredColumn) isIndexAlteration() {}

// DML represents data manipulation language in SQL.
//
// https://cloud.google.com/spanner/docs/data-definition-language
type DML interface {
	Statement
	isDML()
}

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
// SELECT
//
// ================================================================================

// QueryStatement is query statement node.
//
//	{{.Hint | sqlOpt}} {{.With | sqlOpt}} {{.Query | sql}}
//
// https://cloud.google.com/spanner/docs/query-syntax
type QueryStatement struct {
	// pos = (Hint ?? With ?? Query).pos
	// end = Query.end

	Hint  *Hint // optional
	With  *With // optional
	Query QueryExpr
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

	Key   *Ident
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
//	  {{if .Distinct}}DISTINCT{{end}}
//	  {{.As | sqlOpt}}
//	  {{.Results | sqlJoin ","}}
//	  {{.From | sqlOpt}}
//	  {{.Where | sqlOpt}}
//	  {{.GroupBy | sqlOpt}}
//	  {{.Having | sqlOpt}}
//	  {{.OrderBy | sqlOpt}}
//	  {{.Limit | sqlOpt}}
type Select struct {
	// pos = Select
	// end = (Limit ?? OrderBy ?? Having ?? GroupBy ?? Where ?? From ?? Results[$]).end

	Select token.Pos // position of "select" keyword

	Distinct bool
	As       SelectAs     // optional
	Results  []SelectItem // len(Results) > 0
	From     *From        // optional
	Where    *Where       // optional
	GroupBy  *GroupBy     // optional
	Having   *Having      // optional
	OrderBy  *OrderBy     // optional
	Limit    *Limit       // optional
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

// CompoundQuery is query statement node compounded by set operators.
//
//	{{.Queries | sqlJoin (printf "%s %s" .Op or(and(.Distinct, "DISTINCT"), "ALL"))}}
//	  {{.OrderBy | sqlOpt}}
//	  {{.Limit | sqlOpt}}
type CompoundQuery struct {
	// pos = Queries[0].pos
	// end = (Limit ?? OrderBy ?? Queries[$]).end

	Op       SetOp
	Distinct bool
	Queries  []QueryExpr // len(Queries) >= 2
	OrderBy  *OrderBy    // optional
	Limit    *Limit      // optional
}

// SubQuery is subquery statement node.
//
//	({{.Query | sql}}) {{.OrderBy | sqlOpt}} {{.Limit | sqlOpt}}
type SubQuery struct {
	// pos = Lparen
	// end = (Limit ?? OrderBy).end || Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query   QueryExpr
	OrderBy *OrderBy // optional
	Limit   *Limit   // optional
}

// Star is a single * in SELECT result columns list.
//
//	*
type Star struct {
	// pos = Star
	// end = Star + 1

	Star token.Pos // position of "*"
}

// DotStar is expression with * in SELECT result columns list.
//
//	{{.Expr | sql}}.*
type DotStar struct {
	// pos = Expr.pos
	// end = Star + 1

	Star token.Pos // position of "*"

	Expr Expr
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
//	GROUP BY {{.Exprs | sqlJoin ","}}
type GroupBy struct {
	// pos = Group
	// end = Exprs[$].end

	Group token.Pos // position of "GROUP" keyword

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

// GraphTableExpr is GRAPH_TABLE operator
//
//	GRAPH_TABLE({{.PropertyGraphName | sql}} {{.Query | sql}}) {{.As | sqlOpt}}
type GraphTableExpr struct {
	// pos = GraphTable
	// end = As.end || Rparen + 1

	GraphTable        token.Pos
	PropertyGraphName *Ident

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query *GQLMultiLinearQueryStatement
	As    *AsAlias // optional
}

// Join is JOIN expression.
//
//	{{.Left | sql}}
//	  {{.Op}} {{.Method}} {{.Hint | sqlOpt}}
//	  {{.Right | sql}}
//	{{.Cond | sqlOpt}}
type Join struct {
	// pos = Left.pos
	// end = (Cond ?? Right).pos

	Op          JoinOp
	Method      JoinMethod
	Hint        *Hint // optional
	Left, Right TableExpr
	Cond        JoinCondition // nil when Op is CrossJoin, otherwise it must be set.
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
	// end = Right.pos

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

// GQLSubQueryInCondition is GQL subquery at IN condition.
//
//	{{"{"}}{{.Query | sql}}}{{"}"}}
type GQLSubQueryInCondition struct {
	// pos = LBrace
	// end = RBrace + 1

	LBrace, RBrace token.Pos // position of "{" and "}"

	Query *GQLQueryExpr
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

// IsSourceExpr is IS SOURCE expression node.
//
//	{{.Node | sql}} IS {{if .Not}}NOT{{end}} SOURCE OF {{.Edge | sql}}
type IsSourceExpr struct {
	// pos = Node.pos
	// end = Edge.end

	Node, Edge Expr
	Not        bool
}

// IsDestinationExpr is IS DESTINATION expression node.
//
//	{{.Node | sql}} IS {{if .Not}}NOT{{end}} DESTINATION OF {{.Edge | sql}}
type IsDestinationExpr struct {
	// pos = Node.pos
	// end = Edge.end

	Node, Edge Expr
	Not        bool
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
	// end = Ident.pos

	Expr  Expr
	Ident *Ident
}

// IndexExpr is array item access expression node.
//
//	{{.Expr | sql}}[{{if .Ordinal}}ORDINAL{{else}}OFFSET{{end}}({{.Index | sql}})]
type IndexExpr struct {
	// pos = Expr.pos
	// end = Rbrack + 1

	Rbrack token.Pos // position of "]"

	Ordinal     bool
	Expr, Index Expr
}

// CallExpr is function call expression node.
//
//	{{.Func | sql}}(
//		{{if .Distinct}}DISTINCT{{end}}
//		{{.Args | sqlJoin ", "}}
//		{{if len(.Args) > 0 && len(.NamedArgs) > 0}}, {{end}}
//		{{.NamedArgs | sqlJoin ", "}}
//		{{.NullHandling | sqlOpt}}
//		{{.Having | sqlOpt}}
//	)
type CallExpr struct {
	// pos = Func.pos
	// end = Rparen + 1

	Rparen token.Pos // position of ")"

	Func         *Ident
	Distinct     bool
	Args         []Arg
	NamedArgs    []*NamedArg
	NullHandling NullHandlingModifier // optional
	Having       HavingModifier       // optional
}

// ExprArg is argument of the generic function call.
//
//	{{.Expr | sql}}
type ExprArg struct {
	// pos = Expr.pos
	// end = Expr.end

	Expr Expr
}

// IntervalArg is argument of date function call.
//
//	INTERVAL {{.Expr | sql}} {{.Unit | sqlOpt}}
type IntervalArg struct {
	// pos = Interval
	// end = (Unit ?? Expr).end

	Interval token.Pos // position of "INTERVAL" keyword

	Expr Expr
	Unit *Ident // optional
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

// AtTimeZone is AT TIME ZONE clause in EXTRACT call.
//
//	AT TIME ZONE {{.Expr | sql}}
type AtTimeZone struct {
	// pos = At
	// end = Expr.end

	At token.Pos // position of "AT" keyword

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

// ArrayGQLSubQuery is GQL subquery as ARRAY.
//
//	ARRAY {{"{"}}{{.Query | sql}}{{"}"}}
type ArrayGQLSubQuery struct {
	// pos = Array
	// end = RBrace + 1

	Array  token.Pos // position of "ARRAY" keyword
	RBrace token.Pos // position of "}"
	Query  *GQLQueryExpr
}

// ValueGQLSubQuery is GQL subquery as VALUE.
//
//	VALUE {{"{"}}{{.Query | sql}}{{"}"}}
type ValueGQLSubQuery struct {
	// pos = Array
	// end = RBrace + 1

	Array  token.Pos // position of "ARRAY" keyword
	RBrace token.Pos // position of "}"
	Query  *GQLQueryExpr
}

// ExistsGQLSubQuery is GQL subquery as EXISTS.
//
//	EXISTS{{"{"}}{{.Expr | sql}}{{"}"}}
type ExistsGQLSubQuery struct {
	// pos = Exists
	// end = RBrace + 1

	Exists token.Pos // position of "EXISTS" keyword
	RBrace token.Pos // "}"
	Query  GQLExistsExpr
}

type GQLExistsExpr interface {
	Node
	isGQLExistsExpr()
}

func (*GQLQueryExpr) isGQLExistsExpr()      {}
func (*GQLMatchStatement) isGQLExistsExpr() {}
func (*GQLGraphPattern) isGQLExistsExpr()   {}

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
//
//	{{.Idents | sqlJoin "."}}
type Path struct {
	// pos = Idents[0].pos
	// end = Idents[$].end

	Idents []*Ident // len(Idents) >= 2
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

func (*BracedConstructor) isBracedConstructorFieldValue()               {}
func (*BracedConstructorFieldValueExpr) isBracedConstructorFieldValue() {}

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
//	CREATE SCHEMA {{.Name | sql}}
type CreateSchema struct {
	// pos = Create
	// end = Name.end

	Create token.Pos // position of "CREATE" keyword

	Name *Ident
}

// DropSchema is DROP SCHEMA statement node.
//
//	DROP SCHEMA {{.Name | sql}}
type DropSchema struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
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
	// end = Name.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *Options
}

// CreateTable is CREATE TABLE statement node.
//
//	CREATE TABLE {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}} (
//	  {{.Columns | sqlJoin ","}}{{if and .Columns (or .TableConstrains .Synonym)}},{{end}}
//	  {{.TableConstraints | sqlJoin ","}}{{if and .TableConstraints .Synonym}},{{end}}
//	  {{.Synonym | sqlJoin ","}}
//	)
//	PRIMARY KEY ({{.PrimaryKeys | sqlJoin ","}})
//	{{.Cluster | sqlOpt}}
//	{{.CreateRowDeletionPolicy | sqlOpt}}
//
// Spanner SQL allows to mix `Columns` and `TableConstraints` and `Synonyms`,
// however they are separated in AST definition for historical reasons. If you want to get
// the original order of them, please sort them by their `Pos()`.
type CreateTable struct {
	// pos = Create
	// end = RowDeletionPolicy.end || Cluster.end || Rparen + 1

	Create token.Pos // position of "CREATE" keyword
	Rparen token.Pos // position of ")" of PRIMARY KEY clause

	IfNotExists       bool
	Name              *Ident
	Columns           []*ColumnDef
	TableConstraints  []*TableConstraint
	PrimaryKeys       []*IndexKey
	Synonyms          []*Synonym
	Cluster           *Cluster                 // optional
	RowDeletionPolicy *CreateRowDeletionPolicy // optional
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
//	CREATE SEQUENCE {{if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}} }} {{.Options | sql}}
type CreateSequence struct {
	// pos = Create
	// end = Options.end

	Create token.Pos // position of "CREATE" keyword

	Name        *Ident
	IfNotExists bool
	Options     *Options
}

// ColumnDef is column definition in CREATE TABLE.
//
//	{{.Name | sql}}
//	{{.Type | sql}} {{if .NotNull}}NOT NULL{{end}}
//	{{.DefaultExpr | sqlOpt}}
//	{{.GeneratedExpr | sqlOpt}}
//	{{if .Hidden.Invalid | not)}}HIDDEN{{end}}
//	{{.Options | sqlOpt}}
type ColumnDef struct {
	// pos = Name.pos
	// end = Options.end || Hidden + 6 || GeneratedExpr.end || DefaultExpr.end || Null + 4 || Type.end

	Null token.Pos // position of "NULL"

	Name          *Ident
	Type          SchemaType
	NotNull       bool
	DefaultExpr   *ColumnDefaultExpr   // optional
	GeneratedExpr *GeneratedColumnExpr // optional
	Hidden        token.Pos            // InvalidPos if not hidden
	Options       *Options             // optional
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

// ColumnDefOption is options for column definition.
//
//	OPTIONS(allow_commit_timestamp = {{if .AllowCommitTimestamp}}true{{else}null{{end}}})
type ColumnDefOptions struct {
	// pos = Options
	// end = Rparen + 1

	Options token.Pos // position of "OPTIONS" keyword
	Rparen  token.Pos // position of ")"

	AllowCommitTimestamp bool
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
//	FOREIGN KEY ({{.ColumnNames | sqlJoin ","}}) REFERENCES {{.ReferenceTable}} ({{.ReferenceColumns | sqlJoin ","}}) {{.OnDelete}}
type ForeignKey struct {
	// pos = Foreign
	// end = OnDeleteEnd || Rparen + 1

	Foreign     token.Pos // position of "FOREIGN" keyword
	Rparen      token.Pos // position of ")" after reference columns
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	Columns          []*Ident
	ReferenceTable   *Ident
	ReferenceColumns []*Ident       // len(ReferenceColumns) > 0
	OnDelete         OnDeleteAction // optional
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

// Cluster is INTERLEAVE IN PARENT clause in CREATE TABLE.
//
//	, INTERLEAVE IN PARENT {{.TableName | sql}} {{.OnDelete}}
type Cluster struct {
	// pos = Comma
	// end = OnDeleteEnd || TableName.end

	Comma       token.Pos // position of ","
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	TableName *Ident
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

	Name         *Ident
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

	Name *Ident
}

// AlterTable is ALTER TABLE statement node.
//
//	ALTER TABLE {{.Name | sql}} {{.TableAlteration | sql}}
type AlterTable struct {
	// pos = Alter
	// end = TableAlteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name            *Ident
	TableAlteration TableAlteration
}

// AlterIndex is ALTER INDEX statement node.
//
//	ALTER INDEX {{.Name | sql}} {{.IndexAlteration | sql}}
type AlterIndex struct {
	// pos = Alter
	// end = IndexAlteration.end

	Alter token.Pos // position of "ALTER" keyword

	Name            *Ident
	IndexAlteration IndexAlteration
}

// AlterSequence is ALTER SEQUENCE statement node.
type AlterSequence struct {
	// pos = Alter
	// end = Options.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *Options
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

// DropTable is DROP TABLE statement node.
//
//	DROP TABLE {{if .IfExists}}IF NOT EXISTS{{end}} {{.Name | sql}}
type DropTable struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	IfExists bool
	Name     *Ident
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
type CreateIndex struct {
	// pos = Create
	// end = (InterleaveIn ?? Storing).end || Rparen + 1

	Create token.Pos // position of "CREATE" keyword
	Rparen token.Pos // position of ")"

	Unique       bool
	NullFiltered bool
	IfNotExists  bool
	Name         *Ident
	TableName    *Ident
	Keys         []*IndexKey
	Storing      *Storing      // optional
	InterleaveIn *InterleaveIn // optional
}

// CreateVectorIndex is CREATE VECTOR INDEX statement node.
//
//	CREATE VECTOR INDEX {if .IfNotExists}}IF NOT EXISTS{{end}} {{.Name | sql}}
//	ON {{.TableName | sql}}({{.ColumnName | sql}})
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

	// It only allows `WHERE column_name IS NOT NULL` for now, but we still relax the condition
	// by reusing the `parseWhere` function for sake of it may be extended more conditions in the future.
	//
	// Reference: https://cloud.google.com/spanner/docs/reference/standard-sql/data-definition-language#vector_index_statements
	Where   *Where // optional
	Options *Options
}

// VectorIndexOption is OPTIONS record node.
//
//	{{.Key | sql}}={{.Expr | sql}}
type VectorIndexOption struct {
	// pos = Key.pos
	// end = Value.end

	Key   *Ident
	Value Expr
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
	// end = All

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
//	{{.TableName | sql}}{{if .Columns}}({{.Columns | sqlJoin ","}}){{end}}
type ChangeStreamForTable struct {
	// pos = TableName.pos
	// end = TableName.end || Rparen + 1

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
	Name     *Ident
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
	Name     *Ident
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
	Names      []*Ident         // len(Names) > 0
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

	Names []*Ident // len(Names) > 0
}

// SelectPrivilegeOnView is SELECT ON VIEW privilege node in GRANT and REVOKE.
//
//	SELECT ON VIEW {{.Names | sqlJoin ","}}
type SelectPrivilegeOnView struct {
	// pos = Select
	// end = Names[$].end

	Select token.Pos

	Names []*Ident // len(Names) > 0
}

// ExecutePrivilegeOnTableFunction is EXECUTE ON TABLE FUNCTION privilege node in GRANT and REVOKE.
//
//	EXECUTE ON TABLE FUNCTION {{.Names | sqlJoin ","}}
type ExecutePrivilegeOnTableFunction struct {
	// pos = Execute
	// end = Names[$].end

	Execute token.Pos

	Names []*Ident // len(Names) > 0
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

// ================================================================================
//
// Types for Schema
//
// ================================================================================

// ScalarSchemaType is scalar type node in schema.
//
//	{{.Name}}
type ScalarSchemaType struct {
	// pos = NamePos
	// end = NamePos + len(Name)

	NamePos token.Pos // position of this name

	Name ScalarTypeName // except for StringTypeName and BytesTypeName
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
//	ARRAY<{{.Item | sql}}>
type ArraySchemaType struct {
	// pos = Array
	// end = Gt + 1

	Array token.Pos // position of "ARRAY" keyword
	Gt    token.Pos // position of ">"

	Item SchemaType // ScalarSchemaType or SizedSchemaType
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

	Name             *Ident
	TableName        *Ident
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
	Name     *Ident
}

// AlterSearchIndex represents ALTER SEARCH INDEX statement.
//
//	ALTER SEARCH INDEX {{.Name | sql}} {{.IndexAlteration | sql}}
type AlterSearchIndex struct {
	// pos = Alter
	// end = IndexAlteration.end

	Alter           token.Pos
	Name            *Ident
	IndexAlteration IndexAlteration
}

// ================================================================================
//
// DML
//
// ================================================================================

// Insert is INSERT statement node.
//
//	INSERT {{if .InsertOrType}}OR .InsertOrType{{end}}INTO {{.TableName | sql}} ({{.Columns | sqlJoin ","}}) {{.Input | sql}}
type Insert struct {
	// pos = Insert
	// end = Input.end

	Insert token.Pos // position of "INSERT" keyword

	InsertOrType InsertOrType

	TableName *Ident
	Columns   []*Ident
	Input     InsertInput
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
//	DELETE FROM {{.TableName | sql}} {{.As | sqlOpt}} {{.Where | sql}}
type Delete struct {
	// pos = Delete
	// end = Where.end

	Delete token.Pos // position of "DELETE" keyword

	TableName *Ident
	As        *AsAlias // optional
	Where     *Where
}

// Update is UPDATE statement.
//
//	UPDATE {{.TableName | sql}} {{.As | sqlOpt}}
//	SET {{.Updates | sqlJoin ","}} {{.Where | sql}}
type Update struct {
	// pos = Update
	// end = Where.end

	Update token.Pos // position of "UPDATE" keyword

	TableName *Ident
	As        *AsAlias      // optional
	Updates   []*UpdateItem // len(Updates) > 0
	Where     *Where
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
// GQL
//
// https://cloud.google.com/spanner/docs/reference/standard-sql/graph-query-statements
//
// ================================================================================

// GQLGraphQuery is toplevel node of GRAPH query.
//
//	{{.Graph | sql}}
//	{{.MultiLinearQueryStatement | sql}}
type GQLGraphQuery struct {
	// pos = (GraphClause ?? MultiLinearQueryStatement).pos
	// end = MultiLinearQueryStatement.end

	GraphClause               *GQLGraphClause
	MultiLinearQueryStatement *GQLMultiLinearQueryStatement
}

// GQLQueryExpr is similar to GQLGraphQuery,
// but it is appeared in GQL subqueries and it can optionally have GRAPH clause
//
//	{{.Graph | sqlOpt}}
//	{{.MultiLinearQueryStatement | sql}}
type GQLQueryExpr struct {
	// pos = (GraphClause ?? MultiLinearQueryStatement).pos
	// end = MultiLinearQueryStatement.end

	GraphClause               *GQLGraphClause // optional
	MultiLinearQueryStatement *GQLMultiLinearQueryStatement
}

// GQLGraphClause represents `GRAPH property_graph_name`.
//
//	GRAPH {{.PropertyGraphName | sql}}
type GQLGraphClause struct {
	// pos = Graph
	// end = PropertyGraphName.end

	Graph             token.Pos
	PropertyGraphName *Ident
}

// GQLMultiLinearQueryStatement is the body of a GQLGraphClause and GQLQueryExpr.
// It contains a list of LinearQueryStatementList chained together with the NEXT statement.
//
//	{{.LinearQueryStatementList || sqlJoin "\nNEXT\n"}}
type GQLMultiLinearQueryStatement struct {
	// pos = LinearQueryStatementList[0].pos
	// end = LinearQueryStatementList[$].end

	LinearQueryStatementList []GQLLinearQueryStatement
}

type GQLLinearQueryStatement interface {
	Node
	isGQLLinearQueryStatement()
}

func (*GQLSimpleLinearQueryStatement) isGQLLinearQueryStatement()    {}
func (*GQLCompositeLinearQueryStatement) isGQLLinearQueryStatement() {}

// GQLSimpleLinearQueryStatement represents a list of primitive_query_statements that ends with a RETURN statement.
//
//	{{.PrimitiveQueryStatementList | sqlJoin "\n"}}
type GQLSimpleLinearQueryStatement struct {
	// pos = PrimitiveQueryStatementList[0].pos
	// end = PrimitiveQueryStatementList[$].end

	// It contains at least one GQL statements, and It ends with a RETURN statement.
	PrimitiveQueryStatementList []GQLPrimitiveQueryStatement
}

// GQLSimpleLinearQueryStatementWithSetOperator represents GQLSimpleLinearQueryStatement composited with the set operators.
//
// // TODO: {{string(SetOperator)}}
type GQLSimpleLinearQueryStatementWithSetOperator struct {
	// pos = StartPos
	// end = Statement.end

	StartPos      token.Pos
	SetOperator   GQLSetOperatorEnum
	DistinctOrAll GQLAllOrDistinctEnum
	Statement     *GQLSimpleLinearQueryStatement
}

// GQLCompositeLinearQueryStatement represents a list of GQLSimpleLinearQueryStatement composited with the set operators.
//
// {{.HeadSimpleLinearQueryStatement | sql}}
// {{.TailSimpleLinearQueryStatementList | sqlJoin "\n"}}
type GQLCompositeLinearQueryStatement struct {
	// pos = HeadSimpleLinearQueryStatement.pos
	// end = TailSimpleLinearQueryStatementList[$].end

	HeadSimpleLinearQueryStatement     *GQLSimpleLinearQueryStatement
	TailSimpleLinearQueryStatementList []*GQLSimpleLinearQueryStatementWithSetOperator
}

// ================================================================================
//
// GQL statements
//
// ================================================================================

type GQLPrimitiveQueryStatement interface {
	Node
	isGQLPrimitiveQueryStatement()
}

func (*GQLWithStatement) isGQLPrimitiveQueryStatement()    {}
func (*GQLOrderByStatement) isGQLPrimitiveQueryStatement() {}
func (*GQLOffsetStatement) isGQLPrimitiveQueryStatement()  {}
func (*GQLLimitStatement) isGQLPrimitiveQueryStatement()   {}
func (*GQLForStatement) isGQLPrimitiveQueryStatement()     {}
func (*GQLFilterStatement) isGQLPrimitiveQueryStatement()  {}
func (*GQLMatchStatement) isGQLPrimitiveQueryStatement()   {}
func (*GQLLetStatement) isGQLPrimitiveQueryStatement()     {}
func (*GQLReturnStatement) isGQLPrimitiveQueryStatement()  {}

// GQLMatchStatement represents MATCH statement.
//
//	{{if .Optional.Invalid | not}}OPTIONAL {{end-}}
//	MATCH //	{{.MatchHint | sqlOpt}} //	{{.PrefixOrMode | sqlOpt}} {{.GraphPattern | sql}}
type GQLMatchStatement struct {
	// pos = Optional || Match
	// end = GraphPattern.end

	Optional token.Pos //optional
	Match    token.Pos

	MatchHint    *Hint                         // optional
	PrefixOrMode GQLPathSearchPrefixOrPathMode // optional
	GraphPattern *GQLGraphPattern
}

type GQLLimitAndOffsetClause interface {
	Node
	isGQLLimitAndOffsetClause()
}

func (g *GQLLimitClause) isGQLLimitAndOffsetClause()           {}
func (g *GQLOffsetClause) isGQLLimitAndOffsetClause()          {}
func (g *GQLLimitWithOffsetClause) isGQLLimitAndOffsetClause() {}

// GQLFilterStatement represents `FILTER [WHERE] bool_expression`
//
//	FILTER {{if .Where.Invalid | not}}WHERE{{end}} {{.Expr | sql}}
type GQLFilterStatement struct {
	// pos = Filter
	// end = Expr.end

	Filter token.Pos
	Where  token.Pos
	Expr   Expr
}

// GQLForStatement represents GQL FOR statement.
//
//	FOR {{.ElementName | sql}} IN {{.ArrayExpression | sqlJoin ", "}} {{.WithOffsetClause | sqlOpt}}
type GQLForStatement struct {
	// pos = For
	// end = (WithOffsetClause ?? ArrayExpression).end

	For              token.Pos
	ElementName      *Ident
	ArrayExpression  Expr
	WithOffsetClause *GQLWithOffsetClause
}

// GQLWithOffsetClause represents `WITH OFFSET [AS offset_name]` in FOR statement.
//
//	WITH OFFSET {{if isnil .OffsetName | not}}AS {{.OffsetName | sql}}{{end}}
type GQLWithOffsetClause struct {
	// pos = With
	// end = OffsetName.end || Offset + 6

	With       token.Pos
	Offset     token.Pos
	OffsetName *Ident
}

// GQLLimitClause is wrapper of Limit for GQL
//
//	{{.Limit | sql}}
type GQLLimitClause struct {
	// pos = Limit.pos
	// end = Limit.end

	Limit *Limit
}

// GQLOffsetClause is wrapper of Offset for GQL
//
//	{{.Offset | sql}}
type GQLOffsetClause struct {
	// pos = Offset.pos
	// end = Offset.end
	Offset *Offset
}

// GQLLimitWithOffsetClause is wrapper of Limit and Offset
//
//	{{.Offset | sql}} {{.Offset | sql}}
type GQLLimitWithOffsetClause struct {
	// pos = Limit.pos
	// end = Offset.end

	Limit  *Limit
	Offset *Offset
}

// GQLLimitStatement represents LIMIT statement
//
//	LIMIT {{.Count | sql}}
type GQLLimitStatement struct {
	// pos = Limit
	// end = Count.end

	Limit token.Pos
	Count IntValue
}

// GQLOffsetStatement represents OFFSET statement.
// It also represents SKIP statement as the synonym.
//
// {{if IsSkip}}
type GQLOffsetStatement struct {
	// pos = Offset
	// end = Count.end

	Offset token.Pos
	IsSkip bool
	Count  IntValue
}

// GQLOrderByStatement represents ORDER BY statement.
//
//	ORDER BY {{.OrderBySpecificationList | sqlJoin ", "}}
type GQLOrderByStatement struct {
	// pos = Order
	// end = OrderBySpecificationList[$].end

	Order                    token.Pos
	OrderBySpecificationList []*GQLOrderBySpecification
}

// GQLOrderBySpecification represents a single sort criterion for an expression in ORDER BY.
//
// {{.Expr | sql}} {{.CollationSpecification | sqlOpt}} {{if DirectionPos.Invalid | not}}{{string(Direction)}}{{end}}
type GQLOrderBySpecification struct {
	// pos = Expr.pos
	// end = DirectionPos || CollationSpecification.end

	Expr Expr

	CollationSpecification *GQLCollationSpecification // optional
	DirectionPos           token.Pos                  // optional
	Direction              GQLDirectionEnum
}

// GQLCollationSpecification represents `COLLATE collation_specification`
//
//	COLLATE {{.Specification | sql}}
type GQLCollationSpecification struct {
	// pos = Collate
	// end = Specification.end

	Collate       token.Pos
	Specification StringValue
}

// GQLWithStatement represents WITH statement.
//
//	WITH {{.GQLAllOrDistinctEnum | sql}} {{.ReturnItemList | sqlJoin}} {{.GroupBy | sql}}
type GQLWithStatement struct {
	// pos = With
	// end = (GroupByClause ?? ReturnItemList[$]).end

	With           token.Pos
	AllOrDistinct  GQLAllOrDistinctEnum
	ReturnItemList []GQLReturnItem
	GroupByClause  *GroupBy // optional
}

// GQLReturnItem is similar to SelectItem,
// but it don't permit DotStar and AsAlias without AS.
type GQLReturnItem SelectItem

// GQLReturnStatement represents RETURN statement.
//
//	RETURN {{.AllOrDistinct | sql}} {{.ReturnItemList | sqlJoin}}
//	{{.GroupByClause | sqlOpt}}
//	{{.OrderByClause | sqlOpt}}
//	{{.LimitAntOffsetClause | sqlOpt}}
type GQLReturnStatement struct {
	// pos = Return
	// end = (LimitAndOffsetClause ?? OrderByClause ?? GroupByClause ?? ReturnItemList[$]).end

	Return         token.Pos // position of "RETURN" keyword
	AllOrDistinct  GQLAllOrDistinctEnum
	ReturnItemList []GQLReturnItem

	// Use GoogleSQL GroupBy because it is referenced in docs
	GroupByClause *GroupBy //optional

	// Use GoogleSQL OrderBy because it is referenced in docs
	OrderByClause *OrderBy //optional

	LimitAndOffsetClause GQLLimitAndOffsetClause // optional
}

// GQLLinearGraphVariable represents a single `variable_name = value` entry in LET statement.
//
//	{{.VariableName | sql}} = {{.Value | sql}}
type GQLLinearGraphVariable struct {
	// pos = VariableName.pos
	// end = Value.end

	VariableName *Ident
	Value        Expr
}

// GQLLetStatement represents LET statement.
//
//	LET {{.LinearGraphVariableList | sqlJoin ", "}}
type GQLLetStatement struct {
	// pos = Let
	// end = LinearGraphVariableList[$].end

	Let                     token.Pos
	LinearGraphVariableList []*GQLLinearGraphVariable // len(LinearGraphVariableList) > 0
}

// ================================================================================
//
// GQL graph patterns
//
// ================================================================================

// GQLGraphPattern represents is the toplevel node of GQL graph patterns.
//
//	{{.PathPatternList | sqlJoin}} {{.WhereClause | sqlOpt}}
type GQLGraphPattern struct {
	// pos = PathPatternList[0].pos
	// end = (WhereClause ?? PathPatternList[$]).end

	PathPatternList []*GQLTopLevelPathPattern
	WhereClause     *Where // optional
}

// GQLTopLevelPathPattern is a PathPattern optionally prefixed by PathSearchPrefixOrPathMode.
//
//	{{.PathSearchPrefixOrPathMode | sqlOpt}} {{.PathPattern | sql}}
type GQLTopLevelPathPattern struct {
	// pos = (PathSearchPrefixOrPathMode ?? PathPattern).pos
	// end = PathPattern.end

	PathSearchPrefixOrPathMode GQLPathSearchPrefixOrPathMode // optional
	PathPattern                *GQLPathPattern
}

// GQLPathSearchPrefixOrPathMode represents `{ path_search_prefix | path_mode }`
type GQLPathSearchPrefixOrPathMode interface {
	Node
	isGQLPathSearchPrefixOrPathMode()
}

func (*GQLPathMode) isGQLPathSearchPrefixOrPathMode()         {}
func (*GQLPathSearchPrefix) isGQLPathSearchPrefixOrPathMode() {}

// GQLEdgePattern represents edge pattern nodes.
type GQLEdgePattern interface {
	GQLElementPattern
	isGQLEdgePattern()
}

func (*GQLAbbreviatedEdgeLeft) isGQLEdgePattern()  {}
func (*GQLAbbreviatedEdgeAny) isGQLEdgePattern()   {}
func (*GQLFullEdgeRight) isGQLEdgePattern()        {}
func (*GQLFullEdgeLeft) isGQLEdgePattern()         {}
func (*GQLFullEdgeAny) isGQLEdgePattern()          {}
func (*GQLAbbreviatedEdgeRight) isGQLEdgePattern() {}

// GQLFullEdgeAny is node representing`-[pattern_filler]-` .
//
//	-[{{.PatternFiller | sql}}]-
type GQLFullEdgeAny struct {
	// pos = First
	// end = Last + 1

	First, Last   token.Pos
	PatternFiller *GQLPatternFiller
}

// GQLFullEdgeLeft represents `<-[pattern_filler]-`
//
//	<-[{{.PatternFiller | sql}}]-
type GQLFullEdgeLeft struct {
	// pos = First
	// end = Last + 1

	First         token.Pos // position of "<"
	Last          token.Pos // position of the last "-"
	PatternFiller *GQLPatternFiller
}

// GQLFullEdgeRight represents “-[pattern_filler]->
//
//	-[{{.PatternFiller | sql}}]->
type GQLFullEdgeRight struct {
	// pos = First
	// end = Last + 1

	First         token.Pos // position of the first "-"
	Last          token.Pos // position of ">"
	PatternFiller *GQLPatternFiller
}

// GQLAbbreviatedEdgeAny represents `-`.
//
//	-
type GQLAbbreviatedEdgeAny struct {
	// pos = Hyphen
	// end = Hyphen +1

	Hyphen token.Pos // position of "-"
}

// GQLAbbreviatedEdgeLeft represents `<-`.
//
//	<-
type GQLAbbreviatedEdgeLeft struct {
	// pos = First
	// end = Last + 1

	First token.Pos // position of "<"
	Last  token.Pos // position of "-"
}

// GQLAbbreviatedEdgeRight represents `->`.
//
//	->
type GQLAbbreviatedEdgeRight struct {
	// pos = First
	// end = Last + 1

	First token.Pos // position of "-"
	Last  token.Pos // position of ">"
}

// GQLQuantifiablePathTerm represents GQLPathTerm with optional Hint and optional GQLQuantifier..
// NOTE: This node is not documented in spec, but inferred by [quantified_path_primary] and [graph traversal hints].
//
// [graph traversal hints]: https://cloud.google.com/spanner/docs/reference/standard-sql/graph-query-statements#graph_hints
// [quantified_path_primary] https://cloud.google.com/spanner/docs/reference/standard-sql/graph-patterns#quantified_paths
//
//	{{.Hint | sqlOpt}}{{.PathTerm | sql}}{{.Quantifier | sqlOpt}}
type GQLQuantifiablePathTerm struct {
	// pos = (Hint ?? PathTerm).pos
	// end = (Quantifier ?? PathTerm).end

	Hint       *Hint // optional
	PathTerm   GQLPathTerm
	Quantifier GQLQuantifier // optional
}

// GQLPathPattern represents a path pattern that matches paths in a property graph.
//
//	{{.PathTermList | sqlJoin ""}}
type GQLPathPattern struct {
	// pos = PathTermList[0].pos
	// end = PathTermList[$].end

	PathTermList []*GQLQuantifiablePathTerm
}

// GQLPathTerm represents ` { element_pattern | subpath_pattern }`
type GQLPathTerm interface {
	Node
	isGQLPathTerm()
}

func (*GQLSubpathPattern) isGQLPathTerm()       {}
func (*GQLNodePattern) isGQLPathTerm()          {}
func (*GQLAbbreviatedEdgeRight) isGQLPathTerm() {}
func (*GQLAbbreviatedEdgeLeft) isGQLPathTerm()  {}
func (*GQLAbbreviatedEdgeAny) isGQLPathTerm()   {}
func (*GQLFullEdgeRight) isGQLPathTerm()        {}
func (*GQLFullEdgeLeft) isGQLPathTerm()         {}
func (*GQLFullEdgeAny) isGQLPathTerm()          {}

// GQLWhereClause represents `WHERE bool_expression` clause.
//
//	WHERE {{.BoolExpression | sql}}
type GQLWhereClause struct {
	// pos = Where
	// end = BoolExpression.end

	Where          token.Pos
	BoolExpression Expr
}

// GQLElementPattern represents a node pattern or an edge pattern.
type GQLElementPattern interface {
	Node
	GQLPathTerm
	isGQLElementPattern()
}

func (*GQLFullEdgeAny) isGQLElementPattern()          {}
func (*GQLFullEdgeLeft) isGQLElementPattern()         {}
func (*GQLFullEdgeRight) isGQLElementPattern()        {}
func (*GQLAbbreviatedEdgeAny) isGQLElementPattern()   {}
func (*GQLAbbreviatedEdgeLeft) isGQLElementPattern()  {}
func (*GQLAbbreviatedEdgeRight) isGQLElementPattern() {}

// GQLPathMode represents to include or exclude paths that have repeating edges based on the specified mode.
//
//	{{.ModeToken | sql}} {{.PathOrPathsToken | sqlOpt}}
type GQLPathMode struct {
	// pos = ModeToken.pos
	// end = (PathOrPathsToken ?? ModeToken).end

	Mode             GQLPathModeEnum
	ModeToken        *Ident
	PathOrPathsToken *Ident // optional
}

// GQLQuantifier represents `{ fixed_quantifier | bounded_quantifier }`.
type GQLQuantifier interface {
	Node
	isGQLQuantifier()
}

func (g *GQLFixedQuantifier) isGQLQuantifier()   {}
func (g *GQLBoundedQuantifier) isGQLQuantifier() {}

// GQLFixedQuantifier represents the exact number of times the path pattern portion must repeat.
//
//	{{"{"}}{{.Bound | sql}}{{"}"}}
type GQLFixedQuantifier struct {
	// pos = LBrace
	// end = RBrace + 1

	LBrace, RBrace token.Pos
	Bound          IntValue
}

// GQLBoundedQuantifier represents the minimum and maximum number of times the path pattern portion can repeat.
//
//	{{"{"}}{{.LowerBound | sqlOpt}}, {{.UpperBound | sql}}{{"}"}}
type GQLBoundedQuantifier struct {
	// pos = LBrace
	// end = RBrace + 1

	LBrace, RBrace token.Pos
	LowerBound     IntValue // optional
	UpperBound     IntValue
}

// GQLSubpathPattern represents a path pattern enclosed in parentheses.
//
//	({{.PathMode | sqlOpt}} {{.PathPattern | sql}} {{.WhereClause | sqlOpt}})
type GQLSubpathPattern struct {
	// pos = LParen
	// end = RParen + 1

	LParen, RParen token.Pos    // position of "(" and ")"
	PathMode       *GQLPathMode // optional
	PathPattern    *GQLPathPattern
	WhereClause    *Where // optional
}

// GQLNodePattern represents a pattern to match nodes in a property graph.
//
//	({{.PatternFiller | sql}})
type GQLNodePattern struct {
	// pos = LParen
	// end = RParen + 1

	LParen, RParen token.Pos
	PatternFiller  *GQLPatternFiller
}

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

// GQLPatternFiller represents specifications on the node or edge pattern that you want to match.
//
//	{{.Hint | sqlOpt}}
//	{{.GraphPatternVariable | sqlOpt}}
//	{{.IsLabelCondition | sqlOpt}}
//	{{.Filter | sqlOpt}}
type GQLPatternFiller struct {
	// pos = (Hint ?? GraphPatternVariable ?? IsLabelCondition ?? Filter).pos
	// end = (Filter ?? IsLabelCondition ?? GraphPatternVariable ?? Hint).end

	// Hint is graph element hint which is a table hint.
	Hint                 *Hint                  // optional
	GraphPatternVariable *Ident                 // optional
	IsLabelCondition     *GQLIsLabelCondition   // optional
	Filter               GQLPatternFillerFilter // optional
}

// GQLIsLabelCondition represents `{ IS | : } label_expression`.
// It normalizes `IS` to `:`.
//
//	: {{.LabelExpression | sql}}
type GQLIsLabelCondition struct {
	// pos = IsOrColon
	// end = LabelExpression.end

	IsOrColon       token.Pos
	LabelExpression GQLLabelExpression
}

// GQLLabelExpression represents the expression for the label.
// It is formed by combining one or more labels with logical operators (AND, OR, NOT) and parentheses for grouping.
// See https://cloud.google.com/spanner/docs/reference/standard-sql/graph-patterns#label_expression_definition.
type GQLLabelExpression interface {
	Node
	isGQLLabelExpression()
}

// Note: Spanner Graph documentation don't say about paren expression, but there is.
func (g *GQLLabelParenExpression) isGQLLabelExpression() {}
func (g *GQLLabelOrExpression) isGQLLabelExpression()    {}
func (g *GQLLabelAndExpression) isGQLLabelExpression()   {}
func (g *GQLLabelNotExpression) isGQLLabelExpression()   {}
func (g *GQLLabelName) isGQLLabelExpression()            {}

// GQLLabelAndExpression represents `label_expression|label_expression`.
//
//	{{.Left | sql}}|{{.Right | sql}}
type GQLLabelOrExpression struct {
	// pos = Left.pos
	// end = Right.end

	Left, Right GQLLabelExpression
}

// GQLLabelParenExpression represents `(label_expression)`.
//
//	({{.LabelExpr | sql}})
type GQLLabelParenExpression struct {
	// pos = LParen
	// end = RParen + 1

	LParen, RParen token.Pos
	LabelExpr      GQLLabelExpression
}

// GQLLabelAndExpression represents `label_expression&label_expression`.
//
//	{{.Left | sql}}&{{.Right | sql}}
type GQLLabelAndExpression struct {
	// pos = Left.pos
	// end = Right.end

	Left, Right GQLLabelExpression
}

// GQLLabelNotExpression represents `!label_expression`.
//
//	!{{.LabelExpression | sql}}
type GQLLabelNotExpression struct {
	// pos = Not
	// end = LabelExpression.end

	Not             token.Pos // position of "!"
	LabelExpression GQLLabelExpression
}

// GQLLabelName represents the label to match.
//
//	{{if .IsPercent}}%{{else}}{{.LabelName | sql}}{{end}}
type GQLLabelName struct {
	// pos = StartPos
	// end = LabelName.end || StartPos + 1

	StartPos  token.Pos // position of "%" or LabelName
	IsPercent bool
	LabelName *Ident // optional
}

// GQLPatternFillerFilter represents `{where_clause | property_filters}` in GQLPatternFiller.
type GQLPatternFillerFilter interface {
	Node
	isGQLPatternFillerFilter()
}

func (g *GQLPropertyFilters) isGQLPatternFillerFilter() {}
func (g *GQLWhereClause) isGQLPatternFillerFilter()     {}
func (w *Where) isGQLPatternFillerFilter()              {}

// GQLPropertyFilters represents `{ element_property[, ...] }` in GQLPatternFiller.
//
//	{{"{"}}{{.PropertyFilterElemList | sqlJoin ", "}}{{"}"}}
type GQLPropertyFilters struct {
	// pos = LBrace
	// end = RBrace + 1

	LBrace                 token.Pos // position of "{"
	PropertyFilterElemList []*GQLElementProperty
	RBrace                 token.Pos // position of "}"
}

// GQLElementProperty represents an element of GQLPropertyFilters.
//
//	{{.ElementPropertyName | sql}}: {{.ElementPropertyValue | sql}}
type GQLElementProperty struct {
	// pos = ElementPropertyName.pos
	// end = ElementPropertyValue.pos

	ElementPropertyName  *Ident
	ElementPropertyValue Expr
}

// GQLPathSearchPrefix represents `{"ALL" | "ANY" | "ANY SHORTEST"}`.
//
//	{{string(.SearchPrefix)}}
type GQLPathSearchPrefix struct {
	// pos = StartPos
	// end = LastEnd

	StartPos     token.Pos
	LastEnd      token.Pos // end of last token
	SearchPrefix GQLSearchPrefixEnum
}
