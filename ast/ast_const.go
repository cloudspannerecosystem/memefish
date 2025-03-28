package ast

// NOTE: This file defines constants used in AST nodes and they are used for automatic generation,
//       so this file is conventional.
//
// Convention:
//
//   - Each const types should be defined as a string type.
//   - Each value is defined as a string literal.

// AllOrDistinct represents ALL or DISTINCT in SELECT or set operations, etc.
// If it is optional, it may be an empty string, so handle it according to the context.
type AllOrDistinct string

const (
	AllOrDistinctAll      AllOrDistinct = "ALL"
	AllOrDistinctDistinct AllOrDistinct = "DISTINCT"
)

type JoinMethod string

const (
	HashJoinMethod   JoinMethod = "HASH"
	LookupJoinMethod JoinMethod = "LOOKUP" // Undocumented, but the GoogleSQL can parse this value.
)

type PositionKeyword string

const (
	PositionKeywordOffset      PositionKeyword = "OFFSET"
	PositionKeywordSafeOffset  PositionKeyword = "SAFE_OFFSET"
	PositionKeywordOrdinal     PositionKeyword = "ORDINAL"
	PositionKeywordSafeOrdinal PositionKeyword = "SAFE_ORDINAL"
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
	BernoulliSampleMethod TableSampleMethod = "BERNOULLI"
	ReservoirSampleMethod TableSampleMethod = "RESERVOIR"
)

type TableSampleUnit string

const (
	PercentTableSampleUnit TableSampleUnit = "PERCENT"
	RowsTableSampleUnit    TableSampleUnit = "ROWS"
)

type JoinOp string

const (
	CommaJoin      JoinOp = ","
	CrossJoin      JoinOp = "CROSS JOIN"
	InnerJoin      JoinOp = "INNER JOIN"
	FullOuterJoin  JoinOp = "FULL OUTER JOIN"
	LeftOuterJoin  JoinOp = "LEFT OUTER JOIN"
	RightOuterJoin JoinOp = "RIGHT OUTER JOIN"
)

type BinaryOp string

const (
	OpOr            BinaryOp = "OR"
	OpAnd           BinaryOp = "AND"
	OpEqual         BinaryOp = "="
	OpNotEqual      BinaryOp = "!="
	OpLess          BinaryOp = "<"
	OpGreater       BinaryOp = ">"
	OpLessEqual     BinaryOp = "<="
	OpGreaterEqual  BinaryOp = ">="
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
	OpConcat        BinaryOp = "||"
)

type UnaryOp string

const (
	OpNot    UnaryOp = "NOT"
	OpPlus   UnaryOp = "+"
	OpMinus  UnaryOp = "-"
	OpBitNot UnaryOp = "~"
)

type ScalarTypeName string

const (
	BoolTypeName      ScalarTypeName = "BOOL"
	Int64TypeName     ScalarTypeName = "INT64"
	Float32TypeName   ScalarTypeName = "FLOAT32"
	Float64TypeName   ScalarTypeName = "FLOAT64"
	StringTypeName    ScalarTypeName = "STRING"
	BytesTypeName     ScalarTypeName = "BYTES"
	DateTypeName      ScalarTypeName = "DATE"
	TimestampTypeName ScalarTypeName = "TIMESTAMP"
	NumericTypeName   ScalarTypeName = "NUMERIC"
	JSONTypeName      ScalarTypeName = "JSON"
	TokenListTypeName ScalarTypeName = "TOKENLIST"
)

type OnDeleteAction string

const (
	OnDeleteCascade  OnDeleteAction = "ON DELETE CASCADE"
	OnDeleteNoAction OnDeleteAction = "ON DELETE NO ACTION"
)

type Enforcement string

const (
	Enforced    Enforcement = "ENFORCED"
	NotEnforced Enforcement = "NOT ENFORCED"
)

type SecurityType string

const (
	SecurityTypeInvoker SecurityType = "INVOKER"
	SecurityTypeDefiner SecurityType = "DEFINER"
)

type InsertOrType string

const (
	InsertOrTypeUpdate InsertOrType = "UPDATE"
	InsertOrTypeIgnore InsertOrType = "IGNORE"
)
