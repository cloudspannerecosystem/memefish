package parser

type Node interface {
	Pos() Pos
	End() Pos
	SQL() string
}

// Expr repersents an expression in SQL.
type Expr interface {
	Node
	isExpr()
}

// ExprList is expressions separated by comma.
type ExprList []Expr

func (BinaryExpr) isExpr()    {}
func (UnaryExpr) isExpr()     {}
func (InExpr) isExpr()        {}
func (IsNullExpr) isExpr()    {}
func (IsBoolExpr) isExpr()    {}
func (BetweenExpr) isExpr()   {}
func (SelectorExpr) isExpr()  {}
func (IndexExpr) isExpr()     {}
func (CallExpr) isExpr()      {}
func (CountStarExpr) isExpr() {}
func (CastExpr) isExpr()      {}
func (ExtractExpr) isExpr()   {}
func (CaseExpr) isExpr()      {}
func (SubQuery) isExpr()      {}
func (ParenExpr) isExpr()     {}
func (ArrayExpr) isExpr()     {}
func (ExistsExpr) isExpr()    {}
func (Param) isExpr()         {}
func (Ident) isExpr()         {}
func (ArrayLit) isExpr()      {}
func (StructLit) isExpr()     {}
func (NullLit) isExpr()       {}
func (BoolLit) isExpr()       {}
func (IntLit) isExpr()        {}
func (FloatLit) isExpr()      {}
func (StringLit) isExpr()     {}
func (BytesLit) isExpr()      {}
func (DateLit) isExpr()       {}
func (TimestampLit) isExpr()  {}

type QueryExpr interface {
	Node
	setOrderBy(orderBy OrderExprList)
	setLimit(limit *Limit)
}

func (s *Select) setOrderBy(orderBy OrderExprList) {
	s.OrderBy = orderBy
	s.end = orderBy.End()
}

func (s *Select) setLimit(limit *Limit) {
	s.Limit = limit
	s.end = limit.End()
}

func (c *CompoundQuery) setOrderBy(orderBy OrderExprList) {
	c.OrderBy = orderBy
	c.end = orderBy.End()
}

func (c *CompoundQuery) setLimit(limit *Limit) {
	c.Limit = limit
	c.end = limit.End()
}

func (s *SubQueryExpr) setOrderBy(orderBy OrderExprList) {
	s.OrderBy = orderBy
	s.end = orderBy.End()
}

func (s *SubQueryExpr) setLimit(limit *Limit) {
	s.Limit = limit
	s.end = limit.End()
}

type JoinExpr interface {
	Node
	isJoinExpr()
}

func (TableName) isJoinExpr()        {}
func (Unnest) isJoinExpr()           {}
func (PathExpr) isJoinExpr()         {}
func (SubQueryJoinExpr) isJoinExpr() {}
func (ParenJoinExpr) isJoinExpr()    {}
func (Join) isJoinExpr()             {}

type IntValue interface {
	Node
	isIntValue()
}

func (IntLit) isIntValue()       {}
func (Param) isIntValue()        {}
func (CastIntValue) isIntValue() {}

type NumValue interface {
	Node
	isNumValue()
}

func (IntLit) isNumValue()       {}
func (FloatLit) isNumValue()     {}
func (Param) isNumValue()        {}
func (CastNumValue) isNumValue() {}

// {{if .Hint}}{{.Hint | sql}}{{end}}
// {{.Expr | sql}}
type QueryStatement struct {
	Hint *Hint
	Expr QueryExpr
}

type Hint struct {
	pos, end Pos
	Map      map[string]Expr
}

// SELECT
//   {{if .Distinct}}DISTINCT{{end}}
//   {{if .AsStruct}}AS STRUCT{{end}}
//   {{.List | sql}}
//   {{if .From}}FROM {{.From | sql}}{{end}}
//   {{if .Where}}WHERE {{.Where | sql}}{{end}}
//   {{if .GroupBy}}GROUP BY {{.GroupBy | sql}}{{end}}
//   {{if .Having}}HAVING {{.Having | sql}}{{end}}
//   {{if .OrderBy}}ORDER BY {{.OrderBy | sql}}{{end}}
//   {{if .Limit}}LIMIT {{.Limit | sql}}{{end}}
type Select struct {
	pos, end Pos
	Distinct bool
	AsStruct bool
	List     SelectExprList
	From     FromItemList
	Where    Expr
	GroupBy  ExprList // Integer literal on GROUP BY expression has special meaning.
	Having   Expr
	OrderBy  OrderExprList
	Limit    *Limit
}

// {{range $i, $e := .List}}
//   {{if ne($i, 0)}}{{.Op}} {{if .Distinct}}DISTINCT{{else}}ALL{{end}} {{.Right | sql}}{{end}}
//   {{$e | sql}}
// {{end}}
//   {{if .OrderBy}}ORDER BY {{.OrderBy | sql}}{{end}}
//   {{if .Limit}}LIMIT {{.Limit | sql}}{{end}}
type CompoundQuery struct {
	end      Pos
	Op       SetOp
	Distinct bool
	List     []QueryExpr
	OrderBy  OrderExprList
	Limit    *Limit
}

// {{.Expr | sql}}
//   {{if .OrderBy}}ORDER BY {{.OrderBy | sql}}{{end}}
//   {{if .Limit}}LIMIT {{.Limit | sql}}{{end}}
type SubQueryExpr struct {
	end     Pos
	Expr    *SubQuery
	OrderBy OrderExprList
	Limit   *Limit
}

// {{if .Expr}}{{.Expr | sql}}{{end}}{{if .Star}}{{if .Expr}}.{{end}}*{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
type SelectExpr struct {
	pos, end Pos
	Expr     Expr
	Star     bool
	As       *Ident // It must be nil when Star is true
}

type SelectExprList []*SelectExpr

// {{.Expr | sql}}
// {{if .TableSample}}TABLESAMPLE {{.Method}} ({{.Num | sql}} {{if .Rows}}ROWS{{else}}PERCENT{{end}}){{end}}
type FromItem struct {
	end    Pos
	Expr   JoinExpr
	Method TableSampleMethod
	Num    NumValue
	Rows   bool
}

type FromItemList []*FromItem

// {{.Ident | sql}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}} AS {{.As | sql}}{{end}}
type TableName struct {
	end   Pos
	Ident *Ident
	Hint  *Hint
	As    *Ident
}

// {{.Ident | sql}}{{range _, $path := .Paths}}.{{$path | sql}}{{end}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}} AS {{.As | sql}}{{end}}
type PathExpr struct {
	end   Pos
	Ident *Ident
	Paths []*Ident // not IdentList because it is not comma separated.
	Hint  *Hint
	As    *Ident
}

// UNNEST({{.Expr | sql}})
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
//   {{if .WithOffset}}WITH OFFSET {{.WithOffset | sql}}
//     {{if .WithOffsetAs}}{{.WithOffsetAs | sql}}{{end}}{{end}}
type Unnest struct {
	pos, end     Pos
	Expr         Expr
	Hint         *Hint
	As           *Ident
	WithOffset   bool
	WithOffsetAs *Ident
}

// {{.Expr | sql}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
type SubQueryJoinExpr struct {
	end  Pos
	Expr *SubQuery
	Hint *Hint
	As   *Ident
}

// ({{.Expr | sql}})
type ParenJoinExpr struct {
	pos, end Pos
	Expr     JoinExpr // SubQuery or Join
}

//   {{.Left | sql}}
// {{.Op}} {{if .Method}}{{.Method}}{{end}} JOIN
//    {{if .Hint}}{{.Hint | sql}}{{end}}
//    {{.Right | sql}}
// {{if .Cond}}{{.Cond | sql}}{{end}}
type Join struct {
	Op          JoinOp
	Method      JoinMethod
	Hint        *Hint
	Left, Right JoinExpr
	Cond        *JoinCondition
}

// {{if .On}}ON {{.On | sql}}{{end}}
// {{if .Using}}USING ({{.Using | sql}}){{end}}
type JoinCondition struct {
	pos, end Pos
	// Either On or Using must be non-empty.
	On    Expr
	Using IdentList
}

// {{.Expr | sql}} {{.Direction}}
type OrderExpr struct {
	end  Pos
	Expr Expr
	Dir  Direction
}

type OrderExprList []*OrderExpr

// {{.Count | sql}} {{if .Offset}}OFFSET {{.Offset | sql}}{{end}}
type Limit struct {
	Count  IntValue
	Offset IntValue
}

// {{.Left | sql}} {{.Op}} {{.Right | sql}}
type BinaryExpr struct {
	Op          BinaryOp
	Left, Right Expr
}

// {{.Op}} {{.Expr | sql}}
type UnaryExpr struct {
	pos  Pos
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
// {{if .Subqyery}}{{.Subquery | sql}}{{end}}
// {{if .Values}}({{.Values | sql}}){{end}}
type InCondition struct {
	pos, end Pos
	Unnest   Expr
	SubQuery *SubQuery
	Values   ExprList
}

// {{.Left | sql}} IS{{if .Not}} NOT{{end}} NULL
type IsNullExpr struct {
	end  Pos
	Not  bool
	Left Expr
}

// {{.Left | sql}} IS{{if .Not}} NOT{{end}} {{if .Right}}TRUE{{else}}FALSE{{end}}
type IsBoolExpr struct {
	end   Pos
	Not   bool
	Left  Expr
	Right bool
}

// {{.Left | sql}} {{if .Not}}NOT }}BETWEEN {{.RightStart | sql}} AND {{.RightEnd | sql}}
type BetweenExpr struct {
	Not                        bool
	Left, RightStart, RightEnd Expr
}

// {{.Left | sql}} . {{.Right | sql}}
type SelectorExpr struct {
	Left  Expr
	Right *Ident
}

// {{.Left | sql}}{{if .Ordinal}}ORDINAL{{else}}OFFSET{{end}}([{{.Right | sql}})]
type IndexExpr struct {
	end         Pos
	Ordinal     bool
	Left, Right Expr
}

// {{.Func | sql}}({{if .Distinct}}DISTINCT {{end}}{{.Args | sql}})
type CallExpr struct {
	end      Pos
	Func     *Ident
	Distinct bool
	Args     ArgList
}

// {{if .IntervalUnit}}INTERVAL {{.Expr | sql}} {{.IntervalUnit}}{{else}}{{.Expr | sql}}{{end}}
type Arg struct {
	pos, end     Pos
	IntervalUnit ExtractPart
	Expr         Expr
}

type ArgList []*Arg

// COUNT(*)
type CountStarExpr struct {
	pos, end Pos
}

// EXTRACT({{.Part}} FROM {{.Expr | sql}}{{if .TimeZone}} AT TIME ZONE {{.TimeZone | sql}}{{end}})
type ExtractExpr struct {
	pos, end Pos
	Part     ExtractPart
	Expr     Expr
	TimeZone Expr
}

// CAST({{.Expr | sql}} AS {{.Type | sql}})
type CastExpr struct {
	pos, end Pos
	Expr     Expr
	Type     *Type
}

// CASE {{if .Expr}}{{.Expr | sql}}{{end}}
//   {{range .When}}WHEN {{. | sql}}{{end}}
//   {{if .Else}}ELSE {{.Else | sql}}{{end}
// END
type CaseExpr struct {
	pos, end Pos
	Expr     Expr
	When     []*When
	Else     Expr
}

// {{.Cond | sql}} THEN {{.Then | sql}}
type When struct {
	Cond, Then Expr
}

// ({{. | sql}})
type SubQuery struct {
	pos, end Pos
	Expr     QueryExpr
}

// ({{. | sql}})
type ParenExpr struct {
	pos, end Pos
	Expr     Expr
}

// ARRAY{{.Expr | sql}}
type ArrayExpr struct {
	pos  Pos
	Expr *SubQuery
}

// EXISTS {{if .Hint}}{{.Hint | sql}}{{end}} {{.Expr | sql}}
type ExistsExpr struct {
	pos  Pos
	Hint *Hint
	Expr *SubQuery
}

// @{{.Name}}
type Param struct {
	pos  Pos
	Name string
}

// {{.Name}}
type Ident struct {
	pos, end Pos
	Name     string
}

type IdentList []*Ident

// ARRAY{{if .Type}}<{{.Type | sql}}>{{end}}[{{.Values | sql}}]
type ArrayLit struct {
	pos, end Pos
	Type     *Type
	Values   ExprList
}

// STRUCT{{if .Type}}<{{.Fields | sql}}>{{end}}({{.Values | sql}})
type StructLit struct {
	pos, end Pos
	Fields   []*FieldSchema
	Values   ExprList
}

// NULL
type NullLit struct {
	pos Pos
}

// {{if .Value}}TRUE{{else}}FALSE{{end}}
type BoolLit struct {
	pos   Pos
	Value bool
}

// {{.Value}}
type IntLit struct {
	pos, end Pos
	Value    string
}

// {{.Value}}
type FloatLit struct {
	pos, end Pos
	Value    string
}

// {{.Value | sqlQuote}}
type StringLit struct {
	pos, end Pos
	Value    string
}

// B{{. | sqlQuote}}
type BytesLit struct {
	pos, end Pos
	Value    []byte
}

// DATE{{. | sqlQuote}}
type DateLit struct {
	pos, end Pos
	Value    string
}

// TIMESTAMP{{. | sqlQuote}}
type TimestampLit struct {
	pos, end Pos
	Value    string
}

// TODO: separate more accurate types like SimpleType, ArrayType, etc.

type Type struct {
	pos, end Pos
	Name     TypeName
	Fields   []*FieldSchema // for STRUCT<...>
	Value    *Type          // for ARRAY<...>
}

type FieldSchema struct {
	Name *Ident
	Type *Type
}

// CAST({{.Expr | sql}} AS INT64)
type CastIntValue struct {
	pos, end Pos
	Expr     IntValue // IntLit or Param
}

// CAST({{.Expr | sql}} AS {{.Type}})
type CastNumValue struct {
	pos, end Pos
	Expr     NumValue // IntLit, FloatLit or Param
	Type     TypeName // Int64Type or Float64Type
}
