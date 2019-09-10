package analyzer

// Type represents SQL types.
type Type interface {
	String() string
	EqualTo(t Type) bool
	CastTo(t Type) bool
	CoerceTo(t Type) bool
}

// SimpleType is types except for ARRAY and STRUCT.
type SimpleType string

const (
	Int64Type     SimpleType = "INT64"
	Float64Type   SimpleType = "FLOAT64"
	BoolType      SimpleType = "BOOL"
	StringType    SimpleType = "STRING"
	BytesType     SimpleType = "BYTES"
	DateType      SimpleType = "DATE"
	TimestampType SimpleType = "TIMESTAMP"
)

// ArrayType is ARRAY type.
type ArrayType struct {
	// A nested array is not supported in Spanner, so Item never become ArrayType.
	Item Type
}

// StructType is STRUCT type.
type StructType struct {
	Fields []*StructField
}

// StructField is STRUCT field.
type StructField struct {
	Name string
	Type Type
}

func (s SimpleType) String() string {
	return string(s)
}

func (a *ArrayType) String() string {
	return "ARRAY<" + TypeString(a.Item) + ">"
}

func (s *StructType) String() string {
	t := "STRUCT<"
	for i, f := range s.Fields {
		if i != 0 {
			t += ", "
		}
		if f.Name != "" {
			t += f.Name + " "
		}
		t += TypeString(f.Type)
	}
	t += ">"
	return t
}

func (s SimpleType) EqualTo(t Type) bool {
	if t, ok := t.(SimpleType); ok {
		return s == t
	} else {
		return false
	}
}

func (a *ArrayType) EqualTo(t Type) bool {
	if t, ok := t.(*ArrayType); ok {
		return TypeEqual(a.Item, t.Item)
	} else {
		return false
	}
}

func (s *StructType) EqualTo(t Type) bool {
	if t, ok := t.(*StructType); ok {
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if !TypeEqual(f.Type, t.Fields[i].Type) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (s SimpleType) CastTo(t Type) bool {
	if t, ok := t.(SimpleType); !ok {
		return false
	} else {
		// The same types can be cast to each others of course.
		if s == t {
			return true
		}

		// See: https://cloud.google.com/spanner/docs/functions-and-operators#casting
		switch s {
		case Int64Type:
			return t == StringType || t == Float64Type || t == BoolType
		case Float64Type:
			return t == StringType || t == Int64Type
		case BoolType:
			return t == StringType || t == Int64Type
		case StringType:
			return true // StringType can cast to any types via parsing.
		case BytesType:
			return t == StringType
		case DateType:
			return t == StringType || t == TimestampType
		case TimestampType:
			return t == StringType || t == DateType
		}
	}

	panic("unreachable")
}

func (a *ArrayType) CastTo(t Type) bool {
	if t, ok := t.(*ArrayType); ok {
		return TypeEqual(a.Item, t.Item)
	}

	return false
}

func (s *StructType) CastTo(t Type) bool {
	if t, ok := t.(*StructType); ok {
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if !TypeEqual(f.Type, t.Fields[i].Type) {
				return false
			}
		}
		return true
	}

	return false
}

func (s SimpleType) CoerceTo(t Type) bool {
	if t, ok := t.(SimpleType); ok {
		if s == t {
			return true
		}
		switch s {
		case Int64Type:
			return t == Float64Type
		case StringType:
			return t == DateType || t == TimestampType
		}
	}

	return false
}

func (a *ArrayType) CoerceTo(t Type) bool {
	if t, ok := t.(*ArrayType); ok {
		return TypeEqual(a.Item, t.Item)
	}
	return false
}

func (s *StructType) CoerceTo(t Type) bool {
	if t, ok := t.(*StructType); ok {
		if len(s.Fields) != len(t.Fields) {
			return false
		}
		for i, f := range s.Fields {
			if !TypeCoerce(f.Type, t.Fields[i].Type) {
				return false
			}
		}
		return true
	}
	return false
}

// TypeEqual checks s equals to t in structual.
func TypeEqual(s, t Type) bool {
	if s == nil || t == nil {
		return true
	}
	return s.EqualTo(t)
}

// TypeCast checks s can cast to t.
func TypeCast(s, t Type) bool {
	if s == nil || t == nil {
		return true
	}
	return s.CastTo(t)
}

// TypeCoerce checks s convert to t implicitly.
func TypeCoerce(s, t Type) bool {
	if s == nil || t == nil {
		return true
	}
	return s.CoerceTo(t)
}

// TypeString returns string representation of t.
func TypeString(t Type) string {
	if t == nil {
		return "(null)"
	}
	return t.String()
}

// TypeCastArray casts t to ArrayType.
func TypeCastArray(t Type) (*ArrayType, bool) {
	if t == nil {
		return nil, false
	}
	tt, ok := t.(*ArrayType)
	return tt, ok
}

// TypeCastStruct casts t to StructType.
func TypeCastStruct(t Type) (*StructType, bool) {
	if t == nil {
		return nil, false
	}
	tt, ok := t.(*StructType)
	return tt, ok
}

// MergeType merges s and t into a type.
func MergeType(s, t Type) (Type, bool) {
	if s == nil {
		return t, true
	}
	if t == nil {
		return s, true
	}

	if TypeEqual(s, t) {
		return s, true
	}

	s1, sok := s.(*StructType)
	t1, tok := t.(*StructType)
	if sok && tok {
		return mergeStructType(s1, t1)
	}

	if TypeCoerce(t, s) {
		return s, true
	}
	if TypeCoerce(s, t) {
		return t, true
	}

	return nil, false
}

func mergeStructType(s, t *StructType) (Type, bool) {
	if len(s.Fields) != len(t.Fields) {
		return nil, false
	}

	fields := make([]*StructField, len(s.Fields))
	for i, f := range s.Fields {
		t, ok := MergeType(f.Type, t.Fields[i].Type)
		if !ok {
			return nil, false
		}
		fields[i] = &StructField{
			Name: f.Name,
			Type: t,
		}
	}
	return &StructType{Fields: fields}, true
}
