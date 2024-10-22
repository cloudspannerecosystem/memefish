package ast

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				Records: []*OptionsDef{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name: "bool_option",
			want: ptr(true),
		},
		{
			desc: "false",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: false}},
				},
			},
			name: "bool_option",
			want: ptr(false),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "bool_option"}, &NullLiteral{}},
				},
			},
			name: "bool_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "dummy"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			want:    nil,
			wantErr: ErrFieldNotFound,
		},
		{
			desc: "invalid type",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "string_option",
			wantErr: errFieldTypeMismatch,
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
				Records: []*OptionsDef{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name: "string_option",
			want: ptr("foo"),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "string_option"}, &NullLiteral{}},
				},
			},
			name: "string_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "dummy_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "string_field",
			want:    nil,
			wantErr: ErrFieldNotFound,
		},
		{
			desc: "invalid value",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			wantErr: errFieldTypeMismatch,
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
				Records: []*OptionsDef{
					{&Ident{Name: "integer_option"}, &IntLiteral{Value: "7"}},
				},
			},
			name: "integer_option",
			want: ptr(int64(7)),
		},
		{
			desc: "explicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "integer_option"}, &NullLiteral{}},
				},
			},
			name: "integer_option",
			want: nil,
		},
		{
			desc: "implicit null",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "string_option"}, &StringLiteral{Value: "foo"}},
				},
			},
			name:    "integer_option",
			want:    nil,
			wantErr: ErrFieldNotFound,
		},
		{
			desc: "invalid value",
			input: &Options{
				Records: []*OptionsDef{
					{&Ident{Name: "bool_option"}, &BoolLiteral{Value: true}},
				},
			},
			name:    "bool_option",
			wantErr: errFieldTypeMismatch,
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

func ptr[T any](v T) *T {
	return &v
}
