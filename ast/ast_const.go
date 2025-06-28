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
	IntervalTypeName  ScalarTypeName = "INTERVAL"
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

// DateTimePart is used for:
//
//   - INTERVAL literals
//   - EXTRACT for DATE, TIMESTAMP, INTERVAL (TODO)
//   - DATE_TRUNC, TIMESTAMP_TRUNC (TODO)
type DateTimePart string

// Definition of all enum values (including values not used in Spanner)
// https://github.com/google/zetasql/blob/2025.03.1/zetasql/public/functions/datetime.proto#L26
const (
	DateTimePartYear        DateTimePart = "YEAR"
	DateTimePartMonth       DateTimePart = "MONTH"
	DateTimePartDay         DateTimePart = "DAY"
	DateTimePartDayOfWeek   DateTimePart = "DAYOFWEEK"
	DateTimePartDayOfYear   DateTimePart = "DAYOFYEAR"
	DateTimePartQuarter     DateTimePart = "QUARTER"
	DateTimePartHour        DateTimePart = "HOUR"
	DateTimePartMinute      DateTimePart = "MINUTE"
	DateTimePartSecond      DateTimePart = "SECOND"
	DateTimePartMillisecond DateTimePart = "MILLISECOND"
	DateTimePartMicrosecond DateTimePart = "MICROSECOND"
	DateTimePartNanosecond  DateTimePart = "NANOSECOND"
	DateTimePartWeek        DateTimePart = "WEEK"
	DateTimePartISOYear     DateTimePart = "ISOYEAR"
	DateTimePartISOWeek     DateTimePart = "ISOWEEK"
	DateTimePartDate        DateTimePart = "DATE"

	// Not yet used in Spanner?

	DateTimePartDateTime      DateTimePart = "DATETIME"
	DateTimePartTime          DateTimePart = "TIME"
	DateTimePartWeekMonday    DateTimePart = "WEEK(MONDAY)"
	DateTimePartWeekTuesday   DateTimePart = "WEEK(TUESDAY)"
	DateTimePartWeekWednesday DateTimePart = "WEEK(WEDNESDAY)"
	DateTimePartWeekThursday  DateTimePart = "WEEK(THURSDAY)"
	DateTimePartWeekFriday    DateTimePart = "WEEK(FRIDAY)"
	DateTimePartWeekSaturday  DateTimePart = "WEEK(SATURDAY)"
)
