package ast

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
	HashJoinMethod  JoinMethod = "HASH"
	ApplyJoinMethod JoinMethod = "APPLY"
	LoopJoinMethod  JoinMethod = "LOOP" // Undocumented, but the Spanner accept this value at least.
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

type GQLAllOrDistinctEnum string

const (
	GQLAllOrDistinctImplicitAll GQLAllOrDistinctEnum = ""
	GQLAllOrDistinctAll         GQLAllOrDistinctEnum = "ALL"
	GQLAllOrDistinctDistinct    GQLAllOrDistinctEnum = "DISTINCT"
)

type GQLSetOperatorEnum string

const (
	GQLSetOperatorUnion     GQLSetOperatorEnum = "UNION"
	GQLSetOperatorIntersect GQLSetOperatorEnum = "INTERSECT"
	GQLSetOperatorExcept    GQLSetOperatorEnum = "EXCEPT"
)

type GQLDirectionEnum string

const (
	GQLSortOrderUnspecified GQLDirectionEnum = ""
	GQLSortOrderAsc         GQLDirectionEnum = "ASC"
	GQLSortOrderAscending   GQLDirectionEnum = "ASCENDING"
	GQLSortOrderDesc        GQLDirectionEnum = "DESC"
	GQLSortOrderDescending  GQLDirectionEnum = "DESCENDING"
)

type GQLPathModeEnum string

const (
	GQLPathModeUnspecified GQLPathModeEnum = ""
	GQLPathModeWalk        GQLPathModeEnum = "WALK"
	GQLPathModeTrail       GQLPathModeEnum = "TRAIL"
)

type GQLSearchPrefixEnum string

const (
	GQLPathSearchPrefixAll         GQLSearchPrefixEnum = "ALL"
	GQLPathSearchPrefixAny         GQLSearchPrefixEnum = "ANY"
	GQLPathSearchPrefixAnyShortest GQLSearchPrefixEnum = "ANY SHORTEST"
)
