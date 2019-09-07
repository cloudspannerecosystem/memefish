package parser

import (
	"fmt"
	"strings"
)

type Pos int

const InvalidPos Pos = -1

type Position struct {
	FilePath string
	Pos, End Pos
	// Line and Column are 0-origin.
	Line, Column       int
	EndLine, EndColumn int
}

func (loc *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", loc.FilePath, loc.Line+1, loc.Column+1)
}

type File struct {
	FilePath string
	Buffer   string
	lines    []Pos
}

func NewFile(filePath string, buffer string) *File {
	lines := []Pos{0}
	for i, line := range strings.Split(buffer, "\n") {
		lines = append(lines, Pos(int(lines[i])+len(line)+1))
	}
	return &File{
		FilePath: filePath,
		Buffer:   buffer,
		lines:    lines,
	}
}

func (f *File) Position(pos, end Pos) *Position {
	line, column := f.ResovlePos(pos)
	endLine, endColumn := f.ResovlePos(end)
	return &Position{
		FilePath:  f.FilePath,
		Pos:       pos,
		End:       end,
		Line:      line,
		Column:    column,
		EndLine:   endLine,
		EndColumn: endColumn,
	}
}

func (f *File) ResovlePos(pos Pos) (int, int) {
	// TODO: for performance, use binary search instead
	for line := len(f.lines) - 1; line >= 0; line-- {
		linePos := f.lines[line]
		if linePos <= pos {
			return line, int(pos - linePos)
		}
	}
	return -1, -1
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

func (s *StarPath) Pos() Pos { return s.Expr.Pos() }
func (s *StarPath) End() Pos { return s.end }

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

func (s *SubQueryJoinExpr) Pos() Pos {
	return s.Query.Pos()
}

func (s *SubQueryJoinExpr) End() Pos {
	if s.Sample != nil {
		return s.Sample.End()
	}
	if s.As != nil {
		return s.As.End()
	}
	return s.Query.End()
}

func (p *ParenJoinExpr) Pos() Pos { return p.pos }
func (p *ParenJoinExpr) End() Pos { return p.end }

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
