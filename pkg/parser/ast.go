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

func (BinaryExpr) isExpr()   {}
func (UnaryExpr) isExpr()    {}
func (InExpr) isExpr()       {}
func (IsNullExpr) isExpr()   {}
func (IsBoolExpr) isExpr()   {}
func (BetweenExpr) isExpr()  {}
func (SelectorExpr) isExpr() {}
func (IndexExpr) isExpr()    {}
func (CallExpr) isExpr()     {}
func (CastExpr) isExpr()     {}
func (CaseExpr) isExpr()     {}
func (SubQuery) isExpr()     {}
func (ParenExpr) isExpr()    {}
func (ArrayExpr) isExpr()    {}
func (ExistsExpr) isExpr()   {}
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

type QueryExpr interface {
	Node
	isQueryExpr()
}

func (Select) isQueryExpr()        {}
func (CompoundQuery) isQueryExpr() {}

type JoinExpr interface {
	Node
	isJoinExpr()
}

func (TableName) isJoinExpr()        {}
func (Unnest) isJoinExpr()           {}
func (SubQueryJoinExpr) isJoinExpr() {}
func (Join) isJoinExpr()             {}

type IntValue interface {
	Node
	isIntValue()
}

func (IntLit) isIntValue() {}
func (Param) isIntValue()  {}

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
	end  Pos
	Expr Expr
	Star bool
	As   *Ident // It must be nil when Star is true
}

type SelectExprList []*SelectExpr

// {{.Expr | sql}} {{if .TableSample}}TABLESAMPLE {{.TableSanple}}{{end}}
type FromItem struct {
	end         Pos
	Expr        JoinExpr
	TableSample TableSampleMethod
}

type FromItemList []*FromItem

// {{.Name}}
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}} AS {{.As | sql}}{{end}}
type TableName struct {
	pos, end Pos
	Name     string
	Hint     *Hint
	As       *Ident
}

// UNNEST({{.Expr | sql}})
//   {{if .Hint}}{{.Hint | sql}}{{end}}
//   {{if .As}}AS {{.As | sql}}{{end}}
//   {{if .WithOffset}}WITH OFFSET {{.WithOffset | sql}}{{end}}
type Unnest struct {
	pos, end   Pos
	Expr       Expr
	Hint       *Hint
	As         *Ident
	WithOffset *Ident
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

//   {{.Left | sql}}
// {{.Op}} JOIN
//    {{if .Hint}}{{.Hint | sql}}{{end}}
//    {{.Right | sql}}
// {{if .Cond}}{{.Cond | sql}}{{end}}
type Join struct {
	Op          JoinOp
	Hint        *Hint // If this is HASH JOIN, its JOIN_METHOD value is set as HASH_JOIN.
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
	end       Pos
	Expr      Expr
	Direction Direction
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

// {{.Left | sql}}[{{.Right | sql}}]
type IndexExpr struct {
	end         Pos
	Left, Right Expr
}

// {{.Func | sql}}({{.Args | sql}})
type CallExpr struct {
	end  Pos
	Func *Ident
	Args ExprList
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

type IdentList []Ident

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
	pos, end  Pos
	Name      TypeName
	MaxLength bool
	Length    IntValue       // for STRING(n) and BYTES(n)
	Fields    []*FieldSchema // for STRUCT<...>
	Value     *Type          // for ARRAY<...>
}

type FieldSchema struct {
	Name *Ident
	Type *Type
}
