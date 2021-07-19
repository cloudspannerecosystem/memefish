package ast

import (
	"github.com/MakeNowJust/memefish/pkg/token"
)

// ================================================================================
//
// SELECT
//
// ================================================================================

func (q *QueryStatement) Pos() token.Pos {
	if q.Hint != nil {
		return q.Hint.Pos()
	}
	return q.Query.Pos()
}

func (q *QueryStatement) End() token.Pos {
	return q.Query.End()
}

func (h *Hint) Pos() token.Pos { return h.Atmark }
func (h *Hint) End() token.Pos { return h.Rbrace + 1 }

func (h *HintRecord) Pos() token.Pos { return h.Key.Pos() }
func (h *HintRecord) End() token.Pos { return h.Value.End() }

func (s *Select) Pos() token.Pos { return s.Select }

func (s *Select) End() token.Pos {
	if s.Limit != nil {
		return s.Limit.End()
	}
	if s.OrderBy != nil {
		return s.OrderBy.End()
	}
	if s.Having != nil {
		return s.Having.End()
	}
	if s.GroupBy != nil {
		return s.GroupBy.End()
	}
	if s.Where != nil {
		return s.Where.End()
	}
	if s.From != nil {
		return s.From.End()
	}
	return s.Results[len(s.Results)-1].End()
}

func (c *CompoundQuery) Pos() token.Pos {
	return c.Queries[0].Pos()
}

func (c *CompoundQuery) End() token.Pos {
	if c.Limit != nil {
		return c.Limit.End()
	}
	if c.OrderBy != nil {
		return c.OrderBy.End()
	}
	return c.Queries[len(c.Queries)-1].End()
}

func (s *SubQuery) Pos() token.Pos {
	return s.Lparen
}

func (s *SubQuery) End() token.Pos {
	if s.Limit != nil {
		return s.Limit.End()
	}
	if s.OrderBy != nil {
		return s.OrderBy.End()
	}
	return s.Rparen + 1
}

func (s *Star) Pos() token.Pos { return s.Star }
func (s *Star) End() token.Pos { return s.Star + 1 }

func (s *DotStar) Pos() token.Pos { return s.Expr.Pos() }
func (s *DotStar) End() token.Pos { return s.Star + 1 }

func (a *Alias) Pos() token.Pos { return a.Expr.Pos() }
func (a *Alias) End() token.Pos { return a.As.End() }

func (a *AsAlias) Pos() token.Pos {
	if !a.As.Invalid() {
		return a.As
	}
	return a.Alias.Pos()
}

func (a *AsAlias) End() token.Pos {
	return a.Alias.End()
}

func (e *ExprSelectItem) Pos() token.Pos { return e.Expr.Pos() }
func (e *ExprSelectItem) End() token.Pos { return e.Expr.End() }

func (f *From) Pos() token.Pos { return f.From }
func (f *From) End() token.Pos { return f.Source.End() }

func (w *Where) Pos() token.Pos { return w.Where }
func (w *Where) End() token.Pos { return w.Expr.End() }

func (g *GroupBy) Pos() token.Pos { return g.Group }
func (g *GroupBy) End() token.Pos { return g.Exprs[len(g.Exprs)-1].End() }

func (h *Having) Pos() token.Pos { return h.Having }
func (h *Having) End() token.Pos { return h.Expr.End() }

func (o *OrderBy) Pos() token.Pos { return o.Order }
func (o *OrderBy) End() token.Pos { return o.Items[len(o.Items)-1].End() }

func (o *OrderByItem) Pos() token.Pos { return o.Expr.Pos() }

func (o *OrderByItem) End() token.Pos {
	if !o.DirPos.Invalid() {
		return o.DirPos + token.Pos(len(o.Dir))
	}
	if o.Collate != nil {
		return o.Collate.End()
	}
	return o.Expr.End()
}

func (c *Collate) Pos() token.Pos { return c.Collate }
func (c *Collate) End() token.Pos { return c.Value.End() }

func (l *Limit) Pos() token.Pos {
	return l.Limit
}

func (l *Limit) End() token.Pos {
	if l.Offset != nil {
		return l.Offset.End()
	}
	return l.Count.End()
}

func (o *Offset) Pos() token.Pos { return o.Offset }
func (o *Offset) End() token.Pos { return o.Value.End() }

// ================================================================================
//
// JOIN
//
// ================================================================================

func (u *Unnest) Pos() token.Pos {
	if !u.Unnest.Invalid() {
		return u.Unnest
	}
	return u.Expr.Pos()
}

func (u *Unnest) End() token.Pos {
	if u.Sample != nil {
		return u.Sample.End()
	}
	if u.WithOffset != nil {
		return u.WithOffset.End()
	}
	if u.As != nil {
		return u.As.End()
	}
	if u.Hint != nil {
		return u.Hint.End()
	}
	if !u.Rparen.Invalid() {
		return u.Rparen + 1
	}
	return u.Expr.End()
}

func (w *WithOffset) Pos() token.Pos {
	return w.With
}

func (w *WithOffset) End() token.Pos {
	if w.As != nil {
		return w.As.End()
	}
	return w.Offset + 6
}

func (t *TableName) Pos() token.Pos {
	return t.Table.Pos()
}

func (t *TableName) End() token.Pos {
	if t.Sample != nil {
		return t.Sample.End()
	}
	if t.As != nil {
		return t.As.End()
	}
	if t.Hint != nil {
		return t.Hint.End()
	}
	return t.Table.End()
}

func (s *SubQueryTableExpr) Pos() token.Pos {
	return s.Lparen
}

func (s *SubQueryTableExpr) End() token.Pos {
	if s.Sample != nil {
		return s.Sample.End()
	}
	if s.As != nil {
		return s.As.End()
	}
	return s.Rparen + 1
}

func (p *ParenTableExpr) Pos() token.Pos {
	return p.Lparen
}

func (p *ParenTableExpr) End() token.Pos {
	if p.Sample != nil {
		return p.Sample.End()
	}
	return p.Rparen + 1
}

func (j *Join) Pos() token.Pos {
	return j.Left.Pos()
}

func (j *Join) End() token.Pos {
	if j.Cond != nil {
		return j.Cond.End()
	}
	return j.Right.End()
}

func (o *On) Pos() token.Pos { return o.On }
func (o *On) End() token.Pos { return o.Expr.End() }

func (u *Using) Pos() token.Pos { return u.Using }
func (u *Using) End() token.Pos { return u.Rparen + 1 }

func (t *TableSample) Pos() token.Pos { return t.TableSample }
func (t *TableSample) End() token.Pos { return t.Size.End() }

func (t *TableSampleSize) Pos() token.Pos { return t.Lparen }
func (t *TableSampleSize) End() token.Pos { return t.Rparen + 1 }

// ================================================================================
//
// Expr
//
// ================================================================================

func (b *BinaryExpr) Pos() token.Pos { return b.Left.Pos() }
func (b *BinaryExpr) End() token.Pos { return b.Right.End() }

func (u *UnaryExpr) Pos() token.Pos { return u.OpPos }
func (u *UnaryExpr) End() token.Pos { return u.Expr.End() }

func (i *InExpr) Pos() token.Pos { return i.Left.Pos() }
func (i *InExpr) End() token.Pos { return i.Right.End() }

func (u *UnnestInCondition) Pos() token.Pos { return u.Unnest }
func (u *UnnestInCondition) End() token.Pos { return u.Rparen + 1 }

func (s *SubQueryInCondition) Pos() token.Pos { return s.Lparen }
func (s *SubQueryInCondition) End() token.Pos { return s.Rparen + 1 }

func (v *ValuesInCondition) Pos() token.Pos { return v.Lparen }
func (v *ValuesInCondition) End() token.Pos { return v.Rparen + 1 }

func (i *IsNullExpr) Pos() token.Pos { return i.Left.Pos() }
func (i *IsNullExpr) End() token.Pos { return i.Null + 4 }

func (i *IsBoolExpr) Pos() token.Pos {
	return i.Left.Pos()
}

func (i *IsBoolExpr) End() token.Pos {
	if i.Right {
		return i.RightPos + 4
	} else {
		return i.RightPos + 5
	}
}

func (b *BetweenExpr) Pos() token.Pos { return b.Left.Pos() }
func (b *BetweenExpr) End() token.Pos { return b.RightEnd.End() }

func (b *SelectorExpr) Pos() token.Pos { return b.Expr.Pos() }
func (b *SelectorExpr) End() token.Pos { return b.Ident.End() }

func (i *IndexExpr) Pos() token.Pos { return i.Expr.Pos() }
func (i *IndexExpr) End() token.Pos { return i.Rbrack + 1 }

func (c *CallExpr) Pos() token.Pos { return c.Func.Pos() }
func (c *CallExpr) End() token.Pos { return c.Rparen + 1 }

func (a *Arg) Pos() token.Pos {
	if !a.Interval.Invalid() {
		return a.Interval
	}
	return a.Expr.Pos()
}

func (a *Arg) End() token.Pos {
	if a.IntervalUnit != nil {
		return a.IntervalUnit.End()
	}
	return a.Expr.End()
}

func (c *CountStarExpr) Pos() token.Pos { return c.Count }
func (c *CountStarExpr) End() token.Pos { return c.Rparen + 1 }

func (e *ExtractExpr) Pos() token.Pos { return e.Extract }
func (e *ExtractExpr) End() token.Pos { return e.Rparen + 1 }

func (a *AtTimeZone) Pos() token.Pos { return a.At }
func (a *AtTimeZone) End() token.Pos { return a.Expr.End() }

func (c *CastExpr) Pos() token.Pos { return c.Cast }
func (c *CastExpr) End() token.Pos { return c.Rparen + 1 }

func (c *CaseExpr) Pos() token.Pos { return c.Case }
func (c *CaseExpr) End() token.Pos { return c.EndPos + 3 }

func (c *CaseWhen) Pos() token.Pos { return c.When }
func (c *CaseWhen) End() token.Pos { return c.Then.End() }

func (c *CaseElse) Pos() token.Pos { return c.Else }
func (c *CaseElse) End() token.Pos { return c.Expr.End() }

func (p *ParenExpr) Pos() token.Pos { return p.Lparen }
func (p *ParenExpr) End() token.Pos { return p.Rparen + 1 }

func (s *ScalarSubQuery) Pos() token.Pos { return s.Lparen }
func (s *ScalarSubQuery) End() token.Pos { return s.Rparen + 1 }

func (a *ArraySubQuery) Pos() token.Pos { return a.Array }
func (a *ArraySubQuery) End() token.Pos { return a.Rparen + 1 }

func (e *ExistsSubQuery) Pos() token.Pos { return e.Exists }
func (e *ExistsSubQuery) End() token.Pos { return e.Rparen + 1 }

// ================================================================================
//
// Literal
//
// ================================================================================

func (p *Param) Pos() token.Pos { return p.Atmark }
func (p *Param) End() token.Pos { return p.Atmark + 1 + token.Pos(len(p.Name)) }

func (i *Ident) Pos() token.Pos { return i.NamePos }
func (i *Ident) End() token.Pos { return i.NameEnd }

func (p *Path) Pos() token.Pos { return p.Idents[0].Pos() }
func (p *Path) End() token.Pos { return p.Idents[len(p.Idents)-1].End() }

func (a *ArrayLiteral) Pos() token.Pos {
	if !a.Array.Invalid() {
		return a.Array
	}
	return a.Lbrack
}

func (a *ArrayLiteral) End() token.Pos {
	return a.Rbrack + 1
}

func (s *StructLiteral) Pos() token.Pos {
	if !s.Struct.Invalid() {
		return s.Struct
	}
	return s.Lparen
}

func (s *StructLiteral) End() token.Pos {
	return s.Rparen + 1
}

func (n *NullLiteral) Pos() token.Pos { return n.Null }
func (n *NullLiteral) End() token.Pos { return n.Null + 4 }

func (b *BoolLiteral) Pos() token.Pos {
	return b.ValuePos
}

func (b *BoolLiteral) End() token.Pos {
	if b.Value {
		return b.ValuePos + 4
	} else {
		return b.ValuePos + 5
	}
}

func (i *IntLiteral) Pos() token.Pos { return i.ValuePos }
func (i *IntLiteral) End() token.Pos { return i.ValueEnd }

func (f *FloatLiteral) Pos() token.Pos { return f.ValuePos }
func (f *FloatLiteral) End() token.Pos { return f.ValueEnd }

func (s *StringLiteral) Pos() token.Pos { return s.ValuePos }
func (s *StringLiteral) End() token.Pos { return s.ValueEnd }

func (b *BytesLiteral) Pos() token.Pos { return b.ValuePos }
func (b *BytesLiteral) End() token.Pos { return b.ValueEnd }

func (d *DateLiteral) Pos() token.Pos { return d.Date }
func (d *DateLiteral) End() token.Pos { return d.Value.End() }

func (t *TimestampLiteral) Pos() token.Pos { return t.Timestamp }
func (t *TimestampLiteral) End() token.Pos { return t.Value.End() }

func (t *NumericLiteral) Pos() token.Pos { return t.Numeric }
func (t *NumericLiteral) End() token.Pos { return t.Value.End() }

// ================================================================================
//
// Type
//
// ================================================================================

func (s *SimpleType) Pos() token.Pos { return s.NamePos }
func (s *SimpleType) End() token.Pos { return s.NamePos + token.Pos(len(s.Name)) }

func (a *ArrayType) Pos() token.Pos { return a.Array }
func (a *ArrayType) End() token.Pos { return a.Gt + 1 }

func (s *StructType) Pos() token.Pos { return s.Struct }
func (s *StructType) End() token.Pos { return s.Gt + 1 }

func (f *StructField) Pos() token.Pos {
	if f.Ident != nil {
		return f.Ident.Pos()
	}
	return f.Type.Pos()
}

func (f *StructField) End() token.Pos {
	return f.Type.End()
}

// ================================================================================
//
// Cast for Special Cases
//
// ================================================================================

func (c *CastIntValue) Pos() token.Pos { return c.Cast }
func (c *CastIntValue) End() token.Pos { return c.Rparen + 1 }

func (c *CastNumValue) Pos() token.Pos { return c.Cast }
func (c *CastNumValue) End() token.Pos { return c.Rparen + 1 }

// ================================================================================
//
// DDL
//
// ================================================================================

func (c *CreateDatabase) Pos() token.Pos { return c.Create }
func (c *CreateDatabase) End() token.Pos { return c.Name.End() }

func (c *CreateTable) Pos() token.Pos {
	return c.Create
}

func (c *CreateTable) End() token.Pos {
	if c.Cluster != nil {
		return c.Cluster.End()
	}
	return c.Rparen + 1
}

func (c *ColumnDef) Pos() token.Pos {
	return c.Name.Pos()
}

func (c *ColumnDef) End() token.Pos {
	if c.Options != nil {
		return c.Options.End()
	}
	if c.GeneratedExpr != nil {
		return c.GeneratedExpr.End()
	}
	if !c.Null.Invalid() {
		return c.Null + 4
	}
	return c.Type.End()
}

func (g *GeneratedColumnExpr) Pos() token.Pos { return g.As }
func (g *GeneratedColumnExpr) End() token.Pos { return g.Stored }

func (c *ColumnDefOptions) Pos() token.Pos { return c.Options }
func (c *ColumnDefOptions) End() token.Pos { return c.Rparen + 1 }

func (f *ForeignKey) Pos() token.Pos {
	if f.Name == nil {
		return f.Foreign
	}
	return f.Constraint
}
func (f *ForeignKey) End() token.Pos { return f.Rparen + 1 }

func (i *IndexKey) Pos() token.Pos {
	return i.Name.Pos()
}

func (i *IndexKey) End() token.Pos {
	if !i.DirPos.Invalid() {
		return i.DirPos + token.Pos(len(i.Dir))
	}
	return i.Name.End()
}

func (c *Cluster) Pos() token.Pos {
	return c.Comma
}

func (c *Cluster) End() token.Pos {
	if !c.OnDeleteEnd.Invalid() {
		return c.OnDeleteEnd
	}
	return c.TableName.End()
}

func (a *AlterTable) Pos() token.Pos { return a.Alter }
func (a *AlterTable) End() token.Pos { return a.TableAlternation.End() }

func (a *AddColumn) Pos() token.Pos { return a.Add }
func (a *AddColumn) End() token.Pos { return a.Column.End() }

func (a *AddForeignKey) Pos() token.Pos { return a.Add }
func (a *AddForeignKey) End() token.Pos { return a.ForeignKey.End() }

func (d *DropColumn) Pos() token.Pos { return d.Drop }
func (d *DropColumn) End() token.Pos { return d.Name.End() }

func (d *DropConstraint) Pos() token.Pos { return d.Drop }
func (d *DropConstraint) End() token.Pos { return d.Name.End() }

func (s *SetOnDelete) Pos() token.Pos { return s.Set }
func (s *SetOnDelete) End() token.Pos { return s.OnDeleteEnd }

func (a *AlterColumn) Pos() token.Pos {
	return a.Alter
}

func (a *AlterColumn) End() token.Pos {
	if !a.Null.Invalid() {
		return a.Null + 4
	}
	return a.Type.End()
}

func (a *AlterColumnSet) Pos() token.Pos { return a.Alter }
func (a *AlterColumnSet) End() token.Pos { return a.Options.End() }

func (d *DropTable) Pos() token.Pos { return d.Drop }
func (d *DropTable) End() token.Pos { return d.Name.End() }

func (c *CreateIndex) Pos() token.Pos {
	return c.Create
}

func (c *CreateIndex) End() token.Pos {
	if c.Storing != nil {
		return c.Storing.End()
	}
	return c.Rparen + 1
}

func (s *Storing) Pos() token.Pos { return s.Storing }
func (s *Storing) End() token.Pos { return s.Rparen + 1 }

func (i *InterleaveIn) Pos() token.Pos { return i.Comma }
func (i *InterleaveIn) End() token.Pos { return i.TableName.End() }

func (d *DropIndex) Pos() token.Pos { return d.Drop }
func (d *DropIndex) End() token.Pos { return d.Name.End() }

// ================================================================================
//
// Types for Schema
//
// ================================================================================

func (s *ScalarSchemaType) Pos() token.Pos { return s.NamePos }
func (s *ScalarSchemaType) End() token.Pos { return s.NamePos + token.Pos(len(s.Name)) }

func (s *SizedSchemaType) Pos() token.Pos { return s.NamePos }
func (s *SizedSchemaType) End() token.Pos { return s.Rparen + 1 }

func (a *ArraySchemaType) Pos() token.Pos { return a.Array }
func (a *ArraySchemaType) End() token.Pos { return a.Gt + 1 }

// ================================================================================
//
// DML
//
// ================================================================================

func (i *Insert) Pos() token.Pos { return i.Insert }
func (i *Insert) End() token.Pos { return i.Input.End() }

func (v *ValuesInput) Pos() token.Pos { return v.Values }
func (v *ValuesInput) End() token.Pos { return v.Rows[len(v.Rows)-1].End() }

func (v *ValuesRow) Pos() token.Pos { return v.Lparen }
func (v *ValuesRow) End() token.Pos { return v.Rparen + 1 }

func (d *DefaultExpr) Pos() token.Pos {
	if !d.DefaultPos.Invalid() {
		return d.DefaultPos
	}
	return d.Expr.Pos()
}

func (d *DefaultExpr) End() token.Pos {
	if !d.DefaultPos.Invalid() {
		return d.DefaultPos + 7
	}
	return d.Expr.End()
}

func (s *SubQueryInput) Pos() token.Pos { return s.Query.Pos() }
func (s *SubQueryInput) End() token.Pos { return s.Query.End() }

func (d *Delete) Pos() token.Pos { return d.Delete }
func (d *Delete) End() token.Pos { return d.Where.End() }

func (u *Update) Pos() token.Pos { return u.Update }
func (u *Update) End() token.Pos { return u.Where.End() }

func (u *UpdateItem) Pos() token.Pos { return u.Path[0].Pos() }
func (u *UpdateItem) End() token.Pos { return u.Expr.End() }
