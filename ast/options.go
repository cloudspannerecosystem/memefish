package ast

import (
	"errors"
	"fmt"
	"strconv"
)

// Only the `Options` method should be placed in this file.

// Field finds name in Records, and return its value as Expr.
// The second return value indicates that the name was found in Records.
// You should use other *Field methods if possible.
func (o *Options) Field(name string) (expr Expr, found bool) {
	for _, r := range o.Records {
		if r.Name.Name != name {
			continue
		}
		return r.Value, true
	}
	return nil, false
}

// FieldNotFound is returned Options.*Field methods.
// It can be handled as a non-error.
var ErrFieldNotFound = errors.New("field not found")

// fieldTypeMismatch is only defined for test.
// It is intentionally unexported.
var errFieldTypeMismatch = errors.New("type mismatched")

// BoolField finds name in records, and return its value as *bool.
// If Options doesn't have a record with name, it returns FieldNotFound error.
// If record have NullLiteral value, it returns nil.
// If record have BoolLiteral value, it returns pointer of bool value.
// If record have value which is neither NullLiteral nor BoolLiteral, it returns error.
func (o *Options) BoolField(name string) (*bool, error) {
	v, ok := o.Field(name)
	if !ok {
		return nil, ErrFieldNotFound
	}
	switch v := v.(type) {
	case *BoolLiteral:
		return &v.Value, nil
	case *NullLiteral:
		return nil, nil
	default:
		return nil, fmt.Errorf("expect true, false or null, but got unknown type %T: %w", v, errFieldTypeMismatch)
	}
}

// IntegerField finds name in records, and return its value as *int64.
// If Options doesn't have a record with name, it returns FieldNotFound error.
// If record have NullLiteral value, it returns nil.
// If record have IntegerLiteral value, it returns pointer of int64 value.
// If record have value which is neither NullLiteral nor BoolLiteral, it returns error.
func (o *Options) IntegerField(name string) (*int64, error) {
	v, ok := o.Field(name)
	if !ok {
		return nil, ErrFieldNotFound
	}
	switch v := v.(type) {
	case *IntLiteral:
		n, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return &n, nil
	case *NullLiteral:
		return nil, nil
	default:
		return nil, fmt.Errorf("expect integer or null, but got unknown type %T: %w", v, errFieldTypeMismatch)
	}
}

// StringField finds name in records, and return its value as *string.
// If Options doesn't have a record with name, it returns FieldNotFound error.
// If record have NullLiteral value, it returns nil.
// If record have StringLiteral value, it returns pointer of string value.
// If record have value which is neither NullLiteral nor StringLiteral, it returns error.
func (o *Options) StringField(name string) (*string, error) {
	v, ok := o.Field(name)
	if !ok {
		return nil, ErrFieldNotFound
	}
	switch v := v.(type) {
	case *StringLiteral:
		return &v.Value, nil
	case *NullLiteral:
		return nil, nil
	default:
		return nil, fmt.Errorf("expect string literal or null, but got unknown type %T: %w", v, errFieldTypeMismatch)
	}
}
