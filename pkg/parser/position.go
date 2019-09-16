package parser

import (
	"bytes"
	"fmt"
	"strings"
)

// Pos is just source code position.
//
// Internally it is zero-origin offset of the buffer.
type Pos int

const InvalidPos Pos = -1

// Invalid returns whether p is invalid Pos or not.
func (p Pos) Invalid() bool {
	return p < 0
}

// Position is source code position with file path
// and source code around this position.
type Position struct {
	FilePath string
	Pos, End Pos

	// Line and Column are 0-origin.
	Line, Column       int
	EndLine, EndColumn int

	// Source is source code around this position with line number and cursor
	// for detailed error message.
	Source string
}

func (pos *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", pos.FilePath, pos.Line+1, pos.Column+1)
}

// File is input file with source code.
type File struct {
	FilePath string
	Buffer   string

	lines []Pos
}

// Position returns a new Position from pos and end on this File.
func (f *File) Position(pos, end Pos) *Position {
	line, column := f.ResovlePos(pos)
	endLine, endColumn := f.ResovlePos(end)

	// Calculate source coude around this position.
	var source bytes.Buffer
	switch {
	case pos.Invalid() || end.Invalid():
		break
	case line == endLine:
		lineBuffer := f.Buffer[f.lines[line] : f.lines[line+1]-1]
		count := endColumn - column - 1
		if count < 0 {
			count = 0
		}
		fmt.Fprintf(&source, "%3d:  %s\n", line+1, lineBuffer)
		fmt.Fprintf(&source, "      %s^%s\n", strings.Repeat(" ", column), strings.Repeat("~", count))
	case line < endLine:
		for l := line; l <= endLine; l++ {
			lineBuffer := f.Buffer[f.lines[l] : f.lines[l+1]-1]
			fmt.Fprintf(&source, "%3d:  %s\n", l+1, lineBuffer)
		}
	}

	return &Position{
		FilePath:  f.FilePath,
		Pos:       pos,
		End:       end,
		Line:      line,
		Column:    column,
		EndLine:   endLine,
		EndColumn: endColumn,
		Source:    source.String(),
	}
}

// ResolvePos returns line and column number from pos.
func (f *File) ResovlePos(pos Pos) (line int, column int) {
	line, column = -1, -1

	if pos.Invalid() {
		return
	}

	f.init()
	// TODO: for performance, use binary search instead
	for line = len(f.lines) - 1; line >= 0; line-- {
		linePos := f.lines[line]
		if linePos <= pos {
			column = int(pos - linePos)
			return
		}
	}

	return
}

// init initialize f.lines.
func (f *File) init() {
	if f.lines != nil {
		return
	}

	lines := []Pos{0}
	for i, line := range strings.Split(f.Buffer, "\n") {
		lines = append(lines, Pos(int(lines[i])+len(line)+1))
	}
	f.lines = lines
}

// ================================================================================
//
// SELECT
//
// ================================================================================

func (q *QueryStatement) Pos() Pos {
	if q.Hint != nil {
		return q.Hint.Pos()
	}
	return q.Query.Pos()
}

func (q *QueryStatement) End() Pos {
	return q.Query.End()
}

func (h *Hint) Pos() Pos { return h.pos }
func (h *Hint) End() Pos { return h.end }

func (h *HintRecord) Pos() Pos { return h.Key.Pos() }
func (h *HintRecord) End() Pos { return h.Value.End() }

func (s *Select) Pos() Pos { return s.pos }

func (s *Select) End() Pos {
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

func (c *CompoundQuery) Pos() Pos {
	return c.Queries[0].Pos()
}

func (c *CompoundQuery) End() Pos {
	if c.Limit != nil {
		return c.Limit.End()
	}
	if c.OrderBy != nil {
		return c.OrderBy.End()
	}
	return c.Queries[len(c.Queries)-1].End()
}

func (s *SubQuery) Pos() Pos { return s.pos }
func (s *SubQuery) End() Pos { return s.end }

func (s *Star) Pos() Pos { return s.pos }
func (s *Star) End() Pos { return s.pos + 1 }

func (s *DotStar) Pos() Pos { return s.Expr.Pos() }
func (s *DotStar) End() Pos { return s.end }

func (a *Alias) Pos() Pos { return a.Expr.Pos() }
func (a *Alias) End() Pos { return a.As.End() }

func (a *AsAlias) Pos() Pos { return a.pos }
func (a *AsAlias) End() Pos { return a.Alias.End() }

func (e *ExprSelectItem) Pos() Pos { return e.Expr.Pos() }
func (e *ExprSelectItem) End() Pos { return e.Expr.End() }

func (f *From) Pos() Pos { return f.pos }
func (f *From) End() Pos { return f.Source.End() }

func (w *Where) Pos() Pos { return w.pos }
func (w *Where) End() Pos { return w.Expr.End() }

func (g *GroupBy) Pos() Pos { return g.pos }
func (g *GroupBy) End() Pos { return g.Exprs[len(g.Exprs)-1].End() }

func (h *Having) Pos() Pos { return h.pos }
func (h *Having) End() Pos { return h.Expr.End() }

func (o *OrderBy) Pos() Pos { return o.pos }
func (o *OrderBy) End() Pos { return o.Items[len(o.Items)-1].End() }

func (o *OrderByItem) Pos() Pos { return o.Expr.Pos() }

func (o *OrderByItem) End() Pos { return o.end }

func (c *Collate) Pos() Pos { return c.pos }
func (c *Collate) End() Pos { return c.Value.End() }

func (l *Limit) Pos() Pos {
	return l.pos
}

func (l *Limit) End() Pos {
	if l.Offset != nil {
		return l.Offset.End()
	}
	return l.Count.End()
}

func (o *Offset) Pos() Pos { return o.pos }
func (o *Offset) End() Pos { return o.Value.End() }

// ================================================================================
//
// JOIN
//
// ================================================================================

func (u *Unnest) Pos() Pos { return u.pos }
func (u *Unnest) End() Pos { return u.end }

func (w *WithOffset) Pos() Pos { return w.pos }
func (w *WithOffset) End() Pos { return w.end }

func (t *TableName) Pos() Pos {
	return t.Table.Pos()
}

func (t *TableName) End() Pos {
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

func (s *SubQueryTableExpr) Pos() Pos {
	return s.Query.Pos()
}

func (s *SubQueryTableExpr) End() Pos {
	if s.Sample != nil {
		return s.Sample.End()
	}
	if s.As != nil {
		return s.As.End()
	}
	return s.Query.End()
}

func (p *ParenTableExpr) Pos() Pos { return p.pos }
func (p *ParenTableExpr) End() Pos { return p.end }

func (j *Join) Pos() Pos {
	return j.Left.Pos()
}

func (j *Join) End() Pos {
	if j.Cond != nil {
		return j.Cond.End()
	}
	return j.Right.End()
}

func (o *On) Pos() Pos { return o.pos }
func (o *On) End() Pos { return o.Expr.End() }

func (u *Using) Pos() Pos { return u.pos }
func (u *Using) End() Pos { return u.end }

func (t *TableSample) Pos() Pos { return t.pos }
func (t *TableSample) End() Pos { return t.Size.End() }

func (t *TableSampleSize) Pos() Pos { return t.pos }
func (t *TableSampleSize) End() Pos { return t.end }

// ================================================================================
//
// Expr
//
// ================================================================================

func (b *BinaryExpr) Pos() Pos { return b.Left.Pos() }
func (b *BinaryExpr) End() Pos { return b.Right.End() }

func (u *UnaryExpr) Pos() Pos { return u.pos }
func (u *UnaryExpr) End() Pos { return u.Expr.End() }

func (i *InExpr) Pos() Pos { return i.Left.Pos() }
func (i *InExpr) End() Pos { return i.Right.End() }

func (u *UnnestInCondition) Pos() Pos { return u.pos }
func (u *UnnestInCondition) End() Pos { return u.end }

func (s *SubQueryInCondition) Pos() Pos { return s.pos }
func (s *SubQueryInCondition) End() Pos { return s.end }

func (v *ValuesInCondition) Pos() Pos { return v.pos }
func (v *ValuesInCondition) End() Pos { return v.end }

func (i *IsNullExpr) Pos() Pos { return i.Left.Pos() }
func (i *IsNullExpr) End() Pos { return i.end }

func (i *IsBoolExpr) Pos() Pos { return i.Left.Pos() }
func (i *IsBoolExpr) End() Pos { return i.end }

func (b *BetweenExpr) Pos() Pos { return b.Left.Pos() }
func (b *BetweenExpr) End() Pos { return b.RightEnd.End() }

func (b *SelectorExpr) Pos() Pos { return b.Expr.Pos() }
func (b *SelectorExpr) End() Pos { return b.Member.End() }

func (i *IndexExpr) Pos() Pos { return i.Expr.Pos() }
func (i *IndexExpr) End() Pos { return i.end }

func (c *CallExpr) Pos() Pos { return c.Func.Pos() }
func (c *CallExpr) End() Pos { return c.end }

func (a *Arg) Pos() Pos {
	return a.pos
}

func (a *Arg) End() Pos {
	if a.IntervalUnit != nil {
		return a.IntervalUnit.End()
	}
	return a.Expr.End()
}

func (c *CountStarExpr) Pos() Pos { return c.pos }
func (c *CountStarExpr) End() Pos { return c.end }

func (e *ExtractExpr) Pos() Pos { return e.pos }
func (e *ExtractExpr) End() Pos { return e.end }

func (a *AtTimeZone) Pos() Pos { return a.pos }
func (a *AtTimeZone) End() Pos { return a.Expr.End() }

func (c *CastExpr) Pos() Pos { return c.pos }
func (c *CastExpr) End() Pos { return c.end }

func (c *CaseExpr) Pos() Pos { return c.pos }
func (c *CaseExpr) End() Pos { return c.end }

func (c *CaseWhen) Pos() Pos { return c.Cond.Pos() }
func (c *CaseWhen) End() Pos { return c.Then.End() }

func (c *CaseElse) Pos() Pos { return c.pos }
func (c *CaseElse) End() Pos { return c.Expr.End() }

func (p *ParenExpr) Pos() Pos { return p.pos }
func (p *ParenExpr) End() Pos { return p.end }

func (s *ScalarSubQuery) Pos() Pos { return s.pos }
func (s *ScalarSubQuery) End() Pos { return s.end }

func (a *ArraySubQuery) Pos() Pos { return a.pos }
func (a *ArraySubQuery) End() Pos { return a.end }

func (e *ExistsSubQuery) Pos() Pos { return e.pos }
func (e *ExistsSubQuery) End() Pos { return e.end }

// ================================================================================
//
// Literal
//
// ================================================================================

func (p *Param) Pos() Pos { return p.pos }
func (p *Param) End() Pos { return p.pos + 1 + Pos(len(p.Name)) }

func (i *Ident) Pos() Pos { return i.pos }
func (i *Ident) End() Pos { return i.end }

func (p *Path) Pos() Pos { return p.Idents[0].Pos() }
func (p *Path) End() Pos { return p.Idents[len(p.Idents)-1].End() }

func (a *ArrayLiteral) Pos() Pos { return a.pos }
func (a *ArrayLiteral) End() Pos { return a.end }

func (s *StructLiteral) Pos() Pos { return s.pos }
func (s *StructLiteral) End() Pos { return s.end }

func (n *NullLiteral) Pos() Pos { return n.pos }
func (n *NullLiteral) End() Pos { return n.pos + 4 }

func (b *BoolLiteral) Pos() Pos {
	return b.pos
}

func (b *BoolLiteral) End() Pos {
	if b.Value {
		return b.pos + 4
	} else {
		return b.pos + 5
	}
}

func (i *IntLiteral) Pos() Pos { return i.pos }
func (i *IntLiteral) End() Pos { return i.end }

func (f *FloatLiteral) Pos() Pos { return f.pos }
func (f *FloatLiteral) End() Pos { return f.end }

func (s *StringLiteral) Pos() Pos { return s.pos }
func (s *StringLiteral) End() Pos { return s.end }

func (b *BytesLiteral) Pos() Pos { return b.pos }
func (b *BytesLiteral) End() Pos { return b.end }

func (d *DateLiteral) Pos() Pos { return d.pos }
func (d *DateLiteral) End() Pos { return d.end }

func (t *TimestampLiteral) Pos() Pos { return t.pos }
func (t *TimestampLiteral) End() Pos { return t.end }

// ================================================================================
//
// Type
//
// ================================================================================

func (s *SimpleType) Pos() Pos { return s.pos }
func (s *SimpleType) End() Pos { return s.pos + Pos(len(s.Name)) }

func (a *ArrayType) Pos() Pos { return a.pos }
func (a *ArrayType) End() Pos { return a.end }

func (s *StructType) Pos() Pos { return s.pos }
func (s *StructType) End() Pos { return s.end }

func (f *FieldType) Pos() Pos {
	if f.Member != nil {
		return f.Member.Pos()
	}
	return f.Type.Pos()
}

func (f *FieldType) End() Pos {
	return f.Type.End()
}

// ================================================================================
//
// Cast for Special Cases
//
// ================================================================================

func (c *CastIntValue) Pos() Pos { return c.pos }
func (c *CastIntValue) End() Pos { return c.end }

func (c *CastNumValue) Pos() Pos { return c.pos }
func (c *CastNumValue) End() Pos { return c.end }

// ================================================================================
//
// DDL
//
// ================================================================================

func (c *CreateDatabase) Pos() Pos { return c.pos }
func (c *CreateDatabase) End() Pos { return c.Name.End() }

func (c *CreateTable) Pos() Pos { return c.pos }
func (c *CreateTable) End() Pos { return c.end }

func (c *ColumnDef) Pos() Pos { return c.Name.Pos() }
func (c *ColumnDef) End() Pos { return c.end }

func (c *ColumnDefOptions) Pos() Pos { return c.pos }
func (c *ColumnDefOptions) End() Pos { return c.end }

func (i *IndexKey) Pos() Pos { return i.Name.Pos() }
func (i *IndexKey) End() Pos { return i.end }

func (c *Cluster) Pos() Pos { return c.pos }
func (c *Cluster) End() Pos { return c.end }

func (a *AlterTable) Pos() Pos { return a.pos }
func (a *AlterTable) End() Pos { return a.TableAlternation.End() }

func (a *AddColumn) Pos() Pos { return a.pos }
func (a *AddColumn) End() Pos { return a.Column.End() }

func (d *DropColumn) Pos() Pos { return d.pos }
func (d *DropColumn) End() Pos { return d.Name.End() }

func (s *SetOnDelete) Pos() Pos { return s.pos }
func (s *SetOnDelete) End() Pos { return s.end }

func (a *AlterColumn) Pos() Pos { return a.pos }
func (a *AlterColumn) End() Pos { return a.end }

func (a *AlterColumnSet) Pos() Pos { return a.pos }
func (a *AlterColumnSet) End() Pos { return a.Options.End() }

func (d *DropTable) Pos() Pos { return d.pos }
func (d *DropTable) End() Pos { return d.Name.End() }

func (c *CreateIndex) Pos() Pos { return c.pos }
func (c *CreateIndex) End() Pos { return c.end }

func (s *Storing) Pos() Pos { return s.pos }
func (s *Storing) End() Pos { return s.end }

func (i *InterleaveIn) Pos() Pos { return i.pos }
func (i *InterleaveIn) End() Pos { return i.TableName.End() }

func (d *DropIndex) Pos() Pos { return d.pos }
func (d *DropIndex) End() Pos { return d.Name.End() }

// ================================================================================
//
// Types for Schema
//
// ================================================================================

func (s *ScalarSchemaType) Pos() Pos { return s.pos }
func (s *ScalarSchemaType) End() Pos { return s.pos + Pos(len(s.Name)) }

func (s *SizedSchemaType) Pos() Pos { return s.pos }
func (s *SizedSchemaType) End() Pos { return s.end }

func (a *ArraySchemaType) Pos() Pos { return a.pos }
func (a *ArraySchemaType) End() Pos { return a.end }

// ================================================================================
//
// DML
//
// ================================================================================

func (i *Insert) Pos() Pos { return i.pos }
func (i *Insert) End() Pos { return i.Input.End() }

func (v *ValuesInput) Pos() Pos { return v.pos }
func (v *ValuesInput) End() Pos { return v.Rows[len(v.Rows)-1].End() }

func (v *ValuesRow) Pos() Pos { return v.pos }
func (v *ValuesRow) End() Pos { return v.end }

func (d *DefaultExpr) Pos() Pos {
	return d.pos
}

func (d *DefaultExpr) End() Pos {
	if d.Default {
		return d.pos + 7
	} else {
		return d.Expr.End()
	}
}

func (s *SubQueryInput) Pos() Pos { return s.Query.Pos() }
func (s *SubQueryInput) End() Pos { return s.Query.End() }

func (d *Delete) Pos() Pos { return d.pos }
func (d *Delete) End() Pos { return d.Where.End() }

func (u *Update) Pos() Pos { return u.pos }
func (u *Update) End() Pos { return u.Where.End() }

func (u *UpdateItem) Pos() Pos { return u.Path[0].Pos() }
func (u *UpdateItem) End() Pos { return u.Expr.End() }
