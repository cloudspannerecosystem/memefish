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
	BernoulliSampleMethod TableSampleMethod = "BERNOULLI"
	ReservoirSampleMethod TableSampleMethod = "RESERVOIR"
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
