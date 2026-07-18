package ast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cloudspannerecosystem/memefish/token"
	"github.com/cloudspannerecosystem/memefish/tools/util/astcatalog"
)

func TestExprPrecCoversAllExprImplementations(t *testing.T) {
	for _, tc := range exprPrecTestCases {
		t.Run(tc.name, func(t *testing.T) {
			requireNotPanic(t, func() {
				_ = exprPrec(tc.expr)
			})
		})
	}
}

func TestExprPrecCoversAllBinaryOpValues(t *testing.T) {
	values := constValues(t, "BinaryOp")
	for _, v := range values {
		t.Run(v.Name, func(t *testing.T) {
			requireNotPanic(t, func() {
				_ = exprPrec(&BinaryExpr{Op: BinaryOp(v.Value)})
			})
		})
	}
}

func TestExprPrecCoversAllUnaryOpValues(t *testing.T) {
	values := constValues(t, "UnaryOp")
	for _, v := range values {
		t.Run(v.Name, func(t *testing.T) {
			requireNotPanic(t, func() {
				_ = exprPrec(&UnaryExpr{Op: UnaryOp(v.Value)})
			})
		})
	}
}

type unexpectedExpr struct{}

func (*unexpectedExpr) Pos() token.Pos { return token.InvalidPos }
func (*unexpectedExpr) End() token.Pos { return token.InvalidPos }
func (*unexpectedExpr) SQL() string    { return "unexpected" }
func (*unexpectedExpr) isExpr()        {}

func TestExprPrecPanicIncludesUnhandledType(t *testing.T) {
	var got any
	func() {
		defer func() {
			got = recover()
		}()

		_ = exprPrec(&unexpectedExpr{})
	}()

	if got == nil {
		t.Fatal("exprPrec did not panic")
	}

	const want = "exprPrec: unhandled expr type: *ast.unexpectedExpr"
	if !strings.Contains(fmt.Sprint(got), want) {
		t.Fatalf("panic = %v, want to contain %q", got, want)
	}
}

func TestGQLLabelExpressionSQLPreservesPrecedence(t *testing.T) {
	label := func(name string) GQLLabelExpression {
		return &GQLNameLabel{Name: &Ident{Name: name}}
	}

	tests := []struct {
		name string
		expr GQLLabelExpression
		want string
	}{
		{
			name: "and over or",
			expr: &GQLLabelBinaryExpr{
				Op: GQLLabelOpAnd,
				Left: &GQLLabelBinaryExpr{
					Op:    GQLLabelOpOr,
					Left:  label("A"),
					Right: label("B"),
				},
				Right: label("C"),
			},
			want: "(A | B) & C",
		},
		{
			name: "not over or",
			expr: &GQLLabelUnaryExpr{
				Op: GQLLabelOpNot,
				Expr: &GQLLabelBinaryExpr{
					Op:    GQLLabelOpOr,
					Left:  label("A"),
					Right: label("B"),
				},
			},
			want: "!(A | B)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.expr.SQL(); got != test.want {
				t.Fatalf("SQL() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGQLPathSearchPrefixSQLRequiresShortestCount(t *testing.T) {
	counted := &GQLPathSearchPrefix{
		Prefix: GQLSearchPrefixShortest,
		Count:  &IntLiteral{Base: 10, Value: "3"},
	}
	if got, want := counted.SQL(), "SHORTEST 3"; got != want {
		t.Fatalf("SQL() = %q, want %q", got, want)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("SQL() did not panic for SHORTEST without Count")
		}
	}()
	_ = (&GQLPathSearchPrefix{Prefix: GQLSearchPrefixShortest}).SQL()
}

func TestGQLPathModeClauseSQLPreservesSuffix(t *testing.T) {
	tests := []struct {
		name   string
		suffix GQLPathModeSuffix
		want   string
	}{
		{name: "without suffix", want: "WALK"},
		{name: "singular path", suffix: GQLPathModeSuffixPath, want: "WALK PATH"},
		{name: "plural paths", suffix: GQLPathModeSuffixPaths, want: "WALK PATHS"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			n := &GQLPathModeClause{Mode: GQLPathModeWalk, Suffix: test.suffix}
			if got := n.SQL(); got != test.want {
				t.Errorf("SQL() = %q, want %q", got, test.want)
			}
		})
	}
}

func loadCatalog(t *testing.T) *astcatalog.Catalog {
	t.Helper()

	catalog, err := astcatalog.Load("ast.go", "ast_const.go")
	if err != nil {
		t.Fatal(err)
	}
	return catalog
}

func constValues(t *testing.T, name string) []*astcatalog.ConstValueDef {
	t.Helper()

	def := loadCatalog(t).Consts[astcatalog.ConstType(name)]
	if def == nil {
		t.Fatalf("%s is not registered in AST catalog", name)
		return nil
	}
	return def.Values
}

func requireNotPanic(t *testing.T, f func()) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("exprPrec panicked: %v", r)
		}
	}()

	f()
}
