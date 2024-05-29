package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
)

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
	var sql string
	if q.Hint != nil {
		sql += q.Hint.SQL() + " "
	}
	if q.With != nil {
		sql += q.With.SQL() + " "
	}
	sql += q.Query.SQL()
	return sql
}

func (h *Hint) SQL() string {
	sql := "@{" + h.Records[0].SQL()
	for _, r := range h.Records[1:] {
		sql += ", " + r.SQL()
	}
	sql += "}"
	return sql
}

func (h *HintRecord) SQL() string {
	return h.Key.SQL() + "=" + h.Value.SQL()
}

func (w *With) SQL() string {
	sql := "WITH " + w.CTEs[0].SQL()
	for _, c := range w.CTEs[1:] {
		sql += ", " + c.SQL()
	}
	return sql
}

func (c *CTE) SQL() string {
	return c.Name.SQL() + " AS (" + c.QueryExpr.SQL() + ")"
}

func (s *Select) SQL() string {
	sql := "SELECT "
	if s.Distinct {
		sql += "DISTINCT "
	}
	if s.AsStruct {
		sql += "AS STRUCT "
	}
	sql += s.Results[0].SQL()
	for _, r := range s.Results[1:] {
		sql += ", " + r.SQL()
	}
	if s.From != nil {
		sql += " " + s.From.SQL()
	}
	if s.Where != nil {
		sql += " " + s.Where.SQL()
	}
	if s.GroupBy != nil {
		sql += " " + s.GroupBy.SQL()
	}
	if s.Having != nil {
		sql += " " + s.Having.SQL()
	}
	if s.OrderBy != nil {
		sql += " " + s.OrderBy.SQL()
	}
	if s.Limit != nil {
		sql += " " + s.Limit.SQL()
	}
	return sql
}

func (c *CompoundQuery) SQL() string {
	op := string(c.Op)
	if c.Distinct {
		op += " DISTINCT"
	} else {
		op += " ALL"
	}

	sql := c.Queries[0].SQL()
	for _, q := range c.Queries[1:] {
		sql += " " + op + " " + q.SQL()
	}
	if c.OrderBy != nil {
		sql += " " + c.OrderBy.SQL()
	}
	if c.Limit != nil {
		sql += " " + c.Limit.SQL()
	}
	return sql
}

func (s *SubQuery) SQL() string {
	sql := "(" + s.Query.SQL() + ")"
	if s.OrderBy != nil {
		sql += " " + s.OrderBy.SQL()
	}
	if s.Limit != nil {
		sql += " " + s.Limit.SQL()
	}
	return sql
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
	sql := "GROUP BY " + g.Exprs[0].SQL()
	for _, e := range g.Exprs[1:] {
		sql += ", " + e.SQL()
	}
	return sql
}

func (h *Having) SQL() string {
	return "HAVING " + h.Expr.SQL()
}

func (o *OrderBy) SQL() string {
	sql := "ORDER BY " + o.Items[0].SQL()
	for _, item := range o.Items[1:] {
		sql += ", " + item.SQL()
	}
	return sql
}

func (o *OrderByItem) SQL() string {
	sql := o.Expr.SQL()
	if o.Collate != nil {
		sql += " " + o.Collate.SQL()
	}
	if o.Dir != "" {
		sql += " " + string(o.Dir)
	}
	return sql
}

func (c *Collate) SQL() string {
	return "COLLATE " + c.Value.SQL()
}

func (l *Limit) SQL() string {
	sql := "LIMIT " + l.Count.SQL()
	if l.Offset != nil {
		sql += " " + l.Offset.SQL()
	}
	return sql
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
	var sql string
	if u.Implicit {
		sql += u.Expr.SQL()
	} else {
		sql += "UNNEST(" + u.Expr.SQL() + ")"
	}
	if u.Hint != nil {
		sql += " " + u.Hint.SQL()
	}
	if u.As != nil {
		sql += " " + u.As.SQL()
	}
	if u.WithOffset != nil {
		sql += " " + u.WithOffset.SQL()
	}
	if u.Sample != nil {
		sql += " " + u.Sample.SQL()
	}
	return sql
}

func (w *WithOffset) SQL() string {
	sql := "WITH OFFSET"
	if w.As != nil {
		sql += " " + w.As.SQL()
	}
	return sql
}

func (t *TableName) SQL() string {
	sql := t.Table.SQL()
	if t.Hint != nil {
		sql += " " + t.Hint.SQL()
	}
	if t.As != nil {
		sql += " " + t.As.SQL()
	}
	if t.Sample != nil {
		sql += " " + t.Sample.SQL()
	}
	return sql
}

func (s *SubQueryTableExpr) SQL() string {
	sql := "(" + s.Query.SQL() + ")"
	if s.As != nil {
		sql += " " + s.As.SQL()
	}
	if s.Sample != nil {
		sql += " " + s.Sample.SQL()
	}
	return sql
}

func (p *ParenTableExpr) SQL() string {
	sql := "(" + p.Source.SQL() + ")"
	if p.Sample != nil {
		sql += " " + p.Sample.SQL()
	}
	return sql
}

func (j *Join) SQL() string {
	sql := j.Left.SQL()
	if j.Op != CommaJoin {
		sql += " "
	}
	sql += string(j.Op) + " "
	if j.Hint != nil {
		sql += j.Hint.SQL() + " "
	}
	sql += j.Right.SQL()
	if j.Cond != nil {
		sql += " " + j.Cond.SQL()
	}
	return sql
}

func (o *On) SQL() string {
	return "ON " + o.Expr.SQL()
}

func (u *Using) SQL() string {
	sql := "USING (" + u.Idents[0].SQL()
	for _, id := range u.Idents[1:] {
		sql += ", " + id.SQL()
	}
	sql += ")"
	return sql
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
	sql := paren(p, b.Left)
	sql += " " + string(b.Op) + " "
	sql += paren(p, b.Right)
	return sql
}

func (u *UnaryExpr) SQL() string {
	p := exprPrec(u)
	if u.Op == OpNot {
		return "NOT " + paren(p, u.Expr)
	}
	return string(u.Op) + paren(p, u.Expr)
}

func (i *InExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	if i.Not {
		sql += " NOT"
	}
	sql += " IN "
	sql += i.Right.SQL()
	return sql
}

func (u *UnnestInCondition) SQL() string {
	return "UNNEST(" + u.Expr.SQL() + ")"
}

func (s *SubQueryInCondition) SQL() string {
	return "(" + s.Query.SQL() + ")"
}

func (v *ValuesInCondition) SQL() string {
	sql := "(" + v.Exprs[0].SQL()
	for _, e := range v.Exprs[1:] {
		sql += ", " + e.SQL()
	}
	sql += ")"
	return sql
}

func (i *IsNullExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	sql += " IS "
	if i.Not {
		sql += "NOT "
	}
	sql += "NULL"
	return sql
}

func (i *IsBoolExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Left)
	sql += " IS "
	if i.Not {
		sql += "NOT "
	}
	if i.Right {
		sql += "TRUE"
	} else {
		sql += "FALSE"
	}
	return sql
}

func (b *BetweenExpr) SQL() string {
	p := exprPrec(b)
	sql := paren(p, b.Left)
	if b.Not {
		sql += " NOT"
	}
	return sql + " BETWEEN " + paren(p, b.RightStart) + " AND " + paren(p, b.RightEnd)
}

func (s *SelectorExpr) SQL() string {
	p := exprPrec(s)
	return paren(p, s.Expr) + "." + s.Ident.SQL()
}

func (i *IndexExpr) SQL() string {
	p := exprPrec(i)
	sql := paren(p, i.Expr) + "["
	if i.Ordinal {
		sql += "ORDINAL"
	} else {
		sql += "OFFSET"
	}
	sql += "(" + i.Index.SQL() + ")]"
	return sql
}

func (c *CallExpr) SQL() string {
	sql := c.Func.SQL() + "("
	if c.Distinct {
		sql += "DISTINCT "
	}
	for i, a := range c.Args {
		if i != 0 {
			sql += ", "
		}
		sql += a.SQL()
	}
	sql += ")"
	return sql
}

func (o *SequenceOption) SQL() string {
	return o.Name.SQL() + " = " + o.Value.SQL()
}

func (s *ExprArg) SQL() string {
	return s.Expr.SQL()
}

func (i *IntervalArg) SQL() string {
	sql := "INTERVAL " + i.Expr.SQL()
	if i.Unit != nil {
		sql += " " + i.Unit.SQL()
	}
	return sql
}

func (s *SequenceArg) SQL() string {
	return "SEQUENCE " + s.Expr.SQL()
}

func (*CountStarExpr) SQL() string {
	return "COUNT(*)"
}

func (e *ExtractExpr) SQL() string {
	sql := "EXTRACT(" + e.Part.SQL() + " FROM " + e.Expr.SQL()
	if e.AtTimeZone != nil {
		sql += " " + e.AtTimeZone.SQL()
	}
	sql += ")"
	return sql
}

func (a *AtTimeZone) SQL() string {
	return "AT TIME ZONE " + a.Expr.SQL()
}

func (c *CastExpr) SQL() string {
	return "CAST(" + c.Expr.SQL() + " AS " + c.Type.SQL() + ")"
}

func (c *CaseExpr) SQL() string {
	sql := "CASE "
	if c.Expr != nil {
		sql += c.Expr.SQL() + " "
	}
	for _, w := range c.Whens {
		sql += w.SQL() + " "
	}
	if c.Else != nil {
		sql += c.Else.SQL() + " "
	}
	sql += "END"
	return sql
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
	sql := "EXISTS"
	if e.Hint != nil {
		sql += " " + e.Hint.SQL() + " "
	}
	sql += "(" + e.Query.SQL() + ")"
	return sql
}

func (p *Param) SQL() string {
	return "@" + p.Name
}

func (i *Ident) SQL() string {
	return token.QuoteSQLIdent(i.Name)
}

func (p *Path) SQL() string {
	sql := p.Idents[0].SQL()
	for _, id := range p.Idents[1:] {
		sql += "." + id.SQL()
	}
	return sql
}

func (a *ArrayLiteral) SQL() string {
	sql := "ARRAY"
	if a.Type != nil {
		sql += "<" + a.Type.SQL() + ">"
	}
	sql += "["
	for i, v := range a.Values {
		if i != 0 {
			sql += ", "
		}
		sql += v.SQL()
	}
	sql += "]"
	return sql
}

func (s *StructLiteral) SQL() string {
	sql := "STRUCT"
	if s.Fields != nil {
		sql += "<"
		for i, f := range s.Fields {
			if i != 0 {
				sql += ", "
			}
			sql += f.SQL()
		}
		sql += ">"
	}
	sql += "("
	for i, v := range s.Values {
		if i != 0 {
			sql += ", "
		}
		sql += v.SQL()
	}
	sql += ")"
	return sql
}

func (*NullLiteral) SQL() string {
	return "NULL"
}

func (b *BoolLiteral) SQL() string {
	if b.Value {
		return "TRUE"
	} else {
		return "FALSE"
	}
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
	sql := "STRUCT<"
	for i, f := range s.Fields {
		if i != 0 {
			sql += ", "
		}
		sql += f.SQL()
	}
	sql += ">"
	return sql
}

func (f *StructField) SQL() string {
	var sql string
	if f.Ident != nil {
		sql += f.Ident.SQL() + " "
	}
	sql += f.Type.SQL()
	return sql
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

func (c *CreateDatabase) SQL() string {
	return "CREATE DATABASE " + c.Name.SQL()
}

func (c *CreateTable) SQL() string {
	sql := "CREATE TABLE "
	if c.IfNotExists {
		sql += "IF NOT EXISTS "
	}
	sql += c.Name.SQL() + " ("
	for i, c := range c.Columns {
		if i != 0 {
			sql += ", "
		}
		sql += c.SQL()
	}
	for _, c := range c.TableConstraints {
		sql += ", " + c.SQL()
	}
	sql += ") "
	sql += "PRIMARY KEY ("
	for i, k := range c.PrimaryKeys {
		if i != 0 {
			sql += ", "
		}
		sql += k.SQL()
	}
	sql += ")"
	if c.Cluster != nil {
		sql += c.Cluster.SQL()
	}
	if c.RowDeletionPolicy != nil {
		sql += c.RowDeletionPolicy.SQL()
	}
	return sql
}

func (c *CreateSequence) SQL() string {
	sql := "CREATE SEQUENCE "
	if c.IfNotExists {
		sql += "IF NOT EXISTS "
	}
	sql += c.Name.SQL() + " OPTIONS ("
	for i, o := range c.Options {
		if i > 0 {
			sql += ", "
		}
		sql += o.SQL()
	}
	sql += ")"
	return sql
}

func (c *CreateView) SQL() string {
	sql := "CREATE"
	if c.OrReplace {
		sql += " OR REPLACE"
	}
	sql += " VIEW " + c.Name.SQL() + " SQL SECURITY " + string(c.SecurityType) + " AS " + c.Query.SQL()
	return sql
}

func (c *ColumnDef) SQL() string {
	sql := c.Name.SQL() + " " + c.Type.SQL()
	if c.NotNull {
		sql += " NOT NULL"
	}
	if c.DefaultExpr != nil {
		sql += " " + c.DefaultExpr.SQL()
	}
	if c.GeneratedExpr != nil {
		sql += " " + c.GeneratedExpr.SQL()
	}
	if c.Options != nil {
		sql += " " + c.Options.SQL()
	}
	return sql
}

func (c *TableConstraint) SQL() string {
	var sql string
	if c.Name != nil {
		sql += "CONSTRAINT " + c.Name.SQL() + " "
	}
	sql += c.Constraint.SQL()
	return sql
}

func (f *ForeignKey) SQL() string {
	var sql string
	sql += "FOREIGN KEY ("
	for i, k := range f.Columns {
		if i != 0 {
			sql += ", "
		}
		sql += k.SQL()
	}
	sql += ") "
	sql += "REFERENCES " + f.ReferenceTable.SQL() + " ("
	for i, k := range f.ReferenceColumns {
		if i != 0 {
			sql += ", "
		}
		sql += k.SQL()
	}
	sql += ")"
	if f.OnDelete != "" {
		sql += " " + string(f.OnDelete)
	}
	return sql
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

func (c *ColumnDefOptions) SQL() string {
	sql := "OPTIONS(allow_commit_timestamp = "
	if c.AllowCommitTimestamp {
		sql += "true)"
	} else {
		sql += "null)"
	}
	return sql
}

func (i *IndexKey) SQL() string {
	sql := i.Name.SQL()
	if i.Dir != "" {
		sql += " " + string(i.Dir)
	}
	return sql
}

func (c *Cluster) SQL() string {
	sql := ", INTERLEAVE IN PARENT " + c.TableName.SQL()
	if c.OnDelete != "" {
		sql += " " + string(c.OnDelete)
	}
	return sql
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
	sql := "ADD COLUMN "
	if a.IfNotExists {
		sql += "IF NOT EXISTS "
	}
	return sql + a.Column.SQL()
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
	sql := "ALTER COLUMN " + a.Name.SQL() + " " + a.Type.SQL()
	if a.NotNull {
		sql += " NOT NULL"
	}
	if a.DefaultExpr != nil {
		sql += " " + a.DefaultExpr.SQL()
	}
	return sql
}

func (a *AlterColumnSet) SQL() string {
	sql := "ALTER COLUMN " + a.Name.SQL() + " SET "
	if a.Options != nil {
		return sql + a.Options.SQL()
	} else {
		return sql + a.DefaultExpr.SQL()
	}
}

func (d *DropTable) SQL() string {
	sql := "DROP TABLE "
	if d.IfExists {
		sql += "IF EXISTS "
	}
	return sql + d.Name.SQL()
}

func (c *CreateIndex) SQL() string {
	sql := "CREATE "
	if c.Unique {
		sql += "UNIQUE "
	}
	if c.NullFiltered {
		sql += "NULL_FILTERED "
	}
	sql += "INDEX "
	if c.IfNotExists {
		sql += "IF NOT EXISTS "
	}
	sql += c.Name.SQL() + " ON " + c.TableName.SQL() + " ("
	for i, k := range c.Keys {
		if i != 0 {
			sql += ", "
		}
		sql += k.SQL()
	}
	sql += ")"
	if c.Storing != nil {
		sql += " " + c.Storing.SQL()
	}
	if c.InterleaveIn != nil {
		sql += c.InterleaveIn.SQL()
	}
	return sql
}

func (c *CreateChangeStream) SQL() string {
	sql := "CREATE CHANGE STREAM " + c.Name.SQL()
	if c.For != nil {
		sql += " " + c.For.SQL()
	}
	if c.Options != nil {
		sql += c.Options.SQL()
	}
	return sql
}

func (c *ChangeStreamForAll) SQL() string {
	return "FOR ALL"
}

func (c *ChangeStreamForTables) SQL() string {
	sql := "FOR "
	for i, table := range c.Tables {
		if i > 0 {
			sql += ", "
		}
		sql += table.SQL()
	}
	return sql
}

func (c *ChangeStreamOptions) SQL() string {
	sql := " OPTIONS ("
	for i, record := range c.Records {
		if i > 0 {
			sql += ", "
		}
		sql += record.Key.SQL() + "=" + record.Value.SQL()
	}
	sql += ")"
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
	sql := c.TableName.SQL()
	if len(c.Columns) > 0 {
		sql += "("
		for i, id := range c.Columns {
			if i > 0 {
				sql += ", "
			}
			sql += id.SQL()
		}
		sql += ")"
	}
	return sql
}

func (d *DropChangeStream) SQL() string {
	return "DROP CHANGE STREAM " + d.Name.SQL()
}

func (s *Storing) SQL() string {
	sql := "STORING ("
	for i, c := range s.Columns {
		if i != 0 {
			sql += ", "
		}
		sql += c.SQL()
	}
	sql += ")"
	return sql
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
	sql := "DROP INDEX "
	if d.IfExists {
		sql += "IF EXISTS "
	}
	return sql + d.Name.SQL()
}

func (c *CreateRole) SQL() string {
	return "CREATE ROLE " + c.Name.SQL()
}

func (d *DropRole) SQL() string {
	return "DROP ROLE " + d.Name.SQL()
}

func (g *Grant) SQL() string {
	sql := "GRANT "
	sql += g.Privilege.SQL()
	sql += " TO ROLE " + g.Roles[0].SQL()
	for _, id := range g.Roles[1:] {
		sql += ", " + id.SQL()
	}
	return sql
}

func (r *Revoke) SQL() string {
	sql := "REVOKE "
	sql += r.Privilege.SQL()
	sql += " FROM ROLE " + r.Roles[0].SQL()
	for _, id := range r.Roles[1:] {
		sql += ", " + id.SQL()
	}
	return sql
}

func (p *PrivilegeOnTable) SQL() string {
	sql := p.Privileges[0].SQL()
	for _, p := range p.Privileges[1:] {
		sql += ", " + p.SQL()
	}
	sql += " ON TABLE "
	sql += p.Names[0].SQL()
	for _, id := range p.Names[1:] {
		sql += ", " + id.SQL()
	}
	return sql
}

func (s *SelectPrivilege) SQL() string {
	sql := "SELECT"
	if len(s.Columns) > 0 {
		sql += "("
		for i, c := range s.Columns {
			if i > 0 {
				sql += ", "
			}
			sql += c.SQL()
		}
		sql += ")"
	}
	return sql
}

func (i *InsertPrivilege) SQL() string {
	sql := "INSERT"
	if len(i.Columns) > 0 {
		sql += "("
		for j, c := range i.Columns {
			if j > 0 {
				sql += ", "
			}
			sql += c.SQL()
		}
		sql += ")"
	}
	return sql
}

func (u *UpdatePrivilege) SQL() string {
	sql := "UPDATE"
	if len(u.Columns) > 0 {
		sql += "("
		for i, c := range u.Columns {
			if i > 0 {
				sql += ", "
			}
			sql += c.SQL()
		}
		sql += ")"
	}
	return sql
}

func (d *DeletePrivilege) SQL() string {
	return "DELETE"
}

func (s *SelectPrivilegeOnView) SQL() string {
	sql := "SELECT ON VIEW " + s.Names[0].SQL()
	for _, v := range s.Names[1:] {
		sql += ", " + v.SQL()
	}
	return sql
}

func (e *ExecutePrivilegeOnTableFunction) SQL() string {
	sql := "EXECUTE ON TABLE FUNCTION " + e.Names[0].SQL()
	for _, f := range e.Names[1:] {
		sql += ", " + f.SQL()
	}
	return sql
}

func (r *RolePrivilege) SQL() string {
	sql := "ROLE " + r.Names[0].SQL()
	for _, id := range r.Names[1:] {
		sql += ", " + id.SQL()
	}
	return sql
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
	sql := string(s.Name) + "("
	if s.Max {
		sql += "MAX"
	} else {
		sql += s.Size.SQL()
	}
	sql += ")"
	return sql
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
	sql := "INSERT "
	if i.InsertOrType != "" {
		sql += "OR " + string(i.InsertOrType) + " "
	}
	sql += "INTO " + i.TableName.SQL() + " ("
	for i, c := range i.Columns {
		if i != 0 {
			sql += ", "
		}
		sql += c.SQL()
	}
	sql += ") " + i.Input.SQL()
	return sql
}

func (v *ValuesInput) SQL() string {
	sql := "VALUES "
	for i, r := range v.Rows {
		if i != 0 {
			sql += ", "
		}
		sql += r.SQL()
	}
	return sql
}

func (v *ValuesRow) SQL() string {
	sql := "("
	for i, v := range v.Exprs {
		if i != 0 {
			sql += ", "
		}
		sql += v.SQL()
	}
	sql += ")"
	return sql
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
	sql := "DELETE FROM " + d.TableName.SQL()
	if d.As != nil {
		sql += " " + d.As.SQL()
	}
	sql += " " + d.Where.SQL()
	return sql
}

func (u *Update) SQL() string {
	sql := "UPDATE " + u.TableName.SQL()
	if u.As != nil {
		sql += " " + u.As.SQL()
	}
	sql += " SET " + u.Updates[0].SQL()
	for _, item := range u.Updates[1:] {
		sql += ", " + item.SQL()
	}
	sql += " " + u.Where.SQL()
	return sql
}

func (u *UpdateItem) SQL() string {
	sql := u.Path[0].SQL()
	for _, id := range u.Path[1:] {
		sql += "." + id.SQL()
	}
	sql += " = " + u.Expr.SQL()
	return sql
}
