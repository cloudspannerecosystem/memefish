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

func (q *QueryStatement) Pos() Pos {
	if q.Hint != nil {
		return q.Hint.Pos()
	}
	return q.Expr.Pos()
}

func (q *QueryStatement) End() Pos {
	return q.Expr.End()
}

func (h *Hint) Pos() Pos { return h.pos }
func (h *Hint) End() Pos { return h.end }

func (s *Select) Pos() Pos { return s.pos }
func (s *Select) End() Pos { return s.end }

func (c *CompoundQuery) Pos() Pos {
	return c.Left.Pos()
}

func (c *CompoundQuery) End() Pos {
	return c.Right.End()
}

func (s *SelectExpr) Pos() Pos { return s.pos }
func (s *SelectExpr) End() Pos { return s.end }

func (s SelectExprList) Pos() Pos {
	return s[0].Pos()
}

func (s SelectExprList) End() Pos {
	return s[len(s)-1].End()
}

func (f *FromItem) Pos() Pos {
	return f.Expr.Pos()
}
func (f *FromItem) End() Pos { return f.end }

func (f FromItemList) Pos() Pos {
	return f[0].Pos()
}

func (f FromItemList) End() Pos {
	return f[len(f)-1].End()
}

func (t *TableName) Pos() Pos {
	return t.Ident.Pos()
}

func (t *TableName) End() Pos { return t.end }

func (u *Unnest) Pos() Pos { return u.pos }
func (u *Unnest) End() Pos { return u.end }

func (s *SubQueryJoinExpr) Pos() Pos {
	return s.Expr.Pos()
}

func (s *SubQueryJoinExpr) End() Pos { return s.end }

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

func (j *JoinCondition) Pos() Pos { return j.pos }
func (j *JoinCondition) End() Pos { return j.end }

func (o *OrderExpr) Pos() Pos {
	return o.Expr.Pos()
}

func (o *OrderExpr) End() Pos { return o.end }

func (o OrderExprList) Pos() Pos {
	return o[0].Pos()
}

func (o OrderExprList) End() Pos {
	return o[len(o)-1].End()
}

func (l *Limit) Pos() Pos {
	return l.Count.Pos()
}

func (l *Limit) End() Pos {
	if l.Offset != nil {
		return l.Offset.End()
	}
	return l.Count.End()
}

func (b *BinaryExpr) Pos() Pos {
	return b.Left.Pos()
}

func (b *BinaryExpr) End() Pos {
	return b.Right.End()
}

func (u *UnaryExpr) Pos() Pos { return u.pos }

func (u *UnaryExpr) End() Pos {
	return u.Expr.End()
}

func (i *InExpr) Pos() Pos {
	return i.Left.Pos()
}

func (i *InExpr) End() Pos {
	return i.Right.End()
}

func (i *InCondition) Pos() Pos { return i.pos }
func (i *InCondition) End() Pos { return i.end }

func (i *IsNullExpr) Pos() Pos {
	return i.Left.Pos()
}

func (i *IsNullExpr) End() Pos { return i.end }

func (i *IsBoolExpr) Pos() Pos {
	return i.Left.Pos()
}

func (i *IsBoolExpr) End() Pos { return i.end }

func (b *BetweenExpr) Pos() Pos {
	return b.Left.Pos()
}

func (b *BetweenExpr) End() Pos {
	return b.RightEnd.End()
}

func (b *SelectorExpr) Pos() Pos {
	return b.Left.Pos()
}

func (b *SelectorExpr) End() Pos {
	return b.Right.End()
}

func (i *IndexExpr) Pos() Pos {
	return i.Left.Pos()
}

func (i *IndexExpr) End() Pos { return i.end }

func (c *CallExpr) Pos() Pos {
	return c.Func.Pos()
}

func (c *CallExpr) End() Pos { return c.end }

func (c *CastExpr) Pos() Pos { return c.pos }
func (c *CastExpr) End() Pos { return c.end }

func (c *CaseExpr) Pos() Pos { return c.pos }
func (c *CaseExpr) End() Pos { return c.end }

func (w *When) Pos() Pos {
	return w.Cond.Pos()
}

func (w *When) End() Pos {
	return w.Then.End()
}

func (s *SubQuery) Pos() Pos { return s.pos }
func (s *SubQuery) End() Pos { return s.end }

func (p *ParenExpr) Pos() Pos { return p.pos }
func (p *ParenExpr) End() Pos { return p.end }

func (a *ArrayExpr) Pos() Pos { return a.pos }

func (a *ArrayExpr) End() Pos {
	return a.Expr.End()
}

func (e *ExistsExpr) Pos() Pos { return e.pos }

func (e *ExistsExpr) End() Pos {
	return e.Expr.Pos()
}

func (p *Param) Pos() Pos { return p.pos }

func (p *Param) End() Pos {
	return p.pos + 1 + Pos(len(p.Name))
}

func (i *Ident) Pos() Pos { return i.pos }

func (i *Ident) End() Pos { return i.end }

func (i IdentList) Pos() Pos {
	return i[0].Pos()
}

func (i IdentList) End() Pos {
	return i[len(i)-1].End()
}

func (a *ArrayLit) Pos() Pos { return a.pos }
func (a *ArrayLit) End() Pos { return a.end }

func (s *StructLit) Pos() Pos { return s.pos }
func (s *StructLit) End() Pos { return s.end }

func (n *NullLit) Pos() Pos { return n.pos }

func (n *NullLit) End() Pos {
	return n.pos + 4
}

func (b *BoolLit) Pos() Pos { return b.pos }

func (b *BoolLit) End() Pos {
	if b.Value {
		return b.pos + 4
	} else {
		return b.pos + 5
	}
}

func (i *IntLit) Pos() Pos { return i.pos }
func (i *IntLit) End() Pos { return i.end }

func (f *FloatLit) Pos() Pos { return f.pos }
func (f *FloatLit) End() Pos { return f.end }

func (s *StringLit) Pos() Pos { return s.pos }
func (s *StringLit) End() Pos { return s.end }

func (b *BytesLit) Pos() Pos { return b.pos }
func (b *BytesLit) End() Pos { return b.end }

func (d *DateLit) Pos() Pos { return d.pos }
func (d *DateLit) End() Pos { return d.end }

func (t *TimestampLit) Pos() Pos { return t.pos }
func (t *TimestampLit) End() Pos { return t.end }

func (t *Type) Pos() Pos { return t.pos }
func (t *Type) End() Pos { return t.end }
