package memefish_test

import (
	"testing"

	"github.com/cloudspannerecosystem/memefish"
)

func TestParseStatementsPipeDropTrailingCommaBeforeSemicolon(t *testing.T) {
	stmts, err := memefish.ParseStatements("", "FROM Singers |> DROP FirstName,;")
	if err != nil {
		t.Fatalf("ParseStatements() error = %v, want nil", err)
	}
	if len(stmts) != 1 {
		t.Fatalf("ParseStatements() statements = %d, want 1", len(stmts))
	}
	if got := stmts[0].SQL(); got != "FROM Singers |> DROP FirstName" {
		t.Errorf("ParseStatements() SQL() = %q, want %q", got, "FROM Singers |> DROP FirstName")
	}
}
