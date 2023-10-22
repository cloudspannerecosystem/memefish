package ast

import (
	"testing"

	"github.com/cloudspannerecosystem/memefish/token"
)

func TestTableAlteration_Position(t *testing.T) {
	alterColumnDefaultExpr := AlterColumn{
		DefaultExpr: &ColumnDefaultExpr{Rparen: 100},
	}
	if alterColumnDefaultExpr.End() != alterColumnDefaultExpr.DefaultExpr.End() {
		t.Fatalf("Mismatched end postion of the alter column")
	}
	alterColumnNull := AlterColumn{
		Null:    101,
		NotNull: true,
	}
	if alterColumnNull.End() != alterColumnNull.Null+4 {
		t.Fatalf("Mismatched end postion of the alter column")
	}
	alterColumnType := AlterColumn{
		Null: token.InvalidPos,
		Type: &ScalarSchemaType{NamePos: 102},
	}
	if alterColumnType.End() != alterColumnType.Type.End() {
		t.Fatalf("Mismatched end postion of the alter column")
	}

	alterColumnSetDefault := AlterColumnSet{
		DefaultExpr: &ColumnDefaultExpr{Rparen: 103},
	}
	if alterColumnSetDefault.End() != alterColumnSetDefault.DefaultExpr.End() {
		t.Fatalf("Mismatched end postion of the alter column set")
	}
	alterColumnSetOptions := AlterColumnSet{
		Options: &ColumnDefOptions{Rparen: 104},
	}
	if alterColumnSetOptions.End() != alterColumnSetOptions.Options.End() {
		t.Fatalf("Mismatched end postion of the alter column set")
	}
}
