package analyzer

import (
	"strings"

	"github.com/MakeNowJust/memefish/pkg/parser"
)

type functionAnalyzer func(a *Analyzer, e *parser.CallExpr) *TypeInfo

var functionAnalyzers map[string]functionAnalyzer

func init() {
	functionAnalyzers = map[string]functionAnalyzer{
		"CONCAT": concatAnalyzer,
		"SUM":    sumAnalyzer,
	}
}

func (a *Analyzer) analyzeCallExpr(e *parser.CallExpr) *TypeInfo {
	if ana, ok := functionAnalyzers[strings.ToUpper(e.Func.Name)]; ok {
		return ana(a, e)
	}

	panic(a.errorf(e, "unknown function: %s", e.Func.SQL()))
}

func concatAnalyzer(a *Analyzer, e *parser.CallExpr) *TypeInfo {
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

func sumAnalyzer(a *Analyzer, e *parser.CallExpr) *TypeInfo {
	if a.aggregateScope == nil {
		a.panicf(e, "cannot call SUM without aggregate context")
	}

	var oldScope *NameScope
	oldScope, a.scope, a.aggregateScope = a.scope, a.aggregateScope, nil
	defer func() { a.scope, a.aggregateScope = oldScope, a.scope }()

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

func (a *Analyzer) analyzeCountStarExpr(e *parser.CountStarExpr) *TypeInfo {
	if a.aggregateScope == nil {
		a.panicf(e, "cannot call COUNT(*) without aggregate context")
	}
	return &TypeInfo{
		Type: Int64Type,
	}
}
