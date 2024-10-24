package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
	"strconv"
	"strings"
)

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
	case *CallExpr, *CountStarExpr, *CastExpr, *ExtractExpr, *CaseExpr, *ParenExpr, *ScalarSubQuery, *ArraySubQuery, *ExistsSubQuery, *Param, *Ident, *Path, *ArrayLiteral, *StructLiteral, *NullLiteral, *BoolLiteral, *IntLiteral, *FloatLiteral, *StringLiteral, *BytesLiteral, *DateLiteral, *TimestampLiteral, *NumericLiteral:
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

func (q *QueryStatement) SQL() string {
	return sqlOpt("", q.Hint, " ") +
		sqlOpt("", q.With, " ") +
		q.Query.SQL()
}

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
		strOpt(s.Distinct, "DISTINCT ") +
		sqlOpt("", s.As, " ") +
		sqlJoin(s.Results, ", ") +
		sqlOpt(" ", s.From, "") +
		sqlOpt(" ", s.Where, "") +
		sqlOpt(" ", s.GroupBy, "") +
		sqlOpt(" ", s.Having, "") +
		sqlOpt(" ", s.OrderBy, "") +
		sqlOpt(" ", s.Limit, "")
}

func (a *AsStruct) SQL() string { return "AS STRUCT" }

func (a *AsValue) SQL() string { return "AS VALUE" }

func (a *AsTypeName) SQL() string { return "AS " + a.TypeName.SQL() }

func (c *CompoundQuery) SQL() string {
	return sqlJoin(c.Queries,
		" "+string(c.Op)+strIfElse(c.Distinct, " DISTINCT", " ALL")+" ") +
		sqlOpt(" ", c.OrderBy, "") +
		sqlOpt(" ", c.Limit, "")
}

func (s *SubQuery) SQL() string {
	return "(" + s.Query.SQL() + ")" +
		sqlOpt(" ", s.OrderBy, "") +
		sqlOpt(" ", s.Limit, "")
}

func (s *Star) SQL() string {
	return "*"
}

func (s *DotStar) SQL() string {
	return s.Expr.SQL() + ".*"
}

func (a *Alias) SQL() string {
	return a.Expr.SQL() + " " + a.As.SQL()
}

func (a *AsAlias) SQL() string {
	return "AS " + a.Alias.SQL()
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
	return paren(p, i.Expr) + "[" +
		strIfElse(i.Ordinal, "ORDINAL", "OFFSET") +
		"(" + i.Index.SQL() + ")]"
}

func (c *CallExpr) SQL() string {
	return c.Func.SQL() + "(" + strOpt(c.Distinct, "DISTINCT ") +
		sqlJoin(c.Args, ", ") +
		strOpt(len(c.Args) > 0 && len(c.NamedArgs) > 0, ", ") +
		sqlJoin(c.NamedArgs, ", ") +
		sqlOpt(" ", c.NullHandling, "") +
		sqlOpt(" ", c.Having, "") +
		")"
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

func (s *StructLiteral) SQL() string {
	// Both of `STRUCT()` and `STRUCT<>()` are preserved as is.
	return "STRUCT" +
		strOpt(s.Fields != nil, "<"+sqlJoin(s.Fields, ", ")+">") +
		"(" + sqlJoin(s.Values, ", ") + ")"
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

func (c *CreateTable) SQL() string {
	return "CREATE TABLE " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " (" +
		sqlJoin(c.Columns, ", ") + strOpt(len(c.Columns) > 0 && (len(c.TableConstraints) > 0 || len(c.Synonyms) > 0), ", ") +
		sqlJoin(c.TableConstraints, ", ") + strOpt(len(c.TableConstraints) > 0 && len(c.Synonyms) > 0, ", ") +
		sqlJoin(c.Synonyms, ", ") +
		") PRIMARY KEY (" + sqlJoin(c.PrimaryKeys, ", ") + ")" +
		sqlOpt("", c.Cluster, "") +
		sqlOpt("", c.RowDeletionPolicy, "")
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
	return "AS (" + g.Expr.SQL() + ") STORED"
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

func (c *CreateIndex) SQL() string {
	return "CREATE " +
		strOpt(c.Unique, "UNIQUE ") +
		strOpt(c.NullFiltered, "NULL_FILTERED ") +
		"INDEX " +
		strOpt(c.IfNotExists, "IF NOT EXISTS ") +
		c.Name.SQL() + " ON " + c.TableName.SQL() + " (" +
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

func (s *SelectPrivilegeOnView) SQL() string {
	return "SELECT ON VIEW " + sqlJoin(s.Names, ", ")
}

func (e *ExecutePrivilegeOnTableFunction) SQL() string {
	return "EXECUTE ON TABLE FUNCTION " + sqlJoin(e.Names, ", ")
}

func (r *RolePrivilege) SQL() string {
	return "ROLE " + sqlJoin(r.Names, ", ")
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
	return "ARRAY<" + a.Item.SQL() + ">"
}

// ================================================================================
//
// DML
//
// ================================================================================

func (i *Insert) SQL() string {
	return "INSERT " + strOpt(i.InsertOrType != "", "OR "+string(i.InsertOrType)+" ") +
		"INTO " + i.TableName.SQL() + " (" +
		sqlJoin(i.Columns, ", ") + ") " + i.Input.SQL()
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
	return "DELETE FROM " + d.TableName.SQL() +
		sqlOpt(" ", d.As, "") +
		" " + d.Where.SQL()
}

func (u *Update) SQL() string {
	return "UPDATE " + u.TableName.SQL() +
		sqlOpt(" ", u.As, "") +
		" SET " + sqlJoin(u.Updates, ", ") +
		" " + u.Where.SQL()
}

func (u *UpdateItem) SQL() string {
	return sqlJoin(u.Path, ".") + " = " + u.DefaultExpr.SQL()
}
