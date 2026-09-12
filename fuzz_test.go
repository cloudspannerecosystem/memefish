package memefish_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudspannerecosystem/memefish"
	"github.com/cloudspannerecosystem/memefish/ast"
)

// addSeedsFromTestdata seeds the fuzzing corpus with the *.sql files under
// testdata/input/<dir> for each given dir.
func addSeedsFromTestdata(f *testing.F, dirs ...string) {
	f.Helper()

	for _, dir := range dirs {
		paths, err := filepath.Glob(filepath.Join("testdata", "input", dir, "*.sql"))
		if err != nil {
			f.Fatal(err)
		}
		if len(paths) == 0 {
			f.Fatalf("no seed files found in testdata/input/%s", dir)
		}
		for _, path := range paths {
			b, err := os.ReadFile(path)
			if err != nil {
				f.Fatal(err)
			}
			f.Add(string(b))
		}
	}
}

// checkRoundTrip checks that the SQL rendering of a successfully parsed node
// re-parses without error and renders to the same SQL again (idempotence).
func checkRoundTrip[T ast.Node](t *testing.T, parse func(s string) (T, error), input string, node T) {
	t.Helper()

	sql1 := node.SQL()
	node2, err := parse(sql1)
	if err != nil {
		t.Fatalf("failed to re-parse the SQL rendering of a parsed node\ninput: %q\nrendered SQL: %q\nerror: %v", input, sql1, err)
	}
	sql2 := node2.SQL()
	if sql1 != sql2 {
		t.Fatalf("SQL rendering is not idempotent\ninput: %q\nfirst rendering: %q\nsecond rendering: %q", input, sql1, sql2)
	}
}

func FuzzParseStatements(f *testing.F) {
	addSeedsFromTestdata(f, "statement", "query", "ddl", "dml", "gql")

	f.Fuzz(func(t *testing.T, s string) {
		stmts, err := memefish.ParseStatements("", s)
		if err != nil {
			return
		}
		for _, stmt := range stmts {
			checkRoundTrip(t, func(s string) (ast.Statement, error) { return memefish.ParseStatement("", s) }, s, stmt)
		}
	})
}

func FuzzParseQuery(f *testing.F) {
	addSeedsFromTestdata(f, "query", "gql")

	f.Fuzz(func(t *testing.T, s string) {
		stmt, err := memefish.ParseQuery("", s)
		if err != nil {
			return
		}
		checkRoundTrip(t, func(s string) (*ast.QueryStatement, error) { return memefish.ParseQuery("", s) }, s, stmt)
	})
}

func FuzzParseExpr(f *testing.F) {
	addSeedsFromTestdata(f, "expr")

	f.Fuzz(func(t *testing.T, s string) {
		expr, err := memefish.ParseExpr("", s)
		if err != nil {
			return
		}
		checkRoundTrip(t, func(s string) (ast.Expr, error) { return memefish.ParseExpr("", s) }, s, expr)
	})
}

func FuzzParseType(f *testing.F) {
	for _, seed := range []string{
		"INT64",
		"FLOAT64",
		"BOOL",
		"STRING",
		"BYTES",
		"TIMESTAMP",
		"ARRAY<INT64>",
		"ARRAY<ARRAY<STRING>>",
		"STRUCT<>",
		"STRUCT<x INT64, y FLOAT64>",
		"STRUCT<arr ARRAY<STRUCT<n INT64>>>",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		typ, err := memefish.ParseType("", s)
		if err != nil {
			return
		}
		checkRoundTrip(t, func(s string) (ast.Type, error) { return memefish.ParseType("", s) }, s, typ)
	})
}

func FuzzParseSchemaType(f *testing.F) {
	for _, seed := range []string{
		"INT64",
		"BOOL",
		"STRING(MAX)",
		"STRING(42)",
		"BYTES(1024)",
		"NUMERIC",
		"TIMESTAMP",
		"ARRAY<STRING(MAX)>",
		"ARRAY<BYTES(MAX)>",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		typ, err := memefish.ParseSchemaType("", s)
		if err != nil {
			return
		}
		checkRoundTrip(t, func(s string) (ast.SchemaType, error) { return memefish.ParseSchemaType("", s) }, s, typ)
	})
}

func FuzzParseGQLGraphPattern(f *testing.F) {
	addSeedsFromTestdata(f, "gql_graph_pattern")

	f.Fuzz(func(t *testing.T, s string) {
		pattern, err := memefish.ParseGQLGraphPattern("", s)
		if err != nil {
			return
		}
		checkRoundTrip(t, func(s string) (ast.GQLGraphPatternNode, error) { return memefish.ParseGQLGraphPattern("", s) }, s, pattern)
	})
}

func FuzzSplitRawStatements(f *testing.F) {
	addSeedsFromTestdata(f, "statement", "query", "ddl", "dml", "gql")

	f.Fuzz(func(t *testing.T, s string) {
		stmts, err := memefish.SplitRawStatements("", s)
		if err != nil {
			return
		}
		for i, raw := range stmts {
			// The empty fallback result has zero positions and an empty statement.
			if raw.Pos < 0 || raw.End < raw.Pos || len(s) < int(raw.End) {
				t.Fatalf("statement %d has out-of-range positions [%d, %d) for input of length %d\ninput: %q", i, raw.Pos, raw.End, len(s), s)
			}
			if len(stmts) == 1 && raw.Statement == "" && raw.Pos == 0 && raw.End == 0 {
				continue
			}
			if raw.Statement != s[raw.Pos:raw.End] {
				t.Fatalf("statement %d text %q does not match input range [%d, %d) %q\ninput: %q", i, raw.Statement, raw.Pos, raw.End, s[raw.Pos:raw.End], s)
			}
		}
	})
}
