package analyzer

import (
	"fmt"
)

func valueType(v interface{}) (Type, error) {
	switch v := v.(type) {
	case nil:
		return nil, nil
	case bool:
		return BoolType, nil
	case int64:
		return Int64Type, nil
	case float64:
		return Float64Type, nil
	case string:
		return StringType, nil
	case []interface{}:
		item, err := valueType(v[0])
		if err != nil {
			return nil, err
		}
		return &ArrayType{Item: item}, nil
	case *Placeholder:
		// Early resolve when placeholder is filled.
		// It is needed for working such like `select @param.* FROM (select 1) WHERE cast(@param as struct<foo string>).foo = "foo"` SQL.
		// In this SQL, `@param` is filled as `struct` type on `WHER` phase at first, and after,
		// `@param` in `SELECT` is used as `struct`, however `@param` type is `PlaceholderType` in fact, so early resolve is needed.
		if v.PlaceholderType.Type != nil {
			return v.PlaceholderType.Type, nil
		}
		return v.PlaceholderType, nil
	case []map[string]interface{}:
		var fields []*StructField
		for _, kv := range v {
			for name, v := range kv {
				t, err := valueType(v)
				if err != nil {
					return nil, err
				}
				fields = append(fields, &StructField{
					Name: name,
					Type: t,
				})
			}
		}
		return &StructType{Fields: fields}, nil
	default:
		// TODO: support BYTES, DATE and TIMESTAMP type values.
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
