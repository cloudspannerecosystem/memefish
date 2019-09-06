package parser

type Node interface {
	Pos() Pos
	End() Pos
	SQL() string
}

// QueryExpr represents set operator operands.
type QueryExpr interface {
	Node
	setQueryExprSuffix(orderBy *OrderBy, limit *Limit)
}

var _ QueryExpr = &Select{}
var _ QueryExpr = &SubQuery{}
var _ QueryExpr = &CompoundQuery{}

// SelectItem represents expression in SELECT clause result columns list.
type SelectItem interface {
	Node
	isSelectItem()
}

func (Star) isSelectItem()           {}
func (StarPath) isSelectItem()       {}
func (Alias) isSelectItem()          {}
func (ExprSelectItem) isSelectItem() {}

// JoinExpr represents JOIN operands.
type JoinExpr interface {
	Node
	isSimpleJoinExpr() bool
}

var _ JoinExpr = &Unnest{}
var _ JoinExpr = &TableName{}
var _ JoinExpr = &SubQueryJoinExpr{}
var _ JoinExpr = &ParenJoinExpr{}
var _ JoinExpr = &Join{}

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

// IntValue is integer values in SQL.
type IntValue interface {
	Node
	isIntValue()
}

func (Param) isIntValue()        {}
func (IntLiteral) isIntValue()   {}
func (CastIntValue) isIntValue() {}

// NumValue is number values in SQL.
type NumValue interface {
	Node
	isNumValue()
}

func (Param) isNumValue()        {}
func (IntLiteral) isNumValue()   {}
func (FloatLiteral) isNumValue() {}
func (CastNumValue) isNumValue() {}

// StringValue is string value in SQL.
type StringValue interface {
	Node
	isStringValue()
}

func (Param) isStringValue()         {}
func (StringLiteral) isStringValue() {}

// ================================================================================
//
// SELECT
//
// ================================================================================

// QueryStatement is query statement node.
//
//     {{if .Hint}}{{.Hint | sql}}{{end}}
//     {{.Expr | sql}}
type QueryStatement struct {
	// pos = (Hint ?? Expr).pos, end = Expr.end

	Hint  *Hint // optional
	Query QueryExpr
}

// Hint is hint node.
//
//     @{{"{"}}{{.Records | sqlJoin ","}}{{"}"}}
type Hint struct {
	pos, end Pos

	Records []*HintRecord // len(Records) > 0
}

// HintRecord is hint record node.
//
//     {{.Key | sql}}={{.Value | sql}}
type HintRecord struct {
	// pos = Key.pos, end = Value.end

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
	// end = (Limit ?? OrderBy ?? Having ?? GroupBy ?? Where ?? From ?? Results[$]).end
	pos Pos

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

func (s *Select) setQueryExprSuffix(orderBy *OrderBy, limit *Limit) {
	s.OrderBy = orderBy
	s.Limit = limit
}

// CompoundQuery is query statement node compounded by set operators.
//
//     {{.Queries | sqlJoin (printf "%s %s" .Op or(and(.Distinct, "DISTINCT"), "ALL"))}}
//       {{.OrderBy | sqlOpt}}
//       {{.Limit | sqlOpt}}
type CompoundQuery struct {
	// pos = Queries[0].pos, end = (Limit ?? OrderBy ?? Queries[$]).end

	Op       SetOp
	Distinct bool
	Queries  []QueryExpr // len(List) >= 2
	OrderBy  *OrderBy    // optional
	Limit    *Limit      // optional
}

func (c *CompoundQuery) setQueryExprSuffix(orderBy *OrderBy, limit *Limit) {
	c.OrderBy = orderBy
	c.Limit = limit
}

// SubQuery is subquery statement node.
//
//     ({{.Expr | sql}} {{.OrderBy | sqlOpt}} {{.Limit | sqlOpt}})
type SubQuery struct {
	pos, end Pos

	Query   QueryExpr
	OrderBy *OrderBy // optional
	Limit   *Limit   // optional
}

func (s *SubQuery) setQueryExprSuffix(orderBy *OrderBy, limit *Limit) {
	s.OrderBy = orderBy
	s.Limit = limit
}

// Star is a single * in SELECT result columns list.
//
//     *
type Star struct {
	// end = pos + 1
	pos Pos
}

// StarPath is expression with * in SELECT result columns list.
//
//     {{.Expr | sql}}.*
type StarPath struct {
	// pos = Expr.pos
	end Pos

	Expr Expr
}

// Alias is aliased expression by AS clause in SELECT result columns list.
//
//    {{.Expr | sql}} {{.As | sql}}
type Alias struct {
	// pos = Expr.pos, end = As.end

	Expr Expr
	As   *AsAlias
}

// AsAlias is AS clause node for general purpose.
//
// It is used in Alias node and some JoinExpr nodes.
//
// NOTE: Sometime keyword AS can be omited.
//       In this case, it.Pos() == it.Alias.Pos(), so we can detect this.
//
//     AS {{.Alias | sql}}
type AsAlias struct {
	// end = Alias.End
	pos Pos

	Alias *Ident
}

// ExprSelectItem is Expr wrapper for SelectItem.
//
//     {{.Expr | sql}}
type ExprSelectItem struct {
	// pos = Expr.pos, end = Expr.end

	Expr Expr
}

// From is FROM clause node.
//
//     FROM {{.Items | sqlJoin ","}}
type From struct {
	// end = Items[$].end
	pos Pos

	Items []*FromItem // len(Items) > 0
}

// FromItem is FROM clause expression node.
//
//     {{.Source | sql}} {{.TableSample | sqlOpt}}
type FromItem struct {
	// pos = Expr.pos, end = (TableSample ?? Source).end

	Source      JoinExpr
	TableSample *TableSample
}

// TableSample is TABLESAMPLE clause node.
//
//     TABLESAMPLE {{.Method}} {{.Size | sql}}
type TableSample struct {
	// end = Size.end
	pos Pos

	Method TableSampleMethod
	Size   *TableSampleSize
}

// TableSampleSize is size part of TABLESAMPLE clause.
//
//     ({{.Value | sql}} {{.Unit}})
type TableSampleSize struct {
	pos, end Pos

	Value NumValue
	Unit  TableSampleUnit
}

// Where is WHERE clause node.
//
//    WHERE {{.Expr | sql}}
type Where struct {
	// end = Expr.end
	pos Pos

	Expr Expr
}

// GroupBy is GROUP BY clause node.
//
//    GROUP BY {{.Exprs | sqlJoin ","}}
type GroupBy struct {
	// end = Exprs[$].end
	pos Pos

	Exprs []Expr // len(Exprs) > 0
}

// Having is HAVING clause node.
//
//     HAVING {{.Expr | sql}}
type Having struct {
	// end = Expr.end
	pos Pos

	Expr Expr
}

// OrderBy is ORDER BY clause node.
//
//     ORDER BY {{.Items | sqlJoin ","}}
type OrderBy struct {
	// end = Items[$].end
	pos Pos

	Items []*OrderByItem // len(Items) > 0
}

// OrderByItem is expression node in ORDER BY clause list.
//
//     {{.Expr | sql}} {{.Collate | sqlOpt}} {{.Direction}}
type OrderByItem struct {
	// pos = Expr.pos
	end Pos

	Expr    Expr
	Collate *Collate  // optional
	Dir     Direction // optional
}

// Collate is COLLATE clause node in ORDER BY item.
//
//     COLLATE {{.Value | sql}}
type Collate struct {
	// end = Value.end
	pos Pos

	Value StringValue
}

// Limit is LIMIT clause node.
//
//     LIMIT {{.Count | sql}} {{.Offset | sqlOpt}}
type Limit struct {
	// end = (Offset ?? Count).end
	pos Pos

	Count  IntValue
	Offset *Offset // optional
}

// Offset is OFFSET clause node in LIMIT clause.
//
//     OFFSET {{.Value | sql}}
type Offset struct {
	// end = Value.end
	pos Pos

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
type Unnest struct {
	pos, end Pos

	Implicit   bool
	Expr       Expr        // Path or Ident when Implicit is true
	Hint       *Hint       // optional
	As         *AsAlias    // optional
	WithOffset *WithOffset // optional
}

func (Unnest) isSimpleJoinExpr() bool {
	return true
}

// WithOffset is WITH OFFSET clause node after UNNEST call.
//
//     WITH OFFSET {{.As | sqlOpt}}
type WithOffset struct {
	pos, end Pos

	As *AsAlias // optional
}

// TableName is table name node in FROM clause.
//
//     {{.Table | sql}} {{.Hint | sqlOpt}} {{.As | sqlOpt}}
type TableName struct {
	// pos = Table.pos, end = (As ?? Hint ?? Table).end

	Table *Ident
	Hint  *Hint    // optional
	As    *AsAlias // optional
}

func (t *TableName) isSimpleJoinExpr() bool {
	return true
}

// SubQueryJoinExpr is subquery inside JOIN expression.
//
//     ({{.Query | sql}}) {{.As | sqlOpt}}
type SubQueryJoinExpr struct {
	pos, end Pos

	Query QueryExpr
	As    *AsAlias // optional
}

func (s *SubQueryJoinExpr) isSimpleJoinExpr() bool {
	return s.As != nil
}

// ParenJoinExpr is parenthesized JOIN expression.
//
//     ({{.Expr | sql}})
type ParenJoinExpr struct {
	pos, end Pos

	Source JoinExpr // SubQueryJoinExpr (without As) or Join
}

func (ParenJoinExpr) isSimpleJoinExpr() bool {
	return true
}

// Join is JOIN expression.
//
//       {{.Left | sql}}
//     {{.Op}} {{.Method}} JOIN {{.Hint | sqlOpt}}
//        {{.Right | sql}}
//     {{.Cond | sqlOpt}}
type Join struct {
	// pos = Left.pos, end = (Cond ?? Right).pos

	Op          JoinOp
	Method      JoinMethod
	Hint        *Hint // optional
	Left, Right JoinExpr
	Cond        JoinCondition // nil when Op is CrossJoin, otherwise it must be set.
}

func (Join) isSimpleJoinExpr() bool {
	return false
}

// On is ON condition of JOIN expression.
//
//     ON {{.Expr | sql}}
type On struct {
	// end = Expr.end
	pos Pos

	Expr Expr
}

// Using is Using condition of JOIN expression.
//
//     USING ({{Idents | sqlJoin ","}})
type Using struct {
	pos, end Pos

	Idents []*Ident // len(Idents) > 0
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
	// pos = Left.pos, end = Right.pos
	Op BinaryOp

	Left, Right Expr
}

// UnaryExpr is unary operator expression node.
//
//     {{.Op}} {{.Expr | sql}}
type UnaryExpr struct {
	// end = Expr.end
	pos Pos

	Op   UnaryOp
	Expr Expr
}

// InExpr is IN expression node.
//
//     {{.Left | sql}} {{if .Not}}NOT{{end}} IN {{.Right | sql}}
type InExpr struct {
	Not   bool
	Left  Expr
	Right InCondition
}

// UnnestInCondition is UNNEST call at IN condition.
//
//     UNNEST({{.Expr | sql}})
type UnnestInCondition struct {
	pos, end Pos

	Expr Expr
}

// SubQueryInCondition is subquery at IN condition.
//
//     ({{.Query | sql}})
type SubQueryInCondition struct {
	pos, end Pos

	Query QueryExpr
}

// ValuesInCondition is parenthesized values at IN condition.
//
//     ({{.Exprs | sqlJoin ","}})
type ValuesInCondition struct {
	pos, end Pos

	Exprs []Expr // len(Exprs) > 0
}

// IsNullExpr is IS NULL expression node.
//
//     {{.Left | sql}} IS {{if .Not}}NOT{{end}} NULL
type IsNullExpr struct {
	// pos = Expr.pos
	end Pos

	Not  bool
	Left Expr
}

// IsBoolExpr is IS TRUE/FALSE expression node.
//
//     {{.Left | sql}} IS {{if .Not}}NOT{{end}} {{if .Right}}TRUE{{else}}FALSE{{end}}
type IsBoolExpr struct {
	// pos = Expr.pos
	end Pos

	Not   bool
	Left  Expr
	Right bool
}

// BetweenExpr is BETWEEN expression node.
//
//     {{.Left | sql}} {{if .Not}}NOT{{end}} BETWEEN {{.RightStart | sql}} AND {{.RightEnd | sql}}
type BetweenExpr struct {
	// pos = Left.pos, end = RightEnd.end

	Not                        bool
	Left, RightStart, RightEnd Expr
}

// SelectorExpr is struct field access expression node.
//
//     {{.Expr | sql}}.{{.Member | sql}}
type SelectorExpr struct {
	// pos = Expr.pos, end = Member.pos

	Expr   Expr
	Member *Ident
}

// IndexExpr is array item access expression node.
//
//     {{.Expr | sql}}[{{if .Ordinal}}ORDINAL{{else}}OFFSET{{end}}({{.Index | sql}})]
type IndexExpr struct {
	// pos = Expr.pos
	end Pos

	Ordinal     bool
	Expr, Index Expr
}

// CallExpr is function call expression node.
//
//     {{.Func | sql}}({{if .Distinct}}DISTINCT{{end}} {{.Args | sql}})
type CallExpr struct {
	// pos = Func.pos
	end Pos

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
	// end = (IntervalUnit ?? Expr).end
	pos Pos

	Expr         Expr
	IntervalUnit *Ident // optional
}

// CountStarExpr is node just for COUNT(*).
//
//     COUNT(*)
type CountStarExpr struct {
	pos, end Pos
}

// ExtractExpr is EXTRACT call expression node.
//
//     EXTRACT({{.Part | sql}} FROM {{.Expr | sql}} {{.AtTimeZone | sqlOpt}})
type ExtractExpr struct {
	pos, end Pos

	Part       *Ident
	Expr       Expr
	AtTimeZone *AtTimeZone // optional
}

// AtTimeZone is AT TIME ZONE clause in EXTRACT call.
//
//     AT TIME ZONE {{.Expr | sql}}
type AtTimeZone struct {
	// end = Expr.end
	pos Pos

	Expr Expr
}

// CastExpr is CAST call expression node.
//
//     CAST({{.Expr | sql}} AS {{.Type | sql}})
type CastExpr struct {
	pos, end Pos

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
	pos, end Pos

	Expr  Expr        // optional
	Whens []*CaseWhen // len(Whens) > 0
	Else  *CaseElse   // optional
}

// CaseWhen is WHEN clause in CASE expression.
//
//     WHEN {{.Cond | sql}} THEN {{.Then | sql}}
type CaseWhen struct {
	// end = Then.end
	pos Pos

	Cond, Then Expr
}

// CaseElse is ELSE clause in CASE expression.
//
//     ELSE {{.Expr | sql}}
type CaseElse struct {
	// end = Expr.end
	pos Pos

	Expr Expr
}

// ParenExpr is parenthesized expression node.
//
//     ({{. | sql}})
type ParenExpr struct {
	pos, end Pos

	Expr Expr
}

// ScalarSubQuery is subquery in expression.
//
//     ({{.Query | sql}})
type ScalarSubQuery struct {
	pos, end Pos

	Query QueryExpr
}

// ArraySubQuery is subquery in ARRAY call.
//
//     ARRAY({{.Query | sql}})
type ArraySubQuery struct {
	pos, end Pos

	Query QueryExpr
}

// ExistsSubQuery is subquery in EXISTS call.
//
//     EXISTS {{.Hint | sqlOpt}} ({{.Query | sql}})
type ExistsSubQuery struct {
	pos, end Pos
	Hint     *Hint
	Query    QueryExpr
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
	// end = pos + 1 + len(Name)
	pos Pos

	Name string
}

// Ident is identifier node.
//
//     {{.Name | sqlIdentQuote}}
type Ident struct {
	pos, end Pos

	Name string
}

// Path is dot-chained identifier list.
//
//     {{.Idents | sqlJoin "."}}
type Path struct {
	// pos = Idents[0].pos, end = idents[$].end

	Idents []*Ident // len(Idents) >= 2
}

// AraryLiteral is array literal node.
//
//     ARRAY{{if .Type}}<{{.Type | sql}}>{{end}}[{{.Values | sqlJoin ","}}]
type ArrayLiteral struct {
	pos, end Pos

	Type   Type   // optional
	Values []Expr // len(Values) > 0
}

// StructLiteral is struct literal node.
//
//     STRUCT{{if .Type}}<{{.Fields | sqlJoin ","}}>{{end}}({{.Values | sqlJoin ","}})
type StructLiteral struct {
	pos, end Pos

	// NOTE: Distinguish nil from len(Fields) == 0 case.
	//       nil means type is not specified, or empty slice means this struct has 0 fields.
	Fields []*FieldType
	Values []Expr
}

// NullLiteral is just NULL literal.
//
//     NULL
type NullLiteral struct {
	// end = pos + 4
	pos Pos
}

// BoolLiteral is boolean literal node.
//
//     {{if .Value}}TRUE{{else}}FALSE{{end}}
type BoolLiteral struct {
	// end = pos + (Value ? 4 : 5)
	pos Pos

	Value bool
}

// IntLiteral is integer literal node.
//
//     {{.Value}}
type IntLiteral struct {
	pos, end Pos

	Base  int
	Value string
}

// FloatLiteral is floating point number literal node.
//
//     {{.Value}}
type FloatLiteral struct {
	pos, end Pos

	Value string
}

// StringLiteral is string literal node.
//
//     {{.Value | sqlStringQuote}}
type StringLiteral struct {
	pos, end Pos

	Value string
}

// BytesLiteral is bytes literal node.
//
//     B{{.Value | sqlByesQuote}}
type BytesLiteral struct {
	pos, end Pos

	Value []byte
}

// DateLiteral is date literal node.
//
//     DATE {{.Value | sqlStringQuote}}
type DateLiteral struct {
	pos, end Pos

	Value string
}

// TimestampLiteral is timestamp literal node.
//
//     TIMESTAMP {{.Value | sqlStringQuote}}
type TimestampLiteral struct {
	pos, end Pos

	Value string
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
	// end = pos + len(Name)
	pos Pos

	Name TypeName // except for ArrayTypeName and StructTypeName
}

// ArrayType is array type node.
//
//     ARRAY<{{.Item | sql}}>
type ArrayType struct {
	pos, end Pos

	Item Type
}

// StructType is struct type node.
//
//     STRUCT<{{.Fields | sql}}>
type StructType struct {
	pos, end Pos

	Fields []*FieldType
}

// FieldType is field type in struct type node.
//
//     {{if .Member}}{{.Member | sql}}{{}}
type FieldType struct {
	// pos = (Member ?? Type).pos, end = Type.end

	Member *Ident
	Type   Type
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
	pos, end Pos

	Expr IntValue // IntLit or Param
}

// CasrNumValue is cat call in number value context.
//
//     CAST({{.Expr | sql}} AS {{.Type}})
type CastNumValue struct {
	pos, end Pos

	Expr NumValue // IntLit, FloatLit or Param
	Type TypeName // Int64Type or Float64Type
}
