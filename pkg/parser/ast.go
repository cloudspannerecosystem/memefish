package parser

type TableHintKey string

const (
	ForceIndexTableHint     TableHintKey = "FORCE_INDEX"
	GroupScanByOptimization TableHintKey = "GROUPBY_SCAN_OPTIMIZATION"
)

type JoinHintKey string

const (
	ForceJoinOrderJoinHint JoinHintKey = "FORCE_JOIN_ORDER"
	JoinTypeJoinHint       JoinHintKey = "JOIN_TYPE"
)

type JoinMethod string

const (
	HashJoinMethod  JoinMethod = "HASH_JOIN"
	ApplyJoinMethod JoinMethod = "APPLY_JOIN"
	LoopJoinMethod  JoinMethod = "LOOP_JOIN" // Undocumented, but the Spanner accept this value at least.
)

type SetOp string

const (
	SetOpUnion     SetOp = "UNION"
	SetOpIntersect SetOp = "INTERSECT"
	SetOpExcept    SetOp = "EXCEPT"
)

type Direction string

const (
	DirectionAsc  Direction = "ASC"
	DirectionDesc Direction = "DESC"
)

type TableSampleMethod string

const (
	BernoulliSampleMethod TableSampleMethod = "bernoulli"
	ReservoirSampleMethod TableSampleMethod = "reservoir"
)

type JoinOp string

const (
	InnerJoin      JoinOp = "INNER"
	CrossJoin      JoinOp = "CROSS"
	FullOuterJoin  JoinOp = "FULL OUTER"
	LeftOuterJoin  JoinOp = "LEFT OUTER"
	RightOuterJoin JoinOp = "RIGHT OUTER"
)

type BinaryOp string

const (
	OpOr            BinaryOp = "OR"
	OpAnd           BinaryOp = "AND"
	OpEqual         BinaryOp = "="
	OpLessThan      BinaryOp = "<"
	OpGreaterThan   BinaryOp = ">"
	OpLessEqual     BinaryOp = "<="
	OpGreaterEqual  BinaryOp = ">="
	OpNotEqual      BinaryOp = "!="
	OpLike          BinaryOp = "LIKE"
	OpNotLike       BinaryOp = "NOT LIKE"
	OpBitOr         BinaryOp = "|"
	OpBitXor        BinaryOp = "^"
	OpBitAnd        BinaryOp = "&"
	OpBitLeftShift  BinaryOp = "<<"
	OpBitRightShift BinaryOp = ">>"
	OpAdd           BinaryOp = "+"
	OpSub           BinaryOp = "-"
	OpMul           BinaryOp = "*"
	OpDiv           BinaryOp = "/"
)

type UnaryOp string

const (
	OpNot    UnaryOp = "NOT"
	OpPlus   UnaryOp = "+"
	OpMinus  UnaryOp = "-"
	OpBitNot UnaryOp = "~"
)

type TypeName string

const (
	BoolType      TypeName = "BOOL"
	Int64Type     TypeName = "INT64"
	Float64Type   TypeName = "FLOAT64"
	StringType    TypeName = "STRING"
	BytesType     TypeName = "BYTES"
	DateType      TypeName = "DATE"
	TimestampType TypeName = "TIMESTAMP"
	ArrayType     TypeName = "ARRAY"
	StructType    TypeName = "STRUCT"
)

// Expr repersents an expression in SQL.
type Expr interface {
	isExpr()
}

// ExprList is expressions separated by comma.
type ExprList []Expr

func (BinaryExpr) isExpr()   {}
func (UnaryExpr) isExpr()    {}
func (InExpr) isExpr()       {}
func (IsNullExpr) isExpr()   {}
func (IsBoolExpr) isExpr()   {}
func (PathExpr) isExpr()     {}
func (IndexExpr) isExpr()    {}
func (CallExpr) isExpr()     {}
func (CaseExpr) isExpr()     {}
func (SubQuery) isExpr()     {}
func (Param) isExpr()        {}
func (Ident) isExpr()        {}
func (ArrayLit) isExpr()     {}
func (StructLit) isExpr()    {}
func (NullLit) isExpr()      {}
func (BoolLit) isExpr()      {}
func (IntLit) isExpr()       {}
func (FloatLit) isExpr()     {}
func (StringLit) isExpr()    {}
func (BytesLit) isExpr()     {}
func (DateLit) isExpr()      {}
func (TimestampLit) isExpr() {}

// {{if .TableHint}}{{.TableHint | sql}}{{end}}
// {{if .JoinHint}}{{.JoinHint | sql}}{{end}}
// {{.QueryExpr | sql}}
type QueryStatement struct {
	TableHint TableHint
	JoinHint  JoinHint
	QueryExpr QueryExpr
}

type TableHint map[TableHintKey]interface{}

type JoinHint map[JoinHintKey]interface{}

type QueryExpr interface {
	isQueryExpr()
}

func (Select) isQueryExpr()        {}
func (CompoundQuery) isQueryExpr() {}

// SELECT
//   {{if .Distinct}}DISTINCT{{end}}
//   {{if .AsStruct}}AS STRUCT{{end}}
//   {{.List | sql}}
//   {{if .From}}FROM {{.From | sql}}{{end}}
//   {{if .Where}}WHERE {{.Where | sql}}{{end}}
//   {{if .GroupBy}}GROUP BY {{.GroupBy | sql}}{{end}}
//   {{if .OrderBy}}ORDER BY {{.OrderBt | sql}}{{end}}
//   {{if .Limit}}LIMIT {{.Limit | sql}}{{end}}
type Select struct {
	Distinct bool
	AsStruct bool // On top-level, it must be false.
	List     SelectExprList
	From     FromItemList
	Where    Expr
	GroupBy  ExprList
	Having   Expr
	OrderBy  OrderExprList
	Limit    *Limit
}

// {{.Left | sql}} {{.Op}} {{if .Distinct}}DISTINCT{{else}}ALL{{end}} {{.Right | sql}}
type CompoundQuery struct {
	Op          SetOp
	Distinct    bool
	Left, Right QueryExpr
}

// {{if .Expr}}{{.Expr | sql}}{{end}}{{if .Star}}{{if .Expr}}.{{end}}*{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
type SelectExpr struct {
	Expr Expr
	Star bool
	As   Ident // It must be nil when Star is true
}

type SelectExprList []*SelectExprList

type FromItem struct {
	Expr              JoinExpr
	TableSampleMethod TableSampleMethod
}

type FromItemList []*FromItem

type JoinExpr interface {
	isJoinExpr()
}

func (TableName) isJoinExpr()        {}
func (Unnest) isJoinExpr()           {}
func (SubQueryJoinExpr) isJoinExpr() {}
func (Join) isJoinExpr()             {}

// {{.Name}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}} AS {{.As | sql}}{{end}}
type TableName struct {
	Name string
	Hint TableHint
	As   Ident
}

// UNNEST({{.Expr | sql}})
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
//   {{if .WithOffset}}WITH OFFSET {{.WithOffset | sql}}{{end}}
type Unnest struct {
	Expr       Expr
	Hint       TableHint
	As         Ident
	WithOffset Ident
}

// {{.Expr | sql}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
type SubQueryJoinExpr struct {
	Expr *SubQuery
	Hint TableHint
	As   Ident
}

//   {{.Left | sql}}
// {{.Op}} JOIN
//    {{if .Hint}}{{.Hint | sql}}{{end}}
//    {{.Right | sql}}
// {{.Cond | sql}}
type Join struct {
	Op          JoinOp
	Hint        JoinHint // If this is HASH JOIN, its JOIN_METHOD value is set as HASH_JOIN.
	Left, Right JoinExpr
	Cond        *JoinCondition
}

// {{if .On}}ON {{.On | sql}}{{end}}
// {{if .Using}}USING ({{.Using | sql}}){{end}}
type JoinCondition struct {
	// Either On or Using must be non-empty.
	On    Expr
	Using IdentList
}

// {{.Expr | sql}} {{.Direction}}
type OrderExpr struct {
	Expr      Expr
	Direction Direction
}

type OrderExprList []*OrderExpr

// {{.Count | sql}} {{if .Offset}}{{.Offset | sql}}{{end}}
type Limit struct {
	Count  IntValue
	Offset IntValue
}

type IntValue interface {
	isIntValue()
}

func (IntLit) isIntValue() {}
func (Param) isIntValue()  {}

// {{.Left | sql}} {{.Op}} {{.Right | sql}}
type BinaryExpr struct {
	Op          BinaryOp
	Left, Right Expr
}

// {{.Op}} {{.Expr | sql}}
type UnaryExpr struct {
	Op   UnaryOp
	Expr Expr
}

// {{.Left | sql}} {{if .Not}}NOT {{end}}IN {{.Right | sql}}
type InExpr struct {
	Not   bool
	Left  Expr
	Right *InCondition
}

// {{if .Unnest}}UNNEST({{.Unnest | sql}}){{end}}
// {{if .Subqyery}}({{.Subquery | sql}}){{end}}
// {{if .Values}}({{.Values | sql}}){{end}}
type InCondition struct {
	Unnest   Expr
	SubQuery SubQuery
	Values   ExprList
}

// {{.Left | sql}} IS{{if .Not}} NOT{{end}} NULL
type IsNullExpr struct {
	Not  bool
	Left Expr
}

// {{.Left | sql}} IS{{if .Not}} NOT{{end}} {{if .Right}}TRUE{{else}}FALSE{{end}}
type IsBoolExpr struct {
	Not   bool
	Left  Expr
	Right bool
}

// {{.Left | sql}} . {{.Right | sql}}
type PathExpr struct {
	Left  Expr
	Right Ident
}

// {{.Left | sql}}[{{.Right | sql}}]
type IndexExpr struct {
	Left, Right Expr
}

// {{.Func | sql}}({{.Args | sql}})
type CallExpr struct {
	Func Ident
	Args ExprList
}

// CASE {{.Expr | sql}}
//   {{range .When}}WHEN {{. | sql}}{{end}}
//   {{if .Else}}ELSE {{.Else | sql}}{{end}
type CaseExpr struct {
	Expr Expr
	When []*When
	Else Expr
}

// {{.Cond | sql}} THEN {{.Then | sql}}
type When struct {
	Cond, Then Expr
}

// ({{. | sql}})
type SubQuery struct {
	Expr QueryExpr
}

// @{{.Name}}
type Param struct {
	Name string
}

// `{{.Ident}}`
type Ident string

type IdentList string

// ARRAY{{if .Type}}<{{.Type | sql}}>{{end}}[{{.Values | sql}}]
type ArrayLit struct {
	Type   Type
	Values ExprList
}

// STRUCT{{if .Type}}<{{.Type | sql}}>{{end}}({{.Values | sql}})
type StructLit struct {
	Type   Type
	Values ExprList
}

// NULL
type NullLit struct{}

// {{if .}}TRUE{{else}}FALSE{{end}}
type BoolLit bool

// {{.}}
type IntLit int64

// {{.}}
//
// NOTE: when the value is NaN, this text representation is CAST('NaN' AS FLOAT64), and +/-infinity are respectively.
type FloatLit float64

// {{. | sqlQuote}}
type StringLit string

// B{{. | sqlQuote}}
type BytesLit []byte

// DATE{{. | sqlQuote}}
type DateLit string

// TIMESTAMP{{. | sqlQuote}}
type TimestampLit string

type Type struct {
	Name      TypeName
	Length    IntValue // for STRING(n) and BYTES(n)
	Fields    []*Field // for STRUCT<...>
	ValueType *Type    // for ARRAY<...>
}

type Field struct {
	Name Ident
	Type *Type
}
