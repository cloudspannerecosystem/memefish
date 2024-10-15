package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
)

// ================================================================================
//
// Helper functions for Pos(), End()
// These functions are intended for use within this file only.
//
// ================================================================================

// lastNode returns last element of Node slice.
// This function corresponds to NodeSliceVar[$] in ast.go.
func lastNode[T Node](s []T) T {
	return s[len(s)-1]
}

// firstValidEnd returns the first valid Pos() in argument.
// "valid" means the node is not nil and Pos().Invalid() is not true.
// This function corresponds to "(n0 ?? n1 ?? ...).End()"
func firstValidEnd(ns ...Node) token.Pos {
	for _, n := range ns {
		if n != nil && !n.End().Invalid() {
			return n.End()
		}
	}
	return token.InvalidPos
}

// firstPos returns the Pos() of the first node.
// If argument is an empty slice, this function returns token.InvalidPos.
// This function corresponds to NodeSliceVar[0].pos in ast.go.
func firstPos[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return s[0].Pos()
}

// lastEnd returns the End() of the last node.
// If argument is an empty slice, this function returns token.InvalidPos.
// This function corresponds to NodeSliceVar[$].end in ast.go.
func lastEnd[T Node](s []T) token.Pos {
	if len(s) == 0 {
		return token.InvalidPos
	}
	return lastNode(s).End()
}

// ================================================================================
//
// SELECT
//
// ================================================================================

func (q *QueryStatement) Pos() token.Pos {
	if q.Hint != nil {
		return q.Hint.Pos()
	}
	if q.With != nil {
		return q.With.Pos()
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

func (w *With) Pos() token.Pos { return w.With }
func (w *With) End() token.Pos { return w.CTEs[len(w.CTEs)-1].End() }

func (c *CTE) Pos() token.Pos { return c.Name.Pos() }
func (c *CTE) End() token.Pos { return c.Rparen + 1 }

func (s *Select) Pos() token.Pos { return s.Select }

func (s *Select) End() token.Pos {
	return firstValidEnd(s.Limit, s.OrderBy, s.Having, s.GroupBy, s.Where, s.From, lastNode(s.Results))
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

func (e *ExprArg) Pos() token.Pos { return e.Expr.Pos() }
func (e *ExprArg) End() token.Pos { return e.Expr.End() }

func (i *IntervalArg) Pos() token.Pos { return i.Interval }
func (i *IntervalArg) End() token.Pos {
	if i.Unit != nil {
		return i.Unit.End()
	}
	return i.Expr.End()
}

func (s *SequenceArg) Pos() token.Pos { return s.Sequence }
func (s *SequenceArg) End() token.Pos { return s.Expr.End() }

func (n *NamedArg) Pos() token.Pos { return n.Name.Pos() }
func (n *NamedArg) End() token.Pos { return n.Value.End() }

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

func (p *Path) Pos() token.Pos { return firstPos(p.Idents) }
func (p *Path) End() token.Pos { return lastEnd(p.Idents) }

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

func (t *JSONLiteral) Pos() token.Pos { return t.JSON }
func (t *JSONLiteral) End() token.Pos { return t.Value.End() }

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

func (n *NamedType) Pos() token.Pos { return n.Path[0].Pos() }
func (n *NamedType) End() token.Pos { return n.Path[len(n.Path)-1].End() }

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

func (g *GenericOptions) Pos() token.Pos { return g.Options }
func (g *GenericOptions) End() token.Pos { return g.Rparen + 1 }

func (g *GenericOption) Pos() token.Pos { return g.Name.Pos() }
func (g *GenericOption) End() token.Pos { return g.Value.End() }

func (c *CreateDatabase) Pos() token.Pos { return c.Create }
func (c *CreateDatabase) End() token.Pos { return c.Name.End() }

func (c *CreateTable) Pos() token.Pos {
	return c.Create
}

func (c *CreateTable) End() token.Pos {
	if c.RowDeletionPolicy != nil {
		return c.RowDeletionPolicy.End()
	}
	if c.Cluster != nil {
		return c.Cluster.End()
	}
	return c.Rparen + 1
}

func (c *CreateSequence) Pos() token.Pos {
	return c.Create
}

func (c *CreateSequence) End() token.Pos {
	return c.Options.End()
}

func (c *CreateView) Pos() token.Pos {
	return c.Create
}

func (c *CreateView) End() token.Pos {
	return c.Query.End()
}

func (c *ColumnDef) Pos() token.Pos {
	return c.Name.Pos()
}

func (c *ColumnDef) End() token.Pos {
	// TODO: It may be able to be refactored using Pos arithmetic like InvalidPos + n = InvalidPos.
	if c.Options != nil {
		return c.Options.End()
	}
	if !c.Hidden.Invalid() {
		return c.Hidden + 6
	}
	if end := firstValidEnd(c.GeneratedExpr, c.DefaultExpr); !end.Invalid() {
		return end
	}
	if !c.Null.Invalid() {
		return c.Null + 4
	}
	return c.Type.End()
}

func (g *ColumnDefaultExpr) Pos() token.Pos { return g.Default }
func (g *ColumnDefaultExpr) End() token.Pos { return g.Rparen }

func (g *GeneratedColumnExpr) Pos() token.Pos { return g.As }
func (g *GeneratedColumnExpr) End() token.Pos {
	if !g.Stored.Invalid() {
		return g.Stored + 6
	}
	return g.Rparen + 1
}

func (c *ColumnDefOptions) Pos() token.Pos { return c.Options }
func (c *ColumnDefOptions) End() token.Pos { return c.Rparen + 1 }

func (c *TableConstraint) Pos() token.Pos {
	if c.Name != nil {
		return c.ConstraintPos
	}
	return c.Constraint.Pos()
}

func (c *TableConstraint) End() token.Pos { return c.Constraint.End() }

func (f *ForeignKey) Pos() token.Pos {
	return f.Foreign
}
func (f *ForeignKey) End() token.Pos {
	if !f.OnDeleteEnd.Invalid() {
		return f.OnDeleteEnd
	}
	return f.Rparen + 1
}

func (c *Check) Pos() token.Pos { return c.Check }
func (c *Check) End() token.Pos { return c.Rparen + 1 }

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

func (c *CreateRowDeletionPolicy) Pos() token.Pos {
	return c.Comma
}

func (c *CreateRowDeletionPolicy) End() token.Pos {
	return c.RowDeletionPolicy.End()
}

func (r *RowDeletionPolicy) Pos() token.Pos {
	return r.Row
}

func (r *RowDeletionPolicy) End() token.Pos {
	return r.Rparen + 1
}

func (a *AlterTable) Pos() token.Pos { return a.Alter }
func (a *AlterTable) End() token.Pos { return a.TableAlteration.End() }

func (a *AddColumn) Pos() token.Pos { return a.Add }
func (a *AddColumn) End() token.Pos { return a.Column.End() }

func (a *AddRowDeletionPolicy) Pos() token.Pos { return a.Add }
func (a *AddRowDeletionPolicy) End() token.Pos { return a.RowDeletionPolicy.End() }

func (a *AddTableConstraint) Pos() token.Pos { return a.Add }
func (a *AddTableConstraint) End() token.Pos { return a.TableConstraint.End() }

func (d *DropColumn) Pos() token.Pos { return d.Drop }
func (d *DropColumn) End() token.Pos { return d.Name.End() }

func (d *DropConstraint) Pos() token.Pos { return d.Drop }
func (d *DropConstraint) End() token.Pos { return d.Name.End() }

func (d *DropRowDeletionPolicy) Pos() token.Pos { return d.Drop }
func (d *DropRowDeletionPolicy) End() token.Pos { return d.Policy + 6 }

func (d *ReplaceRowDeletionPolicy) Pos() token.Pos { return d.Replace }
func (d *ReplaceRowDeletionPolicy) End() token.Pos { return d.RowDeletionPolicy.End() }

func (s *SetOnDelete) Pos() token.Pos { return s.Set }
func (s *SetOnDelete) End() token.Pos { return s.OnDeleteEnd }

func (a *AlterColumn) Pos() token.Pos {
	return a.Alter
}

func (a *AlterColumn) End() token.Pos {
	if a.DefaultExpr != nil {
		return a.DefaultExpr.End()
	}
	if !a.Null.Invalid() {
		return a.Null + 4
	}
	return a.Type.End()
}

func (a *AlterColumnSet) Pos() token.Pos { return a.Alter }

func (a *AlterColumnSet) End() token.Pos {
	if a.Options != nil {
		return a.Options.End()
	}
	return a.DefaultExpr.End()
}

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

func (c *CreateVectorIndex) Pos() token.Pos {
	return c.Create
}

func (c *CreateVectorIndex) End() token.Pos {
	return c.Options.Rparen + 1
}

func (c *CreateChangeStream) Pos() token.Pos {
	return c.Create
}

func (c *CreateChangeStream) End() token.Pos {
	if c.Options != nil {
		return c.Options.End()
	}
	if c.For != nil {
		return c.For.End()
	}
	return c.Name.End()
}

func (c *ChangeStreamOptions) Pos() token.Pos { return c.Options }
func (c *ChangeStreamOptions) End() token.Pos { return c.Rparen + 1 }

func (c *ChangeStreamForAll) Pos() token.Pos    { return c.For }
func (c *ChangeStreamForAll) End() token.Pos    { return c.All }
func (c *ChangeStreamForTables) Pos() token.Pos { return c.For }
func (c *ChangeStreamForTables) End() token.Pos { return c.Tables[len(c.Tables)-1].End() }

func (s *Storing) Pos() token.Pos { return s.Storing }
func (s *Storing) End() token.Pos { return s.Rparen + 1 }

func (i *InterleaveIn) Pos() token.Pos { return i.Comma }
func (i *InterleaveIn) End() token.Pos { return i.TableName.End() }

func (a *AlterIndex) Pos() token.Pos { return a.Alter }
func (a *AlterIndex) End() token.Pos { return a.IndexAlteration.End() }

func (a *AlterSequence) Pos() token.Pos { return a.Alter }
func (a *AlterSequence) End() token.Pos { return a.Options.End() }

func (a *AddStoredColumn) Pos() token.Pos { return a.Add }
func (a *AddStoredColumn) End() token.Pos { return a.Name.End() }

func (a *DropStoredColumn) Pos() token.Pos { return a.Drop }
func (a *DropStoredColumn) End() token.Pos { return a.Name.End() }

func (d *DropIndex) Pos() token.Pos { return d.Drop }
func (d *DropIndex) End() token.Pos { return d.Name.End() }

func (d *DropVectorIndex) Pos() token.Pos { return d.Drop }
func (d *DropVectorIndex) End() token.Pos { return d.Name.End() }

func (d *DropSequence) Pos() token.Pos { return d.Drop }
func (d *DropSequence) End() token.Pos { return d.Name.End() }

func (c *CreateRole) Pos() token.Pos { return c.Create }
func (c *CreateRole) End() token.Pos { return c.Name.End() }

func (d *DropRole) Pos() token.Pos { return d.Drop }
func (d *DropRole) End() token.Pos { return d.Name.End() }

func (d *DropChangeStream) Pos() token.Pos { return d.Drop }
func (d *DropChangeStream) End() token.Pos { return d.Name.End() }

func (a *AlterChangeStream) Pos() token.Pos { return a.Alter }
func (a *AlterChangeStream) End() token.Pos { return a.ChangeStreamAlteration.End() }

func (a *ChangeStreamSetFor) Pos() token.Pos { return a.Set }
func (a *ChangeStreamSetFor) End() token.Pos { return a.For.End() }

func (a *ChangeStreamDropForAll) Pos() token.Pos { return a.Drop }
func (a *ChangeStreamDropForAll) End() token.Pos { return a.All + 3 }

func (a *ChangeStreamSetOptions) Pos() token.Pos { return a.Set }
func (a *ChangeStreamSetOptions) End() token.Pos { return a.Options.Rparen + 1 }

func (g *Grant) Pos() token.Pos { return g.Grant }
func (g *Grant) End() token.Pos { return g.Roles[len(g.Roles)-1].End() }

func (r *Revoke) Pos() token.Pos { return r.Revoke }
func (r *Revoke) End() token.Pos { return r.Roles[len(r.Roles)-1].End() }

func (p *PrivilegeOnTable) Pos() token.Pos { return p.Privileges[0].Pos() }
func (p *PrivilegeOnTable) End() token.Pos { return p.Names[len(p.Names)-1].End() }

func (s *SelectPrivilege) Pos() token.Pos {
	return s.Select
}

func (s *SelectPrivilege) End() token.Pos {
	if !s.Rparen.Invalid() {
		return s.Rparen + 1
	}
	return s.Select + 6
}

func (i *InsertPrivilege) Pos() token.Pos {
	return i.Insert
}

func (i *InsertPrivilege) End() token.Pos {
	if !i.Rparen.Invalid() {
		return i.Rparen + 1
	}
	return i.Insert + 6
}

func (u *UpdatePrivilege) Pos() token.Pos {
	return u.Update
}

func (u *UpdatePrivilege) End() token.Pos {
	if !u.Rparen.Invalid() {
		return u.Rparen + 1
	}
	return u.Update + 6
}

func (d *DeletePrivilege) Pos() token.Pos { return d.Delete }
func (d *DeletePrivilege) End() token.Pos { return d.Delete + 6 }

func (s *SelectPrivilegeOnView) Pos() token.Pos { return s.Select }
func (s *SelectPrivilegeOnView) End() token.Pos { return s.Names[len(s.Names)-1].End() }

func (e *ExecutePrivilegeOnTableFunction) Pos() token.Pos { return e.Execute }
func (e *ExecutePrivilegeOnTableFunction) End() token.Pos { return e.Names[len(e.Names)-1].End() }

func (r *RolePrivilege) Pos() token.Pos { return r.Role }
func (r *RolePrivilege) End() token.Pos { return r.Names[len(r.Names)-1].End() }

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

func (c *ChangeStreamForTable) End() token.Pos {
	if len(c.Columns) == 0 {
		return c.TableName.End()
	}
	return c.Rparen + 1
}

// ================================================================================
//
// Search Index DDL
//
// ================================================================================

func (c *CreateSearchIndex) Pos() token.Pos { return c.Create }

func (c *CreateSearchIndex) End() token.Pos {
	if end := firstValidEnd(c.Options, c.Interleave, c.Where, c.OrderBy, lastNode(c.PartitionColumns), c.Storing); end != token.InvalidPos {
		return end
	}
	return c.Rparen + 1
}

func (o *SearchIndexOptions) Pos() token.Pos { return (*GenericOptions)(o).Pos() }
func (o *SearchIndexOptions) End() token.Pos { return (*GenericOptions)(o).End() }

func (d *DropSearchIndex) Pos() token.Pos {
	return d.Drop
}

func (d *DropSearchIndex) End() token.Pos {
	return d.Name.End()
}

func (a *AlterSearchIndex) Pos() token.Pos { return a.Alter }
func (a *AlterSearchIndex) End() token.Pos { return a.IndexAlteration.End() }

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

func (o *SequenceOption) Pos() token.Pos { return o.Name.Pos() }
func (o *SequenceOption) End() token.Pos { return o.Value.End() }

func (o *SequenceOptions) Pos() token.Pos { return o.Options }
func (o *SequenceOptions) End() token.Pos { return o.Rparen + 1 }
