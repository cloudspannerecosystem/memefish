package ast

import (
	"github.com/MakeNowJust/memefish/pkg/token"
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

func (QueryStatement) isStatement() {}
func (CreateDatabase) isStatement() {}
func (CreateTable) isStatement()    {}
func (CreateIndex) isStatement()    {}
func (AlterTable) isStatement()     {}
func (DropTable) isStatement()      {}
func (DropIndex) isStatement()      {}
func (Insert) isStatement()         {}
func (Delete) isStatement()         {}
func (Update) isStatement()         {}

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

// TableExpr represents JOIN operands.
type TableExpr interface {
	Node
	isTableExpr()
}

func (Unnest) isTableExpr()            {}
func (TableName) isTableExpr()         {}
func (SubQueryTableExpr) isTableExpr() {}
func (ParenTableExpr) isTableExpr()    {}
func (Join) isTableExpr()              {}

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

func (BinaryExpr) isExpr()       {}
func (UnaryExpr) isExpr()        {}
func (InExpr) isExpr()           {}
func (IsNullExpr) isExpr()       {}
func (IsBoolExpr) isExpr()       {}
func (BetweenExpr) isExpr()      {}
func (SelectorExpr) isExpr()     {}
func (IndexExpr) isExpr()        {}
func (CallExpr) isExpr()         {}
func (CountStarExpr) isExpr()    {}
func (CastExpr) isExpr()         {}
func (ExtractExpr) isExpr()      {}
func (CaseExpr) isExpr()         {}
func (ParenExpr) isExpr()        {}
func (ScalarSubQuery) isExpr()   {}
func (ArraySubQuery) isExpr()    {}
func (ExistsSubQuery) isExpr()   {}
func (Param) isExpr()            {}
func (Ident) isExpr()            {}
func (Path) isExpr()             {}
func (ArrayLiteral) isExpr()     {}
func (StructLiteral) isExpr()    {}
func (NullLiteral) isExpr()      {}
func (BoolLiteral) isExpr()      {}
func (IntLiteral) isExpr()       {}
func (FloatLiteral) isExpr()     {}
func (StringLiteral) isExpr()    {}
func (BytesLiteral) isExpr()     {}
func (DateLiteral) isExpr()      {}
func (TimestampLiteral) isExpr() {}
func (NumericLiteral) isExpr()   {}

// InCondition is right-side value of IN operator.
type InCondition interface {
	Node
	isInCondition()
}

func (UnnestInCondition) isInCondition()   {}
func (SubQueryInCondition) isInCondition() {}
func (ValuesInCondition) isInCondition()   {}

// Type represents type node.
type Type interface {
	Node
	isType()
}

func (SimpleType) isType() {}
func (ArrayType) isType()  {}
func (StructType) isType() {}

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

func (CreateDatabase) isDDL() {}
func (CreateTable) isDDL()    {}
func (AlterTable) isDDL()     {}
func (DropTable) isDDL()      {}
func (CreateIndex) isDDL()    {}
func (DropIndex) isDDL()      {}

// TableAlternation represents ALTER TABLE action.
type TableAlternation interface {
	Node
	isTableAlternation()
}

func (AddColumn) isTableAlternation()      {}
func (AddForeignKey) isTableAlternation()  {}
func (DropColumn) isTableAlternation()     {}
func (DropConstraint) isTableAlternation() {}
func (SetOnDelete) isTableAlternation()    {}
func (AlterColumn) isTableAlternation()    {}
func (AlterColumnSet) isTableAlternation() {}

// SchemaType represents types for schema.
type SchemaType interface {
	Node
	isSchemaType()
}

func (ScalarSchemaType) isSchemaType() {}
func (SizedSchemaType) isSchemaType()  {}
func (ArraySchemaType) isSchemaType()  {}

// DML represents data manipulation language in SQL.
//
//https://cloud.google.com/spanner/docs/data-definition-language
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

// ================================================================================
//
// SELECT
//
// ================================================================================

// QueryStatement is query statement node.
//
//     {{if .Hint}}{{.Hint | sql}}{{end}}
//     {{.Expr | sql}}
//
// https://cloud.google.com/spanner/docs/query-syntax
type QueryStatement struct {
	// pos = (Hint ?? Expr).pos
	// end = Expr.end

	Hint  *Hint // optional
	Query QueryExpr
}

// Hint is hint node.
//
//     @{{"{"}}{{.Records | sqlJoin ","}}{{"}"}}
type Hint struct {
	// pos = Atmark
	// end = Rbrace + 1

	Atmark token.Pos // position of "@"
	Rbrace token.Pos // position of "}"

	Records []*HintRecord // len(Records) > 0
}

// HintRecord is hint record node.
//
//     {{.Key | sql}}={{.Value | sql}}
type HintRecord struct {
	// pos = Key.pos
	// end = Value.end

	Key   *Ident
	Value Expr
}

// Select is SELECT statement node.
//
//     SELECT
//       {{if .Distinct}}DISTINCT{{end}}
//       {{if .AsStruct}}AS STRUCT{{end}}
//       {{.Results | sqlJoin ","}}
//       {{.From | sqlOpt}}
//       {{.Where | sqlOpt}}
//       {{.GroupBy | sqlOpt}}
//       {{.Having | sqlOpt}}
//       {{.OrderBy | sqlOpt}}
//       {{.Limit | sqlOpt}}
type Select struct {
	// pos = Select
	// end = (Limit ?? OrderBy ?? Having ?? GroupBy ?? Where ?? From ?? Results[$]).end

	Select token.Pos // position of "select" keyword

	Distinct bool
	AsStruct bool
	Results  []SelectItem // len(Results) > 0
	From     *From        // optional
	Where    *Where       // optional
	GroupBy  *GroupBy     // optional
	Having   *Having      // optional
	OrderBy  *OrderBy     // optional
	Limit    *Limit       // optional
}

// CompoundQuery is query statement node compounded by set operators.
//
//     {{.Queries | sqlJoin (printf "%s %s" .Op or(and(.Distinct, "DISTINCT"), "ALL"))}}
//       {{.OrderBy | sqlOpt}}
//       {{.Limit | sqlOpt}}
type CompoundQuery struct {
	// pos = Queries[0].pos
	// end = (Limit ?? OrderBy ?? Queries[$]).end

	Op       SetOp
	Distinct bool
	Queries  []QueryExpr // len(List) >= 2
	OrderBy  *OrderBy    // optional
	Limit    *Limit      // optional
}

// SubQuery is subquery statement node.
//
//     ({{.Expr | sql}} {{.OrderBy | sqlOpt}} {{.Limit | sqlOpt}})
type SubQuery struct {
	// pos = Lparen
	// end = (Query ?? Limit).end || Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query   QueryExpr
	OrderBy *OrderBy // optional
	Limit   *Limit   // optional
}

// Star is a single * in SELECT result columns list.
//
//     *
type Star struct {
	// pos = Star
	// end = Star + 1

	Star token.Pos // position of "*"
}

// DotStar is expression with * in SELECT result columns list.
//
//     {{.Expr | sql}}.*
type DotStar struct {
	// pos = Expr.pos
	// end = Star + 1

	Star token.Pos // position of "*"

	Expr Expr
}

// Alias is aliased expression by AS clause in SELECT result columns list.
//
//    {{.Expr | sql}} {{.As | sql}}
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
// NOTE: Sometime keyword AS can be omited.
//       In this case, it.token.Pos() == it.Alias.token.Pos(), so we can detect this.
//
//     AS {{.Alias | sql}}
type AsAlias struct {
	// pos = As || Alias.pos
	// end = Alias.End

	As token.Pos // position of "AS" keyword

	Alias *Ident
}

// ExprSelectItem is Expr wrapper for SelectItem.
//
//     {{.Expr | sql}}
type ExprSelectItem struct {
	// pos = Expr.pos
	// end = Expr.end

	Expr Expr
}

// From is FROM clause node.
//
//     FROM {{.Source | sql}}
type From struct {
	// pos = From
	// end = Source.end

	From token.Pos // position of "FROM" keyword

	Source TableExpr
}

// Where is WHERE clause node.
//
//    WHERE {{.Expr | sql}}
type Where struct {
	// pos = Where
	// end = Expr.end

	Where token.Pos // position of "WHERE" keyword

	Expr Expr
}

// GroupBy is GROUP BY clause node.
//
//    GROUP BY {{.Exprs | sqlJoin ","}}
type GroupBy struct {
	// pos = Group
	// end = Exprs[$].end

	Group token.Pos // position of "GROUP" keyword

	Exprs []Expr // len(Exprs) > 0
}

// Having is HAVING clause node.
//
//     HAVING {{.Expr | sql}}
type Having struct {
	// pos = Having
	// end = Expr.end

	Having token.Pos // position of "HAVING" keyword

	Expr Expr
}

// OrderBy is ORDER BY clause node.
//
//     ORDER BY {{.Items | sqlJoin ","}}
type OrderBy struct {
	// pos = Order
	// end = Items[$].end

	Order token.Pos // position of "ORDER" keyword

	Items []*OrderByItem // len(Items) > 0
}

// OrderByItem is expression node in ORDER BY clause list.
//
//     {{.Expr | sql}} {{.Collate | sqlOpt}} {{.Direction}}
type OrderByItem struct {
	// pos = Expr.pos
	// end = DirPos + len(Dir) || (Collate ?? Expr).pos

	DirPos token.Pos // position of Dir

	Expr    Expr
	Collate *Collate  // optional
	Dir     Direction // optional
}

// Collate is COLLATE clause node in ORDER BY item.
//
//     COLLATE {{.Value | sql}}
type Collate struct {
	// pos = Collate
	// end = Value.end

	Collate token.Pos // position of "COLLATE" keyword

	Value StringValue
}

// Limit is LIMIT clause node.
//
//     LIMIT {{.Count | sql}} {{.Offset | sqlOpt}}
type Limit struct {
	// pos = Limit
	// end = (Offset ?? Count).end

	Limit token.Pos // position of "LIMIT" keyword

	Count  IntValue
	Offset *Offset // optional
}

// Offset is OFFSET clause node in LIMIT clause.
//
//     OFFSET {{.Value | sql}}
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
//     {{if .Implicit}}{{.Expr | sql}}{{else}}UNNEST({{.Expr | sql}}){{end}}
//       {{.Hint | sqlOpt}}
//       {{.As | sqlOpt}}
//       {{.WithOffset | sqlOpt}}
//       {{.Sample | sqlOpt}}
type Unnest struct {
	// pos = Unnest || Expr.pos
	// end = (Sample ?? WithOffset ?? As ?? Hint).end || Rparen + 1 || Expr.end

	Unnest token.Pos // position of "UNNEST"
	Rparen token.Pos // position of ")"

	Implicit   bool
	Expr       Expr         // Path or Ident when Implicit is true
	Hint       *Hint        // optional
	As         *AsAlias     // optional
	WithOffset *WithOffset  // optional
	Sample     *TableSample // optional
}

// WithOffset is WITH OFFSET clause node after UNNEST call.
//
//     WITH OFFSET {{.As | sqlOpt}}
type WithOffset struct {
	// pos = With
	// end = As.end || Offset + 6

	With, Offset token.Pos // position of "WITH" and "OFFSET" keywords

	As *AsAlias // optional
}

// TableName is table name node in FROM clause.
//
//     {{.Table | sql}} {{.Hint | sqlOpt}} {{.As | sqlOpt}} {{.Sample | sqlOpt}}
type TableName struct {
	// pos = Table.pos
	// end = (Sample ?? As ?? Hint ?? Table).end

	Table  *Ident
	Hint   *Hint        // optional
	As     *AsAlias     // optional
	Sample *TableSample // optional
}

// SubQueryTableExpr is subquery inside JOIN expression.
//
//     ({{.Query | sql}}) {{.As | sqlOpt}} {{.Sample | sqlOpt}}
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
//     ({{.Source | sql}}) {{.Sample | sqlOpt}}
type ParenTableExpr struct {
	// pos = Lparen
	// end = Sample.end || Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Source TableExpr    // SubQueryJoinExpr (without As) or Join
	Sample *TableSample // optional
}

// Join is JOIN expression.
//
//     {{.Left | sql}}
//       {{.Op}} {{.Method}} {{.Hint | sqlOpt}}
//       {{.Right | sql}}
//     {{.Cond | sqlOpt}}
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
//     ON {{.Expr | sql}}
type On struct {
	// pos = On
	// end = Expr.end

	On token.Pos // position of "ON" keyword

	Expr Expr
}

// Using is Using condition of JOIN expression.
//
//     USING ({{Idents | sqlJoin ","}})
type Using struct {
	// pos = Using
	// end = Rparen + 1

	Using  token.Pos // position of "USING" keyword
	Rparen token.Pos // position of ")"

	Idents []*Ident // len(Idents) > 0
}

// TableSample is TABLESAMPLE clause node.
//
//     TABLESAMPLE {{.Method}} {{.Size | sql}}
type TableSample struct {
	// pos = TableSample
	// end = Size.end

	TableSample token.Pos // position of "TABLESAMPLE" keyword

	Method TableSampleMethod
	Size   *TableSampleSize
}

// TableSampleSize is size part of TABLESAMPLE clause.
//
//     ({{.Value | sql}} {{.Unit}})
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
//     {{.Left | sql}} {{.Op}} {{.Right | sql}}
type BinaryExpr struct {
	// pos = Left.pos
	// end = Right.pos

	Op BinaryOp

	Left, Right Expr
}

// UnaryExpr is unary operator expression node.
//
//     {{.Op}} {{.Expr | sql}}
type UnaryExpr struct {
	// pos = OpPos
	// end = Expr.end

	OpPos token.Pos // position of Op

	Op   UnaryOp
	Expr Expr
}

// InExpr is IN expression node.
//
//     {{.Left | sql}} {{if .Not}}NOT{{end}} IN {{.Right | sql}}
type InExpr struct {
	// pos = Left.pos
	// end = Right.end

	Not   bool
	Left  Expr
	Right InCondition
}

// UnnestInCondition is UNNEST call at IN condition.
//
//     UNNEST({{.Expr | sql}})
type UnnestInCondition struct {
	// pos = Unnest
	// end = Rparen + 1

	Unnest token.Pos
	Rparen token.Pos

	Expr Expr
}

// SubQueryInCondition is subquery at IN condition.
//
//     ({{.Query | sql}})
type SubQueryInCondition struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query QueryExpr
}

// ValuesInCondition is parenthesized values at IN condition.
//
//     ({{.Exprs | sqlJoin ","}})
type ValuesInCondition struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Exprs []Expr // len(Exprs) > 0
}

// IsNullExpr is IS NULL expression node.
//
//     {{.Left | sql}} IS {{if .Not}}NOT{{end}} NULL
type IsNullExpr struct {
	// pos = Expr.pos
	// end = Null + 4

	Null token.Pos // position of "NULL"

	Not  bool
	Left Expr
}

// IsBoolExpr is IS TRUE/FALSE expression node.
//
//     {{.Left | sql}} IS {{if .Not}}NOT{{end}} {{if .Right}}TRUE{{else}}FALSE{{end}}
type IsBoolExpr struct {
	// pos = Expr.pos
	// end = RightPos + (Right ? 4 : 5)

	RightPos token.Pos // position of Right

	Not   bool
	Left  Expr
	Right bool
}

// BetweenExpr is BETWEEN expression node.
//
//     {{.Left | sql}} {{if .Not}}NOT{{end}} BETWEEN {{.RightStart | sql}} AND {{.RightEnd | sql}}
type BetweenExpr struct {
	// pos = Left.pos
	// end = RightEnd.end

	Not                        bool
	Left, RightStart, RightEnd Expr
}

// SelectorExpr is struct field access expression node.
//
//     {{.Expr | sql}}.{{.Ident | sql}}
type SelectorExpr struct {
	// pos = Expr.pos
	// end = Ident.pos

	Expr  Expr
	Ident *Ident
}

// IndexExpr is array item access expression node.
//
//     {{.Expr | sql}}[{{if .Ordinal}}ORDINAL{{else}}OFFSET{{end}}({{.Index | sql}})]
type IndexExpr struct {
	// pos = Expr.pos
	// end = Rbrack + 1

	Rbrack token.Pos // position of "]"

	Ordinal     bool
	Expr, Index Expr
}

// CallExpr is function call expression node.
//
//     {{.Func | sql}}({{if .Distinct}}DISTINCT{{end}} {{.Args | sql}})
type CallExpr struct {
	// pos = Func.pos
	// end = Rparen + 1

	Rparen token.Pos // position of ")"

	Func     *Ident
	Distinct bool
	Args     []*Arg
}

// Arg is function call argument.
//
//     {{if .IntervalUnit}}
//       INTERVAL {{.Expr | sql}} {{.IntervalUnit | sql}}
//     {{else}}
//       {{.Expr | sql}}
//     {{end}}
type Arg struct {
	// pos = Interval || Expr.pos
	// end = (IntervalUnit ?? Expr).end

	Interval token.Pos // position of "INTERVAL" keyword

	Expr         Expr
	IntervalUnit *Ident // optional
}

// CountStarExpr is node just for COUNT(*).
//
//     COUNT(*)
type CountStarExpr struct {
	// pos = Count
	// end = Rparen + 1

	Count  token.Pos // position of "COUNT"
	Rparen token.Pos // position of ")"
}

// ExtractExpr is EXTRACT call expression node.
//
//     EXTRACT({{.Part | sql}} FROM {{.Expr | sql}} {{.AtTimeZone | sqlOpt}})
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
//     AT TIME ZONE {{.Expr | sql}}
type AtTimeZone struct {
	// pos = At
	// end = Expr.end

	At token.Pos // position of "AT" keyword

	Expr Expr
}

// CastExpr is CAST call expression node.
//
//     CAST({{.Expr | sql}} AS {{.Type | sql}})
type CastExpr struct {
	// pos = Cast
	// end = Rparen + 1

	Cast   token.Pos // position of "CAST" keyword
	Rparen token.Pos // position of ")"

	Expr Expr
	Type Type
}

// CaseExpr is CASE expression node.
//
//     CASE {{.Expr | sqlOpt}}
//       {{.Whens | sqlJoin "\n"}}
//       {{.Else | sqlOpt}}
//     END
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
//     WHEN {{.Cond | sql}} THEN {{.Then | sql}}
type CaseWhen struct {
	// pos = Case
	// end = Then.end

	When token.Pos // position of "WHEN" keyword

	Cond, Then Expr
}

// CaseElse is ELSE clause in CASE expression.
//
//     ELSE {{.Expr | sql}}
type CaseElse struct {
	// pos = Else
	// end = Expr.end

	Else token.Pos // position of "ELSE" keyword

	Expr Expr
}

// ParenExpr is parenthesized expression node.
//
//     ({{. | sql}})
type ParenExpr struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Expr Expr
}

// ScalarSubQuery is subquery in expression.
//
//     ({{.Query | sql}})
type ScalarSubQuery struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Query QueryExpr
}

// ArraySubQuery is subquery in ARRAY call.
//
//     ARRAY({{.Query | sql}})
type ArraySubQuery struct {
	// pos = Array
	// end = Rparen + 1

	Array  token.Pos // position of "ARRAY" keyword
	Rparen token.Pos // position of ")"

	Query QueryExpr
}

// ExistsSubQuery is subquery in EXISTS call.
//
//     EXISTS {{.Hint | sqlOpt}} ({{.Query | sql}})
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
//     @{{.Name}}
type Param struct {
	// pos = Atmark
	// end = pos + 1 + len(Name)

	Atmark token.Pos

	Name string
}

// Ident is identifier node.
//
//     {{.Name | sqlIdentQuote}}
type Ident struct {
	// pos = IdentPos
	// end = IdentEnd

	NamePos, NameEnd token.Pos // position of this name

	Name string
}

// Path is dot-chained identifier list.
//
//     {{.Idents | sqlJoin "."}}
type Path struct {
	// pos = Idents[0].pos
	// end = idents[$].end

	Idents []*Ident // len(Idents) >= 2
}

// AraryLiteral is array literal node.
//
//     ARRAY{{if .Type}}<{{.Type | sql}}>{{end}}[{{.Values | sqlJoin ","}}]
type ArrayLiteral struct {
	// pos = Array || Lbrack
	// end = Rbrack + 1

	Array          token.Pos // position of "ARRAY" keyword
	Lbrack, Rbrack token.Pos // position of "[" and "]"

	Type   Type // optional
	Values []Expr
}

// StructLiteral is struct literal node.
//
//     STRUCT{{if .Type}}<{{.Fields | sqlJoin ","}}>{{end}}({{.Values | sqlJoin ","}})
type StructLiteral struct {
	// pos = Struct || Lparen
	// end = Rparen

	Struct         token.Pos // position of "STRUCT"
	Lparen, Rparen token.Pos // position of "(" and ")"

	// NOTE: Distinguish nil from len(Fields) == 0 case.
	//       nil means type is not specified, or empty slice means this struct has 0 fields.
	Fields []*StructField
	Values []Expr
}

// NullLiteral is just NULL literal.
//
//     NULL
type NullLiteral struct {
	// pos = Null
	// end = Null + 4

	Null token.Pos // position of "NULL"
}

// BoolLiteral is boolean literal node.
//
//     {{if .Value}}TRUE{{else}}FALSE{{end}}
type BoolLiteral struct {
	// pos = ValuePos
	// end = ValuePos + (Value ? 4 : 5)

	ValuePos token.Pos // position of this value

	Value bool
}

// IntLiteral is integer literal node.
//
//     {{.Value}}
type IntLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Base  int // 10 or 16
	Value string
}

// FloatLiteral is floating point number literal node.
//
//     {{.Value}}
type FloatLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value string
}

// StringLiteral is string literal node.
//
//     {{.Value | sqlStringQuote}}
type StringLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value string
}

// BytesLiteral is bytes literal node.
//
//     B{{.Value | sqlByesQuote}}
type BytesLiteral struct {
	// pos = ValuePos
	// end = ValueEnd

	ValuePos, ValueEnd token.Pos // position of this value

	Value []byte
}

// DateLiteral is date literal node.
//
//     DATE {{.Value | sql}}
type DateLiteral struct {
	// pos = Date
	// end = Value.end

	Date token.Pos // position of "DATE"

	Value *StringLiteral
}

// TimestampLiteral is timestamp literal node.
//
//     TIMESTAMP {{.Value | sql}}
type TimestampLiteral struct {
	// pos = Timestamp
	// end = ValueEnd.end

	Timestamp token.Pos // position of "TIMESTAMP"

	Value *StringLiteral
}

// NumericLiteral is numeric literal node.
//
//     NUMERIC {{.Value | sql}}
type NumericLiteral struct {
	// pos = Numeric
	// end = ValueEnd.end

	Numeric token.Pos // position of "NUMERIC"

	Value *StringLiteral
}

// ================================================================================
//
// Type
//
// ================================================================================

// SimpleType is type node having no parameter like INT64, STRING.
//
//    {{.Name}}
type SimpleType struct {
	// pos = NamePos
	// end = NamePos + len(Name)

	NamePos token.Pos // position of this name

	Name ScalarTypeName
}

// ArrayType is array type node.
//
//     ARRAY<{{.Item | sql}}>
type ArrayType struct {
	// pos = Array
	// end = Gt + 1

	Array token.Pos // position of "ARRAY" keyword
	Gt    token.Pos // position of ">"

	Item Type
}

// StructType is struct type node.
//
//     STRUCT<{{.Fields | sqlJoin ","}}>
type StructType struct {
	// pos = Struct
	// end = Gt + 1

	Struct token.Pos // position of "STRUCT"
	Gt     token.Pos // position of ">"

	Fields []*StructField
}

// StructField is field in struct type node.
//
//     {{if .Ident}}{{.Ident | sql}}{{end}} {{.Type | sql}}
type StructField struct {
	// pos = (Ident ?? Type).pos
	// end = Type.end

	Ident *Ident
	Type  Type
}

// ================================================================================
//
// Cast for Special Cases
//
// ================================================================================

// CastIntValue is cast call in integer value context.
//
//     CAST({{.Expr | sql}} AS INT64)
type CastIntValue struct {
	// pos = Cast
	// end = Rparen + 1

	Cast   token.Pos // position of "CAST" keyword
	Rparen token.Pos // position of ")"

	Expr IntValue // IntLit or Param
}

// CasrNumValue is cast call in number value context.
//
//     CAST({{.Expr | sql}} AS {{.Type}})
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

// CreateDatabase is CREATE DATABASE statement node.
//
//     CREATE DATABASE {{.Name | sql}}
type CreateDatabase struct {
	// pos = Create
	// end = Name.end

	Create token.Pos // position of "CREATE" keyword

	Name *Ident
}

// CreateTable is CREATE TABLE statement node.
//
//     CREATE TABLE {{.Name | sql}} (
//       {{.Columns | sqlJoin ","}}
//     )
//     PRIMARY KEY ({{.PrimaryKeys | sqlJoin ","}})
//     {{.Cluster | sqlOpt}}
type CreateTable struct {
	// pos = Create
	// end = Cluster.end || Rparen + 1

	Create token.Pos // position of "CREATE" keyword
	Rparen token.Pos // position of ")" of PRIMARY KEY clause

	Name        *Ident
	Columns     []*ColumnDef
	PrimaryKeys []*IndexKey
	ForeignKeys []*ForeignKey
	Cluster     *Cluster // optional
}

// ColumnDef is column definition in CREATE TABLE.
//
//     {{.Name | sql}}
//     {{.Type | sql}} {{if .NotNull}}NOT NULL{{end}}
//     {{.GeneratedExpr | sqlOpt}}
//     {{.Options | sqlOpt}}
type ColumnDef struct {
	// pos = Name.pos
	// end = Options.end || GeneratedExpr.end || Null + 4 || Type.end

	Null token.Pos // position of "NULL"

	Name          *Ident
	Type          SchemaType
	NotNull       bool
	GeneratedExpr *GeneratedColumnExpr // optional
	Options       *ColumnDefOptions    // optional
}

// GeneratedColumnExpr is generated column expression.
//
//     AS ({{.Expr | sql}}) STORED
type GeneratedColumnExpr struct {
	// pos = As
	// end = Stored

	As     token.Pos // position of "AS" keyword
	Stored token.Pos // position of "STORED" keyword

	Expr Expr
}

// ColumnDefOption is options for column definition.
//
//     OPTIONS(allow_commit_timestamp = {{if .AllowCommitTimestamp}}true{{else}null{{end}}})
type ColumnDefOptions struct {
	// pos = Options
	// end = Rparen + 1

	Options token.Pos // position of "OPTIONS" keyword
	Rparen  token.Pos // position of ")"

	AllowCommitTimestamp bool
}

// ForeignKey is foreign key specifier in CREATE TABLE and ALTER TABLE.
//
//     {{if .Name}}CONSTRAINT {{.Name}}{{end}} FOREIGN KEY ({{.ColumnNames | sqlJoin ","}}) REFERENCES {{.ReferenceTable}}({{.ReferenceColumns | sqlJoin ","}})
type ForeignKey struct {
	// pos = Constraint || Foreign
	// end = Rparen + 1

	Constraint token.Pos // position of "CONSTRAINT" keyword when Implicit is true
	Foreign    token.Pos // position of "FOREIGN" keyword
	Rparen     token.Pos // position of ")" after reference columns

	Name             *Ident // optional
	Columns          []*Ident
	ReferenceTable   *Ident
	ReferenceColumns []*Ident
}

// IndexKey is index key specifier in CREATE TABLE and CREATE INDEX.
//
//     {{.Name | sql}} {{.Dir}}
type IndexKey struct {
	// pos = Name.Pos
	// end = DirPos + len(Dir) || Name.end

	DirPos token.Pos // position of Dir

	Name *Ident
	Dir  Direction // optional
}

// Cluster is INTERLEAVE IN PARENT clause in CREATE TABLE.
//
//     , INTERLEAVE IN PARENT {{.TableName | sql}} {{.OnDelete}}
type Cluster struct {
	// pos = Comma
	// end = OnDeleteEnd || TableName.end

	Comma       token.Pos // position of ","
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	TableName *Ident
	OnDelete  OnDeleteAction // optional
}

// AlterTable is ALTER TABLE statement node.
//
//     ALTER TABLE {{.Name | sql}} {{.TableAlternation | sql}}
type AlterTable struct {
	// pos = Alter
	// end = TableAlternation.end

	Alter token.Pos // position of "ALTER" keyword

	Name             *Ident
	TableAlternation TableAlternation
}

// AddColumn is ADD COLUMN clause in ALTER TABLE.
//
//     ADD COLUMN {{.Column | sql}}
type AddColumn struct {
	// pos = Add
	// end = Column.end

	Add token.Pos // position of "ADD" keyword

	Column *ColumnDef
}

// AddForeignKey is ADD FOREIGN KEY clause in ALTER TABLE.
//
//     ADD {{.ForeignKey | sql}}
type AddForeignKey struct {
	// pos = Add
	// end = ForeignKey.end

	Add token.Pos // position of "ADD" keyword

	ForeignKey *ForeignKey
}

// DropColumn is DROP COLUMN clause in ALTER TABLE.
//
//     DROP COLUMN {{.Name | sql}}
type DropColumn struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of  "DROP" keyword

	Name *Ident
}

// DropConstraint is DROP CONSTRAINT clause in ALTER TABLE.
//
//     DROP CONSTRAINT {{.Name | sql}}
type DropConstraint struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of  "DROP" keyword

	Name *Ident
}

// SetOnDelete is SET ON DELETE clause in ALTER TABLE.
//
//     SET ON DELETE {{.OnDelete}}
type SetOnDelete struct {
	// pos = Set
	// end = OnDeleteEnd

	Set         token.Pos // position of "SET" keyword
	OnDeleteEnd token.Pos // end position of ON DELETE clause

	OnDelete OnDeleteAction
}

// AlterColumn is ALTER COLUMN clause in ALTER TABLE.
//
//     ALTER COLUMN {{.Name | sql}} {{.Type | sql}} {{if .NotNull}}NOT NULL{{end}}
type AlterColumn struct {
	// pos = Alter
	// end = Null + 4 || Type.end

	Alter token.Pos // position of "ALTER" keyword
	Null  token.Pos // position of "NULL"

	Name    *Ident
	Type    SchemaType
	NotNull bool
}

// AlterColumnSet is ALTER COLUMN SET clause in ALTER TABLE.
//
//     ALTER COLUMN {{.Name | sql}} SET {{.Options | sql}}
type AlterColumnSet struct {
	// pos = Alter
	// end = Name.end

	Alter token.Pos // position of "ALTER" keyword

	Name    *Ident
	Options *ColumnDefOptions
}

// DropTable is DROP TABLE statement node.
//
//     DROP TABLE {{.Name | sql}}
type DropTable struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// CreateIndex is CREATE INDEX statement node.
//
//     CREATE
//       {{if .Unique}}UNIQUE{{end}}
//       {{if .NullFiltered}}NULL_FILTERED{{end}}
//       INDEX {{.Name | sql}} ON {{.TableName | sql}} (
//         {{.Keys | sqlJoin ","}}
//       )
//       {{.Storing | sqlOpt}}
//       {{.InterleaveIn | sqlOpt}}
type CreateIndex struct {
	// pos = Create
	// end = (InterleaveIn ?? Storing).end || Rparen + 1

	Create token.Pos // position of "CREATE" keyword
	Rparen token.Pos // position of ")"

	Unique       bool
	NullFiltered bool
	Name         *Ident
	TableName    *Ident
	Keys         []*IndexKey
	Storing      *Storing      // optional
	InterleaveIn *InterleaveIn // optional
}

// Storing is STORING clause in CREATE INDEX.
//
//     STORING ({{.Columns | sqlJoin ","}})
type Storing struct {
	// pos = Storing
	// end = Rparen + 1

	Storing token.Pos // position of "STORING" keyword
	Rparen  token.Pos // position of ")"

	Columns []*Ident
}

// InterleaveIn is INTERLEAVE IN clause in CREATE INDEX.
//
//     , INTERLEAVE IN {{.TableName | sql}}
type InterleaveIn struct {
	// pos = Comma
	// end = TableName.end

	Comma token.Pos // position of ","

	TableName *Ident
}

// DropIndex is DROP INDEX statement node.
//
//     DROP INDEX {{.Name | sql}}
type DropIndex struct {
	// pos = Drop
	// end = Name.end

	Drop token.Pos // position of "DROP" keyword

	Name *Ident
}

// ================================================================================
//
// Types for Schema
//
// ================================================================================

// ScalarSchemaType is scalar type node in schema.
//
//     {{.Name}}
type ScalarSchemaType struct {
	// pos = NamePos
	// end = NamePos + len(Name)

	NamePos token.Pos // position of this name

	Name ScalarTypeName // except for StringTypeName and BytesTypeName
}

// SizedSchemaType is sized type node in schema.
//
//     {{.Name}}({{if .Max}}MAX{{else}}{{.Size | sql}}{{end}})
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
//     ARRAY<{{.Item | sql}}>
type ArraySchemaType struct {
	// pos = Array
	// end = Gt + 1

	Array token.Pos // position of "ARRAY" keyword
	Gt    token.Pos // position of ">"

	Item SchemaType // ScalarSchemaType or SizedSchemaType
}

// ================================================================================
//
// DML
//
// ================================================================================

// Insert is INSERT statement node.
//
//     INSERT INTO {{.TableName | sql}} ({{.Columns | sqlJoin ","}}) {{.Input | sql}}
type Insert struct {
	// pos = Insert
	// end = Input.end

	Insert token.Pos // position of "INSERT" keyword

	TableName *Ident
	Columns   []*Ident
	Input     InsertInput
}

// ValuesInput is VALUES clause in INSERT.
//
//     VALUES {{.Rows | sqlJoin ","}}
type ValuesInput struct {
	// pos = Values
	// end = Rows[$].end

	Values token.Pos // position of "VALUES" keyword

	Rows []*ValuesRow
}

// ValuesRow is row values of VALUES clause.
//
//     ({{.Exprs | sqlJoin ","}})
type ValuesRow struct {
	// pos = Lparen
	// end = Rparen + 1

	Lparen, Rparen token.Pos // position of "(" and ")"

	Exprs []*DefaultExpr
}

// DefaultExpr is DEFAULT or Expr.
//
//     {{if .Default}}DEFAULT{{else}}{{.Expr | sql}}{{end}}
type DefaultExpr struct {
	// pos = DefaultPos || Expr.pos
	// end = DefaultPos + 7 || Expr.end

	DefaultPos token.Pos // position of "DEFAULT"

	Default bool
	Expr    Expr
}

// SubQueryInput is query clause in INSERT.
//
//     {{.Query | sql}}
type SubQueryInput struct {
	// pos = Query.pos
	// end = Query.end

	Query QueryExpr
}

// Delete is DELETE statement.
//
//     DELETE FROM {{.TableName | sql}} {{.As | sqlOpt}} {{.Where | sql}}
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
//     UPDATE {{.TableName | sql}} {{.As | sqlOpt}}
//     SET {{.Updates | sqlJoin ","}} {{.Where | sql}}
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
//     {{.Path | sqlJoin "."}} = {{.Expr | sql}}
type UpdateItem struct {
	// pos = Path[0].pos
	// end = Expr.end

	Path []*Ident // len(Path) > 0
	Expr Expr
}
