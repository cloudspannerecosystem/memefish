package ast

import (
	"strconv"
	"strings"

	"github.com/cloudspannerecosystem/memefish/token"
)

// indentLevel is the whitespace indentation level.
// Currently, memefish does not perform pretty printing in general, it is only used by CreateTable.
// Note: The two-space indentation is the same as GetDatabaseDdl in an actual Spanner instance.
const indentLevel = 2

// indent is a whitespace indentation of indentLevel width.
// You should use this rather than using indentLevel directly.
var indent = strings.Repeat(" ", indentLevel)

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
	var zero T
	if node == zero {
		return ""
	}
	return left + node.SQL() + right
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
	var b strings.Builder
	for i, r := range elems {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(r.SQL())
	}
	return b.String()
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
// Bad Node
//
// ================================================================================

func (b *BadNode) SQL() string {
	var sql string
	for _, tok := range b.Tokens {
		if sql != "" && len(tok.Space) > 0 {
			sql += " "
		}
		sql += tok.Raw
	}
	return sql
}

func (b *BadStatement) SQL() string { return sqlOpt("", b.Hint, " ") + b.BadNode.SQL() }
func (b *BadQueryExpr) SQL() string { return b.BadNode.SQL() }
func (b *BadExpr) SQL() string      { return b.BadNode.SQL() }
func (b *BadType) SQL() string      { return b.BadNode.SQL() }
func (b *BadDDL) SQL() string       { return b.BadNode.SQL() }
func (b *BadDML) SQL() string       { return sqlOpt("", b.Hint, " ") + b.BadNode.SQL() }

// ================================================================================
//
// SELECT
//
// ================================================================================

func (q *QueryStatement) SQL() string {
	return sqlOpt("", q.Hint, " ") + q.Query.SQL()
}

func (q *Query) SQL() string {
	return sqlOpt("", q.With, " ") +
		q.Query.SQL() +
		sqlOpt(" ", q.OrderBy, "") +
		sqlOpt(" ", q.Limit, "") +
		sqlOpt(" ", q.ForUpdate, "") +
		strOpt(len(q.PipeOperators) > 0, " ") +
		sqlJoin(q.PipeOperators, " ")
}

func (f *ForUpdate) SQL() string { return "FOR UPDATE" }

func (h *Hint) SQL() string {
	return "@{" + sqlJoin(h.Records, ", ") + "}"
}

func (h *HintRecord) SQL() string {
	return h.Key.SQL() + "=" + h.Value.SQL()
}

func (w *With) SQL() string {
	return "WITH " + sqlJoin(w.CTEs, ", ")
}

func (c *CTE) SQL() string {
	return c.Name.SQL() + " AS (" + c.QueryExpr.SQL() + ")"
}

func (s *Select) SQL() string {
	return "SELECT " +
		strOpt(s.AllOrDistinct != "", string(s.AllOrDistinct)+" ") +
		sqlOpt("", s.As, " ") +
		sqlJoin(s.Results, ", ") +
		sqlOpt(" ", s.From, "") +
		sqlOpt(" ", s.Where, "") +
		sqlOpt(" ", s.GroupBy, "") +
		sqlOpt(" ", s.Having, "")
}

func (a *AsStruct) SQL() string { return "AS STRUCT" }

func (a *AsValue) SQL() string { return "AS VALUE" }

func (a *AsTypeName) SQL() string { return "AS " + a.TypeName.SQL() }

func (f *FromQuery) SQL() string {
	return f.From.SQL()
}

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

func (f *From) SQL() string {
	return "FROM " + f.Source.SQL()
}

func (w *Where) SQL() string {
	return "WHERE " + w.Expr.SQL()
}

func (g *GroupBy) SQL() string {
	return "GROUP " + sqlOpt("", g.Hint, " ") + "BY " + sqlJoin(g.Exprs, ", ")
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

func (s *SubQueryTableExpr) SQL() string {
	return "(" + s.Query.SQL() + ")" +
		sqlOpt(" ", s.As, "") +
		sqlOpt(" ", s.Sample, "")
}

func (p *ParenTableExpr) SQL() string {
	return "(" + p.Source.SQL() + ")" +
		sqlOpt(" ", p.Sample, "")
}

func (j *Join) SQL() string {
	return j.Left.SQL() +
		strOpt(j.Op != CommaJoin, " ") +
		string(j.Op) + " " +
		sqlOpt("", j.Hint, " ") +
		j.Right.SQL() +
		sqlOpt(" ", j.Cond, "")
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
		sqlOpt(" ", c.OrderBy, "") +
		sqlOpt(" ", c.Limit, "") +
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

func (n *IntervalLiteralSingle) SQL() string {
	return "INTERVAL " + n.Value.SQL() + " " + string(n.DateTimePart)
}

func (n *IntervalLiteralRange) SQL() string {
	return "INTERVAL " + n.Value.SQL() + " " + string(n.StartingDateTimePart) + " TO " + string(n.EndingDateTimePart)
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

func (c *CreateLocalityGroup) SQL() string {
	return "CREATE LOCALITY GROUP " + c.Name.SQL() + sqlOpt(" ", c.Options, "")
}

func (a *AlterLocalityGroup) SQL() string {
	return "ALTER LOCALITY GROUP " + a.Name.SQL() + " SET " + a.Options.SQL()
}

func (d *DropLocalityGroup) SQL() string {
	return "DROP LOCALITY GROUP " + d.Name.SQL()
}

func (s *CreateSchema) SQL() string {
	return "CREATE" + strOpt(s.OrReplace, " OR REPLACE") + " SCHEMA " + strOpt(s.IfNotExists, "IF NOT EXISTS ") + s.Name.SQL()
}

func (s *DropSchema) SQL() string {
	return "DROP SCHEMA " + strOpt(s.IfExists, "IF EXISTS ") + s.Name.SQL()
}

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

func (c *CreateTable) SQL() string {
	return "CREATE TABLE " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " (\n" +
		indent + sqlJoin(c.Columns, ",\n"+indent) + strOpt(len(c.Columns) > 0 && (len(c.TableConstraints) > 0 || len(c.Synonyms) > 0), ",\n") +
		strOpt(len(c.TableConstraints) > 0, indent) + sqlJoin(c.TableConstraints, ",\n"+indent) + strOpt(len(c.TableConstraints) > 0 && len(c.Synonyms) > 0, ",\n") +
		strOpt(len(c.Synonyms) > 0, indent) + sqlJoin(c.Synonyms, ",\n") +
		"\n)" +
		strOpt(len(c.PrimaryKeys) > 0, " PRIMARY KEY ("+sqlJoin(c.PrimaryKeys, ", ")+")") +
		sqlOpt("", c.Cluster, "") +
		sqlOpt("", c.RowDeletionPolicy, "") +
		sqlOpt(", ", c.Options, "")
}

func (s *Synonym) SQL() string { return "SYNONYM (" + s.Name.SQL() + ")" }

func (c *CreateSequence) SQL() string {
	return "CREATE SEQUENCE " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() +
		strOpt(len(c.Params) > 0, " "+sqlJoin(c.Params, "")) +
		sqlOpt(" ", c.Options, "")
}

func (s *SkipRange) SQL() string {
	return "SKIP RANGE " + s.Min.SQL() + ", " + s.Max.SQL()
}

func (s *StartCounterWith) SQL() string {
	return "START COUNTER WITH " + s.Counter.SQL()
}

func (s *BitReversedPositive) SQL() string {
	return "BIT_REVERSED_POSITIVE"
}

func (c *AlterSequence) SQL() string {
	return "ALTER SEQUENCE " + c.Name.SQL() +
		sqlOpt(" SET ", c.Options, "") +
		sqlOpt(" ", c.RestartCounterWith, "") +
		sqlOpt(" ", c.SkipRange, "") +
		sqlOpt(" ", c.NoSkipRange, "")
}

func (c *CreateView) SQL() string {
	return "CREATE" + strOpt(c.OrReplace, " OR REPLACE") + " VIEW " + c.Name.SQL() +
		" SQL SECURITY " + string(c.SecurityType) + " AS " + c.Query.SQL()
}

func (d *DropView) SQL() string { return "DROP VIEW " + d.Name.SQL() }

func (c *ColumnDef) SQL() string {
	return c.Name.SQL() + " " + c.Type.SQL() +
		strOpt(c.NotNull, " NOT NULL") +
		sqlOpt(" ", c.DefaultSemantics, "") +
		strOpt(!c.Hidden.Invalid(), " HIDDEN") +
		strOpt(c.PrimaryKey, " PRIMARY KEY") +
		sqlOpt(" ", c.Options, "")
}

func (c *TableConstraint) SQL() string {
	return sqlOpt("CONSTRAINT ", c.Name, " ") + c.Constraint.SQL()
}

func (f *ForeignKey) SQL() string {
	return "FOREIGN KEY (" + sqlJoin(f.Columns, ", ") + ") " +
		"REFERENCES " + f.ReferenceTable.SQL() + " (" +
		sqlJoin(f.ReferenceColumns, ", ") + ")" +
		strOpt(f.OnDelete != "", " "+string(f.OnDelete)) +
		strOpt(f.Enforcement != "", " "+string(f.Enforcement))
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

func (i *IdentityColumn) SQL() string {
	return "GENERATED BY DEFAULT AS IDENTITY" + strOpt(len(i.Params) > 0, " ("+sqlJoin(i.Params, " ")+")")
}

func (a *AutoIncrement) SQL() string {
	return "AUTO_INCREMENT"
}

func (i *IndexKey) SQL() string {
	return i.Name.SQL() + strOpt(i.Dir != "", " "+string(i.Dir))
}

func (c *Cluster) SQL() string {
	return ",\n" + indent + "INTERLEAVE IN " + strOpt(c.Enforced, "PARENT ") + c.TableName.SQL() +
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

func (s *SetInterleaveIn) SQL() string {
	return "SET INTERLEAVE IN " + strOpt(s.Enforced, "PARENT ") + s.TableName.SQL() +
		strOpt(s.OnDelete != "", " "+string(s.OnDelete))
}

func (a *AlterTableSetOptions) SQL() string { return "SET " + a.Options.SQL() }

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

func (a *AlterColumnAlterIdentity) SQL() string { return "ALTER IDENTITY " + a.Alteration.SQL() }

func (r *RestartCounterWith) SQL() string { return "RESTART COUNTER WITH " + r.Counter.SQL() }

func (s *SetSkipRange) SQL() string { return "SET " + s.SkipRange.SQL() }

func (s *NoSkipRange) SQL() string { return "NO SKIP RANGE" }

func (s *SetNoSkipRange) SQL() string { return "SET NO SKIP RANGE" }

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
		sqlOpt("", c.InterleaveIn, "") +
		sqlOpt(" ", c.Options, "")
}

func (c *CreateVectorIndex) SQL() string {
	return "CREATE VECTOR INDEX " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " ON " + c.TableName.SQL() + " (" + c.ColumnName.SQL() + ") " +
		sqlOpt("", c.Storing, " ") +
		sqlOpt("", c.Where, " ") +
		c.Options.SQL()
}

func (a *AlterVectorIndex) SQL() string {
	return "ALTER VECTOR INDEX " + a.Name.SQL() + " " + a.Alteration.SQL()
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
	return "FOR " + sqlJoin(c.Tables, ", ")
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
	return c.TableName.SQL() + strOpt(!c.Rparen.Invalid(), "("+sqlJoin(c.Columns, ", ")+")")
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

func (p *FunctionParam) SQL() string {
	return p.Name.SQL() + " " + p.Type.SQL() + sqlOpt(" DEFAULT ", p.DefaultExpr, "")
}

func (c *CreateFunction) SQL() string {
	return "CREATE" +
		strOpt(c.OrReplace, " OR REPLACE") +
		" FUNCTION " +
		c.Name.SQL() +
		" (" + sqlJoin(c.Params, ", ") + ")" +
		sqlOpt(" RETURNS ", c.ReturnType, "") +
		strOpt(c.Determinism != "", " "+string(c.Determinism)) +
		strOpt(c.Language != "", " LANGUAGE "+c.Language) +
		strOpt(c.Remote, " REMOTE") +
		strOpt(c.SqlSecurity != "", " SQL SECURITY "+string(c.SqlSecurity)) +
		sqlOpt(" ", c.Options, "") +
		sqlOpt(" AS (", c.Definition, ")")
}

func (d *DropFunction) SQL() string {
	return "DROP" +
		" FUNCTION " +
		strOpt(d.IfExists, "IF EXISTS ") +
		d.Name.SQL()
}

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
// GQL schema statements
//
// ================================================================================

func (c *CreatePropertyGraph) SQL() string {
	return "CREATE " +
		strOpt(c.OrReplace, "OR REPLACE ") +
		"PROPERTY GRAPH " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " " + c.Content.SQL()
}

func (p *PropertyGraphContent) SQL() string {
	return p.NodeTables.SQL() + sqlOpt(" ", p.EdgeTables, "")
}

func (p *PropertyGraphNodeTables) SQL() string {
	return "NODE TABLES " + p.Tables.SQL()
}

func (p *PropertyGraphEdgeTables) SQL() string {
	return "EDGE TABLES " + p.Tables.SQL()
}

func (p *PropertyGraphElementList) SQL() string {
	return "(" + sqlJoin(p.Elements, ", ") + ")"
}

func (p *PropertyGraphElement) SQL() string {
	return p.Name.SQL() +
		sqlOpt(" AS ", p.Alias, "") +
		sqlOpt(" ", p.Keys, "") +
		sqlOpt(" ", p.Properties, "") +
		sqlOpt(" ", p.DynamicLabel, "") +
		sqlOpt(" ", p.DynamicProperties, "")
}

func (p *PropertyGraphSingleProperties) SQL() string { return p.Properties.SQL() }

func (p *PropertyGraphLabelAndPropertiesList) SQL() string {
	return sqlJoin(p.LabelAndProperties, " ")
}

func (p *PropertyGraphLabelAndProperties) SQL() string {
	return p.Label.SQL() + sqlOpt(" ", p.Properties, "")
}

func (p *PropertyGraphElementLabelLabelName) SQL() string { return "LABEL " + p.Name.SQL() }

func (p *PropertyGraphElementLabelDefaultLabel) SQL() string { return "DEFAULT LABEL" }

func (p *PropertyGraphNodeElementKey) SQL() string { return p.Key.SQL() }

func (p *PropertyGraphEdgeElementKeys) SQL() string {
	return sqlOpt("", p.Element, " ") +
		p.Source.SQL() + " " + p.Destination.SQL()
}

func (p *PropertyGraphElementKey) SQL() string { return "KEY " + p.Keys.SQL() }

func (p *PropertyGraphSourceKey) SQL() string {
	return "SOURCE KEY " + p.Keys.SQL() +
		" REFERENCES " + p.ElementReference.SQL() +
		sqlOpt(" ", p.ReferenceColumns, "")
}

func (p *PropertyGraphDestinationKey) SQL() string {
	return "DESTINATION KEY " + p.Keys.SQL() +
		" REFERENCES " + p.ElementReference.SQL() +
		sqlOpt(" ", p.ReferenceColumns, "")
}

func (p *PropertyGraphColumnNameList) SQL() string {
	return "(" + sqlJoin(p.ColumnNameList, ", ") + ")"
}

func (p *PropertyGraphNoProperties) SQL() string {
	return "NO PROPERTIES"
}

func (p *PropertyGraphPropertiesAre) SQL() string {
	return "PROPERTIES ARE ALL COLUMNS" + sqlOpt(" EXCEPT ", p.ExceptColumns, "")
}

func (p *PropertyGraphDerivedPropertyList) SQL() string {
	return "PROPERTIES (" + sqlJoin(p.DerivedProperties, ", ") + ")"
}

func (p *PropertyGraphDerivedProperty) SQL() string {
	return p.Expr.SQL() + sqlOpt(" AS ", p.Alias, "")
}

func (p *PropertyGraphDynamicLabel) SQL() string {
	return "DYNAMIC LABEL (" + p.ColumnName.SQL() + ")"
}

func (p *PropertyGraphDynamicProperties) SQL() string {
	return "DYNAMIC PROPERTIES (" + p.ColumnName.SQL() + ")"
}

func (g *DropPropertyGraph) SQL() string {
	return "DROP PROPERTY GRAPH " +
		strOpt(g.IfExists, "IF EXISTS ") +
		g.Name.SQL()
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
	return sqlOpt("", i.Hint, " ") +
		"INSERT " +
		strOpt(i.InsertOrType != "", "OR "+string(i.InsertOrType)+" ") +
		"INTO " + i.TableName.SQL() +
		sqlOpt("", i.TableHint, "") + " (" +
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
	return sqlOpt("", d.Hint, " ") +
		"DELETE FROM " +
		d.TableName.SQL() +
		sqlOpt("", d.TableHint, "") + " " +
		sqlOpt("", d.As, " ") +
		d.Where.SQL() +
		sqlOpt(" ", d.ThenReturn, "")
}

func (u *Update) SQL() string {
	return sqlOpt("", u.Hint, " ") +
		"UPDATE " + u.TableName.SQL() +
		sqlOpt("", u.TableHint, "") + " " +
		sqlOpt("", u.As, " ") +
		"SET " +
		sqlJoin(u.Updates, ", ") +
		" " + u.Where.SQL() +
		sqlOpt(" ", u.ThenReturn, "")
}

func (u *UpdateItem) SQL() string {
	return sqlJoin(u.Path, ".") + " = " + u.DefaultExpr.SQL()
}

// ================================================================================
//
// Procedural language
//
// ================================================================================

func (c *Call) SQL() string {
	return "CALL " + c.Name.SQL() + "(" + sqlJoin(c.Args, ", ") + ")"
}
