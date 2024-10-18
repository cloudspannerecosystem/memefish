package ast

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestStatement(t *testing.T) {
	Statement(&QueryStatement{}).isStatement()
	Statement(&CreateDatabase{}).isStatement()
	Statement(&CreateTable{}).isStatement()
	Statement(&CreateView{}).isStatement()
	Statement(&CreateIndex{}).isStatement()
	Statement(&CreateSequence{}).isStatement()
	Statement(&CreateRole{}).isStatement()
	Statement(&AlterTable{}).isStatement()
	Statement(&AlterIndex{}).isStatement()
	Statement(&DropTable{}).isStatement()
	Statement(&DropIndex{}).isStatement()
	Statement(&DropVectorIndex{}).isStatement()
	Statement(&DropRole{}).isStatement()
	Statement(&Insert{}).isStatement()
	Statement(&Delete{}).isStatement()
	Statement(&Update{}).isStatement()
	Statement(&Grant{}).isStatement()
	Statement(&Revoke{}).isStatement()
}

func TestQueryExpr(t *testing.T) {
	QueryExpr(&Select{}).isQueryExpr()
	QueryExpr(&SubQuery{}).isQueryExpr()
	QueryExpr(&CompoundQuery{}).isQueryExpr()
}

func TestSelectItem(t *testing.T) {
	SelectItem(&Star{}).isSelectItem()
	SelectItem(&DotStar{}).isSelectItem()
	SelectItem(&Alias{}).isSelectItem()
	SelectItem(&ExprSelectItem{}).isSelectItem()
}

func TestTableExpr(t *testing.T) {
	TableExpr(&Unnest{}).isTableExpr()
	TableExpr(&TableName{}).isTableExpr()
	TableExpr(&SubQueryTableExpr{}).isTableExpr()
	TableExpr(&ParenTableExpr{}).isTableExpr()
	TableExpr(&Join{}).isTableExpr()
}

func TestJoinCondition(t *testing.T) {
	JoinCondition(&On{}).isJoinCondition()
	JoinCondition(&Using{}).isJoinCondition()
}

func TestExpr(t *testing.T) {
	Expr(&BinaryExpr{}).isExpr()
	Expr(&UnaryExpr{}).isExpr()
	Expr(&InExpr{}).isExpr()
	Expr(&IsNullExpr{}).isExpr()
	Expr(&IsBoolExpr{}).isExpr()
	Expr(&BetweenExpr{}).isExpr()
	Expr(&SelectorExpr{}).isExpr()
	Expr(&IndexExpr{}).isExpr()
	Expr(&CallExpr{}).isExpr()
	Expr(&CountStarExpr{}).isExpr()
	Expr(&CastExpr{}).isExpr()
	Expr(&ExtractExpr{}).isExpr()
	Expr(&CaseExpr{}).isExpr()
	Expr(&ParenExpr{}).isExpr()
	Expr(&ScalarSubQuery{}).isExpr()
	Expr(&ArraySubQuery{}).isExpr()
	Expr(&ExistsSubQuery{}).isExpr()
	Expr(&Param{}).isExpr()
	Expr(&Ident{}).isExpr()
	Expr(&Path{}).isExpr()
	Expr(&ArrayLiteral{}).isExpr()
	Expr(&StructLiteral{}).isExpr()
	Expr(&NullLiteral{}).isExpr()
	Expr(&BoolLiteral{}).isExpr()
	Expr(&IntLiteral{}).isExpr()
	Expr(&FloatLiteral{}).isExpr()
	Expr(&StringLiteral{}).isExpr()
	Expr(&BytesLiteral{}).isExpr()
	Expr(&DateLiteral{}).isExpr()
	Expr(&TimestampLiteral{}).isExpr()
	Expr(&NumericLiteral{}).isExpr()
}

func TestArg(t *testing.T) {
	Arg(&IntervalArg{}).isArg()
	Arg(&ExprArg{}).isArg()
	Arg(&SequenceArg{}).isArg()
}

func TestInCondition(t *testing.T) {
	InCondition(&UnnestInCondition{}).isInCondition()
	InCondition(&SubQueryInCondition{}).isInCondition()
	InCondition(&ValuesInCondition{}).isInCondition()
}

func TestType(t *testing.T) {
	Type(&SimpleType{}).isType()
	Type(&ArrayType{}).isType()
	Type(&StructType{}).isType()
}

func TestIntValue(t *testing.T) {
	IntValue(&Param{}).isIntValue()
	IntValue(&IntLiteral{}).isIntValue()
	IntValue(&CastIntValue{}).isIntValue()
}

func TestNumValue(t *testing.T) {
	NumValue(&Param{}).isNumValue()
	NumValue(&IntLiteral{}).isNumValue()
	NumValue(&FloatLiteral{}).isNumValue()
	NumValue(&CastNumValue{}).isNumValue()
}

func TestStringValue(t *testing.T) {
	StringValue(&Param{}).isStringValue()
	StringValue(&StringLiteral{}).isStringValue()
}

func TestDDL(t *testing.T) {
	DDL(&CreateDatabase{}).isDDL()
	DDL(&CreateTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&CreateVectorIndex{}).isDDL()
	DDL(&CreateSequence{}).isDDL()
	DDL(&AlterSequence{}).isDDL()
	DDL(&DropSequence{}).isDDL()
	DDL(&CreateView{}).isDDL()
	DDL(&AlterTable{}).isDDL()
	DDL(&DropTable{}).isDDL()
	DDL(&CreateIndex{}).isDDL()
	DDL(&AlterIndex{}).isDDL()
	DDL(&DropIndex{}).isDDL()
	DDL(&DropVectorIndex{}).isDDL()
	DDL(&CreateRole{}).isDDL()
	DDL(&DropRole{}).isDDL()
	DDL(&Grant{}).isDDL()
	DDL(&Revoke{}).isDDL()
}

func TestConstraint(t *testing.T) {
	Constraint(&ForeignKey{}).isConstraint()
	Constraint(&Check{}).isConstraint()
}

func TestTableAlteration(t *testing.T) {
	TableAlteration(&AddColumn{}).isTableAlteration()
	TableAlteration(&AddTableConstraint{}).isTableAlteration()
	TableAlteration(&DropColumn{}).isTableAlteration()
	TableAlteration(&DropConstraint{}).isTableAlteration()
	TableAlteration(&SetOnDelete{}).isTableAlteration()
	TableAlteration(&AlterColumn{}).isTableAlteration()
	TableAlteration(&AlterColumnSet{}).isTableAlteration()
}

func TestPrivilege(t *testing.T) {
	Privilege(&PrivilegeOnTable{}).isPrivilege()
	Privilege(&SelectPrivilegeOnView{}).isPrivilege()
	Privilege(&ExecutePrivilegeOnTableFunction{}).isPrivilege()
	Privilege(&RolePrivilege{}).isPrivilege()
}

func TestTablePrivilege(t *testing.T) {
	TablePrivilege(&SelectPrivilege{}).isTablePrivilege()
	TablePrivilege(&InsertPrivilege{}).isTablePrivilege()
	TablePrivilege(&UpdatePrivilege{}).isTablePrivilege()
	TablePrivilege(&DeletePrivilege{}).isTablePrivilege()
}

func TestSchemaType(t *testing.T) {
	SchemaType(&ScalarSchemaType{}).isSchemaType()
	SchemaType(&SizedSchemaType{}).isSchemaType()
	SchemaType(&ArraySchemaType{}).isSchemaType()
}

func TestIndexAlteration(t *testing.T) {
	IndexAlteration(&AddStoredColumn{}).isIndexAlteration()
	IndexAlteration(&DropStoredColumn{}).isIndexAlteration()
}

func TestDML(t *testing.T) {
	DML(&Insert{}).isDML()
	DML(&Delete{}).isDML()
	DML(&Update{}).isDML()
}

func TestInsertInput(t *testing.T) {
	InsertInput(&ValuesInput{}).isInsertInput()
	InsertInput(&SubQueryInput{}).isInsertInput()
}

func ptr[T any](v T) *T {
	return &v
}

func TestOptions_BoolField(t *testing.T) {
	tcases := []struct {
		desc    string
		input   *Options
		name    string
		want    *bool
		wantErr error
	}{
		{
			desc: "true",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name: "bool_option",
			want: ptr(true),
		},
		{
			desc: "false",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: false}},
				},
			},
			name: "bool_option",
			want: ptr(false),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "bool_option"}, &NullLiteral{}},
				},
			},
			name: "bool_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "dummy"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			want:    nil,
			wantErr: FieldNotFound,
		},
		{
			desc: "invalid type",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "string_option",
			wantErr: fieldTypeMismatch,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.desc, func(t *testing.T) {
			got, err := tcase.input.BoolField(tcase.name)
			if tcase.wantErr == nil && err != nil {
				t.Errorf("should not fail, but: %v", err)
			}
			if tcase.wantErr != nil && err == nil {
				t.Errorf("should fail, but success")
			}
			if tcase.wantErr != nil && !errors.Is(err, tcase.wantErr) {
				t.Errorf("error differ, want: %v, got: %v", tcase.wantErr, err)
			}
			if diff := cmp.Diff(tcase.want, got); diff != "" {
				t.Errorf("differ: %v", diff)
			}
		})
	}
}

func TestOptions_StringField(t *testing.T) {
	tcases := []struct {
		desc    string
		input   *Options
		name    string
		want    *string
		wantErr error
	}{
		{
			desc: "string",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name: "string_option",
			want: ptr("foo"),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "string_option"}, &NullLiteral{}},
				},
			},
			name: "string_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "dummy_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "string_field",
			want:    nil,
			wantErr: FieldNotFound,
		},
		{
			desc: "invalid value",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			wantErr: fieldTypeMismatch,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.desc, func(t *testing.T) {
			got, err := tcase.input.StringField(tcase.name)
			if tcase.wantErr == nil && err != nil {
				t.Errorf("should not fail, but: %v", err)
			}
			if tcase.wantErr != nil && err == nil {
				t.Errorf("should fail, but success")
			}
			if tcase.wantErr != nil && !errors.Is(err, tcase.wantErr) {
				t.Errorf("error differ, want: %v, got: %v", tcase.wantErr, err)
			}
			if diff := cmp.Diff(tcase.want, got); diff != "" {
				t.Errorf("differ: %v", diff)
			}
		})
	}
}

func TestOptions_IntegerField(t *testing.T) {
	tcases := []struct {
		desc    string
		input   *Options
		name    string
		want    *int64
		wantErr error
	}{
		{
			desc: "integer",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "integer_option"}, &IntLiteral{Value: "7"}},
				},
			},
			name: "integer_option",
			want: ptr(int64(7)),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "integer_option"}, &NullLiteral{}},
				},
			},
			name: "integer_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "integer_option",
			want:    nil,
			wantErr: FieldNotFound,
		},
		{
			desc: "invalid value",
			input: &Options{
				Records: []*OptionsRecord{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			wantErr: fieldTypeMismatch,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.desc, func(t *testing.T) {
			got, err := tcase.input.IntegerField(tcase.name)
			if tcase.wantErr == nil && err != nil {
				t.Errorf("should not fail, but: %v", err)
			}
			if tcase.wantErr != nil && err == nil {
				t.Errorf("should fail, but success")
			}
			if tcase.wantErr != nil && !errors.Is(err, tcase.wantErr) {
				t.Errorf("error differ, want: %v, got: %v", tcase.wantErr, err)
			}
			if diff := cmp.Diff(tcase.want, got); diff != "" {
				t.Errorf("differ: %v", diff)
			}
		})
	}

}
