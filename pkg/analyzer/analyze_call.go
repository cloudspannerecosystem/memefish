package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/ast"
)

type functionAnalyzer func(a *Analyzer, e *ast.CallExpr) *TypeInfo

var functionAnalyzers map[string]functionAnalyzer

func init() {
	functionAnalyzers = map[string]functionAnalyzer{
		"CONCAT": concatAnalyzer,
		"SUM":    sumAnalyzer,
	}
}

func (a *Analyzer) analyzeCallExpr(e *ast.CallExpr) *TypeInfo {
	if ana, ok := functionAnalyzers[strings.ToUpper(e.Func.Name)]; ok {
		return ana(a, e)
	}

	panic(a.errorf(e, "unknown function: %s", e.Func.SQL()))
}

func concatAnalyzer(a *Analyzer, e *ast.CallExpr) *TypeInfo {
	if e.Distinct {
		a.panicf(e, "cannot specify DISTINT to scalar function CONCAT")
	}

	if len(e.Args) == 0 {
		a.panicf(e, "CONCAT needs one or more arguments")
	}

	var baseType Type

	for _, arg := range e.Args {
		if arg.IntervalUnit != nil {
			a.panicf(arg, "CONCAT does not accept INTERVAL argument")
		}
		t := a.analyzeExpr(arg.Expr)
		if baseType == nil {
			if TypeCoerce(t.Type, StringType) || TypeCoerce(t.Type, BytesType) {
				baseType = t.Type
			} else {
				a.panicf(arg, "CONCAT accepts STRING or BYTES argument, but: %s", TypeString(t.Type))
			}
		} else {
			if !TypeCoerce(t.Type, baseType) {
				a.panicf(arg, "CONCAT accepts %s, but: %s", TypeString(baseType), TypeString(t.Type))
			}
		}
	}

	return &TypeInfo{
		Type: baseType,
	}
}

func sumAnalyzer(a *Analyzer, e *ast.CallExpr) *TypeInfo {
	var context *GroupByContext
	a.scope.context, context = nil, a.scope.context
	defer func() { a.scope.context = context }()

	if len(e.Args) != 1 {
		a.panicf(e, "SUM needs just one argument")
	}

	if e.Args[0].IntervalUnit != nil {
		a.panicf(e.Args[0], "CONCAT does not accept INTERVAL argument")
	}

	t := a.analyzeExpr(e.Args[0].Expr)
	if !(TypeCoerce(t.Type, Int64Type) || TypeCoerce(t.Type, Float64Type)) {
		a.panicf(e.Args[0].Expr, "SUM accepts INT64 or FLOAT64, but: %s", TypeString(t.Type))
	}

	return &TypeInfo{
		Type: t.Type,
	}
}

func (a *Analyzer) analyzeCountStarExpr(e *ast.CountStarExpr) *TypeInfo {
	return &TypeInfo{
		Type: Int64Type,
	}
}
