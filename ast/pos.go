// Code generated by tools/gen-ast-pos; DO NOT EDIT.

package ast

import (
	"github.com/cloudspannerecosystem/memefish/token"
)

func (q *QueryStatement) Pos() token.Pos {
	return nodePos(nodeChoice(wrapNode(q.Hint), wrapNode(q.With), wrapNode(q.Query)))
}

func (q *QueryStatement) End() token.Pos {
	return nodeEnd(wrapNode(q.Query))
}

func (h *Hint) Pos() token.Pos {
	return h.Atmark
}

func (h *Hint) End() token.Pos {
	return posAdd(h.Rbrace, 1)
}

func (h *HintRecord) Pos() token.Pos {
	return nodePos(wrapNode(h.Key))
}

func (h *HintRecord) End() token.Pos {
	return nodeEnd(wrapNode(h.Value))
}

func (w *With) Pos() token.Pos {
	return w.With
}

func (w *With) End() token.Pos {
	return nodeEnd(nodeSliceLast(w.CTEs))
}

func (c *CTE) Pos() token.Pos {
	return nodePos(wrapNode(c.Name))
}

func (c *CTE) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (s *Select) Pos() token.Pos {
	return s.Select
}

func (s *Select) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(s.Limit), wrapNode(s.OrderBy), wrapNode(s.Having), wrapNode(s.GroupBy), wrapNode(s.Where), wrapNode(s.From), nodeSliceLast(s.Results)))
}

func (a *AsStruct) Pos() token.Pos {
	return a.As
}

func (a *AsStruct) End() token.Pos {
	return posAdd(a.Struct, 6)
}

func (a *AsValue) Pos() token.Pos {
	return a.As
}

func (a *AsValue) End() token.Pos {
	return posAdd(a.Value, 5)
}

func (a *AsTypeName) Pos() token.Pos {
	return a.As
}

func (a *AsTypeName) End() token.Pos {
	return nodeEnd(wrapNode(a.TypeName))
}

func (c *CompoundQuery) Pos() token.Pos {
	return nodePos(nodeSliceIndex(c.Queries, 0))
}

func (c *CompoundQuery) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(c.Limit), wrapNode(c.OrderBy), nodeSliceLast(c.Queries)))
}

func (s *SubQuery) Pos() token.Pos {
	return s.Lparen
}

func (s *SubQuery) End() token.Pos {
	return posChoice(nodeEnd(nodeChoice(wrapNode(s.Limit), wrapNode(s.OrderBy))), posAdd(s.Rparen, 1))
}

func (s *Star) Pos() token.Pos {
	return s.Star
}

func (s *Star) End() token.Pos {
	return posAdd(s.Star, 1)
}

func (d *DotStar) Pos() token.Pos {
	return nodePos(wrapNode(d.Expr))
}

func (d *DotStar) End() token.Pos {
	return posAdd(d.Star, 1)
}

func (a *Alias) Pos() token.Pos {
	return nodePos(wrapNode(a.Expr))
}

func (a *Alias) End() token.Pos {
	return nodeEnd(wrapNode(a.As))
}

func (a *AsAlias) Pos() token.Pos {
	return posChoice(a.As, nodePos(wrapNode(a.Alias)))
}

func (a *AsAlias) End() token.Pos {
	return nodeEnd(wrapNode(a.Alias))
}

func (e *ExprSelectItem) Pos() token.Pos {
	return nodePos(wrapNode(e.Expr))
}

func (e *ExprSelectItem) End() token.Pos {
	return nodeEnd(wrapNode(e.Expr))
}

func (f *From) Pos() token.Pos {
	return f.From
}

func (f *From) End() token.Pos {
	return nodeEnd(wrapNode(f.Source))
}

func (w *Where) Pos() token.Pos {
	return w.Where
}

func (w *Where) End() token.Pos {
	return nodeEnd(wrapNode(w.Expr))
}

func (g *GroupBy) Pos() token.Pos {
	return g.Group
}

func (g *GroupBy) End() token.Pos {
	return nodeEnd(nodeSliceLast(g.Exprs))
}

func (h *Having) Pos() token.Pos {
	return h.Having
}

func (h *Having) End() token.Pos {
	return nodeEnd(wrapNode(h.Expr))
}

func (o *OrderBy) Pos() token.Pos {
	return o.Order
}

func (o *OrderBy) End() token.Pos {
	return nodeEnd(nodeSliceLast(o.Items))
}

func (o *OrderByItem) Pos() token.Pos {
	return nodePos(wrapNode(o.Expr))
}

func (o *OrderByItem) End() token.Pos {
	return posChoice(posAdd(o.DirPos, len(o.Dir)), nodeEnd(nodeChoice(wrapNode(o.Collate), wrapNode(o.Expr))))
}

func (c *Collate) Pos() token.Pos {
	return c.Collate
}

func (c *Collate) End() token.Pos {
	return nodeEnd(wrapNode(c.Value))
}

func (l *Limit) Pos() token.Pos {
	return l.Limit
}

func (l *Limit) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(l.Offset), wrapNode(l.Count)))
}

func (o *Offset) Pos() token.Pos {
	return o.Offset
}

func (o *Offset) End() token.Pos {
	return nodeEnd(wrapNode(o.Value))
}

func (u *Unnest) Pos() token.Pos {
	return u.Unnest
}

func (u *Unnest) End() token.Pos {
	return posChoice(nodeEnd(nodeChoice(wrapNode(u.Sample), wrapNode(u.WithOffset), wrapNode(u.As), wrapNode(u.Hint))), posAdd(u.Rparen, 1), nodeEnd(wrapNode(u.Expr)))
}

func (w *WithOffset) Pos() token.Pos {
	return w.With
}

func (w *WithOffset) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(w.As)), posAdd(w.Offset, 6))
}

func (t *TableName) Pos() token.Pos {
	return nodePos(wrapNode(t.Table))
}

func (t *TableName) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(t.Sample), wrapNode(t.As), wrapNode(t.Hint), wrapNode(t.Table)))
}

func (p *PathTableExpr) Pos() token.Pos {
	return nodePos(wrapNode(p.Path))
}

func (p *PathTableExpr) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(p.Sample), wrapNode(p.WithOffset), wrapNode(p.As), wrapNode(p.Hint), wrapNode(p.Path)))
}

func (s *SubQueryTableExpr) Pos() token.Pos {
	return s.Lparen
}

func (s *SubQueryTableExpr) End() token.Pos {
	return posChoice(nodeEnd(nodeChoice(wrapNode(s.Sample), wrapNode(s.As))), posAdd(s.Rparen, 1))
}

func (p *ParenTableExpr) Pos() token.Pos {
	return p.Lparen
}

func (p *ParenTableExpr) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(p.Sample)), posAdd(p.Rparen, 1))
}

func (j *Join) Pos() token.Pos {
	return nodePos(wrapNode(j.Left))
}

func (j *Join) End() token.Pos {
	return nodePos(nodeChoice(wrapNode(j.Cond), wrapNode(j.Right)))
}

func (o *On) Pos() token.Pos {
	return o.On
}

func (o *On) End() token.Pos {
	return nodeEnd(wrapNode(o.Expr))
}

func (u *Using) Pos() token.Pos {
	return u.Using
}

func (u *Using) End() token.Pos {
	return posAdd(u.Rparen, 1)
}

func (t *TableSample) Pos() token.Pos {
	return t.TableSample
}

func (t *TableSample) End() token.Pos {
	return nodeEnd(wrapNode(t.Size))
}

func (t *TableSampleSize) Pos() token.Pos {
	return t.Lparen
}

func (t *TableSampleSize) End() token.Pos {
	return posAdd(t.Rparen, 1)
}

func (b *BinaryExpr) Pos() token.Pos {
	return nodePos(wrapNode(b.Left))
}

func (b *BinaryExpr) End() token.Pos {
	return nodePos(wrapNode(b.Right))
}

func (u *UnaryExpr) Pos() token.Pos {
	return u.OpPos
}

func (u *UnaryExpr) End() token.Pos {
	return nodeEnd(wrapNode(u.Expr))
}

func (i *InExpr) Pos() token.Pos {
	return nodePos(wrapNode(i.Left))
}

func (i *InExpr) End() token.Pos {
	return nodeEnd(wrapNode(i.Right))
}

func (u *UnnestInCondition) Pos() token.Pos {
	return u.Unnest
}

func (u *UnnestInCondition) End() token.Pos {
	return posAdd(u.Rparen, 1)
}

func (s *SubQueryInCondition) Pos() token.Pos {
	return s.Lparen
}

func (s *SubQueryInCondition) End() token.Pos {
	return posAdd(s.Rparen, 1)
}

func (v *ValuesInCondition) Pos() token.Pos {
	return v.Lparen
}

func (v *ValuesInCondition) End() token.Pos {
	return posAdd(v.Rparen, 1)
}

func (i *IsNullExpr) Pos() token.Pos {
	return nodePos(wrapNode(i.Left))
}

func (i *IsNullExpr) End() token.Pos {
	return posAdd(i.Null, 4)
}

func (i *IsBoolExpr) Pos() token.Pos {
	return nodePos(wrapNode(i.Left))
}

func (i *IsBoolExpr) End() token.Pos {
	return posAdd(i.RightPos, ifThenElse(i.Right, 4, 5))
}

func (b *BetweenExpr) Pos() token.Pos {
	return nodePos(wrapNode(b.Left))
}

func (b *BetweenExpr) End() token.Pos {
	return nodeEnd(wrapNode(b.RightEnd))
}

func (s *SelectorExpr) Pos() token.Pos {
	return nodePos(wrapNode(s.Expr))
}

func (s *SelectorExpr) End() token.Pos {
	return nodePos(wrapNode(s.Ident))
}

func (i *IndexExpr) Pos() token.Pos {
	return nodePos(wrapNode(i.Expr))
}

func (i *IndexExpr) End() token.Pos {
	return posAdd(i.Rbrack, 1)
}

func (c *CallExpr) Pos() token.Pos {
	return nodePos(wrapNode(c.Func))
}

func (c *CallExpr) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (e *ExprArg) Pos() token.Pos {
	return nodePos(wrapNode(e.Expr))
}

func (e *ExprArg) End() token.Pos {
	return nodeEnd(wrapNode(e.Expr))
}

func (i *IntervalArg) Pos() token.Pos {
	return i.Interval
}

func (i *IntervalArg) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(i.Unit), wrapNode(i.Expr)))
}

func (s *SequenceArg) Pos() token.Pos {
	return s.Sequence
}

func (s *SequenceArg) End() token.Pos {
	return nodeEnd(wrapNode(s.Expr))
}

func (n *NamedArg) Pos() token.Pos {
	return nodePos(wrapNode(n.Name))
}

func (n *NamedArg) End() token.Pos {
	return nodeEnd(wrapNode(n.Value))
}

func (i *IgnoreNulls) Pos() token.Pos {
	return i.Ignore
}

func (i *IgnoreNulls) End() token.Pos {
	return posAdd(i.Nulls, 5)
}

func (r *RespectNulls) Pos() token.Pos {
	return r.Respect
}

func (r *RespectNulls) End() token.Pos {
	return posAdd(r.Nulls, 5)
}

func (h *HavingMax) Pos() token.Pos {
	return h.Having
}

func (h *HavingMax) End() token.Pos {
	return nodeEnd(wrapNode(h.Expr))
}

func (h *HavingMin) Pos() token.Pos {
	return h.Having
}

func (h *HavingMin) End() token.Pos {
	return nodeEnd(wrapNode(h.Expr))
}

func (c *CountStarExpr) Pos() token.Pos {
	return c.Count
}

func (c *CountStarExpr) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (e *ExtractExpr) Pos() token.Pos {
	return e.Extract
}

func (e *ExtractExpr) End() token.Pos {
	return posAdd(e.Rparen, 1)
}

func (a *AtTimeZone) Pos() token.Pos {
	return a.At
}

func (a *AtTimeZone) End() token.Pos {
	return nodeEnd(wrapNode(a.Expr))
}

func (c *CastExpr) Pos() token.Pos {
	return c.Cast
}

func (c *CastExpr) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (c *CaseExpr) Pos() token.Pos {
	return c.Case
}

func (c *CaseExpr) End() token.Pos {
	return posAdd(c.EndPos, 3)
}

func (c *CaseWhen) Pos() token.Pos {
	return c.When
}

func (c *CaseWhen) End() token.Pos {
	return nodeEnd(wrapNode(c.Then))
}

func (c *CaseElse) Pos() token.Pos {
	return c.Else
}

func (c *CaseElse) End() token.Pos {
	return nodeEnd(wrapNode(c.Expr))
}

func (p *ParenExpr) Pos() token.Pos {
	return p.Lparen
}

func (p *ParenExpr) End() token.Pos {
	return posAdd(p.Rparen, 1)
}

func (s *ScalarSubQuery) Pos() token.Pos {
	return s.Lparen
}

func (s *ScalarSubQuery) End() token.Pos {
	return posAdd(s.Rparen, 1)
}

func (a *ArraySubQuery) Pos() token.Pos {
	return a.Array
}

func (a *ArraySubQuery) End() token.Pos {
	return posAdd(a.Rparen, 1)
}

func (e *ExistsSubQuery) Pos() token.Pos {
	return e.Exists
}

func (e *ExistsSubQuery) End() token.Pos {
	return posAdd(e.Rparen, 1)
}

func (p *Param) Pos() token.Pos {
	return p.Atmark
}

func (p *Param) End() token.Pos {
	return posAdd(posAdd(p.Atmark, 1), len(p.Name))
}

func (i *Ident) Pos() token.Pos {
	return i.NamePos
}

func (i *Ident) End() token.Pos {
	return i.NameEnd
}

func (p *Path) Pos() token.Pos {
	return nodePos(nodeSliceIndex(p.Idents, 0))
}

func (p *Path) End() token.Pos {
	return nodeEnd(nodeSliceLast(p.Idents))
}

func (a *ArrayLiteral) Pos() token.Pos {
	return posChoice(a.Array, a.Lbrack)
}

func (a *ArrayLiteral) End() token.Pos {
	return posAdd(a.Rbrack, 1)
}

func (t *TupleStructLiteral) Pos() token.Pos {
	return t.Lparen
}

func (t *TupleStructLiteral) End() token.Pos {
	return posAdd(t.Rparen, 1)
}

func (t *TypedStructLiteral) Pos() token.Pos {
	return t.Struct
}

func (t *TypedStructLiteral) End() token.Pos {
	return posAdd(t.Rparen, 1)
}

func (t *TypelessStructLiteral) Pos() token.Pos {
	return t.Struct
}

func (t *TypelessStructLiteral) End() token.Pos {
	return posAdd(t.Rparen, 1)
}

func (n *NullLiteral) Pos() token.Pos {
	return n.Null
}

func (n *NullLiteral) End() token.Pos {
	return posAdd(n.Null, 4)
}

func (b *BoolLiteral) Pos() token.Pos {
	return b.ValuePos
}

func (b *BoolLiteral) End() token.Pos {
	return posAdd(b.ValuePos, ifThenElse(b.Value, 4, 5))
}

func (i *IntLiteral) Pos() token.Pos {
	return i.ValuePos
}

func (i *IntLiteral) End() token.Pos {
	return i.ValueEnd
}

func (f *FloatLiteral) Pos() token.Pos {
	return f.ValuePos
}

func (f *FloatLiteral) End() token.Pos {
	return f.ValueEnd
}

func (s *StringLiteral) Pos() token.Pos {
	return s.ValuePos
}

func (s *StringLiteral) End() token.Pos {
	return s.ValueEnd
}

func (b *BytesLiteral) Pos() token.Pos {
	return b.ValuePos
}

func (b *BytesLiteral) End() token.Pos {
	return b.ValueEnd
}

func (d *DateLiteral) Pos() token.Pos {
	return d.Date
}

func (d *DateLiteral) End() token.Pos {
	return nodeEnd(wrapNode(d.Value))
}

func (t *TimestampLiteral) Pos() token.Pos {
	return t.Timestamp
}

func (t *TimestampLiteral) End() token.Pos {
	return nodeEnd(wrapNode(t.Value))
}

func (n *NumericLiteral) Pos() token.Pos {
	return n.Numeric
}

func (n *NumericLiteral) End() token.Pos {
	return nodeEnd(wrapNode(n.Value))
}

func (j *JSONLiteral) Pos() token.Pos {
	return j.JSON
}

func (j *JSONLiteral) End() token.Pos {
	return nodeEnd(wrapNode(j.Value))
}

func (n *NewConstructor) Pos() token.Pos {
	return n.New
}

func (n *NewConstructor) End() token.Pos {
	return posAdd(n.Rparen, 1)
}

func (b *BracedNewConstructor) Pos() token.Pos {
	return b.New
}

func (b *BracedNewConstructor) End() token.Pos {
	return nodeEnd(wrapNode(b.Body))
}

func (b *BracedConstructor) Pos() token.Pos {
	return b.Lbrace
}

func (b *BracedConstructor) End() token.Pos {
	return posAdd(b.Rbrace, 1)
}

func (b *BracedConstructorField) Pos() token.Pos {
	return nodePos(wrapNode(b.Name))
}

func (b *BracedConstructorField) End() token.Pos {
	return nodeEnd(wrapNode(b.Value))
}

func (b *BracedConstructorFieldValueExpr) Pos() token.Pos {
	return b.Colon
}

func (b *BracedConstructorFieldValueExpr) End() token.Pos {
	return nodeEnd(wrapNode(b.Expr))
}

func (s *SimpleType) Pos() token.Pos {
	return s.NamePos
}

func (s *SimpleType) End() token.Pos {
	return posAdd(s.NamePos, len(s.Name))
}

func (a *ArrayType) Pos() token.Pos {
	return a.Array
}

func (a *ArrayType) End() token.Pos {
	return posAdd(a.Gt, 1)
}

func (s *StructType) Pos() token.Pos {
	return s.Struct
}

func (s *StructType) End() token.Pos {
	return posAdd(s.Gt, 1)
}

func (s *StructField) Pos() token.Pos {
	return nodePos(nodeChoice(wrapNode(s.Ident), wrapNode(s.Type)))
}

func (s *StructField) End() token.Pos {
	return nodeEnd(wrapNode(s.Type))
}

func (n *NamedType) Pos() token.Pos {
	return nodePos(nodeSliceIndex(n.Path, 0))
}

func (n *NamedType) End() token.Pos {
	return nodeEnd(nodeSliceLast(n.Path))
}

func (c *CastIntValue) Pos() token.Pos {
	return c.Cast
}

func (c *CastIntValue) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (c *CastNumValue) Pos() token.Pos {
	return c.Cast
}

func (c *CastNumValue) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (o *Options) Pos() token.Pos {
	return o.Options
}

func (o *Options) End() token.Pos {
	return posAdd(o.Rparen, 1)
}

func (o *OptionsDef) Pos() token.Pos {
	return nodePos(wrapNode(o.Name))
}

func (o *OptionsDef) End() token.Pos {
	return nodeEnd(wrapNode(o.Value))
}

func (c *CreateSchema) Pos() token.Pos {
	return c.Create
}

func (c *CreateSchema) End() token.Pos {
	return nodeEnd(wrapNode(c.Name))
}

func (d *DropSchema) Pos() token.Pos {
	return d.Drop
}

func (d *DropSchema) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (c *CreateDatabase) Pos() token.Pos {
	return c.Create
}

func (c *CreateDatabase) End() token.Pos {
	return nodeEnd(wrapNode(c.Name))
}

func (a *AlterDatabase) Pos() token.Pos {
	return a.Alter
}

func (a *AlterDatabase) End() token.Pos {
	return nodeEnd(wrapNode(a.Name))
}

func (c *CreatePlacement) Pos() token.Pos {
	return c.Create
}

func (c *CreatePlacement) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(c.Options), wrapNode(c.Name)))
}

func (p *ProtoBundleTypes) Pos() token.Pos {
	return p.Lparen
}

func (p *ProtoBundleTypes) End() token.Pos {
	return posAdd(p.Rparen, 1)
}

func (c *CreateProtoBundle) Pos() token.Pos {
	return c.Create
}

func (c *CreateProtoBundle) End() token.Pos {
	return nodeEnd(wrapNode(c.Types))
}

func (a *AlterProtoBundle) Pos() token.Pos {
	return a.Alter
}

func (a *AlterProtoBundle) End() token.Pos {
	return nodeEnd(wrapNode(a.Alteration))
}

func (a *AlterProtoBundleInsert) Pos() token.Pos {
	return a.Insert
}

func (a *AlterProtoBundleInsert) End() token.Pos {
	return nodeEnd(wrapNode(a.Types))
}

func (a *AlterProtoBundleUpdate) Pos() token.Pos {
	return a.Update
}

func (a *AlterProtoBundleUpdate) End() token.Pos {
	return nodeEnd(wrapNode(a.Types))
}

func (a *AlterProtoBundleDelete) Pos() token.Pos {
	return a.Delete
}

func (a *AlterProtoBundleDelete) End() token.Pos {
	return nodeEnd(wrapNode(a.Types))
}

func (d *DropProtoBundle) Pos() token.Pos {
	return d.Drop
}

func (d *DropProtoBundle) End() token.Pos {
	return posAdd(d.Bundle, 6)
}

func (c *CreateTable) Pos() token.Pos {
	return c.Create
}

func (c *CreateTable) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(c.RowDeletionPolicy)), nodeEnd(wrapNode(c.Cluster)), posAdd(c.Rparen, 1))
}

func (s *Synonym) Pos() token.Pos {
	return s.Synonym
}

func (s *Synonym) End() token.Pos {
	return posAdd(s.Rparen, 1)
}

func (c *CreateSequence) Pos() token.Pos {
	return c.Create
}

func (c *CreateSequence) End() token.Pos {
	return nodeEnd(wrapNode(c.Options))
}

func (c *ColumnDef) Pos() token.Pos {
	return nodePos(wrapNode(c.Name))
}

func (c *ColumnDef) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(c.Options)), posAdd(c.Hidden, 6), nodeEnd(wrapNode(c.GeneratedExpr)), nodeEnd(wrapNode(c.DefaultExpr)), posAdd(c.Null, 4), nodeEnd(wrapNode(c.Type)))
}

func (c *ColumnDefaultExpr) Pos() token.Pos {
	return c.Default
}

func (c *ColumnDefaultExpr) End() token.Pos {
	return c.Rparen
}

func (g *GeneratedColumnExpr) Pos() token.Pos {
	return g.As
}

func (g *GeneratedColumnExpr) End() token.Pos {
	return posChoice(posAdd(g.Stored, 6), posAdd(g.Rparen, 1))
}

func (c *ColumnDefOptions) Pos() token.Pos {
	return c.Options
}

func (c *ColumnDefOptions) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (t *TableConstraint) Pos() token.Pos {
	return posChoice(t.ConstraintPos, nodePos(wrapNode(t.Constraint)))
}

func (t *TableConstraint) End() token.Pos {
	return nodeEnd(wrapNode(t.Constraint))
}

func (f *ForeignKey) Pos() token.Pos {
	return f.Foreign
}

func (f *ForeignKey) End() token.Pos {
	return posChoice(f.OnDeleteEnd, posAdd(f.Rparen, 1))
}

func (c *Check) Pos() token.Pos {
	return c.Check
}

func (c *Check) End() token.Pos {
	return posAdd(c.Rparen, 1)
}

func (i *IndexKey) Pos() token.Pos {
	return nodePos(wrapNode(i.Name))
}

func (i *IndexKey) End() token.Pos {
	return posChoice(posAdd(i.DirPos, len(i.Dir)), nodeEnd(wrapNode(i.Name)))
}

func (c *Cluster) Pos() token.Pos {
	return c.Comma
}

func (c *Cluster) End() token.Pos {
	return posChoice(c.OnDeleteEnd, nodeEnd(wrapNode(c.TableName)))
}

func (c *CreateRowDeletionPolicy) Pos() token.Pos {
	return c.Comma
}

func (c *CreateRowDeletionPolicy) End() token.Pos {
	return nodeEnd(wrapNode(c.RowDeletionPolicy))
}

func (r *RowDeletionPolicy) Pos() token.Pos {
	return r.Row
}

func (r *RowDeletionPolicy) End() token.Pos {
	return posAdd(r.Rparen, 1)
}

func (c *CreateView) Pos() token.Pos {
	return c.Create
}

func (c *CreateView) End() token.Pos {
	return nodeEnd(wrapNode(c.Query))
}

func (d *DropView) Pos() token.Pos {
	return d.Drop
}

func (d *DropView) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (a *AlterTable) Pos() token.Pos {
	return a.Alter
}

func (a *AlterTable) End() token.Pos {
	return nodeEnd(wrapNode(a.TableAlteration))
}

func (a *AlterIndex) Pos() token.Pos {
	return a.Alter
}

func (a *AlterIndex) End() token.Pos {
	return nodeEnd(wrapNode(a.IndexAlteration))
}

func (a *AlterSequence) Pos() token.Pos {
	return a.Alter
}

func (a *AlterSequence) End() token.Pos {
	return nodeEnd(wrapNode(a.Options))
}

func (a *AlterChangeStream) Pos() token.Pos {
	return a.Alter
}

func (a *AlterChangeStream) End() token.Pos {
	return nodeEnd(wrapNode(a.ChangeStreamAlteration))
}

func (a *AddSynonym) Pos() token.Pos {
	return a.Add
}

func (a *AddSynonym) End() token.Pos {
	return nodeEnd(wrapNode(a.Name))
}

func (d *DropSynonym) Pos() token.Pos {
	return d.Drop
}

func (d *DropSynonym) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (r *RenameTo) Pos() token.Pos {
	return r.Rename
}

func (r *RenameTo) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(r.AddSynonym), wrapNode(r.Name)))
}

func (a *AddColumn) Pos() token.Pos {
	return a.Add
}

func (a *AddColumn) End() token.Pos {
	return nodeEnd(wrapNode(a.Column))
}

func (a *AddTableConstraint) Pos() token.Pos {
	return a.Add
}

func (a *AddTableConstraint) End() token.Pos {
	return nodeEnd(wrapNode(a.TableConstraint))
}

func (a *AddRowDeletionPolicy) Pos() token.Pos {
	return a.Add
}

func (a *AddRowDeletionPolicy) End() token.Pos {
	return nodeEnd(wrapNode(a.RowDeletionPolicy))
}

func (d *DropColumn) Pos() token.Pos {
	return d.Drop
}

func (d *DropColumn) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropConstraint) Pos() token.Pos {
	return d.Drop
}

func (d *DropConstraint) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropRowDeletionPolicy) Pos() token.Pos {
	return d.Drop
}

func (d *DropRowDeletionPolicy) End() token.Pos {
	return posAdd(d.Policy, 6)
}

func (r *ReplaceRowDeletionPolicy) Pos() token.Pos {
	return r.Replace
}

func (r *ReplaceRowDeletionPolicy) End() token.Pos {
	return nodeEnd(wrapNode(r.RowDeletionPolicy))
}

func (s *SetOnDelete) Pos() token.Pos {
	return s.Set
}

func (s *SetOnDelete) End() token.Pos {
	return s.OnDeleteEnd
}

func (a *AlterColumn) Pos() token.Pos {
	return a.Alter
}

func (a *AlterColumn) End() token.Pos {
	return nodeEnd(wrapNode(a.Alteration))
}

func (a *AlterColumnType) Pos() token.Pos {
	return nodePos(wrapNode(a.Type))
}

func (a *AlterColumnType) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(a.DefaultExpr)), posAdd(a.Null, 4), nodeEnd(wrapNode(a.Type)))
}

func (a *AlterColumnSetOptions) Pos() token.Pos {
	return a.Set
}

func (a *AlterColumnSetOptions) End() token.Pos {
	return nodeEnd(wrapNode(a.Options))
}

func (a *AlterColumnSetDefault) Pos() token.Pos {
	return a.Set
}

func (a *AlterColumnSetDefault) End() token.Pos {
	return nodeEnd(wrapNode(a.DefaultExpr))
}

func (a *AlterColumnDropDefault) Pos() token.Pos {
	return a.Drop
}

func (a *AlterColumnDropDefault) End() token.Pos {
	return posAdd(a.Default, 7)
}

func (d *DropTable) Pos() token.Pos {
	return d.Drop
}

func (d *DropTable) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (r *RenameTable) Pos() token.Pos {
	return r.Rename
}

func (r *RenameTable) End() token.Pos {
	return nodeEnd(nodeSliceLast(r.Tos))
}

func (r *RenameTableTo) Pos() token.Pos {
	return nodePos(wrapNode(r.Old))
}

func (r *RenameTableTo) End() token.Pos {
	return nodeEnd(wrapNode(r.New))
}

func (c *CreateIndex) Pos() token.Pos {
	return c.Create
}

func (c *CreateIndex) End() token.Pos {
	return posChoice(nodeEnd(nodeChoice(wrapNode(c.InterleaveIn), wrapNode(c.Storing))), posAdd(c.Rparen, 1))
}

func (c *CreateVectorIndex) Pos() token.Pos {
	return c.Create
}

func (c *CreateVectorIndex) End() token.Pos {
	return nodeEnd(wrapNode(c.Options))
}

func (v *VectorIndexOption) Pos() token.Pos {
	return nodePos(wrapNode(v.Key))
}

func (v *VectorIndexOption) End() token.Pos {
	return nodeEnd(wrapNode(v.Value))
}

func (c *CreateChangeStream) Pos() token.Pos {
	return c.Create
}

func (c *CreateChangeStream) End() token.Pos {
	return nodeEnd(nodeChoice(wrapNode(c.Options), wrapNode(c.For), wrapNode(c.Name)))
}

func (c *ChangeStreamForAll) Pos() token.Pos {
	return c.For
}

func (c *ChangeStreamForAll) End() token.Pos {
	return c.All
}

func (c *ChangeStreamForTables) Pos() token.Pos {
	return c.For
}

func (c *ChangeStreamForTables) End() token.Pos {
	return nodeEnd(nodeSliceLast(c.Tables))
}

func (c *ChangeStreamForTable) Pos() token.Pos {
	return nodePos(wrapNode(c.TableName))
}

func (c *ChangeStreamForTable) End() token.Pos {
	return posChoice(nodeEnd(wrapNode(c.TableName)), posAdd(c.Rparen, 1))
}

func (c *ChangeStreamSetFor) Pos() token.Pos {
	return c.Set
}

func (c *ChangeStreamSetFor) End() token.Pos {
	return nodeEnd(wrapNode(c.For))
}

func (c *ChangeStreamDropForAll) Pos() token.Pos {
	return c.Drop
}

func (c *ChangeStreamDropForAll) End() token.Pos {
	return posAdd(c.All, 3)
}

func (c *ChangeStreamSetOptions) Pos() token.Pos {
	return c.Set
}

func (c *ChangeStreamSetOptions) End() token.Pos {
	return nodeEnd(wrapNode(c.Options))
}

func (s *Storing) Pos() token.Pos {
	return s.Storing
}

func (s *Storing) End() token.Pos {
	return posAdd(s.Rparen, 1)
}

func (i *InterleaveIn) Pos() token.Pos {
	return i.Comma
}

func (i *InterleaveIn) End() token.Pos {
	return nodeEnd(wrapNode(i.TableName))
}

func (a *AddStoredColumn) Pos() token.Pos {
	return a.Add
}

func (a *AddStoredColumn) End() token.Pos {
	return nodeEnd(wrapNode(a.Name))
}

func (d *DropStoredColumn) Pos() token.Pos {
	return d.Drop
}

func (d *DropStoredColumn) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropIndex) Pos() token.Pos {
	return d.Drop
}

func (d *DropIndex) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropVectorIndex) Pos() token.Pos {
	return d.Drop
}

func (d *DropVectorIndex) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropSequence) Pos() token.Pos {
	return d.Drop
}

func (d *DropSequence) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (c *CreateRole) Pos() token.Pos {
	return c.Create
}

func (c *CreateRole) End() token.Pos {
	return nodeEnd(wrapNode(c.Name))
}

func (d *DropRole) Pos() token.Pos {
	return d.Drop
}

func (d *DropRole) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (d *DropChangeStream) Pos() token.Pos {
	return d.Drop
}

func (d *DropChangeStream) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (g *Grant) Pos() token.Pos {
	return g.Grant
}

func (g *Grant) End() token.Pos {
	return nodeEnd(nodeSliceLast(g.Roles))
}

func (r *Revoke) Pos() token.Pos {
	return r.Revoke
}

func (r *Revoke) End() token.Pos {
	return nodeEnd(nodeSliceLast(r.Roles))
}

func (p *PrivilegeOnTable) Pos() token.Pos {
	return nodePos(nodeSliceIndex(p.Privileges, 0))
}

func (p *PrivilegeOnTable) End() token.Pos {
	return nodeEnd(nodeSliceLast(p.Names))
}

func (s *SelectPrivilege) Pos() token.Pos {
	return s.Select
}

func (s *SelectPrivilege) End() token.Pos {
	return posChoice(posAdd(s.Rparen, 1), posAdd(s.Select, 6))
}

func (i *InsertPrivilege) Pos() token.Pos {
	return i.Insert
}

func (i *InsertPrivilege) End() token.Pos {
	return posChoice(posAdd(i.Rparen, 1), posAdd(i.Insert, 6))
}

func (u *UpdatePrivilege) Pos() token.Pos {
	return u.Update
}

func (u *UpdatePrivilege) End() token.Pos {
	return posChoice(posAdd(u.Rparen, 1), posAdd(u.Update, 6))
}

func (d *DeletePrivilege) Pos() token.Pos {
	return d.Delete
}

func (d *DeletePrivilege) End() token.Pos {
	return posAdd(d.Delete, 6)
}

func (s *SelectPrivilegeOnChangeStream) Pos() token.Pos {
	return s.Select
}

func (s *SelectPrivilegeOnChangeStream) End() token.Pos {
	return nodeEnd(nodeSliceLast(s.Names))
}

func (s *SelectPrivilegeOnView) Pos() token.Pos {
	return s.Select
}

func (s *SelectPrivilegeOnView) End() token.Pos {
	return nodeEnd(nodeSliceLast(s.Names))
}

func (e *ExecutePrivilegeOnTableFunction) Pos() token.Pos {
	return e.Execute
}

func (e *ExecutePrivilegeOnTableFunction) End() token.Pos {
	return nodeEnd(nodeSliceLast(e.Names))
}

func (r *RolePrivilege) Pos() token.Pos {
	return r.Role
}

func (r *RolePrivilege) End() token.Pos {
	return nodeEnd(nodeSliceLast(r.Names))
}

func (a *AlterStatistics) Pos() token.Pos {
	return a.Alter
}

func (a *AlterStatistics) End() token.Pos {
	return nodeEnd(wrapNode(a.Options))
}

func (a *Analyze) Pos() token.Pos {
	return a.Analyze
}

func (a *Analyze) End() token.Pos {
	return posAdd(a.Analyze, 7)
}

func (s *ScalarSchemaType) Pos() token.Pos {
	return s.NamePos
}

func (s *ScalarSchemaType) End() token.Pos {
	return posAdd(s.NamePos, len(s.Name))
}

func (s *SizedSchemaType) Pos() token.Pos {
	return s.NamePos
}

func (s *SizedSchemaType) End() token.Pos {
	return posAdd(s.Rparen, 1)
}

func (a *ArraySchemaType) Pos() token.Pos {
	return a.Array
}

func (a *ArraySchemaType) End() token.Pos {
	return posAdd(a.Gt, 1)
}

func (c *CreateSearchIndex) Pos() token.Pos {
	return c.Create
}

func (c *CreateSearchIndex) End() token.Pos {
	return posChoice(nodeEnd(nodeChoice(wrapNode(c.Options), wrapNode(c.Interleave), wrapNode(c.Where), wrapNode(c.OrderBy), nodeSliceLast(c.PartitionColumns), wrapNode(c.Storing))), posAdd(c.Rparen, 1))
}

func (d *DropSearchIndex) Pos() token.Pos {
	return d.Drop
}

func (d *DropSearchIndex) End() token.Pos {
	return nodeEnd(wrapNode(d.Name))
}

func (a *AlterSearchIndex) Pos() token.Pos {
	return a.Alter
}

func (a *AlterSearchIndex) End() token.Pos {
	return nodeEnd(wrapNode(a.IndexAlteration))
}

func (i *Insert) Pos() token.Pos {
	return i.Insert
}

func (i *Insert) End() token.Pos {
	return nodeEnd(wrapNode(i.Input))
}

func (v *ValuesInput) Pos() token.Pos {
	return v.Values
}

func (v *ValuesInput) End() token.Pos {
	return nodeEnd(nodeSliceLast(v.Rows))
}

func (v *ValuesRow) Pos() token.Pos {
	return v.Lparen
}

func (v *ValuesRow) End() token.Pos {
	return posAdd(v.Rparen, 1)
}

func (d *DefaultExpr) Pos() token.Pos {
	return posChoice(d.DefaultPos, nodePos(wrapNode(d.Expr)))
}

func (d *DefaultExpr) End() token.Pos {
	return posChoice(posAdd(d.DefaultPos, 7), nodeEnd(wrapNode(d.Expr)))
}

func (s *SubQueryInput) Pos() token.Pos {
	return nodePos(wrapNode(s.Query))
}

func (s *SubQueryInput) End() token.Pos {
	return nodeEnd(wrapNode(s.Query))
}

func (d *Delete) Pos() token.Pos {
	return d.Delete
}

func (d *Delete) End() token.Pos {
	return nodeEnd(wrapNode(d.Where))
}

func (u *Update) Pos() token.Pos {
	return u.Update
}

func (u *Update) End() token.Pos {
	return nodeEnd(wrapNode(u.Where))
}

func (u *UpdateItem) Pos() token.Pos {
	return nodePos(nodeSliceIndex(u.Path, 0))
}

func (u *UpdateItem) End() token.Pos {
	return nodeEnd(wrapNode(u.DefaultExpr))
}
