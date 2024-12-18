package ast

import (
	"strconv"
	"strings"

	"github.com/cloudspannerecosystem/memefish/token"
)

// ================================================================================
//
//  Experimental format with indentation
//
// ================================================================================

// formatOption is container of format configurations.
// It is conceptually immutable.
type formatOption struct {
	newline bool
	indent  int
}

// formatOptionCompact is format option without newline and indentation.
func formatOptionCompact() formatOption {
	return formatOption{}
}

// formatOptionPretty is format option with newline and configured indentation.
func formatOptionPretty(indent int) formatOption {
	return formatOption{newline: true, indent: indent}
}

// FormatContext is container of format option and current indentation.
// Note: All methods of FormatContext must support nil receiver.
type FormatContext struct {
	option        formatOption
	currentIndent int
}

// FormatContextCompact is format context without newline and indentation.
func FormatContextCompact() *FormatContext {
	return &FormatContext{option: formatOptionCompact()}
}

// FormatContextPretty is format context with newline and configured indentation.
func FormatContextPretty(indent int) *FormatContext {
	return &FormatContext{option: formatOptionPretty(indent)}
}

// SQL is entry point of pretty printing.
// If node implements NodeFormat, it calls NodeFormat.sqlContext() instead of Node.SQL().
func (fc *FormatContext) SQL(node Node) string {
	if nodeFormat, ok := node.(NodeFormat); fc != nil && ok {
		return nodeFormat.sqlContext(fc)
	} else {
		return node.SQL()
	}
}

// newlineOr returns newline with indentation if formatOptionPretty is used.
// Otherwise, it returns argument string.
func (fc *FormatContext) newlineOr(s string) string {
	if fc == nil {
		return s
	}

	return strIfElse(fc.option.newline, "\n", s) + strings.Repeat(" ", fc.currentIndent)
}

// indentScope executes function with FormatContext with deeper indentation.
func (fc *FormatContext) indentScope(f func(fc *FormatContext) string) string {
	if fc == nil {
		return f(nil)
	}

	newFc := *fc
	newFc.currentIndent += fc.option.indent
	return f(&newFc)
}

// sqlOptCtx is sqlOpt with FormatContext.
func sqlOptCtx[T interface {
	Node
	comparable
}](fc *FormatContext, left string, node T, right string) string {
	var zero T
	if node == zero {
		return ""
	}
	return left + fc.SQL(node) + right
}

// sqlJoinCtx is sqlJoin with FormatContext.
func sqlJoinCtx[T Node](fc *FormatContext, elems []T, sep string) string {
	var b strings.Builder
	for i, r := range elems {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(fc.SQL(r))
	}
	return b.String()
}

// NodeFormat is Node with FormatContext support.
// If it is implemented, (*FormatContext).SQL calls sqlContext() instead of SQL()
type NodeFormat interface {
	Node

	// sqlContext is Node.SQL() with FormatContext conceptually.
	// It can be called with nil.
	// Note: It would become to Node.SQL() finally.
	sqlContext(fmtCtx *FormatContext) string
}

// ================================================================================
//
// Helper functions for SQL()
// These functions are intended for use within this file only.
//
// ================================================================================

// sqlOpt outputs:
//
//	when node != nil: left + node.SQL() + right
//	else            : empty string
//
// This function corresponds to sqlOpt in ast.go
func sqlOpt[T interface {
	Node
	comparable
}](left string, node T, right string) string {
	return sqlOptCtx(nil, left, node, right)
}

// strOpt outputs:
//
//	when pred == true: s
//	else             : ""
//
// This function corresponds to {{if pred}}s{{end}} in ast.go
func strOpt(pred bool, s string) string {
	if pred {
		return s
	}
	return ""
}

// strIfElse outputs:
//
//	when pred == true: ifStr
//	else             : elseStr
//
// This function corresponds to {{if pred}}ifStr{{else}}elseStr{{end}} in ast.go
func strIfElse(pred bool, ifStr string, elseStr string) string {
	if pred {
		return ifStr
	}
	return elseStr
}

// sqlJoin outputs joined string of SQL() of all elems by sep.
// This function corresponds to sqlJoin in ast.go
func sqlJoin[T Node](elems []T, sep string) string {
	return sqlJoinCtx(nil, elems, sep)
}

// formatBoolUpper formats bool value as uppercase.
func formatBoolUpper(b bool) string {
	return strings.ToUpper(strconv.FormatBool(b))
}

type prec int

const (
	precLit prec = iota
	precSelector
	precUnary
	precMulDiv
	precAddSub
	precBitShift
	precBitAnd
	precBitXor
	precBitOr
	precComparison
	precNot
	precAnd
	precOr
)

func exprPrec(e Expr) prec {
	switch e := e.(type) {
	case *CallExpr, *CountStarExpr, *CastExpr, *ExtractExpr, *CaseExpr, *IfExpr, *ParenExpr, *ScalarSubQuery,
		*ArraySubQuery, *ExistsSubQuery, *Param, *Ident, *Path, *ArrayLiteral, *TupleStructLiteral, *TypedStructLiteral,
		*TypelessStructLiteral, *NullLiteral, *BoolLiteral, *IntLiteral, *FloatLiteral, *StringLiteral, *BytesLiteral,
		*DateLiteral, *TimestampLiteral, *NumericLiteral, *JSONLiteral, *WithExpr:
		return precLit
	case *IndexExpr, *SelectorExpr:
		return precSelector
	case *InExpr, *IsNullExpr, *IsBoolExpr, *BetweenExpr:
		return precComparison
	case *BinaryExpr:
		switch e.Op {
		case OpMul, OpDiv, OpConcat:
			return precMulDiv
		case OpAdd, OpSub:
			return precAddSub
		case OpBitLeftShift, OpBitRightShift:
			return precBitShift
		case OpBitAnd:
			return precBitAnd
		case OpBitXor:
			return precBitXor
		case OpBitOr:
			return precBitOr
		case OpEqual, OpNotEqual, OpLess, OpLessEqual, OpGreater, OpGreaterEqual, OpLike, OpNotLike:
			return precComparison
		case OpAnd:
			return precAnd
		case OpOr:
			return precOr
		}
	case *UnaryExpr:
		switch e.Op {
		case OpPlus, OpMinus, OpBitNot:
			return precUnary
		case OpNot:
			return precNot
		}
	}

	panic("exprPrec: unexpected")
}

func paren(p prec, e Expr) string {
	ep := exprPrec(e)
	if ep <= p {
		return e.SQL()
	} else {
		return "(" + e.SQL() + ")"
	}
}

// ================================================================================
//
// SELECT
//
// ================================================================================

func (q *QueryStatement) sqlContext(fc *FormatContext) string {
	return sqlOptCtx(fc, "", q.Hint, fc.newlineOr(" ")) +
		fc.SQL(q.Query)
}

func (q *QueryStatement) SQL() string {
	return q.sqlContext(nil)
}

func (q *Query) sqlContext(fc *FormatContext) string {
	return sqlOptCtx(fc, "", q.With, fc.newlineOr(" ")) +
		fc.SQL(q.Query) +
		sqlOptCtx(fc, fc.newlineOr(" "), q.OrderBy, "") +
		sqlOptCtx(fc, fc.newlineOr(" "), q.Limit, "") +
		strOpt(len(q.PipeOperators) > 0, fc.newlineOr(" ")) +
		sqlJoinCtx(fc, q.PipeOperators, fc.newlineOr(" "))
}

func (q *Query) SQL() string {
	return q.sqlContext(nil)
}

func (h *Hint) SQL() string {
	return "@{" + sqlJoin(h.Records, ", ") + "}"
}

func (h *HintRecord) SQL() string {
	return h.Key.SQL() + "=" + h.Value.SQL()
}

func (w *With) sqlContext(fc *FormatContext) string {
	return "WITH " + sqlJoinCtx(fc, w.CTEs, ", ")
}

func (w *With) SQL() string {
	return w.sqlContext(nil)
}
func (c *CTE) sqlContext(fc *FormatContext) string {
	return c.Name.SQL() + " AS (" +
		fc.indentScope(func(fc *FormatContext) string {
			return fc.newlineOr("") + fc.SQL(c.QueryExpr)
		}) +
		fc.newlineOr("") + ")"
}

func (c *CTE) SQL() string {
	return c.sqlContext(nil)
}

func (s *Select) sqlContext(fc *FormatContext) string {
	return "SELECT" +
		strOpt(s.AllOrDistinct != "", " "+string(s.AllOrDistinct)) +
		sqlOptCtx(fc, " ", s.As, "") +
		fc.indentScope(func(fc *FormatContext) string {
			if len(s.Results) == 1 {
				return " " + fc.SQL(s.Results[0])
			}
			return fc.newlineOr(" ") + sqlJoinCtx(fc, s.Results, ","+fc.newlineOr(" "))
		}) +
		sqlOptCtx(fc, fc.newlineOr(" "), s.From, "") +
		sqlOptCtx(fc, fc.newlineOr(" "), s.Where, "") +
		sqlOptCtx(fc, fc.newlineOr(" "), s.GroupBy, "") +
		sqlOptCtx(fc, fc.newlineOr(" "), s.Having, "")
}

func (s *Select) SQL() string {
	return s.sqlContext(nil)
}

func (a *AsStruct) SQL() string { return "AS STRUCT" }

func (a *AsValue) SQL() string { return "AS VALUE" }

func (a *AsTypeName) SQL() string { return "AS " + a.TypeName.SQL() }

func (c *CompoundQuery) SQL() string {
	return sqlJoin(c.Queries, " "+string(c.Op)+" "+strOpt(c.AllOrDistinct != "", string(c.AllOrDistinct)+" "))
}

func (s *SubQuery) SQL() string {
	return "(" + s.Query.SQL() + ")"
}

func (s *StarModifierExcept) SQL() string { return "EXCEPT (" + sqlJoin(s.Columns, " ") + ")" }

func (s *StarModifierReplaceItem) SQL() string { return s.Expr.SQL() + " AS " + s.Name.SQL() }

func (s *StarModifierReplace) SQL() string { return "REPLACE (" + sqlJoin(s.Columns, ", ") + ")" }

func (s *Star) SQL() string {
	return "*" + sqlOpt(" ", s.Except, "") + sqlOpt(" ", s.Replace, "")
}

func (s *DotStar) SQL() string {
	return s.Expr.SQL() + ".*" + sqlOpt(" ", s.Except, "") + sqlOpt(" ", s.Replace, "")
}

func (a *Alias) SQL() string {
	return a.Expr.SQL() + " " + a.As.SQL()
}

func (a *AsAlias) SQL() string {
	return strOpt(!a.As.Invalid(), "AS ") + a.Alias.SQL()
}

func (e *ExprSelectItem) SQL() string {
	return e.Expr.SQL()
}

func (f *From) sqlContext(fc *FormatContext) string {
	return "FROM " + fc.SQL(f.Source)
}
func (f *From) SQL() string {
	return f.sqlContext(nil)
}

func (w *Where) SQL() string {
	return "WHERE " + w.Expr.SQL()
}

func (g *GroupBy) SQL() string {
	return "GROUP BY " + sqlJoin(g.Exprs, ", ")
}

func (h *Having) SQL() string {
	return "HAVING " + h.Expr.SQL()
}

func (o *OrderBy) SQL() string {
	return "ORDER BY " + sqlJoin(o.Items, ", ")
}

func (o *OrderByItem) SQL() string {
	return o.Expr.SQL() +
		sqlOpt(" ", o.Collate, "") +
		strOpt(o.Dir != "", " "+string(o.Dir))
}

func (c *Collate) SQL() string {
	return "COLLATE " + c.Value.SQL()
}

func (l *Limit) SQL() string {
	return "LIMIT " + l.Count.SQL() +
		sqlOpt(" ", l.Offset, "")
}

func (o *Offset) SQL() string {
	return "OFFSET " + o.Value.SQL()
}

// ================================================================================
//
// Pipe Operators
//
// ================================================================================

func (p *PipeSelect) SQL() string {
	return "|> SELECT " + strOpt(p.AllOrDistinct != "", string(p.AllOrDistinct)+" ") + sqlOpt("", p.As, " ") + sqlJoin(p.Results, ", ")
}

func (p *PipeWhere) SQL() string {
	return "|> WHERE " + p.Expr.SQL()
}

// ================================================================================
//
// JOIN
//
// ================================================================================

func (u *Unnest) SQL() string {
	return "UNNEST(" + u.Expr.SQL() + ")" +
		sqlOpt("", u.Hint, "") +
		sqlOpt(" ", u.As, "") +
		sqlOpt(" ", u.WithOffset, "") +
		sqlOpt(" ", u.Sample, "")
}

func (w *WithOffset) SQL() string {
	return "WITH OFFSET" + sqlOpt(" ", w.As, "")
}

func (t *TableName) SQL() string {
	return t.Table.SQL() +
		sqlOpt(" ", t.Hint, "") +
		sqlOpt(" ", t.As, "") +
		sqlOpt(" ", t.Sample, "")
}

func (e *PathTableExpr) SQL() string {
	return e.Path.SQL() +
		sqlOpt("", e.Hint, "") +
		sqlOpt(" ", e.As, "") +
		sqlOpt(" ", e.WithOffset, "") +
		sqlOpt(" ", e.Sample, "")
}

func (s *SubQueryTableExpr) sqlContext(fc *FormatContext) string {
	return "(" +
		fc.indentScope(func(fc *FormatContext) string {
			return fc.newlineOr("") + fc.SQL(s.Query)
		}) +
		fc.newlineOr("") + ")" +
		sqlOptCtx(fc, " ", s.As, "") +
		sqlOptCtx(fc, " ", s.Sample, "")
}

func (s *SubQueryTableExpr) SQL() string {
	return s.sqlContext(nil)
}

func (p *ParenTableExpr) SQL() string {
	return "(" + p.Source.SQL() + ")" +
		sqlOpt(" ", p.Sample, "")
}

func (j *Join) sqlContext(fc *FormatContext) string {
	return fc.SQL(j.Left) +
		strOpt(j.Op != CommaJoin, fc.newlineOr(" ")) +
		string(j.Op) + " " +
		sqlOptCtx(fc, "", j.Hint, " ") +
		fc.SQL(j.Right) +
		sqlOptCtx(fc, " ", j.Cond, "")
}

func (j *Join) SQL() string {
	return j.sqlContext(nil)
}
func (o *On) SQL() string {
	return "ON " + o.Expr.SQL()
}

func (u *Using) SQL() string {
	return "USING (" + sqlJoin(u.Idents, ", ") + ")"
}

func (t *TableSample) SQL() string {
	return "TABLESAMPLE " + string(t.Method) + " " + t.Size.SQL()
}

func (t *TableSampleSize) SQL() string {
	return "(" + t.Value.SQL() + " " + string(t.Unit) + ")"
}

// ================================================================================
//
// Expr
//
// ================================================================================

func (b *BinaryExpr) SQL() string {
	p := exprPrec(b)

	return paren(p, b.Left) +
		" " + string(b.Op) + " " +
		paren(p, b.Right)
}

func (u *UnaryExpr) SQL() string {
	p := exprPrec(u)
	return string(u.Op) + strOpt(u.Op == OpNot, " ") + paren(p, u.Expr)
}

func (i *InExpr) SQL() string {
	p := exprPrec(i)
	return paren(p, i.Left) +
		strOpt(i.Not, " NOT") +
		" IN " + i.Right.SQL()
}

func (u *UnnestInCondition) SQL() string {
	return "UNNEST(" + u.Expr.SQL() + ")"
}

func (s *SubQueryInCondition) SQL() string {
	return "(" + s.Query.SQL() + ")"
}

func (v *ValuesInCondition) SQL() string {
	return "(" + sqlJoin(v.Exprs, ", ") + ")"
}

func (i *IsNullExpr) SQL() string {
	p := exprPrec(i)
	return paren(p, i.Left) +
		" IS " + strOpt(i.Not, "NOT ") + "NULL"
}

func (i *IsBoolExpr) SQL() string {
	p := exprPrec(i)
	return paren(p, i.Left) + " IS " + strOpt(i.Not, "NOT ") + formatBoolUpper(i.Right)
}

func (b *BetweenExpr) SQL() string {
	p := exprPrec(b)
	return paren(p, b.Left) +
		strOpt(b.Not, " NOT") +
		" BETWEEN " + paren(p, b.RightStart) + " AND " + paren(p, b.RightEnd)
}

func (s *SelectorExpr) SQL() string {
	p := exprPrec(s)
	return paren(p, s.Expr) + "." + s.Ident.SQL()
}

func (i *IndexExpr) SQL() string {
	p := exprPrec(i)
	return paren(p, i.Expr) + "[" + i.Index.SQL() + "]"
}

func (s *SubscriptSpecifierKeyword) SQL() string {
	return string(s.Keyword) + "(" + s.Expr.SQL() + ")"
}

func (c *CallExpr) SQL() string {
	return c.Func.SQL() + "(" + strOpt(c.Distinct, "DISTINCT ") +
		sqlJoin(c.Args, ", ") +
		strOpt(len(c.Args) > 0 && len(c.NamedArgs) > 0, ", ") +
		sqlJoin(c.NamedArgs, ", ") +
		sqlOpt(" ", c.NullHandling, "") +
		sqlOpt(" ", c.Having, "") +
		")" +
		sqlOpt(" ", c.Hint, "")
}

func (l *LambdaArg) SQL() string {
	// This implementation is not exactly matched with the doc comment for simplicity.
	return strOpt(!l.Lparen.Invalid(), "(") +
		sqlJoin(l.Args, ", ") +
		strOpt(!l.Lparen.Invalid(), ")") +
		" -> " +
		l.Expr.SQL()
}

func (c *TVFCallExpr) SQL() string {
	return c.Name.SQL() + "(" +
		sqlJoin(c.Args, ", ") +
		strOpt(len(c.Args) > 0 && len(c.NamedArgs) > 0, ", ") +
		sqlJoin(c.NamedArgs, ", ") +
		")" +
		sqlOpt(" ", c.Hint, "")
}

func (n *NamedArg) SQL() string { return n.Name.SQL() + " => " + n.Value.SQL() }

func (i *IgnoreNulls) SQL() string { return "IGNORE NULLS" }

func (r *RespectNulls) SQL() string { return "RESPECT NULLS" }

func (h *HavingMax) SQL() string { return "HAVING MAX " + h.Expr.SQL() }

func (h *HavingMin) SQL() string { return "HAVING MIN " + h.Expr.SQL() }

func (s *ExprArg) SQL() string {
	return s.Expr.SQL()
}

func (i *IntervalArg) SQL() string {
	return "INTERVAL " + i.Expr.SQL() + sqlOpt(" ", i.Unit, "")
}

func (s *SequenceArg) SQL() string {
	return "SEQUENCE " + s.Expr.SQL()
}

func (s *ModelArg) SQL() string {
	return "MODEL " + s.Name.SQL()
}

func (s *TableArg) SQL() string {
	return "TABLE " + s.Name.SQL()
}

func (*CountStarExpr) SQL() string {
	return "COUNT(*)"
}

func (e *ExtractExpr) SQL() string {
	return "EXTRACT(" + e.Part.SQL() + " FROM " + e.Expr.SQL() +
		sqlOpt(" ", e.AtTimeZone, "") + ")"
}

func (a *AtTimeZone) SQL() string {
	return "AT TIME ZONE " + a.Expr.SQL()
}

func (r *ReplaceFieldsArg) SQL() string { return r.Expr.SQL() + " AS " + r.Field.SQL() }

func (r *ReplaceFieldsExpr) SQL() string {
	return "REPLACE_FIELDS(" + r.Expr.SQL() + ", " + sqlJoin(r.Fields, ", ") + ")"
}

func (n *WithExprVar) SQL() string { return n.Name.SQL() + " AS " + n.Expr.SQL() }

func (w *WithExpr) SQL() string {
	return "WITH(" + sqlJoin(w.Vars, ", ") + ", " + w.Expr.SQL() + ")"
}

func (c *CastExpr) SQL() string {
	return strOpt(c.Safe, "SAFE_") + "CAST(" + c.Expr.SQL() + " AS " + c.Type.SQL() + ")"
}

func (c *CaseExpr) SQL() string {
	return "CASE " + sqlOpt("", c.Expr, " ") +
		sqlJoin(c.Whens, " ") + " " +
		sqlOpt("", c.Else, " ") +
		"END"
}

func (c *CaseWhen) SQL() string {
	return "WHEN " + c.Cond.SQL() + " THEN " + c.Then.SQL()
}

func (c *CaseElse) SQL() string {
	return "ELSE " + c.Expr.SQL()
}

func (i *IfExpr) SQL() string {
	return "IF(" + i.Expr.SQL() + ", " + i.TrueResult.SQL() + ", " + i.ElseResult.SQL() + ")"
}

func (p *ParenExpr) SQL() string {
	return "(" + p.Expr.SQL() + ")"
}

func (s *ScalarSubQuery) SQL() string {
	return "(" + s.Query.SQL() + ")"
}

func (a *ArraySubQuery) SQL() string {
	return "ARRAY(" + a.Query.SQL() + ")"
}

func (e *ExistsSubQuery) SQL() string {
	return "EXISTS" +
		sqlOpt(" ", e.Hint, " ") +
		"(" + e.Query.SQL() + ")"
}

func (p *Param) SQL() string {
	return "@" + p.Name
}

func (i *Ident) SQL() string {
	return token.QuoteSQLIdent(i.Name)
}

func (p *Path) SQL() string {
	return sqlJoin(p.Idents, ".")
}

func (a *ArrayLiteral) SQL() string {
	return strOpt(!a.Array.Invalid(), "ARRAY") +
		sqlOpt("<", a.Type, ">") +
		"[" + sqlJoin(a.Values, ", ") + "]"
}

func (s *TupleStructLiteral) SQL() string {
	return "(" + sqlJoin(s.Values, ", ") + ")"
}

func (s *TypedStructLiteral) SQL() string {
	return "STRUCT<" + sqlJoin(s.Fields, ", ") + ">(" + sqlJoin(s.Values, ", ") + ")"
}

func (s *TypelessStructLiteral) SQL() string {
	return strOpt(!s.Struct.Invalid(), "STRUCT") + "(" + sqlJoin(s.Values, ", ") + ")"
}

func (*NullLiteral) SQL() string {
	return "NULL"
}

func (b *BoolLiteral) SQL() string {
	return formatBoolUpper(b.Value)
}

func (i *IntLiteral) SQL() string {
	return i.Value
}

func (f *FloatLiteral) SQL() string {
	return f.Value
}

func (s *StringLiteral) SQL() string {
	return token.QuoteSQLString(s.Value)
}

func (b *BytesLiteral) SQL() string {
	return token.QuoteSQLBytes(b.Value)
}

func (d *DateLiteral) SQL() string {
	return "DATE " + d.Value.SQL()
}

func (t *TimestampLiteral) SQL() string {
	return "TIMESTAMP " + t.Value.SQL()
}

func (t *NumericLiteral) SQL() string {
	return "NUMERIC " + t.Value.SQL()
}

func (t *JSONLiteral) SQL() string {
	return "JSON " + t.Value.SQL()
}

// ================================================================================
//
// NEW constructors
//
// ================================================================================

func (n *NewConstructor) SQL() string {
	return "NEW " + n.Type.SQL() + "(" + sqlJoin(n.Args, ", ") + ")"
}

func (b *BracedNewConstructor) SQL() string {
	return "NEW " + b.Type.SQL() + " " + b.Body.SQL()
}

func (b *BracedConstructor) SQL() string {
	return "{" + sqlJoin(b.Fields, ", ") + "}"
}

func (b *BracedConstructorField) SQL() string {
	if _, ok := b.Value.(*BracedConstructor); ok {
		// Name {...}
		return b.Name.SQL() + " " + b.Value.SQL()
	}
	// Name: value
	return b.Name.SQL() + b.Value.SQL()
}

func (b *BracedConstructorFieldValueExpr) SQL() string { return ": " + b.Expr.SQL() }

// ================================================================================
//
// Type
//
// ================================================================================

func (s *SimpleType) SQL() string {
	return string(s.Name)
}

func (a *ArrayType) SQL() string {
	return "ARRAY<" + a.Item.SQL() + ">"
}

func (s *StructType) SQL() string {
	return "STRUCT<" + sqlJoin(s.Fields, ", ") + ">"
}

func (f *StructField) SQL() string {
	return sqlOpt("", f.Ident, " ") + f.Type.SQL()
}

func (n *NamedType) SQL() string {
	return sqlJoin(n.Path, ".")
}

// ================================================================================
//
// Cast for Special Cases
//
// ================================================================================

func (c *CastIntValue) SQL() string {
	return "CAST(" + c.Expr.SQL() + " AS INT64)"
}

func (c *CastNumValue) SQL() string {
	return "CAST(" + c.Expr.SQL() + " AS " + string(c.Type) + ")"
}

// ================================================================================
//
// DDL
//
// ================================================================================

func (g *Options) SQL() string { return "OPTIONS (" + sqlJoin(g.Records, ", ") + ")" }

func (g *OptionsDef) SQL() string {
	// Lowercase "null", "true", "false" is popular in option values.
	var valueSql string
	switch v := g.Value.(type) {
	case *NullLiteral:
		valueSql = "null"
	case *BoolLiteral:
		valueSql = strconv.FormatBool(v.Value)
	default:
		valueSql = g.Value.SQL()
	}
	return g.Name.SQL() + " = " + valueSql
}

func (c *CreateDatabase) SQL() string {
	return "CREATE DATABASE " + c.Name.SQL()
}

func (s *CreateSchema) SQL() string { return "CREATE SCHEMA " + s.Name.SQL() }

func (s *DropSchema) SQL() string { return "DROP SCHEMA " + s.Name.SQL() }

func (d *AlterDatabase) SQL() string {
	return "ALTER DATABASE " + d.Name.SQL() + " SET " + d.Options.SQL()
}

func (c *CreatePlacement) SQL() string {
	return "CREATE PLACEMENT " + c.Name.SQL() + sqlOpt(" ", c.Options, " ")
}

func (p *ProtoBundleTypes) SQL() string { return "(" + sqlJoin(p.Types, ", ") + ")" }

func (b *CreateProtoBundle) SQL() string { return "CREATE PROTO BUNDLE " + b.Types.SQL() }

func (a *AlterProtoBundle) SQL() string {
	return "ALTER PROTO BUNDLE" +
		sqlOpt(" ", a.Insert, "") +
		sqlOpt(" ", a.Update, "") +
		sqlOpt(" ", a.Delete, "")
}

func (a *AlterProtoBundleInsert) SQL() string { return "INSERT " + a.Types.SQL() }

func (a *AlterProtoBundleUpdate) SQL() string { return "UPDATE " + a.Types.SQL() }

func (a *AlterProtoBundleDelete) SQL() string { return "DELETE " + a.Types.SQL() }

func (d *DropProtoBundle) SQL() string { return "DROP PROTO BUNDLE" }

func (c *CreateTable) sqlContext(fc *FormatContext) string {
	return "CREATE TABLE " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		fc.SQL(c.Name) +
		" (" +
		fc.indentScope(func(fc *FormatContext) string {
			return fc.newlineOr("") + sqlJoinCtx(fc, c.Columns, ","+fc.newlineOr(" ")) +
				strOpt(len(c.Columns) > 0 && (len(c.TableConstraints) > 0 || len(c.Synonyms) > 0), ","+fc.newlineOr(" ")) +
				sqlJoinCtx(fc, c.TableConstraints, ","+fc.newlineOr(" ")) +
				strOpt(len(c.TableConstraints) > 0 && len(c.Synonyms) > 0, ","+fc.newlineOr(" ")) +
				sqlJoinCtx(fc, c.Synonyms, ","+fc.newlineOr(" "))
		}) +
		fc.newlineOr("") +
		") PRIMARY KEY (" +
		sqlJoinCtx(fc, c.PrimaryKeys, ", ") +
		")" +
		sqlOptCtx(fc, "", c.Cluster, "") +
		sqlOptCtx(fc, "", c.RowDeletionPolicy, "")
}

func (c *CreateTable) SQL() string {
	return c.sqlContext(nil)
}

func (s *Synonym) SQL() string { return "SYNONYM (" + s.Name.SQL() + ")" }

func (c *CreateSequence) SQL() string {
	return "CREATE SEQUENCE " + strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " " + c.Options.SQL()
}

func (c *AlterSequence) SQL() string {
	return "ALTER SEQUENCE " + c.Name.SQL() + " SET " + c.Options.SQL()
}

func (c *CreateView) SQL() string {
	return "CREATE" + strOpt(c.OrReplace, " OR REPLACE") + " VIEW " + c.Name.SQL() +
		" SQL SECURITY " + string(c.SecurityType) + " AS " + c.Query.SQL()
}

func (d *DropView) SQL() string { return "DROP VIEW " + d.Name.SQL() }

func (c *ColumnDef) SQL() string {
	return c.Name.SQL() + " " + c.Type.SQL() +
		strOpt(c.NotNull, " NOT NULL") +
		sqlOpt(" ", c.DefaultExpr, "") +
		sqlOpt(" ", c.GeneratedExpr, "") +
		strOpt(!c.Hidden.Invalid(), " HIDDEN") +
		sqlOpt(" ", c.Options, "")
}

func (c *TableConstraint) SQL() string {
	return sqlOpt("CONSTRAINT ", c.Name, " ") + c.Constraint.SQL()
}

func (f *ForeignKey) SQL() string {
	return "FOREIGN KEY (" + sqlJoin(f.Columns, ", ") + ") " +
		"REFERENCES " + f.ReferenceTable.SQL() + " (" +
		sqlJoin(f.ReferenceColumns, ", ") + ")" +
		strOpt(f.OnDelete != "", " "+string(f.OnDelete))
}

func (c *Check) SQL() string {
	return "CHECK (" + c.Expr.SQL() + ")"
}

func (c *ColumnDefaultExpr) SQL() string {
	return "DEFAULT (" + c.Expr.SQL() + ")"
}

func (g *GeneratedColumnExpr) SQL() string {
	return "AS (" + g.Expr.SQL() + ")" + strOpt(!g.Stored.Invalid(), " STORED")
}

func (i *IndexKey) SQL() string {
	return i.Name.SQL() + strOpt(i.Dir != "", " "+string(i.Dir))
}

func (c *Cluster) SQL() string {
	return ", INTERLEAVE IN PARENT " + c.TableName.SQL() +
		strOpt(c.OnDelete != "", " "+string(c.OnDelete))
}

func (c *CreateRowDeletionPolicy) SQL() string {
	return ", " + c.RowDeletionPolicy.SQL()
}

func (r *RowDeletionPolicy) SQL() string {
	return "ROW DELETION POLICY ( OLDER_THAN ( " + r.ColumnName.SQL() + ", INTERVAL " + r.NumDays.SQL() + " DAY ))"
}

func (a *AlterTable) SQL() string {
	return "ALTER TABLE " + a.Name.SQL() + " " + a.TableAlteration.SQL()
}

func (s *AddSynonym) SQL() string { return "ADD SYNONYM " + s.Name.SQL() }

func (s *DropSynonym) SQL() string { return "DROP SYNONYM " + s.Name.SQL() }

func (t *RenameTo) SQL() string { return "RENAME TO " + t.Name.SQL() + sqlOpt(", ", t.AddSynonym, "") }

func (a *AddColumn) SQL() string {
	return "ADD COLUMN " + strOpt(a.IfNotExists, "IF NOT EXISTS ") + a.Column.SQL()
}

func (a *AddTableConstraint) SQL() string {
	return "ADD " + a.TableConstraint.SQL()
}

func (a *AddRowDeletionPolicy) SQL() string {
	return "ADD " + a.RowDeletionPolicy.SQL()
}

func (d *DropColumn) SQL() string {
	return "DROP COLUMN " + d.Name.SQL()
}

func (d *DropConstraint) SQL() string {
	return "DROP CONSTRAINT " + d.Name.SQL()
}

func (d *DropRowDeletionPolicy) SQL() string {
	return "DROP ROW DELETION POLICY"
}

func (r *ReplaceRowDeletionPolicy) SQL() string {
	return "REPLACE " + r.RowDeletionPolicy.SQL()
}

func (s *SetOnDelete) SQL() string {
	return "SET " + string(s.OnDelete)
}

func (a *AlterColumn) SQL() string {
	return "ALTER COLUMN " + a.Name.SQL() + " " + a.Alteration.SQL()
}

func (a *AlterColumnType) SQL() string {
	return a.Type.SQL() +
		strOpt(a.NotNull, " NOT NULL") +
		sqlOpt(" ", a.DefaultExpr, "")
}

func (a *AlterColumnSetOptions) SQL() string { return "SET " + a.Options.SQL() }

func (a *AlterColumnSetDefault) SQL() string { return "SET " + a.DefaultExpr.SQL() }

func (a *AlterColumnDropDefault) SQL() string { return "DROP DEFAULT" }

func (d *DropTable) SQL() string {
	return "DROP TABLE " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

func (r *RenameTable) SQL() string { return "RENAME TABLE " + sqlJoin(r.Tos, ", ") }

func (r *RenameTableTo) SQL() string { return r.Old.SQL() + " TO " + r.New.SQL() }

func (c *CreateIndex) SQL() string {
	return "CREATE " +
		strOpt(c.Unique, "UNIQUE ") +
		strOpt(c.NullFiltered, "NULL_FILTERED ") +
		"INDEX " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " ON " + c.TableName.SQL() + "(" +
		sqlJoin(c.Keys, ", ") +
		")" +
		sqlOpt(" ", c.Storing, "") +
		sqlOpt("", c.InterleaveIn, "")
}

func (c *CreateVectorIndex) SQL() string {
	return "CREATE VECTOR INDEX " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " ON " + c.TableName.SQL() + " (" + c.ColumnName.SQL() + ") " +
		sqlOpt("", c.Where, " ") +
		c.Options.SQL()
}

func (c *CreateChangeStream) SQL() string {
	return "CREATE CHANGE STREAM " + c.Name.SQL() +
		sqlOpt(" ", c.For, "") +
		sqlOpt(" ", c.Options, "")
}

func (c *ChangeStreamForAll) SQL() string {
	return "FOR ALL"
}

func (c *ChangeStreamForTables) SQL() string {
	// TODO: Refactor after ChangeStreamForTable implements Node.
	sql := "FOR "
	for i, table := range c.Tables {
		if i > 0 {
			sql += ", "
		}
		sql += table.SQL()
	}
	return sql
}

func (a *AlterChangeStream) SQL() string {
	return "ALTER CHANGE STREAM " + a.Name.SQL() + " " + a.ChangeStreamAlteration.SQL()
}

func (a ChangeStreamSetFor) SQL() string {
	return "SET " + a.For.SQL()
}

func (a ChangeStreamDropForAll) SQL() string {
	return "DROP FOR ALL"
}

func (a ChangeStreamSetOptions) SQL() string {
	return "SET " + a.Options.SQL()
}

func (c *ChangeStreamForTable) SQL() string {
	return c.TableName.SQL() + strOpt(len(c.Columns) > 0, "("+sqlJoin(c.Columns, ", ")+")")
}

func (d *DropChangeStream) SQL() string {
	return "DROP CHANGE STREAM " + d.Name.SQL()
}

func (s *Storing) SQL() string {
	return "STORING (" + sqlJoin(s.Columns, ", ") + ")"
}

func (i *InterleaveIn) SQL() string {
	return ", INTERLEAVE IN " + i.TableName.SQL()
}

func (a *AlterIndex) SQL() string {
	return "ALTER INDEX " + a.Name.SQL() + " " + a.IndexAlteration.SQL()
}

func (a *AddStoredColumn) SQL() string {
	return "ADD STORED COLUMN " + a.Name.SQL()
}

func (a *DropStoredColumn) SQL() string {
	return "DROP STORED COLUMN " + a.Name.SQL()
}

func (d *DropIndex) SQL() string {
	return "DROP INDEX " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

func (d *DropVectorIndex) SQL() string {
	return "DROP VECTOR INDEX " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

func (d *DropSequence) SQL() string {
	return "DROP SEQUENCE " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

func (c *CreateRole) SQL() string {
	return "CREATE ROLE " + c.Name.SQL()
}

func (d *DropRole) SQL() string {
	return "DROP ROLE " + d.Name.SQL()
}

func (g *Grant) SQL() string {
	return "GRANT " + g.Privilege.SQL() + " TO ROLE " + sqlJoin(g.Roles, ", ")
}

func (r *Revoke) SQL() string {
	return "REVOKE " + r.Privilege.SQL() + " FROM ROLE " + sqlJoin(r.Roles, ", ")
}

func (p *PrivilegeOnTable) SQL() string {
	return sqlJoin(p.Privileges, ", ") + " ON TABLE " + sqlJoin(p.Names, ", ")
}

func (s *SelectPrivilege) SQL() string {
	return "SELECT" +
		strOpt(len(s.Columns) > 0, "("+sqlJoin(s.Columns, ", ")+")")
}

func (i *InsertPrivilege) SQL() string {
	return "INSERT" +
		strOpt(len(i.Columns) > 0, "("+sqlJoin(i.Columns, ", ")+")")
}

func (u *UpdatePrivilege) SQL() string {
	return "UPDATE" +
		strOpt(len(u.Columns) > 0, "("+sqlJoin(u.Columns, ", ")+")")
}

func (d *DeletePrivilege) SQL() string {
	return "DELETE"
}

func (p *SelectPrivilegeOnChangeStream) SQL() string {
	return "SELECT ON CHANGE STREAM " + sqlJoin(p.Names, ", ")
}

func (s *SelectPrivilegeOnView) SQL() string {
	return "SELECT ON VIEW " + sqlJoin(s.Names, ", ")
}

func (e *ExecutePrivilegeOnTableFunction) SQL() string {
	return "EXECUTE ON TABLE FUNCTION " + sqlJoin(e.Names, ", ")
}

func (r *RolePrivilege) SQL() string {
	return "ROLE " + sqlJoin(r.Names, ", ")
}

func (s *AlterStatistics) SQL() string {
	return "ALTER STATISTICS " + s.Name.SQL() + " SET " + s.Options.SQL()
}
func (a *Analyze) SQL() string { return "ANALYZE" }

func (c *CreateModelColumn) SQL() string {
	return c.Name.SQL() + " " + c.DataType.SQL() + sqlOpt(" ", c.Options, "")
}

func (c *CreateModelInputOutput) SQL() string {
	return "INPUT (" + sqlJoin(c.InputColumns, ", ") + ") OUTPUT (" + sqlJoin(c.OutputColumns, ", ") + ")"
}

func (c *CreateModel) SQL() string {
	return "CREATE " + strOpt(c.OrReplace, "OR REPLACE ") +
		"MODEL " +
		c.Name.SQL() +
		strOpt(c.IfNotExists, " IF NOT EXISTS") +
		sqlOpt(" ", c.InputOutput, "") +
		" REMOTE" +
		sqlOpt(" ", c.Options, "")
}

func (a *AlterModel) SQL() string {
	return "ALTER MODEL " +
		strOpt(a.IfExists, "IF EXISTS ") +
		a.Name.SQL() +
		" SET " + a.Options.SQL()
}

func (d *DropModel) SQL() string {
	return "DROP MODEL " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

// ================================================================================
//
// Types for Schema
//
// ================================================================================

func (s *ScalarSchemaType) SQL() string {
	return string(s.Name)
}

func (s *SizedSchemaType) SQL() string {
	return string(s.Name) +
		"(" + strIfElse(s.Max, "MAX", sqlOpt("", s.Size, "")) + ")"
}

func (a *ArraySchemaType) SQL() string {
	return "ARRAY<" + a.Item.SQL() + ">" + strOpt(len(a.NamedArgs) > 0, "("+sqlJoin(a.NamedArgs, ", ")+")")
}

// ================================================================================
//
// Search Index DDL
//
// ================================================================================

func (c *CreateSearchIndex) SQL() string {
	return "CREATE SEARCH INDEX " + c.Name.SQL() + " ON " + c.TableName.SQL() +
		"(" + sqlJoin(c.TokenListPart, ", ") + ")" +
		sqlOpt(" ", c.Storing, "") +
		strOpt(len(c.PartitionColumns) > 0, " PARTITION BY "+sqlJoin(c.PartitionColumns, ", ")) +
		sqlOpt(" ", c.OrderBy, "") +
		sqlOpt(" ", c.Where, "") +
		sqlOpt("", c.Interleave, "") +
		sqlOpt(" ", c.Options, "")
}

func (d *DropSearchIndex) SQL() string {
	return "DROP SEARCH INDEX " + strOpt(d.IfExists, "IF EXISTS ") + d.Name.SQL()
}

func (a *AlterSearchIndex) SQL() string {
	return "ALTER SEARCH INDEX " + a.Name.SQL() + " " + a.IndexAlteration.SQL()
}

// ================================================================================
//
// DML
//
// ================================================================================

func (w *WithAction) SQL() string {
	return "WITH ACTION" + sqlOpt(" ", w.Alias, "")
}

func (t *ThenReturn) SQL() string {
	return "THEN RETURN " + sqlOpt("", t.WithAction, " ") + sqlJoin(t.Items, ", ")
}

func (i *Insert) SQL() string {
	return "INSERT " +
		strOpt(i.InsertOrType != "", "OR "+string(i.InsertOrType)+" ") +
		"INTO " + i.TableName.SQL() + " (" +
		sqlJoin(i.Columns, ", ") +
		") " +
		i.Input.SQL() +
		sqlOpt(" ", i.ThenReturn, "")
}

func (v *ValuesInput) SQL() string {
	return "VALUES " + sqlJoin(v.Rows, ", ")
}

func (v *ValuesRow) SQL() string {
	return "(" + sqlJoin(v.Exprs, ", ") + ")"
}

func (d *DefaultExpr) SQL() string {
	if d.Default {
		return "DEFAULT"
	}
	return d.Expr.SQL()
}

func (s *SubQueryInput) SQL() string {
	return s.Query.SQL()
}

func (d *Delete) SQL() string {
	return "DELETE FROM " +
		d.TableName.SQL() + " " +
		sqlOpt("", d.As, " ") +
		d.Where.SQL() +
		sqlOpt(" ", d.ThenReturn, "")
}

func (u *Update) SQL() string {
	return "UPDATE " + u.TableName.SQL() + " " +
		sqlOpt("", u.As, " ") +
		"SET " +
		sqlJoin(u.Updates, ", ") +
		" " + u.Where.SQL() +
		sqlOpt(" ", u.ThenReturn, "")
}

func (u *UpdateItem) SQL() string {
	return sqlJoin(u.Path, ".") + " = " + u.DefaultExpr.SQL()
}
