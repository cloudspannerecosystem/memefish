package ast

import (
	"testing"

	"github.com/cloudspannerecosystem/memefish/token"
)

func TestTableAlteration_Position(t *testing.T) {
	alterColumnDefaultExpr := AlterColumn{
		Alteration: &AlterColumnType{DefaultExpr: &ColumnDefaultExpr{Rparen: 100}},
	}
	if alterColumnDefaultExpr.End() != alterColumnDefaultExpr.Alteration.(*AlterColumnType).DefaultExpr.End() {
		t.Fatalf("Mismatched end postion of the alter column")
	}
	alterColumnNull := AlterColumnType{
		Null:    101,
		NotNull: true,
	}
	if alterColumnNull.End() != alterColumnNull.Null+4 {
		t.Fatalf("Mismatched end postion of the alter column")
	}
	alterColumnType := AlterColumn{Alteration: &AlterColumnType{
		Null: token.InvalidPos,
		Type: &ScalarSchemaType{NamePos: 102},
	}}
	if alterColumnType.End() != alterColumnType.Alteration.(*AlterColumnType).Type.End() {
		t.Fatalf("Mismatched end postion of the alter column")
	}

	alterColumnSetDefault := AlterColumn{Alteration: &AlterColumnSetDefault{
		DefaultExpr: &ColumnDefaultExpr{Rparen: 103},
	}}
	if alterColumnSetDefault.End() != alterColumnSetDefault.Alteration.(*AlterColumnSetDefault).DefaultExpr.End() {
		t.Fatalf("Mismatched end postion of the alter column set")
	}
	alterColumnSetOptions := AlterColumn{Alteration: &AlterColumnSetOptions{
		Options: &Options{Rparen: 104},
	}}
	if alterColumnSetOptions.End() != alterColumnSetOptions.Alteration.(*AlterColumnSetOptions).Options.End() {
		t.Fatalf("Mismatched end postion of the alter column set")
	}
}
